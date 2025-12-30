// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestRemoveEmptyLinesAndComments(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected string
	}{
		{
			name:     "test1",
			script:   "# This is a comment\n\nprint('Hello, World!')\n",
			expected: "print('Hello, World!')",
		},
		{
			name:     "test2",
			script:   "   \n\n\n",
			expected: "",
		},
		{
			name:     "test3",
			script:   "# Another comment\nx = 5\n# Yet another comment\ny = 10\n",
			expected: "x = 5\ny = 10",
		},
		{
			name:     "test4",
			script:   "# Comment\n\nprint('Hello, World!')    \n",
			expected: "print('Hello, World!')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := RemoveEmptyLinesAndComments(tt.script)
			if actual != tt.expected {
				t.Fatalf("expected %q from RemoveEmptyLinesAndComments(), got %q", tt.expected, actual)
			}
		})
	}
}

func TestEncodePythonAsBase64_EmptyInput(t *testing.T) {
	_, err := EncodePythonAsBase64("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestEncodePythonAsBase64(t *testing.T) {
	tests := []struct {
		name     string
		python   string
		expected string
	}{
		{
			name:     "valid python command",
			python:   "print('Hello')",
			expected: "cHJpbnQoJ0hlbGxvJyk=",
		},
		{
			name:     "multiline python",
			python:   "x = 5\ny = 10",
			expected: "eCA9IDUKeSA9IDEw",
		},
		{
			name:     "python with special characters",
			python:   "print('Hello, World!')",
			expected: "cHJpbnQoJ0hlbGxvLCBXb3JsZCEnKQ==",
		},
		{
			name:     "single character",
			python:   "x",
			expected: "eA==",
		},
		{
			name:     "python with whitespace",
			python:   "   print('test')   ",
			expected: "ICAgcHJpbnQoJ3Rlc3QnKSAgIA==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := EncodePythonAsBase64(tt.python)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestFormatInputForPython(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]cty.Value
		whatIf   bool
		expected string
	}{
		{
			name:     "empty input map with whatIf false",
			input:    map[string]cty.Value{},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {}\n",
		},
		{
			name:     "empty input map with whatIf true",
			input:    map[string]cty.Value{},
			whatIf:   true,
			expected: "WHAT_IF = True\nINPUT = {}\n",
		},
		{
			name: "single string value",
			input: map[string]cty.Value{
				"name": cty.StringVal("test"),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"name\": \"test\"}\n",
		},
		{
			name: "single boolean true value",
			input: map[string]cty.Value{
				"enabled": cty.BoolVal(true),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"enabled\": True}\n",
		},
		{
			name: "single boolean false value",
			input: map[string]cty.Value{
				"enabled": cty.BoolVal(false),
			},
			whatIf:   true,
			expected: "WHAT_IF = True\nINPUT = {\"enabled\": False}\n",
		},
		{
			name: "single number value",
			input: map[string]cty.Value{
				"count": cty.NumberIntVal(42),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"count\": 42}\n",
		},
		{
			name: "null value",
			input: map[string]cty.Value{
				"optional": cty.NullVal(cty.String),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"optional\": None}\n",
		},
		{
			name: "list with strings",
			input: map[string]cty.Value{
				"items": cty.ListVal([]cty.Value{
					cty.StringVal("apple"),
					cty.StringVal("banana"),
					cty.StringVal("cherry"),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"items\": [\"apple\", \"banana\", \"cherry\"]}\n",
		},
		{
			name: "list with numbers",
			input: map[string]cty.Value{
				"counts": cty.ListVal([]cty.Value{
					cty.NumberIntVal(10),
					cty.NumberIntVal(20),
					cty.NumberIntVal(30),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"counts\": [10, 20, 30]}\n",
		},
		{
			name: "list with booleans",
			input: map[string]cty.Value{
				"flags": cty.ListVal([]cty.Value{
					cty.BoolVal(true),
					cty.BoolVal(false),
					cty.BoolVal(true),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"flags\": [True, False, True]}\n",
		},
		{
			name: "empty list",
			input: map[string]cty.Value{
				"empty": cty.ListValEmpty(cty.String),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"empty\": []}\n",
		},
		{
			name: "tuple",
			input: map[string]cty.Value{
				"coords": cty.TupleVal([]cty.Value{
					cty.NumberIntVal(10),
					cty.NumberIntVal(20),
					cty.NumberIntVal(30),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"coords\": [10, 20, 30]}\n",
		},
		{
			name: "set with strings",
			input: map[string]cty.Value{
				"tags": cty.SetVal([]cty.Value{
					cty.StringVal("prod"),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"tags\": {\"prod\"}}\n",
		},
		{
			name: "map with string values",
			input: map[string]cty.Value{
				"config": cty.MapVal(map[string]cty.Value{
					"host": cty.StringVal("localhost"),
					"port": cty.StringVal("8080"),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"config\": {\"host\": \"localhost\", \"port\": \"8080\"}}\n",
		},
		{
			name: "object type",
			input: map[string]cty.Value{
				"user": cty.ObjectVal(map[string]cty.Value{
					"name": cty.StringVal("John"),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"user\": {\"name\": \"John\"}}\n",
		},
		{
			name: "nested list",
			input: map[string]cty.Value{
				"matrix": cty.ListVal([]cty.Value{
					cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2)}),
					cty.ListVal([]cty.Value{cty.NumberIntVal(3), cty.NumberIntVal(4)}),
				}),
			},
			whatIf:   false,
			expected: "WHAT_IF = False\nINPUT = {\"matrix\": [[1, 2], [3, 4]]}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := FormatInputForPython(tt.input, tt.whatIf)

			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
