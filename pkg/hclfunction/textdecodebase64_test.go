// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getTextDecodeBase64TestCases() []struct {
	name     string
	input    string
	encoding string
	expected cty.Value
} {
	return []struct {
		name     string
		input    string
		encoding string
		expected cty.Value
	}{
		{
			name:     "UTF-8 decoding",
			input:    base64.StdEncoding.EncodeToString([]byte("Hello, World!")),
			encoding: "UTF-8",
			expected: cty.StringVal("Hello, World!"),
		},
		{
			name:     "ASCII decoding",
			input:    base64.StdEncoding.EncodeToString([]byte("ASCII text")),
			encoding: "US-ASCII",
			expected: cty.StringVal("ASCII text"),
		},
		{
			name:     "empty string",
			input:    "",
			encoding: "UTF-8",
			expected: cty.StringVal(""),
		},
		{
			name:     "unicode characters",
			input:    base64.StdEncoding.EncodeToString([]byte("Hello 世界! 🌍")),
			encoding: "UTF-8",
			expected: cty.StringVal("Hello 世界! 🌍"),
		},
		{
			name:     "special characters",
			input:    base64.StdEncoding.EncodeToString([]byte("!@#$%^&*()_+-=[]{}|;':\",./<>?")),
			encoding: "UTF-8",
			expected: cty.StringVal("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
		},
		{
			name:     "newlines and tabs",
			input:    base64.StdEncoding.EncodeToString([]byte("line1\nline2\tindented")),
			encoding: "UTF-8",
			expected: cty.StringVal("line1\nline2\tindented"),
		},
		{
			name:     "case insensitive encoding name",
			input:    base64.StdEncoding.EncodeToString([]byte("test")),
			encoding: "utf-8",
			expected: cty.StringVal("test"),
		},
		{
			name:     "ISO-8859-1 encoded text",
			input:    "Y2Fm6Q==", // "café" encoded in ISO-8859-1 then base64 (0x63,0x61,0x66,0xE9)
			encoding: "ISO-8859-1",
			expected: cty.StringVal("café"),
		},
		{
			name:     "Windows-1252 encoding",
			input:    base64.StdEncoding.EncodeToString([]byte("Windows text")),
			encoding: "Windows-1252",
			expected: cty.StringVal("Windows text"),
		},
		{
			name:     "base64 with proper padding",
			input:    "SGVsbG8=",
			encoding: "UTF-8",
			expected: cty.StringVal("Hello"),
		},
	}
}

func TestTextDecodeBase64(t *testing.T) {
	tests := getTextDecodeBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := TextDecodeBase64(cty.StringVal(tt.input), cty.StringVal(tt.encoding))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestTextDecodeBase64_InvalidBase64(t *testing.T) {
	input := cty.StringVal("SGVsbG8h@=") // Invalid base64
	encoding := cty.StringVal("UTF-8")

	_, err := TextDecodeBase64(input, encoding)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "the input has invalid base64 character at offset: 8"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTextDecodeBase64_InvalidEncoding(t *testing.T) {
	input := cty.StringVal(base64.StdEncoding.EncodeToString([]byte("test")))
	encoding := cty.StringVal("INVALID_ENCODING")

	_, err := TextDecodeBase64(input, encoding)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := fmt.Sprintf("invalid encoding %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestTextDecodeBase64_NoPadding(t *testing.T) {
	input := cty.StringVal("SGVsbG8") // Base64 without padding
	encoding := cty.StringVal("UTF-8")

	_, err := TextDecodeBase64(input, encoding)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "the input has invalid base64 character at offset: 4"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTextDecodeBase64_TextCannotBeEncoded(t *testing.T) {
	input := cty.StringVal(base64.StdEncoding.EncodeToString([]byte{0xFF, 0xFE, 0x00, 0x48})) // Invalid UTF-8 sequence
	encoding := cty.StringVal("UTF-8")

	_, err := TextDecodeBase64(input, encoding)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := fmt.Sprintf("failed to decode input as %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTextDecodeBase64Func(t *testing.T) {
	tests := getTextDecodeBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := TextDecodeBase64Func.Call([]cty.Value{cty.StringVal(tt.input), cty.StringVal(tt.encoding)})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, tt.expected)
		})
	}
}

func TestTextDecodeBase64Func_InvalidBase64(t *testing.T) {
	_, err := TextDecodeBase64Func.Call([]cty.Value{cty.StringVal("SGVsbG8h@="), cty.StringVal("UTF-8")})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "the input has invalid base64 character at offset: 8"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTextDecodeBase64Func_InvalidEncoding(t *testing.T) {
	input := cty.StringVal(base64.StdEncoding.EncodeToString([]byte("test")))
	encoding := cty.StringVal("INVALID_ENCODING")

	_, err := TextDecodeBase64Func.Call([]cty.Value{input, encoding})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := fmt.Sprintf("invalid encoding %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestTextDecodeBase64Func_NoPadding(t *testing.T) {
	input := cty.StringVal("SGVsbG8") // Base64 without padding
	encoding := cty.StringVal("UTF-8")

	_, err := TextDecodeBase64Func.Call([]cty.Value{input, encoding})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := "the input has invalid base64 character at offset: 4"
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTextDecodeBase64Func_TextCannotBeEncoded(t *testing.T) {
	input := cty.StringVal(base64.StdEncoding.EncodeToString([]byte{0xFF, 0xFE, 0x00, 0x48})) // Invalid UTF-8 sequence
	encoding := cty.StringVal("UTF-8")

	_, err := TextDecodeBase64Func.Call([]cty.Value{input, encoding})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	expectedError := fmt.Sprintf("failed to decode input as %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}
