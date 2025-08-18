// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"time"

	"github.com/trippsoft/forge/pkg/log"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

var (
	// String is a cty.String with no additional validation.
	String Type = &primitiveType{t: cty.String}
	// SensitiveString is a cty.String that should be added to the secret filter.
	SensitiveString Type = &sensitiveStringType{}
	// Number is a cty.Number with no additional validation.
	Number Type = &primitiveType{t: cty.Number}
	// Bool is a cty.Bool with no additional validation.
	Bool Type = &primitiveType{t: cty.Bool}
	// Duration is a cty.String that is expected to be a duration string (e.g. 1m30s).
	Duration Type = &durationType{}
)

// Type represents a data type within an argument spec.
//
// This will wrap a cty.Type and may provide additional validation or conversion logic for custom types.
type Type interface {
	// CtyType returns the cty.Type representation of the type.
	CtyType() cty.Type

	// Convert converts a cty.Value to this type.
	//
	// This will provide an error if the conversion is not possible.
	// Implementations of this function should produce all errors on failure.
	// This function should be called before Validate and the returned value should be provided to it.
	Convert(value cty.Value) (cty.Value, error)

	// Validate checks if a cty.Value is valid for this type.
	//
	// This should be called after Convert with the returned value provided.
	// Implementations should do as much validation as possible on failure to provide as much feedback to the user.
	Validate(value cty.Value) error

	// ValidateSpec checks if the type specification is valid.
	//
	// This function will be called by input specifications to ensure they are valid.
	// This is primarily used for validating nested object types.
	ValidateSpec() error
}

type primitiveType struct {
	t cty.Type
}

// CtyType implements Type.
func (p *primitiveType) CtyType() cty.Type {
	if p == nil {
		return cty.NilType
	}

	return p.t
}

// Convert implements Type.
func (p *primitiveType) Convert(value cty.Value) (cty.Value, error) {
	if p == nil {
		return cty.NilVal, errors.New("primitive type is nil")
	}

	return convertCtyType(value, p.t)
}

// Validate implements Type.
func (p *primitiveType) Validate(value cty.Value) error {
	return nil // No additional validation for primitive types.
}

// ValidateSpec implements Type.
func (p *primitiveType) ValidateSpec() error {
	if p == nil {
		return errors.New("primitive type is nil")
	}

	if p.t.Equals(cty.NilType) {
		return errors.New("primitive type is nil")
	}

	if !p.t.IsPrimitiveType() {
		return fmt.Errorf("Type %q is not a primitive type", p.t.FriendlyName())
	}

	return nil
}

// String represents the primitive type as a friendly string.
func (p *primitiveType) String() string {
	return p.t.FriendlyName()
}

type durationType struct{}

// CtyType implements Type.
func (d *durationType) CtyType() cty.Type {
	return cty.String
}

// Convert implements Type.
func (d *durationType) Convert(value cty.Value) (cty.Value, error) {
	return convertCtyType(value, cty.String)
}

// Validate implements Type.
func (d *durationType) Validate(value cty.Value) error {
	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	if !value.Type().Equals(cty.String) {
		return fmt.Errorf("expected string type for duration, got %s", value.Type().FriendlyName())
	}

	_, err := time.ParseDuration(value.AsString())
	return err
}

// ValidateSpec implements Type.
func (d *durationType) ValidateSpec() error {
	if d == nil {
		return errors.New("duration type is nil")
	}

	return nil
}

// String represents the duration type as a friendly string.
func (d *durationType) String() string {
	return "duration"
}

type sensitiveStringType struct{}

// CtyType implements Type.
func (s *sensitiveStringType) CtyType() cty.Type {
	return cty.String
}

// Convert implements Type.
func (s *sensitiveStringType) Convert(value cty.Value) (cty.Value, error) {
	if s == nil {
		return cty.NilVal, errors.New("sensitive string type is nil")
	}

	v, err := convertCtyType(value, s.CtyType())
	if err == nil {
		s.AddToFilter(v)
	}

	return v, err
}

// Validate implements Type.
func (s *sensitiveStringType) Validate(value cty.Value) error {
	return nil // No additional validation for sensitive string types.
}

// ValidateSpec implements Type.
func (s *sensitiveStringType) ValidateSpec() error {
	if s == nil {
		return errors.New("sensitive string type is nil")
	}

	return nil
}

// String represents the sensitive string type as a friendly string.
func (s *sensitiveStringType) String() string {
	return "sensitive string"
}

// AddToFilter adds the sensitive string value to the secret filter.
func (s *sensitiveStringType) AddToFilter(value cty.Value) {
	if value.IsNull() {
		return // No need to add null values to the filter.
	}

	v := value.AsString()
	if v != "" {
		log.SecretFilter.AddSecret(v)
	}
}

type listType struct {
	elementType Type // The type of elements in the list.
}

// List creates a new listType representing a list of the given element type.
func List(elementType Type) Type {
	return &listType{elementType: elementType}
}

