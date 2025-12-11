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

// ObjectConstraint represents a constraint on an object type.
type ObjectConstraint interface {
	// Validate checks if the given value satisfies the constraint.
	//
	// Implementations are expected to return all errors produced by the validation, if possible.
	Validate(values map[string]cty.Value) error

	// ValidateSpec checks if the constraint is valid for the object type.
	ValidateSpec(t *objectType) error

	// ToProtobuf converts the ObjectConstraint to its protobuf representation.
	ToProtobuf() (*ObjectConstraintPB, error)
}

// ToObjectConstraint converts the protobuf ObjectConstraintPB to an ObjectConstraint.
func (o *ObjectConstraintPB) ToObjectConstraint() (ObjectConstraint, error) {
	if o == nil {
		return nil, errors.New("ObjectConstraintPB is nil")
	}

	switch constraint := o.Constraint.(type) {
	case *ObjectConstraintPB_MutuallyExclusive:
		return constraint.MutuallyExclusive.ToMutuallyExclusiveConstraint()
	case *ObjectConstraintPB_RequiredTogether:
		return constraint.RequiredTogether.ToRequiredTogetherConstraint()
	case *ObjectConstraintPB_RequiredOneOf:
		return constraint.RequiredOneOf.ToRequiredOneOfConstraint()
	case *ObjectConstraintPB_AllowedFieldValues:
		return constraint.AllowedFieldValues.ToAllowedFieldValuesConstraint()
	default:
		return nil, errors.New("unknown ObjectConstraintPB type")
	}
}

type ObjectConstraints []ObjectConstraint

// Validate checks if all constraints are satisfied by the given value.
//
// This function may return more than one error joined together.
func (c ObjectConstraints) Validate(values map[string]cty.Value) error {
	var err error
	for _, constraint := range c {
		if constraint == nil {
			continue // Skip nil constraints
		}

		if e := constraint.Validate(values); e != nil {
			for _, unwrapped := range util.UnwrapErrors(e) {
				err = errors.Join(err, fmt.Errorf("validation failed: %w", unwrapped))
			}
		}
	}

	return err
}

// ValidateSpec checks if all constraints are valid for the object type.
//
// This function may return more than one error joined together.
func (c ObjectConstraints) ValidateSpec(t *objectType) error {
	var err error
	for _, constraint := range c {
		if constraint == nil {
			err = errors.Join(err, errors.New("object type has a nil constraint"))
			continue
		}

		if e := constraint.ValidateSpec(t); e != nil {
			for _, unwrapped := range util.UnwrapErrors(e) {
				err = errors.Join(err, fmt.Errorf("constraint validation failed: %w", unwrapped))
			}
		}
	}

	return err
}

// ToProtobuf converts all ObjectConstraints to their protobuf representations.
func (c ObjectConstraints) ToProtobuf() ([]*ObjectConstraintPB, error) {
	var constraintsPB []*ObjectConstraintPB
	for _, constraint := range c {
		if constraint == nil {
			return nil, errors.New("object type has a nil constraint")
		}

		constraintPB, err := constraint.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint to protobuf: %w", err)
		}

		constraintsPB = append(constraintsPB, constraintPB)
	}

	return constraintsPB, nil
}

type mutuallyExclusiveConstraint struct {
	fields []string // List of field names that are mutually exclusive.
}

