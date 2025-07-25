package transport

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewSSHBuilder(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	if builder == nil {
		t.Fatal("NewSSHBuilder returned nil builder")
	}

	// Check default values
	if builder.port != DefaultSSHPort {
		t.Errorf("Expected default port %d, got %d", DefaultSSHPort, builder.port)
	}

	if builder.connectionTimeout != DefaultSSHConnectionTimeout {
		t.Errorf("Expected default timeout %v, got %v", DefaultSSHConnectionTimeout, builder.connectionTimeout)
	}

	if builder.knownHostsPath == "" {
		t.Error("Expected non-empty default known hosts path")
	}
}

func TestDefaultKnownHostsPath(t *testing.T) {
	path, err := DefaultKnownHostsPath()
	if err != nil {
		t.Fatalf("DefaultKnownHostsPath failed: %v", err)
	}

	if path == "" {
		t.Error("DefaultKnownHostsPath returned empty path")
	}

	// Should end with .ssh/known_hosts
	expectedSuffix := filepath.Join(".ssh", "known_hosts")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}
}

func TestSSHTransportBuilderChaining(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Test method chaining
	result := builder.
		Host("example.com").
		Port(2222).
		User("testuser").
		PasswordAuth("testpass").
		ConnectionTimeout(5 * time.Second)

	if result != builder {
		t.Error("Builder methods should return the same builder instance for chaining")
	}

	// Verify values were set
	if builder.host != "example.com" {
		t.Errorf("Expected host 'example.com', got '%s'", builder.host)
	}
	if builder.port != 2222 {
		t.Errorf("Expected port 2222, got %d", builder.port)
	}
	if builder.user != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", builder.user)
	}
	if !builder.passwordAuth {
		t.Error("Expected password auth to be enabled")
	}
	if builder.password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", builder.password)
	}
	if builder.connectionTimeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", builder.connectionTimeout)
	}
}

func TestSSHTransportBuilderPublicKeyAuth(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Generate a test private key
	privateKey, err := generateTestPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate test private key: %v", err)
	}

	// Test public key auth without passphrase
	builder.PublicKeyAuth(privateKey)
	if !builder.publicKeyAuth {
		t.Error("Expected public key auth to be enabled")
	}
	if string(builder.privateKey) != string(privateKey) {
		t.Error("Private key not set correctly")
	}
	if builder.privateKeyPass != "" {
		t.Errorf("Expected empty passphrase, got '%s'", builder.privateKeyPass)
	}

	// Test public key auth with passphrase
	builder.PublicKeyAuthWithPass(privateKey, "testpass")
	if !builder.publicKeyAuth {
		t.Error("Expected public key auth to be enabled")
	}
	if builder.privateKeyPass != "testpass" {
		t.Errorf("Expected passphrase 'testpass', got '%s'", builder.privateKeyPass)
	}

	// Test disabling public key auth
	builder.NoPublicKeyAuth()
	if builder.publicKeyAuth {
		t.Error("Expected public key auth to be disabled")
	}
	if builder.privateKey != nil {
		t.Error("Expected private key to be cleared")
	}
	if builder.privateKeyPass != "" {
		t.Error("Expected passphrase to be cleared")
	}
}

func TestSSHTransportBuilderKnownHosts(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Test disabling known hosts
	builder.DontUseKnownHosts()
	if builder.useKnownHostsFile {
		t.Error("Expected known hosts to be disabled")
	}

	// Test enabling known hosts with custom path
	customPath := "/custom/known_hosts"
	builder.UseKnownHosts(customPath, true)
	if !builder.useKnownHostsFile {
		t.Error("Expected known hosts to be enabled")
	}
	if builder.knownHostsPath != customPath {
		t.Errorf("Expected known hosts path '%s', got '%s'", customPath, builder.knownHostsPath)
	}
	if !builder.addUnknownHostsToFile {
		t.Error("Expected addUnknownHostsToFile to be true")
	}
}

func TestSSHTransportBuilderValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupBuilder  func(*SSHTransportBuilder)
		expectedError string
	}{
		{
			name:          "empty host",
			setupBuilder:  func(b *SSHTransportBuilder) {},
			expectedError: "host cannot be empty",
		},
		{
			name: "zero port",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(0)
			},
			expectedError: "port must be between 1 and 65535",
		},
		{
			name: "empty user",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22)
			},
			expectedError: "user cannot be empty",
		},
		{
			name: "public key auth without key",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22).User("test")
				b.publicKeyAuth = true
				b.privateKey = nil
			},
			expectedError: "privateKey cannot be empty when public key authentication is enabled",
		},
		{
			name: "password auth without password",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22).User("test")
				b.passwordAuth = true
				b.password = ""
			},
			expectedError: "password cannot be empty when password authentication is enabled",
		},
		{
			name: "known hosts without path",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22).User("test").PasswordAuth("pass")
				b.useKnownHostsFile = true
				b.knownHostsPath = ""
			},
			expectedError: "knownHostsPath cannot be empty when using known hosts",
		},
		{
			name: "zero timeout",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22).User("test").PasswordAuth("pass").ConnectionTimeout(0)
			},
			expectedError: "connectionTimeout must be greater than zero",
		},
		{
			name: "negative timeout",
			setupBuilder: func(b *SSHTransportBuilder) {
				b.Host("example.com").Port(22).User("test").PasswordAuth("pass").ConnectionTimeout(-1 * time.Second)
			},
			expectedError: "connectionTimeout must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := NewSSHBuilder()
			if err != nil {
				t.Fatalf("NewSSHBuilder failed: %v", err)
			}

			tt.setupBuilder(builder)

			_, err = builder.Build()
			if err == nil {
				t.Error("Expected error but got none")
			}

			if err.Error() != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}

func TestSSHTransportBuilderValidBuild(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Create a valid configuration
	transport, err := builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PasswordAuth("testpass").
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if transport == nil {
		t.Fatal("Build returned nil transport")
	}

	if transport.Type() != TransportTypeSSH {
		t.Errorf("Expected transport type %s, got %s", TransportTypeSSH, transport.Type())
	}
}

func TestSSHTransportBuilderWithPrivateKey(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Generate a test private key
	privateKey, err := generateTestPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate test private key: %v", err)
	}

	// Create transport with public key auth
	transport, err := builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PublicKeyAuth(privateKey).
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Build with public key failed: %v", err)
	}

	if transport == nil {
		t.Fatal("Build returned nil transport")
	}
}

func TestSSHTransportBuilderWithEncryptedPrivateKey(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Generate a test encrypted private key
	privateKey, err := generateTestEncryptedPrivateKey("testpass")
	if err != nil {
		t.Fatalf("Failed to generate test encrypted private key: %v", err)
	}

	// Create transport with encrypted public key auth
	transport, err := builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PublicKeyAuthWithPass(privateKey, "testpass").
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Build with encrypted public key failed: %v", err)
	}

	if transport == nil {
		t.Fatal("Build returned nil transport")
	}
}

func TestSSHTransportBuilderInvalidPrivateKey(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Try to build with invalid private key
	_, err = builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PublicKeyAuth([]byte("invalid key data")).
		DontUseKnownHosts().
		Build()

	if err == nil {
		t.Error("Expected error for invalid private key")
	}

	expectedError := "failed to parse private key"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHTransportBuilderInvalidEncryptedPrivateKey(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	// Generate a test encrypted private key
	privateKey, err := generateTestEncryptedPrivateKey("correctpass")
	if err != nil {
		t.Fatalf("Failed to generate test encrypted private key: %v", err)
	}

	// Try to build with wrong passphrase
	_, err = builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PublicKeyAuthWithPass(privateKey, "wrongpass").
		DontUseKnownHosts().
		Build()

	if err == nil {
		t.Error("Expected error for wrong passphrase")
	}

	expectedError := "failed to parse private key with passphrase"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHTransportType(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	transport, err := builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PasswordAuth("testpass").
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if transport.Type() != TransportTypeSSH {
		t.Errorf("Expected transport type %s, got %s", TransportTypeSSH, transport.Type())
	}
}

func TestSSHTransportFileSystem(t *testing.T) {
	builder, err := NewSSHBuilder()
	if err != nil {
		t.Fatalf("NewSSHBuilder failed: %v", err)
	}

	transport, err := builder.
		Host("example.com").
		Port(22).
		User("testuser").
		PasswordAuth("testpass").
		DontUseKnownHosts().
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	fs := transport.FileSystem()
	if fs == nil {
		t.Fatal("FileSystem returned nil")
	}

	// Test that it's actually an sftpFileSystem
	_, ok := fs.(*sftpFileSystem)
	if !ok {
		t.Error("FileSystem did not return an sftpFileSystem instance")
	}

	// Test that multiple calls return the same instance
	fs2 := transport.FileSystem()
	if fs != fs2 {
		t.Error("Multiple FileSystem calls should return the same instance")
	}
}

// Helper functions

func generateTestPrivateKey() ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(privateKeyPEM), nil
}

func generateTestEncryptedPrivateKey(passphrase string) ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	encryptedPEM, err := x509.EncryptPEMBlock(rand.Reader, privateKeyPEM.Type, privateKeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(encryptedPEM), nil
}
