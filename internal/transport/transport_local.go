package transport

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
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

	// pathPrefixes are the prefixes used for file paths.
	pathPrefixes []string
}

func NewLocalTransport() (Transport, error) {

	pathVariable := os.Getenv("PATH")
	if pathVariable == "" {
		return nil, fmt.Errorf("PATH environment variable is not set")
	}

	pathVariable = strings.TrimRight(strings.TrimSpace(pathVariable), string(os.PathListSeparator))
	pathPrefixes := strings.Split(pathVariable, string(os.PathListSeparator))

	for i, prefix := range pathPrefixes {
		if strings.HasSuffix(prefix, string(os.PathSeparator)) {
			pathPrefixes[i] = prefix + string(os.PathSeparator)
		}
	}

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
			pathPrefixes: pathPrefixes,
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
		pathPrefixes:      pathPrefixes,
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

// NewCommand creates a new command to be executed on the managed system.
func (l *localTransport) NewCommand(command string) *Cmd {
	return NewCmd(l, command)
}

// NewPowerShellCommand creates a new PowerShell command to be executed on the managed system.
func (l *localTransport) NewPowerShellCommand(command string) *PowerShellCmd {
	return NewPowerShellCmd(l, command)
}

// executeCommand implements Transport.
func (l *localTransport) executeCommand(ctx context.Context, cmd *Cmd) error {

	args := slices.Clone(l.shellArgs)
	args = append(args, cmd.command)

	execCmd := exec.CommandContext(ctx, l.shellCommand, args...)

	execCmd.Stdout = cmd.Stdout
	execCmd.Stderr = cmd.Stderr

	err := execCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

// executePowerShell implements Transport.
func (l *localTransport) executePowerShell(ctx context.Context, cmd *PowerShellCmd) error {

	if runtime.GOOS != "windows" {
		return fmt.Errorf("PowerShell execution is only supported on Windows")
	}

	args := slices.Clone(l.powershellArgs)
	encodedCommand, err := encodePowerShellAsUTF16LEBase64(cmd.command)
	if err != nil {
		return fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	args = append(args, encodedCommand)

	execCmd := exec.CommandContext(ctx, l.powershellCommand, args...)

	execCmd.Stdout = cmd.Stdout
	execCmd.Stderr = cmd.Stderr

	err = execCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute PowerShell command: %w", err)
	}

	return nil
}

// Stat implements Transport.
func (l *localTransport) Stat(path string) (os.FileInfo, error) {

	fileInfo, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil, nil // Return nil if the file does not exist
	}

	return fileInfo, err
}

// Create implements Transport.
func (l *localTransport) Create(path string) (File, error) {
	return os.Create(path)
}

// Open implements Transport.
func (l *localTransport) Open(path string) (File, error) {

	file, err := os.Open(path)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil, nil // Return nil if the file does not exist
	}

	return file, err
}

// Mkdir implements Transport.
func (l *localTransport) Mkdir(path string) error {

	err := os.Mkdir(path, 0755)
	if errors.Is(err, os.ErrExist) || errors.Is(err, syscall.EEXIST) {
		return nil // Directory already exists, return nil
	}

	return err
}

// MkdirAll implements Transport.
func (l *localTransport) MkdirAll(path string) error {

	return os.MkdirAll(path, 0755)
}

// Remove implements Transport.
func (l *localTransport) Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll implements Transport.
func (l *localTransport) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Join implements Transport.
func (l *localTransport) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// TempDir implements Transport.
func (l *localTransport) TempDir() (string, error) {
	return os.TempDir(), nil
}

// CreateTemp implements Transport.
func (l *localTransport) CreateTemp(dir, pattern string) (File, error) {
	return os.CreateTemp(dir, pattern)
}

// MkdirTemp implements Transport.
func (l *localTransport) MkdirTemp(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// Symlink implements Transport.
func (l *localTransport) Symlink(target, path string) error {
	return os.Symlink(target, path)
}

// ReadLink implements Transport.
func (l *localTransport) ReadLink(path string) (string, error) {
	return os.Readlink(path)
}

// RealPath implements Transport.
func (l *localTransport) RealPath(path string) (string, error) {

	if filepath.IsAbs(path) {
		return path, nil
	}

	for _, prefix := range l.pathPrefixes {
		absPath := filepath.Join(prefix, path)
		if fileInfo, _ := l.Stat(absPath); fileInfo != nil {
			return absPath, nil
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	absPath := filepath.Join(workingDir, path)

	if fileInfo, _ := l.Stat(absPath); fileInfo != nil {
		return absPath, nil
	}

	return "", os.ErrNotExist
}
