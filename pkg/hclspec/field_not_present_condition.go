// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type fieldNotPresentCondition struct {
	field string // The name of the field to check.
}

// IsMet implements [TypeCondition].
func (f *fieldNotPresentCondition) IsMet(value cty.Value) bool {
	if f == nil {
		return false
	}

	values := value.AsValueMap()

	if value, ok := values[f.field]; ok {
		return value.IsWhollyKnown() && value.IsNull()
	}

	return false
}

// Description implements [TypeCondition].
func (f *fieldNotPresentCondition) Description() string {
	if f == nil {
		return ""
	}

	return fmt.Sprintf("field %q is not present", f.field)
}

// ToProtobuf implements [TypeCondition].
func (f *fieldNotPresentCondition) ToProtobuf() (*TypeConditionPB, error) {
	if f == nil {
		return nil, errors.New("field not present condition is nil")
	}

	return &TypeConditionPB{
		Condition: &TypeConditionPB_FieldNotPresent{
			FieldNotPresent: &FieldNotPresentConditionPB{
				Field: f.field,
			},
		},
	}, nil
}

// FieldNotPresent creates a condition that checks if a specified field is not present.
func FieldNotPresent(field string) TypeCondition {
	return &fieldNotPresentCondition{
		field: field,
	}
}

// ToTypeCondition converts the protobuf FieldNotPresentConditionPB to a TypeCondition instance.
func (f *FieldNotPresentConditionPB) ToTypeCondition() (TypeCondition, error) {
	if f == nil {
		return nil, errors.New("FieldNotPresentConditionPB is nil")
	}

	return FieldNotPresent(f.Field), nil
}
