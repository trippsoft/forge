// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getStrContainsTestCases() []struct {
	name      string
	value     cty.Value
	substring cty.Value
	expected  cty.Value
} {
	return []struct {
		name      string
		value     cty.Value
		substring cty.Value
		expected  cty.Value
	}{
		{
			name:      "string contains substring",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("world"),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "string contains substring at beginning",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "string contains substring in middle",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("lo wo"),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "string does not contain substring",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("foo"),
			expected:  cty.BoolVal(false),
		},
		{
			name:      "empty substring",
			value:     cty.StringVal("hello"),
			substring: cty.StringVal(""),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "empty string and empty substring",
			value:     cty.StringVal(""),
			substring: cty.StringVal(""),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "empty string with non-empty substring",
			value:     cty.StringVal(""),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
		},
		{
			name:      "substring longer than string",
			value:     cty.StringVal("hi"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
		},
		{
			name:      "exact match",
			value:     cty.StringVal("hello"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(true),
		},
		{
			name:      "case sensitive",
			value:     cty.StringVal("Hello"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
		},
		{
			name:      "unicode characters",
			value:     cty.StringVal("hello 世界 world"),
			substring: cty.StringVal("世界"),
			expected:  cty.BoolVal(true),
		},
	}
}

func TestStrContains(t *testing.T) {
	tests := getStrContainsTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := StrContains(tt.value, tt.substring)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestStrContainsFunc(t *testing.T) {
	tests := getStrContainsTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := StrContainsFunc.Call([]cty.Value{tt.value, tt.substring})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}
