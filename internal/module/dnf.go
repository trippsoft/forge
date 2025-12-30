// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"github.com/trippsoft/forge/pkg/hclspec"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
)

var (
	Dnf pluginv1.PluginModule = &DnfModule{}
)

type dnfPackageInfo struct {
	Name         string `json:"name" cty:"name"`
	Epoch        string `json:"epoch" cty:"epoch"`
	Version      string `json:"version" cty:"version"`
	Release      string `json:"release" cty:"release"`
	Architecture string `json:"architecture" cty:"architecture"`
	Repo         string `json:"repo" cty:"repo"`
}

// DnfModule is a module for managing DNF packages.
type DnfModule struct{}

// Name implements pluginv1.PluginModule.
func (d *DnfModule) Name() string {
	return "dnf"
}

// Type implements pluginv1.PluginModule.
func (d *DnfModule) Type() pluginv1.ModuleType {
	return pluginv1.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (d *DnfModule) InputSpec() *hclspec.Spec {
	return packageInputSpec
}
