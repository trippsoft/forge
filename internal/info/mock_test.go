package info

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/trippsoft/forge/internal/transport"
)

type commandResponse struct {
	stdout string
	stderr string
	err    error
}

type mockTransport struct {
	transportType transport.TransportType

	connectError error
	closeError   error

	defaultCommandResponse *commandResponse
	commandResponses       map[string]*commandResponse

	defaultPowerShellResponse *commandResponse
	powerShellResponses       map[string]*commandResponse

	fileSystem *mockFileSystem
}

func newMockTransport() *mockTransport {

	return &mockTransport{
		transportType:             transport.TransportTypeSSH,
		defaultCommandResponse:    &commandResponse{},
		commandResponses:          make(map[string]*commandResponse),
		defaultPowerShellResponse: &commandResponse{},
		powerShellResponses:       make(map[string]*commandResponse),
		fileSystem:                newMockFileSystem(),
	}
}

func (m *mockTransport) Type() transport.TransportType {
	return m.transportType
}

func (m *mockTransport) Connect() error {
	return m.connectError
}

func (m *mockTransport) Close() error {
	return m.closeError
}

func (m *mockTransport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {

	response, exists := m.commandResponses[command]
	if exists {
		return response.stdout, response.stderr, response.err
	}
	return m.defaultCommandResponse.stdout, m.defaultCommandResponse.stderr, m.defaultCommandResponse.err
}

func (m *mockTransport) ExecutePowerShell(ctx context.Context, command string) (string, error) {

	response, exists := m.powerShellResponses[command]
	if exists {
		return response.stdout, response.err
	}
	return m.defaultPowerShellResponse.stdout, m.defaultPowerShellResponse.err
}

func (m *mockTransport) FileSystem() transport.FileSystem {
	return m.fileSystem
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

type mockFile struct {
	content io.ReadCloser
	info    *mockFileInfo
}

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

type mockFileSystem struct {
	isNull bool

	connectError error
	closeError   error

	files map[string]*mockFile
	dirs  map[string]*mockFileInfo

	errorPaths map[string]error
}

func newMockFileSystem() *mockFileSystem {
	return &mockFileSystem{
		isNull:     false,
		files:      make(map[string]*mockFile),
		dirs:       make(map[string]*mockFileInfo),
		errorPaths: make(map[string]error),
	}
}

func (m *mockFileSystem) IsNull() bool {
	return m.isNull
}

func (m *mockFileSystem) Connect() error {
	return m.connectError
}

func (m *mockFileSystem) Close() error {
	return m.closeError
}

func (m *mockFileSystem) Stat(path string) (os.FileInfo, error) {
	if err, exists := m.errorPaths[path]; exists {
		return nil, err
	}
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
