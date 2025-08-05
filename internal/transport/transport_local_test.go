package transport

import (
	"bytes"
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
}

func TestNewLocalTransport_Windows(t *testing.T) {

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

func TestNewLocalTransport_NonWindows(t *testing.T) {

	if runtime.GOOS == "windows" {
		t.Skip("Skipping non-Windows test on Windows platform")
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

	cmd := transport.NewCommand(command)

	var outBuf, errBuf bytes.Buffer

	err = cmd.SetStdout(&outBuf)
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	err = cmd.SetStderr(&errBuf)
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	err = cmd.Run(ctx)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}

	stdout := strings.TrimSpace(outBuf.String())
	if !strings.Contains(stdout, "hello") {
		t.Errorf("Expected stdout to contain 'hello', got: %s", stdout)
	}

	if err != nil {
		t.Fatalf("Failed to read stderr: %v", err)
	}

	stderr := strings.TrimSpace(errBuf.String())
	if stderr != "" {
		t.Logf("Stderr (might be empty): %s", stderr)
	}
}

func TestLocalTransportExecuteCommand_Timeout(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var command string
	if runtime.GOOS == "windows" {
		command = "timeout 1"
	} else {
		command = "sleep 1"
	}

	cmd := transport.NewCommand(command)

	err = cmd.Run(ctx)
	if err == nil {
		t.Error("Expected error due to context timeout, but got none")
	}
}

func TestLocalTransportExecuteCommand_Error(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	cmd := transport.NewCommand("nonexistentcommand12345")

	err = cmd.Run(ctx)
	if err == nil {
		t.Error("Expected error for nonexistent command, but got none")
	}
}

func TestLocalTransportExecutePowerShell_Windows(t *testing.T) {

	if runtime.GOOS != "windows" {
		t.Skip("Skipping PowerShell test on non-Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	ctx := context.Background()

	cmd, err := transport.NewPowerShellCommand("Write-Host 'Hello PowerShell'")
	if err != nil {
		t.Fatalf("NewPowerShellCommand failed: %v", err)
	}

	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	err = cmd.Run(ctx)
	if err != nil {
		t.Fatalf("ExecutePowerShell failed: %v", err)
	}

	stdout := strings.TrimSpace(outBuf.String())
	if stdout != "Hello PowerShell" {
		t.Errorf("Expected stdout to contain 'Hello PowerShell', got: %s", stdout)
	}
}

func TestLocalTransportExecutePowerShell_NonWindows(t *testing.T) {

	if runtime.GOOS == "windows" {
		t.Skip("Skipping non-Windows PowerShell test on Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Test PowerShell command on non-Windows (should fail)
	cmd, err := transport.NewPowerShellCommand("Write-Host 'Hello PowerShell'")
	if err == nil {
		t.Fatal("Expected error for NewPowerShellCommand on non-Windows platform, but got a command")
	}

	if cmd != nil {
		t.Fatal("Expected nil command for NewPowerShellCommand on non-Windows platform, but got a command")
	}

	expectedErr := "PowerShell execution is only supported on Windows"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedErr, err)
	}
}

func TestLocalTransportStat_File(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_stat_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test stat on the file
	fileInfo, err := transport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.IsDir() {
		t.Error("Expected file to not be a directory")
	}

	expectedName := filepath.Base(tmpFile.Name())
	if fileInfo.Name() != expectedName {
		t.Errorf("Expected file name '%s', got '%s'", expectedName, fileInfo.Name())
	}
}

func TestLocalTransportStat_Directory(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test stat on the directory
	info, err := transport.Stat(tmpDir)
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

func TestLocalTransportStat_NonExistent(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to stat a non-existent file
	fileInfo, err := transport.Stat("/path/that/does/not/exist/file.txt")
	if err != nil {
		t.Errorf("Expected error when stating non-existent file, got: %v", err)
	}

	if fileInfo != nil {
		t.Error("Expected nil file info when stating non-existent file, got: ", fileInfo)
	}
}

func TestLocalTransportCreate(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_create_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("Initial content")
	tmpFile.Sync()
	tmpFile.Close()

	// Test creating a new file
	file, err := transport.Create(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Write([]byte("New content"))
	file.Sync()
	file.Close()

	if file == nil {
		t.Error("Expected non-nil file after Create, got nil")
	}

	file, err = transport.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open created file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from created file: %v", err)
	}

	if n != len("New content") {
		t.Errorf("Expected to read %d bytes, got %d", len("New content"), n)
	}

	if string(buffer[:n]) != "New content" {
		t.Errorf("Expected to read 'New content', got '%s'", string(buffer[:n]))
	}
}

func TestLocalTransportOpen(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some test data
	testData := "test file content"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test opening the file
	file, err := transport.Open(tmpFile.Name())
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

func TestLocalTransportOpen_NonExistent(t *testing.T) {
	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to open a non-existent file
	file, err := transport.Open("/path/that/does/not/exist/file.txt")
	if err != nil {
		t.Errorf("Expected error when opening non-existent file, got: %v", err)
	}

	if file != nil {
		t.Error("Expected nil file when opening non-existent file, got: ", file)
	}
}

func TestLocalTransportMkdir(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_mkdir_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating a new directory
	newDir := filepath.Join(tmpDir, "newdir")
	err = transport.Mkdir(newDir)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	dir, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("Failed to stat created directory: %v", err)
	}

	if !dir.IsDir() {
		t.Errorf("Expected '%s' to be a directory, but it is not", newDir)
	}
}

func TestLocalTransportMkdir_Existing(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_mkdir_existing_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating an existing directory (should not return an error)
	err = transport.Mkdir(tmpDir)
	if err != nil {
		t.Errorf("Expected no error when creating existing directory, got: %v", err)
	}
}

func TestLocalTransportMkdirAll(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_mkdirall_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating a nested directory structure
	nestedDir := filepath.Join(tmpDir, "nested", "dir")
	err = transport.MkdirAll(nestedDir)
	if err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	dir, err := os.Stat(nestedDir)
	if err != nil {
		t.Fatalf("Failed to stat created nested directory: %v", err)
	}

	if !dir.IsDir() {
		t.Errorf("Expected '%s' to be a directory, but it is not", nestedDir)
	}
}

func TestLocalTransportMkdirAll_Existing(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_mkdirall_existing_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating an existing nested directory (should not return an error)
	err = transport.MkdirAll(tmpDir)
	if err != nil {
		t.Errorf("Expected no error when creating existing directory, got: %v", err)
	}
}

func TestLocalTransportRemove(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_remove_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test removing the file
	err = transport.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}

	file, err := transport.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat removed file: %v", err)
	}

	if file != nil {
		t.Errorf("Expected file to be removed, but it still exists: %v", file)
	}
}

func TestLocalTransportRemove_NonExistent(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to remove a non-existent file
	err = transport.Remove("/path/that/does/not/exist/file.txt")
	if err == nil {
		t.Error("Expected error when removing non-existent file, but got none")
	}

	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		t.Errorf("Expected error to be os.ErrNotExist or syscall.ENOENT, got: %v", err)
	}

	file, err := transport.Stat("/path/that/does/not/exist/file.txt")
	if err != nil {
		t.Errorf("Failed to stat non-existent file: %v", err)
	}

	if file != nil {
		t.Errorf("Expected file to be non-existent, but it still exists: %v", file)
	}
}

func TestLocalTransportRemoveAll(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_local_transport_removeall_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file inside the temporary directory
	tmpFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(tmpFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test removing the directory and its contents
	err = transport.RemoveAll(tmpDir)
	if err != nil {
		t.Fatalf("Failed to remove directory: %v", err)
	}

	fileInfo, err := transport.Stat(tmpDir)
	if err != nil {
		t.Errorf("Failed to stat removed directory: %v", err)
	}

	if fileInfo != nil {
		t.Errorf("Expected directory to be removed, but it still exists")
	}
}

func TestLocalTransportRemoveAll_NonExistent(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to remove a non-existent directory
	err = transport.RemoveAll("/path/that/does/not/exist/directory")
	if err != nil {
		t.Errorf("Expected no error when removing non-existent directory, got: %v", err)
	}

	fileInfo, err := transport.Stat("/path/that/does/not/exist/directory")
	if err != nil {
		t.Errorf("Failed to stat non-existent directory: %v", err)
	}

	if fileInfo != nil {
		t.Errorf("Expected directory to be non-existent, but it still exists: %v", fileInfo)
	}
}

func TestLocalTransportJoin(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Test joining paths
	path1 := "/some/directory"
	path2 := "subdirectory"
	expectedPath := filepath.Join(path1, path2)

	resultPath := transport.Join(path1, path2)
	if resultPath != expectedPath {
		t.Errorf("Expected joined path '%s', got '%s'", expectedPath, resultPath)
	}
}

func TestLocalTransportTempDir(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	tempDir, err := transport.TempDir()
	if err != nil {
		t.Fatalf("TempDir failed: %v", err)
	}

	if tempDir != os.TempDir() {
		t.Errorf("Expected temporary directory '%s', got '%s'", os.TempDir(), tempDir)
	}
}

func TestLocalTransportCreateTemp(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	tempFile, err := transport.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if !strings.HasPrefix(tempFile.Name(), os.TempDir()) {
		t.Errorf("Expected temporary file to be in '%s', got '%s'", os.TempDir(), tempFile.Name())
	}
}

func TestLocalTransportMkdirTemp(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	tempDir, err := transport.MkdirTemp("", "testdir_*")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if !strings.HasPrefix(tempDir, os.TempDir()) {
		t.Errorf("Expected temporary directory to be in '%s', got '%s'", os.TempDir(), tempDir)
	}
}

func TestLocalTransportSymlink(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_symlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a symlink to the temporary file
	symlinkPath := tmpFile.Name() + "_symlink"
	err = transport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}
	defer os.Remove(symlinkPath)

	// Verify the symlink points to the correct target
	target, err := transport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestLocalTransportReadLink(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_local_transport_readlink_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a symlink to the temporary file
	symlinkPath := tmpFile.Name() + "_symlink"
	err = transport.Symlink(tmpFile.Name(), symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}
	defer os.Remove(symlinkPath)

	// Read the symlink
	target, err := transport.ReadLink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}

	if target != tmpFile.Name() {
		t.Errorf("Expected symlink target '%s', got '%s'", tmpFile.Name(), target)
	}
}

func TestLocalTransportReadLink_NonExistent(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to read a non-existent symlink
	_, err = transport.ReadLink("/path/that/does/not/exist/symlink")
	if err == nil {
		t.Error("Expected error when reading non-existent symlink, but got none")
	}

	expectedError := "no such file or directory"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedError, err)
	}
}

