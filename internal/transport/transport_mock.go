package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	posixPathPrefixes = []string{
		"/usr/local/bin/",
		"/usr/bin/",
		"/bin/",
		"/usr/local/sbin/",
		"/usr/sbin/",
		"/sbin/",
	}
)

type MockFileInfo struct {
	FileName     string
	FileSize     int64
	FileMode     os.FileMode
	ModifiedTime time.Time
	IsDirectory  bool
	Target       string // For symlinks, if applicable
}

func (m *MockFileInfo) Name() string       { return m.FileName }
func (m *MockFileInfo) Size() int64        { return m.FileSize }
func (m *MockFileInfo) Mode() os.FileMode  { return m.FileMode }
func (m *MockFileInfo) ModTime() time.Time { return m.ModifiedTime }
func (m *MockFileInfo) IsDir() bool        { return m.IsDirectory }
func (m *MockFileInfo) Sys() any           { return nil }

type MockFile struct {
	Info    *MockFileInfo
	Content []byte
	mutex   sync.RWMutex
	reader  *bytes.Reader
}

// Name implements File.
func (m *MockFile) Name() string {
	return m.Info.Name()
}

// Close implements File.
func (m *MockFile) Close() error {
	m.reader = nil
	return nil
}

// Read implements File.
func (m *MockFile) Read(p []byte) (n int, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.reader == nil {
		return 0, errors.New("file reader is not initialized")
	}

	return m.reader.Read(p)
}

// Sync implements File.
func (m *MockFile) Sync() error {
	return nil
}

// Write implements File.
func (m *MockFile) Write(p []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.Content == nil {
		m.Content = make([]byte, 0)
	}

	m.Content = append(m.Content, p...)
	n = len(p)
	m.Info.ModifiedTime = time.Now()
	m.Info.FileSize += int64(n)
	return n, nil
}

type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

type MockTransport struct {
	TransportType TransportType

	CommandResults map[string]*CommandResult

	ErrorPaths map[string]error
	Files      map[string]*MockFile
	Dirs       map[string]*MockFileInfo
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		TransportType:  TransportTypeSSH,
		CommandResults: make(map[string]*CommandResult),
		ErrorPaths:     make(map[string]error),
		Files:          make(map[string]*MockFile),
		Dirs:           make(map[string]*MockFileInfo),
	}
}

func (w *MockTransport) Type() TransportType {
	return w.TransportType
}

func (w *MockTransport) Connect() error {
	return nil
}

func (w *MockTransport) Close() error {
	return nil
}

func (w *MockTransport) NewCommand(command string) *Cmd {
	return NewCmd(w, command)
}

func (w *MockTransport) NewPowerShellCommand(command string) *PowerShellCmd {
	return NewPowerShellCmd(w, command)
}

func (w *MockTransport) executeCommand(ctx context.Context, cmd *Cmd) error {

	if result, exists := w.CommandResults[cmd.command]; exists {
		cmd.Stdout.Write([]byte(result.Stdout))
		cmd.Stderr.Write([]byte(result.Stderr))
		return result.Err
	}

	return fmt.Errorf("command not found in mock transport: %s", cmd.command)
}

func (w *MockTransport) executePowerShell(ctx context.Context, cmd *PowerShellCmd) error {
	return errors.New("PowerShell execution not supported in mock transport")
}

// Stat implements Transport.
func (w *MockTransport) Stat(path string) (os.FileInfo, error) {

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

	return nil, nil
}

// Create implements Transport.
func (w *MockTransport) Create(path string) (File, error) {

	if err, exists := w.ErrorPaths[path]; exists {
		return nil, err
	}

	if _, exists := w.Dirs[path]; exists {
		return nil, os.ErrExist
	}

	file := &MockFile{
		Info: &MockFileInfo{
			FileName:     path,
			FileMode:     0644,
			ModifiedTime: time.Now(),
		},
		Content: nil,
	}
	w.Files[path] = file
	return file, nil
}

// Open implements Transport.
func (w *MockTransport) Open(path string) (File, error) {

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
		file.reader = bytes.NewReader(file.Content) // Initialize reader for file content
		return file, nil
	}

	return nil, nil
}

// Mkdir implements Transport.
func (w *MockTransport) Mkdir(path string) error {

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
func (w *MockTransport) MkdirAll(path string) error {
	return w.Mkdir(path) // For mock, we treat MkdirAll the same as Mkdir
}

// Remove implements Transport.
func (w *MockTransport) Remove(path string) error {

	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	_, dirExists := w.Dirs[path]

	if dirExists {
		for filePath := range w.Files {
			if strings.HasPrefix(filePath, path+"/") {
				return os.ErrInvalid // Cannot remove directory with files inside
			}
		}

		for dirPath := range w.Dirs {
			if strings.HasPrefix(dirPath, path+"/") {
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
func (w *MockTransport) RemoveAll(path string) error {

	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	_, dirExists := w.Dirs[path]

	if dirExists {
		toDelete := make([]string, 0)
		for filePath := range w.Files {
			if strings.HasPrefix(filePath, path+"/") {
				toDelete = append(toDelete, filePath)
			}
		}
		for _, filePath := range toDelete {
			delete(w.Files, filePath)
		}

		toDelete = make([]string, 0)
		for dirPath := range w.Dirs {
			if strings.HasPrefix(dirPath, path+"/") {
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
func (w *MockTransport) Join(elem ...string) string {

	stringBuilder := &strings.Builder{}
	for i, e := range elem {
		if i > 0 {
			stringBuilder.WriteString("/")
		}
		stringBuilder.WriteString(strings.Trim(e, "/"))
	}

	return stringBuilder.String()
}

// TempDir implements Transport.
func (w *MockTransport) TempDir() (string, error) {
	return "/tmp", nil
}

// CreateTemp implements Transport.
func (w *MockTransport) CreateTemp(dir string, pattern string) (File, error) {

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
	stringBuilder.WriteRune('/')
	stringBuilder.WriteString(prefix)

	randomNumber := fmt.Sprintf("%d", time.Now().UnixNano()%1000000) // Simple random number based on current time
	stringBuilder.WriteString(randomNumber)
	stringBuilder.WriteString(suffix)

	return w.Create(stringBuilder.String())
}

// MkdirTemp implements Transport.
func (w *MockTransport) MkdirTemp(dir string, pattern string) (string, error) {

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
	stringBuilder.WriteRune('/')
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
func (w *MockTransport) Symlink(target string, path string) error {

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
func (w *MockTransport) ReadLink(path string) (string, error) {
	if file, exists := w.Files[path]; exists {
		if file.Info.Target != "" {
			return file.Info.Target, nil // Return the target of the symlink
		}
		return "", os.ErrInvalid // Not a symlink
	}
	return "", os.ErrNotExist
}

// RealPath implements Transport.
func (w *MockTransport) RealPath(path string) (string, error) {

	if _, exists := w.Files[path]; exists {
		return path, nil // Return the path as is for mock transport
	}

	if _, exists := w.Dirs[path]; exists {
		return path, nil // Return the path as is for mock transport
	}

	for _, prefix := range posixPathPrefixes {
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
