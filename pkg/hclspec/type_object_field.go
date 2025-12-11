// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// FieldConstraint represents a constraint on the values of a specific field within an object.
type FieldConstraint interface {
	// Validate checks if the given value satisfies the constraint.
	//
	// Implementations are expected to return all errors produced by the validation, if possible.
	Validate(value cty.Value) error

	// ValidateSpec checks if the constraint is valid for the object field.
	ValidateSpec(field *objectField) error

	// ToProtobuf converts the FieldConstraint to its protobuf representation.
	ToProtobuf() (*FieldConstraintPB, error)
}

// ToFieldConstraint converts a protobuf FieldConstraintPB to a FieldConstraint instance.
func (f *FieldConstraintPB) ToFieldConstraint() (FieldConstraint, error) {
	if f == nil {
		return nil, fmt.Errorf("FieldConstraintPB is nil")
	}

	switch constraint := f.Constraint.(type) {
	case *FieldConstraintPB_AllowedValues:
		return constraint.AllowedValues.ToFieldConstraint()
	default:
		return nil, fmt.Errorf("unknown FieldConstraintPB type: %T", constraint)
	}
}

// FieldConstraints is a slice of FieldConstraint.
type FieldConstraints []FieldConstraint

// Validate checks if the given value satisfies all constraints.
func (c FieldConstraints) Validate(value cty.Value) error {
	var err error
	for _, constraint := range c {
		if e := constraint.Validate(value); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// ValidateSpec checks if all constraints are valid for the object field.
func (c FieldConstraints) ValidateSpec(field *objectField) error {
	if field == nil {
		return fmt.Errorf("field is nil")
	}

	var err error
	for _, constraint := range c {
		if e := constraint.ValidateSpec(field); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// ToProtobuf converts all FieldConstraints to their protobuf representations.
func (c FieldConstraints) ToProtobuf() ([]*FieldConstraintPB, error) {
	pbs := make([]*FieldConstraintPB, 0, len(c))
	for _, constraint := range c {
		pb, err := constraint.ToProtobuf()
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)
	}

	return pbs, nil
}

type allowedValuesConstraint struct {
	values []cty.Value
}

// Validate implements FieldConstraint.
func (a *allowedValuesConstraint) Validate(value cty.Value) error {
	if a == nil {
		return fmt.Errorf("allowed values constraint is nil")
	}

	if value.IsNull() {
		return nil // Skip null values
	}

	for _, v := range a.values {
		if value.Equals(v).True() {
			return nil
		}
	}

	return fmt.Errorf(
		"value %v is not in allowed values: %s",
		util.FormatCtyValueToString(value),
		a.formatAllowedValues(),
	)
}

// ValidateSpec implements FieldConstraint.
func (a *allowedValuesConstraint) ValidateSpec(field *objectField) error {
	if a == nil {
		return fmt.Errorf("allowed values constraint is nil")
	}

	if len(a.values) == 0 {
		return fmt.Errorf("allowed values constraint has no values defined")
	}

	if field == nil {
		return fmt.Errorf("field is nil")
	}

	if field.t == nil || field.t.CtyType().Equals(cty.NilType) {
		return nil // Return early if field type is invalid
	}

	var err error
	for _, v := range a.values {
		if !v.Type().Equals(field.t.CtyType()) {
			err = errors.Join(
				err,
				fmt.Errorf(
					"allowed value %v does not match field type %s",
					v.GoString(),
					field.t.CtyType().FriendlyName(),
				),
			)
			continue
		}

		if e := field.validate(v); e != nil {
			err = errors.Join(err, fmt.Errorf("allowed value %v is invalid: %w", v.GoString(), e))
		}
	}

	return err
}

// ToProtobuf implements FieldConstraint.
func (a *allowedValuesConstraint) ToProtobuf() (*FieldConstraintPB, error) {
	if a == nil {
		return nil, fmt.Errorf("allowed values constraint is nil")
	}

	valuesPB := make([][]byte, 0, len(a.values))
	for _, v := range a.values {
		vPB, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value %v to protobuf: %w", v.GoString(), err)
		}

		valuesPB = append(valuesPB, vPB)
	}

	return &FieldConstraintPB{
		Constraint: &FieldConstraintPB_AllowedValues{
			AllowedValues: &AllowedValuesConstraintPB{
				Values: valuesPB,
			},
		},
	}, nil
}

func (a *allowedValuesConstraint) formatAllowedValues() string {
	if a == nil || len(a.values) == 0 {
		return ""
	}

	allowedValues := make([]string, 0, len(a.values))
	for _, v := range a.values {
		allowedValues = append(allowedValues, util.FormatCtyValueToString(v))
	}

	return strings.Join(allowedValues, ", ")
}

// AllowedValues creates a constraint that only allows the specified values.
//
// This constraint allows null values.  If a null value is not allowed, the field should be marked as required.
func AllowedValues(values ...cty.Value) FieldConstraint {
	return &allowedValuesConstraint{
		values: values,
	}
}

// ToFieldConstraint converts a protobuf AllowedValuesConstraintPB to a FieldConstraint instance.
func (a *AllowedValuesConstraintPB) ToFieldConstraint() (FieldConstraint, error) {
	if a == nil {
		return nil, fmt.Errorf("AllowedValuesConstraintPB is nil")
	}

	values := make([]cty.Value, 0, len(a.Values))
	for _, vPB := range a.Values {
		v, err := json.Unmarshal(vPB, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value from protobuf: %w", err)
		}

		values = append(values, v)
	}

	return AllowedValues(values...), nil
}

type objectField struct {
	name         string
	t            Type
	aliases      []string
	required     bool
	defaultValue cty.Value
	constraints  FieldConstraints
}

// WithAliases sets the alias list to the specified ones.
//
// This function will overwrite any previous aliases.
func (f *objectField) WithAliases(aliases ...string) *objectField {
	f.aliases = aliases
	return f
}

// WithDefaultValue sets the default value to the specified one.
//
// This function will overwrite any previous default value.
func (f *objectField) WithDefaultValue(value cty.Value) *objectField {
	f.defaultValue = value
	return f
}

// WithConstraints sets the constraints on the field to the specified ones.
//
// This function will overwrite any previous constraints.
func (f *objectField) WithConstraints(constraints ...FieldConstraint) *objectField {
	f.constraints = constraints
	return f
}

// ToProtobuf converts the objectField to its protobuf representation.
func (f *objectField) ToProtobuf() (string, *ObjectFieldPB, error) {
	if f == nil {
		return "", nil, fmt.Errorf("object field is nil")
	}

	typePB, err := f.t.ToProtobuf()
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert field %q type to protobuf: %w", f.name, err)
	}

	constraintsPB, err := f.constraints.ToProtobuf()
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert field %q constraints to protobuf: %w", f.name, err)
	}

	defaultValuePB, err := json.Marshal(f.defaultValue, cty.DynamicPseudoType)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert field %q default value to protobuf: %w", f.name, err)
	}

	return f.name, &ObjectFieldPB{
		Type:         typePB,
		Aliases:      f.aliases,
		Required:     f.required,
		DefaultValue: defaultValuePB,
		Constraints:  constraintsPB,
	}, nil
}