// Validate implements ObjectConstraint.
func (m *mutuallyExclusiveConstraint) Validate(values map[string]cty.Value) error {
	if m == nil {
		return errors.New("mutually exclusive constraint is nil")
	}

	var foundFields []string
	for _, field := range m.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}

			foundFields = append(foundFields, field)
		}
	}

	fieldNames := ""
	for i, field := range m.fields {
		if i > 0 {
			fieldNames += ", "
		}

		fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
	}

	if len(foundFields) > 1 {
		return fmt.Errorf("mutually exclusive fields %s are all present", fieldNames)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (m *mutuallyExclusiveConstraint) ValidateSpec(t *objectType) error {
	if m == nil {
		return errors.New("mutually exclusive constraint is nil")
	}

	var err error
	for _, field := range m.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

// ToProtobuf implements ObjectConstraint.
func (m *mutuallyExclusiveConstraint) ToProtobuf() (*ObjectConstraintPB, error) {
	if m == nil {
		return nil, errors.New("mutually exclusive constraint is nil")
	}

	return &ObjectConstraintPB{
		Constraint: &ObjectConstraintPB_MutuallyExclusive{
			MutuallyExclusive: &MutuallyExclusiveConstraintPB{
				Fields: m.fields,
			},
		},
	}, nil
}

// ToMutuallyExclusiveConstraint converts the protobuf MutuallyExclusiveConstraintPB to a mutuallyExclusiveConstraint.
func (m *MutuallyExclusiveConstraintPB) ToMutuallyExclusiveConstraint() (ObjectConstraint, error) {
	if m == nil {
		return nil, errors.New("MutuallyExclusiveConstraintPB is nil")
	}

	return &mutuallyExclusiveConstraint{
		fields: m.Fields,
	}, nil
}

// MutuallyExclusive creates a constraint requiring the specified fields to be mutually exclusive.
func MutuallyExclusive(fields ...string) ObjectConstraint {
	return &mutuallyExclusiveConstraint{fields: fields}
}

type requiredTogetherConstraint struct {
	fields []string
}

// Validate implements ObjectConstraint.
func (r *requiredTogetherConstraint) Validate(values map[string]cty.Value) error {
	if r == nil {
		return errors.New("required together constraint is nil")
	}

	var foundFields []string
	for _, field := range r.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}

			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) > 0 && len(foundFields) != len(r.fields) {
		fieldNames := ""
		for i, field := range foundFields {
			if i > 0 {
				fieldNames += ", "
			}

			fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
		}

		return fmt.Errorf("fields %s are required together, but only %s is present", r.formatFieldNames(), fieldNames)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (r *requiredTogetherConstraint) ValidateSpec(t *objectType) error {
	if r == nil {
		return errors.New("required together constraint is nil")
	}

	var err error
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

// ToProtobuf implements ObjectConstraint.
func (r *requiredTogetherConstraint) ToProtobuf() (*ObjectConstraintPB, error) {
	if r == nil {
		return nil, errors.New("required together constraint is nil")
	}

	return &ObjectConstraintPB{
		Constraint: &ObjectConstraintPB_RequiredTogether{
			RequiredTogether: &RequiredTogetherConstraintPB{
				Fields: r.fields,
			},
		},
	}, nil
}

func (r *requiredTogetherConstraint) formatFieldNames() any {
	if r == nil {
		return ""
	}

	fieldNames := ""
	for i, field := range r.fields {
		if i > 0 {
			fieldNames += ", "
		}

		fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
	}

	return fieldNames
}

// ToRequiredTogetherConstraint converts the protobuf RequiredTogetherConstraintPB to a requiredTogetherConstraint.
func (r *RequiredTogetherConstraintPB) ToRequiredTogetherConstraint() (ObjectConstraint, error) {
	if r == nil {
		return nil, errors.New("RequiredTogetherConstraintPB is nil")
	}

	return &requiredTogetherConstraint{
		fields: r.Fields,
	}, nil
}

// RequiredTogether creates a constraint requiring the specified fields to be present together.
func RequiredTogether(fields ...string) ObjectConstraint {
	return &requiredTogetherConstraint{fields: fields}
}

type requiredOneOfConstraint struct {
	fields []string // List of field names of which at least one is required.
}

// Validate implements ObjectConstraint.
func (r *requiredOneOfConstraint) Validate(values map[string]cty.Value) error {
	if r == nil {
		return errors.New("required one of constraint is nil")
	}

	var foundFields []string
	for _, field := range r.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}

			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) == 0 {
		fieldNames := ""
		for i, fieldName := range r.fields {
			if i > 0 {
				fieldNames += ", "
			}

			fieldNames = fmt.Sprintf("%s%q", fieldNames, fieldName)
		}

		return fmt.Errorf("at least one of the fields %s is required", fieldNames)
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (r *requiredOneOfConstraint) ValidateSpec(t *objectType) error {
	if r == nil {
		return errors.New("required one of constraint is nil")
	}

	var err error
	for _, field := range r.fields {
		if _, ok := t.fields[field]; !ok {
			err = errors.Join(err, fmt.Errorf("field %q is not defined in the object type", field))
		}
	}

	return err
}

// ToProtobuf implements ObjectConstraint.
func (r *requiredOneOfConstraint) ToProtobuf() (*ObjectConstraintPB, error) {
	if r == nil {
		return nil, errors.New("required one of constraint is nil")
	}

	return &ObjectConstraintPB{
		Constraint: &ObjectConstraintPB_RequiredOneOf{
			RequiredOneOf: &RequiredOneOfConstraintPB{
				Fields: r.fields,
			},
		},
	}, nil
}

// ToRequiredOneOfConstraint converts the protobuf RequiredOneOfConstraintPB to a requiredOneOfConstraint.
func (r *RequiredOneOfConstraintPB) ToRequiredOneOfConstraint() (ObjectConstraint, error) {
	if r == nil {
		return nil, errors.New("RequiredOneOfConstraintPB is nil")
	}

	return &requiredOneOfConstraint{
		fields: r.Fields,
	}, nil
}

// RequiredOneOf creates a constraint requiring one of the specified fields to be present.
func RequiredOneOf(fields ...string) ObjectConstraint {
	return &requiredOneOfConstraint{fields: fields}
}

type allowedFieldValuesConstraint struct {
	field  string      // The name of the field to check.
	values []cty.Value // The list of allowed values.
}

// Validate implements ObjectConstraint.
func (a *allowedFieldValuesConstraint) Validate(values map[string]cty.Value) error {
	if a == nil {
		return errors.New("allowed field values constraint is nil")
	}

	value, ok := values[a.field]
	if !ok {
		return fmt.Errorf("field %q is not present", a.field)
	}

	if !value.IsWhollyKnown() {
		return errors.New("cannot validate unknown value")
	}

	if value.IsNull() {
		return nil // Null values are not validated
	}

	for _, allowed := range a.values {
		if value.Equals(allowed).True() {
			return nil
		}
	}

	return fmt.Errorf("field %q has an invalid value, allowed values are: %s", a.field, a.formatAllowedValues())
}

// ValidateSpec implements ObjectConstraint.
func (a *allowedFieldValuesConstraint) ValidateSpec(t *objectType) error {
	if a == nil {
		return errors.New("allowed field values constraint is nil")
	}

	field, ok := t.fields[a.field]
	if !ok {
		return fmt.Errorf("field %q is not defined in the object type", a.field)
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

	return nil
}

// ToProtobuf implements ObjectConstraint.
func (a *allowedFieldValuesConstraint) ToProtobuf() (*ObjectConstraintPB, error) {
	if a == nil {
		return nil, errors.New("allowed field values constraint is nil")
	}

	valuesPB := make([][]byte, 0, len(a.values))
	for _, v := range a.values {
		vPB, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value %v to protobuf: %w", v.GoString(), err)
		}
		valuesPB = append(valuesPB, vPB)
	}

	return &ObjectConstraintPB{
		Constraint: &ObjectConstraintPB_AllowedFieldValues{
			AllowedFieldValues: &AllowedFieldValuesConstraintPB{
				Field:  a.field,
				Values: valuesPB,
			},
		},
	}, nil
}

func (a *allowedFieldValuesConstraint) formatAllowedValues() string {
	if a == nil || len(a.values) == 0 {
		return ""
	}

	allowedValues := make([]string, 0, len(a.values))
	for _, v := range a.values {
		allowedValues = append(allowedValues, util.FormatCtyValueToString(v))
	}

	return strings.Join(allowedValues, ", ")
}

// ToAllowedFieldValuesConstraint converts the protobuf AllowedFieldValuesConstraintPB to an allowedFieldValuesConstraint.
func (a *AllowedFieldValuesConstraintPB) ToAllowedFieldValuesConstraint() (ObjectConstraint, error) {
	if a == nil {
		return nil, errors.New("AllowedFieldValuesConstraintPB is nil")
	}

	var values []cty.Value
	for _, vPB := range a.Values {
		v, err := json.Unmarshal(vPB, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value from protobuf: %w", err)
		}
		values = append(values, v)
	}

	return &allowedFieldValuesConstraint{
		field:  a.Field,
		values: values,
	}, nil
}

// AllowedFieldValues creates a constraint that checks if a field's value is one of the allowed values.
//
// This should be used as part of a conditional constraint.
// Otherwise, the AllowedValues constraint should be placed on the field directly.
func AllowedFieldValues(field string, allowedValues ...cty.Value) ObjectConstraint {
	return &allowedFieldValuesConstraint{
		field:  field,
		values: allowedValues,
	}
}

// ObjectCondition is used by constraints that require a specific condition to be met to apply.
type ObjectCondition interface {
	// IsMet checks if the condition is satisfied by the given values.
	//
	// Note that no error is returned from this method.
	// If there is an error, it means the condition is not met.
	IsMet(values map[string]cty.Value) bool

	// Description provides a human-readable description of the condition for error messages.
	Description() string

	// ValidateSpec checks if the condition is valid for the object type.
	ValidateSpec(t *objectType) error

	// ToProtobuf converts the ObjectCondition to its protobuf representation.
	ToProtobuf() (*ObjectConditionPB, error)
}

// ToObjectCondition converts the protobuf ObjectConditionPB to an ObjectCondition.
func (o *ObjectConditionPB) ToObjectCondition() (ObjectCondition, error) {
	if o == nil {
		return nil, errors.New("ObjectConditionPB is nil")
	}

	switch condition := o.Condition.(type) {
	case *ObjectConditionPB_FieldPresent:
		return condition.FieldPresent.ToFieldPresentCondition()
	case *ObjectConditionPB_FieldNotPresent:
		return condition.FieldNotPresent.ToFieldNotPresentCondition()
	case *ObjectConditionPB_FieldEquals:
		return condition.FieldEquals.ToFieldEqualsCondition()
	default:
		return nil, errors.New("unknown ObjectConditionPB type")
	}
}

type fieldPresentCondition struct {
	field string // The name of the field to check.
}

// IsMet implements ObjectCondition.
func (c *fieldPresentCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && !value.IsNull()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldPresentCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is present", c.field)
}

// ValidateSpec implements ObjectCondition.
func (c *fieldPresentCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field present condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}

	return nil
}

// ToProtobuf implements ObjectCondition.
func (c *fieldPresentCondition) ToProtobuf() (*ObjectConditionPB, error) {
	if c == nil {
		return nil, errors.New("field present condition is nil")
	}

	return &ObjectConditionPB{
		Condition: &ObjectConditionPB_FieldPresent{
			FieldPresent: &FieldPresentConditionPB{
				Field: c.field,
			},
		},
	}, nil
}

// ToFieldPresentCondition converts the protobuf FieldPresentConditionPB to a fieldPresentCondition.
func (f *FieldPresentConditionPB) ToFieldPresentCondition() (ObjectCondition, error) {
	if f == nil {
		return nil, errors.New("FieldPresentConditionPB is nil")
	}

	return &fieldPresentCondition{
		field: f.Field,
	}, nil
}

// FieldPresent creates a condition that checks if a specified field is present.
func FieldPresent(field string) ObjectCondition {
	return &fieldPresentCondition{field: field}
}

type fieldNotPresentCondition struct {
	field string // The name of the field to check.
}

// IsMet implements ObjectCondition.
func (c *fieldNotPresentCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.IsNull()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldNotPresentCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is not present", c.field)
}

// ValidateSpec implements ObjectCondition.
func (c *fieldNotPresentCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field not present condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}

	return nil
}

// ToProtobuf implements ObjectCondition.
func (c *fieldNotPresentCondition) ToProtobuf() (*ObjectConditionPB, error) {
	if c == nil {
		return nil, errors.New("field not present condition is nil")
	}

	return &ObjectConditionPB{
		Condition: &ObjectConditionPB_FieldNotPresent{
			FieldNotPresent: &FieldNotPresentConditionPB{
				Field: c.field,
			},
		},
	}, nil
}

// ToFieldNotPresentCondition converts the protobuf FieldNotPresentConditionPB to a fieldNotPresentCondition.
func (f *FieldNotPresentConditionPB) ToFieldNotPresentCondition() (ObjectCondition, error) {
	if f == nil {
		return nil, errors.New("FieldNotPresentConditionPB is nil")
	}

	return &fieldNotPresentCondition{
		field: f.Field,
	}, nil
}

// FieldNotPresent creates a condition that checks if a specified field is not present.
func FieldNotPresent(field string) ObjectCondition {
	return &fieldNotPresentCondition{field: field}
}

type fieldEqualsCondition struct {
	field string    // The name of the field to check.
	value cty.Value // The value to compare against.
}

// IsMet implements ObjectCondition.
func (c *fieldEqualsCondition) IsMet(values map[string]cty.Value) bool {
	if c == nil {
		return false
	}

	if value, ok := values[c.field]; ok {
		return value.IsWhollyKnown() && value.Equals(c.value).True()
	}

	return false
}

// Description implements ObjectCondition.
func (c *fieldEqualsCondition) Description() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("field %q is equal to %s", c.field, util.FormatCtyValueToString(c.value))
}

