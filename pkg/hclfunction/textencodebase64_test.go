// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"golang.org/x/text/encoding/ianaindex"
)

func getTextEncodeBase64TestCases() []struct {
	name     string
	input    string
	encoding string
} {
	return []struct {
		name     string
		input    string
		encoding string
	}{
		{
			name:     "UTF-8 encoding",
			input:    "Hello, World!",
			encoding: "UTF-8",
		},
		{
			name:     "ASCII encoding",
			input:    "ASCII text",
			encoding: "US-ASCII",
		},
		{
			name:     "ISO-8859-1 encoding",
			input:    "café",
			encoding: "ISO-8859-1",
		},
		{
			name:     "Windows-1252 encoding",
			input:    "Windows text",
			encoding: "Windows-1252",
		},
		{
			name:     "empty string",
			input:    "",
			encoding: "UTF-8",
		},
		{
			name:     "unicode characters",
			input:    "Hello 世界! 🌍",
			encoding: "UTF-8",
		},
		{
			name:     "case insensitive encoding name",
			input:    "test",
			encoding: "utf-8",
		},
		{
			name:     "special characters",
			input:    "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			encoding: "UTF-8",
		},
		{
			name:     "newlines and tabs",
			input:    "line1\nline2\tindented",
			encoding: "UTF-8",
		},
	}
}

func getExpectedBase64Value(t *testing.T, input string, encoding string) cty.Value {
	e, err := ianaindex.IANA.Encoding(encoding)
	if err != nil {
		t.Fatalf("failed to get encoding %q: %v", encoding, err)
	}

	encoder := e.NewEncoder()
	encoded, err := encoder.Bytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to encode input %q as %q: %v", input, encoding, err)
	}

	base64Encoded := base64.StdEncoding.EncodeToString(encoded)

	return cty.StringVal(base64Encoded)
}

func TestTextEncodeBase64(t *testing.T) {
	tests := getTextEncodeBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := getExpectedBase64Value(t, tt.input, tt.encoding)

			input := cty.StringVal(tt.input)
			encoding := cty.StringVal(tt.encoding)

			actual, err := TextEncodeBase64(input, encoding)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}

func TestTextEncodeBase64_InvalidEncoding(t *testing.T) {
	// Test with an invalid encoding
	encoding := cty.StringVal("INVALID_ENCODING")
	input := cty.StringVal("test")

	_, err := TextEncodeBase64(input, encoding)
	if err == nil {
		t.Fatalf("expected error for invalid encoding, got none")
	}

	expectedError := fmt.Sprintf("invalid encoding %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestTextEncodeBase64Func(t *testing.T) {
	tests := getTextEncodeBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := getExpectedBase64Value(t, tt.input, tt.encoding)

			input := cty.StringVal(tt.input)
			encoding := cty.StringVal(tt.encoding)

			actual, err := TextEncodeBase64Func.Call([]cty.Value{input, encoding})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}

func TestTextEncodeBase64Func_InvalidEncoding(t *testing.T) {
	// Test with an invalid encoding
	encoding := cty.StringVal("INVALID_ENCODING")
	input := cty.StringVal("test")

	_, err := TextEncodeBase64Func.Call([]cty.Value{input, encoding})
	if err == nil {
		t.Fatalf("expected error for invalid encoding, got none")
	}

	expectedError := fmt.Sprintf("invalid encoding %q", encoding.AsString())
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}
