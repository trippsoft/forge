// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package hclfunction

import (
	"testing"

	"github.com/zclconf/go-cty/cty"
)

// Helper function to compare cty.Values including unknown values
func assertCtyValueEqual(t *testing.T, actual, expected cty.Value) {
	if !expected.Type().Equals(actual.Type()) {
		t.Fatalf("expected type %q, got %q", expected.Type().FriendlyName(), actual.Type().FriendlyName())
	}

	if !expected.IsKnown() && actual.IsKnown() {
		t.Fatalf("expected unknown value, got %v", actual)
	}

	if expected.IsKnown() && !actual.IsKnown() {
		t.Fatalf("expected known value, got unknown")
	}

	if expected.IsKnown() && !expected.Equals(actual).True() {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
