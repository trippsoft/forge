package info

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/pkg/diag"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

const (
	darwinOSDiscoveryScript = `os_arch="$(uname -m || echo \"\")"; ` +
		`os_version="$(sw_vers -productVersion || echo \"\")"; ` +
		`echo "{\"os_arch\": \"$os_arch\", \"os_version\": \"$os_version\"}"`

	linuxOSDiscoveryScript = `os_arch="$(uname -m || echo \"\")"; ` +
		`if [ -e /etc/os-release ]; then source /etc/os-release; ` +
		`elif [ -L /etc/os-release ]; then source "$(readlink -f /etc/os-release || echo \"\")"; ` +
		`elif [ -e /usr/lib/os-release ]; then source /usr/lib/os-release; ` +
		`elif [ -L /usr/lib/os-release ]; then source "$(readlink -f /usr/lib/os-release || echo \"\")"; ` +
		`fi; ` +
		`if [ -n "$ID" ]; then os_id="$ID"; ` +
		`else os_id="$(lsb_release -si || echo \"\")"; ` +
		`fi; ` +
		`if [ -n "$PRETTY_NAME" ]; then os_friendly_name="$PRETTY_NAME"; ` +
		`else os_friendly_name="$(lsb_release -sd || echo \"\")"; ` +
		`fi; ` +
		`if [ -n "$VERSION_ID" ]; then  os_version="$VERSION_ID"; ` +
		`else os_version="$(lsb_release -sr || echo \"\")"; ` +
		`fi; ` +
		`if [ -n "$VERSION_CODENAME" ]; then os_release="$VERSION_CODENAME"; ` +
		`else os_release="$(lsb_release -sc || echo \"\")"; ` +
		`fi; ` +
		`if [ -n "$VARIANT" ]; then os_edition="$VARIANT"; ` +
		`fi; ` +
		`if [ -n "$VARIANT_ID" ]; then os_edition_id="$VARIANT_ID"; ` +
		`fi; ` +
		`output=$(jq -n ` +
		`--arg os_arch "$os_arch" ` +
		`--arg os_id "$os_id" ` +
		`--arg os_friendly_name "$os_friendly_name" ` +
		`--arg os_release "$os_release" ` +
		`--arg os_version "$os_version" ` +
		`--arg os_edition "$os_edition" ` +
		`--arg os_edition_id "$os_edition_id" ` +
		`'{os_arch: $os_arch, os_id: $os_id, os_friendly_name: $os_friendly_name, os_release: $os_release, os_version: $os_version, os_edition: $os_edition, os_edition_id: $os_edition_id}'); ` +
		`echo "$output"`

	windowsOSDiscoveryScript = `Import-Module -Name Dism; ` +
		`$friendlyName = (Get-CimInstance -ClassName Win32_OperatingSystem).Caption; ` +
		`$version = [System.Environment]::OSVersion.Version.ToString(); ` +
		`$osArch = (Get-CimInstance -ClassName Win32_OperatingSystem).OSArchitecture; ` +
		`$procArch = $env:PROCESSOR_ARCHITECTURE; ` +
		`$output = @{` +
		`os_friendly_name = $friendlyName; ` +
		`os_version = $version; ` +
		`os_bits = $osArch; ` +
		`processor_arch = $procArch; ` +
		`}; ` +
		`$json = $output | ConvertTo-Json -Depth 3; ` +
		`Write-Host $json`
)

