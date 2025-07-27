package transport

import (
	"context"
	"testing"
)

func TestNewNoneTransport(t *testing.T) {
	transport := NewNoneTransport()
	if transport == nil {
		t.Fatal("NewNoneTransport returned nil transport")
	}

	if transport.Type() != TransportTypeNone {
		t.Errorf("Expected transport type %s, got %s", TransportTypeNone, transport.Type())
	}

	// Verify that FileSystem is properly initialized
	fs := transport.FileSystem()
	if fs == nil {
		t.Fatal("FileSystem is nil")
	}

	// Verify it's the correct type
	_, ok := fs.(*noneFileSystem)
	if !ok {
		t.Error("FileSystem did not return a noneFileSystem instance")
	}
}

func TestNoneTransportType(t *testing.T) {
	transport := NewNoneTransport()

	if transport.Type() != TransportTypeNone {
		t.Errorf("Expected transport type %s, got %s", TransportTypeNone, transport.Type())
	}
}

func TestNoneTransportConnect(t *testing.T) {
	transport := NewNoneTransport()

	err := transport.Connect()
	if err != nil {
		t.Errorf("Connect should not return error for none transport, got: %v", err)
	}
}

func TestNoneTransportClose(t *testing.T) {
	transport := NewNoneTransport()

	err := transport.Close()
	if err != nil {
		t.Errorf("Close should not return error for none transport, got: %v", err)
	}
}

func TestNoneTransportExecuteCommand(t *testing.T) {
	transport := NewNoneTransport()
	ctx := context.Background()

	stdout, stderr, err := transport.ExecuteCommand(ctx, "echo hello")
	if err == nil {
		t.Error("Expected error for ExecuteCommand on none transport, but got none")
	}

	expectedError := "no transport available for command execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if stdout != "" {
		t.Errorf("Expected empty stdout, got '%s'", stdout)
	}

	if stderr != "" {
		t.Errorf("Expected empty stderr, got '%s'", stderr)
	}
}

func TestNoneTransportExecuteCommandWithContext(t *testing.T) {
	transport := NewNoneTransport()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	stdout, stderr, err := transport.ExecuteCommand(ctx, "echo hello")
	if err == nil {
		t.Error("Expected error for ExecuteCommand on none transport, but got none")
	}

	expectedError := "no transport available for command execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if stdout != "" {
		t.Errorf("Expected empty stdout, got '%s'", stdout)
	}

	if stderr != "" {
		t.Errorf("Expected empty stderr, got '%s'", stderr)
	}
}

func TestNoneTransportExecutePowerShell(t *testing.T) {
	transport := NewNoneTransport()
	ctx := context.Background()

	stdout, err := transport.ExecutePowerShell(ctx, "Write-Host 'hello'")
	if err == nil {
		t.Error("Expected error for ExecutePowerShell on none transport, but got none")
	}

	expectedError := "no transport available for PowerShell execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if stdout != "" {
		t.Errorf("Expected empty stdout, got '%s'", stdout)
	}
}

func TestNoneTransportExecutePowerShellWithContext(t *testing.T) {
	transport := NewNoneTransport()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	stdout, err := transport.ExecutePowerShell(ctx, "Write-Host 'hello'")
	if err == nil {
		t.Error("Expected error for ExecutePowerShell on none transport, but got none")
	}

	expectedError := "no transport available for PowerShell execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if stdout != "" {
		t.Errorf("Expected empty stdout, got '%s'", stdout)
	}
}

func TestNoneTransportFileSystem(t *testing.T) {
	transport := NewNoneTransport()

	fs := transport.FileSystem()
	if fs == nil {
		t.Fatal("FileSystem returned nil")
	}

	// Test that it's actually a noneFileSystem
	_, ok := fs.(*noneFileSystem)
	if !ok {
		t.Error("FileSystem did not return a noneFileSystem instance")
	}
}

func TestNoneFileSystemIsNull(t *testing.T) {
	fs := &noneFileSystem{}

	if !fs.IsNull() {
		t.Error("IsNull should return true for none file system")
	}
}

func TestNoneFileSystemConnect(t *testing.T) {
	fs := &noneFileSystem{}

	err := fs.Connect()
	if err != nil {
		t.Errorf("Connect should not return error for none file system, got: %v", err)
	}
}

func TestNoneFileSystemClose(t *testing.T) {
	fs := &noneFileSystem{}

	err := fs.Close()
	if err != nil {
		t.Errorf("Close should not return error for none file system, got: %v", err)
	}
}

func TestNoneFileSystemStat(t *testing.T) {
	fs := &noneFileSystem{}

	// Test stat operation (should always fail)
	info, err := fs.Stat("/some/path/file.txt")
	if err == nil {
		t.Error("Expected error for Stat on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if info != nil {
		t.Errorf("Expected nil file info, got %v", info)
	}
}

func TestNoneFileSystemOpen(t *testing.T) {
	fs := &noneFileSystem{}

	// Test open operation (should always fail)
	file, err := fs.Open("/some/path/file.txt")
	if err == nil {
		t.Error("Expected error for Open on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if file != nil {
		t.Errorf("Expected nil file, got %v", file)
	}
}
