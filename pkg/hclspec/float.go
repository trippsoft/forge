// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"math"

	"github.com/zclconf/go-cty/cty"
)

type floatType struct {
	bitSize     int32
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (f *floatType) WithConstraints(constraints ...TypeConstraint) Type {
	if f == nil {
		return nil
	}

	f.constraints = append(f.constraints, constraints...)
	return f
}

// CtyType implements [Type].
func (f *floatType) CtyType() cty.Type {
	return cty.Number
}

// Convert implements [Type].
func (f *floatType) Convert(value cty.Value) (cty.Value, error) {
	if f == nil {
		return cty.NilVal, errors.New("floating-point type is nil")
	}

	return convertCtyType(value, f.CtyType())
}

// Validate implements [Type].
func (f *floatType) Validate(value cty.Value) error {
	if f == nil {
		return errors.New("floating-point type is nil")
	}

	if value.IsNull() {
		return nil
	}

	if value.LessThan(f.minimumValue()).True() || value.GreaterThan(f.maximumValue()).True() {
		return errors.New("value is out of range")
	}

	return f.constraints.Validate(value)
}

// ToProtobuf implements [Type].
func (f *floatType) ToProtobuf() (*TypePB, error) {
	if f == nil {
		return nil, errors.New("floating-point type is nil")
	}

	constraints, err := f.constraints.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Float{
			Float: &FloatTypePB{
				BitSize:     f.bitSize,
				Constraints: constraints,
			},
		},
	}, nil
}

// String returns a string representation of the primitive type.
func (f *floatType) String() string {
	if f == nil {
		return "nil"
	}

	return fmt.Sprintf("%d-bit floating-point", f.bitSize)
}

func (f *floatType) minimumValue() cty.Value {
	switch f.bitSize {
	case 32:
		return cty.NumberFloatVal(-math.MaxFloat32)
	case 64:
		return cty.NumberFloatVal(-math.MaxFloat64)
	}

	panic("invalid float type configuration")
}

func (f *floatType) maximumValue() cty.Value {
	switch f.bitSize {
	case 32:
		return cty.NumberFloatVal(math.MaxFloat32)
	case 64:
		return cty.NumberFloatVal(math.MaxFloat64)
	}

	panic("invalid float type configuration")
}

// Float32 returns a Type representing a 32-bit floating-point number primitive.
func Float32() Type {
	return &floatType{
		bitSize: 32,
	}
}

// Float64 returns a Type representing a 64-bit floating-point number primitive.
func Float64() Type {
	return &floatType{
		bitSize: 64,
	}
}

// ToType converts a FloatTypePB to its Type representation.
func (p *FloatTypePB) ToType() (Type, error) {
	if p == nil {
		return nil, errors.New("floating-point type protobuf is nil")
	}

	constraints := make(TypeConstraints, 0, len(p.Constraints))
	for _, c := range p.Constraints {
		if c == nil {
			continue
		}

		constraint, err := c.ToTypeConstraint()
		if err != nil {
			return nil, fmt.Errorf("failed to convert constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	switch p.BitSize {
	case 32:
		return Float32().WithConstraints(constraints...), nil
	case 64:
		return Float64().WithConstraints(constraints...), nil
	default:
		return nil, fmt.Errorf("unsupported bit size for floating-point: %d", p.BitSize)
	}
}
