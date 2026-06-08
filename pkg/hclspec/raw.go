// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
)

var (
	// Raw is a cty.DynamicPseudoType that does not perform any conversion or validation.
	Raw Type = &rawType{}
)

type rawType struct{}

// WithConstraints implements [Type].
func (r *rawType) WithConstraints(constraints ...TypeConstraint) Type {
	return r // Raw type does not support constraints.
}

// CtyType implements [Type].
func (r *rawType) CtyType() cty.Type {
	return cty.DynamicPseudoType
}

// Convert implements [Type].
func (r *rawType) Convert(value cty.Value) (cty.Value, error) {
	if !value.IsWhollyKnown() {
		return cty.NilVal, errors.New("cannot convert unknown value")
	}

	return value, nil // No conversion for raw type.
}

// Validate implements [Type].
func (r *rawType) Validate(value cty.Value) error {
	return nil
}

// ToProtobuf implements [Type].
func (r *rawType) ToProtobuf() (*TypePB, error) {
	return &TypePB{
		Type: &TypePB_Raw{
			Raw: &RawTypePB{},
		},
	}, nil
}

func (r *rawType) String() string {
	return "raw data"
}

// ToType converts a RawTypePB to its Type representation.
func (r *RawTypePB) ToType() (Type, error) {
	return Raw, nil
}
