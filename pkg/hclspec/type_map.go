// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type mapType struct {
	valueType Type // The type of values in the map.
}

// CtyType implements Type.
func (m *mapType) CtyType() cty.Type {
	if m == nil {
		return cty.NilType
	}

	return cty.Map(m.valueType.CtyType())
}

// Convert implements Type.
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

// Validate implements Type.
func (m *mapType) Validate(value cty.Value) error {
	if m == nil {
		return errors.New("map type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	var err error
	for it.Next() {
		index, elem := it.Element()
		e := m.valueType.Validate(elem)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("element at index %q: %w", index.AsString(), e))
		}
	}

	return err
}

// ValidateSpec implements Type.
func (m *mapType) ValidateSpec() error {
	if m == nil {
		return errors.New("map type is nil")
	}

	return m.valueType.ValidateSpec()
}

// ToProtobuf implements Type.
func (m *mapType) ToProtobuf() (*TypePB, error) {
	if m == nil {
		return nil, errors.New("map type is nil")
	}

	valueTypePB, err := m.valueType.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Map{
			Map: &MapTypePB{
				ElementType: valueTypePB,
			},
		},
	}, nil
}

// String represents the map type as a friendly string.
func (m *mapType) String() string {
	return m.CtyType().FriendlyName()
}

// Map creates a new mapType representing a map of string keys to the given value type.
func Map(valueType Type) Type {
	return &mapType{valueType: valueType}
}

// ToMapType converts a protobuf MapTypePB to a mapType instance.
func (m *MapTypePB) ToMapType() (Type, error) {
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

	return Map(elementType), nil
}
