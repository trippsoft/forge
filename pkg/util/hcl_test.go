// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"testing"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

func TestConvertHCLAttributeToString(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected string
		hasError bool
	}{
		{
			name:     "valid string",
			value:    cty.StringVal("hello world"),
			expected: "hello world",
			hasError: false,
		},
		{
			name:     "empty string",
			value:    cty.StringVal(""),
			expected: "",
			hasError: false,
		},
		{
			name:     "string with special characters",
			value:    cty.StringVal("hello\nworld\t!@#$"),
			expected: "hello\nworld\t!@#$",
			hasError: false,
		},
		{
			name:     "unicode string",
			value:    cty.StringVal("héllo 世界"),
			expected: "héllo 世界",
			hasError: false,
		},
		{
			name:     "number value should error",
			value:    cty.NumberIntVal(123),
			expected: "",
			hasError: true,
		},
		{
			name:     "bool value should error",
			value:    cty.BoolVal(true),
			expected: "",
			hasError: true,
		},
		{
			name:     "null value should error",
			value:    cty.NullVal(cty.String),
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &mockExpr{value: tt.value}
			attr := &hcl.Attribute{
				Name: "test_attr",
				Expr: expr,
			}
			evalCtx := &hcl.EvalContext{}

			result, diags := ConvertHCLAttributeToString(attr, evalCtx)
			if tt.hasError {
				if !diags.HasErrors() {
					t.Errorf("expected error but got none")
				}
			} else {
				if diags.HasErrors() {
					t.Errorf("unexpected error: %v", diags)
				}
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func TestConvertHCLAttributeToUint16(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected uint16
		hasError bool
	}{
		{
			name:     "valid uint16",
			value:    cty.NumberIntVal(12345),
			expected: 12345,
			hasError: false,
		},
		{
			name:     "zero value",
			value:    cty.NumberIntVal(0),
			expected: 0,
			hasError: false,
		},
		{
			name:     "max uint16 value",
			value:    cty.NumberIntVal(65535),
			expected: 65535,
			hasError: false,
		},
		{
			name:     "overflow uint16",
			value:    cty.NumberIntVal(65536),
			expected: 0,
			hasError: true,
		},
		{
			name:     "negative number",
			value:    cty.NumberIntVal(-1),
			expected: 0,
			hasError: true,
		},
		{
			name:     "string value should error",
			value:    cty.StringVal("123"),
			expected: 0,
			hasError: true,
		},
		{
			name:     "bool value should error",
			value:    cty.BoolVal(true),
			expected: 0,
			hasError: true,
		},
		{
			name:     "null value should error",
			value:    cty.NullVal(cty.Number),
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock attribute and eval context
			expr := &mockExpr{value: tt.value}
			attr := &hcl.Attribute{
				Name: "test_attr",
				Expr: expr,
			}
			evalCtx := &hcl.EvalContext{}

			result, diags := ConvertHCLAttributeToUint16(attr, evalCtx)

			if tt.hasError {
				if !diags.HasErrors() {
					t.Errorf("expected error but got none")
				}
			} else {
				if diags.HasErrors() {
					t.Errorf("unexpected error: %v", diags)
				}
				if result != tt.expected {
					t.Errorf("expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestConvertHCLAttributeToBool(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected bool
		hasError bool
	}{
		{
			name:     "true value",
			value:    cty.BoolVal(true),
			expected: true,
			hasError: false,
		},
		{
			name:     "false value",
			value:    cty.BoolVal(false),
			expected: false,
			hasError: false,
		},
		{
			name:     "string value should error",
			value:    cty.StringVal("true"),
			expected: false,
			hasError: true,
		},
		{
			name:     "number value should error",
			value:    cty.NumberIntVal(1),
			expected: false,
			hasError: true,
		},
		{
			name:     "null value should error",
			value:    cty.NullVal(cty.Bool),
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &mockExpr{value: tt.value}
			attr := &hcl.Attribute{
				Name: "test_attr",
				Expr: expr,
			}
			evalCtx := &hcl.EvalContext{}

			result, diags := ConvertHCLAttributeToBool(attr, evalCtx)
			if tt.hasError {
				if !diags.HasErrors() {
					t.Errorf("expected error but got none")
				}
			} else {
				if diags.HasErrors() {
					t.Errorf("unexpected error: %v", diags)
				}
				if result != tt.expected {
					t.Errorf("expected %t, got %t", tt.expected, result)
				}
			}
		})
	}
}

func TestConvertHCLAttributeToDuration(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected time.Duration
		hasError bool
	}{
		{
			name:     "valid duration - seconds",
			value:    cty.StringVal("30s"),
			expected: 30 * time.Second,
			hasError: false,
		},
		{
			name:     "valid duration - minutes",
			value:    cty.StringVal("5m"),
			expected: 5 * time.Minute,
			hasError: false,
		},
		{
			name:     "valid duration - hours",
			value:    cty.StringVal("2h"),
			expected: 2 * time.Hour,
			hasError: false,
		},
		{
			name:     "valid duration - complex",
			value:    cty.StringVal("1h30m45s"),
			expected: 1*time.Hour + 30*time.Minute + 45*time.Second,
			hasError: false,
		},
		{
			name:     "valid duration - milliseconds",
			value:    cty.StringVal("500ms"),
			expected: 500 * time.Millisecond,
			hasError: false,
		},
		{
			name:     "valid duration - microseconds",
			value:    cty.StringVal("100µs"),
			expected: 100 * time.Microsecond,
			hasError: false,
		},
		{
			name:     "valid duration - nanoseconds",
			value:    cty.StringVal("1000ns"),
			expected: 1000 * time.Nanosecond,
			hasError: false,
		},
		{
			name:     "zero duration",
			value:    cty.StringVal("0s"),
			expected: 0,
			hasError: false,
		},
		{
			name:     "invalid duration string",
			value:    cty.StringVal("invalid"),
			expected: 0,
			hasError: true,
		},
		{
			name:     "empty string",
			value:    cty.StringVal(""),
			expected: 0,
			hasError: true,
		},
		{
			name:     "number value should error",
			value:    cty.NumberIntVal(30),
			expected: 0,
			hasError: true,
		},
		{
			name:     "bool value should error",
			value:    cty.BoolVal(true),
			expected: 0,
			hasError: true,
		},
		{
			name:     "null value should error",
			value:    cty.NullVal(cty.String),
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &mockExpr{value: tt.value}
			attr := &hcl.Attribute{
				Name: "test_attr",
				Expr: expr,
			}
			evalCtx := &hcl.EvalContext{}

			result, diags := ConvertHCLAttributeToDuration(attr, evalCtx)
			if tt.hasError {
				if !diags.HasErrors() {
					t.Errorf("expected error but got none")
				}
			} else {
				if diags.HasErrors() {
					t.Errorf("unexpected error: %v", diags)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestConvertHCLAttributeWithExpressionError(t *testing.T) {
	expr := &mockExpr{
		value: cty.StringVal("test"),
		err:   hcl.Diagnostics{&hcl.Diagnostic{Severity: hcl.DiagError, Summary: "test error"}},
	}
	attr := &hcl.Attribute{
		Name: "test_attr",
		Expr: expr,
	}
	evalCtx := &hcl.EvalContext{}

	t.Run("string conversion with expr error", func(t *testing.T) {
		_, diags := ConvertHCLAttributeToString(attr, evalCtx)
		if !diags.HasErrors() {
			t.Errorf("expected error from expression evaluation")
		}
	})

	t.Run("uint16 conversion with expr error", func(t *testing.T) {
		_, diags := ConvertHCLAttributeToUint16(attr, evalCtx)
		if !diags.HasErrors() {
			t.Errorf("expected error from expression evaluation")
		}
	})

	t.Run("bool conversion with expr error", func(t *testing.T) {
		_, diags := ConvertHCLAttributeToBool(attr, evalCtx)
		if !diags.HasErrors() {
			t.Errorf("expected error from expression evaluation")
		}
	})

	t.Run("duration conversion with expr error", func(t *testing.T) {
		_, diags := ConvertHCLAttributeToDuration(attr, evalCtx)
		if !diags.HasErrors() {
			t.Errorf("expected error from expression evaluation")
		}
	})
}

func TestGetAllCtyStrings(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected []string
	}{
		{
			name:     "null value",
			value:    cty.NullVal(cty.String),
			expected: []string{},
		},
		{
			name:     "unknown value",
			value:    cty.UnknownVal(cty.String),
			expected: []string{},
		},
		{
			name:     "empty string",
			value:    cty.StringVal(""),
			expected: []string{},
		},
		{
			name:     "non-empty string",
			value:    cty.StringVal("hello"),
			expected: []string{"hello"},
		},
		{
			name:     "string with spaces",
			value:    cty.StringVal("hello world"),
			expected: []string{"hello world"},
		},
		{
			name:     "string with special characters",
			value:    cty.StringVal("hello\nworld\t!@#$"),
			expected: []string{"hello\nworld\t!@#$"},
		},
		{
			name:     "unicode string",
			value:    cty.StringVal("héllo 世界"),
			expected: []string{"héllo 世界"},
		},
		{
			name:     "non-string primitive (number)",
			value:    cty.NumberIntVal(123),
			expected: []string{},
		},
		{
			name:     "non-string primitive (bool)",
			value:    cty.BoolVal(true),
			expected: []string{},
		},
		{
			name:     "empty list",
			value:    cty.ListValEmpty(cty.String),
			expected: []string{},
		},
		{
			name:     "list with strings",
			value:    cty.ListVal([]cty.Value{cty.StringVal("one"), cty.StringVal("two"), cty.StringVal("three")}),
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "tuple with mixed types (strings and numbers)",
			value:    cty.TupleVal([]cty.Value{cty.StringVal("hello"), cty.NumberIntVal(42), cty.StringVal("world")}),
			expected: []string{"hello", "world"},
		},
		{
			name:     "list with empty strings",
			value:    cty.ListVal([]cty.Value{cty.StringVal("hello"), cty.StringVal(""), cty.StringVal("world")}),
			expected: []string{"hello", "world"},
		},
		{
			name: "tuple with nested list",
			value: cty.TupleVal([]cty.Value{
				cty.ListVal([]cty.Value{cty.StringVal("nested1"), cty.StringVal("nested2")}),
				cty.StringVal("top"),
			}),
			expected: []string{"nested1", "nested2", "top"},
		},
		{
			name:     "empty set",
			value:    cty.SetValEmpty(cty.String),
			expected: []string{},
		},
		{
			name:     "set with strings",
			value:    cty.SetVal([]cty.Value{cty.StringVal("alpha"), cty.StringVal("beta"), cty.StringVal("gamma")}),
			expected: []string{"alpha", "beta", "gamma"},
		},
		{
			name:     "empty tuple",
			value:    cty.TupleVal([]cty.Value{}),
			expected: []string{},
		},
		{
			name:     "tuple with mixed types",
			value:    cty.TupleVal([]cty.Value{cty.StringVal("first"), cty.NumberIntVal(2), cty.StringVal("third")}),
			expected: []string{"first", "third"},
		},
		{
			name:     "empty map",
			value:    cty.MapValEmpty(cty.String),
			expected: []string{},
		},
		{
			name: "map with string values",
			value: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal("value2"),
				"key3": cty.StringVal("value3"),
			}),
			expected: []string{"value1", "value2", "value3"},
		},
		{
			name: "object with mixed value types",
			value: cty.ObjectVal(map[string]cty.Value{
				"str":  cty.StringVal("hello"),
				"num":  cty.NumberIntVal(123),
				"bool": cty.BoolVal(false),
			}),
			expected: []string{"hello"},
		},
		{
			name: "map with empty string values",
			value: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("value1"),
				"key2": cty.StringVal(""),
				"key3": cty.StringVal("value3"),
			}),
			expected: []string{"value1", "value3"},
		},
		{
			name: "object with nested map",
			value: cty.ObjectVal(map[string]cty.Value{
				"outer": cty.MapVal(map[string]cty.Value{
					"inner1": cty.StringVal("nested1"),
					"inner2": cty.StringVal("nested2"),
				}),
				"simple": cty.StringVal("top"),
			}),
			expected: []string{"nested1", "nested2", "top"},
		},
		{
			name:     "empty object",
			value:    cty.EmptyObjectVal,
			expected: []string{},
		},
		{
			name: "object with string attributes",
			value: cty.ObjectVal(map[string]cty.Value{
				"name":        cty.StringVal("test"),
				"description": cty.StringVal("a test object"),
				"version":     cty.StringVal("1.0.0"),
			}),
			expected: []string{"test", "a test object", "1.0.0"},
		},
		{
			name: "object with mixed attribute types",
			value: cty.ObjectVal(map[string]cty.Value{
				"name":    cty.StringVal("test"),
				"count":   cty.NumberIntVal(5),
				"enabled": cty.BoolVal(true),
			}),
			expected: []string{"test"},
		},
		{
			name: "object with empty string attributes",
			value: cty.ObjectVal(map[string]cty.Value{
				"name":        cty.StringVal("test"),
				"description": cty.StringVal(""),
				"version":     cty.StringVal("1.0.0"),
			}),
			expected: []string{"test", "1.0.0"},
		},
		{
			name: "nested object",
			value: cty.ObjectVal(map[string]cty.Value{
				"metadata": cty.ObjectVal(map[string]cty.Value{
					"name":    cty.StringVal("inner"),
					"version": cty.StringVal("2.0.0"),
				}),
				"title": cty.StringVal("outer"),
			}),
			expected: []string{"inner", "2.0.0", "outer"},
		},
		{
			name: "complex nested structure",
			value: cty.ObjectVal(map[string]cty.Value{
				"config": cty.ObjectVal(map[string]cty.Value{
					"servers": cty.ListVal([]cty.Value{
						cty.ObjectVal(map[string]cty.Value{
							"name": cty.StringVal("server1"),
							"port": cty.NumberIntVal(8080),
						}),
						cty.ObjectVal(map[string]cty.Value{
							"name": cty.StringVal("server2"),
							"port": cty.NumberIntVal(8081),
						}),
					}),
					"database": cty.ObjectVal(map[string]cty.Value{
						"host":     cty.StringVal("localhost"),
						"port":     cty.NumberIntVal(5432),
						"database": cty.StringVal("mydb"),
					}),
				}),
				"environment": cty.StringVal("production"),
			}),
			expected: []string{"server1", "server2", "localhost", "mydb", "production"},
		},
		{
			name: "tuple containing objects and maps",
			value: cty.TupleVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"type": cty.StringVal("object"),
					"id":   cty.NumberIntVal(1),
				}),
				cty.MapVal(map[string]cty.Value{
					"type":  cty.StringVal("map"),
					"count": cty.StringVal("2"), // Changed to string to make map homogeneous
				}),
				cty.StringVal("simple"),
			}),
			expected: []string{"object", "map", "2", "simple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAllCtyStrings(tt.value)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d strings, got %d: expected %v, got %v", len(tt.expected), len(result), tt.expected, result)
				return
			}

			expectedMap := make(map[string]int)
			for _, s := range tt.expected {
				expectedMap[s]++
			}

			resultMap := make(map[string]int)
			for _, s := range result {
				resultMap[s]++
			}

			for expected, count := range expectedMap {
				if resultMap[expected] != count {
					t.Errorf("expected string %q to appear %d times, but it appeared %d times", expected, count, resultMap[expected])
				}
			}

			for result, count := range resultMap {
				if expectedMap[result] != count {
					t.Errorf("unexpected string %q appeared %d times", result, count)
				}
			}
		})
	}
}

