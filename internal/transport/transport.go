package transport

import (
	"context"
	"io"
	"os"
)

type TransportType string

const (
	TransportTypeNone  TransportType = "none"
	TransportTypeLocal TransportType = "local"
	TransportTypeSSH   TransportType = "ssh"
	TransportTypeWinRM TransportType = "winrm"
)

// Transport interface defines the methods for interacting with a managed system.
type Transport interface {
	// Type returns the type of transport.
	Type() TransportType

	// Connect establishes a connection to the managed system.
	Connect() error
	// Close closes the connection to the managed system.
	Close() error

	// NewCommand creates a new command to be executed on the managed system.
	NewCommand(command string) Cmd
	// NewPowerShellCommand creates a new PowerShell command to be executed on the managed system.
	NewPowerShellCommand(command string) (Cmd, error)

	// Stat retrieves the file information for the given path on the managed system.
	Stat(path string) (os.FileInfo, error)
	// Create creates or truncates a file with the specified path and flags on the managed system.
	Create(path string) (File, error)
	// Open opens an existing file with the specified path and flags on the managed system.
	Open(path string) (File, error)
	// Mkdir creates a directory with the specified path on the managed system.
	Mkdir(path string) error
	// MkdirAll creates a directory and all necessary parents with the specified path on the managed system.
	MkdirAll(path string) error
	// Remove removes the file or directory at the specified path on the managed system.
	Remove(path string) error
	// RemoveAll removes the file or directory and all its contents at the specified path on the managed system.
	RemoveAll(path string) error

	// Join joins the directory and name into a single path.
	Join(elem ...string) string

	// TempDir returns the default temporary directory for the managed system.
	TempDir() (string, error)
	// CreateTemp creates a temporary file in the managed system's temporary directory.
	CreateTemp(dir, pattern string) (File, error)
	// MkdirTemp creates a temporary directory in the managed system's temporary directory.
	MkdirTemp(dir, pattern string) (string, error)

	// Symlink creates a symbolic link at the specified path pointing to the target.
	Symlink(target, path string) error
	// ReadLink reads the target of a symbolic link at the specified path.
	ReadLink(path string) (string, error)

	// RealPath returns the absolute path of the specified path on the managed system.
	RealPath(path string) (string, error)
}

type Cmd interface {
	// Run executes the command on the managed system.
	Run(ctx context.Context) error
	// Start starts the command execution on the managed system.
	Start(ctx context.Context) error
	// Wait waits for the command to finish execution.
	Wait() error

	// SetStdout sets the standard output for the command.
	SetStdout(stdout io.Writer) error
	// SetStderr sets the standard error for the command.
	SetStderr(stderr io.Writer) error

	// StdoutPipe returns a writer to which the command's standard output will be written.
	StdoutPipe() (io.ReadCloser, error)
	// StderrPipe returns a writer to which the command's standard error will be written.
	StderrPipe() (io.ReadCloser, error)
	// StdinPipe returns a reader from which the command's standard input can be read.
	StdinPipe() (io.WriteCloser, error)
}

// File interface defines methods for file operations.
type File interface {
	io.ReadWriteCloser

	// Name returns the name of the file.
	Name() string
	// Sync synchronizes the file's contents with the underlying storage.
	Sync() error
}
