package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestStrContains(t *testing.T) {
	tests := []struct {
		name      string
		value     cty.Value
		substring cty.Value
		expected  cty.Value
		wantErr   bool
	}{
		{
			name:      "string contains substring",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("world"),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "string contains substring at beginning",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "string contains substring in middle",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("lo wo"),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "string does not contain substring",
			value:     cty.StringVal("hello world"),
			substring: cty.StringVal("foo"),
			expected:  cty.BoolVal(false),
			wantErr:   false,
		},
		{
			name:      "empty substring",
			value:     cty.StringVal("hello"),
			substring: cty.StringVal(""),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "empty string and empty substring",
			value:     cty.StringVal(""),
			substring: cty.StringVal(""),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "empty string with non-empty substring",
			value:     cty.StringVal(""),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
			wantErr:   false,
		},
		{
			name:      "substring longer than string",
			value:     cty.StringVal("hi"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
			wantErr:   false,
		},
		{
			name:      "exact match",
			value:     cty.StringVal("hello"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
		{
			name:      "case sensitive",
			value:     cty.StringVal("Hello"),
			substring: cty.StringVal("hello"),
			expected:  cty.BoolVal(false),
			wantErr:   false,
		},
		{
			name:      "unicode characters",
			value:     cty.StringVal("hello 世界 world"),
			substring: cty.StringVal("世界"),
			expected:  cty.BoolVal(true),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StrContains(tt.value, tt.substring)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrContains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("StrContains() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestStrContainsFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with contained substring",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("world")},
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "function call with non-contained substring",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("foo")},
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StrContainsFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrContainsFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("StrContainsFunc.Call() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
