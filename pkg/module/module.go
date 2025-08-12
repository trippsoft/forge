// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"fmt"
	"maps"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	DefaultTimeout = 10 * time.Minute
)

// CommonConfig holds configuration options common to all modules.
type CommonConfig struct {
	transport.Escalation               // Configuration for privilege escalation.
	Timeout              time.Duration // Maximum duration to wait for a command to complete.
}

// Module defines the interface for a module in the system.
// This is being implemented behind an interface to allow for remote modules eventually.
type Module interface {
	// InputSpec returns the specification for the module's input.
	InputSpec() *hclspec.Spec

	// Validate checks if the module input is valid.
	// This validation is done after ensuring the input matches the InputSpec.
	Validate(host *inventory.Host, input map[string]cty.Value) error
}

// Result holds the result of a module execution.
// It includes whether the module made any changes, any error encountered, and the output data.
type Result struct {
	Changed   bool                 // Indicates if the module made any changes.
	Err       error                // Error encountered during module execution, if any.
	ErrDetail string               // Detailed error message, if any.
	Output    map[string]cty.Value // Output data from the module execution.
}

func NewSuccess(changed bool, output map[string]cty.Value) *Result {
	return &Result{
		Changed: changed,
		Output:  output,
	}
}

func NewFailure(err error, errDetail string) *Result {
	return &Result{
		Err:       err,
		ErrDetail: errDetail,
	}
}

// LocalModule extends the Module interface with a Run method for local execution.
type LocalModule interface {
	Module

	// Run executes the module with the provided host and input.
	Run(host *inventory.Host, common *CommonConfig, input map[string]cty.Value) *Result
}

// Registry manages a collection of modules.
// It allows for registering new modules and looking them up by name.
type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
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

// Lookup retrieves a module by its name from the registry.
func (r *Registry) Lookup(name string) (Module, bool) {
	module, exists := r.modules[name]
	return module, exists
}

// Modules returns a copy of the registered modules in the registry.
// This is used mostly for testing purposes to avoid modifying the original map.
func (r *Registry) Modules() map[string]Module {

	if r.modules == nil {
		return make(map[string]Module)
	}

	return maps.Clone(r.modules)
}
