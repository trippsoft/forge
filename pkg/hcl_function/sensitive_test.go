package hcl_function

import (
	"testing"

	"github.com/trippsoft/forge/pkg/log"
	"github.com/zclconf/go-cty/cty"
)

func getSensitiveTestCases() []struct {
	name  string
	input cty.Value
} {
	return []struct {
		name  string
		input cty.Value
	}{
		{
			name:  "basic string value",
			input: cty.StringVal("secret_password"),
		},
		{
			name:  "empty string",
			input: cty.StringVal(""),
		},
		{
			name:  "string with special characters",
			input: cty.StringVal("p@ssw0rd!#$"),
		},
		{
			name:  "unicode string",
			input: cty.StringVal("пароль密码"),
		},
		{
			name:  "multiline string",
			input: cty.StringVal("line1\nline2\nline3"),
		},
		{
			name:  "complex secret",
			input: cty.StringVal("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
		},
	}
}

func TestSensitiveFunc(t *testing.T) {

	tests := getSensitiveTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Store original secrets and clear for test
			originalSecrets := log.SecretFilter.Secrets()
			log.SecretFilter.Clear()

			defer func() {
				// Restore original secrets after test
				log.SecretFilter.Clear()
				for _, secret := range originalSecrets {
					log.SecretFilter.AddSecret(secret)
				}
			}()

			expected := tt.input
			actual, err := SensitiveFunc.Call([]cty.Value{tt.input})

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)

			// Verify the secret was added to the filter
			inputStr := tt.input.AsString()
			if inputStr != "" {
				filtered := log.SecretFilter.Filter(inputStr)
				if filtered != "<redacted>" {
					t.Errorf("Secret was not properly added to filter. Expected '<redacted>', got '%s'", filtered)
				}
			}
		})
	}
}

func TestSensitive(t *testing.T) {

	tests := getSensitiveTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Store original secrets and clear for test
			originalSecrets := log.SecretFilter.Secrets()
			log.SecretFilter.Clear()

			defer func() {
				// Restore original secrets after test
				log.SecretFilter.Clear()
				for _, secret := range originalSecrets {
					log.SecretFilter.AddSecret(secret)
				}
			}()

			expected := tt.input

			actual, err := Sensitive(tt.input)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)

			// Verify the secret was added to the filter
			inputStr := tt.input.AsString()
			if inputStr != "" {
				filtered := log.SecretFilter.Filter(inputStr)
				if filtered != "<redacted>" {
					t.Errorf("Secret was not properly added to filter. Expected '<redacted>', got '%s'", filtered)
				}
			}
		})
	}
}
