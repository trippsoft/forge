// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"context"
	"testing"

	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/module/shell"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestShellRun_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	host, ok := inv.Host("linux")
	if !ok {
		t.Fatal("Host 'linux' not found in inventory")
	}

	p := &shell.Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	config := &module.RunConfig{
		Transport:  host.Transport(),
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
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

	host, ok := inv.Host("linuxpw")
	if !ok {
		t.Fatal("Host 'linuxpw' not found in inventory")
	}

	p := &shell.Module{}

	escalation := transport.NewEscalation(linuxPWPassword)
	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	config := &module.RunConfig{
		Transport:  host.Transport(),
		HostInfo:   host.Info(),
		Escalation: escalation,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
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

	host, ok := inv.Host("linux")
	if !ok {
		t.Fatal("Host 'linux' not found in inventory")
	}

	p := &shell.Module{}

	escalation := transport.NewNoPasswordEscalation()
	input := map[string]cty.Value{
		"command": cty.StringVal("echo 'Hello, World!'"),
	}

	config := &module.RunConfig{
		Transport:  host.Transport(),
		HostInfo:   host.Info(),
		Escalation: escalation,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
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

	host, ok := inv.Host("windows")
	if !ok {
		t.Fatal("Host 'windows' not found in inventory")
	}

	p := &shell.Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo hello"),
	}

	config := &module.RunConfig{
		Transport:  host.Transport(),
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
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

func TestShellRun_Windows_SSH_Cmd(t *testing.T) {
	setupVagrantEnvironment(t)

	host, ok := inv.Host("cmd")
	if !ok {
		t.Fatal("Host 'cmd' not found in inventory")
	}

	p := &shell.Module{}

	input := map[string]cty.Value{
		"command": cty.StringVal("echo hello"),
	}

	config := &module.RunConfig{
		Transport:  host.Transport(),
		HostInfo:   host.Info(),
		Escalation: nil,
		Input:      input,
	}

	result := p.Run(context.Background(), config)
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
