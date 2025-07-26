package info

import (
	"errors"
	"io"
	"maps"
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestNewServiceManagerInfo(t *testing.T) {
	info := newServiceManagerInfo()
	if info == nil {
		t.Fatal("expected non-nil service manager info")
	}
	if info.name != "" {
		t.Errorf("expected empty name, got %q", info.name)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Windows(t *testing.T) {
	info := newServiceManagerInfo()

	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	transport := newMockTransport()
	fileSystem := transport.fileSystem

	err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.name != "windows-service-manager" {
		t.Errorf("expected name %q, got %q", "windows-service-manager", info.name)
	}
}

func TestServiceManagerInfo_Process1FromVirtualFile(t *testing.T) {
	tests := []struct {
		name          string
		commandOutput string
		expectedName  string
		shouldError   bool
		expectedError string
	}{
		{
			name:          "systemd",
			commandOutput: "systemd\n",
			expectedName:  "systemd",
		},
		{
			name:          "systemd with whitespace",
			commandOutput: "  systemd  \n",
			expectedName:  "systemd",
		},
		{
			name:          "runit-init (corrected)",
			commandOutput: "runit-init",
			expectedName:  "runit",
		},
		{
			name:          "openrc-init (corrected)",
			commandOutput: "openrc-init",
			expectedName:  "openrc",
		},
		{
			name:          "path with systemd",
			commandOutput: "/usr/lib/systemd/systemd",
			expectedName:  "systemd",
		},
		{
			name:          "COMMAND (imprecise)",
			commandOutput: "COMMAND",
			shouldError:   true,
			expectedError: "got imprecise or unexpected value in /proc/1/comm: COMMAND",
		},
		{
			name:          "init (imprecise)",
			commandOutput: "init",
			shouldError:   true,
			expectedError: "got imprecise or unexpected value in /proc/1/comm: init",
		},
		{
			name:          "shell ending",
			commandOutput: "bash",
			shouldError:   true,
			expectedError: "got imprecise or unexpected value in /proc/1/comm: bash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newServiceManagerInfo()

			transport := newMockTransport()
			transport.commandResponses["cat /proc/1/comm"] = &commandResponse{
				stdout: tt.commandOutput,
			}

			err := info.getProcess1FromVirtualFile(transport)

			if tt.shouldError {

				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
				}

			} else {

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if info.name != tt.expectedName {
					t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
				}
			}
		})
	}
}

func TestServiceManagerInfo_GetProcess1FromVirtualFile_CommandError(t *testing.T) {
	info := newServiceManagerInfo()
	transport := newMockTransport()
	transport.defaultCommandResponse = &commandResponse{
		err: io.EOF,
	}

	err := info.getProcess1FromVirtualFile(transport)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedError := "failed to read /proc/1/comm:"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Process1FromVirtualFile(t *testing.T) {
	tests := []struct {
		name          string
		commandOutput string
		expectedName  string
	}{
		{
			name:          "systemd",
			commandOutput: "systemd\n",
			expectedName:  "systemd",
		},
		{
			name:          "systemd with whitespace",
			commandOutput: "  systemd  \n",
			expectedName:  "systemd",
		},
		{
			name:          "runit-init (corrected)",
			commandOutput: "runit-init",
			expectedName:  "runit",
		},
		{
			name:          "openrc-init (corrected)",
			commandOutput: "openrc-init",
			expectedName:  "openrc",
		},
		{
			name:          "path with systemd",
			commandOutput: "/usr/lib/systemd/systemd",
			expectedName:  "systemd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			info := newServiceManagerInfo()

			transport := newMockTransport()
			transport.commandResponses["cat /proc/1/comm"] = &commandResponse{
				stdout: tt.commandOutput,
			}

			fileSystem := transport.fileSystem

			err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
			}
		})
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Process1FromInitLink(t *testing.T) {
	tests := []struct {
		name          string
		commandOutput string
		expectedName  string
	}{
		{
			name:          "systemd path",
			commandOutput: "/usr/lib/systemd/systemd\n",
			expectedName:  "systemd",
		},
		{
			name:          "runit-init (corrected)",
			commandOutput: "/sbin/runit-init\n",
			expectedName:  "runit",
		},
		{
			name:          "openrc-init (corrected)",
			commandOutput: "/sbin/openrc-init\n",
			expectedName:  "openrc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			info := newServiceManagerInfo()

			transport := newMockTransport()
			transport.commandResponses["cat /proc/1/comm"] = &commandResponse{
				stdout: "init\n",
			}
			transport.commandResponses["realpath /sbin/init"] = &commandResponse{
				stdout: tt.commandOutput,
			}

			fileSystem := transport.fileSystem

			err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
			}
		})
	}
}

