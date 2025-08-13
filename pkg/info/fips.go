// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"fmt"

	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run ../../cmd/scriptimport/main.go info fips_linux_discovery.sh
//go:generate go run ../../cmd/scriptimport/main.go info fips_windows_discovery.ps1

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

func (f *FIPSInfo) populateFipsInfo(osInfo *OSInfo, transport transport.Transport) util.Diags {

	if osInfo == nil || osInfo.id == "" {
		return util.Diags{&util.Diag{
			Severity: util.DiagWarning,
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
	return util.Diags{}
}

func (f *FIPSInfo) populateLinuxFipsInfo(t transport.Transport) util.Diags {

	cmd, err := t.NewCommand(fipsLinuxDiscoveryScript, nil)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to create FIPS discovery command",
			Detail:   fmt.Sprintf("Error creating command: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to check FIPS status",
			Detail:   fmt.Sprintf("Error checking FIPS status: %v", err),
		}, &util.Diag{
			Severity: util.DiagDebug,
			Summary:  "Discovery command stderr",
			Detail:   fmt.Sprintf("stderr: %s", stderr),
		}}
	}

	if stdout == "" {
		f.enabled = false
		return util.Diags{}
	}

	f.enabled = stdout == "1"

	return util.Diags{}
}

func (f *FIPSInfo) populateWindowsFipsInfo(t transport.Transport) util.Diags {

	cmd, err := t.NewPowerShellCommand(fipsWindowsDiscoveryScript, nil)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to create PowerShell command",
			Detail:   fmt.Sprintf("Error creating PowerShell command: %v", err),
		}}
	}

	stdout, err := cmd.Output(context.Background())
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to check FIPS status",
			Detail:   fmt.Sprintf("Error checking FIPS status: %v", err),
		}}
	}

	f.known = true
	f.enabled = stdout == "1"

	return util.Diags{}
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
