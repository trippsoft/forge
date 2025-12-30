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

//go:embed dnf_info.py
var dnfInfoScript string

var prunedDnfInfoScript string

type dnfInfoResult struct {
	Packages map[string]dnfPackageInfo `json:"packages,omitempty"`
}

func (d *dnfInfoResult) toModuleResult() *pluginv1.ModuleResult {
	packages := make(map[string]cty.Value, len(d.Packages))
	for name, pkg := range d.Packages {
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

		packages[name] = value
	}

	result, err := pluginv1.NewModuleSuccess(false, map[string]cty.Value{
		"packages": cty.ObjectVal(packages),
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
func (d *DnfInfoModule) RunModule(
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

	encodedCommand, err := util.EncodePythonAsBase64(prunedDnfInfoScript)
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
				fmt.Errorf("failed to run dnf_info script: %s", stderr),
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
			"failed to parse dnf_info script output",
		)
	}

	return result.toModuleResult()
}