var (
	architectureMap = map[string]string{
		"386":     "386",
		"i386":    "386",
		"i486":    "386",
		"i586":    "386",
		"i686":    "386",
		"x86":     "386",
		"amd64":   "amd64",
		"x86_64":  "amd64",
		"arm":     "arm",
		"armv6l":  "arm",
		"armv7l":  "arm",
		"aarch64": "arm64",
		"arm64":   "arm64",
		"mips":    "mips",
		"mips64":  "mips64",
		"ppc64":   "ppc64",
		"ppc64le": "ppc64le",
		"riscv64": "riscv64",
		"s390x":   "s390x",
	}

	architectureBitsMap = map[string]int{
		"386":     32,
		"amd64":   64,
		"arm":     32,
		"arm64":   64,
		"mips":    32,
		"mips64":  64,
		"ppc64":   64,
		"ppc64le": 64,
		"riscv64": 64,
		"s390x":   64,
	}

	osFamiliesMap = map[string][]string{
		"almalinux":     {"el"},
		"amazon":        {"el"},
		"archlinux-arm": {"archlinux"},
		"arcolinux":     {"archlinux"},
		"centos":        {"el"},
		"clearos":       {"el"},
		"cloudlinux":    {"el"},
		"deepin":        {"debian"},
		"devuan":        {"debian"},
		"elementary":    {"debian", "ubuntu"},
		"endeavouros":   {"archlinux"},
		"fedora":        {"el"},
		"kali":          {"debian"},
		"kylin":         {"debian", "ubuntu"},
		"linuxmint":     {"debian", "ubuntu"}, // We will treat Linux Mint as always being Ubuntu-based, despite the existence of Debian-based versions.
		"mageia":        {"mandrake"},
		"manjaro":       {"archlinux"},
		"manjaro-arm":   {"archlinux", "manjaro"},
		"nobara":        {"el", "fedora"},
		"opensuse":      {"suse"},
		"oraclelinux":   {"el"},
		"pop_os":        {"debian", "ubuntu"},
		"raspbian":      {"debian"},
		"rhel":          {"el"},
		"rocky":         {"el"},
		"scientific":    {"el"},
		"sled":          {"suse"},
		"sles":          {"suse"},
		"ubuntu":        {"debian"},
		"virtuozzo":     {"el"},
	}

	osIDCorrectionMap = map[string]string{
		"amzn":           "amazon",
		"arch":           "archlinux",
		"archarm":        "archlinux-arm",
		"clear-linux-os": "clearlinux",
		"cumulus-linux":  "cumuluslinux",
		"pop":            "pop_os",
		"ol":             "oraclelinux",
		"opensuse-leap":  "opensuse",
		"sles_sap":       "sles",
	}
)

type OSInfo struct {
	families *util.Set[string]

	kernel       string
	id           string
	friendlyName string
	release      string
	majorVersion string
	version      string
	edition      string
	editionID    string

	osArch     string
	osArchBits int

	procArch     string
	procArchBits int
}

func newOSInfo() *OSInfo {
	return &OSInfo{
		families: util.NewSet[string](),
	}
}

func (o *OSInfo) Families() util.ReadOnlySet[string] {
	return o.families
}

func (o *OSInfo) Kernel() string {
	return o.kernel
}

func (o *OSInfo) ID() string {
	return o.id
}

func (o *OSInfo) FriendlyName() string {
	return o.friendlyName
}

func (o *OSInfo) Release() string {
	return o.release
}

func (o *OSInfo) MajorVersion() string {
	return o.majorVersion
}

func (o *OSInfo) Version() string {
	return o.version
}

func (o *OSInfo) Edition() string {
	return o.edition
}

func (o *OSInfo) EditionID() string {
	return o.editionID
}

func (o *OSInfo) OSArch() string {
	return o.osArch
}

func (o *OSInfo) OSArchBits() int {
	return o.osArchBits
}

func (o *OSInfo) ProcArch() string {
	return o.procArch
}

func (o *OSInfo) ProcArchBits() int {
	return o.procArchBits
}

func (o *OSInfo) populateOSInfo(transport transport.Transport) diag.Diags {

	cmd := transport.NewCommand("uname -s")

	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	unameErr := cmd.Run(context.Background())
	o.kernel = strings.ToLower(strings.TrimSpace(outBuf.String()))
	if unameErr == nil {
		o.families.Add("posix")
		o.families.Add(o.kernel)

		switch o.kernel {
		case "darwin":
			return o.populateDarwinOSInfo(transport)
		case "linux":
			return o.populateLinuxOSInfo(transport)
		}
	}

	psCmd, psErr := transport.NewPowerShellCommand("Write-Host $PSVersionTable.PSVersion")
	if psErr != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Unsupported OS family",
			Detail:   fmt.Sprintf("uname error: %v, PowerShell error: %v", unameErr, psErr),
		}}
	}

	psErr = psCmd.Run(context.Background())
	if psErr == nil {
		o.families.Add("windows")
		o.kernel = "windows"
		return o.populateWindowsOSInfo(transport)
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagError,
		Summary:  "Unsupported OS family",
		Detail:   fmt.Sprintf("uname error: %v, PowerShell error: %v", unameErr, psErr),
	}}
}

