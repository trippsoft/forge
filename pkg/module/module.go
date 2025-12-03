// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	DefaultTimeout = 10 * time.Minute
)

var (
	localModules = map[string]Module{
		"assert":  &AssertModule{},
		"message": &MessageModule{},
	}
)

// RunConfig provides the context for running a module on a specific host.
type RunConfig struct {
	Transport  transport.Transport   // The transport to use for the host.
	HostInfo   *info.HostInfo        // The host info this context is associated with.
	Escalation *transport.Escalation // Privilege escalation configuration for the host.
	WhatIf     bool                  // If true, the module should not make any changes.
	Input      map[string]cty.Value  // Input variables for the module.
}

// Module abstracts local and plugin modules.
type Module interface {
	// InputSpec returns the specification for the module's input.
	InputSpec() *hclspec.Spec

	// Validate checks if the module input is valid.
	// This validation is done after ensuring the input matches the InputSpec.
	Validate(config *RunConfig) error

	// Run executes the module with the provided host and input.
	Run(ctx context.Context, config *RunConfig) *result.Result
}

// Registry manages a collection of modules.
//
// It allows for registering new modules and looking them up by name.
type Registry struct {
	modules map[string]Module
}

// Register adds a new module to the registry.
func (r *Registry) Register(name string, module Module) error {
	if r.modules == nil {
		r.modules = make(map[string]Module)
	}

	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q is already registered", name)
	}

	r.modules[name] = module
	return nil
}

// RegisterLocalModules registers the built-in local modules to the registry.
func (r *Registry) RegisterLocalModules() error {
	var err error
	for name, module := range localModules {
		regErr := r.Register(name, module)
		if regErr != nil {
			err = errors.Join(err, regErr)
		}
	}

	return err
}

// Lookup retrieves a module by its name from the registry.
func (r *Registry) Lookup(name string) (Module, bool) {
	module, exists := r.modules[name]

	return module, exists
}

// Modules returns a copy of the registered modules in the registry.
//
// This is used mostly for testing purposes to avoid modifying the original map.
func (r *Registry) Modules() map[string]Module {
	if r.modules == nil {
		return make(map[string]Module)
	}

	return maps.Clone(r.modules)
}

// NewRegistry creates a new module registry.
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}
