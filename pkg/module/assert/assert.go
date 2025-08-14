// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package assert

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(map[string]*hclspec.ObjectField{
		"condition":       hclspec.RequiredField(hclspec.Bool),
		"success_message": hclspec.OptionalField(hclspec.String, cty.NullVal(cty.String)),
		"failure_message": hclspec.OptionalField(hclspec.String, cty.NullVal(cty.String)),
	}))
)

type Module struct{}

// InputSpec implements module.Module.
func (m *Module) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (m *Module) Validate(host *inventory.Host, input map[string]cty.Value) error {
	return nil
}

// Run implements module.Module.
func (m *Module) Run(ctx context.Context, config *module.RunConfig) *module.Result {

	condition := config.Input["condition"].True()

	message := ""
	if !condition {
		failureMessage, exists := config.Input["failure_message"]
		if exists && failureMessage.IsWhollyKnown() && !failureMessage.IsNull() {
			message = failureMessage.AsString()
		} else {
			message = "Condition is false"
		}

		return module.NewFailure(errors.New(message), message)
	}

	successMessage, exists := config.Input["success_message"]
	if exists && successMessage.IsWhollyKnown() && !successMessage.IsNull() {
		message = successMessage.AsString()
	} else {
		message = "Condition is true"
	}

	output := map[string]cty.Value{
		"message": cty.StringVal(message),
	}

	result := module.NewSuccess(false, output)

	result.Message = message

	return result
}