func (f *objectField) validate(value cty.Value) error {
	if f == nil {
		return fmt.Errorf("field is nil")
	}

	if f.t == nil {
		return fmt.Errorf("field %q has no type defined", f.name)
	}

	if value.IsNull() {
		if f.required {
			return fmt.Errorf("missing required field %q", f.name)
		}

		return nil
	}

	var err error
	if e := f.t.Validate(value); e != nil {
		err = errors.Join(err, fmt.Errorf("field %q validation failed: %w", f.name, e))
	}

	if e := f.constraints.Validate(value); e != nil {
		err = errors.Join(err, fmt.Errorf("field %q validation failed: %w", f.name, e))
	}

	return err
}

func (f *objectField) validateSpec() error {
	if f == nil {
		return fmt.Errorf("field is nil")
	}

	var err error
	if f.t == nil {
		err = fmt.Errorf("field %q has no type defined", f.name)
	}

	if len(f.aliases) > 0 {
		if slices.Contains(f.aliases, "") {
			err = errors.Join(err, fmt.Errorf("field %q has an empty alias", f.name))
		}
	}

	if f.t == nil {
		return err // Return early if type is not defined, as all further validation requires a type
	}

	e := f.constraints.ValidateSpec(f)
	err = errors.Join(err, e)

	if !f.defaultValue.IsWhollyKnown() {
		err = errors.Join(err, fmt.Errorf("field %q has an unknown default value", f.name))
		return err
	}

	if f.required && !f.defaultValue.IsNull() {
		err = errors.Join(err, fmt.Errorf("field %q is required and has a default value", f.name))
	}

	e = f.t.ValidateSpec()
	if e != nil {
		err = errors.Join(err, e)
		return err // Return early if type is not valid, because it is needed for further validation
	}

	if !f.defaultValue.Type().Equals(f.t.CtyType()) {
		err = errors.Join(
			err,
			fmt.Errorf(
				"field %q default value type mismatch: expected %q, got %q",
				f.name,
				f.t.CtyType().FriendlyName(),
				f.defaultValue.Type().FriendlyName(),
			),
		)
	}

	if e = f.t.Validate(f.defaultValue); e != nil {
		for _, unwrapped := range util.UnwrapErrors(e) {
			err = errors.Join(err, fmt.Errorf("field %q default value validation failed: %w", f.name, unwrapped))
		}
	}

	return err
}

// ToObjectField converts a protobuf ObjectFieldPB to an objectField instance.
func (f *ObjectFieldPB) ToObjectField(name string) (*objectField, error) {
	if f == nil {
		return nil, fmt.Errorf("ObjectFieldPB is nil")
	}

	typeObj, err := f.Type.ToType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert field %q type from protobuf: %w", name, err)
	}

	defaultValue, err := json.Unmarshal(f.DefaultValue, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field %q default value from protobuf: %w", name, err)
	}

	constraints := make(FieldConstraints, 0, len(f.Constraints))
	for _, cPB := range f.Constraints {
		constraint, err := cPB.ToFieldConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q constraint from protobuf: %w", name, err)
		}

		constraints = append(constraints, constraint)
	}

	return &objectField{
		name:         name,
		t:            typeObj,
		aliases:      f.Aliases,
		required:     f.Required,
		defaultValue: defaultValue,
		constraints:  constraints,
	}, nil
}

// RequiredField creates a required object field with the specified name and type.
//
// The default value is set to a null value of the specified type.
// It is initialized with no aliases or constraints.
func RequiredField(name string, t Type) *objectField {
	return &objectField{
		name:         name,
		t:            t,
		required:     true,
		defaultValue: cty.NullVal(t.CtyType()),
	}
}

// OptionalField returns a new optional object field with the specified name and type.
//
// The default value is set to a null value of the specified type.
// It is initialized with no aliases or constraints.
func OptionalField(name string, t Type) *objectField {
	return &objectField{
		name:         name,
		t:            t,
		required:     false,
		defaultValue: cty.NullVal(t.CtyType()),
	}
}