func TestServiceManagerInfo_GetProcess1FromInitLink(t *testing.T) {
	tests := []struct {
		name          string
		commandOutput string
		expectedName  string
		shouldError   bool
		expectedError string
	}{
		{
			name:          "systemd path",
			commandOutput: "/usr/lib/systemd/systemd\n",
			expectedName:  "systemd",
		},
		{
			name:          "runit-init (corrected)",
			commandOutput: "/sbin/runit-init\n",
			expectedName:  "runit",
		},
		{
			name:          "openrc-init (corrected)",
			commandOutput: "/sbin/openrc-init\n",
			expectedName:  "openrc",
		},
		{
			name:          "init (imprecise)",
			commandOutput: "init\n",
			shouldError:   true,
			expectedError: "got imprecise or unexpected value in /sbin/init link: init",
		},
		{
			name:          "shell ending",
			commandOutput: "/bin/sh\n",
			shouldError:   true,
			expectedError: "got imprecise or unexpected value in /sbin/init link: /bin/sh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newServiceManagerInfo()

			transport := newMockTransport()
			transport.commandResponses["realpath /sbin/init"] = &commandResponse{
				stdout: tt.commandOutput,
			}

			err := info.getProcess1FromInitLink(transport)

			if tt.shouldError {

				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
				}

			} else {

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if info.name != tt.expectedName {
					t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
				}
			}
		})
	}
}

func TestServiceManagerInfo_GetProcess1FromInitLink_CommandError(t *testing.T) {
	info := newServiceManagerInfo()
	transport := newMockTransport()
	transport.defaultCommandResponse = &commandResponse{
		err: errors.New("command failed"),
	}

	err := info.getProcess1FromInitLink(transport)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedError := "failed to read /sbin/init link:"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_GetDarwinServiceManager(t *testing.T) {

	tests := []struct {
		name         string
		majorVersion string
		version      string
		expectedName string
	}{
		{
			name:         "macOS 11.0",
			majorVersion: "11",
			version:      "11.0",
			expectedName: "launchd",
		},
		{
			name:         "macOS 12.0",
			majorVersion: "12",
			version:      "12.0",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.15",
			majorVersion: "10",
			version:      "10.15",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.4",
			majorVersion: "10",
			version:      "10.4",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.3",
			majorVersion: "10",
			version:      "10.3",
			expectedName: "systemstarter",
		},
		{
			name:         "macOS 9.0",
			majorVersion: "9",
			version:      "9.0",
			expectedName: "systemstarter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.majorVersion = tt.majorVersion
			osInfo.version = tt.version
			osInfo.families.Add("darwin")

			info := newServiceManagerInfo()

			transport := newMockTransport()
			transport.commandResponses["cat /proc/1/comm"] = &commandResponse{
				err: errors.New("command failed"),
			}
			transport.commandResponses["realpath /sbin/init"] = &commandResponse{
				err: errors.New("command failed"),
			}

			fileSystem := transport.fileSystem

			err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
			}
		})
	}
}

