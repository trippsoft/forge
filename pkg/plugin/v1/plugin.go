// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import (
	"fmt"
	"io"
	"os"

	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// PluginV1 is the implementation of the plugin interface for API version 1.
type PluginV1 struct {
	namespace  string
	pluginName string

	modules map[string]PluginModule
}

func (p *PluginV1) Run() error {
	_, err := fmt.Fprintln(os.Stderr, plugin.PluginReadyMessage)
	if err != nil {
		return err
	}

	if len(os.Args) > 1 && os.Args[1] == "metadata" {
		return p.handleMetadataRequest()
	}

	return p.handleRunModuleRequest()
}

func (p *PluginV1) handleMetadataRequest() error {
	var request plugin.MetadataRequest
	err := plugin.Read(os.Stdin, &request)
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	modules := make(map[string]*plugin.ModuleSpec, len(p.modules))
	for name, mod := range p.modules {
		spec, err := mod.InputSpec().ToProtobuf()
		if err != nil {
			return err
		}

		modules[name] = &plugin.ModuleSpec{
			Type: mod.Type(),
			Spec: spec,
		}
	}

	response := &plugin.MetadataResponse{
		ApiVersion: 1,
		Namespace:  p.namespace,
		PluginName: p.pluginName,
		Modules:    modules,
	}

	return plugin.Write(os.Stdout, response)
}

func (p *PluginV1) handleRunModuleRequest() error {
	var request RunModuleRequest
	err := plugin.Read(os.Stdin, &request)
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	mod, ok := p.modules[request.ModuleName]
	if !ok {
		return fmt.Errorf("unknown module: %s", request.ModuleName)
	}

	input := make(map[string]cty.Value, len(request.Input))
	for k, v := range request.Input {
		val, err := json.Unmarshal(v, cty.DynamicPseudoType)
		if err != nil {
			return err
		}

		input[k] = val
	}

	r := mod.RunModule(request.HostInfo, input, request.WhatIf)

	response := &RunModuleResponse{Result: r}

	return plugin.Write(os.Stdout, response)
}

// NewPluginV1 creates a new PluginV1 instance with the given namespace, plugin name, and modules.
func NewPluginV1(namespace, pluginName string, m ...PluginModule) *PluginV1 {
	modules := make(map[string]PluginModule, len(m))
	for _, mod := range m {
		modules[mod.Name()] = mod
	}

	return &PluginV1{
		namespace:  namespace,
		pluginName: pluginName,
		modules:    modules,
	}
}
