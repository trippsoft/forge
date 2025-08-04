package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/internal/transport/mock"
)

func TestUserInfo_PopulateUserInfo_NoOS(t *testing.T) {

	osInfo := newOSInfo()

	transport := mock.NewMockTransport()

	info := newUserInfo()
	diags := info.populateUserInfo(osInfo, transport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatal("expected warnings, got none")
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Invalid OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping user information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got %q", expectedDetail, warnings[0].Detail)
	}
}

func TestUserInfo_PopulateUserInfo_Posix(t *testing.T) {

	tests := []struct {
		name            string
		output          string
		expectedName    string
		expectedId      string
		expectedGroupId string
		expectedHomeDir string
		expectedShell   string
		expectedGecos   string
	}{
		{
			name: "User 1",
			output: `{
			  "user_name": "mock",
			  "user_id": "1000",
			  "user_gid": "1000",
			  "user_home_dir": "/home/mock",
			  "user_shell": "/bin/bash",
			  "user_gecos": "Mock User"
			}`,
			expectedName:    "mock",
			expectedId:      "1000",
			expectedGroupId: "1000",
			expectedHomeDir: "/home/mock",
			expectedShell:   "/bin/bash",
			expectedGecos:   "Mock User",
		},
		{
			name: "User 2",
			output: `{
			  "user_name": "test",
			  "user_id": "1001",
			  "user_gid": "1001",
			  "user_home_dir": "/home/test",
			  "user_shell": "/bin/zsh",
			  "user_gecos": "Test User"
			}`,
			expectedName:    "test",
			expectedId:      "1001",
			expectedGroupId: "1001",
			expectedHomeDir: "/home/test",
			expectedShell:   "/bin/zsh",
			expectedGecos:   "Test User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("posix")
			osInfo.id = "generic"

			info := newUserInfo()

			transport := mock.NewMockTransport()
			transport.CommandResults[userPosixDiscoveryScript] = &mock.CommandResult{
				Stdout: tt.output,
			}

			diags := info.populateUserInfo(osInfo, transport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("expected name %q, got %q", tt.expectedName, info.Name())
			}
			if info.UserId() != tt.expectedId {
				t.Errorf("expected user ID %q, got %q", tt.expectedId, info.UserId())
			}
			if info.GroupId() != tt.expectedGroupId {
				t.Errorf("expected group ID %q, got %q", tt.expectedGroupId, info.GroupId())
			}
			if info.HomeDir() != tt.expectedHomeDir {
				t.Errorf("expected home directory %q, got %q", tt.expectedHomeDir, info.HomeDir())
			}
			if info.Shell() != tt.expectedShell {
				t.Errorf("expected shell %q, got %q", tt.expectedShell, info.Shell())
			}
			if info.Gecos() != tt.expectedGecos {
				t.Errorf("expected GECOS %q, got %q", tt.expectedGecos, info.Gecos())
			}
		})
	}
}

func TestUserInfo_PopulateUserInfo_Posix_Error(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("posix")
	osInfo.id = "generic"

	info := newUserInfo()

	transport := mock.NewMockTransport()
	transport.CommandResults[userPosixDiscoveryScript] = &mock.CommandResult{
		Err: os.ErrPermission,
	}

	diags := info.populateUserInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("expected name to be empty, got '%s'", info.Name())
	}

	if info.UserId() != "" {
		t.Errorf("expected user ID to be empty, got '%s'", info.UserId())
	}

	if info.GroupId() != "" {
		t.Errorf("expected group ID to be empty, got '%s'", info.GroupId())
	}

	if info.HomeDir() != "" {
		t.Errorf("expected home directory to be empty, got '%s'", info.HomeDir())
	}

	if info.Shell() != "" {
		t.Errorf("expected shell to be empty, got '%s'", info.Shell())
	}

	if info.Gecos() != "" {
		t.Errorf("expected GECOS to be empty, got '%s'", info.Gecos())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}

	expectedSummary := "Failed to get user information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected first error summary %q, got %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error getting user information on POSIX host: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected first error detail %q, got %q", expectedDetail, errors[0].Detail)
	}
}

