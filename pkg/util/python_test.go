// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import "testing"

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
