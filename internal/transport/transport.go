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

	// ExecuteCommand executes a command on the managed system and returns the output.
	ExecuteCommand(ctx context.Context, command string) (stdout string, stderr string, err error)
	// ExecutePowerShell executes a PowerShell command on the managed system and returns the output.
	ExecutePowerShell(ctx context.Context, command string) (stdout string, stderr string, err error)

	// FileSystem returns a FileSystem interface for interacting with the file system of the managed system.
	FileSystem() FileSystem
}

// FileSystem interface defines methods for file system operations.
type FileSystem interface {
	// IsNull checks if the file system is null or not supported.
	IsNull() bool

	// Connect establishes a connection to the file system if needed.
	Connect() error

	// Close closes the file system connection if needed.
	Close() error

	// Stat retrieves the file information for the given path.
	Stat(path string) (os.FileInfo, error)
	// Open opens an existing file with the specified path and flags.
	Open(path string) (File, error)
}

// File interface defines methods for file operations.
type File interface {
	io.ReadWriteCloser

	// Name returns the name of the file.
	Name() string
}
