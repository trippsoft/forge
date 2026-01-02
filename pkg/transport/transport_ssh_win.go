// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/util"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sshWindowsPlatform struct {
	t *sshTransport

	arch string
}

// OS implements sshPlatform.
func (s *sshWindowsPlatform) OS() string {
	return "windows"
}

// Arch implements sshPlatform.
func (s *sshWindowsPlatform) Arch() string {
	return s.arch
}

// PathSeparator implements sshPlatform.
func (s *sshWindowsPlatform) PathSeparator() string {
	return `\`
}

// PluginExtension implements sshPlatform.
func (s *sshWindowsPlatform) PluginExtension() string {
	return ".exe"
}

// GetDefaultTempPath implements sshPlatform.
func (s *sshWindowsPlatform) GetDefaultTempPath() (string, error) {
	session, err := s.t.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}

	homeCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"Write-Host $env:USERPROFILE"`
	homeOutput, err := session.CombinedOutput(homeCmd)
	session.Close()
	if err != nil {
		return "", fmt.Errorf("failed to execute home directory detection command: %w", err)
	}

	homeDir := strings.TrimSpace(string(homeOutput))
	return fmt.Sprintf(`%s\AppData\Local\Temp\Forge`, homeDir), nil
}

// PopulateInfo implements sshPlatform.
func (s *sshWindowsPlatform) PopulateInfo() error {
	err := s.populateWindowsArch()
	if err != nil {
		return err
	}

	return nil
}

// MkdirAll implements sshPlatform.
func (s *sshWindowsPlatform) MkdirAll(path string) error {
	if s.t.sftpClient == nil {
		sftpClient, err := sftp.NewClient(s.t.client)
		if err != nil {
			return fmt.Errorf("failed to create SFTP client: %w", err)
		}
		s.t.sftpClient = sftpClient
	}

	dir, err := s.t.sftpClient.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}

		return &os.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	i := len(path)
	for i > 0 && path[i-1] != '\\' {
		i--
	}

	j := i
	for j > 0 && path[j-1] != '\\' {
		j--
	}

	if j > 3 { // Drive letter and colon, e.g., "C:\"
		err = s.MkdirAll(path[0 : j-1])
		if err != nil {
			return err
		}
	}

	err = s.t.sftpClient.Mkdir(path)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// UploadFile implements sshPlatform.
func (s *sshWindowsPlatform) UploadFile(localPath, remotePath string) error {
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
func (s *sshWindowsPlatform) StartPlugin(
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

	remotePluginPath := fmt.Sprintf(`%s\%s-%s.exe`, s.t.tempPath, namespace, pluginName)

	err = s.MkdirAll(s.t.tempPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.t.tempPath, err)
	}

	err = s.UploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	if escalation != nil {
		return s.startEscalatedPlugin(ctx, remotePluginPath, escalation)
	}

	session, err := s.t.client.NewSession()
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

	cmdlet := fmt.Sprintf(
		`$env:FORGE_PLUGIN_MIN_PORT=%d; $env:FORGE_PLUGIN_MAX_PORT=%d; & "%s"`,
		s.t.minPluginPort,
		s.t.maxPluginPort,
		remotePluginPath,
	)

	encodedCmdlet, err := encodePowerShellAsUTF16LEBase64(cmdlet)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
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

	listener, localPort, err := util.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
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

func (s *sshWindowsPlatform) startEscalatedPlugin(
	ctx context.Context,
	remotePluginPath string,
	escalation *Escalation,
) (*grpc.ClientConn, func(), error) {

	user := escalation.User()
	if user == "" || user == "SYSTEM" || user == `NT AUTHORITY\SYSTEM` {
		return s.startPluginAsSystem(remotePluginPath)
	}

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
		bufferReader := bufio.NewReader(teeReader)
		for {
			line, err := bufferReader.ReadString(':')
			if err != nil {
				return
			}

			if strings.Contains(line, forgeGSudoPrompt) {
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

	cmdlet := fmt.Sprintf(
		`$env:FORGE_PLUGIN_MIN_PORT=%d; $env:FORGE_PLUGIN_MAX_PORT=%d; gsudo -u %s "%s"`,
		s.t.minPluginPort,
		s.t.maxPluginPort,
		user,
		remotePluginPath,
	)

	encodedCmdlet, err := encodePowerShellAsUTF16LEBase64(cmdlet)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
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

		listener, localPort, err := util.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
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

func (s *sshWindowsPlatform) startPluginAsSystem(remotePluginPath string) (*grpc.ClientConn, func(), error) {

	session, err := s.t.client.NewSession()
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

	cmdlet := fmt.Sprintf(
		`$env:FORGE_PLUGIN_MIN_PORT=%d; $env:FORGE_PLUGIN_MAX_PORT=%d; gsudo -s "%s"`,
		s.t.minPluginPort,
		s.t.maxPluginPort,
		remotePluginPath,
	)

	encodedCmdlet, err := encodePowerShellAsUTF16LEBase64(cmdlet)
	if err != nil {
		session.Close()
		return nil, nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
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

	listener, localPort, err := util.GetListenerAndPortInRange(plugin.LocalPluginMinPort, plugin.LocalPluginMaxPort)
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

// FormatCommand implements sshPlatform.
func (s *sshWindowsPlatform) FormatCommand(cmd string, env ...string) string {
	sb := &strings.Builder{}
	sb.WriteString(`powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command "`)

	if len(env) == 0 {
		sb.WriteString("& '")
		sb.WriteString(cmd)
		sb.WriteString(`'"`)
		return sb.String()
	}

	for _, e := range env {
		sb.WriteString("$env.")
		sb.WriteString(e)
		sb.WriteString("; ")
	}

	sb.WriteString("& '")
	sb.WriteString(cmd)
	sb.WriteString(`'"`)
	return sb.String()
}

func (s *sshWindowsPlatform) populateWindowsArch() error {
	session, err := s.t.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	var outBuf, errBuf bytes.Buffer

	session.Stdout = &outBuf
	session.Stderr = &errBuf

	archCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"Write-Host $env:PROCESSOR_ARCHITECTURE"`
	err = session.Run(archCmd)
	if err != nil {
		stderr := strings.TrimSpace(errBuf.String())
		return fmt.Errorf("failed to execute architecture detection command: %w - %s", err, stderr)
	}

	stdout := strings.TrimSpace(outBuf.String())
	arch := strings.TrimSpace(strings.ToLower(stdout))

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

func newSSHWindowsPlatform(t *sshTransport) sshPlatform {
	return &sshWindowsPlatform{
		t: t,
	}
}
