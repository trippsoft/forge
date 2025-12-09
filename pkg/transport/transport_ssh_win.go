// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/sftp"
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

	archCmd := "powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -Command " +
		`"Write-Host $env:PROCESSOR_ARCHITECTURE"`
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

func newSSHWindowsPlatform(t *sshTransport) sshPlatform {
	return &sshWindowsPlatform{
		t: t,
	}
}
