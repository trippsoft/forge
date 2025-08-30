// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

type sshSudoCmd struct {
	transport *sshTransport
	session   *ssh.Session
	ctx       context.Context
	completed bool

	command string

	password string
}

// OutputWithError implements Cmd.
func (s *sshSudoCmd) OutputWithError(ctx context.Context) (string, string, error) {
	err := s.createSession(ctx)
	if err != nil {
		return "", "", err
	}

	defer s.session.Close()

	var outBuf, errBuf bytes.Buffer
	stderrReader, err := s.session.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdinWriter, err := s.session.StdinPipe()
	teeReader := io.TeeReader(stderrReader, &errBuf)
	s.session.Stdout = &outBuf

	commandErrChannel := make(chan error)

	go func() {
		bufferReader := bufio.NewReader(teeReader)
		for {
			line, err := bufferReader.ReadString(':')
			if err != nil {
				if err == io.EOF {
					return
				}
				return
			}
			if strings.Contains(line, sshSudoPrompt) {
				_, err = stdinWriter.Write([]byte(s.password + "\n"))
				if err != nil {
					return
				}
			}
		}
	}()

	go func() {
		err := s.session.Run(s.command)
		commandErrChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return "", "", s.ctx.Err()
	case err = <-commandErrChannel:
		s.session = nil
		s.completed = true
		stdout := strings.TrimSpace(outBuf.String())
		stderr := strings.TrimSpace(errBuf.String())
		return stdout, stderr, err
	}
}

// Output implements Cmd.
func (s *sshSudoCmd) Output(ctx context.Context) (string, error) {
	err := s.createSession(ctx)
	if err != nil {
		return "", err
	}

	defer s.session.Close()

	var outBuf bytes.Buffer
	stderrReader, err := s.session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdinWriter, err := s.session.StdinPipe()
	s.session.Stdout = &outBuf

	commandErrChannel := make(chan error)

	go func() {
		bufferReader := bufio.NewReader(stderrReader)
		for {
			line, err := bufferReader.ReadString(':')
			if err != nil {
				if err == io.EOF {
					return
				}
				return
			}
			if strings.Contains(line, sshSudoPrompt) {
				_, err = stdinWriter.Write([]byte(s.password + "\n"))
				if err != nil {
					return
				}
			}
		}
	}()

	go func() {
		err := s.session.Run(s.command)
		commandErrChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return "", s.ctx.Err()
	case err = <-commandErrChannel:
		s.session = nil
		s.completed = true
		stdout := strings.TrimSpace(outBuf.String())
		return stdout, err
	}
}

// Run implements Cmd.
func (s *sshSudoCmd) Run(ctx context.Context) error {
	err := s.createSession(ctx)
	if err != nil {
		return err
	}

	defer s.session.Close()

	stderrReader, err := s.session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	stdinWriter, err := s.session.StdinPipe()

	commandErrChannel := make(chan error)

	go func() {
		bufferReader := bufio.NewReader(stderrReader)
		for {
			line, err := bufferReader.ReadString(':')
			if err != nil {
				if err == io.EOF {
					return
				}
				return
			}
			if strings.Contains(line, sshSudoPrompt) {
				_, err = stdinWriter.Write([]byte(s.password + "\n"))
				if err != nil {
					return
				}
			}
		}
	}()

	go func() {
		err := s.session.Run(s.command)
		commandErrChannel <- err
	}()

	select {
	case <-s.ctx.Done():
		s.session.Signal(ssh.SIGINT) // Send interrupt signal to the session
		s.session = nil
		s.completed = true
		return s.ctx.Err()
	case err = <-commandErrChannel:
		s.session = nil
		s.completed = true
		return err
	}
}

func (s *sshSudoCmd) createSession(ctx context.Context) error {
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

type sshPosixInfo struct {
	transport          *sshTransport
	pythonInterpreter  string
	cachedPathPrefixes []string
}

// canRunPowerShell implements sshPlatformInfo.
func (s *sshPosixInfo) canRunPowerShell() bool {
	return false
}

// canRunPython implements sshPlatformInfo.
func (s *sshPosixInfo) canRunPython() bool {
	return s.pythonInterpreter != ""
}

// pythonInterpreterPath implements sshPlatformInfo.
func (s *sshPosixInfo) pythonInterpreterPath() string {
	return s.pythonInterpreter
}

// pathSeparator implements sshPlatformInfo.
func (s *sshPosixInfo) pathSeparator() rune {
	return '/'
}

// pathListSeparator implements sshPlatformInfo.
func (s *sshPosixInfo) pathListSeparator() rune {
	return ':'
}

// tempDir implements sshPlatformInfo.
func (s *sshPosixInfo) tempDir() (string, error) {
	return "/tmp", nil
}

// pathPrefixes implements sshPlatformInfo.
func (s *sshPosixInfo) pathPrefixes() ([]string, error) {
	if s.cachedPathPrefixes != nil {
		return s.cachedPathPrefixes, nil // Already populated
	}

	cmd, err := s.transport.NewCommand("echo $PATH", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create command to get PATH: %w", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to run command to get PATH: %w", err)
	}

	pathOutput := strings.TrimRight(stdout, string(s.pathListSeparator()))
	s.cachedPathPrefixes = strings.Split(pathOutput, string(s.pathListSeparator()))
	for i, prefix := range s.cachedPathPrefixes {
		if !strings.HasSuffix(prefix, string(s.pathSeparator())) {
			s.cachedPathPrefixes[i] = prefix + string(s.pathSeparator()) // Ensure each prefix ends with a separator
		}
	}

	return s.cachedPathPrefixes, nil
}

// newCommand implements sshPlatformInfo.
func (s *sshPosixInfo) newCommand(command string, escalateConfig Escalation) (Cmd, error) {
	if escalateConfig == nil {
		return &sshCmd{
			transport: s.transport,
			command:   command,
		}, nil
	}

	username := escalateConfig.User()
	if username == "" {
		username = "root"
	}

	command = fmt.Sprintf("sudo -S -p '%s:' -u %s /bin/sh -c '%s'", sshSudoPrompt, username, command)

	return &sshSudoCmd{
		transport: s.transport,
		command:   command,
		password:  escalateConfig.Pass(),
	}, nil
}
