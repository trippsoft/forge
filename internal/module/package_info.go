// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/zclconf/go-cty/cty"
)

var (
	packageInfoInputSpec = hclspec.NewSpec(hclspec.Object())
	packageInfoModules   = map[string]pluginv1.PluginModule{
		"dnf": DnfInfo,
	}

	PackageInfo pluginv1.PluginModule = &PackageInfoModule{}
)

// PackageInfoModule is a module for retrieving package information, selecting the implementation based on the
// discovered package manager.
type PackageInfoModule struct{}

// Name implements pluginv1.PluginModule.
func (p *PackageInfoModule) Name() string {
	return "package_info"
}

// Type implements pluginv1.PluginModule.
func (p *PackageInfoModule) Type() pluginv1.ModuleType {
	return pluginv1.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (p *PackageInfoModule) InputSpec() *hclspec.Spec {
	return packageInfoInputSpec
}

// RunModule implements pluginv1.PluginModule.
func (p *PackageInfoModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *pluginv1.ModuleResult {

	module, ok := packageInfoModules[hostInfo.PackageManager.Name]
	if !ok {
		return pluginv1.NewModuleFailure(
			fmt.Errorf("unknown package manager: %s", hostInfo.PackageManager.Name),
			"",
		)
	}

	return module.RunModule(hostInfo, input, whatIf)
}
