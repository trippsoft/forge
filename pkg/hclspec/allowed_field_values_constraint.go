// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

type allowedFieldValuesConstraint struct {
	field  string      // The name of the field to check.
	values []cty.Value // The list of allowed values.
}

// Validate implements [TypeConstraint].
func (a *allowedFieldValuesConstraint) Validate(value cty.Value) error {
	if a == nil {
		return errors.New("allowed field values constraint is nil")
	}

	values := value.AsValueMap()

	value, ok := values[a.field]
	if !ok {
		return fmt.Errorf("field %q is not present", a.field)
	}

	if !value.IsWhollyKnown() {
		return errors.New("cannot validate unknown value")
	}

	if value.IsNull() {
		return nil // Null values are not validated
	}

	for _, allowed := range a.values {
		if value.Equals(allowed).True() {
			return nil
		}
	}

	return fmt.Errorf("field %q has an invalid value, allowed values are: %s", a.field, a.formatAllowedValues())
}

// ToProtobuf implements [TypeConstraint].
func (a *allowedFieldValuesConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if a == nil {
		return nil, errors.New("allowed field values constraint is nil")
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
		Constraint: &TypeConstraintPB_AllowedFieldValues{
			AllowedFieldValues: &AllowedFieldValuesConstraintPB{
				Field:  a.field,
				Values: values,
			},
		},
	}, nil
}

func (a *allowedFieldValuesConstraint) formatAllowedValues() string {
	if a == nil || len(a.values) == 0 {
		return ""
	}

	allowedValues := make([]string, 0, len(a.values))
	for _, v := range a.values {
		allowedValues = append(allowedValues, hclutil.FormatCtyValueToString(v))
	}

	return strings.Join(allowedValues, ", ")
}

// AllowedFieldValues creates a constraint that checks if a field's value is one of the allowed values.
//
// This should be used as part of a conditional constraint.
// Otherwise, the AllowedValues constraint should be placed on the field directly.
func AllowedFieldValues(field string, allowedValues ...cty.Value) TypeConstraint {
	return &allowedFieldValuesConstraint{
		field:  field,
		values: allowedValues,
	}
}

// ToTypeConstraint converts the protobuf AllowedFieldValuesConstraintPB to a TypeConstraint instance.
func (a *AllowedFieldValuesConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if a == nil {
		return nil, errors.New("AllowedFieldValuesConstraintPB is nil")
	}

	var values []cty.Value
	for _, vPB := range a.Values {
		v, err := json.Unmarshal(vPB, cty.DynamicPseudoType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert allowed value from protobuf: %w", err)
		}
		values = append(values, v)
	}

	return AllowedFieldValues(a.Field, values...), nil
}