// ValidateSpec implements ObjectCondition.
func (c *fieldEqualsCondition) ValidateSpec(t *objectType) error {
	if c == nil {
		return errors.New("field equals condition is nil")
	}

	if _, ok := t.fields[c.field]; !ok {
		return fmt.Errorf("field %q is not defined in the object type", c.field)
	}

	return nil
}

// ToProtobuf implements ObjectCondition.
func (c *fieldEqualsCondition) ToProtobuf() (*ObjectConditionPB, error) {
	if c == nil {
		return nil, errors.New("field equals condition is nil")
	}

	valuePB, err := json.Marshal(c.value, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field equals value to protobuf: %w", err)
	}

	return &ObjectConditionPB{
		Condition: &ObjectConditionPB_FieldEquals{
			FieldEquals: &FieldEqualsConditionPB{
				Field: c.field,
				Value: valuePB,
			},
		},
	}, nil
}

// ToFieldEqualsCondition converts the protobuf FieldEqualsConditionPB to a fieldEqualsCondition.
func (f *FieldEqualsConditionPB) ToFieldEqualsCondition() (ObjectCondition, error) {
	if f == nil {
		return nil, errors.New("FieldEqualsConditionPB is nil")
	}

	value, err := json.Unmarshal(f.Value, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field equals value from protobuf: %w", err)
	}

	return &fieldEqualsCondition{
		field: f.Field,
		value: value,
	}, nil
}

