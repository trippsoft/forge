// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	_ "embed"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
)

var (
	DnfInfo pluginv1.PluginModule = &DnfInfoModule{}
)

// DnfInfoModule is a module for retrieving DNF package information.
type DnfInfoModule struct{}

// Name implements pluginv1.PluginModule.
func (d *DnfInfoModule) Name() string {
	return "dnf_info"
}

// Type implements pluginv1.PluginModule.
func (d *DnfInfoModule) Type() plugin.ModuleType {
	return plugin.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (d *DnfInfoModule) InputSpec() *hclspec.Spec {
	return packageInfoInputSpec
}
