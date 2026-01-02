// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"errors"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// LocalPluginModule defines a module from a gRPC plugin that runs locally on the controller.
type LocalPluginModule struct {
	basePath string
	info     *ModuleInfo
	spec     *hclspec.Spec
}

// Info implements Module.
func (m *LocalPluginModule) Info() *ModuleInfo {
	return m.info
}

// InputSpec implements Module.
func (m *LocalPluginModule) InputSpec() *hclspec.Spec {
	return m.spec
}

// Run implements Module.
func (m *LocalPluginModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	connection, cleanup, err := transport.LocalTransport.StartPlugin(
		ctx,
		m.basePath,
		m.info.namespace,
		m.info.pluginName,
		config.Escalation,
	)

	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	defer connection.Close()
	defer cleanup()

	client := pluginv1.NewPluginV1ServiceClient(connection)

	input := make(map[string][]byte, len(config.Input))
	for k, v := range config.Input {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return result.NewFailure(err, err.Error())
		}
		input[k] = value
	}

	request, err := createRunModuleRequest(m.info.name, config)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	response, err := client.RunModule(ctx, request)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	r, err := createResultFromRunModuleResponse(response.Result)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	return r
}

// NewLocalPluginModule creates a new LocalPluginModule.
func NewLocalPluginModule(basePath string, info *ModuleInfo, spec *hclspec.Spec) Module {
	return &LocalPluginModule{
		basePath: basePath,
		info:     info,
		spec:     spec,
	}
}

// RemotePluginModule defines a module from a gRPC plugin that runs on the managed host.
type RemotePluginModule struct {
	basePath string
	info     *ModuleInfo
	spec     *hclspec.Spec
}

// Info implements Module.
func (m *RemotePluginModule) Info() *ModuleInfo {
	return m.info
}

// InputSpec implements Module.
func (m *RemotePluginModule) InputSpec() *hclspec.Spec {
	return m.spec
}

// Run implements Module.
func (m *RemotePluginModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	connection, cleanup, err := config.Transport.StartPlugin(
		ctx,
		m.basePath,
		m.info.namespace,
		m.info.pluginName,
		config.Escalation,
	)

	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	defer connection.Close()
	defer cleanup()

	client := pluginv1.NewPluginV1ServiceClient(connection)

	input := make(map[string][]byte, len(config.Input))
	for k, v := range config.Input {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return result.NewFailure(err, err.Error())
		}
		input[k] = value
	}

	request, err := createRunModuleRequest(m.info.name, config)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	response, err := client.RunModule(ctx, request)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	r, err := createResultFromRunModuleResponse(response.Result)
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	return r
}

// NewRemotePluginModule creates a new RemotePluginModule.
func NewRemotePluginModule(basePath string, info *ModuleInfo, spec *hclspec.Spec) Module {
	return &RemotePluginModule{
		basePath: basePath,
		info:     info,
		spec:     spec,
	}
}

func createRunModuleRequest(moduleName string, config *RunConfig) (*pluginv1.RunModuleRequest, error) {
	input := make(map[string][]byte, len(config.Input))
	for k, v := range config.Input {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return nil, err
		}
		input[k] = value
	}

	request := &pluginv1.RunModuleRequest{
		ModuleName: moduleName,
		HostInfo:   config.HostInfo,
		Input:      input,
		WhatIf:     config.WhatIf,
	}

	return request, nil
}

func createResultFromRunModuleResponse(response *pluginv1.ModuleResult) (*result.Result, error) {
	switch response.Result.(type) {
	case *pluginv1.ModuleResult_Success:
		success := response.GetSuccess()
		output, err := json.Unmarshal(success.Output, cty.DynamicPseudoType)
		if err != nil {
			return nil, err
		}
		r := result.NewSuccess(success.Changed, output)
		r.Messages = response.Messages
		r.Warnings = response.Warnings
		return r, nil
	case *pluginv1.ModuleResult_Failure:
		failure := response.GetFailure()
		r := result.NewFailure(errors.New(failure.Error), failure.Details)
		r.Messages = response.Messages
		r.Warnings = response.Warnings
		return r, nil
	default:
		return nil, fmt.Errorf("unknown run module response type")
	}
}
