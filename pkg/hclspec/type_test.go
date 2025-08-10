package hclspec

import (
	"fmt"
	"testing"

	"github.com/trippsoft/forge/pkg/log"
	"github.com/zclconf/go-cty/cty"
)

func TestString(t *testing.T) {

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

func TestString_UnknownValue(t *testing.T) {

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
			verifyFailedConversion(t, String, tt.input, expectedError)
		})
	}
}

func TestString_InvalidValues(t *testing.T) {

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

func TestSensitiveString(t *testing.T) {

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
			log.SecretFilter.Clear()

			verifySuccessfulConversion(t, SensitiveString, tt.input, tt.expected)

			secrets := log.SecretFilter.Secrets()
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

func TestSensitiveString_UnknownValue(t *testing.T) {

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
			verifyFailedConversion(t, SensitiveString, tt.input, expectedError)
		})
	}
}

func TestSensitiveString_InvalidValues(t *testing.T) {

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

			verifyFailedConversion(t, SensitiveString, tt.input, expectedError)
		})
	}
}

func TestNumber(t *testing.T) {

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

func TestNumber_UnknownValue(t *testing.T) {

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
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, Number, tt.input, expectedError)
		})
	}
}

func TestNumber_InvalidValues(t *testing.T) {

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

func TestBool(t *testing.T) {

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

func TestBool_UnknownValue(t *testing.T) {

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
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, Bool, tt.input, expectedError)
		})
	}
}

func TestBool_InvalidValues(t *testing.T) {

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

func TestDuration(t *testing.T) {

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, Duration, tt.input, tt.expected)
		})
	}
}

func TestDuration_NonDurationStringValue(t *testing.T) {

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

func TestDuration_UnknownValue(t *testing.T) {

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

func TestDuration_InvalidValues(t *testing.T) {
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

			verifyFailedConversion(t, Duration, tt.input, expectedError)
		})
	}
}

func TestList(t *testing.T) {

	tests := []struct {
		name     string
		list     Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "list of strings",
			list:     List(String),
			input:    cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
			expected: cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
		},
		{
			name:     "empty list of strings",
			list:     List(String),
			input:    cty.ListValEmpty(cty.String),
			expected: cty.ListValEmpty(cty.String),
		},
		{
			name:     "single string",
			list:     List(String),
			input:    cty.StringVal("value"),
			expected: cty.ListVal([]cty.Value{cty.StringVal("value")}),
		},
		{
			name:     "list of numbers to list of strings",
			list:     List(String),
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2)}),
			expected: cty.ListVal([]cty.Value{cty.StringVal("1"), cty.StringVal("2")}),
		},
		{
			name:     "list of booleans to list of strings",
			list:     List(String),
			input:    cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
			expected: cty.ListVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
		},
		{
			name:     "tuple of strings and numbers to list of strings",
			list:     List(String),
			input:    cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.NumberIntVal(1)}),
			expected: cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("1")}),
		},
		{
			name:     "empty tuple to list of strings",
			list:     List(String),
			input:    cty.EmptyTupleVal,
			expected: cty.ListValEmpty(cty.String),
		},
		{
			name:     "null list of strings",
			list:     List(String),
			input:    cty.NullVal(cty.List(cty.String)),
			expected: cty.NullVal(cty.List(cty.String)),
		},
		{
			name:     "list with null values to list of strings",
			list:     List(String),
			input:    cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.NullVal(cty.String)}),
			expected: cty.ListVal([]cty.Value{cty.StringVal("hello")}),
		},
		{
			name:     "list of numbers",
			list:     List(Number),
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NumberFloatVal(3.14)}),
			expected: cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NumberFloatVal(3.14)}),
		},
		{
			name:     "list of strings to list of numbers",
			list:     List(Number),
			input:    cty.ListVal([]cty.Value{cty.StringVal("42"), cty.StringVal("3.14")}),
			expected: cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NumberFloatVal(3.14)}),
		},
		{
			name:     "empty list of numbers",
			list:     List(Number),
			input:    cty.ListValEmpty(cty.Number),
			expected: cty.ListValEmpty(cty.Number),
		},
		{
			name:     "single number",
			list:     List(Number),
			input:    cty.NumberIntVal(123),
			expected: cty.ListVal([]cty.Value{cty.NumberIntVal(123)}),
		},
		{
			name:     "tuple of number strings and numbers to list of numbers",
			list:     List(Number),
			input:    cty.TupleVal([]cty.Value{cty.StringVal("42"), cty.NumberIntVal(1)}),
			expected: cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NumberIntVal(1)}),
		},
		{
			name:     "empty tuple to list of numbers",
			list:     List(Number),
			input:    cty.EmptyTupleVal,
			expected: cty.ListValEmpty(cty.Number),
		},
		{
			name:     "null list of numbers",
			list:     List(Number),
			input:    cty.NullVal(cty.List(cty.Number)),
			expected: cty.NullVal(cty.List(cty.Number)),
		},
		{
			name:     "list with null values to list of numbers",
			list:     List(Number),
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NullVal(cty.Number)}),
			expected: cty.ListVal([]cty.Value{cty.NumberIntVal(42)}),
		},
		{
			name:     "list of booleans",
			list:     List(Bool),
			input:    cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
			expected: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:     "list of strings to list of booleans",
			list:     List(Bool),
			input:    cty.ListVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
			expected: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:     "empty list of booleans",
			list:     List(Bool),
			input:    cty.ListValEmpty(cty.Bool),
			expected: cty.ListValEmpty(cty.Bool),
		},
		{
			name:     "single boolean",
			list:     List(Bool),
			input:    cty.BoolVal(true),
			expected: cty.ListVal([]cty.Value{cty.BoolVal(true)}),
		},
		{
			name:     "tuple of boolean strings and booleans to list of booleans",
			list:     List(Bool),
			input:    cty.TupleVal([]cty.Value{cty.StringVal("true"), cty.BoolVal(false)}),
			expected: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:     "empty tuple to list of booleans",
			list:     List(Bool),
			input:    cty.EmptyTupleVal,
			expected: cty.ListValEmpty(cty.Bool),
		},
		{
			name:     "null list of booleans",
			list:     List(Bool),
			input:    cty.NullVal(cty.List(cty.Bool)),
			expected: cty.NullVal(cty.List(cty.Bool)),
		},
		{
			name:     "list with null values to list of booleans",
			list:     List(Bool),
			input:    cty.ListVal([]cty.Value{cty.BoolVal(true), cty.NullVal(cty.Bool)}),
			expected: cty.ListVal([]cty.Value{cty.BoolVal(true)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.list, tt.input, tt.expected)
		})
	}
}

