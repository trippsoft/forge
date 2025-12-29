// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/network"
	"github.com/trippsoft/forge/pkg/plugin"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sshPosixPlatform struct {
	t *sshTransport

	os   string
	arch string
}

// OS implements sshPlatform.
func (s *sshPosixPlatform) OS() string {
	return s.os
}

// Arch implements sshPlatform.
func (s *sshPosixPlatform) Arch() string {
	return s.arch
}

// PathSeparator implements sshPlatform.
func (s *sshPosixPlatform) PathSeparator() string {
	return "/"
}

// PluginExtension implements sshPlatform.
func (s *sshPosixPlatform) PluginExtension() string {
	return ""
}

// GetDefaultTempPath implements sshPlatform.
func (s *sshPosixPlatform) GetDefaultTempPath() (string, error) {
	session, err := s.t.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	homeCmd := "echo $HOME"
	homeOutput, err := session.CombinedOutput(homeCmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute home directory detection command: %w", err)
	}

	homeDir := strings.TrimSpace(string(homeOutput))
	return fmt.Sprintf("%s/.local/share/forge-tmp", homeDir), nil
}

// PopulateInfo implements sshPlatform.
func (s *sshPosixPlatform) PopulateInfo() error {
	err := s.populatePosixOS()
	if err != nil {
		return err
	}

	err = s.populatePosixArch()
	if err != nil {
		return err
	}

	return nil
}

// MkdirAll implements sshPlatform.
func (s *sshPosixPlatform) MkdirAll(path string) error {
	if s.t.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.t.client)
		if err != nil {
			return fmt.Errorf("failed to create SFTP client: %w", err)
		}
		s.t.sftpClient = sftpClient
	}

	err := s.t.sftpClient.MkdirAll(path)
	if err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", path, err)
	}

	return nil
}

// UploadFile implements sshPlatform.
func (s *sshPosixPlatform) UploadFile(localPath, remotePath string) error {
	if s.t.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.t.client)
		if err != nil {
			return fmt.Errorf("failed to create SFTP client: %w", err)
		}
		s.t.sftpClient = sftpClient
	}

	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s': %w", localPath, err)
	}
	defer localFile.Close()

	remoteFile, err := s.t.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file '%s': %w", remotePath, err)
	}
	defer remoteFile.Sync()
	defer remoteFile.Close()

	_, err = remoteFile.ReadFrom(localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file to remote path '%s': %w", remotePath, err)
	}

	return nil
}

// StartPlugin implements sshPlatform.
func (s *sshPosixPlatform) StartPlugin(
	ctx context.Context,
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {

	err := s.t.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	localPluginPath, err := plugin.FindPluginPath(basePath, namespace, pluginName, s.OS(), s.Arch())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find local plugin path: %w", err)
	}

	remotePluginPath := fmt.Sprintf("%s/%s-%s", s.t.tempPath, namespace, pluginName)

	err = s.MkdirAll(s.t.tempPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.t.tempPath, err)
	}

	err = s.UploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	err = session.Run(fmt.Sprintf("chmod +x %s", remotePluginPath))
	session.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set execute permission on remote plugin '%s': %w", remotePluginPath, err)
	}

	if escalation != nil {
		return s.startEscalatedPlugin(ctx, remotePluginPath, escalation)
	}

	session, err = s.t.client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrReader, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	cmd := fmt.Sprintf(
		"FORGE_PLUGIN_MIN_PORT=%d FORGE_PLUGIN_MAX_PORT=%d %s",
		s.t.minPluginPort,
		s.t.maxPluginPort,
		remotePluginPath,
	)

	err = session.Start(cmd)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to start remote plugin '%s': %w", remotePluginPath, err)
	}

	scanner := bufio.NewScanner(stdoutReader)
	var portOutput string
	for scanner.Scan() {
		portOutput = scanner.Text()
		break
	}

	if portOutput == "" {
		errOutput, _ := io.ReadAll(stderrReader)
		session.Close()
		stderr := strings.TrimSpace(string(errOutput))
		return nil, nil, fmt.Errorf("no port output from plugin: %s", stderr)
	}

	remotePort, err := strconv.ParseUint(portOutput, 10, 16)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("invalid port output from plugin: %w", err)
	}

	listener, localPort, err := network.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to get local listener: %w", err)
	}

	go s.t.forwardConnections(listener, uint16(remotePort))

	cleanup := func() {
		listener.Close()
		session.Signal(ssh.SIGTERM)
		session.Close()
		session.Wait()
		s.t.sftpClient.Remove(remotePluginPath)
	}

	address := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", localPort))
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failed to create gRPC client connection: %w", err)
	}

	return connection, cleanup, nil
}

