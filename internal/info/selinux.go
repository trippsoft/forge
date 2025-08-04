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
	selinuxDiscoveryScript = `if [ ! -f /etc/selinux/config ]; then	selinux_installed=0; ` +
		`else selinux_installed=1; ` +
		`selinux_status=$(grep -E '^SELINUX=' /etc/selinux/config | cut -d '=' -f 2); ` +
		`selinux_type=$(grep -E '^SELINUXTYPE=' /etc/selinux/config | cut -d '=' -f 2); ` +
		`fi; ` +
		`output=$(jq -n ` +
		`--arg selinux_installed "$selinux_installed" ` +
		`--arg selinux_status "$selinux_status" ` +
		`--arg selinux_type "$selinux_type" ` +
		`'{selinux_installed: $selinux_installed, selinux_status: $selinux_status, selinux_type: $selinux_type}'); ` +
		`echo "$output"`
)

type SELinuxStatus string
type SELinuxType string

const (
	SELinuxEnforcing    SELinuxStatus = "enforcing"
	SELinuxDisabled     SELinuxStatus = "disabled"
	SELinuxPermissive   SELinuxStatus = "permissive"
	SELinuxNotSupported SELinuxStatus = ""
)

const (
	SELinuxTypeTargeted     SELinuxType = "targeted"
	SELinuxTypeMinimum      SELinuxType = "minimum"
	SELinuxTypeMLS          SELinuxType = "mls"
	SELinuxTypeNotSupported SELinuxType = ""
)

type SELinuxInfo struct {
	supported   bool
	installed   bool
	status      SELinuxStatus
	selinuxType SELinuxType
}

func newSELinuxInfo() *SELinuxInfo {
	return &SELinuxInfo{
		supported:   false,
		installed:   false,
		status:      SELinuxNotSupported,
		selinuxType: SELinuxTypeNotSupported,
	}
}

func (s *SELinuxInfo) Supported() bool {
	return s.supported
}

func (s *SELinuxInfo) Installed() bool {
	return s.installed
}

func (s *SELinuxInfo) Status() SELinuxStatus {
	return s.status
}

func (s *SELinuxInfo) SelinuxType() SELinuxType {
	return s.selinuxType
}

func (s *SELinuxInfo) populateSelinuxInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.id == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping SELinux information collection due to missing or invalid OS info",
		}}
	}

	if !osInfo.families.Contains("linux") {
		s.supported = false
		s.installed = false
		s.status = SELinuxNotSupported
		s.selinuxType = SELinuxTypeNotSupported

		return diag.Diags{}
	}

	s.supported = true

	stdout, _, err := transport.ExecuteCommand(context.Background(), selinuxDiscoveryScript)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to execute SELinux discovery script",
			Detail:   fmt.Sprintf("Error executing SELinux discovery script: %v", err),
		}}
	}

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse SELinux discovery script output",
			Detail:   fmt.Sprintf("Error parsing SELinux discovery script output: %v", err),
		}}
	}

	installed, _ := discoveredData["selinux_installed"]
	s.installed = installed == "1"

	if !s.installed {
		s.status = SELinuxNotSupported
		s.selinuxType = SELinuxTypeNotSupported

		return diag.Diags{}
	}

	status, _ := discoveredData["selinux_status"]
	s.status = SELinuxStatus(status)

	if s.status == SELinuxDisabled {
		s.selinuxType = SELinuxTypeNotSupported
		return diag.Diags{}
	}

	selinuxType, _ := discoveredData["selinux_type"]
	s.selinuxType = SELinuxType(selinuxType)

	return diag.Diags{}
}

func (s *SELinuxInfo) toMapOfCtyValues() map[string]cty.Value {

	if !s.supported {
		return map[string]cty.Value{
			"selinux_installed": cty.NullVal(cty.String),
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	if !s.installed {
		return map[string]cty.Value{
			"selinux_installed": cty.False,
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"selinux_installed": cty.True,
		"selinux_status":    cty.StringVal(string(s.status)),
		"selinux_type":      cty.StringVal(string(s.selinuxType)),
	}
}

// String returns a string representation of the SELinux information.
// This is useful for logging or debugging purposes.
func (s *SELinuxInfo) String() string {

	stringBuilder := &strings.Builder{}
	if !s.supported {
		stringBuilder.WriteString("selinux_installed: not supported on this OS\n")
		stringBuilder.WriteString("selinux_status: not supported on this OS\n")
		stringBuilder.WriteString("selinux_type: not supported on this OS")

		return stringBuilder.String()
	}

	if !s.installed {
		stringBuilder.WriteString("selinux_installed: false\n")
		stringBuilder.WriteString("selinux_status: not installed\n")
		stringBuilder.WriteString("selinux_type: not installed\n")

		return stringBuilder.String()
	}

	stringBuilder.WriteString("selinux_installed: true\n")
	stringBuilder.WriteString("selinux_status: ")
	stringBuilder.WriteString(string(s.status))
	stringBuilder.WriteString("\n")
	stringBuilder.WriteString("selinux_type: ")
	stringBuilder.WriteString(string(s.selinuxType))

	return stringBuilder.String()
}
