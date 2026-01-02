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

func TestListCtyType(t *testing.T) {
	tests := []struct {
		name     string
		list     Type
		expected cty.Type
	}{
		{
			name:     "list of strings",
			list:     List(String),
			expected: cty.List(cty.String),
		},
		{
			name:     "list of numbers",
			list:     List(Number),
			expected: cty.List(cty.Number),
		},
		{
			name:     "list of booleans",
			list:     List(Bool),
			expected: cty.List(cty.Bool),
		},
		{
			name:     "list of durations",
			list:     List(Duration),
			expected: cty.List(cty.String),
		},
		{
			name:     "list of list of strings",
			list:     List(List(String)),
			expected: cty.List(cty.List(cty.String)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.list.CtyType()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected.FriendlyName(), actual.FriendlyName())
			}
		})
	}
}

func TestListCtyType_Nil(t *testing.T) {
	var nilList *listType
	expected := cty.NilType
	actual := nilList.CtyType()
	if actual != expected {
		t.Errorf("expected nil, got %q", actual.FriendlyName())
	}
}

func TestListConvert_Success(t *testing.T) {
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
			name:     "null string",
			list:     List(String),
			input:    cty.NullVal(cty.String),
			expected: cty.NullVal(cty.List(cty.String)),
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
		{
			name: "list of duration with invalid element",
			list: List(Duration),
			input: cty.ListVal([]cty.Value{
				cty.StringVal("invalid"),
				cty.StringVal("30s"),
			}),
			expected: cty.ListVal([]cty.Value{
				cty.StringVal("invalid"),
				cty.StringVal("30s"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulConversion(t, tt.list, tt.input, tt.expected)
		})
	}
}

func TestListConvert_SensitiveString(t *testing.T) {
	util.SecretFilter.Clear()
	str := "hello"
	list := List(SensitiveString)
	input := cty.ListVal([]cty.Value{cty.StringVal(str)})
	verifySuccessfulConversion(t, list, input, input)

	secrets := util.SecretFilter.Secrets()
	if !slices.Contains(secrets, str) {
		t.Errorf("expected %q to be filtered", str)
	}
}

func TestListConvert_UnknownValue(t *testing.T) {
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

func TestListConvert_InvalidValues(t *testing.T) {
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
					tt.list.CtyType().FriendlyName(),
				)
			} else {
				expectedError = fmt.Sprintf(
					"cannot convert %q to %q: %s",
					tt.input.Type().FriendlyName(),
					tt.list.CtyType().FriendlyName(),
					tt.conversionError,
				)
			}

			verifyFailedConversion(t, tt.list, tt.input, expectedError)
		})
	}
}

func TestListConvert_Nil(t *testing.T) {
	var nilList *listType

	expectedError := "list type is nil"
	converted, err := nilList.Convert(cty.StringVal("test"))
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if !converted.Equals(cty.NilVal).True() {
		t.Errorf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func TestListValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		list  Type
		input cty.Value
	}{
		{
			name:  "list of strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal("world")}),
		},
		{
			name:  "empty list of strings",
			list:  List(String),
			input: cty.ListValEmpty(cty.String),
		},
		{
			name:  "single string",
			list:  List(String),
			input: cty.StringVal("value"),
		},
		{
			name:  "null string",
			list:  List(String),
			input: cty.NullVal(cty.String),
		},
		{
			name:  "list of numbers to list of strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2)}),
		},
		{
			name:  "list of booleans to list of strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:  "tuple of strings and numbers to list of strings",
			list:  List(String),
			input: cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.NumberIntVal(1)}),
		},
		{
			name:  "empty tuple to list of strings",
			list:  List(String),
			input: cty.EmptyTupleVal,
		},
		{
			name:  "null list of strings",
			list:  List(String),
			input: cty.NullVal(cty.List(cty.String)),
		},
		{
			name:  "list with null values to list of strings",
			list:  List(String),
			input: cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.NullVal(cty.String)}),
		},
		{
			name:  "list of numbers",
			list:  List(Number),
			input: cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NumberFloatVal(3.14)}),
		},
		{
			name:  "list of strings to list of numbers",
			list:  List(Number),
			input: cty.ListVal([]cty.Value{cty.StringVal("42"), cty.StringVal("3.14")}),
		},
		{
			name:  "empty list of numbers",
			list:  List(Number),
			input: cty.ListValEmpty(cty.Number),
		},
		{
			name:  "single number",
			list:  List(Number),
			input: cty.NumberIntVal(123),
		},
		{
			name:  "tuple of number strings and numbers to list of numbers",
			list:  List(Number),
			input: cty.TupleVal([]cty.Value{cty.StringVal("42"), cty.NumberIntVal(1)}),
		},
		{
			name:  "empty tuple to list of numbers",
			list:  List(Number),
			input: cty.EmptyTupleVal,
		},
		{
			name:  "null list of numbers",
			list:  List(Number),
			input: cty.NullVal(cty.List(cty.Number)),
		},
		{
			name:  "list with null values to list of numbers",
			list:  List(Number),
			input: cty.ListVal([]cty.Value{cty.NumberIntVal(42), cty.NullVal(cty.Number)}),
		},
		{
			name:  "list of booleans",
			list:  List(Bool),
			input: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:  "list of strings to list of booleans",
			list:  List(Bool),
			input: cty.ListVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
		},
		{
			name:  "empty list of booleans",
			list:  List(Bool),
			input: cty.ListValEmpty(cty.Bool),
		},
		{
			name:  "single boolean",
			list:  List(Bool),
			input: cty.BoolVal(true),
		},
		{
			name:  "tuple of boolean strings and booleans to list of booleans",
			list:  List(Bool),
			input: cty.TupleVal([]cty.Value{cty.StringVal("true"), cty.BoolVal(false)}),
		},
		{
			name:  "empty tuple to list of booleans",
			list:  List(Bool),
			input: cty.EmptyTupleVal,
		},
		{
			name:  "null list of booleans",
			list:  List(Bool),
			input: cty.NullVal(cty.List(cty.Bool)),
		},
		{
			name:  "list with null values to list of booleans",
			list:  List(Bool),
			input: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.NullVal(cty.Bool)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.list, tt.input)
		})
	}
}

func TestListValidate_InvalidElement(t *testing.T) {
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

func TestListValidate_Nil(t *testing.T) {
	var nilList *listType

	expectedError := "list type is nil"
	err := nilList.Validate(cty.ListValEmpty(cty.String))
	if err == nil {
		t.Fatalf("expected error %q from Validate(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestListValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name  string
		list  Type
		input cty.Value
	}{
		{
			name: "valid list of strings",
			list: List(String),
			input: cty.ListVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
		},
		{
			name: "valid list of objects",
			list: List(Object(RequiredField("foo", String))),
			input: cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"foo": cty.StringVal("bar"),
				}),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.list.ValidateSpec()
			if err != nil {
				t.Errorf("expected no error from ValidateSpec(), got %q", err)
			}
		})
	}
}

func TestListValidateSpec_InvalidObject(t *testing.T) {
	list := List(Object(RequiredField("foo", String)).WithConstraints(MutuallyExclusive("nonexistant")))

	expectedError := `constraint validation failed: field "nonexistant" is not defined in the object type`
	err := list.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestListValidateSpec_Nil(t *testing.T) {
	var nilList *listType

	expectedError := "list type is nil"
	err := nilList.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
