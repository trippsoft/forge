// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import (
	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// NewChanged creates a new result.ModuleResult representing a successful execution that made changes.
func NewChanged(o cty.Value) (*result.ModuleResult, error) {
	output, err := json.Marshal(o, cty.DynamicPseudoType)
	if err != nil {
		return nil, err
	}

	return &result.ModuleResult{
		Result: &result.ModuleResult_Success{
			Success: &result.ModuleSuccess{
				Changed: true,
				Output:  output,
			},
		},
		Messages: []string{},
		Warnings: []string{},
	}, nil
}

// NewNotChanged creates a new result.ModuleResult representing a successful execution that made no changes.
func NewNotChanged(o cty.Value) (*result.ModuleResult, error) {
	output, err := json.Marshal(o, cty.DynamicPseudoType)
	if err != nil {
		return nil, err
	}

	return &result.ModuleResult{
		Result: &result.ModuleResult_Success{
			Success: &result.ModuleSuccess{
				Changed: false,
				Output:  output,
			},
		},
		Messages: []string{},
		Warnings: []string{},
	}, nil
}

// NewModuleFailure creates a new result.ModuleResult representing a failed execution.
func NewFailure(err error, details string) *result.ModuleResult {
	return &result.ModuleResult{
		Result: &result.ModuleResult_Failure{
			Failure: &result.ModuleFailure{
				Error:  err.Error(),
				Detail: details,
			},
		},
		Messages: []string{},
		Warnings: []string{},
	}
}

// PluginModule defines the interface that all plugin modules must implement.
type PluginModule interface {
	// Name returns the name of the plugin module.
	Name() string
	// Type returns the type of the plugin module.
	Type() plugin.ModuleType
	// InputSpec returns the input specification for the plugin module.
	InputSpec() *hclspec.Spec
	// RunModule executes the plugin module with the given input and returns the result.
	RunModule(hostInfo *info.HostInfo, input map[string]cty.Value, whatIf bool) *result.ModuleResult
}
