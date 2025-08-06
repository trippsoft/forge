package hcl_function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		list     cty.Value
		value    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "find element in list",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("b"),
			expected: cty.NumberIntVal(1),
			wantErr:  false,
		},
		{
			name:     "find first element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("a"),
			expected: cty.NumberIntVal(0),
			wantErr:  false,
		},
		{
			name:     "find last element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("c"),
			expected: cty.NumberIntVal(2),
			wantErr:  false,
		},
		{
			name:     "find number in list",
			list:     cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			value:    cty.NumberIntVal(2),
			expected: cty.NumberIntVal(1),
			wantErr:  false,
		},
		{
			name:     "find boolean in list",
			list:     cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false), cty.BoolVal(true)}),
			value:    cty.BoolVal(false),
			expected: cty.NumberIntVal(1),
			wantErr:  false,
		},
		{
			name:     "find element in tuple",
			list:     cty.TupleVal([]cty.Value{cty.StringVal("x"), cty.NumberIntVal(42), cty.BoolVal(true)}),
			value:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(1),
			wantErr:  false,
		},
		{
			name:    "value not found",
			list:    cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:   cty.StringVal("d"),
			wantErr: true,
		},
		{
			name:    "empty list",
			list:    cty.ListValEmpty(cty.String),
			value:   cty.StringVal("a"),
			wantErr: true,
		},
		{
			name:     "unknown list",
			list:     cty.UnknownVal(cty.List(cty.String)),
			value:    cty.StringVal("a"),
			expected: cty.UnknownVal(cty.Number),
			wantErr:  false,
		},
		{
			name:     "list with unknown element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.UnknownVal(cty.String), cty.StringVal("c")}),
			value:    cty.StringVal("b"),
			expected: cty.UnknownVal(cty.Number),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Index(tt.list, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.expected.IsKnown() && result.IsKnown() {
					if !result.Equals(tt.expected).True() {
						t.Errorf("Index() = %v, want %v", result, tt.expected)
					}
				} else if !tt.expected.IsKnown() && !result.IsKnown() && result.Type().Equals(tt.expected.Type()) {
					// Both unknown with same type - this is expected
				} else if tt.expected.IsKnown() != result.IsKnown() {
					t.Errorf("Index() known status = %v, want %v", result.IsKnown(), tt.expected.IsKnown())
				}
			}
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    []cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "function call with valid inputs",
			input:    []cty.Value{cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}), cty.StringVal("b")},
			expected: cty.NumberIntVal(1),
			wantErr:  false,
		},
		{
			name:    "function call with not found value",
			input:   []cty.Value{cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}), cty.StringVal("c")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IndexFunc.Call(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !result.Equals(tt.expected).True() {
					t.Errorf("IndexFunc.Call() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
