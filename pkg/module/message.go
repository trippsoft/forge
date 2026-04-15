// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
)

var (
	messageInputSpec = hclspec.NewSpec(hclspec.Object(hclspec.RequiredField("message", hclspec.Raw)))
	messageID        = NewModuleID("", "", "message")

	message Module = &MessageModule{}
)

// Module defines the message module that displays messages to the console.
type MessageModule struct{}

// Info implements Module.
func (s *MessageModule) ID() *ModuleID {
	return messageID
}

// InputSpec implements Module.
func (s *MessageModule) InputSpec() *hclspec.Spec {
	return messageInputSpec
}

// Run implements ModuleExecutor.
func (s *MessageModule) Run(ctx context.Context, config *RunConfig) result.Result {
	if config == nil {
		return result.NewFailedResult("config cannot be nil", "")
	}

	if config.Input == nil {
		return result.NewFailedResult("input cannot be nil", "")
	}

	messageVal := config.Input["message"]

	var message string
	if messageVal.Type().Equals(cty.String) {
		message = messageVal.AsString()
	} else {
		message = hclutil.FormatCtyValueToIndentedString(config.Input["message"], 0, 4)
	}

	r := result.NewNotChangedResult(cty.EmptyObjectVal)
	r.AddMessages(message)

	return r
}