func (o *OSInfo) populateDarwinOSInfo(transport transport.Transport) diag.Diags {

	o.id = "macos"
	o.families.Add(o.id)

	cmd := transport.NewCommand(darwinOSDiscoveryScript)
	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	err := cmd.Run(context.Background())
	if err != nil {
		o.friendlyName = "macOS"
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get macOS version",
			Detail:   fmt.Sprintf("Error executing discovery command: %v", err),
		}}
	}

	stdout := strings.TrimSpace(outBuf.String())

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		o.friendlyName = "macOS"
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse macOS discovery output",
			Detail:   fmt.Sprintf("Error parsing JSON output: %v", err),
		}}
	}

	diags := o.populatePosixArchitectureInfo(discoveredData)

	moreDiags := o.populateVersionInfo(discoveredData)
	diags = diags.AppendAll(moreDiags)

	o.friendlyName = fmt.Sprintf("macOS %s", o.version)

	switch o.majorVersion {
	case "26":
		o.release = "Tahoe"
	case "15":
		o.release = "Sequoia"
	case "14":
		o.release = "Sonoma"
	case "13":
		o.release = "Ventura"
	case "12":
		o.release = "Monterey"
	case "11":
		o.release = "Big Sur"
	default:
		diags = diags.Append(&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown macOS release",
			Detail:   fmt.Sprintf("Unknown macOS release detected for major version %s", o.majorVersion),
		})
	}

	return diags
}

func (o *OSInfo) populateLinuxOSInfo(transport transport.Transport) diag.Diags {

	cmd := transport.NewCommand(linuxOSDiscoveryScript)
	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	err := cmd.Run(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get Linux OS information",
			Detail:   fmt.Sprintf("Error executing Linux discovery script: %v", err),
		}}
	}

	stdout := strings.TrimSpace(outBuf.String())

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse Linux discovery output",
			Detail:   fmt.Sprintf("Error parsing JSON output: %v", err),
		}}
	}

	diags := o.populatePosixArchitectureInfo(discoveredData)

	moreDiags := o.populateVersionInfo(discoveredData)
	diags = diags.AppendAll(moreDiags)

	o.id, _ = discoveredData["os_id"]

	o.id = strings.Trim(strings.ToLower(o.id), "\"")

	if o.id == "n/a" {
		o.id = ""
	}

	if correctedID, exists := osIDCorrectionMap[o.id]; exists {
		o.id = correctedID
	}

	o.families.Add(o.id)

	additionalFamilies, exists := osFamiliesMap[o.id]
	if exists {
		for _, family := range additionalFamilies {
			o.families.Add(family)
		}
	}

	o.friendlyName, _ = discoveredData["os_friendly_name"]

	o.friendlyName = strings.Trim(strings.TrimSpace(o.friendlyName), "\"")

	if o.friendlyName == "n/a" || o.friendlyName == "" {
		o.friendlyName = o.id
	}

	o.release, _ = discoveredData["os_release"]

	o.release = strings.Trim(strings.TrimSpace(o.release), "\"")

	if o.release == "n/a" {
		o.release = ""
	}

	o.edition, _ = discoveredData["os_edition"]

	o.edition = strings.Trim(strings.TrimSpace(o.edition), "\"")

	if o.edition == "n/a" {
		o.edition = ""
	}

	o.editionID, _ = discoveredData["os_edition_id"]

	o.editionID = strings.Trim(strings.ToLower(strings.TrimSpace(o.editionID)), "\"")

	if o.editionID == "n/a" {
		o.editionID = ""
	}

	return diags
}

