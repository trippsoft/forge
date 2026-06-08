// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type mapType struct {
	valueType   Type
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (m *mapType) WithConstraints(constraints ...TypeConstraint) Type {
	if m == nil {
		return nil
	}

	m.constraints = append(m.constraints, constraints...)
	return m
}

// CtyType implements [Type].
func (m *mapType) CtyType() cty.Type {
	if m == nil {
		return cty.NilType
	}

	return cty.Map(m.valueType.CtyType())
}

// Convert implements [Type].
func (m *mapType) Convert(value cty.Value) (cty.Value, error) {
	if m == nil {
		return cty.NilVal, errors.New("map type is nil")
	}

	converted, err := convertCtyType(value, m.CtyType())
	if err != nil {
		return converted, err
	}

	sensitiveString, ok := m.valueType.(*sensitiveStringType)
	if converted.IsNull() || !ok {
		return converted, nil // A null is presumed valid.
	}

	it := converted.ElementIterator()
	for it.Next() {
		_, elem := it.Element()
		sensitiveString.AddToFilter(elem)
	}

	return converted, nil
}

// Validate implements [Type].
func (m *mapType) Validate(value cty.Value) error {
	if m == nil {
		return errors.New("map type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	err := m.constraints.Validate(value)

	it := value.ElementIterator()
	for it.Next() {
		index, elem := it.Element()
		e := m.valueType.Validate(elem)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("element at index %q: %w", index.AsString(), e))
		}
	}

	return err
}

// ToProtobuf implements [Type].
func (m *mapType) ToProtobuf() (*TypePB, error) {
	if m == nil {
		return nil, errors.New("map type is nil")
	}

	valueType, err := m.valueType.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Map{
			Map: &MapTypePB{
				ElementType: valueType,
			},
		},
	}, nil
}

// String represents the map type as a friendly string.
func (m *mapType) String() string {
	return fmt.Sprintf("map of %s", m.valueType)
}

// Map returns a Type representing a map of the given value type.
func Map(valueType Type) Type {
	return &mapType{valueType: valueType}
}

// ToType converts a protobuf MapTypePB to a mapType instance.
func (m *MapTypePB) ToType() (Type, error) {
	if m == nil {
		return nil, errors.New("MapTypePB is nil")
	}

	if m.ElementType == nil {
		return nil, errors.New("ElementType in MapTypePB is nil")
	}

	elementType, err := m.ElementType.ToType()
	if err != nil {
		return nil, err
	}

	constraints := make(TypeConstraints, 0, len(m.Constraints))
	for _, c := range m.Constraints {
		if c == nil {
			continue
		}

		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	return Map(elementType).WithConstraints(constraints...), nil
}
