// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/trippsoft/forge/pkg/transport"
)

func TestSSHTransportConnect_PowerShell_Password(t *testing.T) {
	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPasswordAuth(windowsPassword).
		DontUseKnownHosts().
		WithConnectionTimeout(30 * time.Second).
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

func TestSSHTransportConnect_PowerShell_PrivateKey(t *testing.T) {
	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPublicKeyAuth(windowsPrivateKey).
		DontUseKnownHosts().
		WithConnectionTimeout(30 * time.Second).
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

func TestSSHTransportConnect_Cmd_Password(t *testing.T) {
	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(cmdHost).
		WithPort(cmdPort).
		WithUser(cmdUser).
		WithPasswordAuth(cmdPassword).
		DontUseKnownHosts().
		WithConnectionTimeout(30 * time.Second).
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

func TestSSHTransportConnect_Cmd_PrivateKey(t *testing.T) {
	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(cmdHost).
		WithPort(cmdPort).
		WithUser(cmdUser).
		WithPublicKeyAuth(cmdPrivateKey).
		DontUseKnownHosts().
		WithConnectionTimeout(30 * time.Second).
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
	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost("192.0.2.1").
		WithPort(22).
		WithUser("testuser").
		WithPasswordAuth("testpass").
		DontUseKnownHosts().
		WithConnectionTimeout(2 * time.Second).
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
	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPasswordAuth(windowsPassword).
		UseStrictKnownHosts(tmpKnownHosts).
		WithConnectionTimeout(30 * time.Second).
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
	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	file, err := os.Create(tmpKnownHosts)
	if err != nil {
		t.Fatalf("Failed to open temp known hosts file: %v", err)
	}

	content := fmt.Sprintf("%s ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n", windowsHost)
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp known hosts file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp known hosts file: %v", err)
	}

	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPublicKeyAuth(windowsPrivateKey). // Use a wrong key
		UseStrictKnownHosts(tmpKnownHosts).
		WithConnectionTimeout(30 * time.Second).
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
	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPasswordAuth(windowsPassword).
		UseKnownHosts(tmpKnownHosts). // Allow auto-adding unknown hosts
		WithConnectionTimeout(30 * time.Second).
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
	sshTransport, err = transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPasswordAuth(windowsPassword).
		UseStrictKnownHosts(tmpKnownHosts).
		WithConnectionTimeout(30 * time.Second).
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
	tmpKnownHosts := createTempKnownHostsFile(t)
	defer cleanupTempFile(t, tmpKnownHosts)

	file, err := os.Create(tmpKnownHosts)
	if err != nil {
		t.Fatalf("Failed to open temp known hosts file: %v", err)
	}

	content := fmt.Sprintf("%s ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n", windowsHost)
	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp known hosts file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("Failed to close temp known hosts file: %v", err)
	}

	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPublicKeyAuth(windowsPrivateKey). // Use a wrong key
		UseKnownHosts(tmpKnownHosts).
		WithConnectionTimeout(30 * time.Second).
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
	nonExistentPath := "/tmp/non_existent_known_hosts"
	cleanupTempFile(t, nonExistentPath)

	sshTransport, err := transport.
		NewSSHBuilder().
		WithHost(windowsHost).
		WithPort(windowsPort).
		WithUser(windowsUser).
		WithPasswordAuth(windowsPassword).
		UseKnownHosts(nonExistentPath).
		WithConnectionTimeout(30 * time.Second).
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
