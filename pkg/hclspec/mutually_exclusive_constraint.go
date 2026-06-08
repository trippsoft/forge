// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"errors"
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

type mutuallyExclusiveConstraint struct {
	fields []string // List of field names that are mutually exclusive.
}

// Validate implements [TypeConstraint].
func (m *mutuallyExclusiveConstraint) Validate(value cty.Value) error {
	if m == nil {
		return errors.New("mutually exclusive constraint is nil")
	}

	values := value.AsValueMap()

	var foundFields []string
	for _, field := range m.fields {
		if value, ok := values[field]; ok {
			if !value.IsWhollyKnown() || value.IsNull() {
				continue
			}

			foundFields = append(foundFields, field)
		}
	}

	fieldNames := ""
	for i, field := range m.fields {
		if i > 0 {
			fieldNames += ", "
		}

		fieldNames = fmt.Sprintf("%s%q", fieldNames, field)
	}

	if len(foundFields) > 1 {
		return fmt.Errorf("mutually exclusive fields %s are all present", fieldNames)
	}

	return nil
}

// ToProtobuf implements [TypeConstraint].
func (m *mutuallyExclusiveConstraint) ToProtobuf() (*TypeConstraintPB, error) {
	if m == nil {
		return nil, errors.New("mutually exclusive constraint is nil")
	}

	return &TypeConstraintPB{
		Constraint: &TypeConstraintPB_MutuallyExclusive{
			MutuallyExclusive: &MutuallyExclusiveConstraintPB{
				Fields: m.fields,
			},
		},
	}, nil
}

// MutuallyExclusive creates a constraint requiring the specified fields to be mutually exclusive.
func MutuallyExclusive(fields ...string) TypeConstraint {
	return &mutuallyExclusiveConstraint{
		fields: fields,
	}
}

// ToTypeConstraint converts the protobuf MutuallyExclusiveConstraintPB to a TypeConstraint instance.
func (m *MutuallyExclusiveConstraintPB) ToTypeConstraint() (TypeConstraint, error) {
	if m == nil {
		return nil, errors.New("MutuallyExclusiveConstraintPB is nil")
	}

	return MutuallyExclusive(m.Fields...), nil
}
