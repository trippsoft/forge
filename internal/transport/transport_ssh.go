package transport

import (
	"bytes"
	"context"
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
)

func DefaultKnownHostsPath() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, ".ssh", "known_hosts"), nil
}

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
}

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

func (b *SSHTransportBuilder) Host(host string) *SSHTransportBuilder {
	b.host = host
	return b
}

func (b *SSHTransportBuilder) Port(port uint16) *SSHTransportBuilder {
	b.port = port
	return b
}

func (b *SSHTransportBuilder) User(user string) *SSHTransportBuilder {

	b.user = user
	return b
}

func (b *SSHTransportBuilder) NoPublicKeyAuth() *SSHTransportBuilder {

	b.publicKeyAuth = false
	b.privateKey = nil
	b.privateKeyPass = ""
	return b
}

func (b *SSHTransportBuilder) PublicKeyAuth(privateKey []byte) *SSHTransportBuilder {

	b.publicKeyAuth = true
	b.privateKey = privateKey
	b.privateKeyPass = ""
	return b
}

func (b *SSHTransportBuilder) PublicKeyAuthWithPass(privateKey []byte, privateKeyPass string) *SSHTransportBuilder {

	b.publicKeyAuth = true
	b.privateKey = privateKey
	b.privateKeyPass = privateKeyPass
	return b
}

func (b *SSHTransportBuilder) PasswordAuth(password string) *SSHTransportBuilder {

	b.passwordAuth = true
	b.password = password
	return b
}

func (b *SSHTransportBuilder) DontUseKnownHosts() *SSHTransportBuilder {
	b.useKnownHostsFile = false
	return b
}

func (b *SSHTransportBuilder) UseKnownHosts(knownHostsPath string, addUnknownHosts bool) *SSHTransportBuilder {

	b.useKnownHostsFile = true
	b.knownHostsPath = knownHostsPath
	b.addUnknownHostsToFile = addUnknownHosts
	return b
}

func (b *SSHTransportBuilder) ConnectionTimeout(timeout time.Duration) *SSHTransportBuilder {

	b.connectionTimeout = timeout
	return b
}

func (b *SSHTransportBuilder) Build() (Transport, error) {

	if b.host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if b.port == 0 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}

	if b.user == "" {
		return nil, fmt.Errorf("user cannot be empty")
	}

	if b.publicKeyAuth && b.privateKey == nil {
		return nil, fmt.Errorf("privateKey cannot be empty when public key authentication is enabled")
	}

	if b.passwordAuth && b.password == "" {
		return nil, fmt.Errorf("password cannot be empty when password authentication is enabled")
	}

	if b.useKnownHostsFile && b.knownHostsPath == "" {
		return nil, fmt.Errorf("knownHostsPath cannot be empty when using known hosts")
	}

	if b.connectionTimeout <= 0 {
		return nil, fmt.Errorf("connectionTimeout must be greater than zero")
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
		host:   b.host,
		port:   b.port,
		config: clientConfig,
	}, nil
}

type sshResult struct {
	stdout string
	stderr string
	err    error
}

type sshTransport struct {
	host string
	port uint16

	config *ssh.ClientConfig

	client *ssh.Client

	fileSystem FileSystem

	hasValidatedPowerShell bool
	canRunPowerShell       bool
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

	session, err := s.client.NewSession()
	if err != nil {
		s.client.Close()
		s.client = nil
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	session.Close() // Close the test session immediately

	return nil
}

// Close implements Transport.
func (s *sshTransport) Close() error {

	if s.fileSystem != nil {
		_ = s.fileSystem.Close() // Close the file system if it exists
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

// ExecuteCommand implements Transport.
func (s *sshTransport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {

	err := s.Connect() // Ensure we are connected
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	outputChannel := make(chan *sshResult)

	go func() {
		var outBuf, errBuf bytes.Buffer
		session.Stdout = &outBuf
		session.Stderr = &errBuf

		err := session.Run(command)
		outputChannel <- &sshResult{
			stdout: outBuf.String(),
			stderr: errBuf.String(),
			err:    err,
		}
	}()

	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		return "", "", ctx.Err()
	case result := <-outputChannel:
		return result.stdout, result.stderr, result.err
	}
}

// ExecutePowerShell implements Transport.
func (s *sshTransport) ExecutePowerShell(ctx context.Context, command string) (string, error) {

	err := s.Connect() // Ensure we are connected
	if err != nil {
		return "", fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	if !s.hasValidatedPowerShell {
		// Check if PowerShell is available on the remote system
		powershellCheckCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command \"Write-Host 'PowerShell is available'\""

		session, err := s.client.NewSession()
		if err != nil {
			return "", fmt.Errorf("failed to create SSH session: %w", err)
		}

		err = session.Run(powershellCheckCmd)
		session.Close()
		if err != nil {
			s.hasValidatedPowerShell = true
			s.canRunPowerShell = false
			return "", fmt.Errorf("PowerShell is not available on the remote system: %w", err)
		} else {
			s.hasValidatedPowerShell = true
			s.canRunPowerShell = true
		}
	}

	if !s.canRunPowerShell {
		return "", fmt.Errorf("PowerShell is not available on the remote system")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	encodedCommand, err := encodePowerShellAsUTF16LEBase64(command)
	if err != nil {
		return "", fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	commandBuilder := &strings.Builder{}

	commandBuilder.WriteString("powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand ")
	commandBuilder.WriteString(encodedCommand)

	command = commandBuilder.String()

	outputChannel := make(chan *sshResult)

	go func() {
		var outBuf bytes.Buffer
		session.Stdout = &outBuf

		err := session.Run(command)
		outputChannel <- &sshResult{
			stdout: outBuf.String(),
			stderr: "",
			err:    err,
		}
	}()

	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		return "", ctx.Err()
	case result := <-outputChannel:
		return result.stdout, result.err
	}
}

// FileSystem implements Transport.
func (s *sshTransport) FileSystem() FileSystem {
	if s.fileSystem == nil {
		s.fileSystem = newSFTPFileSystem(s)
	}

	return s.fileSystem
}

type sftpFileSystem struct {
	transport *sshTransport
	client    *sftp.Client
}

func newSFTPFileSystem(transport *sshTransport) FileSystem {
	return &sftpFileSystem{transport: transport}
}

// IsNull implements FileSystem.
func (s *sftpFileSystem) IsNull() bool {
	return false // SFTP file system is always available
}

// Connect implements FileSystem.
func (s *sftpFileSystem) Connect() error {

	if s.client != nil {
		return nil // Already connected
	}

	err := s.transport.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	client, err := sftp.NewClient(s.transport.client)
	if err != nil {
		s.transport.Close() // Close the SSH client on error
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	s.client = client

	return nil
}

// Close implements FileSystem.
func (s *sftpFileSystem) Close() error {

	if s.client == nil {
		return nil // No client to close
	}

	err := s.client.Close()
	s.client = nil
	if err != nil {
		return fmt.Errorf("failed to close SFTP client: %w", err)
	}

	return nil
}

// Stat implements FileSystem.
func (s *sftpFileSystem) Stat(path string) (os.FileInfo, error) {

	err := s.Connect() // Ensure we are connected
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	return s.client.Stat(path)
}

// Open implements FileSystem.
func (s *sftpFileSystem) Open(path string) (File, error) {

	err := s.Connect() // Ensure we are connected
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	return s.client.Open(path)
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
