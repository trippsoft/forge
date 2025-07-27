package info

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/internal/util"
	"github.com/zclconf/go-cty/cty"
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

	osUnknownError       = errors.New("unknown OS family")
	osNoReleaseFileError = errors.New("neither /etc/os-release nor /usr/lib/os-release found")
)

type osInfo struct {
	families *util.Set[string]

	id           string
	friendlyName string
	release      string
	majorVersion string
	version      string
	edition      string
	editionId    string

	osArch     string
	osArchBits int

	procArch     string
	procArchBits int
}

func newOSInfo() *osInfo {
	return &osInfo{
		families:     util.NewSet[string](),
		id:           "",
		friendlyName: "",
		release:      "",
		majorVersion: "",
		version:      "",
		edition:      "",
		editionId:    "",
		osArch:       "",
		osArchBits:   0,
		procArch:     "",
		procArchBits: 0,
	}
}

func (o *osInfo) Families() *util.Set[string] {
	return o.families
}

func (o *osInfo) Id() string {
	return o.id
}

func (o *osInfo) FriendlyName() string {
	return o.friendlyName
}

func (o *osInfo) Release() string {
	return o.release
}

func (o *osInfo) MajorVersion() string {
	return o.majorVersion
}

func (o *osInfo) Version() string {
	return o.version
}

func (o *osInfo) Edition() string {
	return o.edition
}

func (o *osInfo) EditionId() string {
	return o.editionId
}

func (o *osInfo) OsArch() string {
	return o.osArch
}

func (o *osInfo) OsArchBits() int {
	return o.osArchBits
}

func (o *osInfo) ProcArch() string {
	return o.procArch
}

func (o *osInfo) ProcArchBits() int {
	return o.procArchBits
}

func (o *osInfo) populateOSInfo(transport transport.Transport, fileSystem transport.FileSystem) error {

	_, _, err := transport.ExecuteCommand(context.Background(), "uname -s")
	if err == nil {
		o.families.Add("posix")
		return o.populatePosixOSInfo(transport, fileSystem)
	}

	_, err = transport.ExecutePowerShell(context.Background(), "Write-Host $PSVersionTable.PSVersion")
	if err == nil {
		o.families.Add("windows")
		return o.populateWindowsOSInfo(transport)
	}

	return osUnknownError // Return an error if the shell type is not recognized as POSIX or Windows, subject to future expansion
}

func (o *osInfo) toMapOfCtyValues() map[string]cty.Value {

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

	if o.editionId == "" {
		values["os_edition_id"] = cty.NullVal(cty.String)
	} else {
		values["os_edition_id"] = cty.StringVal(o.editionId)
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
func (o *osInfo) String() string {

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
	stringBuilder.WriteString(o.editionId)
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