func TestList_UnknownValue(t *testing.T) {

	tests := []struct {
		name  string
		list  Type
		input cty.Value
	}{
		{
			name:  "unknown list of strings",
			list:  List(String),
			input: cty.UnknownVal(cty.List(cty.String)),
		},
		{
			name:  "unknown tuple of string",
			list:  List(String),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.String})),
		},
		{
			name:  "unknown string",
			list:  List(String),
			input: cty.UnknownVal(cty.String),
		},
		{
			name:  "unknown object with string",
			list:  List(String),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.String})),
		},
		{
			name:  "list with unknown strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.UnknownVal(cty.String), cty.UnknownVal(cty.String)}),
		},
		{
			name:  "tuple with unknown strings",
			list:  List(String),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.String), cty.UnknownVal(cty.String)}),
		},
		{
			name:  "unknown list of numbers",
			list:  List(Number),
			input: cty.UnknownVal(cty.List(cty.Number)),
		},
		{
			name:  "unknown tuple of number",
			list:  List(Number),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.Number})),
		},
		{
			name:  "unknown number",
			list:  List(Number),
			input: cty.UnknownVal(cty.Number),
		},
		{
			name:  "unknown object with number",
			list:  List(Number),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Number})),
		},
		{
			name:  "list with unknown numbers",
			list:  List(Number),
			input: cty.ListVal([]cty.Value{cty.UnknownVal(cty.Number), cty.UnknownVal(cty.Number)}),
		},
		{
			name:  "tuple with unknown numbers",
			list:  List(Number),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Number), cty.UnknownVal(cty.Number)}),
		},
		{
			name:  "unknown list of booleans",
			list:  List(Bool),
			input: cty.UnknownVal(cty.List(cty.Bool)),
		},
		{
			name:  "unknown tuple of boolean",
			list:  List(Bool),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.Bool})),
		},
		{
			name:  "unknown boolean",
			list:  List(Bool),
			input: cty.UnknownVal(cty.Bool),
		},
		{
			name:  "unknown object with boolean",
			list:  List(Bool),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Bool})),
		},
		{
			name:  "list with unknown booleans",
			list:  List(Bool),
			input: cty.ListVal([]cty.Value{cty.UnknownVal(cty.Bool), cty.UnknownVal(cty.Bool)}),
		},
		{
			name:  "tuple with unknown booleans",
			list:  List(Bool),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Bool), cty.UnknownVal(cty.Bool)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, tt.list, tt.input, expectedError)
		})
	}
}

