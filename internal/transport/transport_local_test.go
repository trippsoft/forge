package transport

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestNewLocalTransport(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	if transport == nil {
		t.Fatal("NewLocalTransport returned nil transport")
	}

	if transport.Type() != TransportTypeLocal {
		t.Errorf("Expected transport type %s, got %s", TransportTypeLocal, transport.Type())
	}

	// Verify that FileSystem is properly initialized
	fs := transport.FileSystem()
	if fs == nil {
		t.Fatal("FileSystem is nil")
	}
}

func TestNewLocalTransportWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	localTransport := transport.(*localTransport)

	// Verify Windows-specific configuration
	if localTransport.shellCommand != "cmd.exe" {
		t.Errorf("Expected shell command 'cmd.exe', got '%s'", localTransport.shellCommand)
	}

	expectedShellArgs := []string{"/C"}
	if len(localTransport.shellArgs) != len(expectedShellArgs) {
		t.Errorf("Expected shell args %v, got %v", expectedShellArgs, localTransport.shellArgs)
	}

	if localTransport.powershellCommand != "powershell.exe" {
		t.Errorf("Expected PowerShell command 'powershell.exe', got '%s'", localTransport.powershellCommand)
	}

	expectedPSArgs := []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-EncodedCommand"}
	if len(localTransport.powershellArgs) != len(expectedPSArgs) {
		t.Errorf("Expected PowerShell args %v, got %v", expectedPSArgs, localTransport.powershellArgs)
	}
}

func TestNewLocalTransportUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	localTransport := transport.(*localTransport)

	// Verify Unix-specific configuration
	if !strings.HasSuffix(localTransport.shellCommand, "sh") {
		t.Errorf("Expected shell command to end with 'sh', got '%s'", localTransport.shellCommand)
	}

	expectedShellArgs := []string{"-c"}
	if len(localTransport.shellArgs) != len(expectedShellArgs) {
		t.Errorf("Expected shell args %v, got %v", expectedShellArgs, localTransport.shellArgs)
	}

	if localTransport.powershellCommand != "" {
		t.Errorf("Expected empty PowerShell command on Unix, got '%s'", localTransport.powershellCommand)
	}

	if localTransport.powershellArgs != nil {
		t.Errorf("Expected nil PowerShell args on Unix, got %v", localTransport.powershellArgs)
	}
}

func TestLocalTransportType(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	if transport.Type() != TransportTypeLocal {
		t.Errorf("Expected transport type %s, got %s", TransportTypeLocal, transport.Type())
	}
}

func TestLocalTransportConnect(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	err = transport.Connect()
	if err != nil {
		t.Errorf("Connect should not return error for local transport, got: %v", err)
	}
}

func TestLocalTransportClose(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	err = transport.Close()
	if err != nil {
		t.Errorf("Close should not return error for local transport, got: %v", err)
	}
}

func TestLocalTransportExecuteCommand(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	// Test simple echo command
	var command string
	if runtime.GOOS == "windows" {
		command = "echo hello"
	} else {
		command = "echo hello"
	}

	stdout, stderr, err := transport.ExecuteCommand(ctx, command)
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v", err)
	}

	if !strings.Contains(stdout, "hello") {
		t.Errorf("Expected stdout to contain 'hello', got: %s", stdout)
	}

	if stderr != "" {
		t.Logf("Stderr (might be empty): %s", stderr)
	}
}

func TestLocalTransportExecuteCommandWithContext(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var command string
	if runtime.GOOS == "windows" {
		command = "timeout 1" // Windows timeout command
	} else {
		command = "sleep 1" // Unix sleep command
	}

	_, _, err = transport.ExecuteCommand(ctx, command)
	if err == nil {
		t.Error("Expected error due to context timeout, but got none")
	}
}

func TestLocalTransportExecuteCommandError(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	// Test command that should fail
	command := "nonexistentcommand12345"
	stdout, stderr, err := transport.ExecuteCommand(ctx, command)
	if err == nil {
		t.Error("Expected error for nonexistent command, but got none")
	}

	t.Logf("Error (expected): %v", err)
	t.Logf("Stdout: %s", stdout)
	t.Logf("Stderr: %s", stderr)
}

func TestLocalTransportExecutePowerShellWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping PowerShell test on non-Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	// Test simple PowerShell command
	command := "Write-Host 'Hello PowerShell'"
	stdout, stderr, err := transport.ExecutePowerShell(ctx, command)
	if err != nil {
		t.Fatalf("ExecutePowerShell failed: %v", err)
	}

	if !strings.Contains(stdout, "Hello PowerShell") {
		t.Errorf("Expected stdout to contain 'Hello PowerShell', got: %s", stdout)
	}

	if stderr != "" {
		t.Logf("Stderr (might be empty): %s", stderr)
	}
}

func TestLocalTransportExecutePowerShellNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping non-Windows PowerShell test on Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	// Test PowerShell command on non-Windows (should fail)
	command := "Write-Host 'Hello PowerShell'"
	_, _, err = transport.ExecutePowerShell(ctx, command)
	if err == nil {
		t.Error("Expected error for PowerShell on non-Windows platform, but got none")
	}

	expectedError := "PowerShell execution is only supported on Windows"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedError, err)
	}
}

func TestLocalTransportFileSystem(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	fs := transport.FileSystem()
	if fs == nil {
		t.Fatal("FileSystem returned nil")
	}

	// Test that it's actually a localFileSystem
	_, ok := fs.(*localFileSystem)
	if !ok {
		t.Error("FileSystem did not return a localFileSystem instance")
	}
}

func TestDetectPosixShell(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping POSIX shell test on Windows")
	}

	shell, args, err := detectPosixShell()
	if err != nil {
		t.Fatalf("detectPosixShell failed: %v", err)
	}
	if shell == "" {
		t.Fatal("detectPosixShell returned empty shell")
	}
	if len(args) == 0 {
		t.Fatal("detectPosixShell returned empty args")
	}

	expectedArgs := []string{"-c"}
	if len(args) != len(expectedArgs) || args[0] != expectedArgs[0] {
		t.Errorf("Expected args %v, got %v", expectedArgs, args)
	}

	if !strings.HasSuffix(shell, "sh") {
		t.Errorf("Expected shell to end with 'sh', got: %s", shell)
	}

	t.Logf("Detected shell: %s with args: %v", shell, args)
}

func TestLocalFileSystemConnect(t *testing.T) {
	fs := newLocalFileSystem()

	err := fs.Connect()
	if err != nil {
		t.Errorf("Connect() failed: %v", err)
	}
}

func TestLocalFileSystemClose(t *testing.T) {
	fs := newLocalFileSystem()

	err := fs.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}

func TestLocalFileSystemIsNull(t *testing.T) {
	fs := newLocalFileSystem()
	if fs.IsNull() {
		t.Error("Expected non-null FileSystem, got null")
	}
}

func TestLocalFileSystemStat(t *testing.T) {
	fs := newLocalFileSystem()

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_fs_stat_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test stat on the file
	info, err := fs.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.IsDir() {
		t.Error("Expected file to not be a directory")
	}

	expectedName := filepath.Base(tmpFile.Name())
	if info.Name() != expectedName {
		t.Errorf("Expected file name '%s', got '%s'", expectedName, info.Name())
	}
}

func TestLocalFileSystemStatNonExistent(t *testing.T) {
	fs := newLocalFileSystem()

	// Try to stat a non-existent file
	_, err := fs.Stat("/path/that/does/not/exist/file.txt")
	if err == nil {
		t.Error("Expected error when stating non-existent file, but got none")
	}

	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		t.Errorf("Expected 'file not found' error, got: %v", err)
	}
}

func TestLocalFileSystemStatDirectory(t *testing.T) {
	fs := newLocalFileSystem()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_fs_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test stat on the directory
	info, err := fs.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected directory to be identified as directory")
	}

	expectedName := filepath.Base(tmpDir)
	if info.Name() != expectedName {
		t.Errorf("Expected directory name '%s', got '%s'", expectedName, info.Name())
	}
}

func TestLocalFileSystemOpen(t *testing.T) {
	fs := newLocalFileSystem()

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_fs_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some test data
	testData := "test file content"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test opening the file
	file, err := fs.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Verify we can read from it
	buffer := make([]byte, len(testData))
	n, err := file.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from file: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to read %d bytes, got %d", len(testData), n)
	}

	if string(buffer) != testData {
		t.Errorf("Expected to read '%s', got '%s'", testData, string(buffer))
	}
}

func TestLocalFileSystemOpenNonExistent(t *testing.T) {
	fs := newLocalFileSystem()

	// Try to open a non-existent file
	_, err := fs.Open("/path/that/does/not/exist/file.txt")
	if err == nil {
		t.Error("Expected error when opening non-existent file, but got none")
	}

	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		t.Errorf("Expected 'file not found' error, got: %v", err)
	}
}
