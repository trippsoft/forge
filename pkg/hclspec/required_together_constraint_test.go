// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestRequiredTogetherValidate_Pass(t *testing.T) {
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
			name:   "all fields present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
			}),
		},
		{
			name:   "all three fields present",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.StringVal("value3"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.Validate(tt.input)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidate_PartialFieldsPresent(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		input         cty.Value
		expectedError string
	}{
		{
			name:   "only first field present",
			fields: []string{"field1", "field2"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
			expectedError: `fields "field1", "field2" are required together, but only "field1" is present`,
		},
		{
			name:   "only last field present",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field3": cty.StringVal("value3"),
			}),
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field3" is present`,
		},
		{
			name:   "some fields present",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field3": cty.StringVal("value3"),
			}),
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field3" is present`,
		},
		{
			name:   "some fields present with null field ignored",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.NullVal(cty.String),
			}),
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field2" is present`,
		},
		{
			name:   "some fields present with unknown field ignored",
			fields: []string{"field1", "field2", "field3"},
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
				"field2": cty.StringVal("value2"),
				"field3": cty.UnknownVal(cty.String),
			}),
			expectedError: `fields "field1", "field2", "field3" are required together, but only "field1", "field2" is present`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := RequiredTogether(tt.fields...)

			err := constraint.Validate(tt.input)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q	from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestRequiredTogetherValidate_Nil(t *testing.T) {
	var constraint *requiredTogetherConstraint

	expectedError := "required together constraint is nil"
	err := constraint.Validate(cty.ObjectVal(map[string]cty.Value{}))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
