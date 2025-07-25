package function

import (
	"encoding/base64"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestTextDecodeBase64(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		encoding    string
		expectedErr bool
		expected    string
	}{
		{
			name:        "UTF-8 decoding",
			input:       base64.StdEncoding.EncodeToString([]byte("Hello, World!")),
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "Hello, World!",
		},
		{
			name:        "ASCII decoding",
			input:       base64.StdEncoding.EncodeToString([]byte("ASCII text")),
			encoding:    "US-ASCII",
			expectedErr: false,
			expected:    "ASCII text",
		},
		{
			name:        "empty string",
			input:       "",
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "",
		},
		{
			name:        "unicode characters",
			input:       base64.StdEncoding.EncodeToString([]byte("Hello 世界! 🌍")),
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "Hello 世界! 🌍",
		},
		{
			name:        "special characters",
			input:       base64.StdEncoding.EncodeToString([]byte("!@#$%^&*()_+-=[]{}|;':\",./<>?")),
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:        "newlines and tabs",
			input:       base64.StdEncoding.EncodeToString([]byte("line1\nline2\tindented")),
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "line1\nline2\tindented",
		},
		{
			name:        "invalid base64",
			input:       "not-valid-base64!@#",
			encoding:    "UTF-8",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "invalid base64 character",
			input:       "SGVsbG8h@",
			encoding:    "UTF-8",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "invalid encoding",
			input:       base64.StdEncoding.EncodeToString([]byte("test")),
			encoding:    "INVALID-ENCODING",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "case insensitive encoding name",
			input:       base64.StdEncoding.EncodeToString([]byte("test")),
			encoding:    "utf-8",
			expectedErr: false,
			expected:    "test",
		},
		{
			name:        "ISO-8859-1 encoded text",
			input:       "Y2Fm6Q==", // "café" encoded in ISO-8859-1 then base64 (0x63,0x61,0x66,0xE9)
			encoding:    "ISO-8859-1",
			expectedErr: false,
			expected:    "café",
		},
		{
			name:        "Windows-1252 encoding",
			input:       base64.StdEncoding.EncodeToString([]byte("Windows text")),
			encoding:    "Windows-1252",
			expectedErr: false,
			expected:    "Windows text",
		},
		{
			name:        "malformed base64 padding",
			input:       "SGVsbG8",
			encoding:    "UTF-8",
			expectedErr: true,
			expected:    "",
		},
		{
			name:        "base64 with proper padding",
			input:       "SGVsbG8=",
			encoding:    "UTF-8",
			expectedErr: false,
			expected:    "Hello",
		},
		{
			name:        "text that cannot be decoded in specified encoding",
			input:       base64.StdEncoding.EncodeToString([]byte{0xFF, 0xFE, 0x00, 0x48}), // Invalid UTF-8 sequence
			encoding:    "UTF-8",
			expectedErr: true,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TextDecodeBase64(cty.StringVal(tt.input), cty.StringVal(tt.encoding))

			if (err != nil) != tt.expectedErr {
				t.Errorf("TextDecodeBase64() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr {
				if result.Type() != cty.String {
					t.Errorf("Expected string type, got %v", result.Type())
				}

				if result.AsString() != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result.AsString())
				}
			}
		})
	}
}

func TestTextDecodeBase64Func(t *testing.T) {
	// Test calling the function directly
	input := base64.StdEncoding.EncodeToString([]byte("test"))
	result, err := TextDecodeBase64Func.Call([]cty.Value{
		cty.StringVal(input),
		cty.StringVal("UTF-8"),
	})

	if err != nil {
		t.Fatalf("TextDecodeBase64Func.Call() failed: %v", err)
	}

	if result.Type() != cty.String {
		t.Errorf("Expected string type, got %v", result.Type())
	}

	if result.AsString() != "test" {
		t.Errorf("Expected 'test', got %q", result.AsString())
	}
}

func TestTextDecodeBase64WithInvalidArgs(t *testing.T) {
	// Test with wrong number of arguments
	_, err := TextDecodeBase64Func.Call([]cty.Value{
		cty.StringVal("dGVzdA=="),
	})
	if err == nil {
		t.Error("Expected error with insufficient arguments")
	}

	// Test with too many arguments
	_, err = TextDecodeBase64Func.Call([]cty.Value{
		cty.StringVal("dGVzdA=="),
		cty.StringVal("UTF-8"),
		cty.StringVal("extra"),
	})
	if err == nil {
		t.Error("Expected error with too many arguments")
	}

	// Test with wrong argument types
	_, err = TextDecodeBase64Func.Call([]cty.Value{
		cty.NumberIntVal(123),
		cty.StringVal("UTF-8"),
	})
	if err == nil {
		t.Error("Expected error with wrong input type")
	}

	_, err = TextDecodeBase64Func.Call([]cty.Value{
		cty.StringVal("dGVzdA=="),
		cty.NumberIntVal(123),
	})
	if err == nil {
		t.Error("Expected error with wrong encoding type")
	}
}

func TestTextDecodeBase64CorruptInputError(t *testing.T) {
	// Test specific error handling for corrupt base64 input
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid character",
			input: "SGVsb@8=",
		},
		{
			name:  "invalid padding",
			input: "SGVsbG8===",
		},
		{
			name:  "truly truncated input",
			input: "S",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TextDecodeBase64(cty.StringVal(tt.input), cty.StringVal("UTF-8"))
			if err == nil {
				t.Error("Expected error for corrupt base64 input")
			}
		})
	}
}

// Test round-trip encoding and decoding
func TestTextEncodeDecodeBase64RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		encoding string
	}{
		{
			name:     "UTF-8 round trip",
			input:    "Hello, World! 世界 🌍",
			encoding: "UTF-8",
		},
		{
			name:     "ASCII round trip",
			input:    "Simple ASCII text",
			encoding: "US-ASCII",
		},
		{
			name:     "empty string round trip",
			input:    "",
			encoding: "UTF-8",
		},
		{
			name:     "special characters round trip",
			input:    "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			encoding: "UTF-8",
		},
		{
			name:     "ISO-8859-1 round trip",
			input:    "café résumé naïve",
			encoding: "ISO-8859-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First encode
			encoded, err := TextEncodeBase64(cty.StringVal(tt.input), cty.StringVal(tt.encoding))
			if err != nil {
				t.Fatalf("TextEncodeBase64() failed: %v", err)
			}

			// Then decode
			decoded, err := TextDecodeBase64(encoded, cty.StringVal(tt.encoding))
			if err != nil {
				t.Fatalf("TextDecodeBase64() failed: %v", err)
			}

			// Verify round trip
			if decoded.AsString() != tt.input {
				t.Errorf("Round trip failed: expected %q, got %q", tt.input, decoded.AsString())
			}
		})
	}
}
