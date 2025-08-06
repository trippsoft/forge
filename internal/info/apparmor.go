package info

import (
	"context"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/pkg/diag"
	"github.com/zclconf/go-cty/cty"
)

const (
	appArmorDiscoveryScript = `if [ -d /sys/kernel/security/apparmor ]; then apparmor_enabled=1; ` +
		`else apparmor_enabled=0; ` +
		`fi; ` +
		`echo "$apparmor_enabled"`
)

type AppArmorInfo struct {
	supported bool
	enabled   bool
}

func newAppArmorInfo() *AppArmorInfo {
	return &AppArmorInfo{
		supported: false,
		enabled:   false,
	}
}

func (a *AppArmorInfo) Supported() bool {
	return a.supported
}

func (a *AppArmorInfo) Enabled() bool {
	return a.enabled
}

func (a *AppArmorInfo) populateAppArmorInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.id == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping AppArmor information collection due to missing or invalid OS info",
		}}
	}

	if !osInfo.families.Contains("linux") {
		a.supported = false
		a.enabled = false
		return diag.Diags{}
	}

	a.supported = true

	cmd := transport.NewCommand(appArmorDiscoveryScript)

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to execute AppArmor discovery script",
			Detail:   fmt.Sprintf("Error executing command: %v", err),
		}}
	}

	stdout := strings.TrimSpace(string(stdoutBytes))

	if stdout == "" {
		a.enabled = false
		return diag.Diags{}
	}

	a.enabled = stdout == "1"
	return diag.Diags{}
}

func (a *AppArmorInfo) toMapOfCtyValues() map[string]cty.Value {

	if !a.supported {
		return map[string]cty.Value{
			"apparmor_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"apparmor_enabled": cty.BoolVal(a.enabled),
	}
}

// String returns a string representation of the AppArmor information.
// This is useful for logging or debugging purposes.
func (a *AppArmorInfo) String() string {

	if !a.supported {
		return "apparmor_enabled: not supported"
	}

	return fmt.Sprintf("apparmor_enabled: %t", a.enabled)
}
