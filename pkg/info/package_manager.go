package info

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/pkg/diag"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run ../../cmd/scriptimport/main.go info package_manager_discovery.sh

type PackageManagerInfo struct {
	name string
	path string
}

func newPackageManagerInfo() *PackageManagerInfo {
	return &PackageManagerInfo{}
}

func (p *PackageManagerInfo) Name() string {
	return p.name
}

func (p *PackageManagerInfo) Path() string {
	return p.path
}

func (p *PackageManagerInfo) populatePackageManagerInfo(osInfo *OSInfo, t transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.ID() == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping package manager information collection due to missing or invalid OS info",
		}}
	}

	if osInfo.Families().Contains("windows") {
		return diag.Diags{} // Windows does not have a traditional package manager like other OS families
	}

	cmd, err := t.NewCommand(packageManagerDiscoveryScript, nil)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to create command for package manager discovery",
			Detail:   fmt.Sprintf("Error creating command: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to check package manager status",
			Detail:   fmt.Sprintf("Error checking package manager status: %v", err),
		}, &diag.Diag{
			Severity: diag.DiagDebug,
			Summary:  "Discovery command stderr",
			Detail:   fmt.Sprintf("stderr: %s", stderr),
		}}
	}

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse package manager information",
			Detail:   fmt.Sprintf("Error parsing package manager information: %v", err),
		}}
	}

	if osInfo.Families().Contains("darwin") {
		return p.populateDarwinPackageManagerInfo(discoveredData)
	}

	diags := diag.Diags{}
	if osInfo.Families().Contains("archlinux") {
		moreDiags := p.populateArchLinuxPackageManagerInfo(discoveredData)
		diags = diags.AppendAll(moreDiags)

		if p.name != "" {
			return diags
		}
	}

	if osInfo.Families().Contains("debian") || osInfo.Families().Contains("altlinux") {
		moreDiags := p.populateDebianPackageManagerInfo(discoveredData)
		diags = diags.AppendAll(moreDiags)

		if p.name != "" {
			return diags
		}
	}

	if osInfo.Families().Contains("el") {
		moreDiags := p.populateEnterpriseLinuxPackageManagerInfo(discoveredData)
		diags = diags.AppendAll(moreDiags)

		if p.name != "" {
			return diags
		}
	}

	if osInfo.Families().Contains("gentoo") {
		moreDiags := p.populateGentooPackageManagerInfo(discoveredData)
		diags = diags.AppendAll(moreDiags)

		if p.name != "" {
			return diags
		}
	}

	if osInfo.Families().Contains("suse") {
		moreDiags := p.populateSusePackageManagerInfo(discoveredData)
		diags = diags.AppendAll(moreDiags)

		if p.name != "" {
			return diags
		}
	}

	moreDiags := p.populateOtherPackageManagerInfo(discoveredData)
	diags = diags.AppendAll(moreDiags)

	return diags
}

func (p *PackageManagerInfo) populateDarwinPackageManagerInfo(data map[string]string) diag.Diags {

	optHomebrewBinBrewExists, _ := data["opt_homebrew_bin_brew_exists"]
	if optHomebrewBinBrewExists == "1" {
		p.name = "homebrew"
		p.path = "/opt/homebrew/bin/brew"
		return diag.Diags{}
	}

	usrLocalBinBrewExists, _ := data["usr_local_bin_brew_exists"]
	if usrLocalBinBrewExists == "1" {
		p.name = "homebrew"
		p.path = "/usr/local/bin/brew"
		return diag.Diags{}
	}

	optLocalBinPortExists, _ := data["opt_local_bin_port_exists"]
	if optLocalBinPortExists == "1" {
		p.name = "macports"
		p.path = "/opt/local/bin/port"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "No package manager found",
		Detail:   "No known package manager found for macOS. Please ensure Homebrew or MacPorts is installed.",
	}}
}

func (p *PackageManagerInfo) populateArchLinuxPackageManagerInfo(data map[string]string) diag.Diags {

	usrBinPacmanExists, _ := data["usr_bin_pacman_exists"]
	if usrBinPacmanExists == "1" {
		p.name = "pacman"
		p.path = "/usr/bin/pacman"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "Primary package manager not found",
		Detail:   "The primary package manager for Arch Linux (pacman) was not found",
	}}
}

func (p *PackageManagerInfo) populateDebianPackageManagerInfo(data map[string]string) diag.Diags {

	usrBinAptGetExists, _ := data["usr_bin_apt_get_exists"]
	if usrBinAptGetExists == "1" {

		aptProvidedByRpmPackage := data["apt_provided_by_rpm_package"]
		if aptProvidedByRpmPackage != "" {
			p.name = "apt-rpm"
			p.path = "/usr/bin/apt-get"
			return diag.Diags{}
		}

		p.name = "apt"
		p.path = "/usr/bin/apt-get"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "Primary package manager not found",
		Detail:   "The primary package manager for Debian or AltLinux (apt) was not found",
	}}
}