// CtyType implements Type.
func (l *listType) CtyType() cty.Type {
	if l == nil {
		return cty.NilType
	}

	return cty.List(l.elementType.CtyType())
}

// Convert implements Type.
func (l *listType) Convert(value cty.Value) (cty.Value, error) {
	if l == nil {
		return cty.NilVal, errors.New("list type is nil")
	}

	converted, err := convertCtyType(value, l.CtyType())
	if err != nil {
		var e error
		converted, e = l.elementType.Convert(value)
		if e != nil {
			return cty.NilVal, err
		}

		return cty.ListVal([]cty.Value{converted}), nil
	}

	if converted.IsNull() {
		return converted, nil
	}

	it := converted.ElementIterator()
	values := make([]cty.Value, 0, converted.LengthInt())
	sensitiveString, isSensitiveString := l.elementType.(*sensitiveStringType)
	for it.Next() {
		_, elem := it.Element()
		if elem.IsNull() {
			continue // Skip null elements.
		}

		if isSensitiveString {
			sensitiveString.AddToFilter(elem)
		}

		values = append(values, elem)
	}

	if len(values) == 0 {
		return cty.ListValEmpty(l.elementType.CtyType()), nil
	}

	return cty.ListVal(values), nil
}

// Validate implements Type.
func (l *listType) Validate(value cty.Value) error {
	if l == nil {
		return errors.New("list type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	var err error
	for it.Next() {
		index, elem := it.Element()
		e := l.elementType.Validate(elem)
		if e != nil {
			i, _ := index.AsBigFloat().Int64()
			err = errors.Join(err, fmt.Errorf("element at index %d: %w", i, e))
		}
	}

	return err
}

// ValidateSpec implements Type.
func (l *listType) ValidateSpec() error {
	if l == nil {
		return errors.New("list type is nil")
	}

	return l.elementType.ValidateSpec()
}

// String represents the list type as a friendly string.
func (l *listType) String() string {
	return l.CtyType().FriendlyName()
}

type mapType struct {
	valueType Type // The type of values in the map.
}

// Map creates a new mapType representing a map of string keys to the given value type.
func Map(valueType Type) Type {
	return &mapType{valueType: valueType}
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

// String represents the map type as a friendly string.
func (m *mapType) String() string {
	return m.CtyType().FriendlyName()
}

// setType represents a set of unique elements of a specific type.
type setType struct {
	elementType Type // The type of elements in the set.
}

// Set creates a new setType representing a set of the given element type.
func Set(elementType Type) Type {
	return &setType{elementType: elementType}
}

// CtyType implements Type.
func (s *setType) CtyType() cty.Type {
	if s == nil {
		return cty.NilType
	}

	return cty.Set(s.elementType.CtyType())
}

// Convert implements Type.
func (s *setType) Convert(value cty.Value) (cty.Value, error) {
	if s == nil {
		return cty.NilVal, errors.New("set type is nil")
	}

	converted, err := convertCtyType(value, s.CtyType())
	if err != nil {
		var e error
		converted, e = s.elementType.Convert(value)
		if e != nil {
			return cty.NilVal, err
		}

		return cty.SetVal([]cty.Value{converted}), nil
	}

	if converted.IsNull() {
		return converted, nil
	}

	it := converted.ElementIterator()
	values := make([]cty.Value, 0, converted.LengthInt())
	sensitiveString, isSensitiveString := s.elementType.(*sensitiveStringType)
	for it.Next() {
		_, elem := it.Element()
		if elem.IsNull() {
			continue
		}

		if isSensitiveString {
			sensitiveString.AddToFilter(elem)
		}

		values = append(values, elem)
	}

	if len(values) == 0 {
		return cty.SetValEmpty(s.elementType.CtyType()), nil
	}

	return cty.SetVal(values), nil
}

// Validate implements Type.
func (s *setType) Validate(value cty.Value) error {
	if s == nil {
		return errors.New("set type is nil")
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	var err error
	for it.Next() {
		_, elem := it.Element()
		e := s.elementType.Validate(elem)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid set element: %w", e))
		}
	}

	return err
}

// ValidateSpec implements Type.
func (s *setType) ValidateSpec() error {
	if s == nil {
		return errors.New("set type is nil")
	}

	return s.elementType.ValidateSpec()
}

// String represents the set type as a friendly string.
func (s *setType) String() string {
	return s.CtyType().FriendlyName()
}

func convertCtyType(value cty.Value, targetType cty.Type) (cty.Value, error) {
	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	if value.IsNull() {
		return cty.NullVal(targetType), nil
	}

	if value.Type().Equals(targetType) {
		return value, nil
	}

	conversion := convert.GetConversionUnsafe(value.Type(), targetType)
	if conversion == nil {
		return cty.NilVal, fmt.Errorf("cannot convert %q to %q", value.Type().FriendlyName(), targetType.FriendlyName())
	}

	converted, err := conversion(value)
	if err != nil {
		return cty.NilVal, fmt.Errorf("cannot convert %q to %q: %w", value.Type().FriendlyName(), targetType.FriendlyName(), err)
	}

	return converted, nil
}
