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

// canRunPowerShell implements sshPlatformInfo.
func (s *sshWindowsInfo) canRunPowerShell() bool {
	return true
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

	cmd, err := s.transport.NewPowerShellCommand("$path = [System.IO.Path]::GetTempPath(); Write-Host $path")
	if err != nil {
		return "", fmt.Errorf("failed to create PowerShell command: %w", err)
	}

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %w", err)
	}

	stdout := strings.TrimRight(strings.TrimSpace(string(stdoutBytes)), string(s.pathSeparator()))

	s.cachedTempDir = stdout

	return s.cachedTempDir, nil
}

// pathPrefixes implements sshPlatformInfo.
func (s *sshWindowsInfo) pathPrefixes() ([]string, error) {

	if s.cachedPathPrefixes != nil {
		return s.cachedPathPrefixes, nil // Already populated
	}

	cmd, err := s.transport.NewPowerShellCommand("Write-Host $env:PATH")
	if err != nil {
		return nil, fmt.Errorf("failed to create PowerShell command: %w", err)
	}

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to run PowerShell command: %w", err)
	}

	stdout := strings.TrimSpace(string(stdoutBytes))

	pathOutput := strings.TrimRight(strings.TrimSpace(stdout), string(s.pathListSeparator()))

	s.cachedPathPrefixes = strings.Split(pathOutput, string(s.pathListSeparator()))

	for i, prefix := range s.cachedPathPrefixes {
		if !strings.HasSuffix(prefix, string(s.pathSeparator())) {
			s.cachedPathPrefixes[i] = prefix + string(s.pathSeparator()) // Ensure each prefix ends with a separator
		}
	}

	return s.cachedPathPrefixes, nil
}

// newEscalatedCommand implements sshPlatformInfo.
func (s *sshWindowsInfo) newEscalatedCommand(command string, config *EscalationConfig) (Cmd, error) {
	return nil, errors.New("escalated commands are not supported on Windows via SSH")
	// Windows over SSH automatically assumes elevated permissions, if the user has them, as of 08/2025
}
