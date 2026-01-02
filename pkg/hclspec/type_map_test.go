// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"slices"
	"testing"

	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func TestMapCtyType(t *testing.T) {
	tests := []struct {
		name     string
		mapType  Type
		expected cty.Type
	}{
		{
			name:     "map of strings",
			mapType:  Map(String),
			expected: cty.Map(cty.String),
		},
		{
			name:     "map of numbers",
			mapType:  Map(Number),
			expected: cty.Map(cty.Number),
		},
		{
			name:     "map of booleans",
			mapType:  Map(Bool),
			expected: cty.Map(cty.Bool),
		},
		{
			name:     "map of lists",
			mapType:  Map(List(String)),
			expected: cty.Map(cty.List(cty.String)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.mapType.CtyType()
			if !actual.Equals(tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected.FriendlyName(), actual.FriendlyName())
			}
		})
	}
}

func TestMapCtyType_Nil(t *testing.T) {
	var nilMap *mapType

	expected := cty.NilType
	actual := nilMap.CtyType()
	if !actual.Equals(expected) {
		t.Errorf("expected nil, got %q", actual.FriendlyName())
	}
}

func TestMapConvert_Success(t *testing.T) {
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

func TestMapConvert_SensitiveString(t *testing.T) {
	util.SecretFilter.Clear()

	str := "hello"
	mapType := Map(SensitiveString)
	input := cty.MapVal(map[string]cty.Value{"key": cty.StringVal(str)})
	verifySuccessfulConversion(t, mapType, input, input)

	secrets := util.SecretFilter.Secrets()
	if !slices.Contains(secrets, str) {
		t.Errorf("expected %q to be filtered", str)
	}
}

func TestMapConvert_UnknownValue(t *testing.T) {
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

func TestMapConvert_InvalidValues(t *testing.T) {
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

func TestMapConvert_Nil(t *testing.T) {
	var nilMap *mapType

	expectedError := "map type is nil"
	converted, err := nilMap.Convert(cty.MapValEmpty(cty.String))
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if converted.Equals(cty.NilVal) != cty.True {
		t.Errorf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func TestMapValidate_Pass(t *testing.T) {
	tests := []struct {
		name    string
		mapType Type
		input   cty.Value
	}{
		{
			name:    "map of strings",
			mapType: Map(String),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
			}),
		},
		{
			name:    "empty map of strings",
			mapType: Map(String),
			input:   cty.MapValEmpty(cty.String),
		},
		{
			name:    "null map of strings",
			mapType: Map(String),
			input:   cty.NullVal(cty.Map(cty.String)),
		},
		{
			name:    "null string to map of strings",
			mapType: Map(String),
			input:   cty.NullVal(cty.String),
		},
		{
			name:    "map with null values",
			mapType: Map(String),
			input: cty.MapVal(map[string]cty.Value{
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
		},
		{
			name:    "object to map of strings",
			mapType: Map(String),
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.NumberIntVal(42),
			}),
		},
		{
			name:    "map of numbers",
			mapType: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberFloatVal(2.5),
			}),
		},
		{
			name:    "empty map of numbers",
			mapType: Map(Number),
			input:   cty.MapValEmpty(cty.Number),
		},
		{
			name:    "null map of numbers",
			mapType: Map(Number),
			input:   cty.NullVal(cty.Map(cty.Number)),
		},
		{
			name:    "null number to map of numbers",
			mapType: Map(Number),
			input:   cty.NullVal(cty.Number),
		},
		{
			name:    "map with null values",
			mapType: Map(Number),
			input: cty.MapVal(map[string]cty.Value{
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
		},
		{
			name:    "object to map of numbers",
			mapType: Map(Number),
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("3.14"),
				"key2": cty.NumberIntVal(42),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.mapType, tt.input)
		})
	}
}

func TestMapValidate_InvalidElement(t *testing.T) {
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

func TestMapValidate_Nil(t *testing.T) {
	var nilMap *mapType

	expectedError := "map type is nil"
	err := nilMap.Validate(cty.MapValEmpty(cty.String))
	if err == nil {
		t.Errorf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestMapValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name    string
		mapType Type
		input   cty.Value
	}{
		{
			name:    "valid map of strings",
			mapType: Map(String),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("foo"),
				"key2": cty.StringVal("bar"),
			}),
		},
		{
			name:    "valid map of objects",
			mapType: Map(Object(RequiredField("foo", String))),
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.ObjectVal(map[string]cty.Value{
					"foo": cty.StringVal("bar"),
				}),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mapType.ValidateSpec()
			if err != nil {
				t.Errorf("expected no error from ValidateSpec(), got %q", err)
			}
		})
	}
}

func TestMapValidateSpec_InvalidObject(t *testing.T) {
	mapType := Map(Object(RequiredField("foo", String)).WithConstraints(MutuallyExclusive("nonexistant")))

	expectedError := `constraint validation failed: field "nonexistant" is not defined in the object type`
	err := mapType.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestMapValidateSpec_Nil(t *testing.T) {
	var nilMap *mapType

	expectedError := "map type is nil"
	err := nilMap.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
