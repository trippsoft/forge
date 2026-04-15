// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"

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
	assertID = NewModuleID("", "", "assert")

	assert Module = &AssertModule{}
)

// AssertModule defines the assert module that checks conditions.
type AssertModule struct{}

// Info implements Module.
func (m *AssertModule) ID() *ModuleID {
	return assertID
}

// InputSpec implements Module.
func (m *AssertModule) InputSpec() *hclspec.Spec {
	return assertInputSpec
}

// Run implements Module.
func (m *AssertModule) Run(ctx context.Context, config *RunConfig) result.Result {
	if config == nil {
		return result.NewFailedResult("config is nil", "")
	}

	if config.Input == nil {
		return result.NewFailedResult("input is nil", "")
	}

	condition := config.Input["condition"].True()
	if !condition {
		failureMessage := config.Input["failure_message"]
		if failureMessage.IsWhollyKnown() && !failureMessage.IsNull() {
			return result.NewFailedResult(failureMessage.AsString(), "")
		}

		return result.NewFailedResult(defaultFailureMessage, "")
	}

	successMessage := config.Input["success_message"]
	var message string
	if successMessage.IsWhollyKnown() && !successMessage.IsNull() {
		message = successMessage.AsString()
	} else {
		message = defaultSuccessMessage
	}

	outputMap := map[string]cty.Value{
		"message": cty.StringVal(message),
	}

	output := cty.ObjectVal(outputMap)

	r := result.NewNotChangedResult(output)
	r.AddMessages(message)

	return r
}
