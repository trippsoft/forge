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
			// Create a mock attribute and eval context
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
			// Create a mock attribute and eval context
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
			// Create a mock attribute and eval context
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
	// Test that expression evaluation errors are properly propagated
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
