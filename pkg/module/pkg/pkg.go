// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/python"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(
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

	packageManagerModules = map[string]module.Module{
		"dnf": &DNFModule{},
	}

	_ module.Module = (*PkgModule)(nil)
)

func init() {
	fullDnfInfoScript = fmt.Sprintf("%s\n%s\n%s\n%s", python.FailFunction, python.SuccessFunction, dnfSetupScript, dnfInfoScript)
	fullDnfInfoScript = util.RemoveEmptyLinesAndComments(fullDnfInfoScript)

	fullDnfAbsentScript = fmt.Sprintf("%s\n%s\n%s\n%s", python.FailFunction, python.SuccessFunction, dnfSetupScript, dnfAbsentScript)
	fullDnfAbsentScript = util.RemoveEmptyLinesAndComments(fullDnfAbsentScript)
	fullDnfPresentScript = fmt.Sprintf("%s\n%s\n%s\n%s", python.FailFunction, python.SuccessFunction, dnfSetupScript, dnfPresentScript)
	fullDnfPresentScript = util.RemoveEmptyLinesAndComments(fullDnfPresentScript)
	fullDnfLatestScript = fmt.Sprintf("%s\n%s\n%s\n%s", python.FailFunction, python.SuccessFunction, dnfSetupScript, dnfLatestScript)
	fullDnfLatestScript = util.RemoveEmptyLinesAndComments(fullDnfLatestScript)
}

// PkgModule is a module for managing packages, selecting the
// implementation based on discovered package manager.
type PkgModule struct{}

// InputSpec implements module.Module.
func (p *PkgModule) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (p *PkgModule) Validate(config *module.RunConfig) error {
	packageManagerName := config.HostInfo.PackageManagerInfo().Name()
	mod, exists := packageManagerModules[packageManagerName]
	if !exists {
		return errors.New("unsupported package manager")
	}

	return mod.Validate(config)
}

// Run implements module.Module.
func (p *PkgModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	packageManagerName := config.HostInfo.PackageManagerInfo().Name()
	mod, exists := packageManagerModules[packageManagerName]
	if !exists {
		return module.NewFailure(nil, "unsupported package manager")
	}

	return mod.Run(ctx, config)
}
