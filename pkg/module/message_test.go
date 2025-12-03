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
	module := &MessageModule{}

	spec := module.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec from InputSpec(), got nil")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got: %q", err.Error())
	}
}

func TestMessageModuleValidate(t *testing.T) {
	m := &MessageModule{}

	input := map[string]cty.Value{
		"message": cty.StringVal("Hello, World!"),
	}

	mockTransport := transport.NewMockTransport()
	escalateConfig := inventory.NewEscalateConfig("")
	hostBuilder := inventory.NewHostBuilder()
	host, _ := hostBuilder.
		WithName("linux").
		WithTransport(mockTransport).
		WithEscalateConfig(escalateConfig).
		Build()

	config := &RunConfig{
		Transport:  mockTransport,
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	err := m.Validate(config)
	if err != nil {
		t.Fatalf("Expected no error from Validate(), got: %q", err.Error())
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

	p := &MessageModule{}

	input := map[string]cty.Value{
		"message": cty.StringVal("Hello, World!"),
	}

	config := &RunConfig{
		Transport:  mockTransport,
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
	if result.Err != nil {
		t.Fatalf("Expected no error from Run(), got: %q", result.Err.Error())
	}

	if result.Changed {
		t.Fatal("Expected module to not indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output from Run(), got nil")
	}

	expectedMessage := "Hello, World!"
	if len(result.Messages) != 1 {
		t.Fatalf("Expected 1 message in result.Messages, got: %d", len(result.Messages))
	}

	if result.Messages[0] != expectedMessage {
		t.Fatalf("Expected message %q, got %q", expectedMessage, result.Messages[0])
	}

	if len(result.Output) != 0 {
		t.Fatalf("Expected output to have 0 keys, got: %d", len(result.Output))
	}
}
