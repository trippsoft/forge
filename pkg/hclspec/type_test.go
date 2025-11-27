// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func TestStringCtyType(t *testing.T) {
	expected := cty.String
	actual := String.CtyType()
	if !actual.Equals(expected) {
		t.Errorf("expected %q from CtyType(), got %q", expected.FriendlyName(), actual.FriendlyName())
	}
}

func TestStringConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid string",
			input:    cty.StringVal("hello"),
			expected: cty.StringVal("hello"),
		},
		{
			name:     "number to string conversion",
			input:    cty.NumberIntVal(123),
			expected: cty.StringVal("123"),
		},
		{
			name:     "bool to string conversion",
			input:    cty.BoolVal(true),
			expected: cty.StringVal("true"),
		},
		{
			name:     "null string",
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.String),
		},
		{
			name:     "null number",
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.String),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, String, tt.input, tt.expected)
		})
	}
}

func TestStringConvert_UnknownValue(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "unknown string",
			input: cty.UnknownVal(cty.String),
		},
		{
			name:  "unknown number",
			input: cty.UnknownVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedConversion(t, String, tt.input, "cannot convert unknown value")
		})
	}
}

func TestStringConvert_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name: "list of strings",
			input: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "map of strings",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name: "map of numbers",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
		},
		{
			name: "map of booleans",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
		},
		{
			name: "tuple of strings",
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "tuple of numbers",
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			name: "object",
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := fmt.Sprintf(
				"cannot convert %q to %q",
				tt.input.Type().FriendlyName(),
				cty.String.FriendlyName())

			verifyFailedConversion(t, String, tt.input, expectedError)
		})
	}
}

func TestStringValidate_Pass(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:  "valid string",
			input: cty.StringVal("hello"),
		},
		{
			name:  "number to string conversion",
			input: cty.NumberIntVal(123),
		},
		{
			name:  "bool to string conversion",
			input: cty.BoolVal(true),
		},
		{
			name:  "null string",
			input: cty.NullVal(cty.String),
		},
		{
			name:  "null number",
			input: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, String, tt.input)
		})
	}
}

func TestStringValidateSpec_Pass(t *testing.T) {
	err := String.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestNumberCtyType(t *testing.T) {
	expected := cty.Number
	actual := Number.CtyType()
	if !actual.Equals(expected) {
		t.Errorf("expected %q from CtyType(), got %q",
			expected.FriendlyName(),
			actual.FriendlyName())
	}
}

func TestNumberConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid integer",
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid float",
			input:    cty.NumberFloatVal(3.14),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "string integer conversion",
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string float conversion",
			input:    cty.StringVal("3.14"),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "null number",
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, Number, tt.input, tt.expected)
		})
	}
}

func TestNumberConvert_UnknownValue(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "unknown number",
			input: cty.UnknownVal(cty.Number),
		},
		{
			name:  "unknown string",
			input: cty.UnknownVal(cty.String),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedConversion(t, Number, tt.input, "cannot convert unknown value")
		})
	}
}

func TestNumberConvert_InvalidValues(t *testing.T) {
	tests := []struct {
		name            string
		input           cty.Value
		conversionError string
	}{
		{
			name:            "invalid string",
			input:           cty.StringVal("not-a-number"),
			conversionError: "a number is required",
		},
		{
			name: "list of strings",
			input: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "map of strings",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name: "map of numbers",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
		},
		{
			name: "map of booleans",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
		},
		{
			name: "tuple of strings",
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "tuple of numbers",
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			name: "object",
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedError string
			if tt.conversionError != "" {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					Number.CtyType().FriendlyName(),
					tt.conversionError)
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					Number.CtyType().FriendlyName())
			}

			verifyFailedConversion(t, Number, tt.input, expectedError)
		})
	}
}

func TestNumberValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "valid integer",
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid float",
			input: cty.NumberFloatVal(3.14),
		},
		{
			name:  "string integer conversion",
			input: cty.StringVal("42"),
		},
		{
			name:  "string float conversion",
			input: cty.StringVal("3.14"),
		},
		{
			name:  "null number",
			input: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, Number, tt.input)
		})
	}
}

func TestNumberValidateSpec_Pass(t *testing.T) {
	err := Number.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err)
	}
}

func TestBoolCtyType(t *testing.T) {
	expected := cty.Bool
	actual := Bool.CtyType()
	if !actual.Equals(expected) {
		t.Errorf("expected %q from CtyType(), got %q",
			expected.FriendlyName(),
			actual.FriendlyName(),
		)
	}
}

func TestBoolConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid true",
			input:    cty.BoolVal(true),
			expected: cty.BoolVal(true),
		},
		{
			name:     "valid false",
			input:    cty.BoolVal(false),
			expected: cty.BoolVal(false),
		},
		{
			name:     "string true conversion",
			input:    cty.StringVal("true"),
			expected: cty.BoolVal(true),
		},
		{
			name:     "string false conversion",
			input:    cty.StringVal("false"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "null bool",
			input:    cty.NullVal(cty.Bool),
			expected: cty.NullVal(cty.Bool),
		},
		{
			name:     "null string",
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.Bool),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, Bool, tt.input, tt.expected)
		})
	}
}

func TestBoolConvert_UnknownValue(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "unknown bool",
			input: cty.UnknownVal(cty.Bool),
		},
		{
			name:  "unknown string",
			input: cty.UnknownVal(cty.String),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedConversion(t, Bool, tt.input, "cannot convert unknown value")
		})
	}
}

