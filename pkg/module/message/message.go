// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package message

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(hclspec.RequiredField("message", hclspec.Raw)))

	_ module.Module = (*Module)(nil)
)

// Module defines the message module that displays messages to the console.
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
	if config == nil {
		return module.NewFailure(errors.New("config cannot be nil"), "")
	}

	if config.Input == nil {
		return module.NewFailure(errors.New("input cannot be nil"), "")
	}

	messageVal, _ := config.Input["message"]

	var message string
	if messageVal.Type().Equals(cty.String) {
		message = messageVal.AsString()
	} else {
		message = hclutil.FormatCtyValueToIndentedString(config.Input["message"], 0, 4)
	}

	output := map[string]cty.Value{}
	result := module.NewSuccess(false, output)
	result.Message = message

	return result
}