func TestServiceManagerInfo_GetDarwinServiceManager(t *testing.T) {

	tests := []struct {
		name         string
		majorVersion string
		version      string
		expectedName string
		shouldError  bool
		errorMsg     string
	}{
		{
			name:         "macOS 11.0",
			majorVersion: "11",
			version:      "11.0",
			expectedName: "launchd",
		},
		{
			name:         "macOS 12.0",
			majorVersion: "12",
			version:      "12.0",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.15",
			majorVersion: "10",
			version:      "10.15",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.4",
			majorVersion: "10",
			version:      "10.4",
			expectedName: "launchd",
		},
		{
			name:         "macOS 10.3",
			majorVersion: "10",
			version:      "10.3",
			expectedName: "systemstarter",
		},
		{
			name:         "macOS 9.0",
			majorVersion: "9",
			version:      "9.0",
			expectedName: "systemstarter",
		},
		{
			name:         "invalid major version",
			majorVersion: "invalid",
			version:      "invalid",
			shouldError:  true,
			errorMsg:     "failed to parse macOS major version:",
		},
		{
			name:         "invalid version format",
			majorVersion: "10",
			version:      "10",
			shouldError:  true,
			errorMsg:     "failed to parse macOS version: 10",
		},
		{
			name:         "invalid minor version",
			majorVersion: "10",
			version:      "10.invalid",
			shouldError:  true,
			errorMsg:     "failed to parse macOS minor version:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newServiceManagerInfo()
			osInfo := newOSInfo()
			osInfo.majorVersion = tt.majorVersion
			osInfo.version = tt.version

			err := info.getDarwinServiceManager(osInfo)

			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if info.name != tt.expectedName {
					t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
				}
			}
		})
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_GetLinuxServiceManager(t *testing.T) {

	tests := []struct {
		name             string
		commandResponses map[string]*commandResponse
		fileSystemDirs   []string
		expectedName     string
	}{
		{
			name: "systemd with /run/systemd/system",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/run/systemd/system"},
			expectedName:   "systemd",
		},
		{
			name: "systemd with /dev/.run/systemd",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/dev/.run/systemd"},
			expectedName:   "systemd",
		},
		{
			name: "systemd with /dev/.systemd",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/dev/.systemd"},
			expectedName:   "systemd",
		},
		{
			name: "upstart",
			commandResponses: map[string]*commandResponse{
				"realpath initctl": {
					stdout: "/sbin/initctl",
				},
			},
			fileSystemDirs: []string{"/etc/init"},
			expectedName:   "upstart",
		},
		{
			name:           "openrc",
			fileSystemDirs: []string{"/sbin/openrc"},
			expectedName:   "openrc",
		},
		{
			name:           "sysvinit",
			fileSystemDirs: []string{"/etc/init.d"},
			expectedName:   "sysvinit",
		},
		{
			name:           "dinit",
			fileSystemDirs: []string{"/etc/dinit.d"},
			expectedName:   "dinit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("linux")

			info := newServiceManagerInfo()

			transport := newMockTransport()

			fileSystem := transport.fileSystem

			if len(tt.commandResponses) > 0 {
				maps.Copy(transport.commandResponses, tt.commandResponses)
			}

			transport.commandResponses["cat /proc/1/comm"] = &commandResponse{
				err: errors.New("command failed"),
			}
			transport.commandResponses["realpath /sbin/init"] = &commandResponse{
				err: errors.New("command failed"),
			}

			if len(tt.fileSystemDirs) > 0 {
				for _, dir := range tt.fileSystemDirs {
					fileSystem.dirs[dir] = &mockFileInfo{
						name:  dir,
						mode:  0755,
						isDir: true,
					}
				}
			}

			err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
			}
		})
	}
}