// FieldEquals creates a condition that checks if a specified field has a specified value.
func FieldEquals(field string, value cty.Value) ObjectCondition {
	return &fieldEqualsCondition{field: field, value: value}
}

type conditionalConstraint struct {
	condition  ObjectCondition
	constraint ObjectConstraint
}

// Validate implements ObjectConstraint.
func (c *conditionalConstraint) Validate(values map[string]cty.Value) error {
	if c == nil {
		return errors.New("conditional constraint is nil")
	}

	if c.condition == nil {
		return nil // No condition means the constraint is never applied
	}

	if c.constraint == nil {
		return nil // No constraint means nothing to validate
	}

	if c.condition.IsMet(values) {
		if err := c.constraint.Validate(values); err != nil {
			return fmt.Errorf("conditional constraint failed: when %s, %w", c.condition.Description(), err)
		}
	}

	return nil
}

// ValidateSpec implements ObjectConstraint.
func (c *conditionalConstraint) ValidateSpec(t *objectType) error {
	var err error
	if e := c.condition.ValidateSpec(t); e != nil {
		err = e
	}

	return errors.Join(err, c.constraint.ValidateSpec(t))
}

// ToProtobuf implements ObjectConstraint.
func (c *conditionalConstraint) ToProtobuf() (*ObjectConstraintPB, error) {
	if c == nil {
		return nil, errors.New("conditional constraint is nil")
	}

	conditionPB, err := c.condition.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert condition to protobuf: %w", err)
	}

	constraintPB, err := c.constraint.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert constraint to protobuf: %w", err)
	}

	return &ObjectConstraintPB{
		Constraint: &ObjectConstraintPB_Conditional{
			Conditional: &ConditionalConstraintPB{
				Condition:  conditionPB,
				Constraint: constraintPB,
			},
		},
	}, nil
}

