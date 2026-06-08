// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestMutuallyExclusiveValidate_Pass(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
		input  cty.Value
	}{
		{
			name:   "no fields present",
			fields: []string{"field1", "field2"},
			input:  cty.ObjectVal(map[string]cty.Value{}),
		},
		{
			name:   "one field present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
		},
		{
			name:   "one field present with null value",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.NullVal(cty.String),
			}),
		},
		{
			name:   "one field present with unknown value",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.UnknownVal(cty.String),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)
			err := constraint.Validate(tt.input)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestMutuallyExclusiveValidate_MultipleFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		input         cty.Value
		expectedError string
	}{
		{
			name:   "two fields present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			}),
			expectedError: `mutually exclusive fields "field1", "field2" are all present`,
		},
		{
			name:   "three fields present",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			}),
			expectedError: `mutually exclusive fields "field1", "field2", "field3" are all present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := MutuallyExclusive(tt.fields...)

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

func TestMutuallyExclusiveValidate_Nil(t *testing.T) {
	var constraint *mutuallyExclusiveConstraint

	err := constraint.Validate(cty.ObjectVal(map[string]cty.Value{}))
	expectedError := "mutually exclusive constraint is nil"
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
