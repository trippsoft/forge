// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type conditionalConstraint struct {
	condition  TypeCondition
	constraint TypeConstraint
}

// Validate implements [TypeConstraint].
func (c *conditionalConstraint) Validate(value cty.Value) error {
	if c == nil {
		return errors.New("conditional constraint is nil")
	}

	if c.condition == nil {
		return nil // No condition means the constraint is never applied
	}

	if c.constraint == nil {
		return nil // No constraint means nothing to validate
	}

	if c.condition.IsMet(value) {
		if err := c.constraint.Validate(value); err != nil {
			return fmt.Errorf("conditional constraint failed: when %s, %w", c.condition.Description(), err)
		}
	}

	return nil
}

// ToProtobuf implements [TypeConstraint].
func (c *conditionalConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if c == nil {
		return nil, errors.New("conditional constraint is nil")
	}

	conditionPB, err := c.condition.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert condition to protobuf: %w", err)
	}

	constraintPB, err := c.constraint.ToProtobuf()
	if err != nil {
		return nil, fmt.Errorf("failed to convert constraint to protobuf: %w", err)
	}

	return &TypeConstraintPB{
		Constraint: &TypeConstraintPB_Conditional{
			Conditional: &ConditionalConstraintPB{
				Condition:  conditionPB,
				Constraint: constraintPB,
			},
		},
	}, nil
}

// ConditionalConstraint creates a constraint that only applies if the specified condition is met.
func ConditionalConstraint(condition TypeCondition, constraint TypeConstraint) TypeConstraint {
	return &conditionalConstraint{
		condition:  condition,
		constraint: constraint,
	}
}

// ToTypeConstraint converts the protobuf ConditionalConstraintPB to a TypeConstraint instance.
func (c *ConditionalConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if c == nil {
		return nil, errors.New("ConditionalConstraintPB is nil")
	}

	condition, err := c.Condition.ToTypeCondition()
	if err != nil {
		return nil, fmt.Errorf("failed to convert condition from protobuf: %w", err)
	}

	constraint, err := c.Constraint.ToTypeConstraint()
	if err != nil {
		return nil, fmt.Errorf("failed to convert constraint from protobuf: %w", err)
	}

	return ConditionalConstraint(condition, constraint), nil
}
