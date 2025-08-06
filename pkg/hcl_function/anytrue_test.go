package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestAnyTrue(t *testing.T) {
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
			input:    cty.ListVal([]cty.Value{cty.False, cty.True, cty.False}),
			expected: cty.True,
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
			expected: cty.False,
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
			name:     "list with null values only",
			input:    cty.ListVal([]cty.Value{cty.NullVal(cty.Bool), cty.NullVal(cty.Bool)}),
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "list with null and true values",
			input:    cty.ListVal([]cty.Value{cty.NullVal(cty.Bool), cty.True, cty.False}),
			expected: cty.True,
			wantErr:  false,
		},
		{
			name:     "list with unknown value only",
			input:    cty.ListVal([]cty.Value{cty.UnknownVal(cty.Bool)}),
			expected: cty.UnknownVal(cty.Bool),
			wantErr:  false,
		},
		{
			name:     "list with unknown and false values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.UnknownVal(cty.Bool), cty.False}),
			expected: cty.UnknownVal(cty.Bool),
			wantErr:  false,
		},
		{
			name:     "list with unknown and true values",
			input:    cty.ListVal([]cty.Value{cty.False, cty.UnknownVal(cty.Bool), cty.True}),
			expected: cty.True,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnyTrue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyTrue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.expected.IsKnown() && result.IsKnown() {
					if !result.Equals(tt.expected).True() {
						t.Errorf("AnyTrue() = %v, want %v", result, tt.expected)
					}
				} else if !tt.expected.IsKnown() && !result.IsKnown() && result.Type().Equals(tt.expected.Type()) {
					// Both unknown with same type - this is expected
				} else if tt.expected.IsKnown() != result.IsKnown() {
					t.Errorf("AnyTrue() known status = %v, want %v", result.IsKnown(), tt.expected.IsKnown())
				}
			}
		})
	}
}

func TestAnyTrueFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with all false values",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.False, cty.False})},
			expected: cty.False,
			wantErr:  false,
		},
		{
			name:     "function call with mixed values",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.False, cty.True})},
			expected: cty.True,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AnyTrueFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnyTrueFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.expected.IsKnown() && result.IsKnown() {
					if !result.Equals(tt.expected).True() {
						t.Errorf("AnyTrueFunc.Call() = %v, want %v", result, tt.expected)
					}
				} else if !tt.expected.IsKnown() && !result.IsKnown() && result.Type().Equals(tt.expected.Type()) {
					// Both unknown with same type - this is expected
				} else if tt.expected.IsKnown() != result.IsKnown() {
					t.Errorf("AnyTrueFunc.Call() known status = %v, want %v", result.IsKnown(), tt.expected.IsKnown())
				}
			}
		})
	}
}