func TestList_InvalidValues(t *testing.T) {

	tests := []struct {
		name            string
		list            Type
		input           cty.Value
		conversionError string
	}{
		{
			name:  "nested list of strings to list of strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.ListVal([]cty.Value{cty.StringVal("hello")})}),
		},
		{
			name: "tuple with strings and lists to list of strings",
			list: List(String),
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.ListVal([]cty.Value{cty.StringVal("world")}),
			}),
		},
		{
			name:  "nested list of numbers to list of numbers",
			list:  List(Number),
			input: cty.ListVal([]cty.Value{cty.ListVal([]cty.Value{cty.NumberIntVal(42)})}),
		},
		{
			name: "tuple with strings and lists to list of strings",
			list: List(Number),
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(42),
				cty.ListVal([]cty.Value{cty.NumberFloatVal(3.14)}),
			}),
		},
		{
			name:            "list with non-number strings to list of numbers",
			list:            List(Number),
			input:           cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("42")}),
			conversionError: "a number is required",
		},
		{
			name:  "nested list of booleans to list of booleans",
			list:  List(Bool),
			input: cty.ListVal([]cty.Value{cty.ListVal([]cty.Value{cty.BoolVal(true)})}),
		},
		{
			name: "tuple with booleans and lists of booleans to list of booleans",
			list: List(Bool),
			input: cty.TupleVal([]cty.Value{
				cty.BoolVal(true),
				cty.ListVal([]cty.Value{cty.BoolVal(false)}),
			}),
		},
		{
			name:            "list with non-boolean strings to list of booleans",
			list:            List(Bool),
			input:           cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("true")}),
			conversionError: "a bool is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var expectedError string
			if tt.conversionError == "" {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					tt.list.CtyType().FriendlyName())
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					tt.list.CtyType().FriendlyName(),
					tt.conversionError)
			}

			verifyFailedConversion(t, tt.list, tt.input, expectedError)
		})
	}
}

func TestList_InvalidElement(t *testing.T) {

	tests := []struct {
		name          string
		list          Type
		input         cty.Value
		expectedError string
	}{
		{
			name: "list of duration with invalid element",
			list: List(Duration),
			input: cty.ListVal([]cty.Value{
				cty.StringVal("invalid"),
				cty.StringVal("30s"),
			}),
			expectedError: fmt.Sprintf("element at index 0: time: invalid duration %q", "invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, tt.list, tt.input, tt.expectedError)
		})
	}
}

func TestMap(t *testing.T) {

	tests := []struct {
		name     string
		mapType  Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name:    "map of strings",
			mapType: Map(String),
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
			name:     "empty map of strings",
			mapType:  Map(String),
			input:    cty.MapValEmpty(cty.String),
			expected: cty.MapValEmpty(cty.String),
		},
		{
			name:     "null map of strings",
			mapType:  Map(String),
			input:    cty.NullVal(cty.Map(cty.String)),
			expected: cty.NullVal(cty.Map(cty.String)),
		},
		{
			name:     "null string to map of strings",
			mapType:  Map(String),
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.Map(cty.String)),
		},
		{
			name:    "map with null values",
			mapType: Map(String),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.NullVal(cty.String),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.NullVal(cty.String),
			}),
		},
		{
			name:    "map of numbers to map of strings",
			mapType: Map(String),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberFloatVal(3.14),
				"key2": cty.NumberIntVal(42),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("3.14"),
				"key2": cty.StringVal("42"),
			}),
		},
		{
			name:    "object to map of strings",
			mapType: Map(String),
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.NumberIntVal(42),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("42"),
			}),
		},
		{
			name:    "map of numbers",
			mapType: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberFloatVal(2.5),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberFloatVal(2.5),
			}),
		},
		{
			name:     "empty map of numbers",
			mapType:  Map(Number),
			input:    cty.MapValEmpty(cty.Number),
			expected: cty.MapValEmpty(cty.Number),
		},
		{
			name:     "null map of numbers",
			mapType:  Map(Number),
			input:    cty.NullVal(cty.Map(cty.Number)),
			expected: cty.NullVal(cty.Map(cty.Number)),
		},
		{
			name:     "null number to map of numbers",
			mapType:  Map(Number),
			input:    cty.NullVal(cty.Number),
			expected: cty.NullVal(cty.Map(cty.Number)),
		},
		{
			name:    "map with null values",
			mapType: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NullVal(cty.Number),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NullVal(cty.Number),
			}),
		},
		{
			name:    "map of strings to map of numbers",
			mapType: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("3.14"),
				"key2": cty.StringVal("42"),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberFloatVal(3.14),
				"key2": cty.NumberIntVal(42),
			}),
		},
		{
			name:    "object to map of numbers",
			mapType: Map(Number),
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("3.14"),
				"key2": cty.NumberIntVal(42),
			}),
			expected: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberFloatVal(3.14),
				"key2": cty.NumberIntVal(42),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.mapType, tt.input, tt.expected)
		})
	}
}

