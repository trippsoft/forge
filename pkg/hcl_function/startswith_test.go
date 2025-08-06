package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		prefix   cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "string starts with prefix",
			value:    cty.StringVal("hello world"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "string does not start with prefix",
			value:    cty.StringVal("hello world"),
			prefix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "empty prefix",
			value:    cty.StringVal("hello"),
			prefix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "empty string and empty prefix",
			value:    cty.StringVal(""),
			prefix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "empty string with non-empty prefix",
			value:    cty.StringVal(""),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "prefix longer than string",
			value:    cty.StringVal("hi"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "exact match",
			value:    cty.StringVal("hello"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "case sensitive",
			value:    cty.StringVal("Hello"),
			prefix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			value:    cty.StringVal("世界 hello"),
			prefix:   cty.StringVal("世界"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StartsWith(tt.value, tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartsWith() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("StartsWith() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestStartsWithFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with matching prefix",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("hello")},
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "function call with non-matching prefix",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("world")},
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StartsWithFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartsWithFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("StartsWithFunc.Call() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