func TestUserInfo_PopulateUserInfo_Posix_NotJSON(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("posix")
	osInfo.id = "generic"

	info := newUserInfo()

	transport := mock.NewMockTransport()
	transport.CommandResults[userPosixDiscoveryScript] = &mock.CommandResult{
		Stdout: "Not a valid JSON output",
	}

	diags := info.populateUserInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("expected name to be empty, got '%s'", info.Name())
	}

	if info.UserId() != "" {
		t.Errorf("expected user ID to be empty, got '%s'", info.UserId())
	}

	if info.GroupId() != "" {
		t.Errorf("expected group ID to be empty, got '%s'", info.GroupId())
	}

	if info.HomeDir() != "" {
		t.Errorf("expected home directory to be empty, got '%s'", info.HomeDir())
	}

	if info.Shell() != "" {
		t.Errorf("expected shell to be empty, got '%s'", info.Shell())
	}

	if info.Gecos() != "" {
		t.Errorf("expected GECOS to be empty, got '%s'", info.Gecos())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}

	expectedSummary := "Failed to parse user information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected first error summary %q, got %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing user information on POSIX host: invalid character 'N' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected first error detail %q, got %q", expectedDetail, errors[0].Detail)
	}
}

func TestUserInfo_PopulateUserInfo_Windows(t *testing.T) {

	tests := []struct {
		name            string
		output          string
		expectedName    string
		expectedId      string
		expectedHomeDir string
	}{
		{
			name: "User 1",
			output: `{
			  "user_name": "mock",
			  "user_id": "S-1-5-21-1234567890-1234567890-1234567890-1001",
			  "user_home_dir": "C:\\Users\\mock"
			}`,
			expectedName:    "mock",
			expectedId:      "S-1-5-21-1234567890-1234567890-1234567890-1001",
			expectedHomeDir: "C:\\Users\\mock",
		},
		{
			name: "User 2",
			output: `{
			  "user_name": "test",
			  "user_id": "S-1-5-21-0987654321-0987654321-0987654321-1002",
			  "user_home_dir": "C:\\Users\\test"
			}`,
			expectedName:    "test",
			expectedId:      "S-1-5-21-0987654321-0987654321-0987654321-1002",
			expectedHomeDir: "C:\\Users\\test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("windows")
			osInfo.id = "windows-server"

			info := newUserInfo()

			transport := mock.NewWinMockTransport()
			transport.PowerShellResults[userWindowsDiscoveryScript] = &mock.CommandResult{
				Stdout: tt.output,
			}

			diags := info.populateUserInfo(osInfo, transport)
			if diags.HasErrors() {
				t.Fatalf("expected no error, got %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("expected name '%s', got '%s'", tt.expectedName, info.Name())
			}
			if info.UserId() != tt.expectedId {
				t.Errorf("expected user ID '%s', got '%s'", tt.expectedId, info.UserId())
			}
			if info.HomeDir() != tt.expectedHomeDir {
				t.Errorf("expected home directory '%s', got '%s'", tt.expectedHomeDir, info.HomeDir())
			}
		})
	}
}

func TestUserInfo_PopulateUserInfo_Windows_Error(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows-server"

	info := newUserInfo()

	transport := mock.NewWinMockTransport()
	transport.PowerShellResults[userWindowsDiscoveryScript] = &mock.CommandResult{
		Err: os.ErrPermission,
	}

	diags := info.populateUserInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("expected name to be empty, got '%s'", info.Name())
	}
	if info.UserId() != "" {
		t.Errorf("expected user ID to be empty, got '%s'", info.UserId())
	}
	if info.HomeDir() != "" {
		t.Errorf("expected home directory to be empty, got '%s'", info.HomeDir())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 diagnostic, got %d", len(errors))
	}

	expectedSummary := "Failed to get user information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected first error summary %q, got %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error getting user information on Windows host: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected first error detail %q, got %q", expectedDetail, errors[0].Detail)
	}
}

func TestUserInfo_PopulateUserInfo_Windows_NotJSON(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows-server"

	info := newUserInfo()

	transport := mock.NewWinMockTransport()
	transport.PowerShellResults[userWindowsDiscoveryScript] = &mock.CommandResult{
		Stdout: "Not a valid JSON output",
	}

	diags := info.populateUserInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("expected name to be empty, got '%s'", info.Name())
	}
	if info.UserId() != "" {
		t.Errorf("expected user ID to be empty, got '%s'", info.UserId())
	}
	if info.HomeDir() != "" {
		t.Errorf("expected home directory to be empty, got '%s'", info.HomeDir())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 diagnostic, got %d", len(errors))
	}

	expectedSummary := "Failed to parse user information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected first error summary %q, got %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing user information on Windows host: invalid character 'N' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected first error detail %q, got %q", expectedDetail, errors[0].Detail)
	}
}

