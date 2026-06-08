// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclspec

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFieldNotPresentIsMet(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		input    cty.Value
		expected bool
	}{
		{
			name:  "field present and not null",
			field: "field1",
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.StringVal("value1"),
			}),
			expected: true,
		},
		{
			name:     "field not in values",
			field:    "field1",
			input:    cty.ObjectVal(map[string]cty.Value{}),
			expected: false,
		},
		{
			name:  "different field in values",
			field: "field1",
			input: cty.ObjectVal(map[string]cty.Value{
				"field2": cty.StringVal("value2"),
			}),
			expected: false,
		},
		{
			name:  "field is null",
			field: "field1",
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.NullVal(cty.String),
			}),
			expected: false,
		},
		{
			name:  "field is unknown",
			field: "field1",
			input: cty.ObjectVal(map[string]cty.Value{
				"field1": cty.UnknownVal(cty.String),
			}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := FieldPresent(tt.field)

			actual := condition.IsMet(tt.input)
			if actual != tt.expected {
				t.Errorf("expected IsMet() to return %t, got %t", tt.expected, actual)
			}
		})
	}
}

func TestFieldNotPresentIsMet_Nil(t *testing.T) {
	var condition *fieldNotPresentCondition

	actual := condition.IsMet(cty.ObjectVal(map[string]cty.Value{}))
	if actual != false {
		t.Errorf("expected IsMet() to return %t, got %t", false, actual)
	}
}

func TestFieldNotPresentDescription(t *testing.T) {
	tests := []struct {
		name      string
		condition TypeCondition
		expected  string
	}{
		{
			name:      "field1",
			condition: FieldNotPresent("field1"),
			expected:  `field "field1" is not present`,
		},
		{
			name:      "field2",
			condition: FieldNotPresent("field2"),
			expected:  `field "field2" is not present`,
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

func TestFieldNotPresentDescription_Nil(t *testing.T) {
	var condition *fieldNotPresentCondition
	actual := condition.Description()
	if actual != "" {
		t.Errorf("expected empty string from Description(), got %q", actual)
	}
}