// ToConditionalConstraint converts the protobuf ConditionalConstraintPB to a conditionalConstraint.
func (c *ConditionalConstraintPB) ToConditionalConstraint() (ObjectConstraint, error) {
	if c == nil {
		return nil, errors.New("ConditionalConstraintPB is nil")
	}

	condition, err := c.Condition.ToObjectCondition()
	if err != nil {
		return nil, fmt.Errorf("failed to convert condition from protobuf: %w", err)
	}

	constraint, err := c.Constraint.ToObjectConstraint()
	if err != nil {
		return nil, fmt.Errorf("failed to convert constraint from protobuf: %w", err)
	}

	return &conditionalConstraint{
		condition:  condition,
		constraint: constraint,
	}, nil
}

// ConditionalConstraint creates a constraint that only applies if the specified condition is met.
func ConditionalConstraint(condition ObjectCondition, constraint ObjectConstraint) ObjectConstraint {
	return &conditionalConstraint{
		condition:  condition,
		constraint: constraint,
	}
}

type objectType struct {
	rawFields   []*objectField
	fields      map[string]*objectField
	constraints ObjectConstraints
}

// WithConstraints sets the constraints on the object to the specified ones.
func (o *objectType) WithConstraints(constraints ...ObjectConstraint) *objectType {
	o.constraints = constraints
	return o
}