func TestServiceManagerInfo_GetLinuxServiceManager(t *testing.T) {

	tests := []struct {
		name             string
		commandResponses map[string]*commandResponse
		fileSystemDirs   []string
		expectedName     string
		expectedError    bool
	}{
		{
			name: "systemd with /run/systemd/system",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/run/systemd/system"},
			expectedName:   "systemd",
		},
		{
			name: "systemd with /dev/.run/systemd",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/dev/.run/systemd"},
			expectedName:   "systemd",
		},
		{
			name: "systemd with /dev/.systemd",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			fileSystemDirs: []string{"/dev/.systemd"},
			expectedName:   "systemd",
		},
		{
			name: "systemctl with no directories",
			commandResponses: map[string]*commandResponse{
				"realpath systemctl": {
					stdout: "/usr/bin/systemctl",
				},
			},
			expectedError: true,
		},
		{
			name: "upstart",
			commandResponses: map[string]*commandResponse{
				"realpath initctl": {
					stdout: "/sbin/initctl",
				},
			},
			fileSystemDirs: []string{"/etc/init"},
			expectedName:   "upstart",
		},
		{
			name:           "openrc",
			fileSystemDirs: []string{"/sbin/openrc"},
			expectedName:   "openrc",
		},
		{
			name:           "sysvinit",
			fileSystemDirs: []string{"/etc/init.d"},
			expectedName:   "sysvinit",
		},
		{
			name:           "dinit",
			fileSystemDirs: []string{"/etc/dinit.d"},
			expectedName:   "dinit",
		},
		{
			name:          "unknown service manager",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newServiceManagerInfo()

			transport := newMockTransport()

			fileSystem := transport.fileSystem

			if len(tt.commandResponses) > 0 {
				maps.Copy(transport.commandResponses, tt.commandResponses)
			}

			if len(tt.fileSystemDirs) > 0 {
				for _, dir := range tt.fileSystemDirs {
					fileSystem.dirs[dir] = &mockFileInfo{
						name:  dir,
						mode:  0755,
						isDir: true,
					}
				}
			}

			err := info.getLinuxServiceManager(transport, fileSystem)

			if tt.expectedError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !strings.Contains(err.Error(), "could not determine Linux service manager") {
					t.Errorf("expected error containing 'could not determine Linux service manager', got %q", err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.name != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.name)
			}
		})
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_UnsupportedOS(t *testing.T) {

	osInfo := newOSInfo()
	info := newServiceManagerInfo()

	transport := newMockTransport()
	transport.defaultCommandResponse = &commandResponse{
		err: errors.New("command failed"),
	}

	fileSystem := transport.fileSystem

	err := info.populateServiceManagerInfo(osInfo, transport, fileSystem)
	if err == nil {
		t.Fatal("expected error for unsupported OS, got nil")
	}

	expectedError := "could not determine service manager"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error containing %q, got %q", expectedError, err.Error())
	}
}

func TestServiceManagerInfo_ToMapOfCtyValues(t *testing.T) {
	tests := []struct {
		name         string
		serviceName  string
		expectedType cty.Type
		expectedVal  interface{}
	}{
		{
			name:         "with service manager name",
			serviceName:  "systemd",
			expectedType: cty.String,
			expectedVal:  "systemd",
		},
		{
			name:         "empty service manager name",
			serviceName:  "",
			expectedType: cty.String,
			expectedVal:  nil, // null value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &serviceManagerInfo{name: tt.serviceName}
			result := info.toMapOfCtyValues()

			if len(result) != 1 {
				t.Fatalf("expected 1 key in result, got %d", len(result))
			}

			value, exists := result["service_manager"]
			if !exists {
				t.Fatal("expected 'service_manager' key in result")
			}

			if !value.Type().Equals(tt.expectedType) {
				t.Errorf("expected type %s, got %s", tt.expectedType.FriendlyName(), value.Type().FriendlyName())
			}

			if tt.expectedVal == nil {
				if !value.IsNull() {
					t.Errorf("expected null value, got %s", value.AsString())
				}
			} else {
				if value.IsNull() {
					t.Errorf("expected non-null value, got null")
				}
				if value.AsString() != tt.expectedVal {
					t.Errorf("expected value %q, got %q", tt.expectedVal, value.AsString())
				}
			}
		})
	}
}
