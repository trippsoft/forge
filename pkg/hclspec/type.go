// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

var ()

// Type represents a data type within an argument spec.
//
// This will wrap a cty.Type and may provide additional validation or conversion logic for custom types.
type Type interface {
	// WithConstraints returns a new Type with the given constraints added.
	WithConstraints(constraints ...TypeConstraint) Type

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

	// ToProtobuf converts the Type to its protobuf representation.
	ToProtobuf() (*TypePB, error)
}

// ToType converts a TypePB to its Type representation.
func (t *TypePB) ToType() (Type, error) {
	switch typ := t.Type.(type) {
	case *TypePB_Raw:
		return typ.Raw.ToType()
	case *TypePB_Simple:
		return typ.Simple.ToType()
	case *TypePB_Set:
		return typ.Set.ToType()
	case *TypePB_List:
		return typ.List.ToType()
	case *TypePB_Map:
		return typ.Map.ToType()
	case *TypePB_Object:
		return typ.Object.ToType()
	default:
		return nil, fmt.Errorf("unsupported type: %T", t.Type)
	}
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