func TestLocalTransportRealPath_Posix(t *testing.T) {

	if runtime.GOOS != "posix" {
		t.Skip("Skipping RealPath test on non-POSIX platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	path, err := transport.RealPath("sh")
	if err != nil {
		t.Fatalf("RealPath failed: %v", err)
	}

	if path != "/bin/sh" && path != "/usr/bin/sh" {
		t.Errorf("Expected RealPath to return '/bin/sh' or '/usr/bin/sh', got '%s'", path)
	}
}

func TestLocalTransportRealPath_NonExistent(t *testing.T) {

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	// Try to get the real path of a non-existent file
	path, err := transport.RealPath("nonexistentfile.txt")
	if err == nil {
		t.Fatal("Expected errors, got none")
	}

	_ = path // Ignore the path variable, we expect an error

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error to be os.ErrNotExist, got: %v", err)
	}
}

func TestLocalTransportRealPath_Windows(t *testing.T) {

	if runtime.GOOS != "windows" {
		t.Skip("Skipping RealPath test on non-Windows platform")
	}

	transport, err := NewLocalTransport()
	if err != nil {
		t.Fatalf("NewLocalTransport failed: %v", err)
	}

	path, err := transport.RealPath("cmd.exe")
	if err != nil {
		t.Fatalf("RealPath failed: %v", err)
	}

	if path != "C:\\Windows\\System32\\cmd.exe" && path != "C:\\Windows\\SysWOW64\\cmd.exe" {
		t.Errorf("Expected RealPath to return 'C:\\Windows\\System32\\cmd.exe' or 'C:\\Windows\\SysWOW64\\cmd.exe', got '%s'", path)
	}
}
