// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
)

const (
	defaultSuccessMessage = ""
	defaultFailureMessage = "condition failed"
)

var (
	assertInputSpec = hclspec.NewSpec(hclspec.Object(
		hclspec.RequiredField("condition", hclspec.Bool),
		hclspec.OptionalField("success_message", hclspec.String).WithDefaultValue(cty.StringVal(defaultSuccessMessage)),
		hclspec.OptionalField("failure_message", hclspec.String).WithDefaultValue(cty.StringVal(defaultFailureMessage)),
	))

	_ Module = (*AssertModule)(nil)
)

// AssertModule defines the assert module that checks conditions.
type AssertModule struct{}

// Info implements Module.
func (m *AssertModule) Info() *ModuleInfo {
	return NewModuleInfo("", "", "assert")
}

// InputSpec implements Module.
func (m *AssertModule) InputSpec() *hclspec.Spec {
	return assertInputSpec
}

// Run implements Module.
func (m *AssertModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	if config == nil {
		return result.NewFailure(errors.New("config is nil"), "")
	}

	if config.Input == nil {
		return result.NewFailure(errors.New("input is nil"), "")
	}

	condition := config.Input["condition"].True()
	if !condition {
		failureMessage := config.Input["failure_message"]
		if failureMessage.IsWhollyKnown() && !failureMessage.IsNull() {
			return result.NewFailure(errors.New(failureMessage.AsString()), "")
		}

		return result.NewFailure(errors.New(defaultFailureMessage), "")
	}

	successMessage := config.Input["success_message"]
	var message string
	if successMessage.IsWhollyKnown() && !successMessage.IsNull() {
		message = successMessage.AsString()
	} else {
		message = defaultSuccessMessage
	}

	output := map[string]cty.Value{
		"message": cty.StringVal(message),
	}

	result := result.NewSuccess(false, output)
	result.Messages = []string{message}

	return result
}
