package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/pkg/transport"
)

func TestAppArmorInfo_PopulateAppArmorInfo_NoOS(t *testing.T) {

	osInfo := newOSInfo()

	transport := transport.NewMockTransport()

	info := newAppArmorInfo()
	diags := info.populateAppArmorInfo(osInfo, transport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatal("expected warnings, got none")
	}

	if info.Supported() {
		t.Error("expected AppArmor to be unsupported on non-Linux system")
	}

	if info.Enabled() {
		t.Error("expected AppArmor to be disabled on non-Linux system")
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Missing OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping AppArmor information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_Linux(t *testing.T) {

	tests := []struct {
		name            string
		output          string
		expectedEnabled bool
	}{
		{
			name:            "AppArmor enabled",
			output:          "1",
			expectedEnabled: true,
		},
		{
			name:            "AppArmor disabled",
			output:          "0",
			expectedEnabled: false,
		},
		{
			name:            "AppArmor not supported",
			output:          "",
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("linux")
			osInfo.id = "ubuntu"

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[apparmorDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newAppArmorInfo()
			diags := info.populateAppArmorInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if !info.Supported() {
				t.Error("expected AppArmor to be supported on Linux system")
			}

			if info.Enabled() != tt.expectedEnabled {
				t.Errorf("expected AppArmor to be %v, got: %v", tt.expectedEnabled, info.Enabled())
			}
		})
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_Linux_Error(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[apparmorDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newAppArmorInfo()
	diags := info.populateAppArmorInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Error("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.Supported() {
		t.Error("expected AppArmor to be unsupported on Linux system with file system error")
	}

	if info.Enabled() {
		t.Error("expected AppArmor to be disabled on Linux system with file system error")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 diagnostic, got: %d", len(errors))
	}

	expectedSummary := "Failed to execute AppArmor discovery script"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing command: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_Windows(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows-server"

	transport := transport.NewMockTransport()

	info := newAppArmorInfo()
	diags := info.populateAppArmorInfo(osInfo, transport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Supported() {
		t.Error("expected AppArmor to be unsupported on Windows system")
	}

	if info.Enabled() {
		t.Error("expected AppArmor to be disabled on Windows system")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_Supported(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = true
	appArmorInfo.enabled = true

	values := appArmorInfo.toMapOfCtyValues()

	if _, exists := values["apparmor_enabled"]; !exists {
		t.Error("expected 'apparmor_enabled' key to be present in values map")
	}

	if !values["apparmor_enabled"].True() {
		t.Error("expected 'apparmor_enabled' to be true")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_SupportedButDisabled(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = true
	appArmorInfo.enabled = false

	values := appArmorInfo.toMapOfCtyValues()

	if _, exists := values["apparmor_enabled"]; !exists {
		t.Error("expected 'apparmor_enabled' key to be present in values map")
	}

	if values["apparmor_enabled"].True() {
		t.Error("expected 'apparmor_enabled' to be false")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_NotSupported(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = false
	appArmorInfo.enabled = false

	values := appArmorInfo.toMapOfCtyValues()

	if value, exists := values["apparmor_enabled"]; exists {
		if !value.IsNull() {
			t.Error("expected 'apparmor_enabled' to be null for unsupported AppArmor")
		}
	} else {
		t.Error("expected 'apparmor_enabled' key to be present in values map")
	}
}
