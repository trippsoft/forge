// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"bytes"
	"strings"
	"testing"
)

func TestTrace(t *testing.T) {
	tests := []struct {
		name          string
		logMessage    string
		verbosity     LogLevel
		shouldContain bool
	}{
		{
			name:          "Trace message at trace verbosity",
			logMessage:    "This is a trace message",
			verbosity:     LevelTrace,
			shouldContain: true,
		},
		{
			name:          "Trace message at debug verbosity",
			logMessage:    "This is a trace message",
			verbosity:     LevelDebug,
			shouldContain: false,
		},
		{
			name:          "Trace message at info verbosity",
			logMessage:    "This is a trace message",
			verbosity:     LevelInfo,
			shouldContain: false,
		},
		{
			name:          "Trace message at warn verbosity",
			logMessage:    "This is a trace message",
			verbosity:     LevelWarn,
			shouldContain: false,
		},
		{
			name:          "Trace message at error verbosity",
			logMessage:    "This is a trace message",
			verbosity:     LevelError,
			shouldContain: false,
		},
	}

	expected := "TRACE:\tThis is a trace message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Trace(tt.logMessage)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, expected) {
				t.Errorf("Expected trace message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, expected) {
				t.Errorf("Trace message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestTracef(t *testing.T) {
	tests := []struct {
		name          string
		logTemplate   string
		args          []any
		verbosity     LogLevel
		shouldContain bool
		expected      string
	}{
		{
			name:          "Tracef with no args at trace verbosity",
			logTemplate:   "This is a trace message with no args",
			args:          []any{},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "TRACE:\tThis is a trace message with no args",
		},
		{
			name:          "Tracef with args at trace verbosity",
			logTemplate:   "This is a trace message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "TRACE:\tThis is a trace message with args: arg1, 42",
		},
		{
			name:          "Tracef with no args at debug verbosity",
			logTemplate:   "This is a trace message with no args",
			args:          []any{},
			verbosity:     LevelDebug,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with no args",
		},
		{
			name:          "Tracef with args at debug verbosity",
			logTemplate:   "This is a trace message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelDebug,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with args: arg1, 42",
		},
		{
			name:          "Tracef with no args at info verbosity",
			logTemplate:   "This is a trace message with no args",
			args:          []any{},
			verbosity:     LevelInfo,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with no args",
		},
		{
			name:          "Tracef with args at info verbosity",
			logTemplate:   "This is a trace message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelInfo,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with args: arg1, 42",
		},
		{
			name:          "Tracef with no args at warn verbosity",
			logTemplate:   "This is a trace message with no args",
			args:          []any{},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with no args",
		},
		{
			name:          "Tracef with args at warn verbosity",
			logTemplate:   "This is a trace message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with args: arg1, 42",
		},
		{
			name:          "Tracef with no args at error verbosity",
			logTemplate:   "This is a trace message with no args",
			args:          []any{},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with no args",
		},
		{
			name:          "Tracef with args at error verbosity",
			logTemplate:   "This is a trace message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "TRACE:\tThis is a trace message with args: arg1, 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Tracef(tt.logTemplate, tt.args...)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected trace message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, tt.expected) {
				t.Errorf("Trace message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name          string
		logMessage    string
		verbosity     LogLevel
		shouldContain bool
	}{
		{
			name:          "Debug message at trace verbosity",
			logMessage:    "This is a debug message",
			verbosity:     LevelTrace,
			shouldContain: true,
		},
		{
			name:          "Debug message at debug verbosity",
			logMessage:    "This is a debug message",
			verbosity:     LevelDebug,
			shouldContain: true,
		},
		{
			name:          "Debug message at info verbosity",
			logMessage:    "This is a debug message",
			verbosity:     LevelInfo,
			shouldContain: false,
		},
		{
			name:          "Debug message at warn verbosity",
			logMessage:    "This is a debug message",
			verbosity:     LevelWarn,
			shouldContain: false,
		},
		{
			name:          "Debug message at error verbosity",
			logMessage:    "This is a debug message",
			verbosity:     LevelError,
			shouldContain: false,
		},
	}

	expected := "DEBUG:\tThis is a debug message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Debug(tt.logMessage)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, expected) {
				t.Errorf("Expected debug message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, expected) {
				t.Errorf("Debug message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestDebugf(t *testing.T) {
	tests := []struct {
		name          string
		logTemplate   string
		args          []any
		verbosity     LogLevel
		shouldContain bool
		expected      string
	}{
		{
			name:          "Debugf with no args at trace verbosity",
			logTemplate:   "This is a debug message with no args",
			args:          []any{},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "DEBUG:\tThis is a debug message with no args",
		},
		{
			name:          "Debugf with args at trace verbosity",
			logTemplate:   "This is a debug message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "DEBUG:\tThis is a debug message with args: arg1, 42",
		},
		{
			name:          "Debugf with no args at debug verbosity",
			logTemplate:   "This is a debug message with no args",
			args:          []any{},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "DEBUG:\tThis is a debug message with no args",
		},
		{
			name:          "Debugf with args at debug verbosity",
			logTemplate:   "This is a debug message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "DEBUG:\tThis is a debug message with args: arg1, 42",
		},
		{
			name:          "Debugf with no args at info verbosity",
			logTemplate:   "This is a debug message with no args",
			args:          []any{},
			verbosity:     LevelInfo,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with no args",
		},
		{
			name:          "Debugf with args at info verbosity",
			logTemplate:   "This is a debug message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelInfo,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with args: arg1, 42",
		},
		{
			name:          "Debugf with no args at warn verbosity",
			logTemplate:   "This is a debug message with no args",
			args:          []any{},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with no args",
		},
		{
			name:          "Debugf with args at warn verbosity",
			logTemplate:   "This is a debug message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with args: arg1, 42",
		},
		{
			name:          "Debugf with no args at error verbosity",
			logTemplate:   "This is a debug message with no args",
			args:          []any{},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with no args",
		},
		{
			name:          "Debugf with args at error verbosity",
			logTemplate:   "This is a debug message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "DEBUG:\tThis is a debug message with args: arg1, 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Debugf(tt.logTemplate, tt.args...)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected debug message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, tt.expected) {
				t.Errorf("Debug message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name          string
		logMessage    string
		verbosity     LogLevel
		shouldContain bool
	}{
		{
			name:          "Info message at trace verbosity",
			logMessage:    "This is an info message",
			verbosity:     LevelTrace,
			shouldContain: true,
		},
		{
			name:          "Info message at debug verbosity",
			logMessage:    "This is an info message",
			verbosity:     LevelDebug,
			shouldContain: true,
		},
		{
			name:          "Info message at info verbosity",
			logMessage:    "This is an info message",
			verbosity:     LevelInfo,
			shouldContain: true,
		},
		{
			name:          "Info message at warn verbosity",
			logMessage:    "This is an info message",
			verbosity:     LevelWarn,
			shouldContain: false,
		},
		{
			name:          "Info message at error verbosity",
			logMessage:    "This is an info message",
			verbosity:     LevelError,
			shouldContain: false,
		},
	}

	expected := "INFO :\tThis is an info message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Info(tt.logMessage)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, expected) {
				t.Errorf("Expected info message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, expected) {
				t.Errorf("Info message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	tests := []struct {
		name          string
		logTemplate   string
		args          []any
		verbosity     LogLevel
		shouldContain bool
		expected      string
	}{
		{
			name:          "Infof with no args at trace verbosity",
			logTemplate:   "This is an info message with no args",
			args:          []any{},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with no args",
		},
		{
			name:          "Infof with args at trace verbosity",
			logTemplate:   "This is an info message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with args: arg1, 42",
		},
		{
			name:          "Infof with no args at debug verbosity",
			logTemplate:   "This is an info message with no args",
			args:          []any{},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with no args",
		},
		{
			name:          "Infof with args at debug verbosity",
			logTemplate:   "This is an info message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with args: arg1, 42",
		},
		{
			name:          "Infof with no args at info verbosity",
			logTemplate:   "This is an info message with no args",
			args:          []any{},
			verbosity:     LevelInfo,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with no args",
		},
		{
			name:          "Infof with args at info verbosity",
			logTemplate:   "This is an info message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelInfo,
			shouldContain: true,
			expected:      "INFO :\tThis is an info message with args: arg1, 42",
		},
		{
			name:          "Infof with no args at warn verbosity",
			logTemplate:   "This is an info message with no args",
			args:          []any{},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "INFO :\tThis is an info message with no args",
		},
		{
			name:          "Infof with args at warn verbosity",
			logTemplate:   "This is an info message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelWarn,
			shouldContain: false,
			expected:      "INFO :\tThis is an info message with args: arg1, 42",
		},
		{
			name:          "Infof with no args at error verbosity",
			logTemplate:   "This is an info message with no args",
			args:          []any{},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "INFO :\tThis is an info message with no args",
		},
		{
			name:          "Infof with args at error verbosity",
			logTemplate:   "This is an info message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "INFO :\tThis is an info message with args: arg1, 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Infof(tt.logTemplate, tt.args...)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected info message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, tt.expected) {
				t.Errorf("Info message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name          string
		logMessage    string
		verbosity     LogLevel
		shouldContain bool
	}{
		{
			name:          "Warn message at trace verbosity",
			logMessage:    "This is a warning message",
			verbosity:     LevelTrace,
			shouldContain: true,
		},
		{
			name:          "Warn message at debug verbosity",
			logMessage:    "This is a warning message",
			verbosity:     LevelDebug,
			shouldContain: true,
		},
		{
			name:          "Warn message at info verbosity",
			logMessage:    "This is a warning message",
			verbosity:     LevelInfo,
			shouldContain: true,
		},
		{
			name:          "Warn message at warn verbosity",
			logMessage:    "This is a warning message",
			verbosity:     LevelWarn,
			shouldContain: true,
		},
		{
			name:          "Warn message at error verbosity",
			logMessage:    "This is a warning message",
			verbosity:     LevelError,
			shouldContain: false,
		},
	}

	expected := "WARN :\tThis is a warning message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Warn(tt.logMessage)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, expected) {
				t.Errorf("Expected warning message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, expected) {
				t.Errorf("Warning message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	tests := []struct {
		name          string
		logTemplate   string
		args          []any
		verbosity     LogLevel
		shouldContain bool
		expected      string
	}{
		{
			name:          "Warnf with no args at trace verbosity",
			logTemplate:   "This is a warning message with no args",
			args:          []any{},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with no args",
		},
		{
			name:          "Warnf with args at trace verbosity",
			logTemplate:   "This is a warning message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelTrace,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with args: arg1, 42",
		},
		{
			name:          "Warnf with no args at debug verbosity",
			logTemplate:   "This is a warning message with no args",
			args:          []any{},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with no args",
		},
		{
			name:          "Warnf with args at debug verbosity",
			logTemplate:   "This is a warning message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelDebug,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with args: arg1, 42",
		},
		{
			name:          "Warnf with no args at info verbosity",
			logTemplate:   "This is a warning message with no args",
			args:          []any{},
			verbosity:     LevelInfo,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with no args",
		},
		{
			name:          "Warnf with args at info verbosity",
			logTemplate:   "This is a warning message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelInfo,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with args: arg1, 42",
		},
		{
			name:          "Warnf with no args at warn verbosity",
			logTemplate:   "This is a warning message with no args",
			args:          []any{},
			verbosity:     LevelWarn,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with no args",
		},
		{
			name:          "Warnf with args at warn verbosity",
			logTemplate:   "This is a warning message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelWarn,
			shouldContain: true,
			expected:      "WARN :\tThis is a warning message with args: arg1, 42",
		},
		{
			name:          "Warnf with no args at error verbosity",
			logTemplate:   "This is a warning message with no args",
			args:          []any{},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "WARN :\tThis is a warning message with no args",
		},
		{
			name:          "Warnf with args at error verbosity",
			logTemplate:   "This is a warning message with args: %s, %d",
			args:          []any{"arg1", 42},
			verbosity:     LevelError,
			shouldContain: false,
			expected:      "WARN :\tThis is a warning message with args: arg1, 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Warnf(tt.logTemplate, tt.args...)

			actual := buf.String()
			if tt.shouldContain && !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected warning message not found in output: %s", actual)
			} else if !tt.shouldContain && strings.Contains(actual, tt.expected) {
				t.Errorf("Warning message should not be logged at this verbosity level: %s", actual)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		logMessage string
		verbosity  LogLevel
	}{
		{
			name:       "Error message at trace verbosity",
			logMessage: "This is an error message",
			verbosity:  LevelTrace,
		},
		{
			name:       "Error message at debug verbosity",
			logMessage: "This is an error message",
			verbosity:  LevelDebug,
		},
		{
			name:       "Error message at info verbosity",
			logMessage: "This is an error message",
			verbosity:  LevelInfo,
		},
		{
			name:       "Error message at warn verbosity",
			logMessage: "This is an error message",
			verbosity:  LevelWarn,
		},
		{
			name:       "Error message at error verbosity",
			logMessage: "This is an error message",
			verbosity:  LevelError,
		},
	}

	expected := "ERROR:\tThis is an error message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Error(tt.logMessage)

			actual := buf.String()
			if !strings.Contains(actual, expected) {
				t.Errorf("Expected error message not found in output: %s", actual)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		name        string
		logTemplate string
		args        []any
		verbosity   LogLevel
		expected    string
	}{
		{
			name:        "Errorf with no args at trace verbosity",
			logTemplate: "This is an error message with no args",
			args:        []any{},
			verbosity:   LevelTrace,
			expected:    "ERROR:\tThis is an error message with no args",
		},
		{
			name:        "Errorf with args at trace verbosity",
			logTemplate: "This is an error message with args: %s, %d",
			args:        []any{"arg1", 42},
			verbosity:   LevelTrace,
			expected:    "ERROR:\tThis is an error message with args: arg1, 42",
		},
		{
			name:        "Errorf with no args at debug verbosity",
			logTemplate: "This is an error message with no args",
			args:        []any{},
			verbosity:   LevelDebug,
			expected:    "ERROR:\tThis is an error message with no args",
		},
		{
			name:        "Errorf with args at debug verbosity",
			logTemplate: "This is an error message with args: %s, %d",
			args:        []any{"arg1", 42},
			verbosity:   LevelDebug,
			expected:    "ERROR:\tThis is an error message with args: arg1, 42",
		},
		{
			name:        "Errorf with no args at info verbosity",
			logTemplate: "This is an error message with no args",
			args:        []any{},
			verbosity:   LevelInfo,
			expected:    "ERROR:\tThis is an error message with no args",
		},
		{
			name:        "Errorf with args at info verbosity",
			logTemplate: "This is an error message with args: %s, %d",
			args:        []any{"arg1", 42},
			verbosity:   LevelInfo,
			expected:    "ERROR:\tThis is an error message with args: arg1, 42",
		},
		{
			name:        "Errorf with no args at warn verbosity",
			logTemplate: "This is an error message with no args",
			args:        []any{},
			verbosity:   LevelWarn,
			expected:    "ERROR:\tThis is an error message with no args",
		},
		{
			name:        "Errorf with args at warn verbosity",
			logTemplate: "This is an error message with args: %s, %d",
			args:        []any{"arg1", 42},
			verbosity:   LevelWarn,
			expected:    "ERROR:\tThis is an error message with args: arg1, 42",
		},
		{
			name:        "Errorf with no args at error verbosity",
			logTemplate: "This is an error message with no args",
			args:        []any{},
			verbosity:   LevelError,
			expected:    "ERROR:\tThis is an error message with no args",
		},
		{
			name:        "Errorf with args at error verbosity",
			logTemplate: "This is an error message with args: %s, %d",
			args:        []any{"arg1", 42},
			verbosity:   LevelError,
			expected:    "ERROR:\tThis is an error message with args: arg1, 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Init(&buf)
			Verbosity = tt.verbosity

			Errorf(tt.logTemplate, tt.args...)

			actual := buf.String()
			if !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected error message not found in output: %s", actual)
			}
		})
	}
}
