// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package message

import (
	"context"
	"testing"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestModuleInputSpec(t *testing.T) {
	module := &Module{}

	spec := module.InputSpec()
	if spec == nil {
		t.Fatal("Expected non-nil input spec from InputSpec(), got nil")
	}

	err := spec.ValidateSpec()
	if err != nil {
		t.Errorf("expected no errors from ValidateSpec(), got: %q", err.Error())
	}
}

func TestModuleValidate(t *testing.T) {
	m := &Module{}

	input := map[string]cty.Value{
		"message": cty.StringVal("Hello, World!"),
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
		t.Fatalf("Expected no error from Validate(), got: %q", err.Error())
	}
}

func TestModuleRun(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("linux", mockTransport, escalateConfig, map[string]cty.Value{})

	p := &Module{}

	input := map[string]cty.Value{
		"message": cty.StringVal("Hello, World!"),
	}

	config := &module.RunConfig{
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
	if result.Message != expectedMessage {
		t.Fatalf("Expected message %q, got %q", expectedMessage, result.Message)
	}

	if len(result.Output) != 0 {
		t.Fatalf("Expected output to have 0 keys, got: %d", len(result.Output))
	}
}
