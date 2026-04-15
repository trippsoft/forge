// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// LocalPluginModule defines a module from a gRPC plugin that runs locally on the controller.
type LocalPluginModule struct {
	basePath string
	id       *ModuleID
	spec     *hclspec.Spec
}

// Info implements Module.
func (m *LocalPluginModule) ID() *ModuleID {
	return m.id
}

// InputSpec implements Module.
func (m *LocalPluginModule) InputSpec() *hclspec.Spec {
	return m.spec
}

// Run implements Module.
func (m *LocalPluginModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	session, err := transport.LocalTransport.StartPluginSession(
		ctx,
		m.basePath,
		m.id.namespace,
		m.id.pluginName,
		config.Escalation,
	)

	if err != nil {
		return result.NewFailure(err, "")
	}

	defer session.Close()

	input := make(map[string][]byte, len(config.Input))
	for k, v := range config.Input {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return result.NewFailure(err, err.Error())
		}

		input[k] = value
	}

	request := &pluginv1.RunModuleRequest{
		ModuleName: m.id.moduleName,
		HostInfo:   config.HostInfo,
		Input:      input,
		WhatIf:     config.WhatIf,
	}

	err = plugin.Write(session.Stdin(), request)
	if err != nil {
		return result.NewFailure(err, "")
	}

	response := &pluginv1.RunModuleResponse{}

	err = plugin.Read(session.Stdout(), response)
	if err != nil {
		return result.NewFailure(err, "")
	}

	return response.Result.ToResult()
}

// NewLocalPluginModule creates a new LocalPluginModule.
func NewLocalPluginModule(basePath string, id *ModuleID, spec *hclspec.Spec) Module {
	return &LocalPluginModule{
		basePath: basePath,
		id:       id,
		spec:     spec,
	}
}

// RemotePluginModule defines a module from a gRPC plugin that runs on the managed host.
type RemotePluginModule struct {
	basePath string
	id       *ModuleID
	spec     *hclspec.Spec
}

// Info implements Module.
func (m *RemotePluginModule) ID() *ModuleID {
	return m.id
}

// InputSpec implements Module.
func (m *RemotePluginModule) InputSpec() *hclspec.Spec {
	return m.spec
}

// Run implements Module.
func (m *RemotePluginModule) Run(ctx context.Context, config *RunConfig) *result.Result {
	session, err := config.Transport.StartPluginSession(
		ctx,
		m.basePath,
		m.id.namespace,
		m.id.pluginName,
		config.Escalation,
	)

	if err != nil {
		return result.NewFailure(err, "")
	}

	defer session.Close()

	input := make(map[string][]byte, len(config.Input))
	for k, v := range config.Input {
		value, err := json.Marshal(v, cty.DynamicPseudoType)
		if err != nil {
			return result.NewFailure(err, err.Error())
		}

		input[k] = value
	}

	request := &pluginv1.RunModuleRequest{
		ModuleName: m.id.moduleName,
		HostInfo:   config.HostInfo,
		Input:      input,
		WhatIf:     config.WhatIf,
	}

	err = plugin.Write(session.Stdin(), request)
	if err != nil {
		return result.NewFailure(err, "")
	}

	response := &pluginv1.RunModuleResponse{}

	err = plugin.Read(session.Stdout(), response)
	if err != nil {
		return result.NewFailure(err, "")
	}

	return response.Result.ToResult()
}

// NewRemotePluginModule creates a new RemotePluginModule.
func NewRemotePluginModule(basePath string, id *ModuleID, spec *hclspec.Spec) Module {
	return &RemotePluginModule{
		basePath: basePath,
		id:       id,
		spec:     spec,
	}
}
