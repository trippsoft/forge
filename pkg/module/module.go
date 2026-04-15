// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

var (
	builtinModules = []Module{}
)

// ModuleID provides identification for a module.
type ModuleID struct {
	namespace  string
	pluginName string
	moduleName string
}

// Namespace returns the namespace of the module.
//
// This will be blank for local modules that have no plugin.
func (id *ModuleID) Namespace() string {
	return id.namespace
}

// PluginName returns the name of the plugin that provides the module.
//
// This will be blank for local modules that have no plugin.
func (id *ModuleID) PluginName() string {
	return id.pluginName
}

// ModuleName returns the name of the module.
func (id *ModuleID) ModuleName() string {
	return id.moduleName
}

// NewModuleID creates a new ModuleID with the given namespace, plugin name, and module name.
func NewModuleID(namespace, pluginName, moduleName string) *ModuleID {
	return &ModuleID{
		namespace:  namespace,
		pluginName: pluginName,
		moduleName: moduleName,
	}
}

// RunConfig provides configuration for running a module.
type RunConfig struct {
	// Transport is the transport to use for running the module.
	Transport transport.Transport

	// HostInfo is the host information to use for running the module.
	HostInfo *info.HostInfo

	// Escalation is the escalation to use for running the module.
	Escalation *transport.Escalation

	// WhatIf indicates whether to run the module in "what if" mode, which simulates the execution without making any
	// changes.
	WhatIf bool

	// Input is the input variables to pass to the module.
	Input map[string]cty.Value
}

// Module abstracts local and plugin modules.
type Module interface {
	// ID returns the module identity.
	ID() *ModuleID

	// InputSpec returns the specification for the module's input.
	InputSpec() *hclspec.Spec

	// Run runs the module with the given context and configuration.
	Run(ctx context.Context, config *RunConfig) result.Result
}
