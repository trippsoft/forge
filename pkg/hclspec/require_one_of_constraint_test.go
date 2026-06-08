// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestRequireOneOfValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		input  cty.Value
	}{
		{
			name:   "first field present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
		},
		{
			name:   "last field present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name:   "multiple fields present",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name:   "field present with null values ignored",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name:   "field present with unknown values ignored",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.UnknownVal(cty.String),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequireOneOf(tt.fields...)

			err := constraint.Validate(tt.input)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestRequireOneOfValidate_NoFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		input         cty.Value
		expectedError string
	}{
		{
			name:          "no fields in values",
			fields:        []string{"field1", "field2"},
			input:         cty.ObjectVal(map[string]cty.Value{}),
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
		{
			name:   "only null values",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
				"field2": cty.NullVal(cty.String),
			}),
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
		{
			name:   "only unknown values",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
				"field2": cty.UnknownVal(cty.String),
			}),
			expectedError: `at least one of the fields "field1", "field2" is required`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequireOneOf(tt.fields...)

			err := constraint.Validate(tt.input)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequireOneOfValidate_Nil(t *testing.T) {
	var constraint *requireOneOfConstraint

	expectedError := "required one of constraint is nil"
	err := constraint.Validate(cty.ObjectVal(map[string]cty.Value{}))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
