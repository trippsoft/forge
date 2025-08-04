package mock

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/trippsoft/forge/internal/transport"
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

type FileInfo struct {
	FileName     string
	FileSize     int64
	FileMode     os.FileMode
	ModifiedTime time.Time
	IsDirectory  bool
	Target       string // For symlinks, if applicable
}

func (m *FileInfo) Name() string       { return m.FileName }
func (m *FileInfo) Size() int64        { return m.FileSize }
func (m *FileInfo) Mode() os.FileMode  { return m.FileMode }
func (m *FileInfo) ModTime() time.Time { return m.ModifiedTime }
func (m *FileInfo) IsDir() bool        { return m.IsDirectory }
func (m *FileInfo) Sys() any           { return nil }

type File struct {
	Info    *FileInfo
	Content []byte
	mutex   sync.RWMutex
	reader  *bytes.Reader
}

// Name implements transport.File.
func (m *File) Name() string {
	return m.Info.Name()
}

// Close implements transport.File.
func (m *File) Close() error {
	m.reader = nil
	return nil
}

// Read implements transport.File.
func (m *File) Read(p []byte) (n int, err error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.reader == nil {
		return 0, errors.New("file reader is not initialized")
	}

	return m.reader.Read(p)
}

// Sync implements transport.File.
func (m *File) Sync() error {
	return nil
}

// Write implements transport.File.
func (m *File) Write(p []byte) (n int, err error) {
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

type Transport struct {
	TransportType transport.TransportType

	CommandResults map[string]*CommandResult

	ErrorPaths map[string]error
	Files      map[string]*File
	Dirs       map[string]*FileInfo
}

func NewMockTransport() *Transport {
	return &Transport{
		TransportType:  transport.TransportTypeSSH,
		CommandResults: make(map[string]*CommandResult),
		ErrorPaths:     make(map[string]error),
		Files:          make(map[string]*File),
		Dirs:           make(map[string]*FileInfo),
	}
}

func (w *Transport) Type() transport.TransportType {
	return w.TransportType
}

func (w *Transport) Connect() error {
	return nil
}

func (w *Transport) Close() error {
	return nil
}

func (w *Transport) ExecuteCommand(ctx context.Context, command string) (string, string, error) {

	if result, exists := w.CommandResults[command]; exists {
		return result.Stdout, result.Stderr, result.Err
	}

	return "", "", fmt.Errorf("command not found in mock transport: %s", command)
}

func (w *Transport) ExecutePowerShell(ctx context.Context, command string) (string, error) {
	return "", errors.New("PowerShell execution not supported in mock transport")
}

// Stat implements transport.Transport.
func (w *Transport) Stat(path string) (os.FileInfo, error) {

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

// Create implements transport.Transport.
func (w *Transport) Create(path string) (transport.File, error) {

	if err, exists := w.ErrorPaths[path]; exists {
		return nil, err
	}

	if _, exists := w.Dirs[path]; exists {
		return nil, os.ErrExist
	}

	file := &File{
		Info: &FileInfo{
			FileName:     path,
			FileMode:     0644,
			ModifiedTime: time.Now(),
		},
		Content: nil,
	}
	w.Files[path] = file
	return file, nil
}

// Open implements transport.Transport.
func (w *Transport) Open(path string) (transport.File, error) {

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

// Mkdir implements transport.Transport.
func (w *Transport) Mkdir(path string) error {

	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	if _, exists := w.Files[path]; exists {
		return os.ErrExist
	}

	if _, exists := w.Dirs[path]; exists {
		return nil // Directory already exists
	}

	w.Dirs[path] = &FileInfo{
		FileName:     path,
		FileSize:     0,
		FileMode:     0755,
		ModifiedTime: time.Now(),
		IsDirectory:  true,
	}
	return nil
}

// MkdirAll implements transport.Transport.
func (w *Transport) MkdirAll(path string) error {
	return w.Mkdir(path) // For mock, we treat MkdirAll the same as Mkdir
}

// Remove implements transport.Transport.
func (w *Transport) Remove(path string) error {

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

// RemoveAll implements transport.Transport.
func (w *Transport) RemoveAll(path string) error {

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

// Join implements transport.Transport.
func (w *Transport) Join(elem ...string) string {

	stringBuilder := &strings.Builder{}
	for i, e := range elem {
		if i > 0 {
			stringBuilder.WriteString("/")
		}
		stringBuilder.WriteString(strings.Trim(e, "/"))
	}

	return stringBuilder.String()
}

// TempDir implements transport.Transport.
func (w *Transport) TempDir() (string, error) {
	return "/tmp", nil
}

// CreateTemp implements transport.Transport.
func (w *Transport) CreateTemp(dir string, pattern string) (transport.File, error) {

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

// MkdirTemp implements transport.Transport.
func (w *Transport) MkdirTemp(dir string, pattern string) (string, error) {

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

// Symlink implements transport.Transport.
func (w *Transport) Symlink(target string, path string) error {

	if err, exists := w.ErrorPaths[path]; exists {
		return err
	}

	if _, exists := w.Files[path]; exists {
		return os.ErrExist // Cannot create symlink to an existing file
	}

	if _, exists := w.Dirs[path]; exists {
		return os.ErrExist // Cannot create symlink to an existing directory
	}

	w.Files[path] = &File{
		Info: &FileInfo{
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

// ReadLink implements transport.Transport.
func (w *Transport) ReadLink(path string) (string, error) {
	if file, exists := w.Files[path]; exists {
		if file.Info.Target != "" {
			return file.Info.Target, nil // Return the target of the symlink
		}
		return "", os.ErrInvalid // Not a symlink
	}
	return "", os.ErrNotExist
}

// RealPath implements transport.Transport.
func (w *Transport) RealPath(path string) (string, error) {

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
