// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestAllowedFieldValuesValidate_Pass(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		allowedValues []cty.Value
		value         cty.Value
	}{
		{
			name:          "field is null",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value:         cty.ObjectVal(map[string]cty.Value{"field1": cty.NullVal(cty.String)}),
		},
		{
			name:          "field is first value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value:         cty.ObjectVal(map[string]cty.Value{"field1": cty.StringVal("value1")}),
		},
		{
			name:          "field is second value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value:         cty.ObjectVal(map[string]cty.Value{"field1": cty.StringVal("value2")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := AllowedFieldValues(tt.field, tt.allowedValues...)

			err := constraint.Validate(tt.value)
			if err != nil {
				t.Fatalf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestAllowedFieldValuesValidate_FieldNotPresent(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		allowedValues []cty.Value
		value         cty.Value
		expectedError string
	}{
		{
			name:          "no field in values",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value:         cty.ObjectVal(map[string]cty.Value{}),
			expectedError: `field "field1" is not present`,
		},
		{
			name:          "not allowed value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value3"),
				"field2": cty.NullVal(cty.String),
			}),
			expectedError: `field "field1" has an invalid value, allowed values are: "value1", "value2"`,
		},
		{
			name:          "unknown value",
			field:         "field1",
			allowedValues: []cty.Value{cty.StringVal("value1"), cty.StringVal("value2")},
			value: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
				"field2": cty.UnknownVal(cty.String),
			}),
			expectedError: `cannot validate unknown value`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := AllowedFieldValues(tt.field, tt.allowedValues...)

			err := constraint.Validate(tt.value)
			if err == nil {
				t.Fatalf("expected error %q from Validate(), got none", tt.expectedError)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q from Validate(), got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestAllowedFieldValuesValidate_Nil(t *testing.T) {
	var constraint *allowedFieldValuesConstraint

	expectedError := "allowed field values constraint is nil"
	err := constraint.Validate(cty.EmptyObjectVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
