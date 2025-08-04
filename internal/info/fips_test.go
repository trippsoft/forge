package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/internal/transport"
)

func TestFipsInfo_PopulateFipsInfo_NoOS(t *testing.T) {

	osInfo := newOSInfo()

	mockTransport := transport.NewMockTransport()

	info := newFipsInfo()
	diags := info.populateFipsInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatal("expected warnings, got none")
	}

	if info.Known() {
		t.Error("expected FIPS to be unknown with missing OS info")
	}

	if info.Enabled() {
		t.Error("expected FIPS to be disabled with missing OS info")
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Missing OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping FIPS information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestFipsInfo_PopulateFipsInfo_Darwin(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("darwin") // macOS
	osInfo.families.Add("macos")
	osInfo.id = "macos"

	mockTransport := transport.NewMockTransport()

	info := newFipsInfo()
	diags := info.populateFipsInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Known() {
		t.Error("expected FIPS to be unknown on macOS")
	}

	if info.Enabled() {
		t.Error("expected FIPS to be disabled on macOS")
	}
}

func TestFipsInfo_PopulateFipsInfo_Linux(t *testing.T) {

	tests := []struct {
		name            string
		output          string
		expectedEnabled bool
	}{
		{
			name:            "FIPS enabled",
			output:          "1",
			expectedEnabled: true,
		},
		{
			name:            "FIPS disabled",
			output:          "0",
			expectedEnabled: false,
		},
		{
			name:            "FIPS not supported",
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
			mockTransport.CommandResults[fipsLinuxDiscoveryScript] = &transport.CommandResult{
				Stdout: tt.output,
			}

			info := newFipsInfo()
			diags := info.populateFipsInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if !info.Known() {
				t.Error("expected FIPS to be known on Linux system")
			}

			if info.Enabled() != tt.expectedEnabled {
				t.Errorf("expected FIPS to be %v, got: %v", tt.expectedEnabled, info.Enabled())
			}
		})
	}
}

func TestFipsInfo_PopulateFipsInfo_Linux_Error(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.id = "ubuntu"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[fipsLinuxDiscoveryScript] = &transport.CommandResult{
		Err: os.ErrPermission,
	}

	info := newFipsInfo()
	diags := info.populateFipsInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.Known() {
		t.Error("expected FIPS to be known on Linux system")
	}

	if info.Enabled() {
		t.Error("expected FIPS to be disabled when command execution fails")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to check FIPS status"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error checking FIPS status: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestFipsInfo_PopulateFipsInfo_Windows(t *testing.T) {

	tests := []struct {
		name            string
		output          string
		expectedEnabled bool
	}{
		{
			name:            "FIPS enabled",
			output:          "1",
			expectedEnabled: true,
		},
		{
			name:            "FIPS disabled",
			output:          "0",
			expectedEnabled: false,
		},
		{
			name:            "FIPS not supported",
			output:          "",
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.families.Add("windows")
			osInfo.id = "windows-server"

			mockTransport := transport.NewWinMockTransport()
			mockTransport.PowerShellResults[fipsWindowsDiscoveryScript] = &transport.CommandResult{
				Stdout: tt.output,
			}

			info := newFipsInfo()
			diags := info.populateFipsInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if !info.Known() {
				t.Error("expected FIPS to be known on Windows system")
			}

			if info.Enabled() != tt.expectedEnabled {
				t.Errorf("expected FIPS to be %v, got: %v", tt.expectedEnabled, info.Enabled())
			}
		})
	}
}

func TestFipsInfo_PopulateFipsInfo_Windows_Error(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")
	osInfo.id = "windows-server"

	mockTransport := transport.NewWinMockTransport()
	mockTransport.PowerShellResults[fipsWindowsDiscoveryScript] = &transport.CommandResult{
		Err: os.ErrPermission,
	}

	info := newFipsInfo()
	diags := info.populateFipsInfo(osInfo, mockTransport)
	if !diags.HasErrors() {
		t.Fatal("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.Known() {
		t.Error("expected FIPS to be known on Windows system")
	}

	if info.Enabled() {
		t.Error("expected FIPS to be disabled when command execution fails")
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to check FIPS status"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error checking FIPS status: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestFipsInfo_ToMapOfCtyValues_Known_Enabled(t *testing.T) {
	info := newFipsInfo()
	info.known = true
	info.enabled = true

	values := info.toMapOfCtyValues()

	if _, exists := values["fips_enabled"]; !exists {
		t.Error("expected fips_enabled key to be present in values map")
	}

	if !values["fips_enabled"].True() {
		t.Error("expected fips_enabled to be true")
	}
}

func TestFipsInfo_ToMapOfCtyValues_Known_Disabled(t *testing.T) {
	info := newFipsInfo()
	info.known = true
	info.enabled = false

	values := info.toMapOfCtyValues()

	if _, exists := values["fips_enabled"]; !exists {
		t.Error("expected fips_enabled key to be present in values map")
	}

	if values["fips_enabled"].True() {
		t.Error("expected fips_enabled to be false")
	}
}

func TestFipsInfo_ToMapOfCtyValues_Unknown(t *testing.T) {

	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "Enabled",
			enabled: true,
		},
		{
			name:    "Disabled",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newFipsInfo()
			info.known = false
			info.enabled = tt.enabled

			values := info.toMapOfCtyValues()

			if value, exists := values["fips_enabled"]; exists {
				if !value.IsNull() {
					t.Error("expected fips_enabled to be null for unknown FIPS")
				}
			} else {
				t.Error("expected fips_enabled key to be present in values map")
			}
		})
	}
}
