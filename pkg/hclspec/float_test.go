// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"math"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFloatCtyType(t *testing.T) {
	tests := []struct {
		name     string
		typ      Type
		expected cty.Type
	}{
		{
			name:     "floating-point with 32-bit precision",
			typ:      Float32(),
			expected: cty.Number,
		},
		{
			name:     "floating-point with 64-bit precision",
			typ:      Float64(),
			expected: cty.Number,
		},
	}

	for _, tt := range tests {
		actual := tt.typ.CtyType()
		if !actual.Equals(tt.expected) {
			t.Errorf("expected %q from CtyType(), got %q",
				tt.expected.FriendlyName(),
				actual.FriendlyName(),
			)
		}
	}
}

func TestFloatConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		typ      Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid floating-point with 32-bit precision (integer)",
			typ:      Float32(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid floating-point with 64-bit precision (integer)",
			typ:      Float64(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid floating-point with 32-bit precision (float)",
			typ:      Float32(),
			input:    cty.NumberFloatVal(3.14),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "valid floating-point with 64-bit precision (float)",
			typ:      Float64(),
			input:    cty.NumberFloatVal(3.14),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "string to floating-point with 32-bit precision (integer)",
			typ:      Float32(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string to floating-point with 64-bit precision (integer)",
			typ:      Float64(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string to floating-point with 32-bit precision (float)",
			typ:      Float32(),
			input:    cty.StringVal("3.14"),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "string to floating-point with 64-bit precision (float)",
			typ:      Float64(),
			input:    cty.StringVal("3.14"),
			expected: cty.NumberFloatVal(3.14),
		},
		{
			name:     "null floating-point with 32-bit precision",
			typ:      Float32(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null floating-point with 64-bit precision",
			typ:      Float64(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.typ, tt.input, tt.expected)
		})
	}
}

func TestFloatConvert_UnknownValue(t *testing.T) {
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
			verifyFailedConversion(t, Float32(), tt.input, "cannot convert unknown value")
		})
	}
}

func TestFloatConvert_InvalidValues(t *testing.T) {
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
					Float32().CtyType().FriendlyName(),
					tt.conversionError)
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					Float32().CtyType().FriendlyName())
			}

			verifyFailedConversion(t, Float32(), tt.input, expectedError)
		})
	}
}

func TestFloatValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		typ   Type
		input cty.Value
	}{
		{
			name:  "valid floating-point with 32-bit precision (integer)",
			typ:   Float32(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid floating-point with 64-bit precision (integer)",
			typ:   Float64(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid floating-point with 32-bit precision (float)",
			typ:   Float32(),
			input: cty.NumberFloatVal(3.14),
		},
		{
			name:  "valid floating-point with 64-bit precision (float)",
			typ:   Float64(),
			input: cty.NumberFloatVal(3.14),
		},
		{
			name:  "valid floating-point with 32-bit precision (string integer conversion)",
			typ:   Float32(),
			input: cty.StringVal("42"),
		},
		{
			name:  "valid floating-point with 64-bit precision (string integer conversion)",
			typ:   Float64(),
			input: cty.StringVal("42"),
		},
		{
			name:  "valid floating-point with 32-bit precision (string float conversion)",
			typ:   Float32(),
			input: cty.StringVal("3.14"),
		},
		{
			name:  "valid floating-point with 64-bit precision (string float conversion)",
			typ:   Float64(),
			input: cty.StringVal("3.14"),
		},
		{
			name:  "valid floating-point with 32-bit precision (null number)",
			typ:   Float32(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "valid floating-point with 64-bit precision (null number)",
			typ:   Float64(),
			input: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.typ, tt.input)
		})
	}
}

func TestFloatValidate_Fail(t *testing.T) {
	tests := []struct {
		name          string
		typ           Type
		input         cty.Value
		expectedError string
	}{
		{
			name:          "above maximum floating-point with 32-bit precision",
			typ:           Float32(),
			input:         cty.NumberFloatVal(math.MaxFloat64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum floating-point with 32-bit precision",
			typ:           Float32(),
			input:         cty.NumberFloatVal(-math.MaxFloat64),
			expectedError: "value is out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, tt.typ, tt.input, tt.expectedError)
		})
	}
}
