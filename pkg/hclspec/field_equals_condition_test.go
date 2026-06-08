// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFieldEqualsIsMet(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    cty.Value
		input    cty.Value
		expected bool
	}{
		{
			name:  "field equals string value",
			field: "field1",
			value: cty.StringVal("expected"),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("expected"),
			}),
			expected: true,
		},
		{
			name:  "field equals number value",
			field: "field1",
			value: cty.NumberIntVal(42),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NumberIntVal(42),
			}),
			expected: true,
		},
		{
			name:  "field equals bool value",
			field: "field1",
			value: cty.BoolVal(true),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.BoolVal(true),
			}),
			expected: true,
		},
		{
			name:     "field not in values",
			field:    "field1",
			value:    cty.StringVal("expected"),
			input:    cty.ObjectVal(map[string]cty.Value{}),
			expected: false,
		},
		{
			name:  "field has different string value",
			field: "field1",
			value: cty.StringVal("expected"),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("actual"),
			}),
			expected: false,
		},
		{
			name:  "field has different type",
			field: "field1",
			value: cty.StringVal("expected"),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NumberIntVal(42),
			}),
			expected: false,
		},
		{
			name:  "field is null",
			field: "field1",
			value: cty.StringVal("expected"),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			}),
			expected: false,
		},
		{
			name:  "field is unknown",
			field: "field1",
			value: cty.StringVal("expected"),
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
			}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldEquals(tt.field, tt.value)

			actual := condition.IsMet(tt.input)
			if actual != tt.expected {
				t.Errorf("expected IsMet() to return %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestFieldEqualsIsMet_Nil(t *testing.T) {
	var condition *fieldEqualsCondition

	actual := condition.IsMet(cty.ObjectVal(map[string]cty.Value{}))
	if actual != false {
		t.Errorf("expected IsMet() to return %t, got %t", false, actual)
	}
}

func TestFieldEqualsDescription(t *testing.T) {
	tests := []struct {
		name      string
		condition TypeCondition
		expected  string
	}{
		{
			name:      "field1",
			condition: FieldEquals("field1", cty.StringVal("value")),
			expected:  `field "field1" is equal to "value"`,
		},
		{
			name:      "field2",
			condition: FieldEquals("field2", cty.StringVal("value2")),
			expected:  `field "field2" is equal to "value2"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.condition.Description()
			if actual != tt.expected {
				t.Errorf("expected Description() to return %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFieldEqualsDescription_Nil(t *testing.T) {
	var condition *fieldEqualsCondition

	actual := condition.Description()
	if actual != "" {
		t.Errorf("expected empty string from Description(), got %q", actual)
	}
}
