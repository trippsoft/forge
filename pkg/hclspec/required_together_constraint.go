// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type requiredTogetherConstraint struct {
	fields []string
}

// Validate implements [TypeConstraint].
func (r *requiredTogetherConstraint) Validate(value cty.Value) error {
	if r == nil {
		return errors.New("required together constraint is nil")
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

	if len(foundFields) > 0 && len(foundFields) != len(r.fields) {
		fieldNames := ""
		for i, field := range foundFields {
			if i > 0 {
				fieldNames += ", "
			}

			fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
		}

		return fmt.Errorf("fields %s are required together, but only %s is present", r.formatFieldNames(), fieldNames)
	}

	return nil
}

// ToProtobuf implements [TypeConstraint].
func (r *requiredTogetherConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if r == nil {
		return nil, errors.New("required together constraint is nil")
	}

	return &TypeConstraintPB{
		Constraint: &TypeConstraintPB_RequiredTogether{
			RequiredTogether: &RequiredTogetherConstraintPB{
				Fields: r.fields,
			},
		},
	}, nil
}

func (r *requiredTogetherConstraint) formatFieldNames() any {
	if r == nil {
		return ""
	}

	fieldNames := ""
	for i, field := range r.fields {
		if i > 0 {
			fieldNames += ", "
		}

		fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
	}

	return fieldNames
}

// RequiredTogether creates a constraint requiring the specified fields to be present together.
func RequiredTogether(fields ...string) TypeConstraint {
	return &requiredTogetherConstraint{
		fields: fields,
	}
}

// ToTypeConstraint converts the protobuf RequiredTogetherConstraintPB to a TypeConstraint instance.
func (r *RequiredTogetherConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if r == nil {
		return nil, errors.New("RequiredTogetherConstraintPB is nil")
	}

	return RequiredTogether(r.Fields...), nil
}
