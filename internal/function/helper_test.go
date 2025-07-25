package function

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

// Helper function to compare cty.Values including unknown values
func assertCtyValueEqual(t *testing.T, result, expected cty.Value, testName string) {
	t.Helper()
	if expected.IsKnown() && result.IsKnown() {
		if !result.Equals(expected).True() {
			t.Errorf("%s = %v, want %v", testName, result, expected)
		}
	} else if expected.IsKnown() != result.IsKnown() {
		t.Errorf("%s known status = %v, want %v", testName, result.IsKnown(), expected.IsKnown())
	} else if !expected.IsKnown() && !result.IsKnown() {
		// Both are unknown, check if they have the same type
		if !result.Type().Equals(expected.Type()) {
			t.Errorf("%s unknown value type = %v, want %v", testName, result.Type(), expected.Type())
		}
	}
}
