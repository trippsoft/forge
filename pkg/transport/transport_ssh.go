// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/discover"
	"github.com/trippsoft/forge/pkg/network"
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

	minLocalPluginPort uint16
	maxLocalPluginPort uint16

	minRemotePluginPort uint16
	maxRemotePluginPort uint16

	os   string
	arch string

	discoveryPluginBasePath string
	tempPath                string

	config     *ssh.ClientConfig
	client     *ssh.Client
	sftpClient *sftp.Client
}

// Type implements Transport.
func (s *sshTransport) Type() TransportType {
	return TransportTypeSSH
}

// OS implements Transport.
func (s *sshTransport) OS() (string, error) {
	if s.os == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.os, nil
}

// Arch implements Transport.
func (s *sshTransport) Arch() (string, error) {
	if s.arch == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.arch, nil
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

// StartDiscovery implements Transport.
func (s *sshTransport) StartDiscovery() (*discover.DiscoveryClient, error) {
	err := s.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	var extension string
	if s.os == "windows" {
		extension = ".exe"
	}

	localPluginPath := fmt.Sprintf(
		"%sforge-discover_%s_%s%s",
		s.discoveryPluginBasePath,
		s.os,
		s.arch,
		extension,
	)

	remotePluginPath := fmt.Sprintf("%s/forge-discover%s", s.tempPath, extension)

	if s.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.client)
		if err != nil {
			return nil, fmt.Errorf("failed to create SFTP client: %w", err)
		}
		s.sftpClient = sftpClient
	}

	err = s.sftpClient.MkdirAll(s.tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.tempPath, err)
	}

	err = s.uploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	if s.os != "windows" {
		session, err := s.client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create SSH session: %w", err)
		}

		err = session.Run(fmt.Sprintf("chmod +x %s", remotePluginPath))
		session.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to make discovery plugin executable: %w", err)
		}
	}

	session, err := s.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the plugin
	err = session.Start(
		fmt.Sprintf(
			"FORGE_DISCOVERY_MIN_PORT=%d FORGE_DISCOVERY_MAX_PORT=%d %s",
			s.minRemotePluginPort,
			s.maxRemotePluginPort,
			remotePluginPath,
		),
	)

	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to start plugin: %w", err)
	}

	// Read port from stdout
	scanner := bufio.NewScanner(stdout)
	var portOutput string
	for scanner.Scan() {
		portOutput = scanner.Text()
		break
	}

	if portOutput == "" {
		errBuf := &bytes.Buffer{}
		errBuf.ReadFrom(stderr)
		session.Close()
		return nil, fmt.Errorf("no port output from plugin: %s", errBuf.String())
	}

	remotePort, err := strconv.ParseUint(portOutput, 10, 16)
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("invalid port output: %w", err)
	}

	// TODO - Implement configurable local port range
	listener, localPort, err := network.GetListenerAndPortInRange(s.minLocalPluginPort, s.maxLocalPluginPort)
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get local listener: %w", err)
	}

	go s.forwardConnections(listener, uint16(remotePort))

	cleanup := func() {
		listener.Close()
		session.Signal(ssh.SIGTERM)
		session.Close()
		s.sftpClient.Remove(remotePluginPath)
	}

	discoveryClient := discover.NewDiscoveryClient(localPort, cleanup)

	return discoveryClient, nil
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
	s.tempPath = fmt.Sprintf(`%s\AppData\Local\Temp\Forge`, homeDir)
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
	s.tempPath = fmt.Sprintf("%s/.local/share/forge-tmp", homeDir)
	return nil
}

func (s *sshTransport) uploadFile(localPath, remotePath string) error {
	if s.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.client)
		if err != nil {
			return fmt.Errorf("failed to create SFTP client: %w", err)
		}
		s.sftpClient = sftpClient
	}

	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s': %w", localPath, err)
	}
	defer localFile.Close()

	remoteFile, err := s.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file '%s': %w", remotePath, err)
	}
	defer remoteFile.Close()

	_, err = remoteFile.ReadFrom(localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file to '%s': %w", remotePath, err)
	}

	return nil
}

func (s *sshTransport) forwardConnections(localListener net.Listener, remotePort uint16) {
	for {
		localConn, err := localListener.Accept()
		if err != nil {
			return // Listener closed
		}

		go func(local net.Conn) {
			defer local.Close()

			remoteConn, err := s.client.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", remotePort))
			if err != nil {
				return
			}
			defer remoteConn.Close()

			// Bidirectional copy
			done := make(chan struct{}, 2)
			go func() {
				io.Copy(remoteConn, local)
				done <- struct{}{}
			}()
			go func() {
				io.Copy(local, remoteConn)
				done <- struct{}{}
			}()
			<-done
		}(localConn)
	}
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
