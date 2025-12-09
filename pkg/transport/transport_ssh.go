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
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/network"
	"github.com/trippsoft/forge/pkg/plugin"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type sshPlatform interface {
	OS() string
	Arch() string
	PathSeparator() string
	PluginExtension() string

	GetDefaultTempPath() (string, error)

	PopulateInfo() error

	MkdirAll(path string) error
	UploadFile(localPath, remotePath string) error
	FormatCommand(cmd string, env ...string) string
}

type sshTransport struct {
	host string
	port uint16

	minPluginPort uint16
	maxPluginPort uint16

	platform sshPlatform

	tempPath string

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
	if s.platform == nil || s.platform.OS() == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.platform.OS(), nil
}

// Arch implements Transport.
func (s *sshTransport) Arch() (string, error) {
	if s.platform == nil || s.platform.Arch() == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.platform.Arch(), nil
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

	if s.platform != nil && s.platform.OS() != "" && s.platform.Arch() != "" && s.tempPath != "" {
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
		s.platform = newSSHWindowsPlatform(s)
	} else {
		s.platform = newSSHPosixPlatform(s)
	}

	err = s.platform.PopulateInfo()
	if err != nil {
		return fmt.Errorf("failed to populate platform info: %w", err)
	}

	if s.tempPath == "" {
		s.tempPath, err = s.platform.GetDefaultTempPath()
		if err != nil {
			return fmt.Errorf("failed to get default temp path: %w", err)
		}
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

// StartPlugin implements Transport.
func (s *sshTransport) StartPlugin(
	namespace string,
	pluginName string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {
	// TODO - handle escalation if needed

	err := s.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	localPluginPath, err := plugin.FindPluginPath(namespace, pluginName, s.platform.OS(), s.platform.Arch())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find local plugin path: %w", err)
	}

	remotePluginPath := fmt.Sprintf(
		"%s%s%s-%s%s",
		s.tempPath,
		s.platform.PathSeparator(),
		namespace,
		pluginName,
		s.platform.PluginExtension(),
	)

	err = s.platform.MkdirAll(s.tempPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.tempPath, err)
	}

	err = s.platform.UploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	envMinPort := fmt.Sprintf("FORGE_PLUGIN_MIN_PORT=%d", s.minPluginPort)
	envMaxPort := fmt.Sprintf("FORGE_PLUGIN_MAX_PORT=%d", s.maxPluginPort)

	if s.platform.OS() != "windows" {
		session, err := s.client.NewSession()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
		}

		err = session.Run(fmt.Sprintf("chmod +x %s", remotePluginPath))
		session.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to make discovery plugin executable: %w", err)
		}
	}

	cmd := s.platform.FormatCommand(remotePluginPath, envMinPort, envMaxPort)

	session, err := s.client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the plugin
	err = session.Start(cmd)

	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to start plugin: %w", err)
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
		return nil, nil, fmt.Errorf("no port output from plugin: %s", errBuf.String())
	}

	remotePort, err := strconv.ParseUint(portOutput, 10, 16)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("invalid port output: %w", err)
	}

	listener, localPort, err := network.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to get local listener: %w", err)
	}

	go s.forwardConnections(listener, uint16(remotePort))

	cleanup := func() {
		listener.Close()
		session.Signal(ssh.SIGTERM)
		session.Close()
		s.sftpClient.Remove(remotePluginPath)
	}

	address := fmt.Sprintf("127.0.0.1:%d", localPort)
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to create gRPC client connection: %w", err)
	}

	return connection, cleanup, nil
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
