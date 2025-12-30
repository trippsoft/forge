// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package module

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trippsoft/forge/pkg/info"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

//go:embed dnf_absent.py
var dnfAbsentScript string

//go:embed dnf_latest.py
var dnfLatestScript string

//go:embed dnf_present.py
var dnfPresentScript string

var prunedDnfAbsentScript string
var prunedDnfLatestScript string
var prunedDnfPresentScript string

type dnfResult struct {
	Installed []dnfPackageInfo `json:"installed,omitempty"`
	Removed   []dnfPackageInfo `json:"removed,omitempty"`
}

func (d *dnfResult) toModuleResult() *pluginv1.ModuleResult {
	installed := make([]cty.Value, 0, len(d.Installed))
	for _, pkg := range d.Installed {
		impliedType, err := gocty.ImpliedType(&pkg)
		if err != nil {
			return pluginv1.NewModuleFailure(
				err,
				"failed to get implied type of dnf package info",
			)
		}

		value, err := gocty.ToCtyValue(pkg, impliedType)
		if err != nil {
			return pluginv1.NewModuleFailure(
				err,
				"failed to convert dnf package info to cty.Value",
			)
		}

		installed = append(installed, value)
	}

	removed := make([]cty.Value, 0, len(d.Removed))
	for _, pkg := range d.Removed {
		impliedType, err := gocty.ImpliedType(&pkg)
		if err != nil {
			return pluginv1.NewModuleFailure(
				err,
				"failed to get implied type of dnf package info",
			)
		}

		value, err := gocty.ToCtyValue(pkg, impliedType)
		if err != nil {
			return pluginv1.NewModuleFailure(
				err,
				"failed to convert dnf package info to cty.Value",
			)
		}

		removed = append(removed, value)
	}

	changed := len(installed) > 0 || len(removed) > 0

	result, err := pluginv1.NewModuleSuccess(changed, map[string]cty.Value{
		"installed": cty.ListVal(installed),
		"removed":   cty.ListVal(removed),
	})

	if err != nil {
		return pluginv1.NewModuleFailure(
			err,
			"failed to create module result from dnf info",
		)
	}

	return result
}

// RunModule implements pluginv1.PluginModule.
func (d *DnfModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *pluginv1.ModuleResult {

	if hostInfo.PackageManager.Name != "dnf" {
		return pluginv1.NewModuleFailure(
			fmt.Errorf(
				"dnf_info can only be run on hosts with dnf as package manager, not %s",
				hostInfo.PackageManager.Name,
			),
			"",
		)
	}

	state := input["state"].AsString()
	var script string
	switch state {
	case "present":
		script = prunedDnfPresentScript
	case "latest":
		script = prunedDnfLatestScript
	case "absent":
		script = prunedDnfAbsentScript
	default:
		return pluginv1.NewModuleFailure(
			fmt.Errorf("unknown state for dnf module: %s", state),
			"",
		)
	}

	header := util.FormatInputForPython(input, whatIf)
	script = header + script

	encodedCommand, err := util.EncodePythonAsBase64(script)
	if err != nil {
		return pluginv1.NewModuleFailure(
			err,
			"failed to encode dnf_info script as base64",
		)
	}

	fullPython := fmt.Sprintf("import base64; exec(base64.b64decode('%s'))", encodedCommand)

	args := []string{"-c", fullPython}
	cmd := exec.Command("python3", args...)

	var outBuf, errBuf bytes.Buffer

	cmd.Stderr = &errBuf
	cmd.Stdout = &outBuf

	err = cmd.Run()
	if err != nil {
		var output map[string]string
		e := json.Unmarshal(errBuf.Bytes(), &output)
		if e != nil {
			stderr := strings.TrimSpace(errBuf.String())
			if stderr == "" {
				stderr = err.Error() + " " + encodedCommand
			}

			return pluginv1.NewModuleFailure(
				fmt.Errorf("failed to run dnf script: %s", stderr),
				"",
			)
		}
		return pluginv1.NewModuleFailure(
			errors.New(output["error"]),
			output["error_details"],
		)
	}

	var result dnfInfoResult
	err = json.Unmarshal(outBuf.Bytes(), &result)
	if err != nil {
		return pluginv1.NewModuleFailure(
			err,
			"failed to parse dnf script output",
		)
	}

	return result.toModuleResult()
}
