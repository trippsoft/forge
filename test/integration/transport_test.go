// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/trippsoft/forge/pkg/transport"
)

func TestSSHTransportConnect_Linux_Password(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()
	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_Linux_PrivateKey(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PublicKeyAuth(linuxPrivateKey).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport with private key: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect with private key: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_WinPowerShell_Password(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_WinPowerShell_PrivateKey(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PublicKeyAuth(windowsPrivateKey).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport with private key: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect with private key: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_WinCmd_Password(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_WinCmd_PrivateKey(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PublicKeyAuth(cmdPrivateKey).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport with private key: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect with private key: %v", err)
	}

	sshTransport.Close()
}

func TestSSHTransportConnect_Failure(t *testing.T) {
	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host("192.0.2.1").
		Port(22).
		User("testuser").
		PasswordAuth("testpass").
		DontUseKnownHosts().
		ConnectionTimeout(2 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err == nil {
		t.Error("Expected connection to fail to unreachable host")
		sshTransport.Close()
	}
}

func TestSSHTransportKnownHosts_Strict(t *testing.T) {
	setupVagrantEnvironment(t)

	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		UseStrictKnownHosts(tmpKnownHosts).
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err == nil {
		t.Error("Expected connection to fail with strict known hosts checking and unknown host")
		sshTransport.Close()
		return
	}

	if !strings.Contains(err.Error(), "key is unknown") {
		t.Errorf("Got expected error (though message could be more specific): %v", err)
	}
}

func TestSSHTransportKnownHosts_Strict_RejectNotMatchingKey(t *testing.T) {
	setupVagrantEnvironment(t)

	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	file, err := os.Create(tmpKnownHosts)
	if err != nil {
		t.Fatalf("Failed to open temp known hosts file: %v", err)
	}

	content := fmt.Sprintf("%s ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n", linuxHost)
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp known hosts file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp known hosts file: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PublicKeyAuth(linuxPrivateKey). // Use a wrong key
		UseStrictKnownHosts(tmpKnownHosts).
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport with wrong key: %v", err)
	}

	err = sshTransport.Connect()
	if err == nil {
		t.Error("Expected connection to fail with wrong key")
		sshTransport.Close()
		return
	}

	if !strings.Contains(err.Error(), "key mismatch") {
		t.Errorf("Got unexpected error: %v", err)
	}
}

func TestSSHTransportKnownHosts_AddUnknown(t *testing.T) {
	setupVagrantEnvironment(t)

	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		UseKnownHosts(tmpKnownHosts). // Allow auto-adding unknown hosts
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Expected connection to succeed with auto-add known hosts: %v", err)
	}

	sshTransport.Close()

	// Test connection again with strict host checking to verify known host was added
	sshTransport, err = builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		UseStrictKnownHosts(tmpKnownHosts).
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build second SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Expected second connection to succeed with known host key: %v", err)
	}

	err = sshTransport.Close()
	if err != nil {
		t.Errorf("Second close failed: %v", err)
	}
}

func TestSSHTransportKnownHosts_AddUnknown_RejectNotMatchingKey(t *testing.T) {
	setupVagrantEnvironment(t)

	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	file, err := os.Create(tmpKnownHosts)
	if err != nil {
		t.Fatalf("Failed to open temp known hosts file: %v", err)
	}

	content := fmt.Sprintf("%s ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n", linuxHost)
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp known hosts file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp known hosts file: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PublicKeyAuth(linuxPrivateKey). // Use a wrong key
		UseKnownHosts(tmpKnownHosts).
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport with wrong key: %v", err)
	}

	err = sshTransport.Connect()
	if err == nil {
		t.Error("Expected connection to fail with wrong key")
		sshTransport.Close()
		return
	}

	if !strings.Contains(err.Error(), "key mismatch") {
		t.Errorf("Got unexpected error: %v", err)
	}
}

func TestSSHTransportKnownHosts_AddUnknown_NonExistentFile(t *testing.T) {
	setupVagrantEnvironment(t)

	nonExistentPath := "/tmp/non_existent_known_hosts"
	cleanupTempFile(t, nonExistentPath)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		UseKnownHosts(nonExistentPath).
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Expected connection to succeed and create known hosts file: %v", err)
	}

	err = sshTransport.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	cleanupTempFile(t, nonExistentPath)
}

