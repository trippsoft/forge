// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package module

import (
	"fmt"
	"runtime"

	"github.com/trippsoft/forge/pkg/info"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
)

// RunModule implements [pluginv1.PluginModule].
func (d *DnfInfoModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *result.ModuleResult {
	return pluginv1.NewFailure(fmt.Errorf("dnf_info cannot be run on %s", runtime.GOOS), "")
}
