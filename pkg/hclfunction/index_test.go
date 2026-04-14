// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getIndexTestCases() []struct {
	name     string
	list     cty.Value
	value    cty.Value
	expected cty.Value
} {
	return []struct {
		name     string
		list     cty.Value
		value    cty.Value
		expected cty.Value
	}{
		{
			name:     "find element in list",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("b"),
			expected: cty.NumberIntVal(1),
		},
		{
			name:     "find first element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("a"),
			expected: cty.NumberIntVal(0),
		},
		{
			name:     "find last element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.StringVal("c"),
			expected: cty.NumberIntVal(2),
		},
		{
			name:     "find number in list",
			list:     cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			value:    cty.NumberIntVal(2),
			expected: cty.NumberIntVal(1),
		},
		{
			name:     "find boolean in list",
			list:     cty.ListVal([]cty.Value{cty.BoolVal(true), cty.BoolVal(false), cty.BoolVal(true)}),
			value:    cty.BoolVal(false),
			expected: cty.NumberIntVal(1),
		},
		{
			name:     "find element in tuple",
			list:     cty.TupleVal([]cty.Value{cty.StringVal("x"), cty.NumberIntVal(42), cty.BoolVal(true)}),
			value:    cty.NumberIntVal(42),
			expected: cty.NumberIntVal(1),
		},
		{
			name:     "unknown list",
			list:     cty.UnknownVal(cty.List(cty.String)),
			value:    cty.StringVal("a"),
			expected: cty.UnknownVal(cty.Number),
		},
		{
			name:     "list with unknown element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.UnknownVal(cty.String), cty.StringVal("c")}),
			value:    cty.StringVal("b"),
			expected: cty.UnknownVal(cty.Number),
		},
		{
			name:     "unknown value with list containing unknown element",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.UnknownVal(cty.String), cty.StringVal("c")}),
			value:    cty.UnknownVal(cty.String),
			expected: cty.UnknownVal(cty.Number),
		},
		{
			name:     "uknown value with list containing all known elements",
			list:     cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")}),
			value:    cty.UnknownVal(cty.String),
			expected: cty.UnknownVal(cty.Number),
		},
	}
}

func TestIndex(t *testing.T) {
	tests := getIndexTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Index(tt.list, tt.value)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestIndex_ValueNotFound(t *testing.T) {
	list := cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")})
	value := cty.StringVal("not_in_list")

	_, err := Index(list, value)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedErr := "index failed: value not found in the list or tuple"
	if err.Error() != expectedErr {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestIndex_EmptyList(t *testing.T) {
	list := cty.ListValEmpty(cty.String)
	value := cty.StringVal("a")

	_, err := Index(list, value)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedErr := "index failed: requires a non-empty list"
	if err.Error() != expectedErr {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestIndexFunc(t *testing.T) {
	tests := getIndexTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := IndexFunc.Call([]cty.Value{tt.list, tt.value})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestIndexFunc_ValueNotFound(t *testing.T) {
	list := cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b"), cty.StringVal("c")})
	value := cty.StringVal("not_in_list")

	_, err := IndexFunc.Call([]cty.Value{list, value})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedErr := "index failed: value not found in the list or tuple"
	if err.Error() != expectedErr {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestIndexFunc_EmptyList(t *testing.T) {
	list := cty.ListValEmpty(cty.String)
	value := cty.StringVal("a")

	_, err := IndexFunc.Call([]cty.Value{list, value})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedErr := "index failed: requires a non-empty list"
	if err.Error() != expectedErr {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