func TestSSHTransportCommand_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewCommand("echo 'Hello from Linux'", nil)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v, stderr: %s", err, stderr)
	}

	expectedStdout := "Hello from Linux"
	if stdout != expectedStdout {
		t.Errorf("Expected stdout to be %q, got %q", expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("Expected stderr to be empty, got %q", stderr)
	}
}

func TestSSHTransportEscalatedCommand_Linux_SudoNoPassword(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	escalateConfig := transport.NewNoPasswordEscalation()
	cmd, err := sshTransport.NewCommand("echo 'Hello from Linux'", escalateConfig)
	if err != nil {
		t.Fatalf("Failed to create escalated command: %v", err)
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v, stderr: %s", err, stderr)
	}

	expectedStdout := "Hello from Linux"
	if stdout != expectedStdout {
		t.Errorf("Expected stdout to be %q, got %q", expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("Expected stderr to be empty, got %q", stderr)
	}
}

func TestSSHTransportEscalatedCommand_Linux_SudoPassword(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxPWHost).
		Port(linuxPWPort).
		User(linuxPWUser).
		PasswordAuth(linuxPWPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	escalateConfig := transport.NewEscalation(linuxPWPassword)
	cmd, err := sshTransport.NewCommand("echo 'Hello from Linux'", escalateConfig)
	if err != nil {
		t.Fatalf("Failed to create escalated command: %v", err)
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v, stderr: %s", err, stderr)
	}

	expectedStdout := "Hello from Linux"
	if stdout != expectedStdout {
		t.Errorf("Expected stdout to be %q, got %q", expectedStdout, stdout)
	}

	if !strings.Contains(stderr, "forge_sudo_prompt:") {
		t.Errorf("Expected stderr to contain sudo prompt, got %q", stderr)
	}
}

func TestSSHTransportCommand_WinPowerShell(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewCommand(`echo "Hello from Windows"`, nil)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v, stderr: %s", err, stderr)
	}

	expectedStdout := "Hello from Windows"
	if stdout != expectedStdout {
		t.Errorf("Expected stdout to be %q, got %q", expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("Expected stderr to be empty, got %q", stderr)
	}
}

func TestSSHTransportCommand_WinCmd(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewCommand("echo Hello from CMD", nil)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v, stderr: %s", err, stderr)
	}

	expectedStdout := "Hello from CMD"
	if stdout != expectedStdout {
		t.Errorf("Expected stdout to be %q, got %q", expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("Expected stderr to be empty, got %q", stderr)
	}
}

func TestSSHTransportPowerShell_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPowerShellCommand("Write-Host 'Hello from PowerShell'", nil)
	if err == nil {
		t.Fatal("Expected error for NewPowerShellCommand on Linux, but got none")
	}

	if cmd != nil {
		t.Fatal("Expected nil command for NewPowerShellCommand on Linux, but got a command")
	}

	expectedErr := "PowerShell is not available on the remote system"
	if err == nil {
		t.Fatal("Expected PowerShell command to fail on Linux, but it succeeded")
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErr, err.Error())
	}
}