func TestMap_UnknownValue(t *testing.T) {

	tests := []struct {
		name    string
		mapType Type
		input   cty.Value
	}{
		{
			name:    "unknown map of strings",
			mapType: Map(String),
			input:   cty.UnknownVal(cty.Map(cty.String)),
		},
		{
			name:    "unknown tuple of string",
			mapType: Map(String),
			input:   cty.UnknownVal(cty.Tuple([]cty.Type{cty.String})),
		},
		{
			name:    "unknown string",
			mapType: Map(String),
			input:   cty.UnknownVal(cty.String),
		},
		{
			name:    "unknown object with string",
			mapType: Map(String),
			input:   cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.String})),
		},
		{
			name:    "list with unknown strings",
			mapType: Map(String),
			input:   cty.MapVal(map[string]cty.Value{"key1": cty.UnknownVal(cty.String), "key2": cty.UnknownVal(cty.String)}),
		},
		{
			name:    "tuple with unknown strings",
			mapType: Map(String),
			input:   cty.TupleVal([]cty.Value{cty.UnknownVal(cty.String), cty.UnknownVal(cty.String)}),
		},
		{
			name:    "unknown map of numbers",
			mapType: Map(Number),
			input:   cty.UnknownVal(cty.Map(cty.Number)),
		},
		{
			name:    "unknown tuple of number",
			mapType: Map(Number),
			input:   cty.UnknownVal(cty.Tuple([]cty.Type{cty.Number})),
		},
		{
			name:    "unknown number",
			mapType: Map(Number),
			input:   cty.UnknownVal(cty.Number),
		},
		{
			name:    "unknown object with number",
			mapType: Map(Number),
			input:   cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Number})),
		},
		{
			name:    "map with unknown numbers",
			mapType: Map(Number),
			input:   cty.MapVal(map[string]cty.Value{"key1": cty.UnknownVal(cty.Number), "key2": cty.UnknownVal(cty.Number)}),
		},
		{
			name:    "tuple with unknown numbers",
			mapType: Map(Number),
			input:   cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Number), cty.UnknownVal(cty.Number)}),
		},
		{
			name:    "unknown map of booleans",
			mapType: Map(Bool),
			input:   cty.UnknownVal(cty.Map(cty.Bool)),
		},
		{
			name:    "unknown tuple of boolean",
			mapType: Map(Bool),
			input:   cty.UnknownVal(cty.Tuple([]cty.Type{cty.Bool})),
		},
		{
			name:    "unknown boolean",
			mapType: Map(Bool),
			input:   cty.UnknownVal(cty.Bool),
		},
		{
			name:    "unknown object with boolean",
			mapType: Map(Bool),
			input:   cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Bool})),
		},
		{
			name:    "map with unknown booleans",
			mapType: Map(Bool),
			input:   cty.MapVal(map[string]cty.Value{"key1": cty.UnknownVal(cty.Bool), "key2": cty.UnknownVal(cty.Bool)}),
		},
		{
			name:    "tuple with unknown booleans",
			mapType: Map(Bool),
			input:   cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Bool), cty.UnknownVal(cty.Bool)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, tt.mapType, tt.input, expectedError)
		})
	}
}

