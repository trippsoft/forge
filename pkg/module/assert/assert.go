// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package assert

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultSuccessMessage = ""
	defaultFailureMessage = "Condition failed"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(
		hclspec.RequiredField("condition", hclspec.Bool),
		hclspec.OptionalField("success_message", hclspec.String).WithDefaultValue(cty.StringVal(defaultSuccessMessage)),
		hclspec.OptionalField("failure_message", hclspec.String).WithDefaultValue(cty.StringVal(defaultFailureMessage)),
	))

	_ module.Module = (*Module)(nil)
)

// Module defines the assert module that checks conditions.
type Module struct{}

// InputSpec implements module.Module.
func (m *Module) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (m *Module) Validate(config *module.RunConfig) error {
	return nil
}

// Run implements module.Module.
func (m *Module) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	if config == nil {
		return module.NewFailure(errors.New("config is nil"), "")
	}

	if config.Input == nil {
		return module.NewFailure(errors.New("input is nil"), "")
	}

	condition := config.Input["condition"].True()
	if !condition {
		failureMessage, _ := config.Input["failure_message"]
		if failureMessage.IsWhollyKnown() && !failureMessage.IsNull() {
			return module.NewFailure(errors.New(failureMessage.AsString()), "")
		}

		return module.NewFailure(errors.New(defaultFailureMessage), "")
	}

	successMessage, _ := config.Input["success_message"]
	var message string
	if successMessage.IsWhollyKnown() && !successMessage.IsNull() {
		message = successMessage.AsString()
	} else {
		message = defaultSuccessMessage
	}

	output := map[string]cty.Value{
		"message": cty.StringVal(message),
	}

	result := module.NewSuccess(false, output)
	result.Message = message

	return result
}
