// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestObjectFieldValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		field *objectField
		value cty.Value
	}{
		{
			name:  "valid string",
			field: RequiredField("test", String()),
			value: cty.StringVal("hello"),
		},
		{
			name:  "valid number",
			field: RequiredField("test", Number()),
			value: cty.NumberIntVal(123),
		},
		{
			name:  "null optional string",
			field: OptionalField("test", String()),
			value: cty.NullVal(cty.String),
		},
		{
			name:  "valid string with constraint",
			field: RequiredField("test", String().WithConstraints(AllowedValues(cty.StringVal("allowed")))),
			value: cty.StringVal("allowed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.field.Validate(tt.value)
			if err != nil {
				t.Errorf("expected no error from validate(), got %q", err.Error())
			}
		})
	}
}

func TestObjectFieldValidate_NoType(t *testing.T) {
	field := &objectField{name: "test"}

	expectedError := `field "test" has no type defined`
	err := field.Validate(cty.StringVal("value"))
	if err == nil {
		t.Fatalf("expected error %q from validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_NullRequired(t *testing.T) {
	field := RequiredField("test", String())

	expectedError := `missing required field "test"`
	err := field.Validate(cty.NullVal(cty.String))
	if err == nil {
		t.Fatalf("expected error %q from validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_UnmetConstraint(t *testing.T) {
	field := RequiredField("test", String().WithConstraints(AllowedValues(cty.StringVal("allowed"))))

	expectedError := `field "test" validation failed: value "not-allowed" is not in allowed values: "allowed"`
	err := field.Validate(cty.StringVal("not-allowed"))
	if err == nil {
		t.Fatalf("expected error %q from validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_Nil(t *testing.T) {
	var field *objectField

	expectedError := "field is nil"
	err := field.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}
