package info

import (
	"os"
	"testing"
	"time"
)

func TestAppArmorInfo_PopulateAppArmorInfo_NonLinux(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	appArmorInfo := newAppArmorInfo()
	err := appArmorInfo.populateAppArmorInfo(osInfo, &mockFileSystem{})

	if err != nil {
		t.Fatalf("expected no error for non-Linux system, got: %v", err)
	}

	if appArmorInfo.supported {
		t.Error("expected AppArmor to be unsupported on non-Linux system")
	}

	if appArmorInfo.enabled {
		t.Error("expected AppArmor to be disabled on non-Linux system")
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_Linux_Enabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockFS := &mockFileSystem{
		dirs: map[string]*mockFileInfo{
			"/sys/kernel/security/apparmor": {
				name:    "apparmor",
				isDir:   true,
				mode:    os.ModeDir | 0755,
				modTime: time.Now(),
			},
		},
	}

	appArmorInfo := newAppArmorInfo()
	err := appArmorInfo.populateAppArmorInfo(osInfo, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !appArmorInfo.supported {
		t.Error("expected AppArmor to be supported on Linux system")
	}

	if !appArmorInfo.enabled {
		t.Error("expected AppArmor to be enabled when directory exists")
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_Linux_Disabled(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockFS := &mockFileSystem{
		dirs: map[string]*mockFileInfo{},
	}

	appArmorInfo := newAppArmorInfo()
	err := appArmorInfo.populateAppArmorInfo(osInfo, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !appArmorInfo.supported {
		t.Error("expected AppArmor to be supported on Linux system")
	}

	if appArmorInfo.enabled {
		t.Error("expected AppArmor to be disabled when directory doesn't exist")
	}
}

func TestAppArmorInfo_PopulateAppArmorInfo_FileSystemError(t *testing.T) {
	// Test AppArmor with file system error (not IsNotExist)
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	mockFS := &errorFileSystem{shouldError: true}

	appArmorInfo := newAppArmorInfo()
	err := appArmorInfo.populateAppArmorInfo(osInfo, mockFS)
	if err == nil {
		t.Error("expected error when file system operations fail")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_Supported(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = true
	appArmorInfo.enabled = true

	values := appArmorInfo.toMapOfCtyValues()

	if _, exists := values["apparmor_enabled"]; !exists {
		t.Error("expected apparmor_enabled key to be present in values map")
	}

	if !values["apparmor_enabled"].True() {
		t.Error("expected apparmor_enabled to be true")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_SupportedButDisabled(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = true
	appArmorInfo.enabled = false

	values := appArmorInfo.toMapOfCtyValues()

	if _, exists := values["apparmor_enabled"]; !exists {
		t.Error("expected apparmor_enabled key to be present in values map")
	}

	if values["apparmor_enabled"].True() {
		t.Error("expected apparmor_enabled to be false")
	}
}

func TestAppArmorInfo_ToMapOfCtyValues_NotSupported(t *testing.T) {
	appArmorInfo := newAppArmorInfo()
	appArmorInfo.supported = false
	appArmorInfo.enabled = false

	values := appArmorInfo.toMapOfCtyValues()

	if value, exists := values["apparmor_enabled"]; exists {
		if !value.IsNull() {
			t.Error("expected apparmor_enabled to be null for unsupported AppArmor")
		}
	} else {
		t.Error("expected apparmor_enabled key to be present in values map")
	}
}
