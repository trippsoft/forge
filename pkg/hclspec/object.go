// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"slices"

	"github.com/zclconf/go-cty/cty"
)

var (
	_ Type = (*objectType)(nil)
)

type objectType struct {
	fields      map[string]*objectField
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (o *objectType) WithConstraints(constraints ...TypeConstraint) Type {
	if o == nil {
		return nil
	}

	o.constraints = append(o.constraints, constraints...)
	return o
}

// CtyType implements [Type].
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

// Convert implements [Type].
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

// Validate implements [Type].
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

	var err error
	for name, field := range o.fields {
		value, ok := valueMap[name]
		if !ok {
			err = errors.Join(err, fmt.Errorf("missing field %q", name))
			continue
		}

		err = errors.Join(err, field.Validate(value))
	}

	return errors.Join(err, o.constraints.Validate(value))
}

// ToProtobuf implements [Type].
func (o *objectType) ToProtobuf() (*TypePB, error) {
	if o == nil {
		return nil, errors.New("object type is nil")
	}

	fields := make(map[string]*ObjectFieldPB, len(o.fields))
	for _, f := range o.fields {
		name, field, err := f.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q to protobuf: %w", f.name, err)
		}

		fields[name] = field
	}

	constraints, err := o.constraints.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert constraints to protobuf: %w", err)
	}

	return &TypePB{
		Type: &TypePB_Object{
			Object: &ObjectTypePB{
				Fields:      fields,
				Constraints: constraints,
			},
		},
	}, nil
}

// String represents the object type as a friendly string.
func (o *objectType) String() string {
	return "object"
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

// Object creates a new object type with the given fields.
func Object(fields ...*objectField) *objectType {
	fieldMap := make(map[string]*objectField, len(fields))
	for _, field := range fields {
		if field == nil {
			continue
		}

		fieldMap[field.name] = field
	}

	return &objectType{
		fields: fieldMap,
	}
}

// ToType converts the protobuf ObjectTypePB to an objectType.
func (o *ObjectTypePB) ToType() (*objectType, error) {
	if o == nil {
		return nil, errors.New("ObjectTypePB is nil")
	}

	var fields []*objectField
	for name, f := range o.Fields {
		field, err := f.ToObjectField(name)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field %q from protobuf: %w", name, err)
		}

		fields = append(fields, field)
	}

	var constraints TypeConstraints
	for _, c := range o.Constraints {
		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint from protobuf: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	obj := Object(fields...)
	obj.WithConstraints(constraints...)

	return obj, nil
}
