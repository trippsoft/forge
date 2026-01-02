// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/zclconf/go-cty/cty"
)

var (
	commandInputSpec = hclspec.NewSpec(hclspec.Object(
		hclspec.RequiredField("name", hclspec.String).WithAliases("path"),
		hclspec.OptionalField("args", hclspec.List(hclspec.String)).WithDefaultValue(cty.ListValEmpty(cty.String)),
	))

	Command pluginv1.PluginModule = &CommandModule{}
)

// CommandModule is a module that provides command execution information.
type CommandModule struct{}

// Name implements pluginv1.PluginModule.
func (c *CommandModule) Name() string {
	return "command"
}

// Type implements pluginv1.PluginModule.
func (c *CommandModule) Type() pluginv1.ModuleType {
	return pluginv1.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (c *CommandModule) InputSpec() *hclspec.Spec {
	return commandInputSpec
}

// RunModule implements pluginv1.PluginModule.
func (c *CommandModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *pluginv1.ModuleResult {

	if whatIf {
		r, err := pluginv1.NewModuleSuccess(
			true,
			cty.ObjectVal(map[string]cty.Value{
				"stdout": cty.NullVal(cty.String),
				"stderr": cty.NullVal(cty.String),
			}),
		)

		if err != nil {
			return pluginv1.NewModuleFailure(
				err,
				"failed to create module result",
			)
		}

		return r
	}

	name := input["name"].AsString()
	args := make([]string, 0, input["args"].LengthInt())
	it := input["args"].ElementIterator()
	for it.Next() {
		_, v := it.Element()
		args = append(args, v.AsString())
	}

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return pluginv1.NewModuleFailure(
			err,
			fmt.Sprintf(
				"failed to execute command: %s; stderr: %s",
				name,
				errBuf.String(),
			),
		)
	}

	stdout := strings.TrimSpace(outBuf.String())
	stderr := strings.TrimSpace(errBuf.String())

	r, err := pluginv1.NewModuleSuccess(
		true,
		cty.ObjectVal(map[string]cty.Value{
			"stdout": cty.StringVal(stdout),
			"stderr": cty.StringVal(stderr),
		}),
	)

	if err != nil {
		return pluginv1.NewModuleFailure(
			err,
			"failed to create module result",
		)
	}

	return r
}
