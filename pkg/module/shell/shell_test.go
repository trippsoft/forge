// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package shell

import (
	"context"
	"fmt"
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
		"command": cty.StringVal("echo 'Hello, World!'"),
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
	command := "echo 'Hello, World!'"
	expectedStdout := "Hello, World!"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[command] = &transport.MockCmd{
		Stdout: fmt.Sprintf("%s\n", expectedStdout),
	}

	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("linux", mockTransport, escalateConfig, map[string]cty.Value{})

	p := &Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
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

	if !result.Changed {
		t.Fatal("Expected module to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output from Run(), got nil")
	}

	if len(result.Output) != 2 {
		t.Fatalf("Expected output to have 2 keys, got: %d", len(result.Output))
	}

	if _, ok := result.Output["stdout"]; !ok {
		t.Fatal("Expected output to contain 'stdout' key")
	}

	if _, ok := result.Output["stderr"]; !ok {
		t.Fatal("Expected output to contain 'stderr' key")
	}

	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}

	actualStderr := result.Output["stderr"].AsString()
	if actualStderr != "" {
		t.Fatalf("Expected stderr to be empty, got %q", actualStderr)
	}
}
