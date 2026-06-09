// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"math"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestIntegerCtyType(t *testing.T) {
	tests := []struct {
		name     string
		typ      Type
		expected cty.Type
	}{
		{
			name:     "signed 8-bit integer",
			typ:      Int8(),
			expected: cty.Number,
		},
		{
			name:     "signed 16-bit integer",
			typ:      Int16(),
			expected: cty.Number,
		},
		{
			name:     "signed 32-bit integer",
			typ:      Int32(),
			expected: cty.Number,
		},
		{
			name:     "signed 64-bit integer",
			typ:      Int64(),
			expected: cty.Number,
		},
		{
			name:     "unsigned 8-bit integer",
			typ:      UInt8(),
			expected: cty.Number,
		},
		{
			name:     "unsigned 16-bit integer",
			typ:      UInt16(),
			expected: cty.Number,
		},
		{
			name:     "unsigned 32-bit integer",
			typ:      UInt32(),
			expected: cty.Number,
		},
		{
			name:     "unsigned 64-bit integer",
			typ:      UInt64(),
			expected: cty.Number,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.typ.CtyType()
			if !actual.Equals(tt.expected) {
				t.Errorf("expected %q from CtyType(), got %q",
					tt.expected.FriendlyName(),
					actual.FriendlyName(),
				)
			}
		})
	}
}

func TestIntegerConvert_Success(t *testing.T) {
	tests := []struct {
		name     string
		typ      Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "valid 8-bit signed integer",
			typ:      Int8(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 16-bit signed integer",
			typ:      Int16(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 32-bit signed integer",
			typ:      Int32(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 64-bit signed integer",
			typ:      Int64(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 8-bit unsigned integer",
			typ:      UInt8(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 16-bit unsigned integer",
			typ:      UInt16(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 32-bit unsigned integer",
			typ:      UInt32(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "valid 64-bit unsigned integer",
			typ:      UInt64(),
			input:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 8-bit signed integer conversion",
			typ:      Int8(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 16-bit signed integer conversion",
			typ:      Int16(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 32-bit signed integer conversion",
			typ:      Int32(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 64-bit signed integer conversion",
			typ:      Int64(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 8-bit unsigned integer conversion",
			typ:      UInt8(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 16-bit unsigned integer conversion",
			typ:      UInt16(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 32-bit unsigned integer conversion",
			typ:      UInt32(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "string 64-bit unsigned integer conversion",
			typ:      UInt64(),
			input:    cty.StringVal("42"),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "null 8-bit signed integer",
			typ:      Int8(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 16-bit signed integer",
			typ:      Int16(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 32-bit signed integer",
			typ:      Int32(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 64-bit signed integer",
			typ:      Int64(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 8-bit unsigned integer",
			typ:      UInt8(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 16-bit unsigned integer",
			typ:      UInt16(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 32-bit unsigned integer",
			typ:      UInt32(),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Number),
		},
		{
			name:     "null 64-bit unsigned integer",
			typ:      UInt64(),
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

func TestIntegerConvert_UnknownValue(t *testing.T) {
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
			verifyFailedConversion(t, Int32(), tt.input, "cannot convert unknown value")
		})
	}
}

func TestIntegerConvert_InvalidValues(t *testing.T) {
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
					Int32().CtyType().FriendlyName(),
					tt.conversionError)
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					Int32().CtyType().FriendlyName())
			}

			verifyFailedConversion(t, Int32(), tt.input, expectedError)
		})
	}
}

func TestIntegerValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		typ   Type
		input cty.Value
	}{
		{
			name:  "valid 8-bit signed integer",
			typ:   Int8(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 16-bit signed integer",
			typ:   Int16(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 32-bit signed integer",
			typ:   Int32(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 64-bit signed integer",
			typ:   Int64(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 8-bit unsigned integer",
			typ:   UInt8(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 16-bit unsigned integer",
			typ:   UInt16(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 32-bit unsigned integer",
			typ:   UInt32(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "valid 64-bit unsigned integer",
			typ:   UInt64(),
			input: cty.NumberIntVal(42),
		},
		{
			name:  "string 8-bit signed integer conversion",
			typ:   Int8(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 16-bit signed integer conversion",
			typ:   Int16(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 32-bit signed integer conversion",
			typ:   Int32(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 64-bit signed integer conversion",
			typ:   Int64(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 8-bit unsigned integer conversion",
			typ:   UInt8(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 16-bit unsigned integer conversion",
			typ:   UInt16(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 32-bit unsigned integer conversion",
			typ:   UInt32(),
			input: cty.StringVal("42"),
		},
		{
			name:  "string 64-bit unsigned integer conversion",
			typ:   UInt64(),
			input: cty.StringVal("42"),
		},
		{
			name:  "null 8-bit signed integer",
			typ:   Int8(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 16-bit signed integer",
			typ:   Int16(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 32-bit signed integer",
			typ:   Int32(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 64-bit signed integer",
			typ:   Int64(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 8-bit unsigned integer",
			typ:   UInt8(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 16-bit unsigned integer",
			typ:   UInt16(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 32-bit unsigned integer",
			typ:   UInt32(),
			input: cty.NullVal(cty.Number),
		},
		{
			name:  "null 64-bit unsigned integer",
			typ:   UInt64(),
			input: cty.NullVal(cty.Number),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.typ, tt.input)
		})
	}
}

func TestIntegerValidate_Fail(t *testing.T) {
	tests := []struct {
		name          string
		typ           Type
		input         cty.Value
		expectedError string
	}{
		{
			name:          "floating-point as 8-bit signed integer",
			typ:           Int8(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 16-bit signed integer",
			typ:           Int16(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 32-bit signed integer",
			typ:           Int32(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 64-bit signed integer",
			typ:           Int64(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 8-bit unsigned integer",
			typ:           UInt8(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 16-bit unsigned integer",
			typ:           UInt16(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 32-bit unsigned integer",
			typ:           UInt32(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "floating-point as 64-bit unsigned integer",
			typ:           UInt64(),
			input:         cty.NumberFloatVal(3.14),
			expectedError: "value is not an integer",
		},
		{
			name:          "above maximum 8-bit signed integer",
			typ:           Int8(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 16-bit signed integer",
			typ:           Int16(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 32-bit signed integer",
			typ:           Int32(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 64-bit signed integer",
			typ:           Int64(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 8-bit unsigned integer",
			typ:           UInt8(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 16-bit unsigned integer",
			typ:           UInt16(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "above maximum 32-bit unsigned integer",
			typ:           UInt32(),
			input:         cty.NumberUIntVal(math.MaxUint64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 8-bit signed integer",
			typ:           Int8(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 16-bit signed integer",
			typ:           Int16(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 32-bit signed integer",
			typ:           Int32(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 8-bit unsigned integer",
			typ:           UInt8(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 16-bit unsigned integer",
			typ:           UInt16(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 32-bit unsigned integer",
			typ:           UInt32(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
		{
			name:          "below minimum 64-bit unsigned integer",
			typ:           UInt64(),
			input:         cty.NumberIntVal(math.MinInt64),
			expectedError: "value is out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, tt.typ, tt.input, tt.expectedError)
		})
	}
}