// CtyType implements Type.
func (o *objectType) CtyType() cty.Type {
	if o == nil {
		return cty.NilType
	}

	if len(o.fields) == 0 {
		return cty.EmptyObject
	}

	fieldTypes := make(map[string]cty.Type, len(o.fields))
	for name, field := range o.fields {
		fieldTypes[name] = field.t.CtyType()
	}

	return cty.Object(fieldTypes)
}

// Convert implements Type.
func (o *objectType) Convert(value cty.Value) (cty.Value, error) {
	if o == nil {
		return cty.NilVal, errors.New("object type is nil")
	}

	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	if value.IsNull() {
		return cty.NullVal(o.CtyType()), nil // Make sure the null value is of the correct type
	}

	if !value.Type().IsObjectType() && !value.Type().IsMapType() {
		err := fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), o.CtyType().FriendlyName())
		return cty.NilVal, err
	}

	valueMap := value.AsValueMap()

	values, err := o.convertMap(valueMap)
	if err != nil {
		return cty.NilVal, err
	}

	return cty.ObjectVal(values), nil
}

func (o *objectType) convertMap(values map[string]cty.Value) (map[string]cty.Value, error) {
	if o == nil {
		return nil, errors.New("object type is nil")
	}

	if values == nil {
		values = make(map[string]cty.Value)
	}

	resultFields := make(map[string]cty.Value, len(o.fields))
	validKeys := []string{}

	for name, field := range o.fields {
		fieldValue := field.defaultValue
		foundAs := []string{}
		validKeys = append(validKeys, name)
		if value, ok := values[name]; ok {
			value, err := field.t.Convert(value)
			if err != nil {
				return nil, fmt.Errorf("cannot convert field %q: %w", name, err)
			}

			foundAs = append(foundAs, name)
			fieldValue = value
		}

		for _, alias := range field.aliases {
			validKeys = append(validKeys, alias)
			if value, ok := values[alias]; ok {
				value, err := field.t.Convert(value)
				if err != nil {
					return nil, fmt.Errorf("cannot convert field %q (alias %q): %w", name, alias, err)
				}

				foundAs = append(foundAs, alias)
				fieldValue = value
			}
		}

		if len(foundAs) > 1 {
			foundAsNames := ""
			for i, name := range foundAs {
				if i > 0 {
					foundAsNames += ", "
				}

				foundAsNames = fmt.Sprintf("%s%q", foundAsNames, name)
			}

			return nil, fmt.Errorf("field %q is defined multiple times as %s", name, foundAsNames)
		}

		resultFields[name] = fieldValue
	}

	invalidIndexes := []string{}
	for key := range values {
		if !slices.Contains(validKeys, key) {
			invalidIndexes = append(invalidIndexes, key)
		}
	}

	if len(invalidIndexes) > 0 {
		indexes := ""
		for i, index := range invalidIndexes {
			if i > 0 {
				indexes += ", "
			}

			indexes += fmt.Sprintf("%q", index)
		}

		return nil, fmt.Errorf("invalid indexes found: %s", indexes)
	}

	return resultFields, nil
}

// Validate implements Type.
func (o *objectType) Validate(value cty.Value) error {
	if o == nil {
		return fmt.Errorf("object type is nil")
	}

	if !value.IsWhollyKnown() {
		return fmt.Errorf("cannot convert unknown value")
	}

	if value.IsNull() {
		return nil // null values are assumed to be valid
	}

	if !value.Type().IsObjectType() && !value.Type().IsMapType() {
		return fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), o.CtyType().FriendlyName())
	}

	valueMap := value.AsValueMap()

	return o.validateMap(valueMap)
}