func (p *PackageManagerInfo) populateEnterpriseLinuxPackageManagerInfo(data map[string]string) diag.Diags {

	usrBinDnf5Exists, _ := data["usr_bin_dnf5_exists"]
	if usrBinDnf5Exists == "1" {
		p.name = "dnf5"
		p.path = "/usr/bin/dnf5"
		return diag.Diags{}
	}

	usrBinDnf3Exists, _ := data["usr_bin_dnf_3_exists"]
	if usrBinDnf3Exists == "1" {
		p.name = "dnf"
		p.path = "/usr/bin/dnf-3"
		return diag.Diags{}
	}

	usrBinDnfExists, _ := data["usr_bin_dnf_exists"]
	if usrBinDnfExists == "1" {
		p.name = "dnf"
		p.path = "/usr/bin/dnf"
		return diag.Diags{}
	}

	usrBinYumExists, _ := data["usr_bin_yum_exists"]
	if usrBinYumExists == "1" {
		p.name = "yum"
		p.path = "/usr/bin/yum"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "Primary package manager not found",
		Detail:   "The primary package manager for Enterprise Linux (dnf or yum) was not found",
	}}
}

func (p *PackageManagerInfo) populateGentooPackageManagerInfo(data map[string]string) diag.Diags {

	usrBinEmergeExists, _ := data["usr_bin_emerge_exists"]
	if usrBinEmergeExists == "1" {
		p.name = "portage"
		p.path = "/usr/bin/emerge"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "Primary package manager not found",
		Detail:   "The primary package manager for Gentoo (portage) was not found",
	}}
}

func (p *PackageManagerInfo) populateSusePackageManagerInfo(data map[string]string) diag.Diags {

	usrBinZypperExists, _ := data["usr_bin_zypper_exists"]
	if usrBinZypperExists == "1" {
		p.name = "zypper"
		p.path = "/usr/bin/zypper"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "Primary package manager not found",
		Detail:   "The primary package manager for SUSE (zypper) was not found",
	}}
}