func (o *OSInfo) populateWindowsOSInfo(transport transport.Transport) diag.Diags {

	cmd, err := transport.NewPowerShellCommand(windowsOSDiscoveryScript)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get Windows OS information",
			Detail:   fmt.Sprintf("Error executing Windows discovery script: %v", err),
		}}
	}

	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	err = cmd.Run(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get Windows OS information",
			Detail:   fmt.Sprintf("Error executing Windows discovery script: %v", err),
		}}
	}

	stdout := strings.TrimSpace(outBuf.String())

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse Windows discovery output",
			Detail:   fmt.Sprintf("Error parsing JSON output: %v", err),
		}}
	}

	diags := o.populateVersionInfo(discoveredData)

	moreDiags := o.populateWindowsArchitectureInfo(discoveredData)
	diags = diags.AppendAll(moreDiags)

	o.friendlyName, _ = discoveredData["os_friendly_name"]

	isServer := strings.Contains(strings.ToLower(o.friendlyName), "server")

	if isServer {
		o.id = "windows-server"
	} else {
		o.id = "windows-client"
	}

	o.families.Add(o.id)

	var friendlyNamePrefix string

	switch {
	case strings.HasPrefix(o.version, "6.1.7600") && isServer:
		o.release = "server-2008-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2"
	case strings.HasPrefix(o.version, "6.1.7600") && !isServer:
		o.release = "7"
		friendlyNamePrefix = "Microsoft Windows 7"
	case strings.HasPrefix(o.version, "6.1.7601") && isServer:
		o.release = "server-2008-r2-sp1"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2 SP1"
	case strings.HasPrefix(o.version, "6.1.7601") && !isServer:
		o.release = "7-sp1"
		friendlyNamePrefix = "Microsoft Windows 7 SP1"
	case strings.HasPrefix(o.version, "6.2.9200") && isServer:
		o.release = "server-2012"
		friendlyNamePrefix = "Microsoft Windows Server 2012"
	case strings.HasPrefix(o.version, "6.2.9200") && !isServer:
		o.release = "8"
		friendlyNamePrefix = "Microsoft Windows 8"
	case strings.HasPrefix(o.version, "6.3.9600") && isServer:
		o.release = "server-2012-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2012 R2"
	case strings.HasPrefix(o.version, "6.3.9600") && !isServer:
		o.release = "8.1"
		friendlyNamePrefix = "Microsoft Windows 8.1"
	case strings.HasPrefix(o.version, "10.0.10240"):
		o.release = "10-1507"
		friendlyNamePrefix = "Microsoft Windows 10 1507"
	case strings.HasPrefix(o.version, "10.0.10586"):
		o.release = "10-1511"
		friendlyNamePrefix = "Microsoft Windows 10 1511"
	case strings.HasPrefix(o.version, "10.0.14393") && isServer:
		o.release = "server-2016"
		friendlyNamePrefix = "Microsoft Windows Server 2016"
	case strings.HasPrefix(o.version, "10.0.14393") && !isServer:
		o.release = "10-1607"
		friendlyNamePrefix = "Microsoft Windows 10 1607"
	case strings.HasPrefix(o.version, "10.0.15063"):
		o.release = "10-1703"
		friendlyNamePrefix = "Microsoft Windows 10 1703"
	case strings.HasPrefix(o.version, "10.0.16299"):
		o.release = "10-1709"
		friendlyNamePrefix = "Microsoft Windows 10 1709"
	case strings.HasPrefix(o.version, "10.0.17134"):
		o.release = "10-1803"
		friendlyNamePrefix = "Microsoft Windows 10 1803"
	case strings.HasPrefix(o.version, "10.0.17763") && isServer:
		o.release = "server-2019"
		friendlyNamePrefix = "Microsoft Windows Server 2019"
	case strings.HasPrefix(o.version, "10.0.17763") && !isServer:
		o.release = "10-1809"
		friendlyNamePrefix = "Microsoft Windows 10 1809"
	case strings.HasPrefix(o.version, "10.0.18362"):
		o.release = "10-1903"
		friendlyNamePrefix = "Microsoft Windows 10 1903"
	case strings.HasPrefix(o.version, "10.0.18363") && isServer:
		o.release = "server-1909"
		friendlyNamePrefix = "Microsoft Windows Server 1909"
	case strings.HasPrefix(o.version, "10.0.18363") && !isServer:
		o.release = "10-1909"
		friendlyNamePrefix = "Microsoft Windows 10 1909"
	case strings.HasPrefix(o.version, "10.0.19041") && isServer:
		o.release = "server-2004"
		friendlyNamePrefix = "Microsoft Windows Server 2004"
	case strings.HasPrefix(o.version, "10.0.19041") && !isServer:
		o.release = "10-2004"
		friendlyNamePrefix = "Microsoft Windows 10 2004"
	case strings.HasPrefix(o.version, "10.0.19042") && isServer:
		o.release = "server-20h2"
		friendlyNamePrefix = "Microsoft Windows Server 20H2"
	case strings.HasPrefix(o.version, "10.0.19042") && !isServer:
		o.release = "10-20h2"
		friendlyNamePrefix = "Microsoft Windows 10 20H2"
	case strings.HasPrefix(o.version, "10.0.19043"):
		o.release = "10-21h1"
		friendlyNamePrefix = "Microsoft Windows 10 21H1"
	case strings.HasPrefix(o.version, "10.0.19044"):
		o.release = "10-21h2"
		friendlyNamePrefix = "Microsoft Windows 10 21H2"
	case strings.HasPrefix(o.version, "10.0.19045"):
		o.release = "10-22h2"
		friendlyNamePrefix = "Microsoft Windows 10 22H2"
	case strings.HasPrefix(o.version, "10.0.20348"):
		o.release = "server-2022"
		friendlyNamePrefix = "Microsoft Windows Server 2022"
	case strings.HasPrefix(o.version, "10.0.22000"):
		o.release = "11-21h2"
		friendlyNamePrefix = "Microsoft Windows 11 21H2"
	case strings.HasPrefix(o.version, "10.0.22621"):
		o.release = "11-22h2"
		friendlyNamePrefix = "Microsoft Windows 11 22H2"
	case strings.HasPrefix(o.version, "10.0.22631"):
		o.release = "11-23h2"
		friendlyNamePrefix = "Microsoft Windows 11 23H2"
	case strings.HasPrefix(o.version, "10.0.25398"):
		o.release = "server-23h2"
		friendlyNamePrefix = "Microsoft Windows Server 23H2"
	case strings.HasPrefix(o.version, "10.0.26100") && isServer:
		o.release = "server-2025"
		friendlyNamePrefix = "Microsoft Windows Server 2025"
	case strings.HasPrefix(o.version, "10.0.26100") && !isServer:
		o.release = "11-24h2"
		friendlyNamePrefix = "Microsoft Windows 11 24H2"
	}

	friendlyNameParts := strings.Split(o.friendlyName, " ")
	for len(friendlyNameParts) > 0 && strings.Contains(friendlyNamePrefix, friendlyNameParts[0]) {
		friendlyNameParts = friendlyNameParts[1:]
	}

	friendlyNameBuilder := strings.Builder{}
	friendlyNameBuilder.WriteString(friendlyNamePrefix)

	editionBuilder := strings.Builder{}
	editionIdBuilder := strings.Builder{}

	for _, part := range friendlyNameParts {
		if strings.Contains(part, "Edition") || strings.Contains(part, "Evaluation") {
			continue
		}

		friendlyNameBuilder.WriteString(" " + part)
		if editionBuilder.Len() == 0 {
			editionBuilder.WriteString(part)
			editionIdBuilder.WriteString(strings.ToLower(part))
		} else {
			editionBuilder.WriteString(" " + part)
			editionIdBuilder.WriteString("-" + strings.ToLower(part))
		}
	}

	o.friendlyName = strings.TrimSpace(friendlyNameBuilder.String())
	o.edition = strings.TrimSpace(editionBuilder.String())
	o.editionID = strings.TrimSpace(editionIdBuilder.String())

	return diags
}

