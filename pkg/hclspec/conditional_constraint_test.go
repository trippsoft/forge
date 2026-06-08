// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestConditionalConstraintValidate_Pass(t *testing.T) {
	tests := []struct {
		name       string
		condition  TypeCondition
		constraint TypeConstraint
		values     cty.Value
	}{
		{
			name:       "condition not met, constraint ignored",
			condition:  FieldPresent("field1"),
			constraint: RequireOneOf("field2", "field3"),
			values: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			}),
		},
		{
			name:       "condition met, constraint passes",
			condition:  FieldPresent("field1"),
			constraint: RequireOneOf("field2", "field3"),
			values: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.Validate(tt.values)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestConditionalConstraintValidate_ConditionMetConstraintFails(t *testing.T) {
	tests := []struct {
		name          string
		condition     TypeCondition
		constraint    TypeConstraint
		values        cty.Value
		expectedError string
	}{
		{
			name:       "condition met, constraint fails",
			condition:  FieldPresent("field1"),
			constraint: RequireOneOf("field2", "field3"),
			values: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
			expectedError: `conditional constraint failed: when field "field1" is present, at least one of the fields "field2", "field3" is required`,
		},
		{
			name:       "condition met, mutually exclusive constraint fails",
			condition:  FieldEquals("field1", cty.StringVal("trigger")),
			constraint: MutuallyExclusive("field2", "field3"),
			values: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("trigger"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			}),
			expectedError: `conditional constraint failed: when field "field1" is equal to "trigger", mutually exclusive fields "field2", "field3" are all present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := ConditionalConstraint(tt.condition, tt.constraint)

			err := constraint.Validate(tt.values)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}