func TestSSHTransportPowerShell_WinPowerShell(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPowerShellCommand("Write-Host 'Hello from PowerShell'", nil)
	if err != nil {
		t.Fatalf("NewPowerShellCommand failed: %v", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		t.Fatalf("ExecutePowerShell failed: %v", err)
	}

	expectedStdout := "Hello from PowerShell"
	if stdout != expectedStdout {
		t.Errorf("Expected PowerShell output to be %q, got %q", expectedStdout, stdout)
	}
}

func TestSSHTransportPowerShell_WinCmd(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPowerShellCommand("Write-Host 'Hello from PowerShell'", nil)
	if err != nil {
		t.Fatalf("NewPowerShellCommand failed: %v", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		t.Fatalf("ExecutePowerShell failed: %v", err)
	}

	expectedStdout := "Hello from PowerShell"
	if stdout != expectedStdout {
		t.Errorf("Expected PowerShell output to be %q, got %q", expectedStdout, stdout)
	}
}

func TestSSHTransportPython_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPythonCommand("/usr/bin/python3", "print('Hello from Python')", nil)
	if err != nil {
		t.Fatalf("NewPythonCommand failed: %v", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		t.Fatalf("ExecutePython failed: %v", err)
	}

	expectedStdout := "Hello from Python"
	if stdout != expectedStdout {
		t.Errorf("Expected Python output to be %q, got %q", expectedStdout, stdout)
	}
}

func TestSSHTransportPython_WinPowerShell(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPythonCommand("", "Write-Host 'Hello from PowerShell'", nil)
	if err == nil {
		t.Fatal("Expected error for NewPythonCommand on Windows, but got none")
	}

	if cmd != nil {
		t.Fatal("Expected nil command for NewPythonCommand on Windows, but got a command")
	}

	expectedErr := "Python is not available on the remote system"
	if err == nil {
		t.Fatal("Expected Python command to fail on Windows, but it succeeded")
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErr, err.Error())
	}
}

func TestSSHTransportPython_WinCmd(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	cmd, err := sshTransport.NewPythonCommand("", "Write-Host 'Hello from PowerShell'", nil)
	if err == nil {
		t.Fatal("Expected error for NewPythonCommand on Windows, but got none")
	}

	if cmd != nil {
		t.Fatal("Expected nil command for NewPythonCommand on Windows, but got a command")
	}

	expectedErr := "Python is not available on the remote system"
	if err == nil {
		t.Fatal("Expected Python command to fail on Windows, but it succeeded")
	}

	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErr, err.Error())
	}
}

func TestSSHTransportStat_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_stat_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	fileInfo, err := sshTransport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if fileInfo == nil {
		t.Error("Expected fileInfo to be non-nil for existing file")
	}
}

func TestSSHTransportStat_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_stat_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	fileInfo, err := sshTransport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if fileInfo == nil {
		t.Error("Expected fileInfo to be non-nil for existing file")
	}
}

func TestSSHTransportCreate_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_create_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	file, err := sshTransport.Create(tmpDir + "/testfile.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = file.Write([]byte("Hello from SSH Create"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	err = file.Sync()
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	file, err = sshTransport.Open(tmpDir + "/testfile.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if file == nil {
		t.Error("Expected file to be non-nil for existing file")
	}

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if string(buffer.String()) != "Hello from SSH Create" {
		t.Errorf("Expected file content to be 'Hello from SSH Create', got: %s", string(buffer.String()))
	}
}

func TestSSHTransportCreate_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_create_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	file, err := sshTransport.Create(tmpDir + "\\testfile.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = file.Write([]byte("Hello from SSH Create"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	err = file.Sync()
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	file, err = sshTransport.Open(tmpDir + "\\testfile.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if file == nil {
		t.Error("Expected file to be non-nil for existing file")
	}

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if buffer.String() != "Hello from SSH Create" {
		t.Errorf("Expected file content to be 'Hello from SSH Create', got: %s", string(buffer.String()))
	}
}

func TestSSHTransportOpen_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_open_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err = tmpFile.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	file, err := sshTransport.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if file == nil {
		t.Error("Expected file to be non-nil for existing file")
	}

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if string(buffer.String()) != "" { // Expecting empty content since we just created the file
		t.Errorf("Expected file content to be empty, got: %s", string(buffer.String()))
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestSSHTransportOpen_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_open_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err = tmpFile.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	file, err := sshTransport.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if file == nil {
		t.Error("Expected file to be non-nil for existing file")
	}

	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if buffer.String() != "" { // Expecting empty content since we just created the file
		t.Errorf("Expected file content to be empty, got: %s", buffer.String())
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestSSHTransportMkdir_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_mkdir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	path := tmpDir + "/newdir"
	err = sshTransport.Mkdir(path)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if fileInfo == nil {
		t.Error("Expected fileInfo to be non-nil for existing directory")
	}

	if !fileInfo.IsDir() {
		t.Error("Expected fileInfo to be a directory")
	}
}

func TestSSHTransportMkdir_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_mkdir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	path := tmpDir + "\\newdir"
	err = sshTransport.Mkdir(path)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	if fileInfo == nil {
		t.Error("Expected fileInfo to be non-nil for existing directory")
	}

	if !fileInfo.IsDir() {
		t.Error("Expected fileInfo to be a directory")
	}
}

