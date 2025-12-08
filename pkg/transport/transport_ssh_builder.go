// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"errors"
	"fmt"
	"time"

	"github.com/trippsoft/forge/pkg/plugin"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHTransportBuilder is a builder for constructing SSH transport instances.
type SSHTransportBuilder struct {
	host string
	port uint16

	user string

	publicKeyAuth  bool
	privateKey     []byte
	privateKeyPass string

	passwordAuth bool
	password     string

	useKnownHostsFile     bool
	knownHostsPath        string
	addUnknownHostsToFile bool

	connectionTimeout time.Duration

	minPluginPort uint16
	maxPluginPort uint16

	tempPath string
}

// WithHost sets the host for the SSH transport.
func (b *SSHTransportBuilder) WithHost(host string) *SSHTransportBuilder {
	b.host = host
	return b
}

// WithPort sets the port for the SSH transport.
func (b *SSHTransportBuilder) WithPort(port uint16) *SSHTransportBuilder {
	b.port = port
	return b
}

// WithUser sets the user for the SSH transport.
func (b *SSHTransportBuilder) WithUser(user string) *SSHTransportBuilder {
	b.user = user
	return b
}

// WithoutPublicKeyAuth disables public key authentication for the SSH transport.
func (b *SSHTransportBuilder) WithoutPublicKeyAuth() *SSHTransportBuilder {
	b.publicKeyAuth = false
	b.privateKey = nil
	b.privateKeyPass = ""
	return b
}

// WithPublicKeyAuth enables public key authentication for the SSH transport.
func (b *SSHTransportBuilder) WithPublicKeyAuth(privateKey []byte) *SSHTransportBuilder {
	b.publicKeyAuth = true
	b.privateKey = privateKey
	b.privateKeyPass = ""
	return b
}

// WithPublicKeyAuthWithPass enables public key authentication with a passphrase for the SSH transport.
func (b *SSHTransportBuilder) WithPublicKeyAuthWithPass(privateKey []byte, privateKeyPass string) *SSHTransportBuilder {
	b.publicKeyAuth = true
	b.privateKey = privateKey
	b.privateKeyPass = privateKeyPass
	return b
}

// WithPasswordAuth disables password authentication for the SSH transport.
func (b *SSHTransportBuilder) WithPasswordAuth(password string) *SSHTransportBuilder {
	b.passwordAuth = true
	b.password = password
	return b
}

// WithoutPasswordAuth disables password authentication for the SSH transport.
func (b *SSHTransportBuilder) WithoutPasswordAuth() *SSHTransportBuilder {
	b.passwordAuth = false
	b.password = ""
	return b
}

// DontUseKnownHosts disables the use of a known hosts file for the SSH transport.
func (b *SSHTransportBuilder) DontUseKnownHosts() *SSHTransportBuilder {
	b.useKnownHostsFile = false
	return b
}

// UseKnownHosts enables the use of a known hosts file for the SSH transport, adding unknown hosts to the file.
func (b *SSHTransportBuilder) UseKnownHosts(knownHostsPath string) *SSHTransportBuilder {
	b.useKnownHostsFile = true
	b.knownHostsPath = knownHostsPath
	b.addUnknownHostsToFile = true
	return b
}

// UseStrictKnownHosts enables the use of a known hosts file for the SSH transport, enforcing strict host key checking.
func (b *SSHTransportBuilder) UseStrictKnownHosts(knownHostsPath string) *SSHTransportBuilder {
	b.useKnownHostsFile = true
	b.knownHostsPath = knownHostsPath
	b.addUnknownHostsToFile = false
	return b
}

// WithConnectionTimeout sets the connection timeout for the SSH transport.
func (b *SSHTransportBuilder) WithConnectionTimeout(timeout time.Duration) *SSHTransportBuilder {
	b.connectionTimeout = timeout
	return b
}

