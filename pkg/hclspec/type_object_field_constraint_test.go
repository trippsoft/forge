// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFieldConstraintsValidate_Pass(t *testing.T) {
	tests := []struct {
		name        string
		constraints FieldConstraints
		value       cty.Value
	}{
		{
			name:        "nil",
			constraints: nil,
			value:       cty.NilVal,
		},
		{
			name: "valid value",
			constraints: FieldConstraints{
				AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
			},
			value: cty.StringVal("value1"),
		},
		{
			name:        "empty constraints",
			constraints: FieldConstraints{},
			value:       cty.NilVal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constraints.Validate(tt.value)
			if err != nil {
				t.Errorf("expected no error from Validate(), got %q", err.Error())
			}
		})
	}
}

func TestFieldConstraintsValidate_Empty(t *testing.T) {
	constraints := FieldConstraints{}

	err := constraints.Validate(cty.NilVal)
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestFieldConstraintsValidate_Fail(t *testing.T) {
	constraints := FieldConstraints{
		AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
	}

	expectedError := `value "value3" is not in allowed values: "value1", "value2"`
	err := constraints.Validate(cty.StringVal("value3"))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestFieldConstraintsValidate_Nil(t *testing.T) {
	var constraints FieldConstraints

	err := constraints.Validate(cty.NilVal)
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestFieldConstraintsValidateSpec_Pass(t *testing.T) {
	constraints := FieldConstraints{
		AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
	}

	field := OptionalField("test", String)

	err := constraints.ValidateSpec(field)
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestFieldConstraintsValidateSpec_Fail(t *testing.T) {
	constraints := FieldConstraints{AllowedValues()}

	field := OptionalField("test", String)

	expectedError := "allowed values constraint has no values defined"
	err := constraints.ValidateSpec(field)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestFieldConstraintsValidateSpec_NilField(t *testing.T) {
	constraints := FieldConstraints{
		AllowedValues(cty.StringVal("value1"), cty.StringVal("value2")),
	}

	expectedError := "field is nil"
	err := constraints.ValidateSpec(nil)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestFieldConstraintsValidateSpec_Empty(t *testing.T) {
	constraints := FieldConstraints{}

	field := OptionalField("test", String)

	err := constraints.ValidateSpec(field)
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestFieldConstraintsValidateSpec_Nil(t *testing.T) {
	var constraints FieldConstraints

	field := OptionalField("test", String)

	err := constraints.ValidateSpec(field)
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestAllowedValuesValidate_Pass(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraints.Validate(cty.StringVal("value1"))
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestAllowedValuesValidate_NullValue(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	err := constraints.Validate(cty.NullVal(cty.String))
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestAllowedValuesValidate_Fail(t *testing.T) {
	constraints := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	expectedError := `value "value3" is not in allowed values: "value1", "value2"`
	err := constraints.Validate(cty.StringVal("value3"))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidate_Nil(t *testing.T) {
	var constraint *allowedValuesConstraint

	expectedError := "allowed values constraint is nil"
	err := constraint.Validate(cty.NilVal)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidateSpec_Pass(t *testing.T) {
	constraint := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	field := OptionalField("test", String)

	err := constraint.ValidateSpec(field)
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestAllowedValuesValidateSpec_EmptyValues(t *testing.T) {
	constraint := AllowedValues()

	expectedError := "allowed values constraint has no values defined"
	err := constraint.ValidateSpec(nil)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidateSpec_NilField(t *testing.T) {
	constraint := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))

	expectedError := "field is nil"
	err := constraint.ValidateSpec(nil)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidateSpec_InvalidFieldType(t *testing.T) {
	constraint := AllowedValues(cty.StringVal("value1"), cty.StringVal("value2"))
	field := &objectField{t: &primitiveType{t: cty.NilType}}

	err := constraint.ValidateSpec(field)
	if err != nil {
		t.Fatalf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestAllowedValuesValidateSpec_ValueOfWrongType(t *testing.T) {
	constraint := AllowedValues(cty.StringVal("value1"))
	field := OptionalField("test", Number)

	expectedError := `allowed value cty.StringVal("value1") does not match field type number`
	err := constraint.ValidateSpec(field)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidateSpec_InvalidValue(t *testing.T) {
	constraint := AllowedValues(cty.StringVal("not-a-duration"))
	field := OptionalField("test", Duration)

	expectedError := `allowed value cty.StringVal("not-a-duration") is invalid: field "test" validation failed: time: invalid duration "not-a-duration"`
	err := constraint.ValidateSpec(field)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestAllowedValuesValidateSpec_Nil(t *testing.T) {
	var constraint *allowedValuesConstraint

	expectedError := "allowed values constraint is nil"
	err := constraint.ValidateSpec(nil)
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