func TestSSHTransportRemove_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_remove_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err = tmpFile.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	err = sshTransport.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Stat failed after remove: %v", err)
	}

	if fileInfo != nil {
		t.Error("Expected fileInfo to be nil for removed file")
	}
}

func TestSSHTransportRemove_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_remove_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	err = tmpFile.Sync()
	if err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	err = sshTransport.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Stat failed after remove: %v", err)
	}

	if fileInfo != nil {
		t.Error("Expected fileInfo to be nil for removed file")
	}
}

func TestSSHTransportRemoveAll_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_removeall_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	file, err := sshTransport.Create(tmpDir + "/testfile.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	err = sshTransport.RemoveAll(tmpDir)
	if err != nil {
		t.Fatalf("RemoveAll failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Stat failed after RemoveAll: %v", err)
	}

	if fileInfo != nil {
		t.Error("Expected fileInfo to be nil for removed directory")
	}
}

func TestSSHTransportRemoveAll_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_removeall_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	file, err := sshTransport.Create(tmpDir + "\\testfile.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	err = sshTransport.RemoveAll(tmpDir)
	if err != nil {
		t.Fatalf("RemoveAll failed: %v", err)
	}

	fileInfo, err := sshTransport.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Stat failed after RemoveAll: %v", err)
	}

	if fileInfo != nil {
		t.Error("Expected fileInfo to be nil for removed directory")
	}
}

func TestSSHTransportJoin_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	path := sshTransport.Join("/tmp", "testfile.txt")
	expectedPath := "/tmp/testfile.txt"

	if path != expectedPath {
		t.Errorf("Expected joined path to be '%s', got '%s'", expectedPath, path)
	}
}

func TestSSHTransportJoin_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	path := sshTransport.Join("C:\\temp", "testfile.txt")
	expectedPath := "C:\\temp\\testfile.txt"

	if path != expectedPath {
		t.Errorf("Expected joined path to be '%s', got '%s'", expectedPath, path)
	}
}

func TestSSHTransportTempDir_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.TempDir()
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	if tmpDir != "/tmp" {
		t.Errorf("Expected temp dir to be '/tmp', got: %s", tmpDir)
	}
}

func TestSSHTransportTempDir_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.TempDir()
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	if tmpDir != "C:\\Users\\vagrant\\AppData\\Local\\Temp" {
		t.Errorf("Expected temp dir to be 'C:\\Users\\vagrant\\AppData\\Local\\Temp', got: %s", tmpDir)
	}
}

func TestSSHTransportCreateTemp_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_create_temp_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	if !strings.HasPrefix(tmpFile.Name(), "/tmp/test_ssh_create_temp_") {
		t.Errorf("Expected temp file to start with '/tmp/test_ssh_create_temp_', got: %s", tmpFile.Name())
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
}

func TestSSHTransportCreateTemp_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_create_temp_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	if !strings.HasPrefix(tmpFile.Name(), "C:\\Users\\vagrant\\AppData\\Local\\Temp\\test_ssh_create_temp_") {
		t.Errorf("Expected temp file to start with 'C:\\Users\\vagrant\\AppData\\Local\\Temp\\test_ssh_create_temp_', got: %s", tmpFile.Name())
	}

	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
}

func TestSSHTransportMkdirTemp_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_mkdir_temp_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	if !strings.HasPrefix(tmpDir, "/tmp/test_ssh_mkdir_temp_") {
		t.Errorf("Expected temp dir to start with '/tmp/test_ssh_mkdir_temp_', got: %s", tmpDir)
	}
}

func TestSSHTransportMkdirTemp_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpDir, err := sshTransport.MkdirTemp("", "test_ssh_mkdir_temp_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer sshTransport.RemoveAll(tmpDir)

	if !strings.HasPrefix(tmpDir, "C:\\Users\\vagrant\\AppData\\Local\\Temp\\test_ssh_mkdir_temp_") {
		t.Errorf("Expected temp dir to start with 'C:\\Users\\vagrant\\AppData\\Local\\Temp\\test_ssh_mkdir_temp_', got: %s", tmpDir)
	}
}