func TestUserInfo_PopulateUserInfo_UnknownOS(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("unknown")
	osInfo.id = "unknown-os"

	info := newUserInfo()

	transport := mock.NewMockTransport()

	diags := info.populateUserInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	expectedSummary := "Unsupported OS family"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "User information collection is not supported on this OS"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got %q", expectedDetail, errors[0].Detail)
	}
}

func TestUserInfo_ToMapOfCtyValues_EmptyValues(t *testing.T) {

	info := newUserInfo()

	result := info.toMapOfCtyValues()
	if result == nil {
		t.Fatal("expected non-nil result, got nil")
	}

	if len(result) != 6 {
		t.Errorf("expected 6 keys, got %d", len(result))
	}

	if name, exists := result["user_name"]; exists {
		if !name.IsNull() {
			t.Errorf("expected 'user_name' to be null, got %s", name)
		}
	} else {
		t.Error("expected 'user_name' key to exist")
	}

	if userId, exists := result["user_id"]; exists {
		if !userId.IsNull() {
			t.Errorf("expected 'user_id' to be null, got %s", userId)
		}
	} else {
		t.Error("expected 'user_id' key to exist")
	}

	if groupId, exists := result["user_gid"]; exists {
		if !groupId.IsNull() {
			t.Errorf("expected 'user_gid' to be null, got %s", groupId)
		}
	} else {
		t.Error("expected 'user_gid' key to exist")
	}

	if homeDir, exists := result["user_home_dir"]; exists {
		if !homeDir.IsNull() {
			t.Errorf("expected 'user_home_dir' to be null, got %s", homeDir)
		}
	} else {
		t.Error("expected 'user_home_dir' key to exist")
	}

	if shell, exists := result["user_shell"]; exists {
		if !shell.IsNull() {
			t.Errorf("expected 'user_shell' to be null, got %s", shell)
		}
	} else {
		t.Error("expected 'user_shell' key to exist")
	}

	if gecos, exists := result["user_gecos"]; exists {
		if !gecos.IsNull() {
			t.Errorf("expected 'user_gecos' to be null, got %s", gecos)
		}
	} else {
		t.Error("expected 'user_gecos' key to exist")
	}
}

func TestUserInfo_ToMapOfCtyValues_PopulatedValues(t *testing.T) {

	info := newUserInfo()
	info.name = "mock"
	info.userId = "1000"
	info.groupId = "1000"
	info.homeDir = "/home/mock"
	info.shell = "/bin/bash"
	info.gecos = "Mock User"

	result := info.toMapOfCtyValues()
	if result == nil {
		t.Fatal("expected non-nil result, got nil")
	}

	if len(result) != 6 {
		t.Errorf("expected 6 keys, got %d", len(result))
	}

	if name, exists := result["user_name"]; exists {
		if name.AsString() != "mock" {
			t.Errorf("expected 'user_name' to be 'mock', got %s", name)
		}
	} else {
		t.Error("expected 'user_name' key to exist")
	}

	if userId, exists := result["user_id"]; exists {
		if userId.AsString() != "1000" {
			t.Errorf("expected 'user_id' to be '1000', got %s", userId)
		}
	} else {
		t.Error("expected 'user_id' key to exist")
	}

	if groupId, exists := result["user_gid"]; exists {
		if groupId.AsString() != "1000" {
			t.Errorf("expected 'user_gid' to be '1000', got %s", groupId)
		}
	} else {
		t.Error("expected 'user_gid' key to exist")
	}

	if homeDir, exists := result["user_home_dir"]; exists {
		if homeDir.AsString() != "/home/mock" {
			t.Errorf("expected 'user_home_dir' to be '/home/mock', got %s", homeDir)
		}
	} else {
		t.Error("expected 'user_home_dir' key to exist")
	}

	if shell, exists := result["user_shell"]; exists {
		if shell.AsString() != "/bin/bash" {
			t.Errorf("expected 'user_shell' to be '/bin/bash', got %s", shell)
		}
	} else {
		t.Error("expected 'user_shell' key to exist")
	}

	if gecos, exists := result["user_gecos"]; exists {
		if gecos.AsString() != "Mock User" {
			t.Errorf("expected 'user_gecos' to be 'Mock User', got %s", gecos)
		}
	} else {
		t.Error("expected 'user_gecos' key to exist")
	}
}