func TestModifyUnexpectedElementDiags(t *testing.T) {
	tests := []struct {
		name     string
		diags    hcl.Diagnostics
		location string
		expected hcl.Diagnostics
	}{
		{
			name:     "nil",
			diags:    nil,
			location: "in test file",
			expected: nil,
		},
		{
			name:     "no diagnostics",
			diags:    hcl.Diagnostics{},
			location: "in test file",
			expected: hcl.Diagnostics{},
		},
		{
			name: "non-matching diagnostic",
			diags: hcl.Diagnostics{&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "original error",
				Detail:   "original error",
			}},
			location: "in test file",
			expected: hcl.Diagnostics{&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "original error",
				Detail:   "original error",
			}},
		},
		{
			name: "matching diagnostics with location",
			diags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"test\" is not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"test\" are not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unexpected \"test\" block",
					Detail:   "Blocks are not allowed here.",
				},
			},
			location: "in a test file",
			expected: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"test\" is not expected in a test file.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"test\" are not expected in a test file.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unexpected \"test\" block",
					Detail:   "Blocks are not allowed in a test file.",
				},
			},
		},
		{
			name: "matching diagnostics with no location",
			diags: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"test\" is not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"test\" are not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unexpected \"test\" block",
					Detail:   "Blocks are not allowed here.",
				},
			},
			location: "",
			expected: hcl.Diagnostics{
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unsupported argument",
					Detail:   "An argument named \"test\" is not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unsupported block type",
					Detail:   "Blocks of type \"test\" are not expected here.",
				},
				&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Unexpected \"test\" block",
					Detail:   "Blocks are not allowed here.",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ModifyUnexpectedElementDiags(tt.diags, tt.location)
			for i := range tt.diags {
				if tt.diags[i].Severity != tt.expected[i].Severity ||
					tt.diags[i].Summary != tt.expected[i].Summary ||
					tt.diags[i].Detail != tt.expected[i].Detail {
					t.Errorf("expected %v, got %v", tt.expected[i], tt.diags[i])
				}
			}
		})
	}
}

