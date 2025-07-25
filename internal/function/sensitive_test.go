package function

import (
	"testing"

	"github.com/trippsoft/forge/internal/log"
	"github.com/zclconf/go-cty/cty"
)

func TestSensitiveFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "basic string value",
			input:    cty.StringVal("secret_password"),
			expected: cty.StringVal("secret_password"),
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    cty.StringVal(""),
			expected: cty.StringVal(""),
			wantErr:  false,
		},
		{
			name:     "string with special characters",
			input:    cty.StringVal("p@ssw0rd!#$"),
			expected: cty.StringVal("p@ssw0rd!#$"),
			wantErr:  false,
		},
		{
			name:     "unicode string",
			input:    cty.StringVal("пароль密码"),
			expected: cty.StringVal("пароль密码"),
			wantErr:  false,
		},
		{
			name:     "multiline string",
			input:    cty.StringVal("line1\nline2\nline3"),
			expected: cty.StringVal("line1\nline2\nline3"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original secrets and clear for test
			originalSecrets := log.LogSecretFilter.Secrets()
			log.LogSecretFilter.Clear()

			defer func() {
				// Restore original secrets after test
				log.LogSecretFilter.Clear()
				for _, secret := range originalSecrets {
					log.LogSecretFilter.AddSecret(secret)
				}
			}()

			result, err := SensitiveFunc.Call([]cty.Value{tt.input})

			if (err != nil) != tt.wantErr {
				t.Errorf("SensitiveFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !result.RawEquals(tt.expected) {
					t.Errorf("SensitiveFunc.Call() = %v, want %v", result, tt.expected)
				}

				// Verify the secret was added to the filter
				inputStr := tt.input.AsString()
				if inputStr != "" {
					filtered := log.LogSecretFilter.Filter(inputStr)
					if filtered != "<redacted>" {
						t.Errorf("Secret was not properly added to filter. Expected '<redacted>', got '%s'", filtered)
					}
				}
			}
		})
	}
}

func TestSensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    cty.Value
		expected cty.Value
		wantErr  bool
	}{
		{
			name:     "basic string value",
			input:    cty.StringVal("api_key_12345"),
			expected: cty.StringVal("api_key_12345"),
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    cty.StringVal(""),
			expected: cty.StringVal(""),
			wantErr:  false,
		},
		{
			name:     "complex secret",
			input:    cty.StringVal("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			expected: cty.StringVal("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original secrets and clear for test
			originalSecrets := log.LogSecretFilter.Secrets()
			log.LogSecretFilter.Clear()

			defer func() {
				// Restore original secrets after test
				log.LogSecretFilter.Clear()
				for _, secret := range originalSecrets {
					log.LogSecretFilter.AddSecret(secret)
				}
			}()

			result, err := Sensitive(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Sensitive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !result.RawEquals(tt.expected) {
					t.Errorf("Sensitive() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestSensitiveFuncSecretFiltering(t *testing.T) {
	// Test that secrets are properly added to the secret filter
	secrets := []string{
		"password123",
		"api_key_abcdef",
		"token_xyz789",
	}

	// Store original secrets and clear for test
	originalSecrets := log.LogSecretFilter.Secrets()
	log.LogSecretFilter.Clear()

	defer func() {
		// Restore original secrets after test
		log.LogSecretFilter.Clear()
		for _, secret := range originalSecrets {
			log.LogSecretFilter.AddSecret(secret)
		}
	}()

	// Add secrets using the Sensitive function
	for _, secret := range secrets {
		_, err := SensitiveFunc.Call([]cty.Value{cty.StringVal(secret)})
		if err != nil {
			t.Fatalf("Failed to call SensitiveFunc with secret '%s': %v", secret, err)
		}
	}

	// Test that all secrets are properly filtered
	testMessage := "The password123 and api_key_abcdef and token_xyz789 should be redacted"
	filtered := log.LogSecretFilter.Filter(testMessage)
	expected := "The <redacted> and <redacted> and <redacted> should be redacted"

	if filtered != expected {
		t.Errorf("Secret filtering failed.\nExpected: %s\nGot: %s", expected, filtered)
	}
}

func TestSensitiveFuncEmptyString(t *testing.T) {
	// Test that empty strings don't cause issues
	originalSecrets := log.LogSecretFilter.Secrets()
	log.LogSecretFilter.Clear()

	defer func() {
		// Restore original secrets after test
		log.LogSecretFilter.Clear()
		for _, secret := range originalSecrets {
			log.LogSecretFilter.AddSecret(secret)
		}
	}()

	result, err := SensitiveFunc.Call([]cty.Value{cty.StringVal("")})
	if err != nil {
		t.Fatalf("SensitiveFunc.Call() with empty string failed: %v", err)
	}

	if !result.RawEquals(cty.StringVal("")) {
		t.Errorf("SensitiveFunc.Call() with empty string = %v, want %v", result, cty.StringVal(""))
	}

	// Test that filtering with empty string doesn't break anything
	testMessage := "This message should not be affected"
	filtered := log.LogSecretFilter.Filter(testMessage)
	if filtered != testMessage {
		t.Errorf("Empty string secret affected filtering: got %s, want %s", filtered, testMessage)
	}
}

func TestSensitiveFuncReturnType(t *testing.T) {
	// Test that the function returns the correct type
	input := cty.StringVal("test_secret")
	result, err := SensitiveFunc.Call([]cty.Value{input})

	if err != nil {
		t.Fatalf("SensitiveFunc.Call() failed: %v", err)
	}

	if result.Type() != cty.String {
		t.Errorf("SensitiveFunc.Call() returned wrong type: got %v, want %v", result.Type(), cty.String)
	}

	// Test that the result is not null (as specified by RefineResult)
	if result.IsNull() {
		t.Error("SensitiveFunc.Call() returned null value, but RefineResult specifies NotNull()")
	}
}

func TestSensitiveFuncSecretsAreTracked(t *testing.T) {
	// Test that multiple secrets are properly tracked
	originalSecrets := log.LogSecretFilter.Secrets()
	log.LogSecretFilter.Clear()

	defer func() {
		// Restore original secrets after test
		log.LogSecretFilter.Clear()
		for _, secret := range originalSecrets {
			log.LogSecretFilter.AddSecret(secret)
		}
	}()

	testSecrets := []string{
		"secret1",
		"another_secret",
		"third_secret",
	}

	// Add secrets one by one
	for _, secret := range testSecrets {
		_, err := SensitiveFunc.Call([]cty.Value{cty.StringVal(secret)})
		if err != nil {
			t.Fatalf("Failed to add secret '%s': %v", secret, err)
		}
	}

	// Verify all secrets are tracked
	trackedSecrets := log.LogSecretFilter.Secrets()
	if len(trackedSecrets) != len(testSecrets) {
		t.Errorf("Expected %d secrets to be tracked, got %d", len(testSecrets), len(trackedSecrets))
	}

	// Verify each secret is tracked
	secretsMap := make(map[string]bool)
	for _, secret := range trackedSecrets {
		secretsMap[secret] = true
	}

	for _, expected := range testSecrets {
		if !secretsMap[expected] {
			t.Errorf("Secret '%s' was not found in tracked secrets", expected)
		}
	}
}

func TestSensitiveFuncEmptyStringIsIgnored(t *testing.T) {
	// Test that empty strings don't get added to the secret filter
	originalSecrets := log.LogSecretFilter.Secrets()
	log.LogSecretFilter.Clear()

	defer func() {
		// Restore original secrets after test
		log.LogSecretFilter.Clear()
		for _, secret := range originalSecrets {
			log.LogSecretFilter.AddSecret(secret)
		}
	}()

	// Add an empty string as a "secret"
	_, err := SensitiveFunc.Call([]cty.Value{cty.StringVal("")})
	if err != nil {
		t.Fatalf("Failed to call SensitiveFunc with empty string: %v", err)
	}

	// Verify no secrets are tracked (empty string should be ignored)
	trackedSecrets := log.LogSecretFilter.Secrets()
	if len(trackedSecrets) != 0 {
		t.Errorf("Expected 0 secrets to be tracked when adding empty string, got %d", len(trackedSecrets))
	}
}
