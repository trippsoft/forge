package info

import (
	"testing"

	"github.com/trippsoft/forge/internal/util"
)

func TestPackageManagerInfo_PopulatePackageManagerInfo_Windows(t *testing.T) {
	osInfo := &osInfo{
		families: util.NewSet("windows"),
	}

	mockTransport := &mockTransport{}
	mockFS := &mockFileSystem{
		files: make(map[string]*mockFile),
		dirs:  make(map[string]*mockFileInfo),
	}

	pmInfo := newPackageManagerInfo()
	err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pmInfo.supported {
		t.Error("expected package manager to be unsupported on Windows")
	}

	if pmInfo.name != "" {
		t.Errorf("expected empty name, got %q", pmInfo.name)
	}

	if pmInfo.path != "" {
		t.Errorf("expected empty path, got %q", pmInfo.path)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_EL(t *testing.T) {
	testCases := []struct {
		name       string
		packageMgr string
		expectName string
		osTreeBoot bool
	}{
		{
			name:       "dnf5",
			packageMgr: "/usr/bin/dnf5",
			expectName: "dnf5",
			osTreeBoot: false,
		},
		{
			name:       "dnf",
			packageMgr: "/usr/bin/dnf",
			expectName: "dnf",
			osTreeBoot: false,
		},
		{
			name:       "yum",
			packageMgr: "/usr/bin/yum",
			expectName: "yum",
			osTreeBoot: false,
		},
		{
			name:       "ostree_booted",
			packageMgr: "/usr/bin/dnf",
			expectName: "",
			osTreeBoot: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osInfo := &osInfo{
				families: util.NewSet("linux", "el"),
			}

			mockFS := &mockFileSystem{
				files: make(map[string]*mockFile),
				dirs:  make(map[string]*mockFileInfo),
			}

			// Add package manager file
			if tc.packageMgr != "" {
				mockFS.files[tc.packageMgr] = &mockFile{
					info: &mockFileInfo{
						name:  "packagemgr",
						isDir: false,
						mode:  0755,
					},
				}
			}

			// Add or omit ostree-booted file
			if tc.osTreeBoot {
				mockFS.files["/run/ostree-booted"] = &mockFile{
					info: &mockFileInfo{
						name:  "ostree-booted",
						isDir: false,
						mode:  0644,
					},
				}
			}

			mockTransport := &mockTransport{}

			pmInfo := newPackageManagerInfo()
			err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.osTreeBoot {
				if pmInfo.supported {
					t.Error("expected package manager to be unsupported on OSTree system")
				}
				return
			}

			if tc.expectName == "" {
				if pmInfo.supported {
					t.Error("expected package manager to be unsupported")
				}
				return
			}

			if !pmInfo.supported {
				t.Error("expected package manager to be supported")
			}

			if pmInfo.name != tc.expectName {
				t.Errorf("expected name %q, got %q", tc.expectName, pmInfo.name)
			}

			if pmInfo.path != tc.packageMgr {
				t.Errorf("expected path %q, got %q", tc.packageMgr, pmInfo.path)
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Debian(t *testing.T) {
	osInfo := &osInfo{
		families: util.NewSet("linux", "debian"),
	}

	mockFS := &mockFileSystem{
		files: map[string]*mockFile{
			"/usr/bin/apt-get": {
				info: &mockFileInfo{
					name:  "apt-get",
					isDir: false,
					mode:  0755,
				},
			},
		},
		dirs: make(map[string]*mockFileInfo),
	}

	mockTransport := &mockTransport{}

	pmInfo := newPackageManagerInfo()
	err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pmInfo.supported {
		t.Error("expected package manager to be supported")
	}

	if pmInfo.name != "apt" {
		t.Errorf("expected name 'apt', got %q", pmInfo.name)
	}

	if pmInfo.path != "/usr/bin/apt-get" {
		t.Errorf("expected path '/usr/bin/apt-get', got %q", pmInfo.path)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_AltLinux(t *testing.T) {
	osInfo := &osInfo{
		families: util.NewSet("linux", "altlinux"),
	}

	mockFS := &mockFileSystem{
		files: map[string]*mockFile{
			"/usr/bin/apt-get": {
				info: &mockFileInfo{
					name:  "apt-get",
					isDir: false,
					mode:  0755,
				},
			},
		},
		dirs: make(map[string]*mockFileInfo),
	}

	mockTransport := &mockTransport{}

	pmInfo := newPackageManagerInfo()
	err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pmInfo.supported {
		t.Error("expected package manager to be supported")
	}

	if pmInfo.name != "apt-rpm" {
		t.Errorf("expected name 'apt-rpm', got %q", pmInfo.name)
	}

	if pmInfo.path != "/usr/bin/apt-get" {
		t.Errorf("expected path '/usr/bin/apt-get', got %q", pmInfo.path)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Generic(t *testing.T) {
	testCases := []struct {
		name       string
		packageMgr string
		expectName string
		families   []string
	}{
		{
			name:       "pacman",
			packageMgr: "/usr/bin/pacman",
			expectName: "pacman",
			families:   []string{"linux", "archlinux"},
		},
		{
			name:       "zypper",
			packageMgr: "/usr/bin/zypper",
			expectName: "zypper",
			families:   []string{"linux", "suse"},
		},
		{
			name:       "emerge",
			packageMgr: "/usr/bin/emerge",
			expectName: "portage",
			families:   []string{"linux", "gentoo"},
		},
		{
			name:       "homebrew",
			packageMgr: "/opt/homebrew/bin/brew",
			expectName: "homebrew",
			families:   []string{"posix", "darwin"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osInfo := &osInfo{
				families: util.NewSet(tc.families...),
			}

			mockFS := &mockFileSystem{
				files: map[string]*mockFile{
					tc.packageMgr: {
						info: &mockFileInfo{
							name:  "packagemgr",
							isDir: false,
							mode:  0755,
						},
					},
				},
				dirs: make(map[string]*mockFileInfo),
			}

			// Test RPM detection for apt
			if tc.expectName == "apt" {
				mockFS.files["/usr/bin/rpm"] = &mockFile{
					info: &mockFileInfo{
						name:  "rpm",
						isDir: false,
						mode:  0755,
					},
				}
			}

			mockTransport := &mockTransport{
				commandOutput: "",
				shouldError:   tc.expectName == "apt", // Simulate RPM check failure for apt
			}

			pmInfo := newPackageManagerInfo()
			err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !pmInfo.supported {
				t.Error("expected package manager to be supported")
			}

			if pmInfo.name != tc.expectName {
				t.Errorf("expected name %q, got %q", tc.expectName, pmInfo.name)
			}

			if pmInfo.path != tc.packageMgr {
				t.Errorf("expected path %q, got %q", tc.packageMgr, pmInfo.path)
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_AptRpmBacked(t *testing.T) {
	osInfo := &osInfo{
		families: util.NewSet("linux", "unknown"),
	}

	mockFS := &mockFileSystem{
		files: map[string]*mockFile{
			"/usr/bin/apt-get": {
				info: &mockFileInfo{
					name:  "apt-get",
					isDir: false,
					mode:  0755,
				},
			},
			"/usr/bin/rpm": {
				info: &mockFileInfo{
					name:  "rpm",
					isDir: false,
					mode:  0755,
				},
			},
		},
		dirs: make(map[string]*mockFileInfo),
	}

	mockTransport := &mockTransport{
		commandOutput: "some-package-provides-apt",
		shouldError:   false, // Simulate successful RPM check
	}

	pmInfo := newPackageManagerInfo()
	err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pmInfo.supported {
		t.Error("expected package manager to be supported")
	}

	if pmInfo.name != "apt-rpm" {
		t.Errorf("expected name 'apt-rpm', got %q", pmInfo.name)
	}

	if pmInfo.path != "/usr/bin/apt-get" {
		t.Errorf("expected path '/usr/bin/apt-get', got %q", pmInfo.path)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_NoPackageManager(t *testing.T) {
	osInfo := &osInfo{
		families: util.NewSet("linux"),
	}

	mockFS := &mockFileSystem{
		files: make(map[string]*mockFile),
		dirs:  make(map[string]*mockFileInfo),
	}

	mockTransport := &mockTransport{}

	pmInfo := newPackageManagerInfo()
	err := pmInfo.populatePackageManagerInfo(osInfo, mockTransport, mockFS)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pmInfo.supported {
		t.Error("expected package manager to be unsupported when no package managers found")
	}

	if pmInfo.name != "" {
		t.Errorf("expected empty name, got %q", pmInfo.name)
	}

	if pmInfo.path != "" {
		t.Errorf("expected empty path, got %q", pmInfo.path)
	}
}

func TestPackageManagerInfo_GetFirstMatchingPackageManager(t *testing.T) {
	testCases := []struct {
		name            string
		managers        []string
		existingFiles   map[string]bool
		existingDirs    map[string]bool
		expectedManager string
		shouldError     bool
	}{
		{
			name:            "first_match",
			managers:        []string{"/usr/bin/dnf", "/usr/bin/yum"},
			existingFiles:   map[string]bool{"/usr/bin/dnf": true},
			expectedManager: "/usr/bin/dnf",
		},
		{
			name:            "second_match",
			managers:        []string{"/usr/bin/dnf", "/usr/bin/yum"},
			existingFiles:   map[string]bool{"/usr/bin/yum": true},
			expectedManager: "/usr/bin/yum",
		},
		{
			name:            "no_match",
			managers:        []string{"/usr/bin/dnf", "/usr/bin/yum"},
			existingFiles:   map[string]bool{},
			expectedManager: "",
		},
		{
			name:            "skip_directories",
			managers:        []string{"/usr/bin/dnf", "/usr/bin/yum"},
			existingDirs:    map[string]bool{"/usr/bin/dnf": true},
			existingFiles:   map[string]bool{"/usr/bin/yum": true},
			expectedManager: "/usr/bin/yum",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFS := &mockFileSystem{
				files: make(map[string]*mockFile),
				dirs:  make(map[string]*mockFileInfo),
			}

			// Add files
			for path := range tc.existingFiles {
				mockFS.files[path] = &mockFile{
					info: &mockFileInfo{
						name:  "file",
						isDir: false,
						mode:  0755,
					},
				}
			}

			// Add directories
			for path := range tc.existingDirs {
				mockFS.dirs[path] = &mockFileInfo{
					name:  "dir",
					isDir: true,
					mode:  0755,
				}
			}

			pmInfo := newPackageManagerInfo()
			result, err := pmInfo.getFirstMatchingPackageManager(mockFS, tc.managers)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tc.expectedManager {
				t.Errorf("expected %q, got %q", tc.expectedManager, result)
			}
		})
	}
}

func TestPackageManagerInfo_IsOSTreeBooted(t *testing.T) {
	testCases := []struct {
		name        string
		fileExists  bool
		shouldError bool
		expected    bool
	}{
		{
			name:       "ostree_booted",
			fileExists: true,
			expected:   true,
		},
		{
			name:       "not_ostree_booted",
			fileExists: false,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFS := &mockFileSystem{
				files: make(map[string]*mockFile),
				dirs:  make(map[string]*mockFileInfo),
			}

			if tc.fileExists {
				mockFS.files["/run/ostree-booted"] = &mockFile{
					info: &mockFileInfo{
						name:  "ostree-booted",
						isDir: false,
						mode:  0644,
					},
				}
			}

			result, err := isOSTreeBooted(mockFS)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestPackageManagerInfo_ToMapOfCtyValues(t *testing.T) {
	testCases := []struct {
		name      string
		supported bool
		pmName    string
		pmPath    string
	}{
		{
			name:      "supported",
			supported: true,
			pmName:    "dnf",
			pmPath:    "/usr/bin/dnf",
		},
		{
			name:      "not_supported",
			supported: false,
			pmName:    "",
			pmPath:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pmInfo := &packageManagerInfo{
				supported: tc.supported,
				name:      tc.pmName,
				path:      tc.pmPath,
			}

			values := pmInfo.toMapOfCtyValues()

			expectedKeys := []string{"package_manager_name", "package_manager_path"}
			for _, key := range expectedKeys {
				if _, exists := values[key]; !exists {
					t.Errorf("expected key %q to be present in values map", key)
				}
			}

			if tc.supported {
				if values["package_manager_name"].AsString() != tc.pmName {
					t.Errorf("expected name %q, got %q", tc.pmName, values["package_manager_name"].AsString())
				}

				if values["package_manager_path"].AsString() != tc.pmPath {
					t.Errorf("expected path %q, got %q", tc.pmPath, values["package_manager_path"].AsString())
				}
			} else {
				if !values["package_manager_name"].IsNull() {
					t.Error("expected package_manager_name to be null for unsupported")
				}

				if !values["package_manager_path"].IsNull() {
					t.Error("expected package_manager_path to be null for unsupported")
				}
			}
		})
	}
}

func TestNewPackageManagerInfo(t *testing.T) {
	pmInfo := newPackageManagerInfo()

	if pmInfo == nil {
		t.Fatal("expected non-nil packageManagerInfo")
	}

	if pmInfo.supported {
		t.Error("expected supported to be false initially")
	}

	if pmInfo.name != "" {
		t.Errorf("expected empty name initially, got %q", pmInfo.name)
	}

	if pmInfo.path != "" {
		t.Errorf("expected empty path initially, got %q", pmInfo.path)
	}
}

func TestPackageManagerConstants(t *testing.T) {
	// Test that packageManagerMap has entries for all expected package managers
	expectedPackageManagers := []string{
		"/usr/bin/dnf", "/usr/bin/yum", "/usr/bin/apt-get", "/usr/bin/pacman",
		"/usr/bin/zypper", "/usr/bin/emerge", "/opt/homebrew/bin/brew",
	}

	for _, pm := range expectedPackageManagers {
		if _, exists := packageManagerMap[pm]; !exists {
			t.Errorf("expected package manager %q to be in packageManagerMap", pm)
		}
	}

	// Test that EL package managers are all in the main map
	for _, pm := range elPackageManagers {
		if _, exists := packageManagerMap[pm]; !exists {
			t.Errorf("expected EL package manager %q to be in packageManagerMap", pm)
		}
	}

	// Test that Debian package managers are all in the main map
	for _, pm := range debianPackageManagers {
		if _, exists := packageManagerMap[pm]; !exists {
			t.Errorf("expected Debian package manager %q to be in packageManagerMap", pm)
		}
	}
}
