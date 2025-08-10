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
	String          Type = &primitiveType{t: cty.String}
	SensitiveString Type = &sensitiveStringType{}
	Number          Type = &primitiveType{t: cty.Number}
	Bool            Type = &primitiveType{t: cty.Bool}
	Duration        Type = &durationType{}
)

// Type represents an HCL type within an argument spec.
type Type interface {
	CtyType() cty.Type                          // CtyType returns the cty.Type representation of the HCL type.
	Convert(value cty.Value) (cty.Value, error) // Convert converts a cty.Value to this type.
	Validate(value cty.Value) error             // Validate checks if a cty.Value is valid for this type.
	ValidateSpec() []error                      // ValidateSpec checks if the type specification is valid.
}

// primitiveType represents a primitive HCL type.
type primitiveType struct {
	t cty.Type
}

// CtyType implements Type.
func (p *primitiveType) CtyType() cty.Type {
	return p.t
}

// Convert implements Type.
func (p *primitiveType) Convert(value cty.Value) (cty.Value, error) {
	return convertCtyType(value, p.t)
}

// Validate implements Type.
func (p *primitiveType) Validate(value cty.Value) error {

	_, err := p.Convert(value)
	if err != nil {
		return err
	}

	return nil
}

// ValidateSpec implements Type.
func (p *primitiveType) ValidateSpec() []error {
	return []error{}
}

// String represents the primitive type as a friendly string.
func (p *primitiveType) String() string {
	return p.t.FriendlyName()
}

// durationType represents a time.Duration represented in HCL as a string.
type durationType struct{}

// CtyType implements Type.
func (d *durationType) CtyType() cty.Type {
	return cty.String
}

// Convert implements Type.
func (d *durationType) Convert(value cty.Value) (cty.Value, error) {
	return convertCtyType(value, d.CtyType())
}

// Validate implements Type.
func (d *durationType) Validate(value cty.Value) error {

	var err error
	value, err = d.Convert(value)
	if err != nil {
		return err
	}

	_, err = time.ParseDuration(value.AsString())
	return err
}

// ValidateSpec implements Type.
func (d *durationType) ValidateSpec() []error {
	return []error{}
}

// String represents the duration type as a friendly string.
func (d *durationType) String() string {
	return "duration"
}

// sensitiveStringType represents a sensitive string type.
type sensitiveStringType struct{}

// CtyType implements Type.
func (s *sensitiveStringType) CtyType() cty.Type {
	return cty.String
}

// Convert implements Type.
func (s *sensitiveStringType) Convert(value cty.Value) (cty.Value, error) {
	v, err := convertCtyType(value, s.CtyType())
	if err == nil && !v.IsNull() {
		valStr := v.AsString()
		if valStr != "" {
			log.SecretFilter.AddSecret(valStr)
		}
	}
	return v, err
}

// Validate implements Type.
func (s *sensitiveStringType) Validate(value cty.Value) error {

	var err error
	value, err = s.Convert(value)
	if err != nil {
		return err
	}

	return nil
}

// ValidateSpec implements Type.
func (s *sensitiveStringType) ValidateSpec() []error {
	return []error{}
}

// String represents the sensitive string type as a friendly string.
func (s *sensitiveStringType) String() string {
	return "sensitive string"
}

// listType represents a list of elements of a specific type.
type listType struct {
	elementType Type // The type of elements in the list.
}

// List creates a new listType representing a list of the given element type.
func List(elementType Type) Type {
	return &listType{elementType: elementType}
}

// CtyType implements Type.
func (l *listType) CtyType() cty.Type {
	return cty.List(l.elementType.CtyType())
}