func (s *sshPosixPlatform) startEscalatedPlugin(
	ctx context.Context,
	remotePluginPath string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {

	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	var errBuf strings.Builder

	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrReader, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdinWriter, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	teeReader := io.TeeReader(stderrReader, &errBuf)

	outputChannel := make(chan string)

	go func() {
		scanner := bufio.NewScanner(teeReader)
		promptsAnswered := 0

		for scanner.Scan() {
			if promptsAnswered > 3 {
				session.Signal(ssh.SIGKILL)
				session.Close()
				return
			}

			line := scanner.Text()
			if strings.Contains(line, forgeSudoPrompt) {
				promptsAnswered++
				_, err = stdinWriter.Write([]byte(escalation.Pass() + "\n"))
				if err != nil {
					session.Signal(ssh.SIGKILL)
					session.Close()
					return
				}
			}
		}
	}()

	go func() {
		defer close(outputChannel)
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			line := scanner.Text()
			outputChannel <- line
			return
		}
	}()

	user := escalation.User()
	if user == "" {
		user = "root"
	}

	cmd := fmt.Sprintf(
		"sudo -S -p '%s:' -u %s FORGE_PLUGIN_MIN_PORT=%d FORGE_PLUGIN_MAX_PORT=%d %s",
		forgeSudoPrompt,
		user,
		s.t.minPluginPort,
		s.t.maxPluginPort,
		remotePluginPath,
	)

	err = session.Start(cmd)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to start remote plugin '%s': %w", remotePluginPath, err)
	}

	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("context cancelled while starting plugin at '%s': %w", remotePluginPath, ctx.Err())
	case portOutput := <-outputChannel:

		if portOutput == "" {
			errOutput, _ := io.ReadAll(stderrReader)
			session.Close()
			stderr := strings.TrimSpace(string(errOutput))
			return nil, nil, fmt.Errorf("no port output from plugin: %s", stderr)
		}

		remotePort, err := strconv.ParseUint(portOutput, 10, 16)
		if err != nil {
			session.Close()
			return nil, nil, fmt.Errorf("invalid port output from plugin: %w", err)
		}

		listener, localPort, err := network.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
		if err != nil {
			session.Close()
			return nil, nil, fmt.Errorf("failed to get local listener: %w", err)
		}

		go s.t.forwardConnections(listener, uint16(remotePort))

		cleanup := func() {
			listener.Close()
			session.Signal(ssh.SIGTERM)
			session.Close()
			session.Wait()
			s.t.sftpClient.Remove(remotePluginPath)
		}

		address := net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", localPort))
		connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to create gRPC client connection: %w", err)
		}

		return connection, cleanup, nil
	}
}

// FormatCommand implements sshPlatform.
func (s *sshPosixPlatform) FormatCommand(cmd string, env ...string) string {
	if len(env) == 0 {
		return cmd
	}

	sb := &strings.Builder{}

	for _, e := range env {
		sb.WriteString(e)
		sb.WriteRune(' ')
	}

	sb.WriteString(cmd)
	return sb.String()
}

func (s *sshPosixPlatform) populatePosixOS() error {
	session, err := s.t.client.NewSession()
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

func (s *sshPosixPlatform) populatePosixArch() error {
	session, err := s.t.client.NewSession()
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

func newSSHPosixPlatform(t *sshTransport) sshPlatform {
	return &sshPosixPlatform{
		t: t,
	}
}
