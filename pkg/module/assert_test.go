// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"testing"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestAssertModuleRun(t *testing.T) {
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

			hostBuilder := inventory.NewHostBuilder()
			host, _ := hostBuilder.
				WithName("linux").
				WithTransport(mockTransport).
				WithEscalateConfig(escalateConfig).
				Build()

			ctx := context.Background()

			config := &RunConfig{
				HostInfo:   host.Info(),
				Escalation: nil,
				Input:      cty.ObjectVal(tt.input),
			}

			r := assert.Run(ctx, config)
			if r == nil {
				t.Fatalf("Expected non-nil result from Run(), got nil")
			}

			if r.Failed != tt.expectedFailed {
				t.Errorf("Expected result from Run() to be %t, got %t", tt.expectedFailed, r.Failed)
			}
		})
	}
}
