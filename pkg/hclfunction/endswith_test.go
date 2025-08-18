// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getEndsWithTestCases() []struct {
	name     string
	value    cty.Value
	suffix   cty.Value
	expected cty.Value
} {
	return []struct {
		name     string
		value    cty.Value
		suffix   cty.Value
		expected cty.Value
	}{
		{
			name:     "string ends with suffix",
			value:    cty.StringVal("hello world"),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(true),
		},
		{
			name:     "string does not end with suffix",
			value:    cty.StringVal("hello world"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "empty suffix",
			value:    cty.StringVal("hello"),
			suffix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
		},
		{
			name:     "empty string and empty suffix",
			value:    cty.StringVal(""),
			suffix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
		},
		{
			name:     "empty string with non-empty suffix",
			value:    cty.StringVal(""),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "suffix longer than string",
			value:    cty.StringVal("hi"),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "exact match",
			value:    cty.StringVal("hello"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
		},
		{
			name:     "case sensitive",
			value:    cty.StringVal("Hello"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "unicode characters",
			value:    cty.StringVal("hello 世界"),
			suffix:   cty.StringVal("世界"),
			expected: cty.BoolVal(true),
		},
	}
}

func TestEndsWith(t *testing.T) {
	tests := getEndsWithTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual, err := EndsWith(tt.value, tt.suffix)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestEndsWithFunc(t *testing.T) {
	tests := getEndsWithTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := EndsWithFunc.Call([]cty.Value{tt.value, tt.suffix})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
