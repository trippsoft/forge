// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"math/big"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getSumTestCases() []struct {
	name     string
	input    cty.Value
	expected cty.Value
} {
	return []struct {
		name     string
		input    cty.Value
		expected cty.Value
	}{
		{
			name:     "sum of positive integers",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
		},
		{
			name:     "sum of negative integers",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(-1), cty.NumberIntVal(-2), cty.NumberIntVal(-3)}),
			expected: cty.NumberIntVal(-6),
		},
		{
			name:     "sum of mixed positive and negative",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(5), cty.NumberIntVal(-2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
		},
		{
			name: "sum of floats",
			input: cty.ListVal([]cty.Value{
				cty.NumberFloatVal(1.5),
				cty.NumberFloatVal(2.5),
				cty.NumberFloatVal(3.0),
			}),
			expected: cty.NumberFloatVal(7.0),
		},
		{
			name:     "sum of mixed integers and floats",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5), cty.NumberIntVal(3)}),
			expected: cty.NumberFloatVal(6.5),
		},
		{
			name:     "single element",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(42)}),
			expected: cty.NumberIntVal(42),
		},
		{
			name:     "sum with zero",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(0), cty.NumberIntVal(5), cty.NumberIntVal(0)}),
			expected: cty.NumberIntVal(5),
		},
		{
			name:     "sum of big numbers",
			input:    cty.ListVal([]cty.Value{cty.NumberVal(big.NewFloat(1e10)), cty.NumberVal(big.NewFloat(2e10))}),
			expected: cty.NumberVal(big.NewFloat(3e10)),
		},
		{
			name:     "sum of tuple",
			input:    cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
		},
		{
			name:     "sum of set",
			input:    cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
		},
		{
			name:     "unknown list",
			input:    cty.UnknownVal(cty.List(cty.Number)),
			expected: cty.UnknownVal(cty.Number),
		},
	}
}

func TestSum(t *testing.T) {
	tests := getSumTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Sum(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestSum_EmptyList(t *testing.T) {
	// Test case for an empty list
	input := cty.ListValEmpty(cty.Number)

	_, err := Sum(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a non-empty iterable"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSum_ListWithNull(t *testing.T) {
	// Test case for a list with null values
	input := cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NullVal(cty.Number), cty.NumberIntVal(2)})

	_, err := Sum(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a list, set, or tuple of numbers, got null"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSum_TupleWithNonNumber(t *testing.T) {
	// Test case for a tuple with non-number values
	input := cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.StringVal("not_a_number"), cty.NumberIntVal(2)})

	_, err := Sum(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a list, set, or tuple of numbers"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSum_NonIterable(t *testing.T) {
	// Test case for a non-iterable input
	input := cty.StringVal("not_an_iterable")

	_, err := Sum(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires an iterable type"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSumFunc(t *testing.T) {
	tests := getSumTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := SumFunc.Call([]cty.Value{tt.input})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestSumFunc_EmptyList(t *testing.T) {
	// Test case for an empty list
	input := cty.ListValEmpty(cty.Number)

	_, err := SumFunc.Call([]cty.Value{input})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a non-empty iterable"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSumFunc_ListWithNull(t *testing.T) {
	// Test case for a list with null values
	input := cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NullVal(cty.Number), cty.NumberIntVal(2)})

	_, err := SumFunc.Call([]cty.Value{input})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a list, set, or tuple of numbers, got null"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSumFunc_TupleWithNonNumber(t *testing.T) {
	// Test case for a tuple with non-number values
	input := cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.StringVal("not_a_number"), cty.NumberIntVal(2)})

	_, err := SumFunc.Call([]cty.Value{input})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires a list, set, or tuple of numbers"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSumFunc_NonIterable(t *testing.T) {
	// Test case for a non-iterable input
	input := cty.StringVal("not_an_iterable")

	_, err := SumFunc.Call([]cty.Value{input})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "sum function requires an iterable type"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}
