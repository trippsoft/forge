// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
)

// TypeCondition is used by constraints that require a specific condition to be met to apply.
type TypeCondition interface {
	// IsMet checks if the condition is satisfied by the given value.
	//
	// Note that no error is returned from this method.
	// If there is an error, it means the condition is not met.
	IsMet(value cty.Value) bool

	// Description provides a human-readable description of the condition for error messages.
	Description() string

	// ToProtobuf converts the TypeCondition to its protobuf representation.
	ToProtobuf() (*TypeConditionPB, error)
}

// ToTypeCondition converts the protobuf TypeConditionPB to a TypeCondition.
func (o *TypeConditionPB) ToTypeCondition() (TypeCondition, error) {
	if o == nil {
		return nil, errors.New("TypeConditionPB is nil")
	}

	switch condition := o.Condition.(type) {
	case *TypeConditionPB_FieldPresent:
		return condition.FieldPresent.ToTypeCondition()
	case *TypeConditionPB_FieldNotPresent:
		return condition.FieldNotPresent.ToTypeCondition()
	case *TypeConditionPB_FieldEquals:
		return condition.FieldEquals.ToTypeCondition()
	}

	return nil, errors.New("unknown TypeConditionPB type")
}
