// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

type fieldEqualsCondition struct {
	field string    // The name of the field to check.
	value cty.Value // The value to compare against.
}

// IsMet implements [TypeCondition].
func (f *fieldEqualsCondition) IsMet(value cty.Value) bool {
	if f == nil {
		return false
	}

	values := value.AsValueMap()

	if value, ok := values[f.field]; ok {
		return value.IsWhollyKnown() && value.Equals(f.value).True()
	}

	return false
}

// Description implements [TypeCondition].
func (f *fieldEqualsCondition) Description() string {
	if f == nil {
		return ""
	}

	return fmt.Sprintf("field %q is equal to %s", f.field, hclutil.FormatCtyValueToString(f.value))
}

// ToProtobuf implements [TypeCondition].
func (f *fieldEqualsCondition) ToProtobuf() (*TypeConditionPB, error) {
	if f == nil {
		return nil, errors.New("field equals condition is nil")
	}

	valuePB, err := json.Marshal(f.value, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field equals value to protobuf: %w", err)
	}

	return &TypeConditionPB{
		Condition: &TypeConditionPB_FieldEquals{
			FieldEquals: &FieldEqualsConditionPB{
				Field: f.field,
				Value: valuePB,
			},
		},
	}, nil
}

// FieldEquals creates a condition that checks if a specified field has a specified value.
func FieldEquals(field string, value cty.Value) TypeCondition {
	return &fieldEqualsCondition{
		field: field,
		value: value,
	}
}

// ToTypeCondition converts the protobuf FieldEqualsConditionPB to a TypeCondition instance.
func (f *FieldEqualsConditionPB) ToTypeCondition() (TypeCondition, error) {
	if f == nil {
		return nil, errors.New("FieldEqualsConditionPB is nil")
	}

	value, err := json.Unmarshal(f.Value, cty.DynamicPseudoType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert field equals value from protobuf: %w", err)
	}

	return FieldEquals(f.Field, value), nil
}