func TestSSHTransportSymlink_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	// Create a temporary file for testing
	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_symlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	// Create a symlink to the temporary file
	symlinkPath := tmpFile.Name() + "_symlink"
	err = sshTransport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	defer sshTransport.Remove(symlinkPath)

	// Verify the symlink points to the correct target
	target, err := sshTransport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestSSHTransportSymlink_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	// Create a temporary file for testing
	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_symlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	// Create a symlink to the temporary file
	symlinkPath := tmpFile.Name() + "_symlink"
	err = sshTransport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	defer sshTransport.Remove(symlinkPath)

	// Verify the symlink points to the correct target
	target, err := sshTransport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestSSHTransportReadLink_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_readlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	symlinkPath := tmpFile.Name() + "_symlink"
	err = sshTransport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	defer sshTransport.Remove(symlinkPath)

	target, err := sshTransport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestSSHTransportReadLink_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	tmpFile, err := sshTransport.CreateTemp("", "test_ssh_readlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer sshTransport.Remove(tmpFile.Name())

	symlinkPath := tmpFile.Name() + "_symlink"
	err = sshTransport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	defer sshTransport.Remove(symlinkPath)

	target, err := sshTransport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestSSHTransportRealPath_Linux(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	path, err := sshTransport.RealPath("sh")
	if err != nil {
		t.Fatalf("RealPath failed: %v", err)
	}

	if path != "/bin/sh" && path != "/usr/bin/sh" {
		t.Errorf("Expected real path to be '/bin/sh' or '/usr/bin/sh', got '%s'", path)
	}
}

func TestSSHTransportRealPath_Linux_NotFound(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	_, err = sshTransport.RealPath("nonexistent_command")
	if err == nil {
		t.Fatal("Expected RealPath to return an error for nonexistent command, but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error to be os.ErrNotExist, got: %v", err)
	}
}

func TestSSHTransportRealPath_Windows(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	path, err := sshTransport.RealPath("cmd.exe")
	if err != nil {
		t.Fatalf("RealPath failed: %v", err)
	}

	if path != "C:\\WINDOWS\\system32\\cmd.exe" {
		t.Errorf("Expected real path 'C:\\WINDOWS\\system32\\cmd.exe', got '%s'", path)
	}
}

func TestSSHTransportRealPath_Windows_NotFound(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	_, err = sshTransport.RealPath("nonexistent_command")
	if err == nil {
		t.Fatal("Expected RealPath to return an error for nonexistent command, but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error to be os.ErrNotExist, got: %v", err)
	}
}

func TestSSHTransportRealPath_Cmd(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	path, err := sshTransport.RealPath("cmd.exe")
	if err != nil {
		t.Fatalf("RealPath failed: %v", err)
	}

	if path != "C:\\WINDOWS\\system32\\cmd.exe" {
		t.Errorf("Expected real path 'C:\\WINDOWS\\system32\\cmd.exe', got '%s'", path)
	}
}

func TestSSHTransportRealPath_Cmd_NotFound(t *testing.T) {
	setupVagrantEnvironment(t)

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	sshTransport, err := builder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
		DontUseKnownHosts().
		ConnectionTimeout(30 * time.Second).
		Build()

	if err != nil {
		t.Fatalf("Failed to build SSH transport: %v", err)
	}

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	defer sshTransport.Close()

	_, err = sshTransport.RealPath("nonexistent_command")
	if err == nil {
		t.Fatal("Expected RealPath to return an error for nonexistent command, but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error to be os.ErrNotExist, got: %v", err)
	}
}

// Helper functions for known hosts testing

func createTempKnownHostsFile(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_known_hosts_empty_*")
	if err != nil {
		t.Fatalf("Failed to create empty temp known hosts file: %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

func cleanupTempFile(t *testing.T, path string) {
	t.Helper()

	err := os.Remove(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		t.Logf("Warning: failed to cleanup temp file %s: %v", path, err)
	}
}
