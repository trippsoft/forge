package function

import (
	"os"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestEnvFunc(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		envValue    string
		setEnvVar   bool
		expectedVal cty.Value
		wantErr     bool
	}{
		{
			name:        "existing environment variable",
			envVar:      "TEST_ENV_VAR",
			envValue:    "test_value",
			setEnvVar:   true,
			expectedVal: cty.StringVal("test_value"),
			wantErr:     false,
		},
		{
			name:        "non-existing environment variable",
			envVar:      "NON_EXISTING_ENV_VAR",
			envValue:    "",
			setEnvVar:   false,
			expectedVal: cty.StringVal(""),
			wantErr:     false,
		},
		{
			name:        "empty environment variable",
			envVar:      "EMPTY_ENV_VAR",
			envValue:    "",
			setEnvVar:   true,
			expectedVal: cty.StringVal(""),
			wantErr:     false,
		},
		{
			name:        "environment variable with spaces",
			envVar:      "SPACE_ENV_VAR",
			envValue:    "  value with spaces  ",
			setEnvVar:   true,
			expectedVal: cty.StringVal("  value with spaces  "),
			wantErr:     false,
		},
		{
			name:        "environment variable with special characters",
			envVar:      "SPECIAL_ENV_VAR",
			envValue:    "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./ ",
			setEnvVar:   true,
			expectedVal: cty.StringVal("!@#$%^&*()_+-={}[]|\\:;\"'<>?,./ "),
			wantErr:     false,
		},
		{
			name:        "environment variable with unicode",
			envVar:      "UNICODE_ENV_VAR",
			envValue:    "héllo 世界 🌍",
			setEnvVar:   true,
			expectedVal: cty.StringVal("héllo 世界 🌍"),
			wantErr:     false,
		},
		{
			name:        "environment variable with newlines",
			envVar:      "MULTILINE_ENV_VAR",
			envValue:    "line1\nline2\nline3",
			setEnvVar:   true,
			expectedVal: cty.StringVal("line1\nline2\nline3"),
			wantErr:     false,
		},
		{
			name:        "environment variable with only whitespace",
			envVar:      "WHITESPACE_ENV_VAR",
			envValue:    "   ",
			setEnvVar:   true,
			expectedVal: cty.StringVal("   "),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variable if needed
			if tt.setEnvVar {
				originalValue, wasSet := os.LookupEnv(tt.envVar)
				os.Setenv(tt.envVar, tt.envValue)
				defer func() {
					if wasSet {
						os.Setenv(tt.envVar, originalValue)
					} else {
						os.Unsetenv(tt.envVar)
					}
				}()
			} else {
				// Ensure the env var is not set
				originalValue, wasSet := os.LookupEnv(tt.envVar)
				os.Unsetenv(tt.envVar)
				defer func() {
					if wasSet {
						os.Setenv(tt.envVar, originalValue)
					}
				}()
			}

			result, err := EnvFunc.Call([]cty.Value{cty.StringVal(tt.envVar)})

			if (err != nil) != tt.wantErr {
				t.Errorf("EnvFunc.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !result.RawEquals(tt.expectedVal) {
					t.Errorf("EnvFunc.Call() = %v, want %v", result, tt.expectedVal)
				}
			}
		})
	}
}

func TestEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVar      string
		envValue    string
		setEnvVar   bool
		expectedVal cty.Value
		wantErr     bool
	}{
		{
			name:        "existing environment variable",
			envVar:      "TEST_ENV_HELPER",
			envValue:    "helper_value",
			setEnvVar:   true,
			expectedVal: cty.StringVal("helper_value"),
			wantErr:     false,
		},
		{
			name:        "non-existing environment variable",
			envVar:      "NON_EXISTING_HELPER",
			envValue:    "",
			setEnvVar:   false,
			expectedVal: cty.StringVal(""),
			wantErr:     false,
		},
		{
			name:        "complex environment variable value",
			envVar:      "COMPLEX_ENV_HELPER",
			envValue:    "complex/path/with:colons;and=equals&ampersands",
			setEnvVar:   true,
			expectedVal: cty.StringVal("complex/path/with:colons;and=equals&ampersands"),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variable if needed
			if tt.setEnvVar {
				originalValue, wasSet := os.LookupEnv(tt.envVar)
				os.Setenv(tt.envVar, tt.envValue)
				defer func() {
					if wasSet {
						os.Setenv(tt.envVar, originalValue)
					} else {
						os.Unsetenv(tt.envVar)
					}
				}()
			} else {
				// Ensure the env var is not set
				originalValue, wasSet := os.LookupEnv(tt.envVar)
				os.Unsetenv(tt.envVar)
				defer func() {
					if wasSet {
						os.Setenv(tt.envVar, originalValue)
					}
				}()
			}

			result, err := Env(cty.StringVal(tt.envVar))

			if (err != nil) != tt.wantErr {
				t.Errorf("Env() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !result.RawEquals(tt.expectedVal) {
					t.Errorf("Env() = %v, want %v", result, tt.expectedVal)
				}
			}
		})
	}
}

func TestEnvFuncReturnType(t *testing.T) {
	// Test that the function returns the correct type when successful
	os.Setenv("TEST_TYPE_VAR", "test")
	defer os.Unsetenv("TEST_TYPE_VAR")

	result, err := EnvFunc.Call([]cty.Value{cty.StringVal("TEST_TYPE_VAR")})
	if err != nil {
		t.Fatalf("EnvFunc.Call() failed: %v", err)
	}

	if result.Type() != cty.String {
		t.Errorf("EnvFunc.Call() returned wrong type: got %v, want %v", result.Type(), cty.String)
	}
}

func TestEnvFuncWithExistingSystemEnvVars(t *testing.T) {
	// Test with actual system environment variables that should exist
	systemVars := []string{"PATH", "HOME", "USER"}

	for _, varName := range systemVars {
		t.Run("system_var_"+varName, func(t *testing.T) {
			if os.Getenv(varName) == "" {
				t.Skipf("System environment variable %s is not set, skipping", varName)
			}

			result, err := EnvFunc.Call([]cty.Value{cty.StringVal(varName)})
			if err != nil {
				t.Errorf("EnvFunc.Call() failed for system var %s: %v", varName, err)
			}

			if result.IsNull() {
				t.Errorf("EnvFunc.Call() returned nil for existing system var %s", varName)
			}

			if result.Type() != cty.String {
				t.Errorf("EnvFunc.Call() returned wrong type for %s: got %v, want %v", varName, result.Type(), cty.String)
			}

			// Verify the value matches what os.Getenv returns
			expectedValue := os.Getenv(varName)
			if result.AsString() != expectedValue {
				t.Errorf("EnvFunc.Call() returned %q for %s, but os.Getenv returned %q", result.AsString(), varName, expectedValue)
			}
		})
	}
}

func TestEnvFuncEmptyVsUnset(t *testing.T) {
	// Test the difference between an unset variable and an empty variable
	testVar := "TEST_EMPTY_VS_UNSET"

	// Test unset variable
	t.Run("unset_variable", func(t *testing.T) {
		os.Unsetenv(testVar)
		result, err := EnvFunc.Call([]cty.Value{cty.StringVal(testVar)})
		if err != nil {
			t.Errorf("EnvFunc.Call() failed: %v", err)
		}
		if result.IsNull() {
			t.Errorf("EnvFunc.Call() should return empty string for unset variable, got nil value")
		}
		if !result.IsNull() && result.AsString() != "" {
			t.Errorf("EnvFunc.Call() should return empty string for unset variable, got %v", result.AsString())
		}
	})

	// Test empty variable
	t.Run("empty_variable", func(t *testing.T) {
		os.Setenv(testVar, "")
		defer os.Unsetenv(testVar)

		result, err := EnvFunc.Call([]cty.Value{cty.StringVal(testVar)})
		if err != nil {
			t.Errorf("EnvFunc.Call() failed: %v", err)
		}
		if result.IsNull() {
			t.Errorf("EnvFunc.Call() should return empty string for empty variable, got nil value")
		}
		if !result.IsNull() && result.AsString() != "" {
			t.Errorf("EnvFunc.Call() should return empty string for empty variable, got %v", result.AsString())
		}
	})
}