func (p *PackageManagerInfo) populateOtherPackageManagerInfo(data map[string]string) diag.Diags {

	qopenSysPkgsBinYumExists, _ := data["qopensys_pkgs_bin_yum_exists"]
	if qopenSysPkgsBinYumExists == "1" {
		p.name = "yum"
		p.path = "/QOpenSys/pkgs/bin/yum"
		return diag.Diags{}
	}

	usrBinInstallpExists, _ := data["usr_bin_installp_exists"]
	if usrBinInstallpExists == "1" {
		p.name = "installp"
		p.path = "/usr/bin/installp"
		return diag.Diags{}
	}

	usrSbinSorceryExists, _ := data["usr_sbin_sorcery_exists"]
	if usrSbinSorceryExists == "1" {
		p.name = "sorcery"
		p.path = "/usr/sbin/sorcery"
		return diag.Diags{}
	}

	usrBinSwupdExists, _ := data["usr_bin_swupd_exists"]
	if usrBinSwupdExists == "1" {
		p.name = "swupd"
		p.path = "/usr/bin/swupd"
		return diag.Diags{}
	}

	usrLocalSbinPkgExists, _ := data["usr_local_sbin_pkg_exists"]
	if usrLocalSbinPkgExists == "1" {
		p.name = "pkgng"
		p.path = "/usr/local/sbin/pkg"
		return diag.Diags{}
	}

	usrBinXbpsInstallExists, _ := data["usr_bin_xbps_install_exists"]
	if usrBinXbpsInstallExists == "1" {
		p.name = "xbps"
		p.path = "/usr/bin/xbps-install"
		return diag.Diags{}
	}

	usrBinPkgExists, _ := data["usr_bin_pkg_exists"]
	if usrBinPkgExists == "1" {
		p.name = "pkg5"
		p.path = "/usr/bin/pkg"
		return diag.Diags{}
	}

	usrSbinPkgaddExists, _ := data["usr_sbin_pkgadd_exists"]
	if usrSbinPkgaddExists == "1" {
		p.name = "svr4pkg"
		p.path = "/usr/sbin/pkgadd"
		return diag.Diags{}
	}

	usrBinEmergeExists, _ := data["usr_bin_emerge_exists"]
	if usrBinEmergeExists == "1" {
		p.name = "portage"
		p.path = "/usr/bin/emerge"
		return diag.Diags{}
	}

	usrSbinSwlistExists, _ := data["usr_sbin_swlist_exists"]
	if usrSbinSwlistExists == "1" {
		p.name = "swdepot"
		p.path = "/usr/sbin/swlist"
		return diag.Diags{}
	}

	usrSbinPkgExists, _ := data["usr_sbin_pkg_exists"]
	if usrSbinPkgExists == "1" {
		p.name = "pkgng"
		p.path = "/usr/sbin/pkg"
		return diag.Diags{}
	}

	sbinApkExists, _ := data["sbin_apk_exists"]
	if sbinApkExists == "1" {
		p.name = "apk"
		p.path = "/sbin/apk"
		return diag.Diags{}
	}

	optHomebrewBinBrewExists, _ := data["opt_homebrew_bin_brew_exists"]
	if optHomebrewBinBrewExists == "1" {
		p.name = "homebrew"
		p.path = "/opt/homebrew/bin/brew"
		return diag.Diags{}
	}

	usrLocalBinBrewExists, _ := data["usr_local_bin_brew_exists"]
	if usrLocalBinBrewExists == "1" {
		p.name = "homebrew"
		p.path = "/usr/local/bin/brew"
		return diag.Diags{}
	}

	optLocalBinPortExists, _ := data["opt_local_bin_port_exists"]
	if optLocalBinPortExists == "1" {
		p.name = "macports"
		p.path = "/opt/local/bin/port"
		return diag.Diags{}
	}

	optToolsBinPkginExists, _ := data["opt_tools_bin_pkgin_exists"]
	if optToolsBinPkginExists == "1" {
		p.name = "pkgin"
		p.path = "/opt/tools/bin/pkgin"
		return diag.Diags{}
	}

	optLocalBinPkginExists, _ := data["opt_local_bin_pkgin_exists"]
	if optLocalBinPkginExists == "1" {
		p.name = "pkgin"
		p.path = "/opt/local/bin/pkgin"
		return diag.Diags{}
	}

	usrPkgBinPkginExists, _ := data["usr_pkg_bin_pkgin_exists"]
	if usrPkgBinPkginExists == "1" {
		p.name = "pkgin"
		p.path = "/usr/pkg/bin/pkgin"
		return diag.Diags{}
	}

	binOpkgExists, _ := data["bin_opkg_exists"]
	if binOpkgExists == "1" {
		p.name = "opkg"
		p.path = "/bin/opkg"
		return diag.Diags{}
	}

	usrBinPacmanExists, _ := data["usr_bin_pacman_exists"]
	if usrBinPacmanExists == "1" {
		p.name = "pacman"
		p.path = "/usr/bin/pacman"
		return diag.Diags{}
	}

	usrSbinUrpmiExists, _ := data["usr_sbin_urpmi_exists"]
	if usrSbinUrpmiExists == "1" {
		p.name = "urpmi"
		p.path = "/usr/sbin/urpmi"
		return diag.Diags{}
	}

	usrBinZypperExists, _ := data["usr_bin_zypper_exists"]
	if usrBinZypperExists == "1" {
		p.name = "zypper"
		p.path = "/usr/bin/zypper"
		return diag.Diags{}
	}

	usrBinAptGetExists, _ := data["usr_bin_apt_get_exists"]
	if usrBinAptGetExists == "1" {
		aptProvidedByRpmPackage, _ := data["apt_provided_by_rpm_package"]
		if aptProvidedByRpmPackage != "" {
			p.name = "apt-rpm"
			p.path = "/usr/bin/apt-get"
			return diag.Diags{}
		}
		p.name = "apt"
		p.path = "/usr/bin/apt-get"
		return diag.Diags{}
	}

	usrBinDnf5Exists, _ := data["usr_bin_dnf5_exists"]
	if usrBinDnf5Exists == "1" {
		p.name = "dnf5"
		p.path = "/usr/bin/dnf5"
		return diag.Diags{}
	}

	usrBinDnf3Exists, _ := data["usr_bin_dnf_3_exists"]
	if usrBinDnf3Exists == "1" {
		p.name = "dnf"
		p.path = "/usr/bin/dnf-3"
		return diag.Diags{}
	}

	usrBinDnfExists, _ := data["usr_bin_dnf_exists"]
	if usrBinDnfExists == "1" {
		p.name = "dnf"
		p.path = "/usr/bin/dnf"
		return diag.Diags{}
	}

	usrBinYumExists, _ := data["usr_bin_yum_exists"]
	if usrBinYumExists == "1" {
		p.name = "yum"
		p.path = "/usr/bin/yum"
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagWarning,
		Summary:  "No package manager found",
		Detail:   "No known package manager found",
	}}
}

func (p *PackageManagerInfo) toMapOfCtyValues() map[string]cty.Value {

	values := make(map[string]cty.Value)

	if p.name != "" {
		values["package_manager_name"] = cty.StringVal(p.name)
	} else {
		values["package_manager_name"] = cty.NullVal(cty.String)
	}

	if p.path != "" {
		values["package_manager_path"] = cty.StringVal(p.path)
	} else {
		values["package_manager_path"] = cty.NullVal(cty.String)
	}

	return values
}

// String returns a string representation of the package manager information.
// This is useful for logging or debugging purposes.
func (p *PackageManagerInfo) String() string {

	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString("package_manager_name: ")
	if p.name != "" {
		stringBuilder.WriteString(p.name)
	} else {
		stringBuilder.WriteString("unknown")
	}
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("package_manager_path: ")

	if p.path != "" {
		stringBuilder.WriteString(p.path)
	} else {
		stringBuilder.WriteString("unknown")
	}

	return stringBuilder.String()
}