func TestFormatCtyValueToIndentedString(t *testing.T) {
	tests := []struct {
		name           string
		value          cty.Value
		startingIndent int
		indentSize     int
		expected       string
	}{
		{
			name:           "null value",
			value:          cty.NullVal(cty.String),
			startingIndent: 0,
			indentSize:     2,
			expected:       "null",
		},
		{
			name:           "unknown value",
			value:          cty.UnknownVal(cty.String),
			startingIndent: 0,
			indentSize:     2,
			expected:       "null",
		},
		{
			name:           "string value",
			value:          cty.StringVal("test"),
			startingIndent: 0,
			indentSize:     2,
			expected:       "\"test\"",
		},
		{
			name:           "integer value",
			value:          cty.NumberIntVal(42),
			startingIndent: 0,
			indentSize:     2,
			expected:       "42",
		},
		{
			name:           "float value",
			value:          cty.NumberFloatVal(3.14),
			startingIndent: 0,
			indentSize:     2,
			expected:       "3.14",
		},
		{
			name:           "boolean value",
			value:          cty.BoolVal(true),
			startingIndent: 0,
			indentSize:     2,
			expected:       "true",
		},
		{
			name:           "list of strings value",
			value:          cty.ListVal([]cty.Value{cty.StringVal("test1"), cty.StringVal("test2")}),
			startingIndent: 0,
			indentSize:     2,
			expected:       "[\n  \"test1\",\n  \"test2\"\n]",
		},
		{
			name:           "tuple value",
			value:          cty.TupleVal([]cty.Value{cty.StringVal("test1"), cty.NumberIntVal(42)}),
			startingIndent: 0,
			indentSize:     2,
			expected:       "[\n  \"test1\",\n  42\n]",
		},
		{
			name:           "map of strings value",
			value:          cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			startingIndent: 0,
			indentSize:     2,
			expected:       "{\n  \"key1\": \"value1\",\n  \"key2\": \"value2\"\n}",
		},
		{
			name:           "object value",
			value:          cty.ObjectVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.NumberIntVal(42)}),
			startingIndent: 0,
			indentSize:     2,
			expected:       "{\n  \"key1\": \"value1\",\n  \"key2\": 42\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FormatCtyValueToIndentedString(tt.value, tt.startingIndent, tt.indentSize)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFormatCtyValueToString(t *testing.T) {
	tests := []struct {
		name     string
		value    cty.Value
		expected string
	}{
		{
			name:     "null value",
			value:    cty.NullVal(cty.String),
			expected: "null",
		},
		{
			name:     "unknown value",
			value:    cty.UnknownVal(cty.String),
			expected: "null",
		},
		{
			name:     "string value",
			value:    cty.StringVal("test"),
			expected: "\"test\"",
		},
		{
			name:     "integer value",
			value:    cty.NumberIntVal(42),
			expected: "42",
		},
		{
			name:     "float value",
			value:    cty.NumberFloatVal(3.14),
			expected: "3.14",
		},
		{
			name:     "boolean value",
			value:    cty.BoolVal(true),
			expected: "true",
		},
		{
			name:     "list of strings value",
			value:    cty.ListVal([]cty.Value{cty.StringVal("test1"), cty.StringVal("test2")}),
			expected: "[\"test1\", \"test2\"]",
		},
		{
			name:     "tuple value",
			value:    cty.TupleVal([]cty.Value{cty.StringVal("test1"), cty.NumberIntVal(42)}),
			expected: "[\"test1\", 42]",
		},
		{
			name:     "map of strings value",
			value:    cty.MapVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.StringVal("value2")}),
			expected: "{\"key1\": \"value1\", \"key2\": \"value2\"}",
		},
		{
			name:     "object value",
			value:    cty.ObjectVal(map[string]cty.Value{"key1": cty.StringVal("value1"), "key2": cty.NumberIntVal(42)}),
			expected: "{\"key1\": \"value1\", \"key2\": 42}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FormatCtyValueToString(tt.value)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

// mockExpr is a simple mock implementation of hcl.Expression for testing
type mockExpr struct {
	value cty.Value
	err   hcl.Diagnostics
}

func (m *mockExpr) Value(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return m.value, m.err
}

func (m *mockExpr) Variables() []hcl.Traversal {
	return nil
}

func (m *mockExpr) Range() hcl.Range {
	return hcl.Range{}
}

func (m *mockExpr) StartRange() hcl.Range {
	return hcl.Range{}
}
