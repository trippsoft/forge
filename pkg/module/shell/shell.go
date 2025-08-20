// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package shell

import (
	"context"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(hclspec.RequiredField("command", hclspec.String)))

	_ module.Module = (*Module)(nil)
)

// Module defines the shell module that executes shell commands.
type Module struct{}

// InputSpec implements module.Module.
func (s *Module) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (s *Module) Validate(config *module.RunConfig) error {
	return nil // No specific validation needed for this module.
}

// Run implements module.Module.
func (s *Module) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	if config.WhatIf {
		return module.NewSuccess(true, map[string]cty.Value{
			"stdout": cty.NullVal(cty.String),
			"stderr": cty.NullVal(cty.String),
		})
	}

	t := config.Transport

	command := config.Input["command"].AsString()
	cmd, err := t.NewCommand(command, config.Escalation)
	if err != nil {
		return module.NewFailure(err, "failed to create command")
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
