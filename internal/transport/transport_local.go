package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
)

type localCmd struct {
	cmd *exec.Cmd

	command string
	args    []string

	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
}

// Run implements Cmd.
func (c *localCmd) Run(ctx context.Context) error {

	if c.cmd != nil {
		return errors.New("command already started")
	}

	c.cmd = c.createCommand(ctx)

	err := c.cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command '%s %s': %w", c.command, strings.Join(c.args, " "), err)
	}

	return nil
}

// Start implements Cmd.
func (c *localCmd) Start(ctx context.Context) error {

	if c.cmd != nil {
		return errors.New("command already started")
	}

	c.cmd = c.createCommand(ctx)

	err := c.cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command '%s %s': %w", c.command, strings.Join(c.args, " "), err)
	}

	return nil
}

// Wait implements Cmd.
func (c *localCmd) Wait() error {

	if c.cmd == nil {
		return errors.New("command not started")
	}

	return c.cmd.Wait()
}

// SetStdout implements Cmd.
func (c *localCmd) SetStdout(stdout io.Writer) error {

	if c.cmd != nil {
		return errors.New("command already started")
	}

	if c.stdout != nil {
		return errors.New("stdout already set")
	}

	c.stdout = stdout
	return nil
}

// SetStderr implements Cmd.
func (c *localCmd) SetStderr(stderr io.Writer) error {

	if c.cmd != nil {
		return errors.New("command already started")
	}

	if c.stderr != nil {
		return errors.New("stderr already set")
	}

	c.stderr = stderr
	return nil
}

// StdoutPipe implements Cmd.
func (c *localCmd) StdoutPipe() (io.ReadCloser, error) {

	if c.cmd != nil {
		return nil, errors.New("command already started")
	}

	if c.stdout != nil {
		return nil, errors.New("stdout already set")
	}

	pipeReader, pipeWriter := io.Pipe()

	c.stdout = pipeWriter

	return pipeReader, nil
}

// StderrPipe implements Cmd.
func (c *localCmd) StderrPipe() (io.ReadCloser, error) {

	if c.cmd != nil {
		return nil, errors.New("command already started")
	}

	if c.stderr != nil {
		return nil, errors.New("stderr already set")
	}

	pipeReader, pipeWriter := io.Pipe()

	c.stderr = pipeWriter

	return pipeReader, nil
}

// StdinPipe implements Cmd.
func (c *localCmd) StdinPipe() (io.WriteCloser, error) {

	if c.cmd != nil {
		return nil, errors.New("command already started")
	}

	if c.stdin != nil {
		return nil, errors.New("stdin pipe already created")
	}

	pipeReader, pipeWriter := io.Pipe()

	c.stdin = pipeReader

	return pipeWriter, nil
}

func (c *localCmd) createCommand(ctx context.Context) *exec.Cmd {

	cmd := exec.CommandContext(ctx, c.command, c.args...)

	if c.stdout != nil {
		cmd.Stdout = c.stdout
	}

	if c.stderr != nil {
		cmd.Stderr = c.stderr
	}

	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}

	return cmd
}

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
func (l *localTransport) NewCommand(command string) Cmd {

	args := slices.Clone(l.shellArgs)
	args = append(args, command)

	return &localCmd{
		command: l.shellCommand,
		args:    args,
	}
}

// NewPowerShellCommand creates a new PowerShell command to be executed on the managed system.
func (l *localTransport) NewPowerShellCommand(command string) (Cmd, error) {

	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("PowerShell execution is only supported on Windows")
	}

	args := slices.Clone(l.powershellArgs)
	encodedCommand, err := encodePowerShellAsUTF16LEBase64(command)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PowerShell command: %w", err)
	}

	args = append(args, encodedCommand)

	return &localCmd{
		command: l.powershellCommand,
		args:    args,
	}, nil
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
