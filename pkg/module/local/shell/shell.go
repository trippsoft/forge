// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package shell

import (
	"context"
	"fmt"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(map[string]*hclspec.ObjectField{
		"command": hclspec.RequiredField(hclspec.String),
	}))

	_ module.Module = &Module{} // Ensure Module implements the module.Module interface.
)

type Module struct{}

// InputSpec implements module.Module.
func (s *Module) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (s *Module) Validate(host *inventory.Host, input map[string]cty.Value) error {
	return nil // No specific validation needed for this module.
}

// Run implements module.Module.
func (s *Module) Run(host *inventory.Host, common *module.CommonConfig, input map[string]cty.Value) *module.Result {

	t := host.Transport()

	command := input["command"].AsString()
	cmd, err := t.NewCommand(command, common.Escalation)
	if err != nil {
		return module.NewFailure(err, "failed to create command")
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
		return module.NewFailure(err, fmt.Sprintf("failed to execute command: %s", stderr))
	}

	output := map[string]cty.Value{
		"stdout": cty.StringVal(stdout),
		"stderr": cty.StringVal(stderr),
	}

	return module.NewSuccess(true, output)
}
