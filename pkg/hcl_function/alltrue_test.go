package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestAllTrue(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "all true values",
			input:    cty.ListVal([]cty.Value{cty.True, cty.True, cty.True}),
			expected: cty.True,
			wantErr:  false,
		},
		{
			name:     "mixed true and false values",
			input:    cty.ListVal([]cty.Value{cty.True, cty.False, cty.True}),
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "all false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.False, cty.False}),
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "empty list",
			input:    cty.ListValEmpty(cty.Bool),
			expected: cty.True,
			wantErr:  false,
		},
		{
			name:     "single true value",
			input:    cty.ListVal([]cty.Value{cty.True}),
			expected: cty.True,
			wantErr:  false,
		},
		{
			name:     "single false value",
			input:    cty.ListVal([]cty.Value{cty.False}),
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "list with null value",
			input:    cty.ListVal([]cty.Value{cty.True, cty.NullVal(cty.Bool), cty.True}),
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "list with unknown value",
			input:    cty.ListVal([]cty.Value{cty.True, cty.UnknownVal(cty.Bool), cty.True}),
			expected: cty.UnknownVal(cty.Bool),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AllTrue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllTrue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assertCtyValueEqual(t, result, tt.expected, "AllTrue()")
			}
		})
	}
}

func TestAllTrueFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with all true values",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.True, cty.True})},
			expected: cty.True,
			wantErr:  false,
		},
		{
			name:     "function call with mixed values",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.True, cty.False})},
			expected: cty.False,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AllTrueFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllTrueFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assertCtyValueEqual(t, result, tt.expected, "AllTrueFunc.Call()")
			}
		})
	}
}
