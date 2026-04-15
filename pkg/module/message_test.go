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

func TestMessageModuleInputSpec(t *testing.T) {
	spec := message.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec from InputSpec(), got nil")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got: %q", err.Error())
	}
}

func TestMessageModuleRun(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	escalateConfig := inventory.NewEscalateConfig("")
	hostBuilder := inventory.NewHostBuilder()
	host, _ := hostBuilder.
		WithName("linux").
		WithTransport(mockTransport).
		WithEscalateConfig(escalateConfig).
		Build()

	input := map[string]cty.Value{
		"message": cty.StringVal("Hello, World!"),
	}

	config := &RunConfig{
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	r := message.Run(context.Background(), config)
	if r.Error != nil {
		t.Fatalf("Expected no error from Run(), got: %q", r.Error.Error())
	}

	if r.Changed {
		t.Fatal("Expected module to not indicate changes were made")
	}

	if !r.Output.IsWhollyKnown() || r.Output.IsNull() {
		t.Fatal("Expected non-nil output from Run(), got nil")
	}

	expectedMessage := "Hello, World!"
	if len(r.Messages) != 1 {
		t.Fatalf("Expected 1 message in result.Messages, got: %d", len(r.Messages))
	}

	if r.Messages[0] != expectedMessage {
		t.Fatalf("Expected message %q, got %q", expectedMessage, r.Messages[0])
	}

	outputMap := r.Output.AsValueMap()
	if len(outputMap) != 0 {
		t.Fatalf("Expected empty output map, got: %v", outputMap)
	}
}
