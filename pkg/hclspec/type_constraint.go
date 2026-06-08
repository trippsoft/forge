// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"

	"github.com/zclconf/go-cty/cty"
)

// TypeConstraint represents a constraint on a Type.
type TypeConstraint interface {
	// Validate checks if the given value satisfies the constraint.
	//
	// Implementations are expected to return all errors produced by the validation, if possible.
	Validate(value cty.Value) error

	// ToProtobuf converts the TypeConstraint to its protobuf representation.
	ToProtobuf() (*TypeConstraintPB, error)
}

// TypeConstraints is a slice of TypeConstraint that can be used to validate a value against multiple constraints.
type TypeConstraints []TypeConstraint

// Validate checks if the given value satisfies all constraints.
func (c TypeConstraints) Validate(value cty.Value) error {
	var err error
	for _, constraint := range c {
		e := constraint.Validate(value)
		err = errors.Join(err, e)
	}

	return err
}

// ToProtobuf converts the TypeConstraints to their protobuf representation.
func (c TypeConstraints) ToProtobuf() ([]*TypeConstraintPB, error) {
	if len(c) == 0 {
		return nil, nil
	}

	constraints := make([]*TypeConstraintPB, 0, len(c))
	var err error
	for _, constraint := range c {
		pb, e := constraint.ToProtobuf()
		if e != nil {
			err = errors.Join(err, e)
			continue
		}

		constraints = append(constraints, pb)
	}

	if err != nil {
		return nil, err
	}

	return constraints, nil
}

// ToTypeConstraint converts a TypeConstraintPB to its TypeConstraint representation.
func (c *TypeConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	switch constraint := c.GetConstraint().(type) {
	case *TypeConstraintPB_AllowedValues:
		return constraint.AllowedValues.ToTypeConstraint()
	case *TypeConstraintPB_MutuallyExclusive:
		return constraint.MutuallyExclusive.ToTypeConstraint()
	case *TypeConstraintPB_RequiredTogether:
		return constraint.RequiredTogether.ToTypeConstraint()
	case *TypeConstraintPB_RequireOneOf:
		return constraint.RequireOneOf.ToTypeConstraint()
	case *TypeConstraintPB_AllowedFieldValues:
		return constraint.AllowedFieldValues.ToTypeConstraint()
	case *TypeConstraintPB_Conditional:
		return constraint.Conditional.ToTypeConstraint()
	}

	return nil, errors.New("unsupported constraint type")
}
