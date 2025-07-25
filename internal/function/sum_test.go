package function

import (
	"math/big"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestSumFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "sum of positive integers",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
			wantErr:  false,
		},
		{
			name:     "sum of negative integers",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(-1), cty.NumberIntVal(-2), cty.NumberIntVal(-3)}),
			expected: cty.NumberIntVal(-6),
			wantErr:  false,
		},
		{
			name:     "sum of mixed positive and negative",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(5), cty.NumberIntVal(-2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
			wantErr:  false,
		},
		{
			name:     "sum of floats",
			input:    cty.ListVal([]cty.Value{cty.NumberFloatVal(1.5), cty.NumberFloatVal(2.5), cty.NumberFloatVal(3.0)}),
			expected: cty.NumberFloatVal(7.0),
			wantErr:  false,
		},
		{
			name:     "sum of mixed integers and floats",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberFloatVal(2.5), cty.NumberIntVal(3)}),
			expected: cty.NumberFloatVal(6.5),
			wantErr:  false,
		},
		{
			name:     "single element",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(42)}),
			expected: cty.NumberIntVal(42),
			wantErr:  false,
		},
		{
			name:     "sum with zero",
			input:    cty.ListVal([]cty.Value{cty.NumberIntVal(0), cty.NumberIntVal(5), cty.NumberIntVal(0)}),
			expected: cty.NumberIntVal(5),
			wantErr:  false,
		},
		{
			name:     "sum of big numbers",
			input:    cty.ListVal([]cty.Value{cty.NumberVal(big.NewFloat(1e10)), cty.NumberVal(big.NewFloat(2e10))}),
			expected: cty.NumberVal(big.NewFloat(3e10)),
			wantErr:  false,
		},
		{
			name:     "sum of tuple",
			input:    cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
			wantErr:  false,
		},
		{
			name:     "sum of set",
			input:    cty.SetVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			expected: cty.NumberIntVal(6),
			wantErr:  false,
		},
		{
			name:    "empty list",
			input:   cty.ListValEmpty(cty.Number),
			wantErr: true,
		},
		{
			name:    "list with null value",
			input:   cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NullVal(cty.Number), cty.NumberIntVal(3)}),
			wantErr: true,
		},
		{
			name:    "list with non-numeric value as tuple",
			input:   cty.TupleVal([]cty.Value{cty.NumberIntVal(1), cty.StringVal("hello"), cty.NumberIntVal(3)}),
			wantErr: true,
		},
		{
			name:    "non-iterable type",
			input:   cty.StringVal("hello"),
			wantErr: true,
		},
		{
			name:     "unknown list",
			input:    cty.UnknownVal(cty.List(cty.Number)),
			expected: cty.UnknownVal(cty.Number),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SumFunc.Call([]cty.Value{tt.input})
			if (err != nil) != tt.wantErr {
				t.Errorf("SumFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.expected.IsKnown() && result.IsKnown() {
					if !result.Equals(tt.expected).True() {
						t.Errorf("SumFunc.Call() = %v, want %v", result, tt.expected)
					}
				} else if !tt.expected.IsKnown() && !result.IsKnown() && result.Type().Equals(tt.expected.Type()) {
					// Both unknown with same type - this is expected
				} else if tt.expected.IsKnown() != result.IsKnown() {
					t.Errorf("SumFunc.Call() known status = %v, want %v", result.IsKnown(), tt.expected.IsKnown())
				}
			}
		})
	}
}
