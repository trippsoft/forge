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
		t.Errorf("expected no error from Connect(), got: %v", err)
	}
}

func TestNoneTransportClose(t *testing.T) {

	transport := TransportNone

	err := transport.Close()
	if err != nil {
		t.Errorf("expected no error from Close(), got: %v", err)
	}
}

func TestNoneTransport_Command(t *testing.T) {
	transport := TransportNone

	expectedError := "no transport available for command execution"
	cmd, err := transport.NewCommand("echo hello", nil)
	if err == nil {
		t.Fatalf("expected error %q from NewCommand(), got none", expectedError)
	}

	if cmd != nil {
		t.Errorf("expected nil command from NewCommand(), got %v", cmd)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from NewCommand(), got %q", expectedError, err.Error())
	}
}

func TestNoneTransport_PowerShell(t *testing.T) {
	transport := TransportNone

	expectedError := "no transport available for PowerShell execution"
	cmd, err := transport.NewPowerShellCommand("Write-Host 'hello'", nil)
	if err == nil {
		t.Fatalf("expected error %q from NewPowerShellCommand(), got none", expectedError)
	}

	if cmd != nil {
		t.Errorf("expected nil command from NewPowerShellCommand(), got %v", cmd)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from NewPowerShellCommand(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemStat(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	info, err := transport.Stat("/some/path/file.txt")
	if err == nil {
		t.Errorf("expected error %q from Stat(), got none", expectedError)
	}

	if info != nil {
		t.Errorf("expected nil file info from Stat(), got %v", info)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Stat(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemCreate(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	file, err := transport.Create("/some/path/file.txt")
	if err == nil {
		t.Errorf("expected error %q from Create(), got none", expectedError)
	}

	if file != nil {
		t.Errorf("expected nil file from Create(), got %v", file)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Create(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemOpen(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	file, err := transport.Open("/some/path/file.txt")
	if err == nil {
		t.Errorf("expected error %q from Open(), got none", expectedError)
	}

	if file != nil {
		t.Errorf("expected nil file from Open(), got %v", file)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Open(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemMkdir(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	err := transport.Mkdir("/some/path/directory")
	if err == nil {
		t.Errorf("expected error %q from Mkdir(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Mkdir(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemMkdirAll(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	err := transport.MkdirAll("/some/path/directory")
	if err == nil {
		t.Errorf("expected error %q from MkdirAll(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from MkdirAll(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemRemove(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	err := transport.Remove("/some/path/file.txt")
	if err == nil {
		t.Errorf("expected error %q from Remove(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Remove(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemRemoveAll(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	err := transport.RemoveAll("/some/path/directory")
	if err == nil {
		t.Errorf("expected error %q from RemoveAll(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from RemoveAll(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemJoin(t *testing.T) {
	transport := TransportNone

	path := transport.Join("some", "path", "file.txt")
	if path != "" {
		t.Errorf("expected empty path from Join on none file system, got %q", path)
	}
}

func TestNoneFileSystemTempDir(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	dir, err := transport.TempDir()
	if err == nil {
		t.Errorf("expected error %q from TempDir(), got none", expectedError)
	}

	if dir != "" {
		t.Errorf("expected empty temp dir from TempDir(), got %q", dir)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from TempDir(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemCreateTemp(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	file, err := transport.CreateTemp("", "tempfile")
	if err == nil {
		t.Errorf("expected error %q from CreateTemp(), got none", expectedError)
	}

	if file != nil {
		t.Errorf("expected nil file from CreateTemp(), got %v", file)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from CreateTemp(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemMkdirTemp(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	dir, err := transport.MkdirTemp("", "tempdir")
	if err == nil {
		t.Errorf("expected error %q from MkdirTemp(), got none", expectedError)
	}

	if dir != "" {
		t.Errorf("expected empty temp dir from MkdirTemp(), got %q", dir)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from MkdirTemp(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemSymlink(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	err := transport.Symlink("/target/path", "/link/path")
	if err == nil {
		t.Errorf("expected error %q from Symlink(), got none", expectedError)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from Symlink(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemReadLink(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	target, err := transport.ReadLink("/link/path")
	if err == nil {
		t.Errorf("expected error %q from ReadLink(), got none", expectedError)
	}

	if target != "" {
		t.Errorf("expected empty target from ReadLink(), got %q", target)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from ReadLink(), got %q", expectedError, err.Error())
	}
}

func TestNoneFileSystemRealPath(t *testing.T) {
	transport := TransportNone

	expectedError := "no file system available"
	realPath, err := transport.RealPath("/some/path")
	if err == nil {
		t.Errorf("expected error %q from RealPath(), got none", expectedError)
	}

	if realPath != "" {
		t.Errorf("expected empty real path from RealPath(), got %q", realPath)
	}

	if err.Error() != expectedError {
		t.Errorf("expected error %q from RealPath(), got %q", expectedError, err.Error())
	}
}