func TestBoolConvert_InvalidValues(t *testing.T) {
	tests := []struct {
		name            string
		input           cty.Value
		conversionError string
	}{
		{
			name:            "invalid string",
			input:           cty.StringVal("not-a-bool"),
			conversionError: "a bool is required",
		},
		{
			name:  "number",
			input: cty.NumberIntVal(1),
		},
		{
			name: "list of strings",
			input: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "map of strings",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name: "map of numbers",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
		},
		{
			name: "map of booleans",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
		},
		{
			name: "tuple of strings",
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "tuple of numbers",
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			name: "object",
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedError string
			if tt.conversionError != "" {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					Bool.CtyType().FriendlyName(),
					tt.conversionError)
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					Bool.CtyType().FriendlyName())
			}

			verifyFailedConversion(t, Bool, tt.input, expectedError)
		})
	}
}

func TestBoolValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "valid true",
			input: cty.BoolVal(true),
		},
		{
			name:  "valid false",
			input: cty.BoolVal(false),
		},
		{
			name:  "string true conversion",
			input: cty.StringVal("true"),
		},
		{
			name:  "string false conversion",
			input: cty.StringVal("false"),
		},
		{
			name:  "null bool",
			input: cty.NullVal(cty.Bool),
		},
		{
			name:  "null string",
			input: cty.NullVal(cty.String),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, Bool, tt.input)
		})
	}
}

func TestBoolValidateSpec_Pass(t *testing.T) {
	err := Bool.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestPrimitiveTypeCtyType_Nil(t *testing.T) {
	var primitiveType *primitiveType

	actual := primitiveType.CtyType()
	if !actual.Equals(cty.NilType) {
		t.Errorf("expected %q from CtyType(), got %q", cty.NilType.FriendlyName(), actual.FriendlyName())
	}
}

func TestPrimitiveTypeConvert_Nil(t *testing.T) {
	var primitiveType *primitiveType

	expectedError := "primitive type is nil"
	converted, err := primitiveType.Convert(cty.StringVal("test"))
	if err == nil {
		t.Errorf("expected error %q from Convert(), got none", expectedError)
	}

	if !converted.Equals(cty.NilVal).True() {
		t.Errorf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func TestPrimitiveTypeValidate_Nil(t *testing.T) {
	var primitiveType *primitiveType

	err := primitiveType.Validate(cty.StringVal("test"))
	if err != nil {
		t.Errorf("expected no error from Validate(), got %q", err.Error())
	}
}

func TestPrimitiveTypeValidateSpec_Nil(t *testing.T) {
	var primitiveType *primitiveType

	expectedError := "primitive type is nil"
	err := primitiveType.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestPrimitiveTypeValidateSpec_NilType(t *testing.T) {
	primitiveType := &primitiveType{}

	expectedError := "primitive type is nil"
	err := primitiveType.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestPrimitiveTypeValidateSpec_NonPrimitive(t *testing.T) {
	primitiveType := &primitiveType{cty.List(cty.String)}

	expectedError := `Type "list of string" is not a primitive type`
	err := primitiveType.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestRawCtyType(t *testing.T) {
	expected := cty.DynamicPseudoType
	actual := Raw.CtyType()
	if !actual.Equals(expected) {
		t.Errorf("expected %q from CtyType(), got %q", expected.FriendlyName(), actual.FriendlyName())
	}
}

func TestRawConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid string",
			input:    cty.StringVal("hello"),
			expected: cty.StringVal("hello"),
		},
		{
			name:     "number to string conversion",
			input:    cty.NumberIntVal(123),
			expected: cty.NumberIntVal(123),
		},
		{
			name:     "bool to string conversion",
			input:    cty.BoolVal(true),
			expected: cty.BoolVal(true),
		},
		{
			name:     "null string",
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.String),
		},
		{
			name:     "null number",
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name: "list of strings",
			input: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "map of strings",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name: "map of numbers",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
		},
		{
			name: "map of booleans",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
		},
		{
			name: "tuple of strings",
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
			expected: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "tuple of numbers",
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
			expected: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			name: "object",
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
			expected: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, Raw, tt.input, tt.expected)
		})
	}
}

func TestRawConvert_UnknownValue(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "unknown string",
			input: cty.UnknownVal(cty.String),
		},
		{
			name:  "unknown number",
			input: cty.UnknownVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedConversion(t, Raw, tt.input, "cannot convert unknown value")
		})
	}
}

func TestRawValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "valid string",
			input: cty.StringVal("hello"),
		},
		{
			name:  "number to string conversion",
			input: cty.NumberIntVal(123),
		},
		{
			name:  "bool to string conversion",
			input: cty.BoolVal(true),
		},
		{
			name:  "null string",
			input: cty.NullVal(cty.String),
		},
		{
			name:  "null number",
			input: cty.NullVal(cty.Number),
		},
		{
			name: "list of strings",
			input: cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "map of strings",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name: "map of numbers",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
		},
		{
			name: "map of booleans",
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.BoolVal(true),
				"key2": cty.BoolVal(false),
			}),
		},
		{
			name: "tuple of strings",
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.StringVal("world"),
			}),
		},
		{
			name: "tuple of numbers",
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			name: "object",
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, Raw, tt.input)
		})
	}
}

func TestRawValidateSpec_Pass(t *testing.T) {
	err := Raw.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestRawValidateSpec_Nil(t *testing.T) {
	var rawType *rawType

	expectedError := "raw type is nil"
	err := rawType.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