func (o *OSInfo) populatePosixArchitectureInfo(data map[string]string) diag.Diags {

	archString, exists := data["os_arch"]
	if !exists {
		o.procArch = ""
		o.procArchBits = 0
		o.osArch = ""
		o.osArchBits = 0
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture",
			Detail:   "No architecture information found in discovery output",
		}}
	}

	archString = strings.Trim(strings.ToLower(strings.TrimSpace(archString)), "\"")

	arch, exists := architectureMap[archString]
	if !exists {
		o.procArch = archString
		o.procArchBits = 0
		o.osArch = archString
		o.osArchBits = 0
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture",
			Detail:   fmt.Sprintf("Unknown architecture %q detected, using it as is", archString),
		}}
	}

	o.procArch = arch
	o.osArch = arch

	archBits, exists := architectureBitsMap[arch]
	if !exists {
		o.procArchBits = 0
		o.osArchBits = 0

		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture bits",
			Detail:   fmt.Sprintf("Unknown architecture bits for %q detected", arch),
		}}
	}

	o.procArchBits = archBits
	o.osArchBits = archBits

	return diag.Diags{}
}

func (o *OSInfo) populateWindowsArchitectureInfo(data map[string]string) diag.Diags {

	procArchString, exists := data["processor_arch"]
	if !exists {
		o.procArch = ""
		o.procArchBits = 0
		o.osArch = ""
		o.osArchBits = 0
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture",
			Detail:   "No architecture information found in discovery output",
		}}
	}

	procArchString = strings.ToLower(strings.TrimSpace(procArchString))

	procArch, exists := architectureMap[procArchString]
	if !exists {
		o.procArch = procArchString
		o.procArchBits = 0
		o.osArch = procArchString
		o.osArchBits = 0
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture",
			Detail:   fmt.Sprintf("Unknown architecture %q detected, using it as is", procArchString),
		}}
	}

	o.procArch = procArch

	procArchBits, exists := architectureBitsMap[procArch]
	if !exists {
		o.procArchBits = 0
		o.osArch = procArch
		o.osArchBits = 0

		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown architecture bits",
			Detail:   fmt.Sprintf("Unknown architecture bits for %q detected", procArch),
		}}
	}

	o.procArchBits = procArchBits

	osArchString, exists := data["os_bits"]
	if !exists {
		o.osArch = procArch
		o.osArchBits = procArchBits
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown OS architecture",
			Detail:   "No OS architecture found in discovery output, using processor architecture",
		}}
	}

	switch {
	case o.procArch == "amd64" && osArchString == "64-bit":
		o.osArch = "amd64"
		o.osArchBits = 64
	case o.procArch == "amd64" && osArchString == "32-bit":
		o.osArch = "386"
		o.osArchBits = 32
	case o.procArch == "386" && osArchString == "32-bit":
		o.osArch = "386"
		o.osArchBits = 32
	case o.procArch == "arm64" && osArchString == "64-bit":
		o.osArch = "arm64"
		o.osArchBits = 64
	case o.procArch == "arm64" && osArchString == "32-bit":
		o.osArch = "arm"
		o.osArchBits = 32
	case o.procArch == "arm" && osArchString == "32-bit":
		o.osArch = "arm"
		o.osArchBits = 32
	default:
		o.osArch = procArch
		o.osArchBits = procArchBits
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Unknown OS architecture",
			Detail:   fmt.Sprintf("Unknown OS architecture %q detected, using processor architecture", procArch),
		}}
	}

	return diag.Diags{}
}

