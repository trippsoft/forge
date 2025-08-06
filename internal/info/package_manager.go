package info

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/pkg/diag"
	"github.com/zclconf/go-cty/cty"
)

const (
	packageManagerDiscoveryScript = `qopensys_pkgs_bin_yum_exists="0"; ` +
		`usr_bin_installp_exists="0"; ` +
		`usr_sbin_sorcery_exists="0"; ` +
		`usr_bin_swupd_exists="0"; ` +
		`usr_local_sbin_pkg_exists="0"; ` +
		`usr_bin_xbps_install_exists="0"; ` +
		`usr_bin_pkg_exists="0"; ` +
		`usr_sbin_pkgadd_exists="0"; ` +
		`usr_bin_emerge_exists="0"; ` +
		`usr_sbin_swlist_exists="0"; ` +
		`usr_sbin_pkg_exists="0"; ` +
		`sbin_apk_exists="0"; ` +
		`opt_homebrew_bin_brew_exists="0"; ` +
		`usr_local_bin_brew_exists="0"; ` +
		`opt_local_bin_port_exists="0"; ` +
		`opt_tools_bin_pkgin_exists="0"; ` +
		`opt_local_bin_pkgin_exists="0"; ` +
		`usr_pkg_bin_pkgin_exists="0"; ` +
		`bin_opkg_exists="0"; ` +
		`usr_bin_pacman_exists="0"; ` +
		`usr_sbin_urpmi_exists="0"; ` +
		`usr_bin_zypper_exists="0"; ` +
		`usr_bin_apt_get_exists="0"; ` +
		`usr_bin_dnf5_exists="0"; ` +
		`usr_bin_dnf_3_exists="0"; ` +
		`usr_bin_dnf_exists="0"; ` +
		`usr_bin_yum_exists="0"; ` +
		`apt_provided_by_rpm_package=""; ` +
		`if [ -x /QOpenSys/pkgs/bin/yum ]; ` +
		`then qopensys_pkgs_bin_yum_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/installp ]; ` +
		`then usr_bin_installp_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/sbin/sorcery ]; ` +
		`then usr_sbin_sorcery_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/swupd ]; ` +
		`then usr_bin_swupd_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/local/sbin/pkg ]; ` +
		`then usr_local_sbin_pkg_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/xbps-install ]; ` +
		`then usr_bin_xbps_install_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/pkg ]; ` +
		`then usr_bin_pkg_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/sbin/pkgadd ]; ` +
		`then usr_sbin_pkgadd_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/emerge ]; ` +
		`then usr_bin_emerge_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/sbin/swlist ]; ` +
		`then usr_sbin_swlist_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/sbin/pkg ]; ` +
		`then usr_sbin_pkg_exists="1"; ` +
		`fi; ` +
		`if [ -x /sbin/apk ]; ` +
		`then sbin_apk_exists="1"; ` +
		`fi; ` +
		`if [ -x /opt/homebrew/bin/brew ]; ` +
		`then opt_homebrew_bin_brew_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/local/bin/brew ]; ` +
		`then usr_local_bin_brew_exists="1"; ` +
		`fi; ` +
		`if [ -x /opt/local/bin/port ]; ` +
		`then opt_local_bin_port_exists="1"; ` +
		`fi; ` +
		`if [ -x /opt/tools/bin/pkgin ]; ` +
		`then opt_tools_bin_pkgin_exists="1"; ` +
		`fi; ` +
		`if [ -x /opt/local/bin/pkgin ]; ` +
		`then opt_local_bin_pkgin_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/pkg/bin/pkgin ]; ` +
		`then usr_pkg_bin_pkgin_exists="1"; ` +
		`fi; ` +
		`if [ -x /bin/opkg ]; ` +
		`then bin_opkg_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/pacman ]; ` +
		`then usr_bin_pacman_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/sbin/urpmi ]; ` +
		`then usr_sbin_urpmi_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/zypper ]; ` +
		`then usr_bin_zypper_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/apt-get ]; ` +
		`then usr_bin_apt_get_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/dnf5 ]; ` +
		`then usr_bin_dnf5_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/dnf-3 ]; ` +
		`then usr_bin_dnf_3_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/dnf ]; ` +
		`then usr_bin_dnf_exists="1"; ` +
		`fi; ` +
		`if [ -x /usr/bin/yum ]; ` +
		`then usr_bin_yum_exists="1"; ` +
		`fi; ` +
		`if [ "$usr_bin_apt_get_exists" -eq 1 ] && [ -x /usr/bin/rpm ]; ` +
		`then apt_provided_by_rpm_package=$(/usr/bin/rpm -q --whatprovides /usr/bin/apt-get || echo ""); ` +
		`fi; ` +
		`output=$(jq -n ` +
		`--arg qopensys_pkgs_bin_yum_exists "$qopensys_pkgs_bin_yum_exists" ` +
		`--arg usr_bin_installp_exists "$usr_bin_installp_exists" ` +
		`--arg usr_sbin_sorcery_exists "$usr_sbin_sorcery_exists" ` +
		`--arg usr_bin_swupd_exists "$usr_bin_swupd_exists" ` +
		`--arg usr_local_sbin_pkg_exists "$usr_local_sbin_pkg_exists" ` +
		`--arg usr_bin_xbps_install_exists "$usr_bin_xbps_install_exists" ` +
		`--arg usr_bin_pkg_exists "$usr_bin_pkg_exists" ` +
		`--arg usr_sbin_pkgadd_exists "$usr_sbin_pkgadd_exists" ` +
		`--arg usr_bin_emerge_exists "$usr_bin_emerge_exists" ` +
		`--arg usr_sbin_swlist_exists "$usr_sbin_swlist_exists" ` +
		`--arg usr_sbin_pkg_exists "$usr_sbin_pkg_exists" ` +
		`--arg sbin_apk_exists "$sbin_apk_exists" ` +
		`--arg opt_homebrew_bin_brew_exists "$opt_homebrew_bin_brew_exists" ` +
		`--arg usr_local_bin_brew_exists "$usr_local_bin_brew_exists" ` +
		`--arg opt_local_bin_port_exists "$opt_local_bin_port_exists" ` +
		`--arg opt_tools_bin_pkgin_exists "$opt_tools_bin_pkgin_exists" ` +
		`--arg opt_local_bin_pkgin_exists "$opt_local_bin_pkgin_exists" ` +
		`--arg usr_pkg_bin_pkgin_exists "$usr_pkg_bin_pkgin_exists" ` +
		`--arg bin_opkg_exists "$bin_opkg_exists" ` +
		`--arg usr_bin_pacman_exists "$usr_bin_pacman_exists" ` +
		`--arg usr_sbin_urpmi_exists "$usr_sbin_urpmi_exists" ` +
		`--arg usr_bin_zypper_exists "$usr_bin_zypper_exists" ` +
		`--arg usr_bin_apt_get_exists "$usr_bin_apt_get_exists" ` +
		`--arg usr_bin_dnf5_exists "$usr_bin_dnf5_exists" ` +
		`--arg usr_bin_dnf_3_exists "$usr_bin_dnf_3_exists" ` +
		`--arg usr_bin_dnf_exists "$usr_bin_dnf_exists" ` +
		`--arg usr_bin_yum_exists "$usr_bin_yum_exists" ` +
		`--arg apt_provided_by_rpm_package "$apt_provided_by_rpm_package" ` +
		`'{` +
		`qopensys_pkgs_bin_yum_exists: $qopensys_pkgs_bin_yum_exists, ` +
		`usr_bin_installp_exists: $usr_bin_installp_exists, ` +
		`usr_sbin_sorcery_exists: $usr_sbin_sorcery_exists, ` +
		`usr_bin_swupd_exists: $usr_bin_swupd_exists, ` +
		`usr_local_sbin_pkg_exists: $usr_local_sbin_pkg_exists, ` +
		`usr_bin_xbps_install_exists: $usr_bin_xbps_install_exists, ` +
		`usr_bin_pkg_exists: $usr_bin_pkg_exists, ` +
		`usr_sbin_pkgadd_exists: $usr_sbin_pkgadd_exists, ` +
		`usr_bin_emerge_exists: $usr_bin_emerge_exists, ` +
		`usr_sbin_swlist_exists: $usr_sbin_swlist_exists, ` +
		`usr_sbin_pkg_exists: $usr_sbin_pkg_exists, ` +
		`sbin_apk_exists: $sbin_apk_exists, ` +
		`opt_homebrew_bin_brew_exists: $opt_homebrew_bin_brew_exists, ` +
		`usr_local_bin_brew_exists: $usr_local_bin_brew_exists, ` +
		`opt_local_bin_port_exists: $opt_local_bin_port_exists, ` +
		`opt_tools_bin_pkgin_exists: $opt_tools_bin_pkgin_exists, ` +
		`opt_local_bin_pkgin_exists: $opt_local_bin_pkgin_exists, ` +
		`usr_pkg_bin_pkgin_exists: $usr_pkg_bin_pkgin_exists, ` +
		`bin_opkg_exists: $bin_opkg_exists, ` +
		`usr_bin_pacman_exists: $usr_bin_pacman_exists, ` +
		`usr_sbin_urpmi_exists: $usr_sbin_urpmi_exists, ` +
		`usr_bin_zypper_exists: $usr_bin_zypper_exists, ` +
		`usr_bin_apt_get_exists: $usr_bin_apt_get_exists, ` +
		`usr_bin_dnf5_exists: $usr_bin_dnf5_exists, ` +
		`usr_bin_dnf_3_exists: $usr_bin_dnf_3_exists, ` +
		`usr_bin_dnf_exists: $usr_bin_dnf_exists, ` +
		`usr_bin_yum_exists: $usr_bin_yum_exists, ` +
		`apt_provided_by_rpm_package: $apt_provided_by_rpm_package}'); ` +
		`echo "$output"`
)

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

func (p *PackageManagerInfo) populatePackageManagerInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

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

	cmd := transport.NewCommand(packageManagerDiscoveryScript)

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to check package manager status",
			Detail:   fmt.Sprintf("Error checking package manager status: %v", err),
		}}
	}

	stdout := strings.TrimSpace(string(stdoutBytes))

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
