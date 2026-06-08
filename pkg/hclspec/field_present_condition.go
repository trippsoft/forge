// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type fieldPresentCondition struct {
	field string // The name of the field to check.
}

// IsMet implements [TypeCondition].
func (f *fieldPresentCondition) IsMet(value cty.Value) bool {
	if f == nil {
		return false
	}

	values := value.AsValueMap()

	if value, ok := values[f.field]; ok {
		return value.IsWhollyKnown() && !value.IsNull()
	}

	return false
}

// Description implements [TypeCondition].
func (f *fieldPresentCondition) Description() string {
	if f == nil {
		return ""
	}

	return fmt.Sprintf("field %q is present", f.field)
}

// ToProtobuf implements [TypeCondition].
func (f *fieldPresentCondition) ToProtobuf() (*TypeConditionPB, error) {
	if f == nil {
		return nil, errors.New("field present condition is nil")
	}

	return &TypeConditionPB{
		Condition: &TypeConditionPB_FieldPresent{
			FieldPresent: &FieldPresentConditionPB{
				Field: f.field,
			},
		},
	}, nil
}

// FieldPresent creates a condition that checks if a specified field is present.
func FieldPresent(field string) TypeCondition {
	return &fieldPresentCondition{
		field: field,
	}
}

// ToTypeCondition converts the protobuf FieldPresentConditionPB to a TypeCondition instance.
func (f *FieldPresentConditionPB) ToTypeCondition() (TypeCondition, error) {
	if f == nil {
		return nil, errors.New("FieldPresentConditionPB is nil")
	}

	return FieldPresent(f.Field), nil
}
