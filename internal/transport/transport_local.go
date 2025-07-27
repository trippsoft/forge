package transport

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
)

type localTransport struct {
	// shellCommand is the command used to execute shell commands.
	shellCommand string
	// shellArgs are the arguments passed to the shell command before the actual command.
	shellArgs []string

	// powershellCommand is the command used to execute PowerShell commands.
	powershellCommand string
	// powershellArgs are the arguments passed to the PowerShell command before the actual command.
	powershellArgs []string

	// fileSystem is the local file system.
	fileSystem FileSystem
}

func NewLocalTransport() (Transport, error) {

	if runtime.GOOS == "windows" {
		return &localTransport{
			shellCommand:      "cmd.exe",
			shellArgs:         []string{"/C"},
			powershellCommand: "powershell.exe",
			powershellArgs: []string{
				"-NoProfile",
				"-NonInteractive",
				"-ExecutionPolicy", "Bypass",
				"-EncodedCommand",
			},
			fileSystem: newLocalFileSystem(),
		}, nil
	}

	shellCommand, err := filepath.EvalSymlinks("/bin/sh")
	if err != nil {
		return nil, fmt.Errorf("failed to stat /bin/sh: %w", err)
	}

	return &localTransport{
		shellCommand:      shellCommand,
		shellArgs:         []string{"-c"},
		powershellCommand: "",
		powershellArgs:    nil,
		fileSystem:        newLocalFileSystem(),
	}, nil
}

// Type implements Transport.
func (l *localTransport) Type() TransportType {
	return TransportTypeLocal
}

// Connect implements Transport.
func (l *localTransport) Connect() error {
	return nil // No connection needed for local transport
}

// Close implements Transport.
func (l *localTransport) Close() error {
	return nil // No connection needed for local transport
}

// ExecuteCommand implements Transport.
func (l *localTransport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {

	args := slices.Clone(l.shellArgs)
	args = append(args, command)

	cmd := exec.CommandContext(ctx, l.shellCommand, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return outBuf.String(), errBuf.String(), fmt.Errorf("failed to execute command: %w", err)
	}

	return outBuf.String(), errBuf.String(), nil
}

// ExecutePowerShell implements Transport.
func (l *localTransport) ExecutePowerShell(ctx context.Context, command string) (string, error) {

	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("PowerShell execution is only supported on Windows")
	}

	args := slices.Clone(l.powershellArgs)
	encodedCommand, err := encodePowerShellAsUTF16LEBase64(command)
	if err != nil {
		return "", fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	args = append(args, encodedCommand)

	cmd := exec.CommandContext(ctx, l.powershellCommand, args...)

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	err = cmd.Run()
	if err != nil {
		return outBuf.String(), fmt.Errorf("failed to execute PowerShell command: %w", err)
	}

	return outBuf.String(), nil
}

// FileSystem implements Transport.
func (l *localTransport) FileSystem() FileSystem {
	return l.fileSystem
}

func detectPosixShell() (string, []string, error) {
	shellTarget, err := filepath.EvalSymlinks("/bin/sh")
	if err != nil {
		return "", nil, fmt.Errorf("failed to stat /bin/sh: %w", err)
	}

	return shellTarget, []string{"-c"}, nil
}

type localFileSystem struct{}

func newLocalFileSystem() FileSystem {
	return &localFileSystem{}
}

// Connect implements FileSystem.
func (l *localFileSystem) Connect() error {
	return nil // No connection needed for local file system
}

// Close implements FileSystem.
func (l *localFileSystem) Close() error {
	return nil // No connection needed for local file system
}

// IsNull implements FileSystem.
func (l *localFileSystem) IsNull() bool {
	return false // Local file system is always available
}

// Open implements FileSystem.
func (l *localFileSystem) Open(path string) (File, error) {
	return os.Open(path)
}

// Stat implements FileSystem.
func (l *localFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