func (o *OSInfo) populateVersionInfo(data map[string]string) diag.Diags {

	o.version, _ = data["os_version"]

	o.version = strings.Trim(strings.TrimSpace(o.version), "\"")

	versionParts := strings.Split(o.version, ".")
	o.majorVersion = versionParts[0]

	return diag.Diags{}
}

func (o *OSInfo) toMapOfCtyValues() map[string]cty.Value {

	values := make(map[string]cty.Value)

	if o.families.Size() > 0 {

		families := make([]cty.Value, 0, o.families.Size())
		for _, family := range o.families.Items() {
			families = append(families, cty.StringVal(family))
		}

		values["os_families"] = cty.SetVal(families)

	} else {
		values["os_families"] = cty.NullVal(cty.Set(cty.String))
	}

	if o.kernel == "" {
		values["os_kernel"] = cty.NullVal(cty.String)
	} else {
		values["os_kernel"] = cty.StringVal(o.kernel)
	}

	if o.id == "" {
		values["os_id"] = cty.NullVal(cty.String)
	} else {
		values["os_id"] = cty.StringVal(o.id)
	}

	if o.friendlyName == "" {
		values["os_friendly_name"] = cty.NullVal(cty.String)
	} else {
		values["os_friendly_name"] = cty.StringVal(o.friendlyName)
	}

	if o.release == "" {
		values["os_release"] = cty.NullVal(cty.String)
	} else {
		values["os_release"] = cty.StringVal(o.release)
	}

	if o.majorVersion == "" {
		values["os_major_version"] = cty.NullVal(cty.String)
	} else {
		values["os_major_version"] = cty.StringVal(o.majorVersion)
	}

	if o.version == "" {
		values["os_version"] = cty.NullVal(cty.String)
	} else {
		values["os_version"] = cty.StringVal(o.version)
	}

	if o.edition == "" {
		values["os_edition"] = cty.NullVal(cty.String)
	} else {
		values["os_edition"] = cty.StringVal(o.edition)
	}

	if o.editionID == "" {
		values["os_edition_id"] = cty.NullVal(cty.String)
	} else {
		values["os_edition_id"] = cty.StringVal(o.editionID)
	}

	if o.osArch == "" {
		values["os_architecture"] = cty.NullVal(cty.String)
	} else {
		values["os_architecture"] = cty.StringVal(o.osArch)
	}

	if o.osArchBits == 0 {
		values["os_architecture_bits"] = cty.NullVal(cty.Number)
	} else {
		values["os_architecture_bits"] = cty.NumberIntVal(int64(o.osArchBits))
	}

	if o.procArch == "" {
		values["processor_architecture"] = cty.NullVal(cty.String)
	} else {
		values["processor_architecture"] = cty.StringVal(o.procArch)
	}

	if o.procArchBits == 0 {
		values["processor_architecture_bits"] = cty.NullVal(cty.Number)
	} else {
		values["processor_architecture_bits"] = cty.NumberIntVal(int64(o.procArchBits))
	}

	return values
}

