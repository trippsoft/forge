// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pkg

import (
	"context"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
)

var (
	infoInputSpec = hclspec.NewSpec(hclspec.Object())

	packageManagerInfoModules = map[string]module.Module{
		"dnf": &DNFInfoModule{},
	}

	_ module.Module = (*PkgInfoModule)(nil)
)

// PkgInfoModule is a module for retrieving package information, selecting the implementation
// based on discovered package manager.
type PkgInfoModule struct{}

// InputSpec implements module.Module.
func (m *PkgInfoModule) InputSpec() *hclspec.Spec {
	return infoInputSpec
}

// Validate implements module.Module.
func (m *PkgInfoModule) Validate(config *module.RunConfig) error {
	return nil
}

// Run implements module.Module.
func (m *PkgInfoModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	packageManagerName := config.HostInfo.PackageManagerInfo().Name()
	mod, exists := packageManagerInfoModules[packageManagerName]
	if !exists {
		return module.NewFailure(nil, "unsupported package manager")
	}

	return mod.Run(ctx, config)
}
