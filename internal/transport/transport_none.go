package transport

import (
	"context"
	"errors"
	"os"
)

var (
	TransportNone Transport = &noneTransport{}
	cmdNone       Cmd       = &noneCmd{}
)

type noneCmd struct{}

// OutputWithError implements Cmd.
func (n *noneCmd) OutputWithError(ctx context.Context) (stdout []byte, stderr []byte, err error) {
	return nil, nil, errors.New("no transport available for command execution")
}

// Output implements Cmd.
func (n *noneCmd) Output(ctx context.Context) ([]byte, error) {
	return nil, errors.New("no transport available for command execution")
}

// Run implements Cmd.
func (n *noneCmd) Run(ctx context.Context) error {
	return errors.New("no transport available for command execution")
}

type noneTransport struct{}

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

// NewCommand creates a new command to be executed on the managed system.
func (n *noneTransport) NewCommand(command string) Cmd {
	return cmdNone
}

// NewEscalatedCommand creates a new command to be executed with privilege escalation.
func (n *noneTransport) NewEscalatedCommand(command string, escalationConfig *EscalationConfig) (Cmd, error) {
	return cmdNone, errors.New("no transport available for escalated command execution")
}

// NewPowerShellCommand creates a new PowerShell command to be executed on the managed system.
func (n *noneTransport) NewPowerShellCommand(command string) (Cmd, error) {
	return nil, errors.New("no transport available for PowerShell execution")
}

// NewEscalatedPowerShellCommand creates a new PowerShell command to be executed with privilege escalation.
func (n *noneTransport) NewEscalatedPowerShellCommand(command string, escalationConfig *EscalationConfig) (Cmd, error) {
	return nil, errors.New("no transport available for escalated PowerShell execution")
}

// Stat implements Transport.
func (n *noneTransport) Stat(path string) (os.FileInfo, error) {
	return nil, errors.New("no file system available")
}

// Create implements Transport.
func (n *noneTransport) Create(path string) (File, error) {
	return nil, errors.New("no file system available")
}

// Open implements Transport.
func (n *noneTransport) Open(path string) (File, error) {
	return nil, errors.New("no file system available")
}

// Mkdir implements Transport.
func (n *noneTransport) Mkdir(path string) error {
	return errors.New("no file system available")
}

// MkdirAll implements Transport.
func (n *noneTransport) MkdirAll(path string) error {
	return errors.New("no file system available")
}

// Remove implements Transport.
func (n *noneTransport) Remove(path string) error {
	return errors.New("no file system available")
}

// RemoveAll implements Transport.
func (n *noneTransport) RemoveAll(path string) error {
	return errors.New("no file system available")
}

// Join implements Transport.
func (n *noneTransport) Join(elem ...string) string {
	return ""
}

// TempDir implements Transport.
func (n *noneTransport) TempDir() (string, error) {
	return "", errors.New("no file system available")
}

// CreateTemp implements Transport.
func (n *noneTransport) CreateTemp(dir string, pattern string) (File, error) {
	return nil, errors.New("no file system available")
}

// MkdirTemp implements Transport.
func (n *noneTransport) MkdirTemp(dir string, pattern string) (string, error) {
	return "", errors.New("no file system available")
}

// Symlink implements Transport.
func (n *noneTransport) Symlink(target, path string) error {
	return errors.New("no file system available")
}

// Readlink implements Transport.
func (n *noneTransport) ReadLink(path string) (string, error) {
	return "", errors.New("no file system available")
}

// RealPath implements Transport.
func (n *noneTransport) RealPath(path string) (string, error) {
	return "", errors.New("no file system available")
}
