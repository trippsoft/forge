package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/internal/transport/mock"
	"github.com/zclconf/go-cty/cty"
)

func TestSelinuxInfo_PopulateSelinuxInfo_NoOS(t *testing.T) {
	osInfo := newOSInfo()

	transport := mock.NewMockTransport()

	info := newSELinuxInfo()
	diags := info.populateSelinuxInfo(osInfo, transport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags)
	}

	if !diags.HasWarnings() {
		t.Fatal("expected warnings, got none")
	}

	if info.Supported() {
		t.Error("expected SELinux to be unsupported on non-Linux system")
	}

	if info.Installed() {
		t.Error("expected SELinux to be uninstalled on non-Linux system")
	}

	if info.Status() != SELinuxNotSupported {
		t.Errorf("expected status to be %q, got %q", SELinuxNotSupported, info.Status())
	}

	if info.SelinuxType() != SELinuxTypeNotSupported {
		t.Errorf("expected type to be %q, got %q", SELinuxTypeNotSupported, info.SelinuxType())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Missing OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected warning summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping SELinux information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected warning detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Windows(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows-server"

	transport := mock.NewWinMockTransport()

	info := newSELinuxInfo()
	diags := info.populateSelinuxInfo(osInfo, transport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Supported() {
		t.Error("expected SELinux to be unsupported on non-Linux system")
	}

	if info.Installed() {
		t.Error("expected SELinux to be uninstalled on non-Linux system")
	}

	if info.Status() != SELinuxNotSupported {
		t.Errorf("expected status to be %q, got %q", SELinuxNotSupported, info.Status())
	}

	if info.SelinuxType() != SELinuxTypeNotSupported {
		t.Errorf("expected type to be %q, got %q", SELinuxTypeNotSupported, info.SelinuxType())
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Linux(t *testing.T) {

	tests := []struct {
		name              string
		output            string
		expectedInstalled bool
		expectedStatus    SELinuxStatus
		expectedType      SELinuxType
	}{
		{
			name: "not installed",
			output: `{
				  "selinux_installed": "0",
				  "selinux_status": "",
				  "selinux_type": ""
				}`,
			expectedInstalled: false,
			expectedStatus:    SELinuxNotSupported,
			expectedType:      SELinuxTypeNotSupported,
		},
		{
			name: "disabled",
			output: `{
				  "selinux_installed": "1",
				  "selinux_status": "disabled",
				  "selinux_type": ""
				}`,
			expectedInstalled: true,
			expectedStatus:    SELinuxDisabled,
			expectedType:      SELinuxTypeNotSupported,
		},
		{
			name: "enforcing targeted",
			output: `{
				  "selinux_installed": "1",
				  "selinux_status": "enforcing",
				  "selinux_type": "targeted"
				}`,
			expectedInstalled: true,
			expectedStatus:    SELinuxEnforcing,
			expectedType:      SELinuxTypeTargeted,
		},
		{
			name: "permissive minimum",
			output: `{
				  "selinux_installed": "1",
				  "selinux_status": "permissive",
				  "selinux_type": "minimum"
				}`,
			expectedInstalled: true,
			expectedStatus:    SELinuxPermissive,
			expectedType:      SELinuxTypeMinimum,
		},
		{
			name: "enforcing mls",
			output: `{
				  "selinux_installed": "1",
				  "selinux_status": "enforcing",
				  "selinux_type": "mls"
				}`,
			expectedInstalled: true,
			expectedStatus:    SELinuxEnforcing,
			expectedType:      SELinuxTypeMLS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("linux")
			osInfo.id = "ubuntu"

			transport := mock.NewMockTransport()
			transport.CommandResults[selinuxDiscoveryScript] = &mock.CommandResult{
				Stdout: tt.output,
			}

			info := newSELinuxInfo()
			diags := info.populateSelinuxInfo(osInfo, transport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if !info.Supported() {
				t.Error("expected SELinux to be supported on Linux system")
			}

			if info.Installed() != tt.expectedInstalled {
				t.Errorf("expected SELinux installed to be %v, got: %v", tt.expectedInstalled, info.Installed())
			}

			if info.Status() != tt.expectedStatus {
				t.Errorf("expected status to be %q, got %q", tt.expectedStatus, info.Status())
			}

			if info.SelinuxType() != tt.expectedType {
				t.Errorf("expected type to be %q, got %q", tt.expectedType, info.SelinuxType())
			}
		})
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Linux_Error(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	transport := mock.NewMockTransport()
	transport.CommandResults[selinuxDiscoveryScript] = &mock.CommandResult{
		Err: os.ErrPermission,
	}

	info := newSELinuxInfo()
	diags := info.populateSelinuxInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.supported {
		t.Error("expected SELinux to be supported on Linux system")
	}

	if info.Installed() {
		t.Error("expected SELinux to be not installed on Linux system with config file stat error")
	}

	if info.Status() != SELinuxNotSupported {
		t.Errorf("expected status to be %q, got %q", SELinuxNotSupported, info.Status())
	}

	if info.SelinuxType() != SELinuxTypeNotSupported {
		t.Errorf("expected type to be %q, got %q", SELinuxTypeNotSupported, info.SelinuxType())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to execute SELinux discovery script"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing SELinux discovery script: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Linux_NotJSON(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	transport := mock.NewMockTransport()
	transport.CommandResults[selinuxDiscoveryScript] = &mock.CommandResult{
		Stdout: "Not a valid JSON output",
	}

	info := newSELinuxInfo()
	diags := info.populateSelinuxInfo(osInfo, transport)

	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.supported {
		t.Error("expected SELinux to be supported on Linux system")
	}

	if info.Installed() {
		t.Error("expected SELinux to be not installed on Linux system with config file stat error")
	}

	if info.Status() != SELinuxNotSupported {
		t.Errorf("expected status to be %q, got %q", SELinuxNotSupported, info.Status())
	}

	if info.SelinuxType() != SELinuxTypeNotSupported {
		t.Errorf("expected type to be %q, got %q", SELinuxTypeNotSupported, info.SelinuxType())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse SELinux discovery script output"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing SELinux discovery script output: invalid character 'N' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestSelinuxInfo_ToMapOfCtyValues_Supported(t *testing.T) {

	info := &SELinuxInfo{
		supported:   true,
		installed:   true,
		status:      SELinuxEnforcing,
		selinuxType: SELinuxTypeTargeted,
	}

	values := info.toMapOfCtyValues()

	expectedKeys := []string{"selinux_installed", "selinux_status", "selinux_type"}
	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	if values["selinux_installed"].True() != info.installed {
		t.Errorf("expected selinux_installed to be true, got %s", values["selinux_installed"].GoString())
	}

	if values["selinux_status"].AsString() != string(SELinuxEnforcing) {
		t.Errorf("expected selinux_status to be %q, got %q", SELinuxEnforcing, values["selinux_status"].AsString())
	}

	if values["selinux_type"].AsString() != string(SELinuxTypeTargeted) {
		t.Errorf("expected selinux_type to be %q, got %q", SELinuxTypeTargeted, values["selinux_type"].AsString())
	}
}

func TestSelinuxInfo_ToMapOfCtyValues_NotInstalled(t *testing.T) {

	info := &SELinuxInfo{
		supported:   true,
		installed:   false,
		status:      SELinuxEnforcing,    // Value doesn't matter here and should be ignored
		selinuxType: SELinuxTypeTargeted, // Value doesn't matter here and should be ignored
	}

	values := info.toMapOfCtyValues()

	expectedKeys := []string{"selinux_installed", "selinux_status", "selinux_type"}
	for _, key := range expectedKeys {
		if value, exists := values[key]; exists {
			if key == "selinux_installed" && value != cty.False {
				t.Errorf("expected selinux_installed to be false for SELinux not installed, got %s", value.GoString())
			}
			if key != "selinux_installed" && !value.IsNull() {
				t.Errorf("expected key %q to be null for SELinux not installed, got %s", key, value.GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}
}

func TestSelinuxInfo_ToMapOfCtyValues_NotSupported(t *testing.T) {

	info := &SELinuxInfo{
		supported:   false,
		installed:   true,                    // Value doesn't matter here and should be ignored
		status:      SELinuxNotSupported,     // Value doesn't matter here and should be ignored
		selinuxType: SELinuxTypeNotSupported, // Value doesn't matter here and should be ignored
	}

	values := info.toMapOfCtyValues()

	expectedKeys := []string{"selinux_status", "selinux_type"}
	for _, key := range expectedKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected key %q to be null for unsupported SELinux, got %s", key, value.GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}
}
