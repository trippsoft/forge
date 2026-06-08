// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type requireOneOfConstraint struct {
	fields []string // List of field names of which at least one is required.
}

// Validate implements [TypeConstraint].
func (r *requireOneOfConstraint) Validate(value cty.Value) error {
	if r == nil {
		return errors.New("required one of constraint is nil")
	}

	values := value.AsValueMap()

	var foundFields []string
	for _, field := range r.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}

			foundFields = append(foundFields, field)
		}
	}

	if len(foundFields) == 0 {
		fieldNames := ""
		for i, fieldName := range r.fields {
			if i > 0 {
				fieldNames += ", "
			}

			fieldNames = fmt.Sprintf("%s%q", fieldNames, fieldName)
		}

		return fmt.Errorf("at least one of the fields %s is required", fieldNames)
	}

	return nil
}

// ToProtobuf implements [TypeConstraint].
func (r *requireOneOfConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if r == nil {
		return nil, errors.New("required one of constraint is nil")
	}

	return &TypeConstraintPB{
		Constraint: &TypeConstraintPB_RequireOneOf{
			RequireOneOf: &RequireOneOfConstraintPB{
				Fields: r.fields,
			},
		},
	}, nil
}

// RequireOneOf creates a constraint requiring one of the specified fields to be present.
func RequireOneOf(fields ...string) TypeConstraint {
	return &requireOneOfConstraint{
		fields: fields,
	}
}

// ToTypeConstraint converts the protobuf RequireOneOfConstraintPB to a TypeConstraint instance.
func (r *RequireOneOfConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if r == nil {
		return nil, errors.New("RequireOneOfConstraintPB is nil")
	}

	return RequireOneOf(r.Fields...), nil
}
