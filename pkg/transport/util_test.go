// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"encoding/base64"
	"fmt"
	"testing"

	"golang.org/x/text/encoding/unicode"
)

func TestEncodePowerShellAsUTF16LEBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple command",
			input:    "Write-Host 'Hello'",
			expected: "VwByAGkAdABlAC0ASABvAHMAdAAgACcASABlAGwAbABvACcA",
		},

		{
			name:     "single character",
			input:    "A",
			expected: "QQA=",
		},
		{
			name:     "unicode characters",
			input:    "Hello 世界",
			expected: "SABlAGwAbABvACAAT^QL5Q==",
		},
		{
			name:     "special characters",
			input:    "Get-Process | Where-Object {$_.Name -eq 'notepad'}",
			expected: "RwBlAHQALQBQAHIAbwBjAGUAcwBzACAAfAAgAFcAaABlAHIAZQAtAE8AYgBqAGUAYwB0ACAAewAkAF8ALgBOAGEAbQBlACAALQBlAHEAIAAnAG4AbwB0AGUAcABhAGQAJwB9AA==",
		},
		{
			name:     "powershell with quotes",
			input:    `Write-Host "Hello World"`,
			expected: "VwByAGkAdABlAC0ASABvAHMAdAAgACIASABlAGwAbABvACAAVwBvAHIAbABkACIA",
		},
		{
			name:     "multiline script",
			input:    "Get-Service\nGet-Process",
			expected: "RwBlAHQALQBTAGUAcgB2AGkAYwBlAAoARwBlAHQALQBQAHIAbwBjAGUAcwBzAA==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePowerShellAsUTF16LEBase64(tt.input)
			if err != nil {
				t.Fatalf("encodePowerShellAsUTF16LEBase64() error = %v", err)
			}

			// Verify the result is not empty for valid inputs
			if result == "" {
				t.Error("Expected non-empty result for non-empty input")
				return
			}

			// Verify it's valid base64
			decoded, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Errorf("Result is not valid base64: %v", err)
				return
			}

			// Verify it decodes back to the original string when interpreted as UTF-16LE
			decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
			decodedString, err := decoder.String(string(decoded))
			if err != nil {
				t.Errorf("Failed to decode UTF-16LE: %v", err)
				return
			}

			if decodedString != tt.input {
				t.Errorf("Round trip failed: expected %q, got %q", tt.input, decodedString)
			}

			t.Logf("Input: %q -> Base64: %q", tt.input, result)
		})
	}
}

func TestEncodePowerShellAsUTF16LEBase64EmptyString(t *testing.T) {
	_, err := encodePowerShellAsUTF16LEBase64("")
	if err == nil {
		t.Error("Expected error for empty input, but got none")
	}

	expectedError := "input PowerShell command cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEncodePowerShellAsUTF16LEBase64WhitespaceOnly(t *testing.T) {
	// Test various whitespace-only inputs that should be valid
	whitespaceInputs := []string{
		" ",      // single space
		"  ",     // multiple spaces
		"\t",     // tab
		"\n",     // newline
		"\r\n",   // Windows line ending
		" \t\n ", // mixed whitespace
	}

	for _, input := range whitespaceInputs {
		t.Run(fmt.Sprintf("whitespace_%q", input), func(t *testing.T) {
			result, err := encodePowerShellAsUTF16LEBase64(input)
			if err != nil {
				t.Errorf("Expected no error for whitespace input %q, got: %v", input, err)
			}
			if result == "" {
				t.Errorf("Expected non-empty result for whitespace input %q", input)
			}

			// Verify round trip
			decoded, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Errorf("Result is not valid base64: %v", err)
				return
			}

			decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
			decodedString, err := decoder.String(string(decoded))
			if err != nil {
				t.Errorf("Failed to decode UTF-16LE: %v", err)
				return
			}

			if decodedString != input {
				t.Errorf("Round trip failed: expected %q, got %q", input, decodedString)
			}
		})
	}
}

