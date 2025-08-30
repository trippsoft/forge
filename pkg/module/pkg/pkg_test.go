// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pkg

import "testing"

func TestModuleInputSpec(t *testing.T) {
	module := &PkgModule{}

	spec := module.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec from InputSpec(), got nil")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got: %q", err.Error())
	}
}
