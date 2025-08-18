// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getStartsWithTestCases() []struct {
	name     string
	value    cty.Value
	prefix   cty.Value
	expected cty.Value
} {
	return []struct {
		name     string
		value    cty.Value
		prefix   cty.Value
		expected cty.Value
	}{
		{
			name:     "string starts with prefix",
			value:    cty.StringVal("hello world"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
		},
		{
			name:     "string does not start with prefix",
			value:    cty.StringVal("hello world"),
			prefix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "empty prefix",
			value:    cty.StringVal("hello"),
			prefix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
		},
		{
			name:     "empty string and empty prefix",
			value:    cty.StringVal(""),
			prefix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
		},
		{
			name:     "empty string with non-empty prefix",
			value:    cty.StringVal(""),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "prefix longer than string",
			value:    cty.StringVal("hi"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "exact match",
			value:    cty.StringVal("hello"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
		},
		{
			name:     "case sensitive",
			value:    cty.StringVal("Hello"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
		},
		{
			name:     "unicode characters",
			value:    cty.StringVal("世界 hello"),
			prefix:   cty.StringVal("世界"),
			expected: cty.BoolVal(true),
		},
	}
}

func TestStartsWith(t *testing.T) {
	tests := getStartsWithTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := StartsWith(tt.value, tt.prefix)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestStartsWithFunc(t *testing.T) {
	tests := getStartsWithTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := StartsWithFunc.Call([]cty.Value{tt.value, tt.prefix})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
