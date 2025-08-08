package transport

import (
	"context"
	"io"
	"os"
)

type TransportType string

const (
	TransportTypeNone TransportType = "none"
	TransportTypeSSH  TransportType = "ssh"
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
	NewCommand(command string, escalateConfig EscalateConfig) (Cmd, error)
	// NewPowerShellCommand creates a new PowerShell command to be executed on the managed system.
	NewPowerShellCommand(command string, escalateConfig EscalateConfig) (Cmd, error)

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

// Cmd interface defines methods for executing commands on the managed system.
type Cmd interface {
	// OutputWithError executes the command and returns its combined standard output and standard error.
	OutputWithError(ctx context.Context) (stdout string, stderr string, err error)
	// Output executes the command and returns its standard output.
	Output(ctx context.Context) (string, error)
	// Run executes the command on the managed system with no output.
	Run(ctx context.Context) error
}

// EscalationConfig holds configuration for privilege escalation.
type EscalationConfig struct {
	User     string // User specifies the user as which to run the command.
	Password string // Password specifies the password for the user.
}

// File interface defines methods for file operations.
type File interface {
	io.ReadWriteCloser

	// Name returns the name of the file.
	Name() string
	// Sync synchronizes the file's contents with the underlying storage.
	Sync() error
}
