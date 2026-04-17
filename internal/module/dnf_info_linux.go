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
	"github.com/trippsoft/forge/pkg/python"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

//go:embed dnf_info.py
var dnfInfoScript string

var prunedDnfInfoScript string

type dnfInfoOutput struct {
	Packages map[string]dnfPackageInfo `json:"packages,omitempty" cty:"packages"`
}

// RunModule implements [pluginv1.PluginModule].
func (d *DnfInfoModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *result.ModuleResult {
	if hostInfo.PackageManager.Name != "dnf" {
		return pluginv1.NewFailure(
			fmt.Errorf(
				"dnf_info can only be run on hosts with dnf as package manager, not %s",
				hostInfo.PackageManager.Name,
			),
			"",
		)
	}

	encodedCommand, err := python.EncodePythonAsBase64(prunedDnfInfoScript)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to encode dnf_info script as base64")
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

			return pluginv1.NewFailure(fmt.Errorf("failed to run dnf_info script: %s", stderr), "")
		}

		return pluginv1.NewFailure(errors.New(output["error"]), output["error_details"])
	}

	var output dnfInfoOutput
	err = json.Unmarshal(outBuf.Bytes(), &output)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to parse dnf_info script output")
	}

	impliedType, err := gocty.ImpliedType(output)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to get implied type for dnf_info result")
	}

	outputValue, err := gocty.ToCtyValue(output, impliedType)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to convert dnf_info result to cty.Value")
	}

	moduleResult, err := pluginv1.NewNotChanged(outputValue)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to create module result for dnf_info")
	}

	return moduleResult
}