func TestMap_InvalidValues(t *testing.T) {

	tests := []struct {
		name            string
		list            Type
		input           cty.Value
		conversionError string
	}{
		{
			name: "nested map of strings to map of strings",
			list: Map(String),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.MapVal(map[string]cty.Value{
					"nestedKey": cty.StringVal("value"),
				}),
			}),
		},
		{
			name: "tuple with strings and lists to map of strings",
			list: Map(String),
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.ListVal([]cty.Value{cty.StringVal("world")}),
			}),
		},
		{
			name: "nested map of numbers to map of numbers",
			list: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.MapVal(map[string]cty.Value{
					"nestedKey": cty.NumberIntVal(42),
				}),
			}),
		},
		{
			name: "tuple with strings and lists to map of strings",
			list: Map(String),
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(42),
				cty.ListVal([]cty.Value{cty.NumberFloatVal(3.14)}),
			}),
		},
		{
			name: "map with non-number strings to map of numbers",
			list: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("hello"),
				"key2": cty.StringVal("42"),
			}),
			conversionError: "a number is required",
		},
		{
			name: "nested map of booleans to map of booleans",
			list: Map(Bool),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.MapVal(map[string]cty.Value{"nestedKey": cty.BoolVal(true)}),
			}),
		},
		{
			name: "tuple with booleans and lists of booleans to map of booleans",
			list: Map(Bool),
			input: cty.TupleVal([]cty.Value{
				cty.BoolVal(true),
				cty.ListVal([]cty.Value{cty.BoolVal(false)}),
			}),
		},
		{
			name: "map with non-boolean strings to map of booleans",
			list: Map(Bool),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("hello"),
				"key2": cty.StringVal("true"),
			}),
			conversionError: "a bool is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var expectedError string
			if tt.conversionError == "" {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					tt.list.CtyType().FriendlyName())
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					tt.list.CtyType().FriendlyName(),
					tt.conversionError)
			}

			verifyFailedConversion(t, tt.list, tt.input, expectedError)
		})
	}
}

func TestMap_InvalidElement(t *testing.T) {

	tests := []struct {
		name          string
		list          Type
		input         cty.Value
		expectedError string
	}{
		{
			name: "map of duration with invalid element",
			list: Map(Duration),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("invalid"),
				"key2": cty.StringVal("30s"),
			}),
			expectedError: fmt.Sprintf("element at index %q: time: invalid duration %q", "key1", "invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, tt.list, tt.input, tt.expectedError)
		})
	}
}

func TestSet(t *testing.T) {

	tests := []struct {
		name     string
		set      Type
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "set of strings",
			set:      Set(String),
			input:    cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
			expected: cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		},
		{
			name:     "empty set of strings",
			set:      Set(String),
			input:    cty.SetValEmpty(cty.String),
			expected: cty.SetValEmpty(cty.String),
		},
		{
			name:     "null set of strings",
			set:      Set(String),
			input:    cty.NullVal(cty.Set(cty.String)),
			expected: cty.NullVal(cty.Set(cty.String)),
		},
		{
			name:     "set of strings with null values",
			set:      Set(String),
			input:    cty.SetVal([]cty.Value{cty.StringVal("a"), cty.NullVal(cty.String)}),
			expected: cty.SetVal([]cty.Value{cty.StringVal("a")}),
		},
		{
			name:     "single string to set of strings",
			set:      Set(String),
			input:    cty.StringVal("a"),
			expected: cty.SetVal([]cty.Value{cty.StringVal("a")}),
		},
		{
			name:     "list of strings to set of strings",
			set:      Set(String),
			input:    cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
			expected: cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		},
		{
			name:     "set of numbers to set of strings",
			set:      Set(String),
			input:    cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
			expected: cty.SetVal([]cty.Value{cty.StringVal("1"), cty.StringVal("2.5")}),
		},
		{
			name:     "set of numbers",
			set:      Set(Number),
			input:    cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
			expected: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:     "empty set of numbers",
			set:      Set(Number),
			input:    cty.SetValEmpty(cty.Number),
			expected: cty.SetValEmpty(cty.Number),
		},
		{
			name:     "null set of numbers",
			set:      Set(Number),
			input:    cty.NullVal(cty.Set(cty.Number)),
			expected: cty.NullVal(cty.Set(cty.Number)),
		},
		{
			name:     "set of numbers with null values",
			set:      Set(Number),
			input:    cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NullVal(cty.Number)}),
			expected: cty.SetVal([]cty.Value{cty.NumberIntVal(1)}),
		},
		{
			name:     "single number to set of numbers",
			set:      Set(Number),
			input:    cty.NumberIntVal(1),
			expected: cty.SetVal([]cty.Value{cty.NumberIntVal(1)}),
		},
		{
			name:     "list of numbers to set of numbers",
			set:      Set(Number),
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
			expected: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:     "set of number-strings to set of numbers",
			set:      Set(Number),
			input:    cty.SetVal([]cty.Value{cty.StringVal("1"), cty.StringVal("2.5")}),
			expected: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:     "set of booleans",
			set:      Set(Bool),
			input:    cty.SetVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
			expected: cty.SetVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:     "empty set of booleans",
			set:      Set(Bool),
			input:    cty.SetValEmpty(cty.Bool),
			expected: cty.SetValEmpty(cty.Bool),
		},
		{
			name:     "null set of booleans",
			set:      Set(Bool),
			input:    cty.NullVal(cty.Set(cty.Bool)),
			expected: cty.NullVal(cty.Set(cty.Bool)),
		},
		{
			name:     "set of booleans with null values",
			set:      Set(Bool),
			input:    cty.SetVal([]cty.Value{cty.BoolVal(true), cty.NullVal(cty.Bool)}),
			expected: cty.SetVal([]cty.Value{cty.BoolVal(true)}),
		},
		{
			name:     "single boolean to set of booleans",
			set:      Set(Bool),
			input:    cty.BoolVal(true),
			expected: cty.SetVal([]cty.Value{cty.BoolVal(true)}),
		},
		{
			name:     "list of booleans to set of booleans",
			set:      Set(Bool),
			input:    cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
			expected: cty.SetVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:     "set of boolean-strings to set of booleans",
			set:      Set(Bool),
			input:    cty.SetVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
			expected: cty.SetVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.set, tt.input, tt.expected)
		})
	}
}

