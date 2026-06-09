// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
)

type primitiveType struct {
	t           cty.Type
	constraints TypeConstraints
}

// WithConstraints implements [Type].
func (p *primitiveType) WithConstraints(constraints ...TypeConstraint) Type {
	if p == nil {
		return nil
	}

	p.constraints = append(p.constraints, constraints...)
	return p
}

// CtyType implements [Type].
func (p *primitiveType) CtyType() cty.Type {
	if p == nil {
		return cty.NilType
	}

	return p.t
}

// Convert implements [Type].
func (p *primitiveType) Convert(value cty.Value) (cty.Value, error) {
	if p == nil {
		return cty.NilVal, errors.New("primitive type is nil")
	}

	return convertCtyType(value, p.CtyType())
}

// Validate implements [Type].
func (p *primitiveType) Validate(value cty.Value) error {
	if p == nil {
		return errors.New("primitive type is nil")
	}

	return p.constraints.Validate(value)
}

// ToProtobuf implements [Type].
func (p *primitiveType) ToProtobuf() (*TypePB, error) {
	if p == nil {
		return nil, errors.New("primitive type is nil")
	}

	var data SimpleTypeDataPB
	switch {
	case p.t.Equals(cty.Bool):
		data = SimpleTypeDataPB_BOOL
	case p.t.Equals(cty.String):
		data = SimpleTypeDataPB_STRING
	default:
		return nil, errors.New("unsupported primitive type")
	}

	constraints, err := p.constraints.ToProtobuf()
	if err != nil {
		return nil, err
	}

	return &TypePB{
		Type: &TypePB_Simple{
			Simple: &SimpleTypePB{
				Data:        data,
				Constraints: constraints,
			},
		},
	}, nil
}

// String returns a string representation of the primitive type.
func (p *primitiveType) String() string {
	if p == nil {
		return "nil"
	}

	return p.CtyType().FriendlyName()
}

// Bool returns a Type representing a boolean primitive.
func Bool() Type {
	return &primitiveType{t: cty.Bool}
}

// String returns a Type representing a string primitive.
func String() Type {
	return &primitiveType{t: cty.String}
}
