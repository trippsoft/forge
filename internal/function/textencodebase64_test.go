package function

import (
	"encoding/base64"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestTextEncodeBase64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		encoding    string
		expectedErr bool
		validate    func(t *testing.T, result cty.Value)
	}{
		{
			name:        "UTF-8 encoding",
			input:       "Hello, World!",
			encoding:    "UTF-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Decode the base64 result to verify
				decoded, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
				if string(decoded) != "Hello, World!" {
					t.Errorf("Expected 'Hello, World!', got %q", string(decoded))
				}
			},
		},
		{
			name:        "ASCII encoding",
			input:       "ASCII text",
			encoding:    "US-ASCII",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Verify it's valid base64
				decoded, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
				if string(decoded) != "ASCII text" {
					t.Errorf("Expected 'ASCII text', got %q", string(decoded))
				}
			},
		},
		{
			name:        "ISO-8859-1 encoding",
			input:       "café",
			encoding:    "ISO-8859-1",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Verify it's valid base64
				_, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
			},
		},
		{
			name:        "Windows-1252 encoding",
			input:       "Windows text",
			encoding:    "Windows-1252",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Verify it's valid base64
				_, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
			},
		},
		{
			name:        "empty string",
			input:       "",
			encoding:    "UTF-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Empty input should result in empty base64
				if result.AsString() != "" {
					t.Errorf("Expected empty string for empty input, got %q", result.AsString())
				}
			},
		},
		{
			name:        "unicode characters",
			input:       "Hello 世界! 🌍",
			encoding:    "UTF-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Decode and verify unicode is preserved
				decoded, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
				if string(decoded) != "Hello 世界! 🌍" {
					t.Errorf("Expected 'Hello 世界! 🌍', got %q", string(decoded))
				}
			},
		},
		{
			name:        "invalid encoding",
			input:       "test",
			encoding:    "INVALID-ENCODING",
			expectedErr: true,
			validate:    nil,
		},
		{
			name:        "case insensitive encoding name",
			input:       "test",
			encoding:    "utf-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
			},
		},
		{
			name:        "special characters",
			input:       "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			encoding:    "UTF-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Decode and verify special characters are preserved
				decoded, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
				if string(decoded) != "!@#$%^&*()_+-=[]{}|;':\",./<>?" {
					t.Errorf("Special characters not preserved correctly")
				}
			},
		},
		{
			name:        "newlines and tabs",
			input:       "line1\nline2\tindented",
			encoding:    "UTF-8",
			expectedErr: false,
			validate: func(t *testing.T, result cty.Value) {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}
				// Decode and verify whitespace is preserved
				decoded, err := base64.StdEncoding.DecodeString(result.AsString())
				if err != nil {
					t.Errorf("Failed to decode base64 result: %v", err)
				}
				if string(decoded) != "line1\nline2\tindented" {
					t.Errorf("Whitespace not preserved correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TextEncodeBase64(cty.StringVal(tt.input), cty.StringVal(tt.encoding))

			if (err != nil) != tt.expectedErr {
				t.Errorf("TextEncodeBase64() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestTextEncodeBase64Func(t *testing.T) {
	// Test calling the function directly
	result, err := TextEncodeBase64Func.Call([]cty.Value{
		cty.StringVal("test"),
		cty.StringVal("UTF-8"),
	})

	if err != nil {
		t.Fatalf("TextEncodeBase64Func.Call() failed: %v", err)
	}

	if result.Type() != cty.String {
		t.Errorf("Expected string type, got %v", result.Type())
	}

	// Verify the result is valid base64
	decoded, err := base64.StdEncoding.DecodeString(result.AsString())
	if err != nil {
		t.Errorf("Result is not valid base64: %v", err)
	}

	if string(decoded) != "test" {
		t.Errorf("Expected 'test', got %q", string(decoded))
	}
}

func TestTextEncodeBase64WithInvalidArgs(t *testing.T) {
	// Test with wrong number of arguments
	_, err := TextEncodeBase64Func.Call([]cty.Value{
		cty.StringVal("test"),
	})
	if err == nil {
		t.Error("Expected error with insufficient arguments")
	}

	// Test with too many arguments
	_, err = TextEncodeBase64Func.Call([]cty.Value{
		cty.StringVal("test"),
		cty.StringVal("UTF-8"),
		cty.StringVal("extra"),
	})
	if err == nil {
		t.Error("Expected error with too many arguments")
	}

	// Test with wrong argument types
	_, err = TextEncodeBase64Func.Call([]cty.Value{
		cty.NumberIntVal(123),
		cty.StringVal("UTF-8"),
	})
	if err == nil {
		t.Error("Expected error with wrong input type")
	}

	_, err = TextEncodeBase64Func.Call([]cty.Value{
		cty.StringVal("test"),
		cty.NumberIntVal(123),
	})
	if err == nil {
		t.Error("Expected error with wrong encoding type")
	}
}
