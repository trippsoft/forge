// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/powershell"
)

type sshWindowsPlatform struct {
	t *sshTransport

	arch string
}

// OS implements [sshPlatform].
func (s *sshWindowsPlatform) OS() string {
	return "windows"
}

// Arch implements [sshPlatform].
func (s *sshWindowsPlatform) Arch() string {
	return s.arch
}

// PathSeparator implements [sshPlatform].
func (s *sshWindowsPlatform) PathSeparator() string {
	return `\`
}

// PluginExtension implements [sshPlatform].
func (s *sshWindowsPlatform) PluginExtension() string {
	return ".exe"
}

// GetDefaultTempPath implements [sshPlatform].
func (s *sshWindowsPlatform) GetDefaultTempPath() (string, error) {
	err := s.t.connectSFTP()
	if err != nil {
		return "", err
	}

	fileInfo, err := s.t.sftpClient.Stat(`C:\Windows\Temp`) // Try Windows temp folder first
	if err == nil && fileInfo.IsDir() {
		err = s.MkdirAll(`C:\Windows\Temp\Forge`)
		if err == nil {
			return `C:\Windows\Temp\Forge`, nil
		}
	}

	fileInfo, err = s.t.sftpClient.Stat(`C:\ProgramData`) // Try ProgramData next
	if err == nil && fileInfo.IsDir() {
		err = s.MkdirAll(`C:\ProgramData\Forge\tmp`)
		if err == nil {
			return `C:\ProgramData\Forge\tmp`, nil
		}
	}

	session, err := s.t.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	homeCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"Write-Host $env:USERPROFILE"`
	homeOutput, err := session.CombinedOutput(homeCmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute home directory detection command: %w", err)
	}

	homeDir := strings.TrimSpace(string(homeOutput))
	fileInfo, err = s.t.sftpClient.Stat(homeDir)
	if err == nil && fileInfo.IsDir() {
		tmpPath := fmt.Sprintf(`%s\AppData\Local\Temp\Forge`, homeDir)
		err = s.MkdirAll(tmpPath)
		if err == nil {
			return tmpPath, nil
		}
	}

	return "", errors.New("failed to determine suitable remote temp path")
}

// PopulateInfo implements [sshPlatform].
func (s *sshWindowsPlatform) PopulateInfo() error {
	err := s.populateWindowsArch()
	if err != nil {
		return err
	}

	return nil
}

// MkdirAll implements [sshPlatform].
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

// UploadFile implements [sshPlatform].
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

// StartPluginSession implements [sshPlatform].
func (s *sshWindowsPlatform) StartPluginSession(
	ctx context.Context,
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (plugin.Session, error) {
	err := s.t.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	localPluginPath, err := plugin.FindPluginPath(basePath, namespace, pluginName, s.OS(), s.Arch())
	if err != nil {
		return nil, fmt.Errorf("failed to find local plugin path: %w", err)
	}

	remotePluginPath := fmt.Sprintf(`%s\%s-%s.exe`, s.t.tempPath, namespace, pluginName)

	err = s.MkdirAll(s.t.tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.t.tempPath, err)
	}

	err = s.UploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	if escalation != nil {
		return s.startEscalatedPluginSession(ctx, remotePluginPath, escalation)
	}

	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	cmdlet := fmt.Sprintf(`& "%s"`, remotePluginPath)

	encodedCmdlet, err := powershell.EncodeAsUTF16LEBase64(cmdlet)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
	)

	stderrPipeReader, stderrPipeWriter := io.Pipe()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		defer stderrPipeWriter.Close()

		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)

		for {
			n, readErr := stderr.Read(buf)
			if n > 0 {
				accumulatedStderr.Write(buf[:n])
				text := accumulatedStderr.String()

				if strings.Contains(text, plugin.PluginReadyMessage) {
					close(readyChan)
					io.Copy(stderrPipeWriter, stderr)
					return
				}
			}

			if readErr != nil {
				if readErr != io.EOF {
					errChan <- fmt.Errorf("error reading stderr for plugin at '%s': %w", remotePluginPath, readErr)
					stdin.Close()
					session.Close()
					session.Wait()
				}
				return
			}
		}
	}()

	err = session.Start(cmd)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to start remote plugin '%s': %w", remotePluginPath, err)
	}

	select {
	case <-ctx.Done():
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, fmt.Errorf("context cancelled while starting plugin at '%s': %w", remotePluginPath, ctx.Err())
	case err := <-errChan:
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, err
	case <-readyChan:
		return &sshPluginSession{
			session: session,
			stdout:  stdout,
			stderr:  stderrPipeReader,
			stdin:   stdin,
		}, nil
	}
}

