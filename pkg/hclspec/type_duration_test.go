// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestDurationCtyType(t *testing.T) {
	durationType := &durationType{}
	actual := durationType.CtyType()
	if !actual.Equals(cty.String) {
		t.Errorf("expected %q, got %q", cty.String.FriendlyName(), actual.FriendlyName())
	}
}

func TestDurationCtyType_Nil(t *testing.T) {
	var durationType *durationType
	actual := durationType.CtyType()
	if !actual.Equals(cty.String) {
		t.Errorf("expected %q, got %q", cty.String.FriendlyName(), actual.FriendlyName())
	}
}

func TestDurationConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "nil duration",
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.String),
		},
		{
			name:     "valid duration",
			input:    cty.StringVal("5m"),
			expected: cty.StringVal("5m"),
		},
		{
			name:     "valid complex duration",
			input:    cty.StringVal("1h30m45s"),
			expected: cty.StringVal("1h30m45s"),
		},
		{
			name:     "invalid duration format",
			input:    cty.StringVal("not-a-duration"),
			expected: cty.StringVal("not-a-duration"),
		},
		{
			name:     "number conversion to string",
			input:    cty.NumberIntVal(123),
			expected: cty.StringVal("123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, Duration, tt.input, tt.expected)
		})
	}
}

func TestDurationConvert_UnknownValue(t *testing.T) {
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
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, Duration, tt.input, expectedError)
		})
	}
}

func TestDurationConvert_InvalidValues(t *testing.T) {
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

			verifyFailedConversion(t, Duration, tt.input, expectedError)
		})
	}
}

func TestDurationConvert_Nil(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid duration",
			input:    cty.StringVal("5m"),
			expected: cty.StringVal("5m"),
		},
		{
			name:     "valid complex duration",
			input:    cty.StringVal("1h30m45s"),
			expected: cty.StringVal("1h30m45s"),
		},
		{
			name:     "invalid duration format",
			input:    cty.StringVal("not-a-duration"),
			expected: cty.StringVal("not-a-duration"),
		},
		{
			name:     "number conversion to string",
			input:    cty.NumberIntVal(123),
			expected: cty.StringVal("123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nilDuration *durationType
			verifySuccessfulConversion(t, nilDuration, tt.input, tt.expected)
		})
	}
}

func TestDurationValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "nil duration",
			input: cty.NullVal(cty.String),
		},
		{
			name:  "valid duration",
			input: cty.StringVal("5m"),
		},
		{
			name:  "valid complex duration",
			input: cty.StringVal("1h30m45s"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, Duration, tt.input)
		})
	}
}

func TestDurationValidate_NonDurationStringValue(t *testing.T) {
	tests := []struct {
		name          string
		input         cty.Value
		expectedError string
	}{
		{
			name:          "invalid duration format",
			input:         cty.StringVal("not-a-duration"),
			expectedError: "time: invalid duration \"not-a-duration\"",
		},
		{
			name:          "number conversion to string",
			input:         cty.NumberIntVal(123),
			expectedError: "time: missing unit in duration \"123\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, Duration, tt.input, tt.expectedError)
		})
	}
}

func TestDurationValidateSpec_Pass(t *testing.T) {
	err := Duration.ValidateSpec()
	if err != nil {
		t.Errorf("expected no error from ValidateSpec(), got error %q", err.Error())
	}
}

func TestDurationValidateSpec_Nil(t *testing.T) {
	var nilDuration *durationType

	expectedError := "duration type is nil"
	err := nilDuration.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