// WithPluginPortRange sets the plugin port range for the SSH transport.
func (b *SSHTransportBuilder) WithPluginPortRange(minPluginPort, maxPluginPort uint16) *SSHTransportBuilder {
	b.minPluginPort = minPluginPort
	b.maxPluginPort = maxPluginPort
	return b
}

// WithTempPath sets the temporary path for the SSH transport.
func (b *SSHTransportBuilder) WithTempPath(tempPath string) *SSHTransportBuilder {
	b.tempPath = tempPath
	return b
}

// Build constructs the SSHTransport based on the builder's configuration.
func (b *SSHTransportBuilder) Build() (Transport, error) {
	if b.host == "" {
		return nil, errors.New("host cannot be empty")
	}

	if b.port == 0 {
		return nil, errors.New("port must be between 1 and 65535")
	}

	if b.user == "" {
		return nil, errors.New("user cannot be empty")
	}

	if b.publicKeyAuth && b.privateKey == nil {
		return nil, errors.New("privateKey cannot be empty when public key authentication is enabled")
	}

	if b.passwordAuth && b.password == "" {
		return nil, errors.New("password cannot be empty when password authentication is enabled")
	}

	if b.useKnownHostsFile && b.knownHostsPath == "" {
		return nil, errors.New("knownHostsPath cannot be empty when using known hosts")
	}

	if b.connectionTimeout <= 0 {
		return nil, errors.New("connectionTimeout must be greater than zero")
	}

	if b.minPluginPort > 0 && b.minPluginPort < 1024 {
		return nil, errors.New("minPluginPort must be at least 1024")
	}

	if b.maxPluginPort > 0 && b.maxPluginPort < 1024 {
		return nil, errors.New("maxPluginPort must be at least 1024")
	}

	if b.minPluginPort == 0 {
		if b.maxPluginPort != 0 && b.maxPluginPort < plugin.DefaultRemotePluginMinPort {
			b.minPluginPort = b.maxPluginPort
		} else {
			b.minPluginPort = plugin.DefaultRemotePluginMinPort
		}
	}

	if b.maxPluginPort == 0 {
		if b.minPluginPort != 0 && b.minPluginPort > plugin.DefaultRemotePluginMaxPort {
			b.maxPluginPort = b.minPluginPort
		} else {
			b.maxPluginPort = plugin.DefaultRemotePluginMaxPort
		}
	}

	authMethods := make([]ssh.AuthMethod, 0, 2)
	if b.publicKeyAuth {
		var signer ssh.Signer
		var err error
		if b.privateKeyPass != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(b.privateKey, []byte(b.privateKeyPass))
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key with passphrase: %w", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(b.privateKey)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if b.passwordAuth {
		authMethods = append(authMethods, ssh.Password(b.password))
	}

	hostKeyCallback := ssh.InsecureIgnoreHostKey() // Default to insecure host key checking
	var err error
	if b.useKnownHostsFile {
		if b.addUnknownHostsToFile {
			hostKeyCallback, err = newHostKeyAddingCallback(b.knownHostsPath)
		} else {
			hostKeyCallback, err = knownhosts.New(b.knownHostsPath)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create host key adding callback: %w", err)
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:            b.user,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         b.connectionTimeout,
	}

	return &sshTransport{
		host:          b.host,
		port:          b.port,
		config:        clientConfig,
		minPluginPort: b.minPluginPort,
		maxPluginPort: b.maxPluginPort,
		tempPath:      b.tempPath,
	}, nil
}

// NewSSHBuilder creates a new SSHTransportBuilder with default settings.
func NewSSHBuilder() (*SSHTransportBuilder, error) {
	knownHostsPath, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get default known hosts path: %w", err)
	}

	return &SSHTransportBuilder{
		port:              22,               // Default SSH port
		connectionTimeout: 10 * time.Second, // Default connection timeout
		knownHostsPath:    knownHostsPath,
	}, nil
}
