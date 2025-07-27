package transport

import (
	"context"
	"errors"
	"os"
)

type noneTransport struct {
	fileSystem FileSystem
}

func NewNoneTransport() Transport {
	return &noneTransport{
		fileSystem: &noneFileSystem{},
	}
}

// Type implements Transport.
func (n *noneTransport) Type() TransportType {
	return TransportTypeNone
}

// Connect implements Transport.
func (n *noneTransport) Connect() error {
	return nil
}

// Close implements Transport.
func (n *noneTransport) Close() error {
	return nil
}

// ExecuteCommand implements Transport.
func (n *noneTransport) ExecuteCommand(ctx context.Context, command string) (stdout string, stderr string, err error) {
	return "", "", errors.New("no transport available for command execution")
}

// ExecutePowerShell implements Transport.
func (n *noneTransport) ExecutePowerShell(ctx context.Context, command string) (stdout string, err error) {
	return "", errors.New("no transport available for PowerShell execution")
}

// FileSystem implements Transport.
func (n *noneTransport) FileSystem() FileSystem {
	return n.fileSystem
}

type noneFileSystem struct{}

// Connect implements FileSystem.
func (n *noneFileSystem) Connect() error {
	return nil // No connection needed for none file system
}

// Close implements FileSystem.
func (n *noneFileSystem) Close() error {
	return nil // No connection needed for none file system
}

// IsNull implements FileSystem.
func (n *noneFileSystem) IsNull() bool {
	return true // None file system is always null
}

// Stat implements FileSystem.
func (n *noneFileSystem) Stat(path string) (os.FileInfo, error) {
	return nil, errors.New("no file system available")
}

// Open implements FileSystem.
func (n *noneFileSystem) Open(path string) (File, error) {
	return nil, errors.New("no file system available")
}
