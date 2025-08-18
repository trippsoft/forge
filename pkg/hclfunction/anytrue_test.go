// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getAnyTrueTestCases() []struct {
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
			name:     "all true values",
			input:    cty.ListVal([]cty.Value{cty.True, cty.True, cty.True}),
			expected: cty.True,
		},
		{
			name:     "mixed true and false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.True, cty.False}),
			expected: cty.True,
		},
		{
			name:     "all false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.False, cty.False}),
			expected: cty.False,
		},
		{
			name:     "empty list",
			input:    cty.ListValEmpty(cty.Bool),
			expected: cty.False,
		},
		{
			name:     "single true value",
			input:    cty.ListVal([]cty.Value{cty.True}),
			expected: cty.True,
		},
		{
			name:     "single false value",
			input:    cty.ListVal([]cty.Value{cty.False}),
			expected: cty.False,
		},
		{
			name:     "list with null values only",
			input:    cty.ListVal([]cty.Value{cty.NullVal(cty.Bool), cty.NullVal(cty.Bool)}),
			expected: cty.False,
		},
		{
			name:     "list with null and true values",
			input:    cty.ListVal([]cty.Value{cty.NullVal(cty.Bool), cty.True, cty.False}),
			expected: cty.True,
		},
		{
			name:     "list with unknown value only",
			input:    cty.ListVal([]cty.Value{cty.UnknownVal(cty.Bool)}),
			expected: cty.UnknownVal(cty.Bool),
		},
		{
			name:     "list with unknown and false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.UnknownVal(cty.Bool), cty.False}),
			expected: cty.UnknownVal(cty.Bool),
		},
		{
			name:     "list with unknown and true values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.UnknownVal(cty.Bool), cty.True}),
			expected: cty.True,
		},
	}
}

func TestAnyTrue(t *testing.T) {
	tests := getAnyTrueTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := AnyTrue(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestAnyTrueFunc(t *testing.T) {
	tests := getAnyTrueTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := AnyTrueFunc.Call([]cty.Value{tt.input})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
