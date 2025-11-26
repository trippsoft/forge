// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getAllTrueTestCases() []struct {
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
			input:    cty.ListVal([]cty.Value{cty.True, cty.False, cty.True}),
			expected: cty.False,
		},
		{
			name:     "all false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.False, cty.False}),
			expected: cty.False,
		},
		{
			name:     "empty list",
			input:    cty.ListValEmpty(cty.Bool),
			expected: cty.True,
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
			name:     "list with null value",
			input:    cty.ListVal([]cty.Value{cty.True, cty.NullVal(cty.Bool), cty.True}),
			expected: cty.False,
		},
		{
			name:     "list with unknown value",
			input:    cty.ListVal([]cty.Value{cty.True, cty.UnknownVal(cty.Bool), cty.True}),
			expected: cty.UnknownVal(cty.Bool),
		},
	}
}

func TestAllTrue(t *testing.T) {
	tests := getAllTrueTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := AllTrue(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestAllTrueFunc(t *testing.T) {
	tests := getAllTrueTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := AllTrueFunc.Call([]cty.Value{tt.input})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
