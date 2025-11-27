// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func TestSensitiveStringCtyType(t *testing.T) {
	expected := cty.String
	actual := SensitiveString.CtyType()
	if actual != expected {
		t.Errorf("expected %q from CtyType(), got %q", expected.FriendlyName(), actual.FriendlyName())
	}
}

func TestSensitiveStringCtyType_Nil(t *testing.T) {
	var sensitiveString *sensitiveStringType

	expected := cty.String
	actual := sensitiveString.CtyType()
	if actual != expected {
		t.Errorf("expected %q from CtyType(), got %q", expected.FriendlyName(), actual.FriendlyName())
	}
}

func TestSensitiveStringConvert_Success(t *testing.T) {
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
			ui.SecretFilter.Clear()

			verifySuccessfulConversion(t, SensitiveString, tt.input, tt.expected)

			secrets := ui.SecretFilter.Secrets()
			if tt.expected.IsNull() && len(secrets) != 0 {
				t.Errorf("expected no secrets to be added, got %v", secrets)
			}

			if !tt.expected.IsNull() {
				if len(secrets) != 1 {
					t.Fatal("expected secrets to be added, got none")
				}

				if secrets[0] != tt.expected.AsString() {
					t.Errorf("expected secret %q, got %q", tt.expected.AsString(), secrets[0])
				}
			}
		})
	}
}

func TestSensitiveStringConvert_UnknownValue(t *testing.T) {
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
			verifyFailedConversion(t, SensitiveString, tt.input, "cannot convert unknown value")
		})
	}
}

func TestSensitiveStringConvert_InvalidValues(t *testing.T) {
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
				cty.String.FriendlyName(),
			)

			verifyFailedConversion(t, SensitiveString, tt.input, expectedError)
		})
	}
}

func TestSensitiveStringConvert_Nil(t *testing.T) {
	var sensitiveString *sensitiveStringType

	expectedError := "sensitive string type is nil"
	converted, err := sensitiveString.Convert(cty.StringVal("test"))
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if !converted.Equals(cty.NilVal).True() {
		t.Errorf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestSensitiveStringValidate_Pass(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, SensitiveString, tt.input)
		})
	}
}

func TestSensitiveStringValidateSpec_Pass(t *testing.T) {
	err := SensitiveString.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got %q", err.Error())
	}
}

func TestSensitiveStringValidateSpec_Nil(t *testing.T) {
	var sensitiveString *sensitiveStringType

	expectedError := "sensitive string type is nil"
	err := sensitiveString.ValidateSpec()
	if err == nil {
		t.Fatalf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
