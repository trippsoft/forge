// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package shell

import (
	"context"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(map[string]*hclspec.ObjectField{
		"command": {Type: hclspec.String, Required: true, DefaultValue: cty.NullVal(cty.String)},
	}))

	_ plugin.LocalPlugin = &Plugin{} // Ensure Plugin implements the plugin.LocalPlugin interface.
)

type Plugin struct{}

// InputSpec implements plugin.Plugin.
func (s *Plugin) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements plugin.Plugin.
func (s *Plugin) Validate(host *inventory.Host, input map[string]cty.Value) error {
	return nil // No specific validation needed for this plugin.
}

// Run implements plugin.Plugin.
func (s *Plugin) Run(host *inventory.Host, common *plugin.CommonConfig, input map[string]cty.Value) *plugin.Result {

	t := host.Transport()

	command := input["command"].AsString()
	cmd, err := t.NewCommand(command, common.Escalation)
	if err != nil {
		return plugin.NewFailure(err)
	}

	timeout := common.Timeout

	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
	}

	stdout, stderr, err := cmd.OutputWithError(ctx)
	if err != nil {
		return plugin.NewFailure(err)
	}

	output := map[string]cty.Value{
		"stdout": cty.StringVal(stdout),
		"stderr": cty.StringVal(stderr),
	}

	return plugin.NewSuccess(true, output)
}
