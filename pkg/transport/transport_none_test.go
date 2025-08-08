// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"testing"
)

func TestNewNoneTransport(t *testing.T) {

	transport := TransportNone
	if transport == nil {
		t.Fatal("NewNoneTransport returned nil transport")
	}

	if transport.Type() != TransportTypeNone {
		t.Errorf("Expected transport type %s, got %s", TransportTypeNone, transport.Type())
	}
}

func TestNoneTransportType(t *testing.T) {

	transport := TransportNone

	if transport.Type() != TransportTypeNone {
		t.Errorf("Expected transport type %s, got %s", TransportTypeNone, transport.Type())
	}
}

func TestNoneTransportConnect(t *testing.T) {

	transport := TransportNone

	err := transport.Connect()
	if err != nil {
		t.Errorf("Connect should not return error for none transport, got: %v", err)
	}
}

func TestNoneTransportClose(t *testing.T) {

	transport := TransportNone

	err := transport.Close()
	if err != nil {
		t.Errorf("Close should not return error for none transport, got: %v", err)
	}
}

func TestNoneTransport_Command(t *testing.T) {

	transport := TransportNone

	cmd, err := transport.NewCommand("echo hello", nil)
	if err == nil {
		t.Fatal("Expected error for ExecuteCommand on none transport, but got none")
	}

	if cmd != nil {
		t.Error("Expected nil command for ExecuteCommand on none transport, but got a command")
	}

	expectedError := "no transport available for command execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneTransport_PowerShell(t *testing.T) {

	transport := TransportNone

	cmd, err := transport.NewPowerShellCommand("Write-Host 'hello'", nil)
	if err == nil {
		t.Fatal("Expected error for NewPowerShellCommand on none transport, but got none")
	}

	if cmd != nil {
		t.Error("Expected nil command for NewPowerShellCommand on none transport, but got a command")
	}

	expectedError := "no transport available for PowerShell execution"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemStat(t *testing.T) {

	transport := TransportNone

	// Test stat operation (should always fail)
	info, err := transport.Stat("/some/path/file.txt")
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

func TestNoneFileSystemCreate(t *testing.T) {

	transport := TransportNone

	// Test create operation (should always fail)
	file, err := transport.Create("/some/path/file.txt")
	if err == nil {
		t.Error("Expected error for Create on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if file != nil {
		t.Errorf("Expected nil file, got %v", file)
	}
}

func TestNoneFileSystemOpen(t *testing.T) {

	transport := TransportNone

	// Test open operation (should always fail)
	file, err := transport.Open("/some/path/file.txt")
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

func TestNoneFileSystemMkdir(t *testing.T) {

	transport := TransportNone

	err := transport.Mkdir("/some/path/directory")
	if err == nil {
		t.Error("Expected error for Mkdir on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemMkdirAll(t *testing.T) {

	transport := TransportNone

	err := transport.MkdirAll("/some/path/directory")
	if err == nil {
		t.Error("Expected error for MkdirAll on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemRemove(t *testing.T) {

	transport := TransportNone

	err := transport.Remove("/some/path/file.txt")
	if err == nil {
		t.Error("Expected error for Remove on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemRemoveAll(t *testing.T) {

	transport := TransportNone

	err := transport.RemoveAll("/some/path/directory")
	if err == nil {
		t.Error("Expected error for RemoveAll on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemJoin(t *testing.T) {

	transport := TransportNone

	path := transport.Join("some", "path", "file.txt")
	if path != "" {
		t.Errorf("Expected empty path from Join on none file system, got '%s'", path)
	}
}

func TestNoneFileSystemTempDir(t *testing.T) {

	transport := TransportNone

	dir, err := transport.TempDir()
	if err == nil {
		t.Error("Expected error for TempDir on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if dir != "" {
		t.Errorf("Expected empty temp dir, got '%s'", dir)
	}
}

func TestNoneFileSystemCreateTemp(t *testing.T) {

	transport := TransportNone

	file, err := transport.CreateTemp("", "tempfile")
	if err == nil {
		t.Error("Expected error for CreateTemp on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if file != nil {
		t.Errorf("Expected nil file from CreateTemp, got %v", file)
	}
}

func TestNoneFileSystemMkdirTemp(t *testing.T) {

	transport := TransportNone

	dir, err := transport.MkdirTemp("", "tempdir")
	if err == nil {
		t.Error("Expected error for MkdirTemp on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if dir != "" {
		t.Errorf("Expected empty temp dir from MkdirTemp, got '%s'", dir)
	}
}

func TestNoneFileSystemSymlink(t *testing.T) {

	transport := TransportNone

	err := transport.Symlink("/target/path", "/link/path")
	if err == nil {
		t.Error("Expected error for Symlink on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNoneFileSystemReadLink(t *testing.T) {

	transport := TransportNone

	target, err := transport.ReadLink("/link/path")
	if err == nil {
		t.Error("Expected error for ReadLink on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if target != "" {
		t.Errorf("Expected empty target from Readlink, got '%s'", target)
	}
}

func TestNoneFileSystemRealPath(t *testing.T) {

	transport := TransportNone

	realPath, err := transport.RealPath("/some/path")
	if err == nil {
		t.Error("Expected error for RealPath on none file system, but got none")
	}

	expectedError := "no file system available"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if realPath != "" {
		t.Errorf("Expected empty real path from RealPath, got '%s'", realPath)
	}
}
