package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestPackageManagerInfo_PopulatePackageManagerInfo_NoOS(t *testing.T) {

	osInfo := newOSInfo()

	mockTransport := transport.NewMockTransport()

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Missing OS information"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Skipping package manager information collection due to missing or invalid OS info"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Windows(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "windows-server"
	osInfo.families.Add("windows")

	mockTransport := transport.NewMockTransport()

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Errorf("Expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Darwin(t *testing.T) {

	tests := []struct {
		name         string
		output       string
		expectedName string
		expectedPath string
	}{
		{
			name: "/opt/homebrew/bin/brew",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "1",
			  "usr_local_bin_brew_exists": "1",
			  "opt_local_bin_port_exists": "1",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "homebrew",
			expectedPath: "/opt/homebrew/bin/brew",
		},
		{
			name: "/usr/local/bin/brew",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "1",
			  "opt_local_bin_port_exists": "1",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "homebrew",
			expectedPath: "/usr/local/bin/brew",
		},
		{
			name: "/opt/local/bin/port",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "1",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "macports",
			expectedPath: "/opt/local/bin/port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.id = "macos"
			osInfo.families.Add("darwin")
			osInfo.families.Add("macos")

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newPackageManagerInfo()
			diags := info.populatePackageManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Errorf("Expected no error, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Errorf("Expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("Expected package manager name %q, got: %q", tt.expectedName, info.Name())
			}

			if info.Path() != tt.expectedPath {
				t.Errorf("Expected package manager path %q, got: %q", tt.expectedPath, info.Path())
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Darwin_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "macos"
	osInfo.families.Add("darwin")
	osInfo.families.Add("macos")

	mockTransport := transport.NewMockTransport()

	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Errorf("Expected warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "No package manager found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "No known package manager found for macOS. Please ensure Homebrew or MacPorts is installed."
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_ArchLinux(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "archlinux"
	osInfo.families.Add("archlinux")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "1",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Errorf("Expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "pacman"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/usr/bin/pacman"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_ArchLinux_NotPacman(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "archlinux"
	osInfo.families.Add("archlinux")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Arch Linux (pacman) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_ArchLinux_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "archlinux"
	osInfo.families.Add("archlinux")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Arch Linux (pacman) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Debian(t *testing.T) {

	tests := []struct {
		name         string
		output       string
		expectedName string
		expectedPath string
	}{
		{
			name: "apt",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "1",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "apt",
			expectedPath: "/usr/bin/apt-get",
		},
		{
			name: "apt-rpm",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "1",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": "apt-2.9.8-1.el9.x86_64"
			}
			`,
			expectedName: "apt-rpm",
			expectedPath: "/usr/bin/apt-get",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.id = "debian"
			osInfo.families.Add("debian")

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newPackageManagerInfo()
			diags := info.populatePackageManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Errorf("Expected no error, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Errorf("Expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("Expected package manager name %q, got: %q", tt.expectedName, info.Name())
			}

			if info.Path() != tt.expectedPath {
				t.Errorf("Expected package manager path %q, got: %q", tt.expectedPath, info.Path())
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Debian_NotApt(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "debian"
	osInfo.families.Add("debian")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Debian or AltLinux (apt) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Debian_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "debian"
	osInfo.families.Add("debian")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Debian or AltLinux (apt) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_AltLinux(t *testing.T) {

	tests := []struct {
		name         string
		output       string
		expectedName string
		expectedPath string
	}{
		{
			name: "apt",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "1",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "apt",
			expectedPath: "/usr/bin/apt-get",
		},
		{
			name: "apt-rpm",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "1",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": "apt-2.9.8-1.el9.x86_64"
			}
			`,
			expectedName: "apt-rpm",
			expectedPath: "/usr/bin/apt-get",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.id = "altlinux"
			osInfo.families.Add("altlinux")

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newPackageManagerInfo()
			diags := info.populatePackageManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Errorf("Expected no error, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Errorf("Expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("Expected package manager name %q, got: %q", tt.expectedName, info.Name())
			}

			if info.Path() != tt.expectedPath {
				t.Errorf("Expected package manager path %q, got: %q", tt.expectedPath, info.Path())
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_AltLinux_NotApt(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "altlinux"
	osInfo.families.Add("altlinux")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Debian or AltLinux (apt) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_AltLinux_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "altlinux"
	osInfo.families.Add("altlinux")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Debian or AltLinux (apt) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_EL(t *testing.T) {

	tests := []struct {
		name         string
		output       string
		expectedName string
		expectedPath string
	}{
		{
			name: "dnf5",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "1",
			  "usr_bin_dnf_3_exists": "1",
			  "usr_bin_dnf_exists": "1",
			  "usr_bin_yum_exists": "1",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "dnf5",
			expectedPath: "/usr/bin/dnf5",
		},
		{
			name: "dnf-3",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "1",
			  "usr_bin_dnf_exists": "1",
			  "usr_bin_yum_exists": "1",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "dnf",
			expectedPath: "/usr/bin/dnf-3",
		},
		{
			name: "dnf",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "1",
			  "usr_bin_yum_exists": "1",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "dnf",
			expectedPath: "/usr/bin/dnf",
		},
		{
			name: "yum",
			output: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "1",
			  "apt_provided_by_rpm_package": ""
			}
			`,
			expectedName: "yum",
			expectedPath: "/usr/bin/yum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.id = "rhel"
			osInfo.families.Add("el")

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newPackageManagerInfo()
			diags := info.populatePackageManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Errorf("Expected no error, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Errorf("Expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("Expected package manager name %q, got: %q", tt.expectedName, info.Name())
			}

			if info.Path() != tt.expectedPath {
				t.Errorf("Expected package manager path %q, got: %q", tt.expectedPath, info.Path())
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_EL_NotDnfOrYum(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "rhel"
	osInfo.families.Add("el")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Enterprise Linux (dnf or yum) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_EL_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "rhel"
	osInfo.families.Add("el")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Enterprise Linux (dnf or yum) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Gentoo(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "gentoo"
	osInfo.families.Add("gentoo")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "1",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Errorf("Expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "portage"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/usr/bin/emerge"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Gentoo_NotPortage(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "gentoo"
	osInfo.families.Add("gentoo")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Gentoo (portage) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Gentoo_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "gentoo"
	osInfo.families.Add("gentoo")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for Gentoo (portage) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_SUSE(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "suse"
	osInfo.families.Add("suse")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "1",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if diags.HasWarnings() {
		t.Errorf("Expected no warnings, got: %v", diags.Warnings())
	}

	expectedName := "zypper"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/usr/bin/zypper"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_SUSE_NotZypper(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "suse"
	osInfo.families.Add("suse")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "1",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	expectedName := "yum"
	if info.Name() != expectedName {
		t.Errorf("Expected package manager name %q, got: %q", expectedName, info.Name())
	}

	expectedPath := "/QOpenSys/pkgs/bin/yum"
	if info.Path() != expectedPath {
		t.Errorf("Expected package manager path %q, got: %q", expectedPath, info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for SUSE (zypper) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_SUSE_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "suse"
	osInfo.families.Add("suse")

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 2 {
		t.Fatalf("Expected 2 warnings, got: %d", len(warnings))
	}

	expectedSummary := "Primary package manager not found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "The primary package manager for SUSE (zypper) was not found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}

	expectedSummary = "No package manager found"
	if warnings[1].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[1].Summary)
	}

	expectedDetail = "No known package manager found"
	if warnings[1].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[1].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Generic(t *testing.T) {

	tests := []struct {
		name         string
		output       string
		expectedName string
		expectedPath string
	}{
		{
			name: "/QOpenSys/pkgs/bin/yum",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "1",
				"usr_bin_installp_exists": "1",
				"usr_sbin_sorcery_exists": "1",
				"usr_bin_swupd_exists": "1",
				"usr_local_sbin_pkg_exists": "1",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "yum",
			expectedPath: "/QOpenSys/pkgs/bin/yum",
		},
		{
			name: "/usr/bin/installp",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "1",
				"usr_sbin_sorcery_exists": "1",
				"usr_bin_swupd_exists": "1",
				"usr_local_sbin_pkg_exists": "1",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "installp",
			expectedPath: "/usr/bin/installp",
		},
		{
			name: "/usr/sbin/sorcery",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "1",
				"usr_bin_swupd_exists": "1",
				"usr_local_sbin_pkg_exists": "1",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "sorcery",
			expectedPath: "/usr/sbin/sorcery",
		},
		{
			name: "/usr/bin/swupd",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "1",
				"usr_local_sbin_pkg_exists": "1",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "swupd",
			expectedPath: "/usr/bin/swupd",
		},
		{
			name: "/usr/local/sbin/pkg",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "1",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkgng",
			expectedPath: "/usr/local/sbin/pkg",
		},
		{
			name: "/usr/bin/xbps-install",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "1",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "xbps",
			expectedPath: "/usr/bin/xbps-install",
		},
		{
			name: "/usr/bin/pkg",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "1",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkg5",
			expectedPath: "/usr/bin/pkg",
		},
		{
			name: "/usr/sbin/pkgadd",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "1",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "svr4pkg",
			expectedPath: "/usr/sbin/pkgadd",
		},
		{
			name: "/usr/bin/emerge",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "1",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "portage",
			expectedPath: "/usr/bin/emerge",
		},
		{
			name: "/usr/sbin/swlist",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "1",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "swdepot",
			expectedPath: "/usr/sbin/swlist",
		},
		{
			name: "/usr/sbin/pkg",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "1",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkgng",
			expectedPath: "/usr/sbin/pkg",
		},
		{
			name: "/sbin/apk",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "1",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "apk",
			expectedPath: "/sbin/apk",
		},
		{
			name: "/opt/homebrew/bin/brew",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "1",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "homebrew",
			expectedPath: "/opt/homebrew/bin/brew",
		},
		{
			name: "/usr/local/bin/brew",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "1",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "homebrew",
			expectedPath: "/usr/local/bin/brew",
		},
		{
			name: "/opt/local/bin/port",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "1",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "macports",
			expectedPath: "/opt/local/bin/port",
		},
		{
			name: "/opt/tools/bin/pkgin",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "1",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkgin",
			expectedPath: "/opt/tools/bin/pkgin",
		},
		{
			name: "/opt/local/bin/pkgin",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "1",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkgin",
			expectedPath: "/opt/local/bin/pkgin",
		},
		{
			name: "/usr/pkg/bin/pkgin",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "1",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pkgin",
			expectedPath: "/usr/pkg/bin/pkgin",
		},
		{
			name: "/bin/opkg",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "1",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "opkg",
			expectedPath: "/bin/opkg",
		},
		{
			name: "/usr/bin/pacman",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "1",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "pacman",
			expectedPath: "/usr/bin/pacman",
		},
		{
			name: "/usr/sbin/urpmi",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "1",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "urpmi",
			expectedPath: "/usr/sbin/urpmi",
		},
		{
			name: "/usr/bin/zypper",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "1",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "zypper",
			expectedPath: "/usr/bin/zypper",
		},
		{
			name: "/usr/bin/apt-get",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "apt",
			expectedPath: "/usr/bin/apt-get",
		},
		{
			name: "/usr/bin/apt-get (apt-rpm)",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "1",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": "apt-0.5.15lorg3.95a"
				}
				`,
			expectedName: "apt-rpm",
			expectedPath: "/usr/bin/apt-get",
		},
		{
			name: "/usr/bin/dnf5",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "0",
				"usr_bin_dnf5_exists": "1",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "dnf5",
			expectedPath: "/usr/bin/dnf5",
		},
		{
			name: "/usr/bin/dnf-3",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "0",
				"usr_bin_dnf5_exists": "0",
				"usr_bin_dnf_3_exists": "1",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "dnf",
			expectedPath: "/usr/bin/dnf-3",
		},
		{
			name: "/usr/bin/dnf",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "0",
				"usr_bin_dnf5_exists": "0",
				"usr_bin_dnf_3_exists": "0",
				"usr_bin_dnf_exists": "1",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "dnf",
			expectedPath: "/usr/bin/dnf",
		},
		{
			name: "/usr/bin/yum",
			output: `{
				"qopensys_pkgs_bin_yum_exists": "0",
				"usr_bin_installp_exists": "0",
				"usr_sbin_sorcery_exists": "0",
				"usr_bin_swupd_exists": "0",
				"usr_local_sbin_pkg_exists": "0",
				"usr_bin_xbps_install_exists": "0",
				"usr_bin_pkg_exists": "0",
				"usr_sbin_pkgadd_exists": "0",
				"usr_bin_emerge_exists": "0",
				"usr_sbin_swlist_exists": "0",
				"usr_sbin_pkg_exists": "0",
				"sbin_apk_exists": "0",
				"opt_homebrew_bin_brew_exists": "0",
				"usr_local_bin_brew_exists": "0",
				"opt_local_bin_port_exists": "0",
				"opt_tools_bin_pkgin_exists": "0",
				"opt_local_bin_pkgin_exists": "0",
				"usr_pkg_bin_pkgin_exists": "0",
				"bin_opkg_exists": "0",
				"usr_bin_pacman_exists": "0",
				"usr_sbin_urpmi_exists": "0",
				"usr_bin_zypper_exists": "0",
				"usr_bin_apt_get_exists": "0",
				"usr_bin_dnf5_exists": "0",
				"usr_bin_dnf_3_exists": "0",
				"usr_bin_dnf_exists": "0",
				"usr_bin_yum_exists": "1",
				"apt_provided_by_rpm_package": ""
				}
				`,
			expectedName: "yum",
			expectedPath: "/usr/bin/yum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()
			osInfo.id = "generic"

			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newPackageManagerInfo()
			diags := info.populatePackageManagerInfo(osInfo, mockTransport)

			if diags.HasErrors() {
				t.Errorf("Expected no error, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Errorf("Expected no warnings, got: %v", diags.Warnings())
			}

			if info.Name() != tt.expectedName {
				t.Errorf("Expected package manager name %q, got: %q", tt.expectedName, info.Name())
			}

			if info.Path() != tt.expectedPath {
				t.Errorf("Expected package manager path %q, got: %q", tt.expectedPath, info.Path())
			}
		})
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Generic_NotFound(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "generic"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			  "qopensys_pkgs_bin_yum_exists": "0",
		      "usr_bin_installp_exists": "0",
		      "usr_sbin_sorcery_exists": "0",
		      "usr_bin_swupd_exists": "0",
		      "usr_local_sbin_pkg_exists": "0",
		      "usr_bin_xbps_install_exists": "0",
		      "usr_bin_pkg_exists": "0",
		      "usr_sbin_pkgadd_exists": "0",
		      "usr_bin_emerge_exists": "0",
		      "usr_sbin_swlist_exists": "0",
		      "usr_sbin_pkg_exists": "0",
		      "sbin_apk_exists": "0",
			  "opt_homebrew_bin_brew_exists": "0",
			  "usr_local_bin_brew_exists": "0",
			  "opt_local_bin_port_exists": "0",
			  "opt_tools_bin_pkgin_exists": "0",
			  "opt_local_bin_pkgin_exists": "0",
			  "usr_pkg_bin_pkgin_exists": "0",
			  "bin_opkg_exists": "0",
			  "usr_bin_pacman_exists": "0",
			  "usr_sbin_urpmi_exists": "0",
			  "usr_bin_zypper_exists": "0",
			  "usr_bin_apt_get_exists": "0",
			  "usr_bin_dnf5_exists": "0",
			  "usr_bin_dnf_3_exists": "0",
			  "usr_bin_dnf_exists": "0",
			  "usr_bin_yum_exists": "0",
			  "apt_provided_by_rpm_package": ""
			}
			`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if diags.HasErrors() {
		t.Errorf("Expected no error, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Error("Expected warnings, got none")
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("Expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "No package manager found"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("Expected summary %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "No known package manager found"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("Expected detail %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Generic_Error(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "generic"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Errorf("Expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("Expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to check package manager status"
	if errors[0].Summary != expectedSummary {
		t.Errorf("Expected error summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error checking package manager status: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("Expected error detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestPackageManagerInfo_PopulatePackageManagerInfo_Generic_NotJSON(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.id = "generic"

	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults[packageManagerDiscoveryScript] = &transport.MockCmd{
		Stdout: `This is not JSON output`,
	}

	info := newPackageManagerInfo()
	diags := info.populatePackageManagerInfo(osInfo, mockTransport)

	if !diags.HasErrors() {
		t.Errorf("Expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("Expected no warnings, got: %v", diags.Warnings())
	}

	if info.Name() != "" {
		t.Errorf("Expected empty package manager name, got: %q", info.Name())
	}

	if info.Path() != "" {
		t.Errorf("Expected empty package manager path, got: %q", info.Path())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("Expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse package manager information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("Expected error summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing package manager information: invalid character 'T' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("Expected error detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestPackageManagerInfo_ToMapOfCtyValues(t *testing.T) {

	tests := []struct {
		name               string
		packageManagerName string
		packageManagerPath string
	}{
		{
			name:               "Homebrew",
			packageManagerName: "homebrew",
			packageManagerPath: "/opt/homebrew/bin/brew",
		},
		{
			name:               "MacPorts",
			packageManagerName: "macports",
			packageManagerPath: "/opt/local/bin/port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			info := newPackageManagerInfo()
			info.name = tt.packageManagerName
			info.path = tt.packageManagerPath

			result := info.toMapOfCtyValues()

			if result["package_manager_name"].AsString() != tt.packageManagerName {
				t.Errorf("Expected package_manager_name to be %q, got: %q", tt.packageManagerName, result["package_manager_name"].AsString())
			}

			if result["package_manager_path"].AsString() != tt.packageManagerPath {
				t.Errorf("Expected package_manager_path to be %q, got: %q", tt.packageManagerPath, result["package_manager_path"].AsString())
			}
		})
	}
}

func TestPackageManagerInfo_ToMapOfCtyValues_Empty(t *testing.T) {

	info := newPackageManagerInfo()

	result := info.toMapOfCtyValues()

	if result["package_manager_name"].Type() != cty.String {
		t.Errorf("Expected package_manager_name to be of type string, got: %q", result["package_manager_name"].Type())
	}

	if !result["package_manager_name"].IsNull() {
		t.Errorf("Expected package_manager_name to be null, got: %q", result["package_manager_name"])
	}

	if result["package_manager_path"].Type() != cty.String {
		t.Errorf("Expected package_manager_path to be of type string, got: %q", result["package_manager_path"].Type())
	}

	if !result["package_manager_path"].IsNull() {
		t.Errorf("Expected package_manager_path to be null, got: %q", result["package_manager_path"])
	}
}