func TestSet_UnknownValue(t *testing.T) {

	tests := []struct {
		name  string
		set   Type
		input cty.Value
	}{
		{
			name:  "unknown set of strings",
			set:   Set(String),
			input: cty.UnknownVal(cty.Set(cty.String)),
		},
		{
			name:  "unknown tuple of string",
			set:   Set(String),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.String})),
		},
		{
			name:  "unknown string",
			set:   Set(String),
			input: cty.UnknownVal(cty.String),
		},
		{
			name:  "unknown object with string",
			set:   Set(String),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.String})),
		},
		{
			name:  "set with unknown strings",
			set:   Set(String),
			input: cty.SetVal([]cty.Value{cty.UnknownVal(cty.String), cty.UnknownVal(cty.String)}),
		},
		{
			name:  "tuple with unknown strings",
			set:   Set(String),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.String), cty.UnknownVal(cty.String)}),
		},
		{
			name:  "unknown set of numbers",
			set:   Set(Number),
			input: cty.UnknownVal(cty.Set(cty.Number)),
		},
		{
			name:  "unknown tuple of number",
			set:   Set(Number),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.Number})),
		},
		{
			name:  "unknown number",
			set:   Set(Number),
			input: cty.UnknownVal(cty.Number),
		},
		{
			name:  "unknown object with number",
			set:   Set(Number),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Number})),
		},
		{
			name:  "set with unknown numbers",
			set:   Set(Number),
			input: cty.SetVal([]cty.Value{cty.UnknownVal(cty.Number), cty.UnknownVal(cty.Number)}),
		},
		{
			name:  "tuple with unknown numbers",
			set:   Set(Number),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Number), cty.UnknownVal(cty.Number)}),
		},
		{
			name:  "unknown set of booleans",
			set:   Set(Bool),
			input: cty.UnknownVal(cty.Set(cty.Bool)),
		},
		{
			name:  "unknown tuple of boolean",
			set:   Set(Bool),
			input: cty.UnknownVal(cty.Tuple([]cty.Type{cty.Bool})),
		},
		{
			name:  "unknown boolean",
			set:   Set(Bool),
			input: cty.UnknownVal(cty.Bool),
		},
		{
			name:  "unknown object with boolean",
			set:   Set(Bool),
			input: cty.UnknownVal(cty.Object(map[string]cty.Type{"key": cty.Bool})),
		},
		{
			name:  "set with unknown booleans",
			set:   Set(Bool),
			input: cty.SetVal([]cty.Value{cty.UnknownVal(cty.Bool), cty.UnknownVal(cty.Bool)}),
		},
		{
			name:  "tuple with unknown booleans",
			set:   Set(Bool),
			input: cty.TupleVal([]cty.Value{cty.UnknownVal(cty.Bool), cty.UnknownVal(cty.Bool)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedError := "cannot convert unknown value"
			verifyFailedConversion(t, tt.set, tt.input, expectedError)
		})
	}
}

