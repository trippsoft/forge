// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"os"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func getEnvTestCases() []struct {
	name     string
	envVar   string
	envValue string
} {
	return []struct {
		name     string
		envVar   string
		envValue string
	}{
		{
			name:     "existing environment variable",
			envVar:   "TEST_ENV_VAR",
			envValue: "test_value",
		},
		{
			name:     "empty environment variable",
			envVar:   "EMPTY_ENV_VAR",
			envValue: "",
		},
		{
			name:     "environment variable with spaces",
			envVar:   "SPACE_ENV_VAR",
			envValue: "  value with spaces  ",
		},
		{
			name:     "environment variable with special characters",
			envVar:   "SPECIAL_ENV_VAR",
			envValue: "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./ ",
		},
		{
			name:     "environment variable with unicode",
			envVar:   "UNICODE_ENV_VAR",
			envValue: "héllo 世界 🌍",
		},
		{
			name:     "environment variable with newlines",
			envVar:   "MULTILINE_ENV_VAR",
			envValue: "line1\nline2\nline3",
		},
		{
			name:     "environment variable with only whitespace",
			envVar:   "WHITESPACE_ENV_VAR",
			envValue: "   ",
		},
	}
}

func TestEnv(t *testing.T) {
	tests := getEnvTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue, wasSet := os.LookupEnv(tt.envVar)
			os.Setenv(tt.envVar, tt.envValue)
			defer func() {
				if wasSet {
					os.Setenv(tt.envVar, originalValue)
					return
				}
				os.Unsetenv(tt.envVar)
			}()

			expected := cty.StringVal(tt.envValue)

			actual, err := Env(cty.StringVal(tt.envVar))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}

func TestEnvFunc(t *testing.T) {
	tests := getEnvTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue, wasSet := os.LookupEnv(tt.envVar)
			os.Setenv(tt.envVar, tt.envValue)
			defer func() {
				if wasSet {
					os.Setenv(tt.envVar, originalValue)
					return
				}
				os.Unsetenv(tt.envVar)
			}()

			expected := cty.StringVal(tt.envValue)

			actual, err := EnvFunc.Call([]cty.Value{cty.StringVal(tt.envVar)})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			assertCtyValueEqual(t, actual, expected)
		})
	}
}
