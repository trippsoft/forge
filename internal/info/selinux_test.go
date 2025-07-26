package info

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/zclconf/go-cty/cty"
)

func TestSelinuxInfo_PopulateSelinuxInfo_NonLinux(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	fileSystem := newMockFileSystem()

	selinuxInfo := &selinuxInfo{}

	err := selinuxInfo.populateSelinuxInfo(osInfo, fileSystem)

	if err != nil {
		t.Fatalf("expected no error for non-Linux system, got: %v", err)
	}

	if selinuxInfo.supported {
		t.Error("expected SELinux to be unsupported on non-Linux system")
	}

	if selinuxInfo.status != SelinuxNotSupported {
		t.Errorf("expected status to be %q, got %q", SelinuxNotSupported, selinuxInfo.status)
	}

	if selinuxInfo.selinuxType != SelinuxTypeNotSupported {
		t.Errorf("expected type to be %q, got %q", SelinuxTypeNotSupported, selinuxInfo.selinuxType)
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Linux_NoConfigFile(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")

	fileSystem := newMockFileSystem()

	selinuxInfo := &selinuxInfo{}
	err := selinuxInfo.populateSelinuxInfo(osInfo, fileSystem)

	if err != nil {
		t.Fatalf("expected no error when config file doesn't exist, got: %v", err)
	}

	if !selinuxInfo.supported {
		t.Error("expected SELinux to be supported on Linux system")
	}
}

func TestSelinuxInfo_PopulateSelinuxInfo_Linux_WithConfigFile(t *testing.T) {
	testCases := []struct {
		name           string
		configContent  string
		expectedStatus selinuxStatus
		expectedType   selinuxType
	}{
		{
			name:           "disabled",
			configContent:  "SELINUX=disabled\nSELINUXTYPE=targeted\n",
			expectedStatus: SelinuxDisabled,
			expectedType:   SelinuxTypeNotSupported,
		},
		{
			name:           "enforcing_targeted",
			configContent:  "SELINUX=enforcing\nSELINUXTYPE=targeted\n",
			expectedStatus: SelinuxEnforcing,
			expectedType:   SelinuxTypeTargeted,
		},
		{
			name:           "permissive_minimum",
			configContent:  "SELINUX=permissive\nSELINUXTYPE=minimum\n",
			expectedStatus: SelinuxPermissive,
			expectedType:   SelinuxTypeMinimum,
		},
		{
			name:           "enforcing_mls",
			configContent:  "SELINUX=enforcing\nSELINUXTYPE=mls\n",
			expectedStatus: SelinuxEnforcing,
			expectedType:   SelinuxTypeMLS,
		},
		{
			name:           "with_comments_and_empty_lines",
			configContent:  "# This is a comment\n\nSELINUX=enforcing\n# Another comment\nSELINUXTYPE=targeted\n\n",
			expectedStatus: SelinuxEnforcing,
			expectedType:   SelinuxTypeTargeted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osInfo := newOSInfo()
			osInfo.families.Add("linux")

			fileSystem := newMockFileSystem()
			fileSystem.files["/etc/selinux/config"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tc.configContent)),
				info: &mockFileInfo{
					name:    "config",
					size:    int64(len(tc.configContent)),
					mode:    0644,
					modTime: time.Now(),
				},
			}

			selinuxInfo := &selinuxInfo{}
			err := selinuxInfo.populateSelinuxInfo(osInfo, fileSystem)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !selinuxInfo.supported {
				t.Error("expected SELinux to be supported on Linux system")
			}

			if selinuxInfo.status != tc.expectedStatus {
				t.Errorf("expected status %q, got %q", tc.expectedStatus, selinuxInfo.status)
			}

			if selinuxInfo.selinuxType != tc.expectedType {
				t.Errorf("expected type %q, got %q", tc.expectedType, selinuxInfo.selinuxType)
			}
		})
	}
}

func TestSelinuxInfo_ToMapOfCtyValues_Supported(t *testing.T) {
	selinuxInfo := &selinuxInfo{
		supported:   true,
		installed:   true,
		status:      SelinuxEnforcing,
		selinuxType: SelinuxTypeTargeted,
	}

	values := selinuxInfo.toMapOfCtyValues()

	expectedKeys := []string{"selinux_installed", "selinux_status", "selinux_type"}
	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	if values["selinux_installed"].True() != selinuxInfo.installed {
		t.Errorf("expected selinux_installed to be true, got %s", values["selinux_installed"].GoString())
	}

	if values["selinux_status"].AsString() != string(SelinuxEnforcing) {
		t.Errorf("expected selinux_status to be %q, got %q", SelinuxEnforcing, values["selinux_status"].AsString())
	}

	if values["selinux_type"].AsString() != string(SelinuxTypeTargeted) {
		t.Errorf("expected selinux_type to be %q, got %q", SelinuxTypeTargeted, values["selinux_type"].AsString())
	}
}

func TestSelinuxInfo_ToMapOfCtyValues_NotInstalled(t *testing.T) {
	selinuxInfo := &selinuxInfo{
		supported:   true,
		installed:   false,
		status:      SelinuxEnforcing,    // Value doesn't matter here and should be ignored
		selinuxType: SelinuxTypeTargeted, // Value doesn't matter here and should be ignored
	}

	values := selinuxInfo.toMapOfCtyValues()

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
	selinuxInfo := &selinuxInfo{
		supported:   false,
		installed:   true,                    // Value doesn't matter here and should be ignored
		status:      SelinuxNotSupported,     // Value doesn't matter here and should be ignored
		selinuxType: SelinuxTypeNotSupported, // Value doesn't matter here and should be ignored
	}

	values := selinuxInfo.toMapOfCtyValues()

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

func TestSelinuxConstants(t *testing.T) {
	// Test status constants
	statusValues := []selinuxStatus{
		SelinuxEnforcing,
		SelinuxDisabled,
		SelinuxPermissive,
		SelinuxNotSupported,
	}

	expectedStatuses := []string{"enforcing", "disabled", "permissive", ""}

	for i, status := range statusValues {
		if string(status) != expectedStatuses[i] {
			t.Errorf("expected status constant %d to be %q, got %q", i, expectedStatuses[i], string(status))
		}
	}

	// Test type constants
	typeValues := []selinuxType{
		SelinuxTypeTargeted,
		SelinuxTypeMinimum,
		SelinuxTypeMLS,
		SelinuxTypeNotSupported,
	}

	expectedTypes := []string{"targeted", "minimum", "mls", ""}

	for i, selinuxType := range typeValues {
		if string(selinuxType) != expectedTypes[i] {
			t.Errorf("expected type constant %d to be %q, got %q", i, expectedTypes[i], string(selinuxType))
		}
	}
}
