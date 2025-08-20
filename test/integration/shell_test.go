// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/module/shell"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestShellRun_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()

	sshTransport, err := builder.Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PublicKeyAuth(linuxPrivateKey).
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}

	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("linux", sshTransport, escalateConfig, map[string]cty.Value{})

	p := &shell.Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := &module.RunConfig{
		Transport:  sshTransport,
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(ctx, config)
	if result.Err != nil {
		t.Fatalf("Expected no error, got: %v", result.Err)
	}

	if !result.Changed {
		t.Fatal("Expected plugin to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output")
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

	expectedStdout := "Hello, World!"
	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}

	expectedStderr := ""
	actualStderr := result.Output["stderr"].AsString()
	if actualStderr != expectedStderr {
		t.Fatalf("Expected stderr %q, got %q", expectedStderr, actualStderr)
	}
}

func TestShellRun_Linux_SudoPassword(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()

	sshTransport, err := builder.Host(linuxPWHost).
		Port(linuxPWPort).
		User(linuxPWUser).
		PublicKeyAuth(linuxPWPrivateKey).
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}

	escalateConfig := inventory.NewEscalateConfig(linuxPWPassword)
	host := inventory.NewHost("linux", sshTransport, escalateConfig, map[string]cty.Value{})

	p := &shell.Module{}

	escalation := transport.NewEscalation(linuxPWPassword)
	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := &module.RunConfig{
		Transport:  sshTransport,
		HostInfo:   host.Info(),
		Escalation: escalation,
		Input:      input,
	}

	result := p.Run(ctx, config)
	if result.Err != nil {
		t.Fatalf("Expected no error, got: %v", result.Err)
	}

	if !result.Changed {
		t.Fatal("Expected plugin to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output")
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

	expectedStdout := "Hello, World!"
	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}
}

func TestShellRun_Linux_NoSudoPassword(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()

	sshTransport, err := builder.Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PublicKeyAuth(linuxPrivateKey).
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}

	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("linux", sshTransport, escalateConfig, map[string]cty.Value{})

	p := &shell.Module{}

	escalation := transport.NewNoPasswordEscalation()
	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := &module.RunConfig{
		Transport:  sshTransport,
		HostInfo:   host.Info(),
		Escalation: escalation,
		Input:      input,
	}

	result := p.Run(ctx, config)
	if result.Err != nil {
		t.Fatalf("Expected no error, got: %v", result.Err)
	}

	if !result.Changed {
		t.Fatal("Expected plugin to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output")
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

	expectedStdout := "Hello, World!"
	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}

	expectedStderr := ""
	actualStderr := result.Output["stderr"].AsString()
	if actualStderr != expectedStderr {
		t.Fatalf("Expected stderr %q, got %q", expectedStderr, actualStderr)
	}
}

func TestShellRun_Windows_SSH_PowerShell(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()

	sshTransport, err := builder.Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PublicKeyAuth(windowsPrivateKey).
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Failed to create transport: %v", err)
	}

	escalateConfig := inventory.NewEscalateConfig("")
	host := inventory.NewHost("windows", sshTransport, escalateConfig, map[string]cty.Value{})

	p := &shell.Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo hello"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := &module.RunConfig{
		Transport:  sshTransport,
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(ctx, config)
	if result.Err != nil {
		t.Fatalf("Expected no error, got: %v", result.Err)
	}

	if !result.Changed {
		t.Fatal("Expected plugin to indicate changes were made")
	}

	if result.Output == nil {
		t.Fatal("Expected non-nil output")
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

	expectedStdout := "hello"
	actualStdout := result.Output["stdout"].AsString()
	if actualStdout != expectedStdout {
		t.Fatalf("Expected stdout %q, got %q", expectedStdout, actualStdout)
	}

	expectedStderr := ""
	actualStderr := result.Output["stderr"].AsString()
	if actualStderr != expectedStderr {
		t.Fatalf("Expected stderr %q, got %q", expectedStderr, actualStderr)
	}
}
