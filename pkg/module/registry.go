// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"errors"
	"maps"
)

// Registry manages a collection of modules.
//
// It allows for registering new modules and looking them up by name.
type Registry struct {
	modules map[string]Module
}

// Register adds a new module to the registry.
//
// If a module with the same name already exists, it will be overwritten.
func (r *Registry) Register(module Module) error {
	if r.modules == nil {
		r.modules = make(map[string]Module)
	}

	if module.ID() == nil {
		return errors.New("module info cannot be nil")
	}

	var name string
	if module.ID().namespace == "forge" && module.ID().pluginName == "core" {
		r.modules[module.ID().moduleName] = module
		return nil
	}

	if module.ID().namespace != "" {
		name = module.ID().namespace + "/"
	}

	if module.ID().pluginName != "" {
		name += module.ID().pluginName + "/"
	}

	name += module.ID().moduleName

	r.modules[name] = module
	return nil
}

// RegisterBuiltinModules registers the built-in core modules to the registry.
func (r *Registry) RegisterBuiltinModules() error {
	var err error
	for _, module := range builtinModules {
		e := r.Register(module)
		if e != nil {
			err = errors.Join(err, e)
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