func (o *objectType) validateMap(values map[string]cty.Value) error {
	if o == nil {
		return fmt.Errorf("object type is nil")
	}

	if values == nil {
		values = make(map[string]cty.Value)
	}

	var err error
	for name, field := range o.fields {
		value, ok := values[name]
		if !ok {
			err = errors.Join(err, fmt.Errorf("missing field %q", name))
		}

		if e := field.validate(value); e != nil {
			err = errors.Join(err, e)
		}
	}

	e := o.constraints.Validate(values)
	return errors.Join(err, e)
}

// ValidateSpec implements Type.
func (o *objectType) ValidateSpec() error {
	if o == nil {
		return errors.New("object type is nil")
	}

	var err error
	if len(o.fields) != len(o.rawFields) {
		fieldNamesDefined := make(map[string]int, len(o.fields))
		for _, field := range o.rawFields {
			if field == nil {
				err = errors.Join(err, errors.New("nil field definition found"))
				continue
			}

			if _, ok := fieldNamesDefined[field.name]; !ok {
				fieldNamesDefined[field.name] = 0
			}

			fieldNamesDefined[field.name]++
		}

		for name, count := range fieldNamesDefined {
			if count > 1 {
				err = errors.Join(err, fmt.Errorf("field %q is defined multiple times", name))
			}
		}

		if err != nil {
			return err
		}
	}

	definedNames := map[string][]string{}
	for name, field := range o.fields {
		e := field.validateSpec()
		err = errors.Join(err, e)

		definedNames[name] = append(definedNames[name], name)
		for _, alias := range field.aliases {
			definedNames[alias] = append(definedNames[alias], name)
		}
	}

	for name, definitions := range definedNames {
		if len(definitions) > 1 {
			aliases := ""
			for i, alias := range definitions {
				if i > 0 {
					aliases += ", "
				}

				aliases = fmt.Sprintf("%s%q", aliases, alias)
			}

			err = errors.Join(err, fmt.Errorf("field %q is defined multiple times (aliases: %s)", name, aliases))
		}
	}

	if o.constraints == nil {
		return err // No constraints to validate
	}

	e := o.constraints.ValidateSpec(o)
	err = errors.Join(err, e)

	return err
}

// ToProtobuf implements Type.
func (o *objectType) ToProtobuf() (*TypePB, error) {
	if o == nil {
		return nil, errors.New("object type is nil")
	}

	fieldsPB := make(map[string]*ObjectFieldPB, len(o.rawFields))
	for _, field := range o.rawFields {
		name, fieldPB, err := field.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q to protobuf: %w", field.name, err)
		}

		fieldsPB[name] = fieldPB
	}

	constraintsPB := make([]*ObjectConstraintPB, 0, len(o.constraints))
	for _, constraint := range o.constraints {
		constraintPB, err := constraint.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint to protobuf: %w", err)
		}

		constraintsPB = append(constraintsPB, constraintPB)
	}

	return &TypePB{
		Type: &TypePB_Object{
			Object: &ObjectTypePB{
				Fields:      fieldsPB,
				Constraints: constraintsPB,
			},
		},
	}, nil
}

// ToObjectType converts the protobuf ObjectTypePB to an objectType.
func (o *ObjectTypePB) ToObjectType() (*objectType, error) {
	if o == nil {
		return nil, errors.New("ObjectTypePB is nil")
	}

	var fields []*objectField
	for name, fieldPB := range o.Fields {
		field, err := fieldPB.ToObjectField(name)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q from protobuf: %w", name, err)
		}

		fields = append(fields, field)
	}

	var constraints ObjectConstraints
	for _, constraintPB := range o.Constraints {
		constraint, err := constraintPB.ToObjectConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint from protobuf: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	return Object(fields...).WithConstraints(constraints...), nil
}

// Object creates a new object type with the specified fields.
func Object(fields ...*objectField) *objectType {
	fieldMap := make(map[string]*objectField, len(fields))
	for _, field := range fields {
		if field == nil {
			continue
		}

		fieldMap[field.name] = field
	}

	return &objectType{
		rawFields: fields,
		fields:    fieldMap,
	}
}
