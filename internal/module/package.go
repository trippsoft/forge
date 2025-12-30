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
	packageInputSpec = hclspec.NewSpec(
		hclspec.Object(
			hclspec.RequiredField("names", hclspec.List(hclspec.String)),
			hclspec.OptionalField("state", hclspec.String).
				WithDefaultValue(cty.StringVal("present")).
				WithConstraints(hclspec.AllowedValues(
					cty.StringVal("present"),
					cty.StringVal("absent"),
					cty.StringVal("latest"),
				)),
			hclspec.OptionalField("update_cache", hclspec.Bool).
				WithDefaultValue(cty.BoolVal(false)),
			hclspec.OptionalField("autoremove", hclspec.Bool).
				WithDefaultValue(cty.BoolVal(false)),
		),
	)

	packageModules = map[string]pluginv1.PluginModule{
		"dnf": Dnf,
	}

	Package pluginv1.PluginModule = &PackageModule{}
)

// PackageModule is a module for managing packages, selecting the implementation based on the discovered package
// manager.
type PackageModule struct{}

// Name implements pluginv1.PluginModule.
func (p *PackageModule) Name() string {
	return "package"
}

// Type implements pluginv1.PluginModule.
func (p *PackageModule) Type() pluginv1.ModuleType {
	return pluginv1.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (p *PackageModule) InputSpec() *hclspec.Spec {
	return packageInputSpec
}

// RunModule implements pluginv1.PluginModule.
func (p *PackageModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *pluginv1.ModuleResult {

	module, ok := packageModules[hostInfo.PackageManager.Name]
	if !ok {
		return pluginv1.NewModuleFailure(
			fmt.Errorf("unknown package manager: %s", hostInfo.PackageManager.Name),
			"",
		)
	}

	return module.RunModule(hostInfo, input, whatIf)
}