func TestEncodePowerShellAsUTF16LEBase64LongScript(t *testing.T) {
	// Test with a longer PowerShell script
	longScript := `
$computers = Get-Content "C:\computers.txt"
foreach ($computer in $computers) {
    if (Test-Connection -ComputerName $computer -Count 1 -Quiet) {
        Write-Host "$computer is online" -ForegroundColor Green
        $services = Get-Service -ComputerName $computer | Where-Object {$_.Status -eq "Running"}
        Write-Host "Running services on $computer : $($services.Count)"
    } else {
        Write-Host "$computer is offline" -ForegroundColor Red
    }
}
`

	result, err := encodePowerShellAsUTF16LEBase64(longScript)
	if err != nil {
		t.Fatalf("encodePowerShellAsUTF16LEBase64() error = %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result for long script")
	}

	// Verify it's valid base64
	decoded, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		t.Errorf("Result is not valid base64: %v", err)
		return
	}

	// Verify round trip
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	decodedString, err := decoder.String(string(decoded))
	if err != nil {
		t.Errorf("Failed to decode UTF-16LE: %v", err)
		return
	}

	if decodedString != longScript {
		t.Errorf("Round trip failed for long script")
	}
}

func TestEncodePowerShellAsUTF16LEBase64SpecialCharacters(t *testing.T) {
	// Test with various special characters that might appear in PowerShell
	testCases := []string{
		`Write-Host "Hello\nWorld"`, // Escape sequences
		`$env:PATH`,                 // Environment variables
		`Get-ChildItem -Path "C:\Program Files (x86)"`,   // Paths with spaces and parentheses
		`[System.Environment]::GetFolderPath("Desktop")`, // .NET calls
		`$(Get-Date).ToString("yyyy-MM-dd")`,             // Complex expressions
		`Write-Host $([char]0x1F600)`,                    // Unicode emoji (if supported)
		`@"
This is a here-string
with multiple lines
and special characters: !@#$%^&*()
"@`, // Here-strings
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("special_chars_%d", i), func(t *testing.T) {
			result, err := encodePowerShellAsUTF16LEBase64(testCase)
			if err != nil {
				t.Fatalf("encodePowerShellAsUTF16LEBase64() error = %v", err)
			}

			if result == "" {
				t.Error("Expected non-empty result")
				return
			}

			// Verify round trip
			decoded, err := base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Errorf("Result is not valid base64: %v", err)
				return
			}

			decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
			decodedString, err := decoder.String(string(decoded))
			if err != nil {
				t.Errorf("Failed to decode UTF-16LE: %v", err)
				return
			}

			if decodedString != testCase {
				t.Errorf("Round trip failed: expected %q, got %q", testCase, decodedString)
			}
		})
	}
}

func TestEncodePowerShellAsUTF16LEBase64Encoding(t *testing.T) {
	input := "Hi"
	expected := []byte{
		0x48, 0x00, // 'H' in UTF-16LE
		0x69, 0x00, // 'i' in UTF-16LE
	}

	expectedBase64 := base64.StdEncoding.EncodeToString(expected)

	result, err := encodePowerShellAsUTF16LEBase64(input)
	if err != nil {
		t.Fatalf("encodePowerShellAsUTF16LEBase64() error = %v", err)
	}

	if result != expectedBase64 {
		t.Errorf("Expected %q, got %q", expectedBase64, result)
	}
}

func TestEncodePowerShellAsUTF16LEBase64Consistency(t *testing.T) {
	input := "Get-Process | Select-Object Name"

	results := make([]string, 5)
	for i := range results {
		result, err := encodePowerShellAsUTF16LEBase64(input)
		if err != nil {
			t.Fatalf("encodePowerShellAsUTF16LEBase64() error on call %d: %v", i+1, err)
		}

		results[i] = result
	}

	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Inconsistent results: call 1 returned %q, call %d returned %q", results[0], i+1, results[i])
		}
	}
}
