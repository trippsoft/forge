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

//go:generate go run ../../cmd/scriptimport/main.go info apparmor_discovery.sh

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

func (a *AppArmorInfo) populateAppArmorInfo(osInfo *OSInfo, t transport.Transport) util.Diags {
	if osInfo == nil || osInfo.id == "" {
		return util.Diags{&util.Diag{
			Severity: util.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping AppArmor information collection due to missing or invalid OS info",
		}}
	}

	if !osInfo.families.Contains("linux") {
		a.supported = false
		a.enabled = false
		return util.Diags{}
	}

	a.supported = true

	cmd, err := t.NewCommand(apparmorDiscoveryScript, nil)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to create AppArmor discovery command",
			Detail:   fmt.Sprintf("Error creating command: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return util.Diags{
			&util.Diag{
				Severity: util.DiagError,
				Summary:  "Failed to execute AppArmor discovery script",
				Detail:   fmt.Sprintf("Error executing command: %v", err),
			},
			&util.Diag{
				Severity: util.DiagDebug,
				Summary:  "Discovery command stderr",
				Detail:   fmt.Sprintf("stderr: %s", stderr),
			},
		}
	}

	if stdout == "" {
		a.enabled = false
		return util.Diags{}
	}

	a.enabled = stdout == "1"
	return util.Diags{}
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