func (s *sshWindowsPlatform) startEscalatedPluginSession(
	ctx context.Context,
	path string,
	escalation *Escalation,
) (plugin.Session, error) {
	user := escalation.User()
	if user == "" || user == "SYSTEM" || user == `NT AUTHORITY\SYSTEM` {
		return s.startPluginSessionAsSystem(ctx, path)
	}

	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	cmdlet := fmt.Sprintf(`gsudo -u %s "%s"`, user, path)
	encodedCmdlet, err := powershell.EncodeAsUTF16LEBase64(cmdlet)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
	)

	stderrPipeReader, stderrPipeWriter := io.Pipe()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		defer stderrPipeWriter.Close()

		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)
		promptsAnswered := 0

		for {
			n, readErr := stderr.Read(buf)
			if n > 0 {
				accumulatedStderr.Write(buf[:n])
				text := accumulatedStderr.String()

				if strings.Contains(text, plugin.PluginReadyMessage) {
					close(readyChan)
					io.Copy(stderrPipeWriter, stderr)
					return
				}

				if strings.Contains(text, forgeGSudoPrompt) {
					if promptsAnswered >= 3 {
						errChan <- fmt.Errorf("too many gsudo password attempts for plugin at '%s'", path)
						stdin.Close()
						session.Close()
						session.Wait()
						return
					}

					promptsAnswered++
					_, err = stdin.Write([]byte(escalation.Pass() + "\n"))
					if err != nil {
						errChan <- fmt.Errorf("failed to write password to stdin for plugin at '%s': %w", path, err)
						stdin.Close()
						session.Close()
						session.Wait()
						return
					}

					accumulatedStderr.Reset()
				}
			}

			if readErr != nil {
				if readErr != io.EOF {
					errChan <- fmt.Errorf("error reading stderr for plugin at '%s': %w", path, readErr)
					stdin.Close()
					session.Close()
					session.Wait()
				}
				return
			}
		}
	}()

	err = session.Start(cmd)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to start remote plugin '%s': %w", path, err)
	}

	select {
	case <-ctx.Done():
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, fmt.Errorf("context cancelled while starting plugin at '%s': %w", path, ctx.Err())
	case err := <-errChan:
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, err
	case <-readyChan:
		return &sshPluginSession{
			session: session,
			stdout:  stdout,
			stderr:  stderrPipeReader,
			stdin:   stdin,
		}, nil
	}
}

func (s *sshWindowsPlatform) startPluginSessionAsSystem(
	ctx context.Context,
	remotePluginPath string,
) (plugin.Session, error) {
	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	cmdlet := fmt.Sprintf(`gsudo -s "%s"`, remotePluginPath)
	encodedCmdlet, err := powershell.EncodeAsUTF16LEBase64(cmdlet)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	cmd := fmt.Sprintf(
		"powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s",
		encodedCmdlet,
	)

	stderrPipeReader, stderrPipeWriter := io.Pipe()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		defer stderrPipeWriter.Close()

		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)

		for {
			n, readErr := stderr.Read(buf)
			if n > 0 {
				accumulatedStderr.Write(buf[:n])
				text := accumulatedStderr.String()

				if strings.Contains(text, plugin.PluginReadyMessage) {
					close(readyChan)
					io.Copy(stderrPipeWriter, stderr)
					return
				}
			}

			if readErr != nil {
				if readErr != io.EOF {
					errChan <- fmt.Errorf("error reading stderr for plugin at '%s': %w", remotePluginPath, readErr)
					stdin.Close()
					session.Close()
					session.Wait()
				}
				return
			}
		}
	}()

	err = session.Start(cmd)
	if err != nil {
		stdin.Close()
		session.Close()
		return nil, fmt.Errorf("failed to start remote plugin '%s': %w", remotePluginPath, err)
	}

	select {
	case <-ctx.Done():
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, fmt.Errorf("context cancelled while starting plugin at '%s': %w", remotePluginPath, ctx.Err())
	case err := <-errChan:
		stdin.Close()
		session.Close()
		session.Wait()
		return nil, err
	case <-readyChan:
		return &sshPluginSession{
			session: session,
			stdout:  stdout,
			stderr:  stderrPipeReader,
			stdin:   stdin,
		}, nil
	}
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
