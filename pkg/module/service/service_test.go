// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package service

import "testing"

func TestModuleInputSpec(t *testing.T) {
	module := &ServiceModule{}

	spec := module.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec from InputSpec(), got nil")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got: %q", err.Error())
	}
}
