// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

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

	sshSudoPrompt = "forge_sudo_prompt"
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

func (b *SSHTransportBuilder) UseKnownHosts(knownHostsPath string) *SSHTransportBuilder {

	b.useKnownHostsFile = true
	b.knownHostsPath = knownHostsPath
	b.addUnknownHostsToFile = true
	return b
}

func (b *SSHTransportBuilder) UseStrictKnownHosts(knownHostsPath string) *SSHTransportBuilder {

	b.useKnownHostsFile = true
	b.knownHostsPath = knownHostsPath
	b.addUnknownHostsToFile = false
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

type sshCmd struct {
	transport *sshTransport
	session   *ssh.Session
	ctx       context.Context
	completed bool

	command string
}

// OutputWithError implements Cmd.
func (s *sshCmd) OutputWithError(ctx context.Context) (string, string, error) {

	err := s.createSession(ctx)
	if err != nil {
		return "", "", err
	}
	defer s.session.Close()

	var outBuf, errBuf bytes.Buffer

	s.session.Stdout = &outBuf
	s.session.Stderr = &errBuf

	errChannel := make(chan error)

	go func() {
		err := s.session.Run(s.command)
		errChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return "", "", s.ctx.Err()
	case err = <-errChannel:
		s.session = nil
		s.completed = true
		stdout := strings.TrimSpace(outBuf.String())
		stderr := strings.TrimSpace(errBuf.String())
		return stdout, stderr, nil
	}
}

// Output implements Cmd.
func (s *sshCmd) Output(ctx context.Context) (string, error) {

	err := s.createSession(ctx)
	if err != nil {
		return "", err
	}
	defer s.session.Close()

	var outBuf bytes.Buffer

	s.session.Stdout = &outBuf

	errChannel := make(chan error)

	go func() {
		err := s.session.Run(s.command)
		errChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return "", s.ctx.Err()
	case err = <-errChannel:
		s.session = nil
		s.completed = true
		stdout := strings.TrimSpace(outBuf.String())
		return stdout, err
	}
}

// Run implements Cmd.
func (s *sshCmd) Run(ctx context.Context) error {

	err := s.createSession(ctx)
	if err != nil {
		return err
	}
	defer s.session.Close()

	errChannel := make(chan error)

	go func() {
		err := s.session.Run(s.command)
		errChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return s.ctx.Err()
	case err = <-errChannel:
		s.session = nil
		s.completed = true
		return err
	}
}

func (s *sshCmd) createSession(ctx context.Context) error {

	if s.completed {
		return errors.New("command already completed")
	}

	if s.session != nil {
		return errors.New("command already started")
	}

	err := s.transport.Connect()
	if err != nil {
		return err
	}

	s.session, err = s.transport.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}

	s.ctx = ctx

	return nil
}

// sshPlatformInfo provides platform-specific information for SSH connections.
type sshPlatformInfo interface {
	// canRunPowerShell indicates if PowerShell is available on the platform.
	canRunPowerShell() bool

	// pathSeparator returns the path separator for the platform.
	pathSeparator() rune
	// pathListSeparator returns the path list separator for the platform.
	pathListSeparator() rune
	// tempDir returns the temporary directory for the platform.
	tempDir() (string, error)
	// pathPrefixes returns the path prefixes for the platform.
	pathPrefixes() ([]string, error)

	// newCommand creates a new command with the specified escalation configuration.
	newCommand(command string, config Escalation) (Cmd, error)
}

