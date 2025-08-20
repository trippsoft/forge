// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"errors"
	"os"
)

var (
	TransportNone Transport = &noneTransport{}
)

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

// NewCommand implements Transport.
func (n *noneTransport) NewCommand(command string, escalateConfig Escalation) (Cmd, error) {
	return nil, errors.New("no transport available for command execution")
}

// NewPowerShellCommand implements Transport.
func (n *noneTransport) NewPowerShellCommand(command string, escalateConfig Escalation) (Cmd, error) {
	return nil, errors.New("no transport available for PowerShell execution")
}

// NewPythonCommand implements Transport.
func (n *noneTransport) NewPythonCommand(interpreter, command string, escalateConfig Escalation) (Cmd, error) {
	return nil, errors.New("no transport available for Python execution")
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
