// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/trippsoft/forge/pkg/plugin"
)

type sshPosixPlatform struct {
	t *sshTransport

	os   string
	arch string
}

// OS implements [sshPlatform].
func (s *sshPosixPlatform) OS() string {
	return s.os
}

// Arch implements [sshPlatform].
func (s *sshPosixPlatform) Arch() string {
	return s.arch
}

// PathSeparator implements [sshPlatform].
func (s *sshPosixPlatform) PathSeparator() string {
	return "/"
}

// PluginExtension implements [sshPlatform].
func (s *sshPosixPlatform) PluginExtension() string {
	return ""
}

// GetDefaultTempPath implements [sshPlatform].
func (s *sshPosixPlatform) GetDefaultTempPath() (string, error) {
	err := s.t.connectSFTP()
	if err != nil {
		return "", err
	}

	// Many security policies restrict execution in /tmp or /var/tmp.  Is this a sensible default?

	fileInfo, err := s.t.sftpClient.Stat("/tmp") // Attempt /tmp first
	if err == nil && fileInfo.IsDir() {
		err = s.MkdirAll("/tmp/forge-tmp")
		if err == nil {
			err = s.t.sftpClient.Chmod("/tmp", 01777)
			if err == nil {
				return "/tmp/forge-tmp", nil
			}
		}
	}

	fileInfo, err = s.t.sftpClient.Stat("/var/tmp") // Attempt /var/tmp second
	if err == nil && fileInfo.IsDir() {
		err = s.MkdirAll("/var/tmp/forge-tmp")
		if err == nil {
			err = s.t.sftpClient.Chmod("/var/tmp", 01777)
			if err == nil {
				return "/var/tmp/forge-tmp", nil
			}
		}
	}

	// Fallback to user home directory (this is less ideal if impersonation is used)
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
	fileInfo, err = s.t.sftpClient.Stat(homeDir)
	if err == nil && fileInfo.IsDir() {
		tmpPath := fmt.Sprintf("%s/.local/share/forge-tmp", homeDir)
		err = s.MkdirAll(tmpPath)
		if err == nil {
			return tmpPath, nil
		}
	}

	return "", errors.New("failed to determine suitable remote temp path")
}

// PopulateInfo implements [sshPlatform].
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

// MkdirAll implements [sshPlatform].
func (s *sshPosixPlatform) MkdirAll(path string) error {
	err := s.t.connectSFTP()
	if err != nil {
		return err
	}

	err = s.t.sftpClient.MkdirAll(path)
	if err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", path, err)
	}

	return nil
}

// UploadFile implements [sshPlatform].
func (s *sshPosixPlatform) UploadFile(localPath, remotePath string) error {
	err := s.t.connectSFTP()
	if err != nil {
		return err
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

// StartPlugin implements [sshPlatform].
func (s *sshPosixPlatform) StartPluginSession(
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

	remotePluginPath := fmt.Sprintf("%s/%s-%s", s.t.tempPath, namespace, pluginName)

	err = s.MkdirAll(s.t.tempPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote temp path '%s': %w", s.t.tempPath, err)
	}

	err = s.UploadFile(localPluginPath, remotePluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload discovery plugin to remote path '%s': %w", remotePluginPath, err)
	}

	session, err := s.t.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	err = session.Run(fmt.Sprintf("chmod +x %s", remotePluginPath))
	session.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to set execute permission on remote plugin '%s': %w", remotePluginPath, err)
	}

	if escalation != nil {
		return s.startEscalatedPluginSession(ctx, remotePluginPath, escalation)
	}

	session, err = s.t.client.NewSession()
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

	err = session.Start(remotePluginPath)
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

func (s *sshPosixPlatform) startEscalatedPluginSession(
	ctx context.Context,
	path string,
	escalation *Escalation,
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

	err = session.Start(path)
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to start plugin at '%s': %w", path, err)
	}

	stderrPipeReader, stderrPipeWriter := io.Pipe()

	readyChan := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		defer stderrPipeWriter.Close()

		var accumulatedStderr strings.Builder
		buf := make([]byte, 4096)
		promptsAnswered := 0
		sudoPromptSuffix := forgeSudoPrompt + ":"

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

				if strings.HasSuffix(text, sudoPromptSuffix) {
					if promptsAnswered >= 3 {
						errChan <- fmt.Errorf("too many sudo password attempts for plugin at '%s'", path)
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
