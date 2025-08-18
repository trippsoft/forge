// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/trippsoft/forge/pkg/errorwrap"
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
			field: RequiredField("test", String),
			value: cty.StringVal("hello"),
		},
		{
			name:  "valid number",
			field: RequiredField("test", Number),
			value: cty.NumberIntVal(123),
		},
		{
			name:  "null optional string",
			field: OptionalField("test", String),
			value: cty.NullVal(cty.String),
		},
		{
			name:  "valid string with constraint",
			field: RequiredField("test", String).WithConstraints(AllowedValues(cty.StringVal("allowed"))),
			value: cty.StringVal("allowed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.field.validate(tt.value)
			if err != nil {
				t.Errorf("expected no error from validate(), got: %v", err)
			}
		})
	}
}

func TestObjectFieldValidate_NoType(t *testing.T) {
	field := &objectField{name: "test"}
	err := field.validate(cty.StringVal("value"))
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := `field "test" has no type defined`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_NullRequired(t *testing.T) {
	field := RequiredField("test", String)
	err := field.validate(cty.NullVal(cty.String))
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := `missing required field "test"`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_InvalidValue(t *testing.T) {
	field := RequiredField("test", Duration)
	err := field.validate(cty.StringVal("not-a-duration"))
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := `field "test" validation failed: time: invalid duration "not-a-duration"`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_UnmetConstraint(t *testing.T) {
	field := RequiredField("test", String).WithConstraints(AllowedValues(cty.StringVal("allowed")))
	err := field.validate(cty.StringVal("not-allowed"))
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := `field "test" validation failed: value "not-allowed" is not in allowed values: "allowed"`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidate_Nil(t *testing.T) {
	var field *objectField
	err := field.validate(cty.NilVal)
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := "field is nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_Pass(t *testing.T) {

	tests := []struct {
		name  string
		field *objectField
	}{
		{
			name:  "valid required field",
			field: RequiredField("test", String),
		},
		{
			name:  "valid optional field with default",
			field: OptionalField("test", Number).WithDefaultValue(cty.NumberIntVal(42)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.field.validateSpec()
			if err != nil {
				t.Errorf("expected no errors from validateSpec(), got: %v", err)
			}
		})
	}
}

func TestObjectFieldValidateSpec_NoType(t *testing.T) {
	field := &objectField{name: "test"}
	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validate(), got none")
	}

	expectedError := `field "test" has no type defined`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validate(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_EmptyAlias(t *testing.T) {
	field := RequiredField("test", String).WithAliases("")
	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	expectedError := `field "test" has an empty alias`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_UnknownDefaultValue(t *testing.T) {

	field := OptionalField("test", String).WithDefaultValue(cty.UnknownVal(cty.String))

	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	expectedError := `field "test" has an unknown default value`
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_InvalidType(t *testing.T) {
	field := RequiredField("test", &primitiveType{t: cty.NilType})
	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	expectedError := "primitive type is nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_RequiredWithDefault(t *testing.T) {

	field := RequiredField("test", String).WithDefaultValue(cty.StringVal("default"))

	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	errs := errorwrap.UnwrapErrors(err)

	expectedError := `field "test" is required and has a default value`
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}

		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_ValueOfWrongType(t *testing.T) {

	field := OptionalField("test", Number).WithDefaultValue(cty.StringVal("not-a-number"))

	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	errs := errorwrap.UnwrapErrors(err)

	expectedError := `field "test" default value type mismatch: expected "number", got "string"`
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}

		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_InvalidValue(t *testing.T) {

	field := OptionalField("test", Duration).WithDefaultValue(cty.StringVal("not-a-duration"))

	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	errs := errorwrap.UnwrapErrors(err)

	expectedError := `field "test" default value validation failed: time: invalid duration "not-a-duration"`
	for _, e := range errs {
		if e.Error() == expectedError {
			return
		}

		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}

func TestObjectFieldValidateSpec_Nil(t *testing.T) {
	var field *objectField
	err := field.validateSpec()
	if err == nil {
		t.Fatal("expected error from validateSpec(), got none")
	}

	expectedError := "field is nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q from validateSpec(), got %q", expectedError, err.Error())
	}
}
