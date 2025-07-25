package function

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFileBase64(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "filebase64_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		setupFile   func() (string, string) // returns (filepath, content)
		expectedErr bool
	}{
		{
			name: "text file",
			setupFile: func() (string, string) {
				content := "Hello, World!"
				filepath := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(filepath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return filepath, content
			},
			expectedErr: false,
		},
		{
			name: "empty file",
			setupFile: func() (string, string) {
				content := ""
				filepath := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(filepath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return filepath, content
			},
			expectedErr: false,
		},
		{
			name: "binary file",
			setupFile: func() (string, string) {
				content := "\x00\x01\x02\x03\xFF"
				filepath := filepath.Join(tempDir, "binary.bin")
				err := os.WriteFile(filepath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return filepath, content
			},
			expectedErr: false,
		},
		{
			name: "unicode content",
			setupFile: func() (string, string) {
				content := "Hello 世界! 🌍"
				filepath := filepath.Join(tempDir, "unicode.txt")
				err := os.WriteFile(filepath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return filepath, content
			},
			expectedErr: false,
		},
		{
			name: "non-existent file",
			setupFile: func() (string, string) {
				return filepath.Join(tempDir, "nonexistent.txt"), ""
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filepath, originalContent := tt.setupFile()

			result, err := FileBase64(cty.StringVal(filepath))
			if (err != nil) != tt.expectedErr {
				t.Errorf("FileBase64() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr {
				// Verify it's a string
				if result.Type() != cty.String {
					t.Errorf("FileBase64() returned type %v, want %v", result.Type(), cty.String)
				}

				// Decode the base64 and verify it matches original content
				encodedContent := result.AsString()
				decodedBytes, decodeErr := base64.StdEncoding.DecodeString(encodedContent)
				if decodeErr != nil {
					t.Errorf("FileBase64() returned invalid base64: %v", decodeErr)
				}

				if string(decodedBytes) != originalContent {
					t.Errorf("FileBase64() decoded content = %v, want %v", string(decodedBytes), originalContent)
				}
			}
		})
	}
}

func TestFileBase64Func(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "filebase64_func_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := "test content"
	filepath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		input       []cty.Value
		expectedErr bool
	}{
		{
			name:        "function call with valid file",
			input:       []cty.Value{cty.StringVal(filepath)},
			expectedErr: false,
		},
		{
			name:        "function call with non-existent file",
			input:       []cty.Value{cty.StringVal("/nonexistent/file.txt")},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FileBase64Func.Call(tt.input)
			if (err != nil) != tt.expectedErr {
				t.Errorf("FileBase64Func.Call() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if !tt.expectedErr {
				// Verify it's a string
				if result.Type() != cty.String {
					t.Errorf("FileBase64Func.Call() returned type %v, want %v", result.Type(), cty.String)
				}

				// Decode and verify content matches
				encodedContent := result.AsString()
				decodedBytes, decodeErr := base64.StdEncoding.DecodeString(encodedContent)
				if decodeErr != nil {
					t.Errorf("FileBase64Func.Call() returned invalid base64: %v", decodeErr)
				}

				if string(decodedBytes) != content {
					t.Errorf("FileBase64Func.Call() decoded content = %v, want %v", string(decodedBytes), content)
				}
			}
		})
	}
}