// String returns a string representation of the OS information.
// This is useful for logging or debugging purposes.
func (o *OSInfo) String() string {

	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString("os_families: ")

	if o.families.Size() == 0 {
		stringBuilder.WriteString("none\n")
	} else {
		for i, family := range o.families.Items() {

			if i > 0 {
				stringBuilder.WriteString(", ")
			}

			stringBuilder.WriteString(family)
		}

		stringBuilder.WriteString("\n")
	}

	stringBuilder.WriteString("os_kernel: ")
	stringBuilder.WriteString(o.kernel)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_id: ")
	stringBuilder.WriteString(o.id)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_friendly_name: ")
	stringBuilder.WriteString(o.friendlyName)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_release: ")
	stringBuilder.WriteString(o.release)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_major_version: ")
	stringBuilder.WriteString(o.majorVersion)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_version: ")
	stringBuilder.WriteString(o.version)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_edition: ")
	stringBuilder.WriteString(o.edition)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_edition_id: ")
	stringBuilder.WriteString(o.editionID)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("os_architecture: ")
	stringBuilder.WriteString(o.osArch)
	stringBuilder.WriteString("\n")

	if o.osArchBits == 0 {
		stringBuilder.WriteString("os_architecture_bits: unknown\n")
	} else {
		stringBuilder.WriteString("os_architecture_bits: ")
		fmt.Fprintf(stringBuilder, "%d", o.osArchBits)
		stringBuilder.WriteString("\n")
	}

	stringBuilder.WriteString("processor_architecture: ")
	stringBuilder.WriteString(o.procArch)
	stringBuilder.WriteString("\n")

	if o.procArchBits == 0 {
		stringBuilder.WriteString("processor_architecture_bits: unknown")
	} else {
		stringBuilder.WriteString("processor_architecture_bits: ")
		fmt.Fprintf(stringBuilder, "%d", o.procArchBits)
	}

	return stringBuilder.String()
}
