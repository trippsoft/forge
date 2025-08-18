// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"fmt"
	"testing"
)

func TestDiagError(t *testing.T) {
	summary := "Test error"
	detail := "This is a test error detail."
	diag := &Diag{
		Severity: DiagError,
		Summary:  summary,
		Detail:   detail,
	}

	expectedError := fmt.Sprintf("%s; %s", summary, detail)
	if diag.Error() != expectedError {
		t.Errorf("Expected %q, got %q", expectedError, diag.Error())
	}
}

func TestDiagsError(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected string
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: "no diagnostics to report",
		},
		{
			name: "Single diagnostic",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Single error",
					Detail:   "This is a single error detail",
				},
			},
			expected: "Single error; This is a single error detail",
		},
		{
			name: "Multiple diagnostics",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "First error",
					Detail:   "This is the first error detail",
				},
				{
					Severity: DiagWarning,
					Summary:  "Second warning",
					Detail:   "This is the second warning detail",
				},
			},
			expected: "First error; This is the first error detail, and 1 more diagnostics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diags.Error(); got != tt.expected {
				t.Errorf("Diags.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected bool
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: false,
		},
		{
			name: "Only warnings",
			diags: Diags{
				{
					Severity: DiagWarning,
					Summary:  "Warning",
					Detail:   "This is a warning",
				},
			},
			expected: false,
		},
		{
			name: "Contains error",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error",
					Detail:   "This is an error",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diags.HasErrors(); got != tt.expected {
				t.Errorf("Diags.HasErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsErrors(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected int
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: 0,
		},
		{
			name: "Only warnings",
			diags: Diags{
				{
					Severity: DiagWarning,
					Summary:  "Warning",
					Detail:   "This is a warning",
				},
			},
			expected: 0,
		},
		{
			name: "Contains errors",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error 1",
					Detail:   "This is the first error",
				},
				{
					Severity: DiagError,
					Summary:  "Error 2",
					Detail:   "This is the second error",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.diags.Errors()); got != tt.expected {
				t.Errorf("Diags.Errors() length = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsHasWarnings(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected bool
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: false,
		},
		{
			name: "Only errors",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error",
					Detail:   "This is an error",
				},
			},
			expected: false,
		},
		{
			name: "Contains warning",
			diags: Diags{
				{
					Severity: DiagWarning,
					Summary:  "Warning",
					Detail:   "This is a warning",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diags.HasWarnings(); got != tt.expected {
				t.Errorf("Diags.HasWarnings() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsWarnings(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected int
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: 0,
		},
		{
			name: "Only errors",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error",
					Detail:   "This is an error",
				},
			},
			expected: 0,
		},
		{
			name: "Contains warnings",
			diags: Diags{
				{
					Severity: DiagWarning,
					Summary:  "Warning 1",
					Detail:   "This is the first warning",
				},
				{
					Severity: DiagWarning,
					Summary:  "Warning 2",
					Detail:   "This is the second warning",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.diags.Warnings()); got != tt.expected {
				t.Errorf("Diags.Warnings() length = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsHasDebug(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected bool
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: false,
		},
		{
			name: "Only errors and warnings",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error",
					Detail:   "This is an error",
				},
				{
					Severity: DiagWarning,
					Summary:  "Warning",
					Detail:   "This is a warning",
				},
			},
			expected: false,
		},
		{
			name: "Contains debug",
			diags: Diags{
				{
					Severity: DiagDebug,
					Summary:  "Debug info",
					Detail:   "This is debug information",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.diags.HasDebug(); got != tt.expected {
				t.Errorf("Diags.HasDebug() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsDebugs(t *testing.T) {
	tests := []struct {
		name     string
		diags    Diags
		expected int
	}{
		{
			name:     "No diagnostics",
			diags:    Diags{},
			expected: 0,
		},
		{
			name: "Only errors and warnings",
			diags: Diags{
				{
					Severity: DiagError,
					Summary:  "Error",
					Detail:   "This is an error",
				},
				{
					Severity: DiagWarning,
					Summary:  "Warning",
					Detail:   "This is a warning",
				},
			},
			expected: 0,
		},
		{
			name: "Contains debug diagnostics",
			diags: Diags{
				{
					Severity: DiagDebug,
					Summary:  "Debug 1",
					Detail:   "This is the first debug detail",
				},
				{
					Severity: DiagDebug,
					Summary:  "Debug 2",
					Detail:   "This is the second debug detail",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.diags.Debugs()); got != tt.expected {
				t.Errorf("Diags.Debugs() length = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagsAppend(t *testing.T) {
	diag1 := &Diag{
		Severity: DiagError,
		Summary:  "Error 1",
		Detail:   "This is the first error",
	}
	diag2 := &Diag{
		Severity: DiagWarning,
		Summary:  "Warning 1",
		Detail:   "This is the first warning",
	}

	diags := Diags{diag1}
	diags = diags.Append(diag2)

	if len(diags) != 2 {
		t.Errorf("Expected 2 diagnostics, got %d", len(diags))
	}
	if diags[0] != diag1 || diags[1] != diag2 {
		t.Error("Diagnostics were not appended correctly")
	}
}

func TestDiagsAppend_Null(t *testing.T) {
	diag := &Diag{
		Severity: DiagError,
		Summary:  "Error",
		Detail:   "This is an error",
	}
	diags := Diags{diag}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected not to panic when appending nil Diag, but did panic")
		}
	}()

	diags = diags.Append(nil)
	if len(diags) != 1 {
		t.Errorf("Expected 1 diagnostic after appending nil, got %d", len(diags))
	}
	if diags[0] != diag {
		t.Error("Original diagnostic was modified after appending nil")
	}
}

func TestDiagsAppendAll(t *testing.T) {
	diag1 := &Diag{
		Severity: DiagError,
		Summary:  "Error 1",
		Detail:   "This is the first error",
	}
	diag2 := &Diag{
		Severity: DiagWarning,
		Summary:  "Warning 1",
		Detail:   "This is the first warning",
	}

	diags := Diags{diag1}
	moreDiags := Diags{diag2}

	diags = diags.AppendAll(moreDiags)

	if len(diags) != 2 {
		t.Errorf("Expected 2 diagnostics, got %d", len(diags))
	}
	if diags[0] != diag1 || diags[1] != diag2 {
		t.Error("Diagnostics were not appended correctly")
	}
}

func TestDiagsAppendAll_Null(t *testing.T) {
	diag := &Diag{
		Severity: DiagError,
		Summary:  "Error",
		Detail:   "This is an error",
	}
	diags := Diags{diag}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected not to panic when appending nil Diags, but did panic")
		}
	}()

	diags = diags.AppendAll(nil)
	if len(diags) != 1 {
		t.Errorf("Expected 1 diagnostic after appending nil, got %d", len(diags))
	}
	if diags[0] != diag {
		t.Error("Original diagnostic was modified after appending nil")
	}
}
