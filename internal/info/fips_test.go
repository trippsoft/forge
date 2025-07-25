package info

import (
	"testing"
)

func TestFipsInfo_PopulateFipsInfo_Linux_Enabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockTransport := &mockTransport{commandOutput: "1\n"}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fipsInfo.known {
		t.Error("expected FIPS to be known on Linux system")
	}

	if !fipsInfo.enabled {
		t.Error("expected FIPS to be enabled when file contains '1'")
	}
}

func TestFipsInfo_PopulateFipsInfo_Linux_Disabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockTransport := &mockTransport{commandOutput: "0\n"}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fipsInfo.known {
		t.Error("expected FIPS to be known on Linux system")
	}

	if fipsInfo.enabled {
		t.Error("expected FIPS to be disabled when file contains '0'")
	}
}

func TestFipsInfo_PopulateFipsInfo_Windows_Enabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	mockTransport := &mockTransport{
		powershellOutput: "1",
	}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fipsInfo.known {
		t.Error("expected FIPS to be known on Windows system")
	}

	if !fipsInfo.enabled {
		t.Error("expected FIPS to be enabled when PowerShell returns '1'")
	}
}

func TestFipsInfo_PopulateFipsInfo_Windows_Disabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	mockTransport := &mockTransport{
		powershellOutput: "0",
	}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fipsInfo.known {
		t.Error("expected FIPS to be known on Windows system")
	}

	if fipsInfo.enabled {
		t.Error("expected FIPS to be disabled when PowerShell returns '0'")
	}
}

func TestFipsInfo_PopulateFipsInfo_UnknownOS(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("darwin") // macOS

	mockTransport := &mockTransport{}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fipsInfo.known {
		t.Error("expected FIPS to be unknown on unsupported OS")
	}

	if fipsInfo.enabled {
		t.Error("expected FIPS to be disabled on unsupported OS")
	}
}

func TestFipsInfo_PopulateFipsInfo_CommandError(t *testing.T) {
	// Test FIPS with file system read error
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockTransport := &mockTransport{}

	fipsInfo := &fipsInfo{}
	err := fipsInfo.populateFipsInfo(osInfo, mockTransport)
	if err == nil {
		t.Error("expected error when reading FIPS file fails")
	}
}

func TestFipsInfo_ToMapOfCtyValues_Known(t *testing.T) {
	fipsInfo := &fipsInfo{
		known:   true,
		enabled: true,
	}

	values := fipsInfo.toMapOfCtyValues()

	if _, exists := values["fips_enabled"]; !exists {
		t.Error("expected fips_enabled key to be present in values map")
	}

	if !values["fips_enabled"].True() {
		t.Error("expected fips_enabled to be true")
	}
}

func TestFipsInfo_ToMapOfCtyValues_KnownButDisabled(t *testing.T) {
	fipsInfo := &fipsInfo{
		known:   true,
		enabled: false,
	}

	values := fipsInfo.toMapOfCtyValues()

	if _, exists := values["fips_enabled"]; !exists {
		t.Error("expected fips_enabled key to be present in values map")
	}

	if values["fips_enabled"].True() {
		t.Error("expected fips_enabled to be false")
	}
}

func TestFipsInfo_ToMapOfCtyValues_Unknown(t *testing.T) {
	fipsInfo := &fipsInfo{
		known:   false,
		enabled: false,
	}

	values := fipsInfo.toMapOfCtyValues()

	if value, exists := values["fips_enabled"]; exists {
		if !value.IsNull() {
			t.Error("expected fips_enabled to be null for unknown FIPS")
		}
	} else {
		t.Error("expected fips_enabled key to be present in values map")
	}
}
