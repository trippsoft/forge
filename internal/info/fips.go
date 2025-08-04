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
	fipsLinuxDiscoveryScript = `if [ -f /proc/sys/crypto/fips_enabled ]; ` +
		`then fips_enabled=$(cat /proc/sys/crypto/fips_enabled); ` +
		`else fips_enabled=0; ` +
		`fi; ` +
		`echo "$fips_enabled"`
	fipsWindowsDiscoveryScript = `$value = Get-ItemPropertyValue -LiteralPath 'HKLM:\SYSTEM\CurrentControlSet\Control\Lsa\FipsAlgorithm' -Name 'Enabled' -ErrorAction SilentlyContinue; Write-Host $value`
)

type FIPSInfo struct {
	known   bool
	enabled bool
}

func newFipsInfo() *FIPSInfo {
	return &FIPSInfo{
		known:   false,
		enabled: false,
	}
}

func (f *FIPSInfo) Known() bool {
	return f.known
}

func (f *FIPSInfo) Enabled() bool {
	return f.enabled
}

func (f *FIPSInfo) populateFipsInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.id == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping FIPS information collection due to missing or invalid OS info",
		}}
	}

	if osInfo.families.Contains("linux") {
		f.known = true
		return f.populateLinuxFipsInfo(transport)
	}

	if osInfo.families.Contains("windows") {
		f.known = true
		return f.populateWindowsFipsInfo(transport)
	}

	f.known = false
	f.enabled = false
	return diag.Diags{}
}

func (f *FIPSInfo) populateLinuxFipsInfo(transport transport.Transport) diag.Diags {

	stdout, _, err := transport.ExecuteCommand(context.Background(), fipsLinuxDiscoveryScript)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to check FIPS status",
			Detail:   fmt.Sprintf("Error checking FIPS status: %v", err),
		}}
	}

	if stdout == "" {
		f.enabled = false
		return diag.Diags{}
	}

	f.enabled = stdout == "1"

	return diag.Diags{}
}

func (f *FIPSInfo) populateWindowsFipsInfo(transport transport.Transport) diag.Diags {

	stdout, err := transport.ExecutePowerShell(context.Background(), fipsWindowsDiscoveryScript)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to check FIPS status",
			Detail:   fmt.Sprintf("Error checking FIPS status: %v", err),
		}}
	}

	stdout = strings.TrimSpace(stdout)

	f.known = true
	f.enabled = stdout == "1"

	return diag.Diags{}
}

func (f *FIPSInfo) toMapOfCtyValues() map[string]cty.Value {

	if !f.known {
		return map[string]cty.Value{
			"fips_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"fips_enabled": cty.BoolVal(f.enabled),
	}
}

// String returns a string representation of the FIPS information.
// This is useful for logging or debugging purposes.
func (f *FIPSInfo) String() string {

	if !f.known {
		return "fips_enabled: unknown on this OS"
	}

	return fmt.Sprintf("fips_enabled: %t", f.enabled)
}
