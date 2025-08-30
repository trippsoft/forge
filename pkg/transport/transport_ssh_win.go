// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type sshWindowsInfo struct {
	transport *sshTransport

	cachedTempDir      string
	cachedPathPrefixes []string
}

// pathSeparator implements sshPlatformInfo.
func (s *sshWindowsInfo) pathSeparator() rune {
	return '\\'
}

// pathListSeparator implements sshPlatformInfo.
func (s *sshWindowsInfo) pathListSeparator() rune {
	return ';'
}

// tempDir implements sshPlatformInfo.
func (s *sshWindowsInfo) tempDir() (string, error) {
	if s.cachedTempDir != "" {
		return s.cachedTempDir, nil // Return cached temp dir if available
	}

	err := s.transport.Connect() // Ensure we are connected
	if err != nil {
		return "", fmt.Errorf("failed to connect to SSH transport: %w", err)
	}

	cmd, err := s.transport.NewPowerShellCommand("$path = [System.IO.Path]::GetTempPath(); Write-Host $path", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create PowerShell command: %w", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %w", err)
	}

	stdout = strings.TrimRight(stdout, string(s.pathSeparator()))
	s.cachedTempDir = stdout

	return s.cachedTempDir, nil
}

// pathPrefixes implements sshPlatformInfo.
func (s *sshWindowsInfo) pathPrefixes() ([]string, error) {
	if s.cachedPathPrefixes != nil {
		return s.cachedPathPrefixes, nil // Already populated
	}

	cmd, err := s.transport.NewPowerShellCommand("Write-Host $env:PATH", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create PowerShell command: %w", err)
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to run PowerShell command: %w", err)
	}

	pathOutput := strings.TrimRight(strings.TrimSpace(stdout), string(s.pathListSeparator()))
	s.cachedPathPrefixes = strings.Split(pathOutput, string(s.pathListSeparator()))
	for i, prefix := range s.cachedPathPrefixes {
		if !strings.HasSuffix(prefix, string(s.pathSeparator())) {
			s.cachedPathPrefixes[i] = prefix + string(s.pathSeparator()) // Ensure each prefix ends with a separator
		}
	}

	return s.cachedPathPrefixes, nil
}

// newCommand implements sshPlatformInfo.
func (s *sshWindowsInfo) newCommand(command string, escalateConfig Escalation) (Cmd, error) {
	return s.newCommandImpl(command, escalateConfig)
}

// newPowerShellCommand implements sshPlatformInfo.
func (s *sshWindowsInfo) newPowerShellCommand(command string, escalateConfig Escalation) (Cmd, error) {
	encodedCommand, err := encodePowerShellAsUTF16LEBase64(command)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	command = fmt.Sprintf("powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand %s", encodedCommand)

	return s.newCommandImpl(command, escalateConfig)
}

// newPythonCommand implements sshPlatformInfo.
func (s *sshWindowsInfo) newPythonCommand(command string, escalateConfig Escalation) (Cmd, error) {
	return nil, errors.New("Python is not supported on Windows SSH transport")
}

func (s *sshWindowsInfo) newCommandImpl(command string, escalateConfig Escalation) (Cmd, error) {
	if escalateConfig == nil {
		return &sshCmd{
			transport: s.transport,
			command:   command,
		}, nil
	}

	return nil, errors.New("escalation is not supported for Windows SSH transport")
	// Windows SSH has the highest privileges available to the user without escalation, so runas does not apply.
}
