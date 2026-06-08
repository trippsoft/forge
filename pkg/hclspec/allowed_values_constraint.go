// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"strings"

	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

type allowedValuesConstraint struct {
	values []cty.Value
}

// Validate implements [TypeConstraint].
func (a *allowedValuesConstraint) Validate(value cty.Value) error {
	if a == nil {
		return fmt.Errorf("allowed values constraint is nil")
	}

	if value.IsNull() {
		return nil // Skip null values
	}

	for _, v := range a.values {
		if value.Equals(v).True() {
			return nil
		}
	}

	return fmt.Errorf(
		"value %v is not in allowed values: %s",
		hclutil.FormatCtyValueToString(value),
		a.formatAllowedValues(),
	)
}

// ToProtobuf implements [TypeConstraint].
func (a *allowedValuesConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if a == nil {
		return nil, fmt.Errorf("allowed values constraint is nil")
	}

	values := make([][]byte, 0, len(a.values))
	for _, v := range a.values {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value %v to protobuf: %w", v.GoString(), err)
		}

		values = append(values, value)
	}

	return &TypeConstraintPB{
		Constraint: &TypeConstraintPB_AllowedValues{
			AllowedValues: &AllowedValuesConstraintPB{
				Values: values,
			},
		},
	}, nil
}

func (a *allowedValuesConstraint) formatAllowedValues() string {
	if a == nil || len(a.values) == 0 {
		return ""
	}

	allowedValues := make([]string, 0, len(a.values))
	for _, v := range a.values {
		allowedValues = append(allowedValues, hclutil.FormatCtyValueToString(v))
	}

	return strings.Join(allowedValues, ", ")
}

// AllowedValues creates a constraint that only allows the specified values.
//
// This constraint allows null values.  If a null value is not allowed, the field should be marked as required.
func AllowedValues(values ...cty.Value) TypeConstraint {
	return &allowedValuesConstraint{
		values: values,
	}
}

// ToTypeConstraint converts a protobuf AllowedValuesConstraintPB to a TypeConstraint instance.
func (a *AllowedValuesConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if a == nil {
		return nil, fmt.Errorf("AllowedValuesConstraintPB is nil")
	}

	values := make([]cty.Value, 0, len(a.Values))
	for _, vPB := range a.Values {
		v, err := json.Unmarshal(vPB, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value from protobuf: %w", err)
		}

		values = append(values, v)
	}

	return AllowedValues(values...), nil
}
