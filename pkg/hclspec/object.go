// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"slices"

	"github.com/zclconf/go-cty/cty"
)

type ObjectField struct {
	t            Type
	aliases      []string
	required     bool
	defaultValue cty.Value
}

func RequiredField(t Type, aliases ...string) *ObjectField {
	return &ObjectField{
		t:            t,
		aliases:      aliases,
		required:     true,
		defaultValue: cty.NullVal(t.CtyType()),
	}
}

func OptionalField(t Type, defaultValue cty.Value, aliases ...string) *ObjectField {
	return &ObjectField{
		t:            t,
		aliases:      aliases,
		required:     false,
		defaultValue: defaultValue,
	}
}

type objectType struct {
	fields      map[string]*ObjectField
	constraints ObjectConstraints
}

func Object(fields map[string]*ObjectField, constraints ...ObjectConstraint) *objectType {
	return &objectType{
		fields:      fields,
		constraints: constraints,
	}
}

func (o *ObjectField) validateSpec(name string) error {

	var err error
	if o.t == nil {
		err = fmt.Errorf("field %q has no type defined", name)
	}

	if len(o.aliases) > 0 {
		if slices.Contains(o.aliases, "") {
			err = errors.Join(err, fmt.Errorf("field %q has an empty alias", name))
		}
	}

	e := o.t.ValidateSpec()
	err = errors.Join(err, e)

	if !o.defaultValue.IsWhollyKnown() {
		err = errors.Join(err, fmt.Errorf("field %q has an unknown default value", name))
		return err
	}

	if o.required && !o.defaultValue.IsNull() {
		err = errors.Join(err, fmt.Errorf("field %q is required and has a default value", name))
	}

	if !o.defaultValue.Type().Equals(o.t.CtyType()) {
		err = errors.Join(err, fmt.Errorf("field %q default value type mismatch: expected %q, got %q", name, o.t.CtyType().FriendlyName(), o.defaultValue.Type().FriendlyName()))
	}

	e = o.t.Validate(o.defaultValue)
	if e != nil {
		err = errors.Join(err, fmt.Errorf("field %q default value validation failed: %w", name, e))
	}

	return err
}

// CtyType implements Type.
func (o *objectType) CtyType() cty.Type {

	fieldTypes := make(map[string]cty.Type, len(o.fields))
	for name, field := range o.fields {
		fieldTypes[name] = field.t.CtyType()
	}

	return cty.Object(fieldTypes)
}

// Convert implements Type.
func (o *objectType) Convert(value cty.Value) (cty.Value, error) {

	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	if value.IsNull() {
		return cty.NullVal(o.CtyType()), nil // Make sure the null value is of the correct type
	}

	if !value.Type().IsObjectType() && !value.Type().IsMapType() {
		return cty.NilVal, fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), o.CtyType().FriendlyName())
	}

	valueMap := value.AsValueMap()

	if valueMap == nil {
		valueMap = make(map[string]cty.Value)
	}

	values, err := o.convert(valueMap)
	if err != nil {
		return cty.NilVal, err
	}

	return cty.ObjectVal(values), nil
}

// convert converts a map of values to match the object type.
func (o *objectType) convert(values map[string]cty.Value) (map[string]cty.Value, error) {

	if values == nil {
		return nil, errors.New("cannot convert nil map")
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
			return nil, fmt.Errorf("field %q is defined multiple times as %v", name, foundAs)
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
		return nil, fmt.Errorf("invalid indexes found: %v", invalidIndexes)
	}

	return resultFields, nil
}

// Validate implements Type.
func (o *objectType) Validate(value cty.Value) error {

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

	if valueMap == nil {
		valueMap = make(map[string]cty.Value)
	}

	return o.validate(valueMap)
}

func (o *objectType) validate(values map[string]cty.Value) error {

	converted, err := o.convert(values)
	if err != nil {
		return err
	}

	for name, field := range o.fields {
		value, ok := converted[name]
		if !ok {
			return fmt.Errorf("missing field %q", name)
		}

		if field.required && value.IsNull() {
			return fmt.Errorf("missing required field %q", name)
		}

		if err := field.t.Validate(value); err != nil {
			return fmt.Errorf("field %q validation failed: %w", name, err)
		}
	}

	return o.validateConstraints(converted)
}

func (o *objectType) validateConstraints(values map[string]cty.Value) error {

	if len(o.constraints) == 0 {
		return nil // No constraints to validate
	}

	for _, constraint := range o.constraints {
		if err := constraint.Validate(values); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// ValidateSpec implements Type.
func (o *objectType) ValidateSpec() error {
	var err error
	if o.fields == nil {
		return errors.New("object type has no fields defined")
	}

	definedNames := map[string][]string{}

	for name, field := range o.fields {
		if field == nil {
			err = errors.Join(err, fmt.Errorf("field %q has no definition", name))
			continue
		}

		e := field.validateSpec(name)
		err = errors.Join(err, e)

		definedNames[name] = append(definedNames[name], name)

		for _, alias := range field.aliases {
			definedNames[alias] = append(definedNames[alias], name)
		}
	}

	for name, definitions := range definedNames {
		if len(definitions) > 1 {
			err = errors.Join(err, fmt.Errorf("field %q is defined multiple times (aliases: %v)", name, definitions))
		}
	}

	if o.constraints == nil {
		return err // No constraints to validate
	}

	for _, constraint := range o.constraints {
		if constraint == nil {
			err = errors.Join(err, errors.New("object type has a nil constraint"))
			continue
		}

		if e := constraint.ValidateSpec(o); e != nil {
			err = errors.Join(err, fmt.Errorf("constraint validation failed: %w", e))
		}
	}

	return err
}
