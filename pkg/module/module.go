// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	DefaultTimeout = 10 * time.Minute
)

var (
	localModules = []Module{
		&AssertModule{},
		&MessageModule{},
	}
)

// ModuleInfo provides metadata about a module.
type ModuleInfo struct {
	namespace  string
	pluginName string
	name       string
}

// NewModuleInfo creates a new ModuleInfo instance.
func NewModuleInfo(namespace, pluginName, name string) *ModuleInfo {
	return &ModuleInfo{
		namespace:  namespace,
		pluginName: pluginName,
		name:       name,
	}
}

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
	// Info returns the module information.
	Info() *ModuleInfo

	// InputSpec returns the specification for the module's input.
	InputSpec() *hclspec.Spec

	// Run runs the module with the given context and configuration.
	Run(ctx context.Context, config *RunConfig) *result.Result
}

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

	if module.Info() == nil {
		return errors.New("module info cannot be nil")
	}

	var name string
	if module.Info().namespace != "" {
		name = module.Info().namespace + "/"
	}

	if module.Info().pluginName != "" {
		name += module.Info().pluginName + "/"
	}

	name += module.Info().name

	r.modules[name] = module
	return nil
}

// RegisterCoreModules registers the built-in core modules to the registry.
func (r *Registry) RegisterCoreModules() error {
	var err error
	for _, module := range localModules {
		e := r.Register(module)
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// RegisterPluginModules registers modules provided by plugins to the registry.
func (r *Registry) RegisterPluginModules() error {
	var err error
	e := r.registerPluginModulesAtBasePath(plugin.SharedPluginBasePath)
	if e != nil {
		err = errors.Join(err, e)
	}

	e = r.registerPluginModulesAtBasePath(plugin.UserPluginBasePath)
	if e != nil {
		err = errors.Join(err, e)
	}

	return err
}

func (r *Registry) registerPluginModulesAtBasePath(basePath string) error {
	fileInfo, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read base path info for %q: %w", basePath, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("base path %q is not a directory", basePath)
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read directory entries for %q: %w", basePath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directory entries
		}

		e := r.registerPluginModulesAtNamespacePath(basePath, entry.Name())
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func (r *Registry) registerPluginModulesAtNamespacePath(basePath, namespace string) error {
	namespacePath := path.Join(basePath, namespace)
	entries, err := os.ReadDir(namespacePath)
	if err != nil {
		return fmt.Errorf("failed to read namespace directory entries for %q: %w", namespacePath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directory entries
		}

		e := r.registerPluginModulesAtPluginPath(namespacePath, namespace, entry.Name())
		if e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

func (r *Registry) registerPluginModulesAtPluginPath(basePath, namespace, pluginName string) error {
	if namespace == "forge" && pluginName == "discover" {
		// Skip the discover plugin.
		return nil
	}

	connection, cleanup, err := transport.LocalTransport.StartPlugin(
		context.Background(),
		basePath,
		namespace,
		pluginName,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to start plugin %q/%q: %w", namespace, pluginName, err)
	}

	defer connection.Close()
	defer cleanup()

	client := pluginv1.NewPluginV1ServiceClient(connection)
	response, err := client.GetModules(context.Background(), &pluginv1.GetModulesRequest{})
	if err != nil {
		return fmt.Errorf("failed to get modules from plugin %q/%q: %w", namespace, pluginName, err)
	}

	for name, moduleSpec := range response.Modules {
		spec, err := moduleSpec.Spec.ToSpec()
		if err != nil {
			return fmt.Errorf("failed to parse module spec for %q/%q/%q: %w", namespace, pluginName, name, err)
		}

		moduleInfo := NewModuleInfo(namespace, pluginName, name)

		var module Module
		switch moduleSpec.Type {
		case pluginv1.ModuleType_LOCAL:
			module = NewLocalPluginModule(basePath, moduleInfo, spec)
		case pluginv1.ModuleType_REMOTE:
			module = NewRemotePluginModule(basePath, moduleInfo, spec)
		default:
			return fmt.Errorf(
				"unknown module type %q for %q/%q/%q",
				moduleSpec.Type.String(),
				namespace,
				pluginName,
				name,
			)
		}

		e := r.Register(module)
		if e != nil {
			return errors.Join(
				err,
				fmt.Errorf("failed to register module %q/%q/%q: %w", namespace, pluginName, name, e),
			)
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
