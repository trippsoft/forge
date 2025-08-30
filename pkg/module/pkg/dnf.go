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

//go:embed dnf_setup.py
var dnfSetupScript string

//go:embed dnf_present.py
var dnfPresentScript string

//go:embed dnf_absent.py
var dnfAbsentScript string

//go:embed dnf_latest.py
var dnfLatestScript string

var (
	fullDnfPresentScript string
	fullDnfAbsentScript  string
	fullDnfLatestScript  string
)

type DNFModule struct{}

// InputSpec implements module.Module.
func (m *DNFModule) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (m *DNFModule) Validate(config *module.RunConfig) error {
	return nil
}

// Run implements module.Module.
func (m *DNFModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	t := config.Transport

	state := config.Input["state"].AsString()
	var script string
	switch state {
	case "present":
		script = fullDnfPresentScript
	case "absent":
		script = fullDnfAbsentScript
	case "latest":
		script = fullDnfLatestScript
	}

	if script == "" {
		return module.NewFailure(fmt.Errorf("state %q is not valid", state), "")
	}

	script = fmt.Sprintf("%s\n%s", config.FormatInputForPython(), script)

	cmd, err := t.NewPythonCommand("", script, config.Escalation)
	if err != nil {
		return module.NewFailure(err, "")
	}

	stdout, stderr, err := cmd.OutputWithError(ctx)
	if err != nil {
		return module.NewFailure(err, fmt.Sprintf("failed to execute DNF command: %s", stderr))
	}

	result := &dnfResult{}
	if err := json.Unmarshal([]byte(stdout), result); err != nil {
		return module.NewFailure(err, fmt.Sprintf("failed to parse DNF command output: %s", stdout))
	}

	return result.toModuleResult()
}

// DNFPackageInfo represents information about a package in the DNF package manager.
type DNFPackageInfo struct {
	Name         string `json:"name"`
	Epoch        string `json:"epoch"`
	Version      string `json:"version"`
	Release      string `json:"release"`
	Architecture string `json:"architecture"`
	Repo         string `json:"repo"`
}

// dnfResult represents the result of running the DNF python script.
type dnfResult struct {
	Err       string           `json:"error,omitempty"`
	ErrDetail string           `json:"error_detail,omitempty"`
	Installed []DNFPackageInfo `json:"installed,omitempty"`
	Removed   []DNFPackageInfo `json:"removed,omitempty"`
}

func (r *dnfResult) toModuleResult() *module.Result {
	if r.Err != "" {
		return module.NewFailure(errors.New(r.Err), r.ErrDetail)
	}

	changed := false
	if len(r.Installed) > 0 || len(r.Removed) > 0 {
		changed = true
	}

	return module.NewSuccess(changed, r.createOutput())
}

func (r *dnfResult) createOutput() map[string]cty.Value {
	installedPackages := make([]cty.Value, 0, len(r.Installed))
	for _, pkg := range r.Installed {
		pkgMap := map[string]cty.Value{
			"name":         cty.StringVal(pkg.Name),
			"epoch":        cty.StringVal(pkg.Epoch),
			"version":      cty.StringVal(pkg.Version),
			"release":      cty.StringVal(pkg.Release),
			"architecture": cty.StringVal(pkg.Architecture),
			"repo":         cty.StringVal(pkg.Repo),
		}

		installedPackages = append(installedPackages, cty.ObjectVal(pkgMap))
	}

	removedPackages := make([]cty.Value, 0, len(r.Removed))
	for _, pkg := range r.Removed {
		pkgMap := map[string]cty.Value{
			"name":         cty.StringVal(pkg.Name),
			"epoch":        cty.StringVal(pkg.Epoch),
			"version":      cty.StringVal(pkg.Version),
			"release":      cty.StringVal(pkg.Release),
			"architecture": cty.StringVal(pkg.Architecture),
			"repo":         cty.StringVal(pkg.Repo),
		}

		removedPackages = append(removedPackages, cty.ObjectVal(pkgMap))
	}

	output := map[string]cty.Value{}

	if len(installedPackages) == 0 {
		output["installed"] = cty.ListValEmpty(cty.EmptyObject)
	} else {
		output["installed"] = cty.ListVal(installedPackages)
	}

	if len(removedPackages) == 0 {
		output["removed"] = cty.ListValEmpty(cty.EmptyObject)
	} else {
		output["removed"] = cty.ListVal(removedPackages)
	}

	return output
}