// Convert implements Type.
func (l *listType) Convert(value cty.Value) (cty.Value, error) {

	convertedValue, listErr := convertCtyType(value, l.CtyType())

	if listErr != nil {
		var err error
		convertedValue, err = l.elementType.Convert(value)
		if err != nil {
			return cty.NilVal, listErr
		}

		convertedValue = cty.ListVal([]cty.Value{convertedValue})
	}

	if convertedValue.IsNull() {
		return convertedValue, nil
	}

	it := convertedValue.ElementIterator()
	values := make([]cty.Value, 0, convertedValue.LengthInt())
	for it.Next() {
		_, elem := it.Element()
		if !elem.IsNull() {
			values = append(values, elem)
		}
	}

	if len(values) == 0 {
		return cty.ListValEmpty(l.elementType.CtyType()), nil
	}

	return cty.ListVal(values), nil
}

// Validate implements Type.
func (l *listType) Validate(value cty.Value) error {

	var err error
	value, err = l.Convert(value)
	if err != nil {
		return err
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	for it.Next() {
		index, elem := it.Element()
		err := l.elementType.Validate(elem)
		if err != nil {
			i, _ := index.AsBigFloat().Int64()
			return fmt.Errorf("element at index %d: %w", i, err)
		}
	}

	return nil
}

// ValidateSpec implements Type.
func (l *listType) ValidateSpec() []error {
	return l.elementType.ValidateSpec()
}

// String represents the list type as a friendly string.
func (l *listType) String() string {
	return l.CtyType().FriendlyName()
}

// mapType represents a map of string keys to values of a specific type.
type mapType struct {
	valueType Type // The type of values in the map.
}

// Map creates a new mapType representing a map of string keys to the given value type.
func Map(valueType Type) Type {
	return &mapType{valueType: valueType}
}

// CtyType implements Type.
func (m *mapType) CtyType() cty.Type {
	return cty.Map(m.valueType.CtyType())
}

// Convert implements Type.
func (m *mapType) Convert(value cty.Value) (cty.Value, error) {
	return convertCtyType(value, m.CtyType())
}

// Validate implements Type.
func (m *mapType) Validate(value cty.Value) error {

	var err error
	value, err = m.Convert(value)
	if err != nil {
		return err
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	for it.Next() {
		index, elem := it.Element()
		err := m.valueType.Validate(elem)
		if err != nil {
			return fmt.Errorf("element at index %q: %w", index.AsString(), err)
		}
	}

	return nil
}

// ValidateSpec implements Type.
func (m *mapType) ValidateSpec() []error {
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
	return cty.Set(s.elementType.CtyType())
}

// Convert implements Type.
func (s *setType) Convert(value cty.Value) (cty.Value, error) {

	convertedValue, setErr := convertCtyType(value, s.CtyType())

	if setErr != nil {
		var err error
		convertedValue, err = s.elementType.Convert(value)
		if err != nil {
			return cty.NilVal, setErr
		}

		convertedValue = cty.SetVal([]cty.Value{convertedValue})
	}

	if convertedValue.IsNull() {
		return convertedValue, nil
	}

	it := convertedValue.ElementIterator()
	values := make([]cty.Value, 0, convertedValue.LengthInt())
	for it.Next() {
		_, elem := it.Element()
		if !elem.IsNull() {
			values = append(values, elem)
		}
	}

	if len(values) == 0 {
		return cty.SetValEmpty(s.elementType.CtyType()), nil
	}

	return cty.SetVal(values), nil
}

// Validate implements Type.
func (s *setType) Validate(value cty.Value) error {

	var err error
	value, err = s.Convert(value)
	if err != nil {
		return err
	}

	if value.IsNull() {
		return nil // A null is presumed valid.
	}

	it := value.ElementIterator()
	for it.Next() {
		_, elem := it.Element()
		err := s.elementType.Validate(elem)
		if err != nil {
			return fmt.Errorf("invalid set element: %w", err)
		}
	}

	return nil
}

// ValidateSpec implements Type.
func (s *setType) ValidateSpec() []error {
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

	convertedValue, err := conversion(value)
	if err != nil {
		return cty.NilVal, fmt.Errorf("cannot convert %q to %q: %w", value.Type().FriendlyName(), targetType.FriendlyName(), err)
	}

	return convertedValue, nil
}
