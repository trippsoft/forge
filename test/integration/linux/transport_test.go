// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package linux

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/trippsoft/forge/pkg/transport"
)

type containerConfig struct {
	name string
	port uint16
}

func getContainerConfigs(t *testing.T) []containerConfig {
	t.Helper()

	return []containerConfig{
		{name: "debian13", port: 2201},
		{name: "debian12", port: 2202},
		{name: "fedora42", port: 2211},
		{name: "fedora41", port: 2212},
		{name: "rocky10", port: 2221},
		{name: "rocky9", port: 2222},
		{name: "rocky8", port: 2223},
		{name: "ubuntu2404", port: 2231},
		{name: "ubuntu2204", port: 2232},
	}
}

func TestSSHTransportConnect_Password(t *testing.T) {
	for _, config := range getContainerConfigs(t) {
		t.Run(config.name, func(t *testing.T) {
			sshTransport, err := transport.
				NewSSHBuilder().
				WithHost("127.0.0.1").
				WithPort(config.port).
				WithUser("forge").
				WithPasswordAuth("forge").
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
		})
	}
}

func TestSSHTransportConnect_PrivateKey(t *testing.T) {
	for _, config := range getContainerConfigs(t) {
		t.Run(config.name, func(t *testing.T) {
			sshTransport, err := transport.
				NewSSHBuilder().
				WithHost("127.0.0.1").
				WithPort(config.port).
				WithUser("forge").
				WithPublicKeyAuth(privateKeyContent).
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
		})
	}
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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPasswordAuth("forge").
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

	content := "127.0.0.1:2201 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n"
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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPublicKeyAuth(privateKeyContent). // Use a wrong key
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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPasswordAuth("forge").
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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPasswordAuth("forge").
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

	content := "127.0.0.1:2201 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBH3kZFNyb8iutKv6WzIA5Z1W+TqjwLU/kxRnFBLnLjBo5sXGkbAwUZd8xN7u4nF+OPdFwk9yfJ5ZHzvlsYYXowI=\n"

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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPublicKeyAuth(privateKeyContent). // Use a wrong key
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
		WithHost("127.0.0.1").
		WithPort(2201).
		WithUser("forge").
		WithPasswordAuth("forge").
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
