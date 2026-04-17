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

//go:embed dnf_absent.py
var dnfAbsentScript string

//go:embed dnf_latest.py
var dnfLatestScript string

//go:embed dnf_present.py
var dnfPresentScript string

var prunedDnfAbsentScript string
var prunedDnfLatestScript string
var prunedDnfPresentScript string

type dnfPackageInfo struct {
	Name         string `json:"name" cty:"name"`
	Epoch        string `json:"epoch" cty:"epoch"`
	Version      string `json:"version" cty:"version"`
	Release      string `json:"release" cty:"release"`
	Architecture string `json:"architecture" cty:"architecture"`
	Repo         string `json:"repo" cty:"repo"`
}

type dnfOutput struct {
	Installed []dnfPackageInfo `json:"installed,omitempty" cty:"installed"`
	Removed   []dnfPackageInfo `json:"removed,omitempty" cty:"removed"`
}

// RunModule implements [pluginv1.PluginModule].
func (d *DnfModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *result.ModuleResult {
	if hostInfo.PackageManager.Name != "dnf" {
		return pluginv1.NewFailure(
			fmt.Errorf(
				"dnf can only be run on hosts with dnf as package manager, not %s",
				hostInfo.PackageManager.Name,
			),
			"",
		)
	}

	sb := &strings.Builder{}
	sb.WriteString(python.FormatInputForPython(input, whatIf))

	state := input["state"].AsString()
	switch state {
	case "present":
		sb.WriteString(prunedDnfPresentScript)
	case "latest":
		sb.WriteString(prunedDnfLatestScript)
	case "absent":
		sb.WriteString(prunedDnfAbsentScript)
	default:
		return pluginv1.NewFailure(fmt.Errorf("unknown state for dnf module: %s", state), "")
	}

	encodedCommand, err := python.EncodePythonAsBase64(sb.String())
	if err != nil {
		return pluginv1.NewFailure(err, "failed to encode dnf script as base64")
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

			return pluginv1.NewFailure(fmt.Errorf("failed to run dnf script: %s", stderr), "")
		}

		return pluginv1.NewFailure(errors.New(output["error"]), output["error_details"])
	}

	var output dnfOutput
	err = json.Unmarshal(outBuf.Bytes(), &output)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to parse dnf script output")
	}

	changed := len(output.Installed) > 0 || len(output.Removed) > 0

	impliedType, err := gocty.ImpliedType(output)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to get implied type for dnf result")
	}

	outputValue, err := gocty.ToCtyValue(output, impliedType)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to convert dnf result to cty.Value")
	}

	if changed {
		moduleResult, err := pluginv1.NewChanged(outputValue)
		if err != nil {
			return pluginv1.NewFailure(err, "failed to create module result for dnf")
		}

		return moduleResult
	}

	moduleResult, err := pluginv1.NewNotChanged(outputValue)
	if err != nil {
		return pluginv1.NewFailure(err, "failed to create module result for dnf")
	}

	return moduleResult
}