func TestSet_InvalidValues(t *testing.T) {

	tests := []struct {
		name            string
		set             Type
		input           cty.Value
		conversionError string
	}{
		{
			name:  "nested set of strings to set of strings",
			set:   Set(String),
			input: cty.SetVal([]cty.Value{cty.SetVal([]cty.Value{cty.StringVal("hello")})}),
		},
		{
			name: "tuple with strings and sets to set of strings",
			set:  Set(String),
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.SetVal([]cty.Value{cty.StringVal("world")}),
			}),
		},
		{
			name:  "nested set of numbers to set of numbers",
			set:   Set(Number),
			input: cty.SetVal([]cty.Value{cty.SetVal([]cty.Value{cty.NumberIntVal(42)})}),
		},
		{
			name: "tuple with strings and sets to set of strings",
			set:  Set(Number),
			input: cty.TupleVal([]cty.Value{
				cty.NumberIntVal(42),
				cty.SetVal([]cty.Value{cty.NumberFloatVal(3.14)}),
			}),
		},
		{
			name:            "set with non-number strings to set of numbers",
			set:             Set(Number),
			input:           cty.SetVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("42")}),
			conversionError: "a number is required",
		},
		{
			name:  "nested set of booleans to set of booleans",
			set:   Set(Bool),
			input: cty.SetVal([]cty.Value{cty.SetVal([]cty.Value{cty.BoolVal(true)})}),
		},
		{
			name: "tuple with booleans and sets of booleans to set of booleans",
			set:  Set(Bool),
			input: cty.TupleVal([]cty.Value{
				cty.BoolVal(true),
				cty.SetVal([]cty.Value{cty.BoolVal(false)}),
			}),
		},
		{
			name:            "set with non-boolean strings to set of booleans",
			set:             Set(Bool),
			input:           cty.SetVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("true")}),
			conversionError: "a bool is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var expectedError string
			if tt.conversionError == "" {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q",
					tt.input.Type().FriendlyName(),
					tt.set.CtyType().FriendlyName())
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					tt.set.CtyType().FriendlyName(),
					tt.conversionError)
			}

			verifyFailedConversion(t, tt.set, tt.input, expectedError)
		})
	}
}

func TestSet_InvalidElement(t *testing.T) {

	tests := []struct {
		name          string
		list          Type
		input         cty.Value
		expectedError string
	}{
		{
			name: "set of duration with invalid element",
			list: Set(Duration),
			input: cty.SetVal([]cty.Value{
				cty.StringVal("invalid"),
				cty.StringVal("30s"),
			}),
			expectedError: fmt.Sprintf("invalid set element: time: invalid duration %q", "invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFailedValidation(t, tt.list, tt.input, tt.expectedError)
		})
	}
}

func verifySuccessfulConversion(t *testing.T, ty Type, input, expected cty.Value) {

	actual, err := ty.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got %v", err)
	}

	if !actual.Type().Equals(ty.CtyType()) {
		t.Errorf(
			"expected Convert() to produce type %q, got %q",
			ty.CtyType().FriendlyName(),
			actual.Type().FriendlyName())
	}

	if actual.Equals(expected) != cty.True {
		t.Errorf(
			"expected Convert() to produce value %q, got %q",
			expected.GoString(),
			actual.GoString())
	}

	err = ty.Validate(actual)
	if err != nil {
		t.Errorf("expected no error from Validate(), got %v", err)
	}
}

func verifyFailedConversion(t *testing.T, ty Type, input cty.Value, expectedError string) {

	_, err := ty.Convert(input)
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}

	err = ty.Validate(input)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func verifyFailedValidation(t *testing.T, ty Type, input cty.Value, expectedError string) {

	_, err := ty.Convert(input)
	if err != nil {
		t.Fatalf("expected no error from Convert(), got %v", err)
	}

	err = ty.Validate(input)
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}
