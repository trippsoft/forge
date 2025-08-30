// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pkg

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

//go:embed dnf_info.py
var dnfInfoScript string

var fullDnfInfoScript string

var (
	_ module.Module = (*DNFInfoModule)(nil)
)

// DNFInfoModule is a module for returning info about the DNF packages installed.
type DNFInfoModule struct{}

// InputSpec implements module.Module.
func (d *DNFInfoModule) InputSpec() *hclspec.Spec {
	return infoInputSpec
}

// Validate implements module.Module.
func (d *DNFInfoModule) Validate(config *module.RunConfig) error {
	return nil
}

// Run implements module.Module.
func (d *DNFInfoModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	t := config.Transport

	script := fmt.Sprintf("%s\n%s", config.FormatInputForPython(), fullDnfInfoScript)
	cmd, err := t.NewPythonCommand("", script, config.Escalation)
	if err != nil {
		return module.NewFailure(err, "failed to create DNF info command")
	}

	stdout, stderr, err := cmd.OutputWithError(ctx)
	if err != nil {
		return module.NewFailure(err, fmt.Sprintf("failed to execute DNF info command: %s", stderr))
	}

	result := &dnfInfoResult{}
	if err := json.Unmarshal([]byte(stdout), result); err != nil {
		return module.NewFailure(err, "failed to unmarshal DNF info command output")
	}

	return result.toModuleResult()
}

// dnfInfoResult represents the result of running the DNF info python script.
type dnfInfoResult struct {
	Err       string                    `json:"error,omitempty"`
	ErrDetail string                    `json:"error_detail,omitempty"`
	Packages  map[string]DNFPackageInfo `json:"packages,omitempty"`
}

func (r *dnfInfoResult) toModuleResult() *module.Result {
	if r.Err != "" {
		return module.NewFailure(errors.New(r.Err), r.ErrDetail)
	}

	return module.NewSuccess(false, r.createOutput())
}

func (r *dnfInfoResult) createOutput() map[string]cty.Value {
	packages := make(map[string]cty.Value, len(r.Packages))
	for name, pkg := range r.Packages {
		pkgMap := map[string]cty.Value{
			"name":         cty.StringVal(pkg.Name),
			"epoch":        cty.StringVal(pkg.Epoch),
			"version":      cty.StringVal(pkg.Version),
			"release":      cty.StringVal(pkg.Release),
			"architecture": cty.StringVal(pkg.Architecture),
			"repo":         cty.StringVal(pkg.Repo),
		}

		packages[name] = cty.ObjectVal(pkgMap)
	}

	if len(packages) == 0 {
		return map[string]cty.Value{
			"packages": cty.EmptyObjectVal,
		}
	}

	return map[string]cty.Value{
		"packages": cty.ObjectVal(packages),
	}
}
