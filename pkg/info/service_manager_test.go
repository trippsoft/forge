package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/pkg/transport"
)

func TestServiceManagerInfo_PopulateServiceManagerInfo_NoOS(t *testing.T) {

	osInfo := newOSInfo()

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatal("expected warnings, got none")
	}

	if info.Name() != "" {
		t.Error("expected service manager name to be empty with missing OS info")
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Missing OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping service manager information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Darwin(t *testing.T) {

	tests := []struct {
		name         string
		majorVersion string
		version      string
		expected     string
	}{
		{
			name:         "macOS 12",
			majorVersion: "12",
			version:      "", // not used in this test
			expected:     "launchd",
		},
		{
			name:         "macOS 11",
			majorVersion: "11",
			version:      "", // not used in this test
			expected:     "launchd",
		},
		{
			name:         "macOS 10.0",
			majorVersion: "10",
			version:      "10.0",
			expected:     "systemstarter",
		},
		{
			name:         "macOS 10.3",
			majorVersion: "10",
			version:      "10.3",
			expected:     "systemstarter",
		},
		{
			name:         "macOS 10.4",
			majorVersion: "10",
			version:      "10.4",
			expected:     "launchd",
		},
		{
			name:         "macOS 9",
			majorVersion: "9",
			version:      "", // not used in this test
			expected:     "systemstarter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("darwin") // macOS
			osInfo.id = "macos"
			osInfo.majorVersion = tt.majorVersion
			osInfo.version = tt.version

			mockTransport := transport.NewMockTransport()

			info := newServiceManagerInfo()
			diags := info.populateServiceManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatal("expected no warnings, got some")
			}

			if info.Name() != tt.expected {
				t.Errorf("expected service manager name to be %q, got: %q", tt.expected, info.Name())
			}
		})
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Darwin_InvalidMajorVersion(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("darwin") // macOS
	osInfo.id = "macos"
	osInfo.majorVersion = "invalid" // invalid major version

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "launchd" // default name before error handling
	if info.Name() != expectedName {
		t.Errorf("expected service manager name to be %q, got: %q", expectedName, info.Name())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Invalid OS major version"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing OS major version: strconv.Atoi: parsing \"invalid\": invalid syntax"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Darwin_NotEnoughVersionParts(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("darwin") // macOS
	osInfo.id = "macos"
	osInfo.majorVersion = "10"
	osInfo.version = "invalid" // invalid version format

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "launchd" // default name before error handling
	if info.Name() != expectedName {
		t.Errorf("expected service manager name to be %q, got: %q", expectedName, info.Name())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Invalid OS version format"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "OS version does not contain enough parts to determine service manager"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Darwin_InvalidVersion(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("darwin") // macOS
	osInfo.id = "macos"
	osInfo.majorVersion = "10"
	osInfo.version = "10.invalid" // invalid minor version

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "launchd" // default name before error handling
	if info.Name() != expectedName {
		t.Errorf("expected service manager name to be %q, got: %q", expectedName, info.Name())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Invalid OS minor version"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing OS minor version: strconv.Atoi: parsing \"invalid\": invalid syntax"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Windows(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows"

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Fatal("expected no warnings, got some")
	}

	expectedName := "windows-service-manager"
	if info.Name() != expectedName {
		t.Errorf("expected service manager name to be %q, got: %q", expectedName, info.Name())
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Linux(t *testing.T) {

	tests := []struct {
		name     string
		output   string
		expected string
	}{
		{
			name: "systemd - /run/systemd/system",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "1",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "systemd"
				}
				`,
			expected: "systemd",
		},
		{
			name: "systemd - /dev/.run/systemd",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "1",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "systemd"
			}`,
			expected: "systemd",
		},
		{
			name: "systemd - /dev/.systemd",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "1",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "systemd"
			}`,
			expected: "systemd",
		},
		{
			name: "upstart",
			output: `{
				"systemctl_exists": "0",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "1",
				"etc_init_exists": "1",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "init"
			}`,
			expected: "upstart",
		},
		{
			name: "openrc",
			output: `{
				"systemctl_exists": "0",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "1",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "openrc"
			}`,
			expected: "openrc",
		},
		{
			name: "systemd - init link target",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "/lib/systemd",
				"etc_init_d_exists": "0",
				"proc1_comm": "systemd"
			}`,
			expected: "systemd",
		},
		{
			name: "sysvinit",
			output: `{
				"systemctl_exists": "0",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "1",
				"proc1_comm": "sysvinit"
			}`,
			expected: "sysvinit",
		},
		{
			name: "systemd - /proc/1/comm",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "systemd"
			}`,
			expected: "systemd",
		},
		{
			name: "openrc-init - /proc/1/comm",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "",
				"etc_init_d_exists": "0",
				"proc1_comm": "openrc-init"
			}`,
			expected: "openrc",
		},
		{
			name: "openrc-init - init link target",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "/sbin/openrc-init",
				"etc_init_d_exists": "0",
				"proc1_comm": ""
			}`,
			expected: "openrc",
		},
		{
			name: "sysvinit - init link target",
			output: `{
				"systemctl_exists": "1",
				"run_systemd_system_exists": "0",
				"dev_run_systemd_exists": "0",
				"dev_systemd_exists": "0",
				"initctl_exists": "0",
				"etc_init_exists": "0",
				"openrc_exists": "0",
				"init_link_target": "/sbin/sysvinit",
				"etc_init_d_exists": "0",
				"proc1_comm": ""
			}`,
			expected: "sysvinit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("linux")
			osInfo.id = "ubuntu"

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[linuxServiceManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newServiceManagerInfo()
			diags := info.populateServiceManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expected {
				t.Errorf("expected service manager name to be %q, got: %q", tt.expected, info.Name())
			}
		})
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Linux_Error(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[linuxServiceManagerDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Error("expected service manager name to be empty with no output")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to execute service manager discovery script"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing service manager discovery script: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Linux_NoOutput(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[linuxServiceManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: "",
	}

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Error("expected service manager name to be empty with no output")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse service manager discovery script output"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing service manager discovery script output: unexpected end of JSON input"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_Linux_NoServiceManager(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[linuxServiceManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			"systemctl_exists": "0",
			"run_systemd_system_exists": "0",
			"dev_run_systemd_exists": "0",
			"dev_systemd_exists": "0",
			"initctl_exists": "0",
			"etc_init_exists": "0",
			"openrc_exists": "0",
			"init_link_target": "",
			"etc_init_d_exists": "0",
			"proc1_comm": "/bin/bash"
			}
			`,
	}

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Error("expected service manager name to be empty with no service manager found")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to determine service manager"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Could not identify the service manager for the current Linux system"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_PopulateServiceManagerInfo_UnknownOS(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "generic"

	mockTransport := transport.NewMockTransport()

	info := newServiceManagerInfo()
	diags := info.populateServiceManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("expected service manager name to be empty for unknown OS, got: %s", info.Name())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Unsupported OS family"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Service manager information collection is not supported for this OS family"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestServiceManagerInfo_ToMapOfCtyValues(t *testing.T) {

	tests := []struct {
		name        string
		serviceName string
	}{
		{
			name:        "systemd",
			serviceName: "systemd",
		},
		{
			name:        "upstart",
			serviceName: "upstart",
		},
		{
			name:        "openrc",
			serviceName: "openrc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newServiceManagerInfo()
			info.name = tt.serviceName

			values := info.toMapOfCtyValues()

			if tt.serviceName != values["service_manager"].AsString() {
				t.Errorf("expected service_manager to be %q, got: %q", tt.serviceName, values["service_manager"].AsString())
			}
		})
	}
}

func TestServiceManagerInfo_ToMapOfCtyValues_Empty(t *testing.T) {

	info := newServiceManagerInfo()

	values := info.toMapOfCtyValues()

	if !values["service_manager"].IsNull() {
		t.Error("expected service_manager to be null, got a value")
	}
}
