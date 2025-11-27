// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"fmt"
	"slices"
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

func TestSetCtyType(t *testing.T) {
	tests := []struct {
		name     string
		set      Type
		expected cty.Type
	}{
		{
			name:     "set of strings",
			set:      Set(String),
			expected: cty.Set(cty.String),
		},
		{
			name:     "set of numbers",
			set:      Set(Number),
			expected: cty.Set(cty.Number),
		},
		{
			name:     "set of booleans",
			set:      Set(Bool),
			expected: cty.Set(cty.Bool),
		},
		{
			name:     "set of durations",
			set:      Set(Duration),
			expected: cty.Set(cty.String),
		},
		{
			name:     "set of set of strings",
			set:      Set(Set(String)),
			expected: cty.Set(cty.Set(cty.String)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.set.CtyType()
			if !actual.Equals(tt.expected) {
				t.Errorf("expected %q from CtyType(), got %q", tt.expected.FriendlyName(), actual.FriendlyName())
			}
		})
	}
}

func TestSetCtyType_Nil(t *testing.T) {
	var nilSet *setType
	if !nilSet.CtyType().Equals(cty.NilType) {
		t.Errorf("expected nil type from CtyType(), got %q", nilSet.CtyType().FriendlyName())
	}
}

func TestSetConvert_Success(t *testing.T) {
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

func TestSetConvert_SensitiveString(t *testing.T) {
	ui.SecretFilter.Clear()
	str := "hello"
	set := Set(SensitiveString)
	input := cty.SetVal([]cty.Value{cty.StringVal(str)})
	verifySuccessfulConversion(t, set, input, input)

	secrets := ui.SecretFilter.Secrets()
	if !slices.Contains(secrets, str) {
		t.Errorf("expected %q to be filtered", str)
	}
}

func TestSetConvert_UnknownValue(t *testing.T) {
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
			verifyFailedConversion(t, tt.set, tt.input, "cannot convert unknown value")
		})
	}
}

func TestSetConvert_InvalidValues(t *testing.T) {
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

func TestSetConvert_Nil(t *testing.T) {
	var nilSet *setType

	expectedError := "set type is nil"
	converted, err := nilSet.Convert(cty.SetValEmpty(cty.String))
	if err == nil {
		t.Fatalf("expected error %q from Convert(), got none", expectedError)
	}

	if converted.Equals(cty.NilVal) != cty.True {
		t.Fatalf("expected nil value from Convert(), got %s", util.FormatCtyValueToString(converted))
	}

	if expectedError != err.Error() {
		t.Errorf("expected error %q from Convert(), got %q", expectedError, err.Error())
	}
}

func TestSetValidate_Pass(t *testing.T) {
	tests := []struct {
		name  string
		set   Type
		input cty.Value
	}{
		{
			name:  "set of strings",
			set:   Set(String),
			input: cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		},
		{
			name:  "empty set of strings",
			set:   Set(String),
			input: cty.SetValEmpty(cty.String),
		},
		{
			name:  "null set of strings",
			set:   Set(String),
			input: cty.NullVal(cty.Set(cty.String)),
		},
		{
			name:  "set of strings with null values",
			set:   Set(String),
			input: cty.SetVal([]cty.Value{cty.StringVal("a"), cty.NullVal(cty.String)}),
		},
		{
			name:  "single string to set of strings",
			set:   Set(String),
			input: cty.StringVal("a"),
		},
		{
			name:  "list of strings to set of strings",
			set:   Set(String),
			input: cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
		},
		{
			name:  "set of numbers to set of strings",
			set:   Set(String),
			input: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:  "set of numbers",
			set:   Set(Number),
			input: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:  "empty set of numbers",
			set:   Set(Number),
			input: cty.SetValEmpty(cty.Number),
		},
		{
			name:  "null set of numbers",
			set:   Set(Number),
			input: cty.NullVal(cty.Set(cty.Number)),
		},
		{
			name:  "set of numbers with null values",
			set:   Set(Number),
			input: cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NullVal(cty.Number)}),
		},
		{
			name:  "single number to set of numbers",
			set:   Set(Number),
			input: cty.NumberIntVal(1),
		},
		{
			name:  "list of numbers to set of numbers",
			set:   Set(Number),
			input: cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5)}),
		},
		{
			name:  "set of number-strings to set of numbers",
			set:   Set(Number),
			input: cty.SetVal([]cty.Value{cty.StringVal("1"), cty.StringVal("2.5")}),
		},
		{
			name:  "set of booleans",
			set:   Set(Bool),
			input: cty.SetVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:  "empty set of booleans",
			set:   Set(Bool),
			input: cty.SetValEmpty(cty.Bool),
		},
		{
			name:  "null set of booleans",
			set:   Set(Bool),
			input: cty.NullVal(cty.Set(cty.Bool)),
		},
		{
			name:  "set of booleans with null values",
			set:   Set(Bool),
			input: cty.SetVal([]cty.Value{cty.BoolVal(true), cty.NullVal(cty.Bool)}),
		},
		{
			name:  "single boolean to set of booleans",
			set:   Set(Bool),
			input: cty.BoolVal(true),
		},
		{
			name:  "list of booleans to set of booleans",
			set:   Set(Bool),
			input: cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false)}),
		},
		{
			name:  "set of boolean-strings to set of booleans",
			set:   Set(Bool),
			input: cty.SetVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifySuccessfulValidation(t, tt.set, tt.input)
		})
	}
}

func TestSetValidate_InvalidElement(t *testing.T) {

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

func TestSetValidate_Nil(t *testing.T) {
	var nilSet *setType

	expectedError := "set type is nil"
	err := nilSet.Validate(cty.SetValEmpty(cty.String))
	if err == nil {
		t.Errorf("expected error %q from Validate(), got nil", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Validate(), got %q", expectedError, err.Error())
	}
}

func TestSetValidateSpec_Pass(t *testing.T) {
	tests := []struct {
		name  string
		set   Type
		input cty.Value
	}{
		{
			name: "valid set of strings",
			set:  Set(String),
			input: cty.SetVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
		},
		{
			name: "valid set of objects",
			set:  Set(Object(RequiredField("foo", String))),
			input: cty.SetVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"foo": cty.StringVal("bar"),
				}),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.set.ValidateSpec()
			if err != nil {
				t.Errorf("expected no error from ValidateSpec(), got %q", err)
			}
		})
	}
}

func TestSetValidateSpec_InvalidObject(t *testing.T) {
	set := Set(Object(RequiredField("foo", String)).WithConstraints(MutuallyExclusive("nonexistant")))

	expectedError := `constraint validation failed: field "nonexistant" is not defined in the object type`
	err := set.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}

func TestSetValidateSpec_Nil(t *testing.T) {
	var nilSet *setType

	expectedError := "set type is nil"
	err := nilSet.ValidateSpec()
	if err == nil {
		t.Errorf("expected error %q from ValidateSpec(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ValidateSpec(), got %q", expectedError, err.Error())
	}
}
