// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getLengthTestCases() []struct {
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
			name:     "string length",
			input:    cty.StringVal("hello"),
			expected: cty.NumberIntVal(5),
		},
		{
			name:     "empty string",
			input:    cty.StringVal(""),
			expected: cty.NumberIntVal(0),
		},
		{
			name:     "unicode string",
			input:    cty.StringVal("hello 世界"),
			expected: cty.NumberIntVal(8),
		},
		{
			name:     "list length",
			input:    cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			expected: cty.NumberIntVal(3),
		},
		{
			name:     "empty list",
			input:    cty.ListValEmpty(cty.String),
			expected: cty.NumberIntVal(0),
		},
		{
			name:     "map length",
			input:    cty.MapVal(map[string]cty.Value{"a": cty.StringVal("1"), "b": cty.StringVal("2")}),
			expected: cty.NumberIntVal(2),
		},
		{
			name:     "empty map",
			input:    cty.MapValEmpty(cty.String),
			expected: cty.NumberIntVal(0),
		},
		{
			name:     "set length",
			input:    cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
			expected: cty.NumberIntVal(2),
		},
		{
			name:     "tuple length",
			input:    cty.TupleVal([]cty.Value{cty.StringVal("a"), cty.NumberIntVal(1), cty.BoolVal(true)}),
			expected: cty.NumberIntVal(3),
		},
		{
			name:     "object length",
			input:    cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("test"), "age": cty.NumberIntVal(25)}),
			expected: cty.NumberIntVal(2),
		},
		{
			name:     "dynamic type returns unknown",
			input:    cty.UnknownVal(cty.DynamicPseudoType),
			expected: cty.UnknownVal(cty.Number),
		},
	}
}

func TestLength(t *testing.T) {
	tests := getLengthTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Length(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestLengthFunc(t *testing.T) {
	tests := getLengthTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := LengthFunc.Call([]cty.Value{tt.input})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
