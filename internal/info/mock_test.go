package info

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/trippsoft/forge/internal/transport"
)

// Mock transport for testing
type mockTransport struct {
	powershellOutput string
	commandOutput    string
	shouldError      bool
}

func (m *mockTransport) Type() transport.TransportType {
	return transport.TransportTypeNone
}

func (m *mockTransport) Connect() error {
	return nil
}

func (m *mockTransport) Close() error {
	return nil
}

func (m *mockTransport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {
	if m.shouldError {
		return "", "", io.EOF
	}
	return m.commandOutput, "", nil
}

func (m *mockTransport) ExecutePowerShell(ctx context.Context, command string) (string, string, error) {
	if m.shouldError {
		return "", "", io.EOF
	}
	return m.powershellOutput, "", nil
}

func (m *mockTransport) FileSystem() transport.FileSystem {
	return &mockFileSystem{}
}

// Mock file system for testing
type mockFileSystem struct {
	files map[string]*mockFile
	dirs  map[string]*mockFileInfo
}

type mockFile struct {
	content io.ReadCloser
	info    *mockFileInfo
}

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	isDir   bool
	modTime time.Time
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func (m *mockFile) Read(p []byte) (n int, err error) {
	return m.content.Read(p)
}

func (m *mockFile) Write(p []byte) (n int, err error) {
	// For testing purposes, we don't need to implement write
	return len(p), nil
}

func (m *mockFile) Close() error {
	return m.content.Close()
}

func (m *mockFile) Name() string {
	if m.info != nil {
		return m.info.Name()
	}
	return ""
}

func (m *mockFileSystem) IsNull() bool {
	return false
}

func (m *mockFileSystem) Stat(path string) (os.FileInfo, error) {
	if file, exists := m.files[path]; exists {
		return file.info, nil
	}
	if dir, exists := m.dirs[path]; exists {
		return dir, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockFileSystem) Open(path string) (transport.File, error) {
	if file, exists := m.files[path]; exists {
		return file, nil
	}
	return nil, os.ErrNotExist
}

type errorTransport struct {
	shouldError bool
}

func (e *errorTransport) Type() transport.TransportType {
	return transport.TransportTypeNone
}

func (e *errorTransport) Connect() error {
	if e.shouldError {
		return errors.New("connection failed")
	}
	return nil
}

func (e *errorTransport) Close() error {
	return nil
}

func (e *errorTransport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {
	if e.shouldError {
		return "", "", errors.New("command failed")
	}
	return "", "", nil
}

func (e *errorTransport) ExecutePowerShell(ctx context.Context, command string) (string, string, error) {
	if e.shouldError {
		return "", "", errors.New("powershell failed")
	}
	return "", "", nil
}

func (e *errorTransport) FileSystem() transport.FileSystem {
	return &errorFileSystem{shouldError: e.shouldError}
}

type errorFileSystem struct {
	shouldError bool
}

func (e *errorFileSystem) IsNull() bool {
	return e.shouldError
}

func (e *errorFileSystem) Stat(path string) (os.FileInfo, error) {
	if e.shouldError {
		return nil, errors.New("stat failed")
	}
	return nil, os.ErrNotExist
}

func (e *errorFileSystem) Open(path string) (transport.File, error) {
	if e.shouldError {
		return nil, errors.New("open failed")
	}
	return nil, os.ErrNotExist
}

type errorReadCloser struct {
	shouldError bool
}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	if e.shouldError {
		return 0, errors.New("read failed")
	}
	return 0, io.EOF
}

func (e *errorReadCloser) Close() error {
	return nil
}
