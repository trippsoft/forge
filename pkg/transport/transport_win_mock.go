// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	winPathPrefixes = []string{
		"C:\\Windows\\System32\\",
		"C:\\Windows\\",
		"C:\\Windows\\System32\\Wbem\\",
		"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\",
	}
)

type MockWinTransport struct {
	TransportType TransportType

	CommandResults    map[string]*MockCmd
	PowerShellResults map[string]*MockCmd

	ErrorPaths map[string]error
	Files      map[string]*MockFile
	Dirs       map[string]*MockFileInfo
}

func NewWinMockTransport() *MockWinTransport {
	return &MockWinTransport{
		TransportType:     TransportTypeSSH,
		CommandResults:    make(map[string]*MockCmd),
		PowerShellResults: make(map[string]*MockCmd),
	}
}

func (w *MockWinTransport) Type() TransportType {
	return w.TransportType
}

func (w *MockWinTransport) Connect() error {
	return nil
}

func (w *MockWinTransport) Close() error {
	return nil
}

func (w *MockWinTransport) NewCommand(command string, escalateConfig Escalation) (Cmd, error) {
	if cmd, exists := w.CommandResults[command]; exists {
		cmd.completed = false // Reset completion status for reuse
		cmd.stdin = nil       // Reset stdin for new command execution
		return cmd, nil
	}

	return &MockCmd{
		Err: fmt.Errorf("command not found in mock transport: %s", command),
	}, nil
}

func (w *MockWinTransport) NewPowerShellCommand(command string, escalateConfig Escalation) (Cmd, error) {
	if cmd, exists := w.PowerShellResults[command]; exists {
		cmd.completed = false // Reset completion status for reuse
		cmd.stdin = nil       // Reset stdin for new command execution
		return cmd, nil
	}

	return &MockCmd{
		Err: fmt.Errorf("PowerShell command not found in mock transport: %s", command),
	}, nil
}

func (w *MockWinTransport) NewPythonCommand(interpreter, command string, escalateConfig Escalation) (Cmd, error) {
	return nil, errors.New("Python execution not supported in mock transport")
}

// Stat implements Transport.
func (w *MockWinTransport) Stat(path string) (os.FileInfo, error) {
	if err, exists := w.ErrorPaths[path]; exists {
		return nil, err
	}

	if dir, exists := w.Dirs[path]; exists {
		return dir, nil
	}

	if file, exists := w.Files[path]; exists {
		if file.Info.Target != "" {
			return w.Stat(file.Info.Target) // Follow symlink if it exists
		}
		return file.Info, nil
	}

	return nil, os.ErrNotExist
}

// Create implements Transport.
func (w *MockWinTransport) Create(path string) (File, error) {
	if err, exists := w.ErrorPaths[path]; exists {
		return nil, err
	}

	if _, exists := w.Dirs[path]; exists {
		return nil, os.ErrExist
	}

	file := &MockFile{
		Info: &MockFileInfo{
			FileName:     path,
			FileSize:     0,
			FileMode:     0644,
			ModifiedTime: time.Now(),
			IsDirectory:  false,
		},
		Content: nil,
	}
	w.Files[path] = file
	return file, nil
}

// Open implements Transport.
func (w *MockWinTransport) Open(path string) (File, error) {
	if err, exists := w.ErrorPaths[path]; exists {
		return nil, err
	}

	if _, exists := w.Dirs[path]; exists {
		return nil, os.ErrExist
	}

	if file, exists := w.Files[path]; exists {
		if file.Info.Target != "" {
			return w.Open(file.Info.Target) // Follow symlink if it exists
		}
		return file, nil
	}

	return nil, os.ErrNotExist
}

// Mkdir implements Transport.
func (w *MockWinTransport) Mkdir(path string) error {
	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	if _, exists := w.Files[path]; exists {
		return os.ErrExist
	}

	if _, exists := w.Dirs[path]; exists {
		return nil // Directory already exists
	}

	w.Dirs[path] = &MockFileInfo{
		FileName:     path,
		FileSize:     0,
		FileMode:     0755,
		ModifiedTime: time.Now(),
		IsDirectory:  true,
	}

	return nil
}

// MkdirAll implements Transport.
func (w *MockWinTransport) MkdirAll(path string) error {
	return w.Mkdir(path) // For mock, we treat MkdirAll the same as Mkdir
}

// Remove implements Transport.
func (w *MockWinTransport) Remove(path string) error {
	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	_, dirExists := w.Dirs[path]

	if dirExists {
		for filePath := range w.Files {
			if strings.HasPrefix(filePath, path+"\\") {
				return os.ErrInvalid // Cannot remove directory with files inside
			}
		}

		for dirPath := range w.Dirs {
			if strings.HasPrefix(dirPath, path+"\\") {
				return os.ErrInvalid // Cannot remove directory with subdirectories
			}
		}

		delete(w.Dirs, path)
		return nil
	}

	if file, exists := w.Files[path]; exists {
		delete(w.Files, path)
		file.Content = nil // Clear content on remove
		return nil
	}

	return os.ErrNotExist
}

