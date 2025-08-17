// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package assert

import (
	"context"
	"testing"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestModuleInputSpec(t *testing.T) {

	m := &Module{}

	spec := m.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got %v", err)
	}
}

func TestModuleValidate(t *testing.T) {

	m := &Module{}

	input := map[string]cty.Value{
		"condition": cty.BoolVal(true),
	}

	mockTransport := transport.NewMockTransport()
	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("linux", mockTransport, escalateConfig, map[string]cty.Value{})

	config := &module.RunConfig{
		Transport:  mockTransport,
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	err := m.Validate(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestModuleRun(t *testing.T) {

	tests := []struct {
		name           string
		input          map[string]cty.Value
		expectedFailed bool
	}{
		{
			name: "condition is true",
			input: map[string]cty.Value{
				"condition": cty.BoolVal(true),
			},
			expectedFailed: false,
		},
		{
			name: "condition is false",
			input: map[string]cty.Value{
				"condition": cty.BoolVal(false),
			},
			expectedFailed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockTransport := transport.NewMockTransport()
			escalateConfig := inventory.NewEscalateConfig("")
			host := inventory.NewHost("linux", mockTransport, escalateConfig, map[string]cty.Value{})

			ctx := context.Background()

			config := &module.RunConfig{
				Transport:  mockTransport,
				HostInfo:   host.Info(),
				Escalation: nil,
				Input:      tt.input,
			}

			m := &Module{}

			result := m.Run(ctx, config)

			if result == nil {
				t.Fatalf("Expected non-nil result, got nil")
			}

			if result.Failed != tt.expectedFailed {
				t.Errorf("Expected failed: %t, got: %t", tt.expectedFailed, result.Failed)
			}
		})
	}
}
