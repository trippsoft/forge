// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pkg/sftp"
	"github.com/trippsoft/forge/pkg/plugin"
	"golang.org/x/crypto/ssh"
)

type sshPlatform interface {
	OS() string
	Arch() string
	PathSeparator() string
	PluginExtension() string

	GetDefaultTempPath() (string, error)

	PopulateInfo() error

	MkdirAll(path string) error
	UploadFile(localPath, remotePath string) error

	StartPluginSession(
		ctx context.Context,
		basePath string,
		namespace string,
		pluginName string,
		escalation *Escalation,
	) (plugin.Session, error)
}

type sshTransport struct {
	host string
	port uint16

	platform sshPlatform

	tempPath string

	config     *ssh.ClientConfig
	client     *ssh.Client
	sftpClient *sftp.Client
}

// Type implements [Transport].
func (s *sshTransport) Type() TransportType {
	return TransportTypeSSH
}

// OS implements [Transport].
func (s *sshTransport) OS() (string, error) {
	if s.platform == nil || s.platform.OS() == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.platform.OS(), nil
}

// Arch implements [Transport].
func (s *sshTransport) Arch() (string, error) {
	if s.platform == nil || s.platform.Arch() == "" {
		err := s.Connect()
		if err != nil {
			return "", err
		}
	}

	return s.platform.Arch(), nil
}

// Connect implements [Transport].
func (s *sshTransport) Connect() error {
	if s.client != nil {
		return nil // Already connected
	}

	address := net.JoinHostPort(s.host, fmt.Sprintf("%d", s.port))
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

// Close implements [Transport].
func (s *sshTransport) Close() error {
	if s.sftpClient != nil {
		s.sftpClient.Close() // Close the SFTP client if it exists
		s.sftpClient = nil
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

// StartPluginSession implements [Transport].
func (s *sshTransport) StartPluginSession(
	ctx context.Context,
	basePath string,
	namespace string,
	pluginName string,
	escalation *Escalation,
) (plugin.Session, error) {
	err := s.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect before starting plugin: %w", err)
	}

	return s.platform.StartPluginSession(ctx, basePath, namespace, pluginName, escalation)
}

func (s *sshTransport) connectSFTP() error {
	if s.sftpClient != nil {
		return nil // Already connected
	}

	if s.client == nil {
		err := s.Connect()
		if err != nil {
			return err
		}
	}

	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	s.sftpClient = sftpClient
	return nil
}

type sshPluginSession struct {
	session *ssh.Session
	stdout  io.Reader
	stderr  io.ReadCloser
	stdin   io.WriteCloser
}

// Close implements [plugin.Session].
func (s *sshPluginSession) Close() error {
	s.stdin.Close()
	s.stderr.Close()
	err := s.session.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH session: %w", err)
	}

	return s.session.Wait()
}

// Stdout implements [plugin.Session].
func (s *sshPluginSession) Stdout() io.Reader {
	return s.stdout
}

// Stderr implements [plugin.Session].
func (s *sshPluginSession) Stderr() io.Reader {
	return s.stderr
}

// Stdin implements [plugin.Session].
func (s *sshPluginSession) Stdin() io.Writer {
	return s.stdin
}
