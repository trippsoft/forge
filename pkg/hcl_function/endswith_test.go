package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestEndsWith(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		suffix   cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "string ends with suffix",
			value:    cty.StringVal("hello world"),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "string does not end with suffix",
			value:    cty.StringVal("hello world"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "empty suffix",
			value:    cty.StringVal("hello"),
			suffix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "empty string and empty suffix",
			value:    cty.StringVal(""),
			suffix:   cty.StringVal(""),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "empty string with non-empty suffix",
			value:    cty.StringVal(""),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "suffix longer than string",
			value:    cty.StringVal("hi"),
			suffix:   cty.StringVal("world"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "exact match",
			value:    cty.StringVal("hello"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "case sensitive",
			value:    cty.StringVal("Hello"),
			suffix:   cty.StringVal("hello"),
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			value:    cty.StringVal("hello 世界"),
			suffix:   cty.StringVal("世界"),
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EndsWith(tt.value, tt.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndsWith() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("EndsWith() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestEndsWithFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with matching suffix",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("world")},
			expected: cty.BoolVal(true),
			wantErr:  false,
		},
		{
			name:     "function call with non-matching suffix",
			input:    []cty.Value{cty.StringVal("hello world"), cty.StringVal("hello")},
			expected: cty.BoolVal(false),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EndsWithFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndsWithFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("EndsWithFunc.Call() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
