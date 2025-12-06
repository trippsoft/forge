// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	DefaultSSHPort               uint16        = 22
	DefaultUseKnownHostsFile     bool          = true
	DefaultAddUnknownHostsToFile bool          = true
	DefaultSSHConnectionTimeout  time.Duration = 10 * time.Second

	DefaultMinPluginPort uint16 = 25000
	DefaultMaxPluginPort uint16 = 40000

	sshSudoPrompt = "forge_sudo_prompt"
)

func DefaultKnownHostsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".ssh", "known_hosts"), nil
}

type sshTransport struct {
	host string
	port uint16

	minPluginPort uint16
	maxPluginPort uint16

	os   string
	arch string

	tempPath string

	config     *ssh.ClientConfig
	client     *ssh.Client
	sftpClient *sftp.Client
}

// Type implements Transport.
func (s *sshTransport) Type() TransportType {
	return TransportTypeSSH
}

// Connect implements Transport.
func (s *sshTransport) Connect() error {
	if s.client != nil {
		return nil // Already connected
	}

	address := fmt.Sprintf("%s:%d", s.host, s.port)
	client, err := ssh.Dial("tcp", address, s.config)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server at %s: %w", address, err)
	}

	s.client = client

	if s.os != "" && s.arch != "" && s.tempPath != "" {
		session, err := s.client.NewSession()
		if err != nil {
			s.client.Close()
			s.client = nil
			return fmt.Errorf("failed to create SSH session: %w", err)
		}
		session.Close()

		return nil
	}

	err = s.populatePlatformInfo()
	if err != nil {
		s.client.Close()
		s.client = nil
		return fmt.Errorf("failed to populate platform info: %w", err)
	}

	return nil
}

// Close implements Transport.
func (s *sshTransport) Close() error {
	if s.sftpClient != nil {
		_ = s.sftpClient.Close() // Close the SFTP client if it exists
	}

	if s.client == nil {
		return nil // No client to close
	}

	err := s.client.Close()
	s.client = nil
	if err != nil {
		return fmt.Errorf("failed to close SSH client: %w", err)
	}

	return nil
}

func (s *sshTransport) populatePlatformInfo() error {

	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}

	powershellCheckCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"Write-Host 'PowerShell is available'"`
	psErr := session.Run(powershellCheckCmd)
	session.Close()
	if psErr == nil {
		return s.populateWindowsInfo()
	}

	return s.populatePosixInfo()
}

func (s *sshTransport) populateWindowsInfo() error {
	s.os = "windows"
	err := s.populateWindowsArch()
	if err != nil {
		return err
	}

	if s.tempPath == "" {
		err = s.populateWindowsTempPath()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sshTransport) populateWindowsArch() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	archCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"$env:PROCESSOR_ARCHITECTURE"`
	archOutput, err := session.CombinedOutput(archCmd)
	if err != nil {
		return fmt.Errorf("failed to execute architecture detection command: %w", err)
	}

	arch := strings.TrimSpace(strings.ToLower(string(archOutput)))

	switch arch {
	case "amd64", "arm64":
		s.arch = arch
		return nil
	case "x86":
		s.arch = "386"
		return nil
	}

	return fmt.Errorf("unknown or unsupported architecture: %s", arch)
}

func (s *sshTransport) populateWindowsTempPath() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}

	homeCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"$env:USERPROFILE"`
	homeOutput, err := session.CombinedOutput(homeCmd)
	session.Close()
	if err != nil {
		return fmt.Errorf("failed to execute home directory detection command: %w", err)
	}

	homeDir := strings.TrimSpace(string(homeOutput))
	s.tempPath = fmt.Sprintf("%s\\AppData\\Local\\Temp", homeDir)
	return nil
}

func (s *sshTransport) populatePosixInfo() error {
	err := s.populatePosixOS()
	if err != nil {
		return err
	}

	err = s.populatePosixArch()
	if err != nil {
		return err
	}

	if s.tempPath == "" {
		err = s.populatePosixTempPath()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sshTransport) populatePosixOS() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	osCmd := "uname -s"
	osOutput, err := session.CombinedOutput(osCmd)
	if err != nil {
		return fmt.Errorf("failed to execute OS detection command: %w", err)
	}

	os := strings.TrimSpace(strings.ToLower(string(osOutput)))

	switch os {
	case "aix", "darwin", "dragonfly", "freebsd", "illumos", "linux", "netbsd", "openbsd", "plan9", "solaris":
		s.os = os
	default:
		return fmt.Errorf("unknown or unsupported POSIX OS: %s", os)
	}

	return nil
}

func (s *sshTransport) populatePosixArch() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	archCmd := "uname -m"
	archOutput, err := session.CombinedOutput(archCmd)
	if err != nil {
		return fmt.Errorf("failed to execute architecture detection command: %w", err)
	}

	arch := strings.TrimSpace(strings.ToLower(string(archOutput)))

	switch arch {
	case "x86_64":
		s.arch = "amd64"
	case "aarch64":
		s.arch = "arm64"
	case "i386", "i486", "i586", "i686", "i786", "x86":
		s.arch = "386"
	case "armv6l", "armv7l":
		s.arch = "arm"
	case "386", "amd64", "arm", "arm64", "mips", "mips64", "ppc64", "ppc64le", "riscv64", "s390x":
		s.arch = arch
	default:
		return fmt.Errorf("unknown or unsupported architecture: %s", arch)
	}

	return nil
}

func (s *sshTransport) populatePosixTempPath() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	homeCmd := "echo $HOME"
	homeOutput, err := session.CombinedOutput(homeCmd)
	if err != nil {
		return fmt.Errorf("failed to execute home directory detection command: %w", err)
	}

	homeDir := strings.TrimSpace(string(homeOutput))
	s.tempPath = fmt.Sprintf("%s/.local/share", homeDir)
	return nil
}

// GetOSAndArch implements Transport.
func (s *sshTransport) GetOSAndArch() (string, string, error) {
	err := s.Connect()
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	return s.os, s.arch, nil
}

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
		if b.maxPluginPort != 0 && b.maxPluginPort < DefaultMinPluginPort {
			b.minPluginPort = b.maxPluginPort
		} else {
			b.minPluginPort = DefaultMinPluginPort
		}
	}

	if b.maxPluginPort == 0 {
		if b.minPluginPort != 0 && b.minPluginPort > DefaultMaxPluginPort {
			b.maxPluginPort = b.minPluginPort
		} else {
			b.maxPluginPort = DefaultMaxPluginPort
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

func newHostKeyAddingCallback(path string) (ssh.HostKeyCallback, error) {
	_, err := os.Stat(path)
	if err != nil && (errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT)) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create known hosts file %s: %w", path, err)
		}
		file.Close()
	}

	checkingCallback, err := knownhosts.New(path)
	if err != nil {
		return nil, err
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err = checkingCallback(hostname, remote, key)
		if err == nil {
			return nil // Host key is already known
		}

		var keyErr *knownhosts.KeyError
		if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
			return keyErr // Host has known hosts entry, but key does not match
		}

		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open known hosts file %s: %w", path, err)
		}

		defer file.Close()

		remoteNormalized := knownhosts.Normalize(remote.String())
		hostNormalized := knownhosts.Normalize(hostname)
		addresses := []string{remoteNormalized}

		if remoteNormalized != hostNormalized {
			addresses = append(addresses, hostNormalized)
		}

		_, err = file.WriteString(knownhosts.Line(addresses, key) + "\n")

		return err
	}, nil
}
