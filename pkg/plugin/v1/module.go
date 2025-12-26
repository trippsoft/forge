// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import (
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// ToResult converts the ModuleResult to a result.Result.
func (mr *ModuleResult) ToResult() (*result.Result, error) {
	switch mr.Result.(type) {
	case *ModuleResult_Failure:
		r := result.NewFailure(errors.New(mr.GetFailure().Error), mr.GetFailure().Details)
		r.Messages = mr.Messages
		r.Warnings = mr.Warnings
		return r, nil
	case *ModuleResult_Success:
		output := make(map[string]cty.Value, len(mr.GetSuccess().Output))
		for key, value := range mr.GetSuccess().Output {
			ctyValue, err := json.Unmarshal(value, cty.DynamicPseudoType)
			if err != nil {
				return nil, err
			}
			output[key] = ctyValue
		}
		r := result.NewSuccess(mr.GetSuccess().Changed, output)
		r.Messages = mr.Messages
		r.Warnings = mr.Warnings
		return r, nil
	}

	return nil, errors.New("unknown ModuleResult type")
}

// NewModuleSuccess creates a new ModuleResult representing a successful execution.
func NewModuleSuccess(changed bool, output map[string]cty.Value) (*ModuleResult, error) {
	outputPB := make(map[string][]byte, len(output))
	for key, value := range output {
		jsonBytes, err := json.Marshal(value, cty.DynamicPseudoType)
		if err != nil {
			return nil, err
		}
		outputPB[key] = jsonBytes
	}

	return &ModuleResult{
		Result: &ModuleResult_Success{
			Success: &ModuleSuccess{
				Changed: changed,
				Output:  outputPB,
			},
		},
		Messages: []string{},
		Warnings: []string{},
	}, nil
}

// NewModuleFailure creates a new ModuleResult representing a failed execution.
func NewModuleFailure(err error, details string) *ModuleResult {
	return &ModuleResult{
		Result: &ModuleResult_Failure{
			Failure: &ModuleFailure{
				Error:   err.Error(),
				Details: details,
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
	Type() ModuleType
	// InputSpec returns the input specification for the plugin module.
	InputSpec() *hclspec.Spec
	// RunModule executes the plugin module with the given input and returns the result.
	RunModule(hostInfo *info.HostInfo, input map[string]cty.Value, whatIf bool) *ModuleResult
}