// RemoveAll implements Transport.
func (w *MockWinTransport) RemoveAll(path string) error {
	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	_, dirExists := w.Dirs[path]

	if dirExists {
		toDelete := make([]string, 0)
		for filePath := range w.Files {
			if strings.HasPrefix(filePath, path+"\\") {
				toDelete = append(toDelete, filePath)
			}
		}
		for _, filePath := range toDelete {
			delete(w.Files, filePath)
		}

		toDelete = make([]string, 0)
		for dirPath := range w.Dirs {
			if strings.HasPrefix(dirPath, path+"\\") {
				toDelete = append(toDelete, dirPath)
			}
		}
		for _, dirPath := range toDelete {
			delete(w.Dirs, dirPath)
		}
		delete(w.Dirs, path)

		return nil
	}

	if file, exists := w.Files[path]; exists {
		delete(w.Files, path)
		file.Content = nil // Clear content on remove
	}

	return nil
}

// Join implements Transport.
func (w *MockWinTransport) Join(elem ...string) string {
	stringBuilder := &strings.Builder{}
	for i, e := range elem {
		if i > 0 {
			stringBuilder.WriteString("\\")
		}
		stringBuilder.WriteString(strings.Trim(e, "\\"))
	}

	return stringBuilder.String()
}

// TempDir implements Transport.
func (w *MockWinTransport) TempDir() (string, error) {
	return "C:\\Users\\mock\\AppData\\Local\\Temp", nil
}

// CreateTemp implements Transport.
func (w *MockWinTransport) CreateTemp(dir string, pattern string) (File, error) {
	if dir == "" {
		dir, _ = w.TempDir()
	}

	_ = w.Mkdir(dir)

	splitPattern := strings.Split(pattern, "*")
	if len(splitPattern) > 2 {
		return nil, fmt.Errorf("pattern must contain at most one wildcard (*)")
	}

	var prefix, suffix string
	if len(splitPattern) == 1 {
		prefix = splitPattern[0]
	} else {
		prefix = splitPattern[0]
		suffix = splitPattern[1]
	}

	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(dir)
	stringBuilder.WriteRune('\\')
	stringBuilder.WriteString(prefix)

	randomNumber := fmt.Sprintf("%d", time.Now().UnixNano()%1000000) // Simple random number based on current time
	stringBuilder.WriteString(randomNumber)
	stringBuilder.WriteString(suffix)

	return w.Create(stringBuilder.String())
}

// MkdirTemp implements Transport.
func (w *MockWinTransport) MkdirTemp(dir string, pattern string) (string, error) {
	if dir == "" {
		dir, _ = w.TempDir()
	}

	_ = w.Mkdir(dir)

	splitPattern := strings.Split(pattern, "*")
	if len(splitPattern) > 2 {
		return "", fmt.Errorf("pattern must contain at most one wildcard (*)")
	}

	var prefix, suffix string
	if len(splitPattern) == 1 {
		prefix = splitPattern[0]
	} else {
		prefix = splitPattern[0]
		suffix = splitPattern[1]
	}

	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString(dir)
	stringBuilder.WriteRune('\\')
	stringBuilder.WriteString(prefix)

	randomNumber := fmt.Sprintf("%d", time.Now().UnixNano()%1000000) // Simple random number based on current time
	stringBuilder.WriteString(randomNumber)
	stringBuilder.WriteString(suffix)

	err := w.Mkdir(stringBuilder.String())
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}

// Symlink implements Transport.
func (w *MockWinTransport) Symlink(target string, path string) error {
	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	if _, exists := w.Files[path]; exists {
		return os.ErrExist // Cannot create symlink to an existing file
	}

	if _, exists := w.Dirs[path]; exists {
		return os.ErrExist // Cannot create symlink to an existing directory
	}

	w.Files[path] = &MockFile{
		Info: &MockFileInfo{
			FileName:     path,
			FileSize:     0,
			FileMode:     0777, // Symlinks are typically executable
			ModifiedTime: time.Now(),
			IsDirectory:  false,
			Target:       target,
		},
		Content: nil,
	}

	return nil
}

// ReadLink implements Transport.
func (w *MockWinTransport) ReadLink(path string) (string, error) {
	if file, exists := w.Files[path]; exists {
		if file.Info.Target != "" {
			return file.Info.Target, nil // Return the target of the symlink
		}
		return "", os.ErrInvalid // Not a symlink
	}
	return "", os.ErrNotExist
}

// RealPath implements Transport.
func (w *MockWinTransport) RealPath(path string) (string, error) {
	if _, exists := w.Files[path]; exists {
		return path, nil // Return the path as is for mock transport
	}

	if _, exists := w.Dirs[path]; exists {
		return path, nil // Return the path as is for mock transport
	}

	for _, prefix := range winPathPrefixes {
		newPath := prefix + path
		if err, exists := w.ErrorPaths[newPath]; exists {
			return "", err // Return error if path is in error map
		}
		if _, exists := w.Files[newPath]; exists {
			return newPath, nil // Return the first matching path
		}
		if _, exists := w.Dirs[newPath]; exists {
			return newPath, nil // Return the first matching directory
		}
	}

	return "", os.ErrNotExist // No matching file or directory found
}
