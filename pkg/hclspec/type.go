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
	// Raw is a cty.DynamicPseudoType that does not perform any conversion or validation.
	Raw Type = &rawType{}
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

	// ToProtobuf converts the Type to its protobuf representation.
	ToProtobuf() (*TypePB, error)
}

// ToType converts a TypePB to its Type representation.
func (t *TypePB) ToType() (Type, error) {
	switch typ := t.Type.(type) {
	case *TypePB_Object:
		return typ.Object.ToObjectType()
	case *TypePB_List:
		return typ.List.ToListType()
	case *TypePB_Set:
		return typ.Set.ToSetType()
	case *TypePB_Map:
		return typ.Map.ToMapType()
	case *TypePB_Primitive:
		return typ.Primitive.ToPrimitiveType()
	case *TypePB_Duration:
		return typ.Duration.ToDurationType()
	case *TypePB_SensitiveString:
		return typ.SensitiveString.ToSensitiveStringType()
	case *TypePB_Raw:
		return typ.Raw.ToRawType()
	default:
		return nil, fmt.Errorf("unsupported type: %T", t.Type)
	}
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

// ToProtobuf implements Type.
func (p *primitiveType) ToProtobuf() (*TypePB, error) {
	if p == nil {
		return nil, errors.New("primitive type is nil")
	}

	var prim PrimitiveTypePB
	switch {
	case p.t.Equals(cty.String):
		prim = PrimitiveTypePB_STRING
	case p.t.Equals(cty.Number):
		prim = PrimitiveTypePB_NUMBER
	case p.t.Equals(cty.Bool):
		prim = PrimitiveTypePB_BOOL
	default:
		return nil, fmt.Errorf("unsupported primitive type: %q", p.t.FriendlyName())
	}

	return &TypePB{
		Type: &TypePB_Primitive{
			Primitive: prim,
		},
	}, nil
}

// String represents the primitive type as a friendly string.
func (p *primitiveType) String() string {
	return p.t.FriendlyName()
}

// ToPrimitiveType converts a PrimitiveTypePB to its Type representation.
func (p *PrimitiveTypePB) ToPrimitiveType() (Type, error) {
	switch *p {
	case PrimitiveTypePB_STRING:
		return &primitiveType{t: cty.String}, nil
	case PrimitiveTypePB_NUMBER:
		return &primitiveType{t: cty.Number}, nil
	case PrimitiveTypePB_BOOL:
		return &primitiveType{t: cty.Bool}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive type: %v", *p)
	}
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

// ToProtobuf implements Type.
func (d *durationType) ToProtobuf() (*TypePB, error) {
	if d == nil {
		return nil, errors.New("duration type is nil")
	}

	return &TypePB{
		Type: &TypePB_Duration{
			Duration: &DurationTypePB{},
		},
	}, nil
}

// String represents the duration type as a friendly string.
func (d *durationType) String() string {
	return "duration"
}

// ToDurationType converts a DurationTypePB to its Type representation.
func (d *DurationTypePB) ToDurationType() (Type, error) {
	return &durationType{}, nil
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

// ToProtobuf implements Type.
func (s *sensitiveStringType) ToProtobuf() (*TypePB, error) {
	if s == nil {
		return nil, errors.New("sensitive string type is nil")
	}

	return &TypePB{
		Type: &TypePB_SensitiveString{
			SensitiveString: &SensitiveStringTypePB{},
		},
	}, nil
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

// ToSensitiveStringType converts a SensitiveStringTypePB to its Type representation.
func (s *SensitiveStringTypePB) ToSensitiveStringType() (Type, error) {
	return &sensitiveStringType{}, nil
}

type rawType struct{}

// CtyType implements Type.
func (r *rawType) CtyType() cty.Type {
	return cty.DynamicPseudoType
}

// Convert implements Type.
func (r *rawType) Convert(value cty.Value) (cty.Value, error) {
	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	return value, nil // No conversion for raw type.
}

// Validate implements Type.
func (r *rawType) Validate(value cty.Value) error {
	return nil
}

// ValidateSpec implements Type.
func (r *rawType) ValidateSpec() error {
	if r == nil {
		return errors.New("raw type is nil")
	}

	return nil
}

// ToProtobuf implements Type.
func (r *rawType) ToProtobuf() (*TypePB, error) {
	if r == nil {
		return nil, errors.New("raw type is nil")
	}

	return &TypePB{
		Type: &TypePB_Raw{
			Raw: &RawTypePB{},
		},
	}, nil
}

// String represents the raw type as a friendly string.
func (r *rawType) String() string {
	return r.CtyType().FriendlyName()
}

// ToRawType converts a RawTypePB to its Type representation.
func (r *RawTypePB) ToRawType() (Type, error) {
	return &rawType{}, nil
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
		return cty.NilVal, fmt.Errorf(
			"cannot convert %q to %q: %w",
			value.Type().FriendlyName(),
			targetType.FriendlyName(),
			err,
		)
	}

	return converted, nil
}
