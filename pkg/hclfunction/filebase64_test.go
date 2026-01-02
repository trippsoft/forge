// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func setupTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "filebase64_test*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

func getFileBase64TestCases() []struct {
	name    string
	content []byte
} {
	return []struct {
		name    string
		content []byte
	}{
		{
			name:    "text file",
			content: []byte("Hello, World!"),
		},
		{
			name:    "empty file",
			content: []byte(""),
		},
		{
			name:    "binary file",
			content: []byte("\x00\x01\x02\x03\xFF"),
		},
		{
			name:    "unicode content",
			content: []byte("Hello 世界! 🌍"),
		},
	}
}

func createTempFile(t *testing.T, dir string, content []byte) string {
	file, err := os.CreateTemp(dir, "testfile_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	err = file.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return file.Name()
}

func encodeExpectedAsBase64(content []byte) cty.Value {
	encoded := base64.StdEncoding.EncodeToString(content)
	return cty.StringVal(encoded)
}

func TestFileBase64(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := getFileBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTempFile(t, tempDir, tt.content)

			expected := encodeExpectedAsBase64(tt.content)

			actual, err := FileBase64(cty.StringVal(path))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}

func TestFileBase64_NonExistentFile(t *testing.T) {
	// Test with a non-existent file
	_, err := FileBase64(cty.StringVal("non_existent_file.txt"))
	if err == nil {
		t.Fatal("expected an error for non-existent file, got nil")
	}

	expectedErr := "filebase64 failed: failed to read file "
	if !strings.HasPrefix(err.Error(), expectedErr) {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestFileBase64Func(t *testing.T) {
	tempDir := setupTempDir(t)
	defer os.RemoveAll(tempDir)

	tests := getFileBase64TestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := createTempFile(t, tempDir, tt.content)

			expected := encodeExpectedAsBase64(tt.content)

			workingDir, _ := os.Getwd()
			actual, err := MakeFileBase64Func(workingDir).Call([]cty.Value{cty.StringVal(path)})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}

func TestFileBase64Func_NonExistentFile(t *testing.T) {
	// Test with a non-existent file
	workingDir, _ := os.Getwd()
	_, err := MakeFileBase64Func(workingDir).Call([]cty.Value{cty.StringVal("non_existent_file.txt")})
	if err == nil {
		t.Fatal("expected an error for non-existent file, got nil")
	}

	expectedErr := "filebase64 failed: failed to read file "
	if !strings.HasPrefix(err.Error(), expectedErr) {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