type sshTransport struct {
	host string
	port uint16

	config     *ssh.ClientConfig
	client     *ssh.Client
	sftpClient *sftp.Client
	platform   sshPlatformInfo
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

	// Check if PowerShell is available on the remote system
	powershellCheckCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command \"Write-Host 'PowerShell is available'\""

	psErr := session.Run(powershellCheckCmd)
	session.Close()
	if psErr != nil {

		session, unameErr := s.client.NewSession()
		if unameErr != nil {
			return fmt.Errorf("failed to create SSH session for uname check: %w", unameErr)
		}
		defer session.Close()

		unameErr = session.Run("uname -s")
		if unameErr != nil {
			return fmt.Errorf("failed to check for PowerShell or uname command; PowerShell Error: %w; uname Error: %w", psErr, unameErr)
		}

		s.platform = &sshPosixInfo{transport: s} // For now, we will assume non-Windows is POSIX.
		return nil
	}

	s.platform = &sshWindowsInfo{transport: s}
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

// NewCommand implements Transport.
func (s *sshTransport) NewCommand(command string, escalateConfig Escalation) (Cmd, error) {

	err := s.Connect() // Connect to ensure that the platform detection is done
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return s.platform.newCommand(command, escalateConfig)
}

// NewPowerShellCommand implements Transport.
func (s *sshTransport) NewPowerShellCommand(command string, escalateConfig Escalation) (Cmd, error) {

	err := s.Connect() // Connect to ensure that the platform detection is done
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if !s.platform.canRunPowerShell() {
		return nil, errors.New("PowerShell is not available on the remote system")
	}

	encodedCommand, err := encodePowerShellAsUTF16LEBase64(command)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	command = fmt.Sprintf("powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s", encodedCommand)

	return s.platform.newCommand(command, escalateConfig)
}

// Stat implements Transport.
func (s *sshTransport) Stat(path string) (os.FileInfo, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	fileInfo, err := s.sftpClient.Stat(path)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil, nil // Return nil if the file does not exist
	}

	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// Create implements Transport.
func (s *sshTransport) Create(path string) (File, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	file, err := s.sftpClient.Create(path)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil, nil // Return nil if the file does not exist
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Open implements Transport.
func (s *sshTransport) Open(path string) (File, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	file, err := s.sftpClient.Open(path)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil, nil // Return nil if the file does not exist
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Mkdir implements Transport.
func (s *sshTransport) Mkdir(path string) error {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	err = s.sftpClient.Mkdir(path)
	if errors.Is(err, os.ErrExist) || errors.Is(err, syscall.EEXIST) {
		return nil // Directory already exists, return nil
	}

	return err
}

// MkdirAll implements Transport.
func (s *sshTransport) MkdirAll(path string) error {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	err = s.sftpClient.MkdirAll(path)
	if errors.Is(err, os.ErrExist) || errors.Is(err, syscall.EEXIST) {
		return nil // Directory already exists, return nil
	}

	return err
}

// Remove implements Transport.
func (s *sshTransport) Remove(path string) error {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	return s.sftpClient.Remove(path)
}

// RemoveAll implements Transport.
func (s *sshTransport) RemoveAll(path string) error {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	return s.sftpClient.RemoveAll(path)
}

// Join implements Transport.
func (s *sshTransport) Join(elem ...string) string {

	if len(elem) == 0 {
		return ""
	}

	stringBuilder := &strings.Builder{}

	for i, e := range elem {
		if i > 0 {
			stringBuilder.WriteRune(s.platform.pathSeparator())
		}

		if strings.HasSuffix(e, string(s.platform.pathSeparator())) {
			e = strings.TrimSuffix(e, string(s.platform.pathSeparator()))
		}

		stringBuilder.WriteString(e)
	}

	return stringBuilder.String()
}

// TempDir implements Transport.
func (s *sshTransport) TempDir() (string, error) {
	return s.platform.tempDir()
}

// CreateTemp implements Transport.
func (s *sshTransport) CreateTemp(dir, pattern string) (File, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	if dir == "" {
		dir, err = s.TempDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get temp dir: %w", err)
		}
	}

	splitPattern := strings.Split(pattern, "*")
	if len(splitPattern) > 2 {
		return nil, fmt.Errorf("pattern must contain at most one wildcard (*)")
	}

	var prefix, suffix string
	if len(splitPattern) == 1 {
		prefix = splitPattern[0]
	} else {
		prefix = splitPattern[0]
		suffix = splitPattern[1]
	}

	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(dir)
	stringBuilder.WriteRune(s.platform.pathSeparator())
	stringBuilder.WriteString(prefix)

	randomNumber := fmt.Sprintf("%d", time.Now().UnixNano()%1000000) // Simple random number based on current time
	stringBuilder.WriteString(randomNumber)
	stringBuilder.WriteString(suffix)

	return s.Create(stringBuilder.String())
}

// MkdirTemp implements Transport.
func (s *sshTransport) MkdirTemp(dir, pattern string) (string, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return "", fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	if dir == "" {
		dir, err = s.TempDir()
		if err != nil {
			return "", fmt.Errorf("failed to get temp dir: %w", err)
		}
	}

	splitPattern := strings.Split(pattern, "*")
	if len(splitPattern) > 2 {
		return "", fmt.Errorf("pattern must contain at most one wildcard (*)")
	}

	var prefix, suffix string
	if len(splitPattern) == 1 {
		prefix = splitPattern[0]
	} else {
		prefix = splitPattern[0]
		suffix = splitPattern[1]
	}

	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(dir)
	stringBuilder.WriteRune(s.platform.pathSeparator())
	stringBuilder.WriteString(prefix)

	randomNumber := fmt.Sprintf("%d", time.Now().UnixNano()%1000000) // Simple random number based on current time
	stringBuilder.WriteString(randomNumber)
	stringBuilder.WriteString(suffix)

	tempDirPath := stringBuilder.String()

	err = s.Mkdir(tempDirPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return tempDirPath, nil
}

// Symlink implements Transport.
func (s *sshTransport) Symlink(target, path string) error {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	return s.sftpClient.Symlink(target, path)
}

// ReadLink implements Transport.
func (s *sshTransport) ReadLink(path string) (string, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return "", fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	target, err := s.sftpClient.ReadLink(path)
	if err != nil {
		return "", err
	}

	if s.platform.pathSeparator() != '/' {
		target = strings.ReplaceAll(target, "/", string(s.platform.pathSeparator())) // Normalize path for Windows
		target = strings.Trim(target, string(s.platform.pathSeparator()))
	}

	return target, nil
}

// RealPath implements Transport.
func (s *sshTransport) RealPath(path string) (string, error) {

	err := s.connectSFTP() // Ensure we are connected to SFTP
	if err != nil {
		return "", fmt.Errorf("failed to connect to SFTP client: %w", err)
	}

	prefixes, err := s.platform.pathPrefixes()
	if err != nil {
		return "", fmt.Errorf("failed to get path prefixes: %w", err)
	}

	for _, prefix := range prefixes {
		newPath := prefix + path
		fileInfo, _ := s.Stat(newPath) // Ignore error, just check if file exists

		if fileInfo != nil {
			return newPath, nil // Return the first valid path found
		}
	}

	realPath, err := s.sftpClient.RealPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to get real path: %w", err)
	}

	fileInfo, _ := s.Stat(realPath) // Ignore error, just check if file exists
	if fileInfo != nil {
		return realPath, nil // Return the absolute path if it exists
	}

	return "", os.ErrNotExist // Return error if no valid path found
}

func (s *sshTransport) connectSFTP() error {

	if s.sftpClient != nil {
		return nil // Already connected
	}

	if s.client == nil {
		err := s.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect to SSH transport: %w", err)
		}
	}

	client, err := sftp.NewClient(s.client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	s.sftpClient = client

	return nil
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
