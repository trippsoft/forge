// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package module

import (
	"fmt"
	"runtime"

	"github.com/trippsoft/forge/pkg/info"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/zclconf/go-cty/cty"
)

// RunModule implements pluginv1.PluginModule.
func (d *DnfInfoModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *pluginv1.ModuleResult {

	return pluginv1.NewModuleFailure(fmt.Errorf("dnf_info cannot be run on %s", runtime.GOOS), "")
}
