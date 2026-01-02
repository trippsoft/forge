// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(hclspec.Object(hclspec.RequiredField("message", hclspec.Raw)))

	_ Module = (*MessageModule)(nil)
)

// Module defines the message module that displays messages to the console.
type MessageModule struct{}

// Info implements Module.
func (s *MessageModule) Info() *ModuleInfo {
	return NewModuleInfo("", "", "message")
}

// InputSpec implements Module.
func (s *MessageModule) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Run implements ModuleExecutor.
func (s *MessageModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	if config == nil {
		return result.NewFailure(errors.New("config cannot be nil"), "")
	}

	if config.Input == nil {
		return result.NewFailure(errors.New("input cannot be nil"), "")
	}

	messageVal := config.Input["message"]

	var message string
	if messageVal.Type().Equals(cty.String) {
		message = messageVal.AsString()
	} else {
		message = util.FormatCtyValueToIndentedString(config.Input["message"], 0, 4)
	}

	result := result.NewSuccess(false, cty.EmptyObjectVal)
	result.Messages = []string{message}

	return result
}
