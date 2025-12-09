// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/sftp"
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
