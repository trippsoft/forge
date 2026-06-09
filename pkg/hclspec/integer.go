// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"math"

	"github.com/zclconf/go-cty/cty"
)

type integerType struct {
	signed      bool
	bitSize     int32
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (i *integerType) WithConstraints(constraints ...TypeConstraint) Type {
	if i == nil {
		return nil
	}

	i.constraints = append(i.constraints, constraints...)
	return i
}

// CtyType implements [Type].
func (i *integerType) CtyType() cty.Type {
	return cty.Number
}

// Convert implements [Type].
func (i *integerType) Convert(value cty.Value) (cty.Value, error) {
	if i == nil {
		return cty.NilVal, errors.New("integer type is nil")
	}

	return convertCtyType(value, i.CtyType())
}

// Validate implements [Type].
func (i *integerType) Validate(value cty.Value) error {
	if i == nil {
		return errors.New("integer type is nil")
	}

	if value.IsNull() {
		return nil
	}

	if !value.AsBigFloat().IsInt() {
		return errors.New("value is not an integer")
	}

	if value.LessThan(i.minimumValue()).True() || value.GreaterThan(i.maximumValue()).True() {
		return errors.New("value is out of range")
	}

	return i.constraints.Validate(value)
}

// ToProtobuf implements [Type].
func (i *integerType) ToProtobuf() (*TypePB, error) {
	if i == nil {
		return nil, errors.New("integer type is nil")
	}

	constraints, err := i.constraints.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Integer{
			Integer: &IntegerTypePB{
				Signed:      i.signed,
				BitSize:     i.bitSize,
				Constraints: constraints,
			},
		},
	}, nil
}

// String returns a string representation of the primitive type.
func (i *integerType) String() string {
	if i == nil {
		return "nil"
	}

	signedStr := "unsigned"
	if i.signed {
		signedStr = "signed"
	}

	return fmt.Sprintf("%d-bit %s integer", i.bitSize, signedStr)
}

func (i *integerType) minimumValue() cty.Value {
	if i.signed {
		switch i.bitSize {
		case 8:
			return cty.NumberIntVal(math.MinInt8)
		case 16:
			return cty.NumberIntVal(math.MinInt16)
		case 32:
			return cty.NumberIntVal(math.MinInt32)
		case 64:
			return cty.NumberIntVal(math.MinInt64)
		}
	} else {
		return cty.Zero
	}

	panic("invalid integer type configuration")
}

func (i *integerType) maximumValue() cty.Value {
	if i.signed {
		switch i.bitSize {
		case 8:
			return cty.NumberIntVal(math.MaxInt8)
		case 16:
			return cty.NumberIntVal(math.MaxInt16)
		case 32:
			return cty.NumberIntVal(math.MaxInt32)
		case 64:
			return cty.NumberIntVal(math.MaxInt64)
		}
	} else {
		switch i.bitSize {
		case 8:
			return cty.NumberUIntVal(math.MaxUint8)
		case 16:
			return cty.NumberUIntVal(math.MaxUint16)
		case 32:
			return cty.NumberUIntVal(math.MaxUint32)
		case 64:
			return cty.NumberUIntVal(math.MaxUint64)
		}
	}

	panic("invalid integer type configuration")
}

// Int8 returns a Type representing an 8-bit signed integer.
func Int8() Type {
	return &integerType{
		signed:  true,
		bitSize: 8,
	}
}

// Int16 returns a Type representing a 16-bit signed integer.
func Int16() Type {
	return &integerType{
		signed:  true,
		bitSize: 16,
	}
}

// Int32 returns a Type representing a 32-bit signed integer.
func Int32() Type {
	return &integerType{
		signed:  true,
		bitSize: 32,
	}
}

// Int64 returns a Type representing a 64-bit signed integer.
func Int64() Type {
	return &integerType{
		signed:  true,
		bitSize: 64,
	}
}

// UInt8 returns a Type representing an 8-bit unsigned integer.
func UInt8() Type {
	return &integerType{
		signed:  false,
		bitSize: 8,
	}
}

// UInt16 returns a Type representing a 16-bit unsigned integer.
func UInt16() Type {
	return &integerType{
		signed:  false,
		bitSize: 16,
	}
}

// UInt32 returns a Type representing a 32-bit unsigned integer.
func UInt32() Type {
	return &integerType{
		signed:  false,
		bitSize: 32,
	}
}

// UInt64 returns a Type representing a 64-bit unsigned integer.
func UInt64() Type {
	return &integerType{
		signed:  false,
		bitSize: 64,
	}
}

// ToType converts a IntegerTypePB to its Type representation.
func (p *IntegerTypePB) ToType() (Type, error) {
	if p == nil {
		return nil, errors.New("integer type protobuf is nil")
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

	if p.Signed {
		switch p.BitSize {
		case 8:
			return Int8().WithConstraints(constraints...), nil
		case 16:
			return Int16().WithConstraints(constraints...), nil
		case 32:
			return Int32().WithConstraints(constraints...), nil
		case 64:
			return Int64().WithConstraints(constraints...), nil
		default:
			return nil, fmt.Errorf("unsupported bit size for signed integer: %d", p.BitSize)
		}
	}

	switch p.BitSize {
	case 8:
		return UInt8().WithConstraints(constraints...), nil
	case 16:
		return UInt16().WithConstraints(constraints...), nil
	case 32:
		return UInt32().WithConstraints(constraints...), nil
	case 64:
		return UInt64().WithConstraints(constraints...), nil
	default:
		return nil, fmt.Errorf("unsupported bit size for unsigned integer: %d", p.BitSize)
	}
}
