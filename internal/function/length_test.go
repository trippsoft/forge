package function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestLength(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "string length",
			input:    cty.StringVal("hello"),
			expected: cty.NumberIntVal(5),
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    cty.StringVal(""),
			expected: cty.NumberIntVal(0),
			wantErr:  false,
		},
		{
			name:     "unicode string",
			input:    cty.StringVal("hello 世界"),
			expected: cty.NumberIntVal(8),
			wantErr:  false,
		},
		{
			name:     "list length",
			input:    cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			expected: cty.NumberIntVal(3),
			wantErr:  false,
		},
		{
			name:     "empty list",
			input:    cty.ListValEmpty(cty.String),
			expected: cty.NumberIntVal(0),
			wantErr:  false,
		},
		{
			name:     "map length",
			input:    cty.MapVal(map[string]cty.Value{"a": cty.StringVal("1"), "b": cty.StringVal("2")}),
			expected: cty.NumberIntVal(2),
			wantErr:  false,
		},
		{
			name:     "empty map",
			input:    cty.MapValEmpty(cty.String),
			expected: cty.NumberIntVal(0),
			wantErr:  false,
		},
		{
			name:     "set length",
			input:    cty.SetVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
			expected: cty.NumberIntVal(2),
			wantErr:  false,
		},
		{
			name:     "tuple length",
			input:    cty.TupleVal([]cty.Value{cty.StringVal("a"), cty.NumberIntVal(1), cty.BoolVal(true)}),
			expected: cty.NumberIntVal(3),
			wantErr:  false,
		},
		{
			name:     "object length",
			input:    cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("test"), "age": cty.NumberIntVal(25)}),
			expected: cty.NumberIntVal(2),
			wantErr:  false,
		},
		{
			name:     "dynamic type returns unknown",
			input:    cty.UnknownVal(cty.DynamicPseudoType),
			expected: cty.UnknownVal(cty.Number),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Length(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Length() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.expected.IsKnown() && result.IsKnown() {
					if !result.Equals(tt.expected).True() {
						t.Errorf("Length() = %v, want %v", result, tt.expected)
					}
				} else if !tt.expected.IsKnown() && !result.IsKnown() && result.Type().Equals(tt.expected.Type()) {
					// Both unknown with same type - this is expected
				} else if tt.expected.IsKnown() != result.IsKnown() {
					t.Errorf("Length() known status = %v, want %v", result.IsKnown(), tt.expected.IsKnown())
				}
			}
		})
	}
}

func TestLengthFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with string",
			input:    []cty.Value{cty.StringVal("test")},
			expected: cty.NumberIntVal(4),
			wantErr:  false,
		},
		{
			name:     "function call with list",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")})},
			expected: cty.NumberIntVal(2),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LengthFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("LengthFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("LengthFunc.Call() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
