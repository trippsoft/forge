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

	stdout, stderr, err := transport.ExecutePowerShell(ctx, "Write-Host 'hello'")
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

	if stderr != "" {
		t.Errorf("Expected empty stderr, got '%s'", stderr)
	}
}

func TestNoneTransportExecutePowerShellWithContext(t *testing.T) {
	transport := NewNoneTransport()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	stdout, stderr, err := transport.ExecutePowerShell(ctx, "Write-Host 'hello'")
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

	if stderr != "" {
		t.Errorf("Expected empty stderr, got '%s'", stderr)
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

func TestNoneFileSystemStatEmptyPath(t *testing.T) {
	fs := &noneFileSystem{}

	// Test stat with empty path (should still fail with same error)
	info, err := fs.Stat("")
	if err == nil {
		t.Error("Expected error for Stat on none file system with empty path, but got none")
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

func TestNoneFileSystemOpenEmptyPath(t *testing.T) {
	fs := &noneFileSystem{}

	// Test open with empty path (should still fail with same error)
	file, err := fs.Open("")
	if err == nil {
		t.Error("Expected error for Open on none file system with empty path, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if file != nil {
		t.Errorf("Expected nil file, got %v", file)
	}
}

func TestNoneFileSystemDirectCreation(t *testing.T) {
	// Test creating noneFileSystem directly (not through transport)
	fs := &noneFileSystem{}

	// Test both methods fail appropriately
	info, err := fs.Stat("test")
	if err == nil {
		t.Error("Expected error for direct noneFileSystem Stat, but got none")
	}
	if info != nil {
		t.Error("Expected nil info for direct noneFileSystem Stat")
	}

	file, err := fs.Open("test")
	if err == nil {
		t.Error("Expected error for direct noneFileSystem Open, but got none")
	}
	if file != nil {
		t.Error("Expected nil file for direct noneFileSystem Open")
	}
}

func TestNoneTransportMultipleOperations(t *testing.T) {
	// Test that multiple operations on the same transport instance work consistently
	transport := NewNoneTransport()
	ctx := context.Background()

	// Test multiple command executions
	for i := 0; i < 3; i++ {
		stdout, stderr, err := transport.ExecuteCommand(ctx, "test command")
		if err == nil {
			t.Errorf("Expected error on attempt %d, but got none", i+1)
		}
		if stdout != "" || stderr != "" {
			t.Errorf("Expected empty output on attempt %d, got stdout='%s', stderr='%s'", i+1, stdout, stderr)
		}
	}

	// Test multiple PowerShell executions
	for i := 0; i < 3; i++ {
		stdout, stderr, err := transport.ExecutePowerShell(ctx, "test powershell")
		if err == nil {
			t.Errorf("Expected error on PowerShell attempt %d, but got none", i+1)
		}
		if stdout != "" || stderr != "" {
			t.Errorf("Expected empty output on PowerShell attempt %d, got stdout='%s', stderr='%s'", i+1, stdout, stderr)
		}
	}

	// Test multiple file system operations
	fs := transport.FileSystem()
	for i := 0; i < 3; i++ {
		info, err := fs.Stat("test")
		if err == nil {
			t.Errorf("Expected error on Stat attempt %d, but got none", i+1)
		}
		if info != nil {
			t.Errorf("Expected nil info on Stat attempt %d", i+1)
		}

		file, err := fs.Open("test")
		if err == nil {
			t.Errorf("Expected error on Open attempt %d, but got none", i+1)
		}
		if file != nil {
			t.Errorf("Expected nil file on Open attempt %d", i+1)
		}
	}
}
