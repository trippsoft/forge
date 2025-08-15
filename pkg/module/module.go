// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	DefaultTimeout = 10 * time.Minute
)

// RunConfig provides the context for running a module on a specific host.
type RunConfig struct {
	Transport  transport.Transport  // The transport to use for the host.
	HostInfo   *info.HostInfo       // The host info this context is associated with.
	Escalation transport.Escalation // Privilege escalation configuration for the host.
	WhatIf     bool                 // If true, the module should not make any changes.
	Input      map[string]cty.Value // Input variables for the module.
}

// Module defines the interface for a module in the system.
// This is being implemented behind an interface to allow for remote modules eventually.
type Module interface {
	// InputSpec returns the specification for the module's input.
	InputSpec() *hclspec.Spec

	// Validate checks if the module input is valid.
	// This validation is done after ensuring the input matches the InputSpec.
	Validate(config *RunConfig) error

	// Run executes the module with the provided host and input.
	Run(ctx context.Context, config *RunConfig) *Result
}

// Local wraps a local module implementation.
type Local struct {
	module Module
}

// NewLocal creates a new Local module.
func NewLocal(module Module) Module {
	return &Local{
		module: module,
	}
}

func (l *Local) InputSpec() *hclspec.Spec {
	return l.module.InputSpec()
}

func (l *Local) Validate(config *RunConfig) error {
	return l.module.Validate(config)
}

func (l *Local) Run(ctx context.Context, config *RunConfig) *Result {

	outputChannel := make(chan *Result)
	go func(ctx context.Context) {
		outputChannel <- l.module.Run(ctx, config)
	}(ctx)

	select {
	case <-ctx.Done():
		return NewFailure(ctx.Err(), "module run timed out")
	case result := <-outputChannel:
		return result
	}
}

type MockModule struct {
	inputSpec    *hclspec.Spec
	validateFunc func(config *RunConfig) error
	Result       *Result
}

func NewMockModule(spec *hclspec.Spec, validate func(config *RunConfig) error) *MockModule {
	return &MockModule{
		inputSpec:    spec,
		validateFunc: validate,
	}
}

func (m *MockModule) InputSpec() *hclspec.Spec {
	return m.inputSpec
}

func (m *MockModule) Validate(config *RunConfig) error {

	if m.validateFunc == nil {
		return nil
	}

	return m.validateFunc(config)
}

func (m *MockModule) Run(ctx context.Context, config *RunConfig) *Result {
	return m.Result
}

// Result holds the result of a module execution.
// It includes whether the module made any changes, any error encountered, and the output data.
type Result struct {
	Failed         bool                 // Indicates if the module execution failed.
	IgnoredFailure bool                 // Indicates if the failure was ignored.
	Skipped        bool                 // Indicates if the module was skipped.
	Changed        bool                 // Indicates if the module made any changes.
	Err            error                // Error encountered during module execution, if any.
	ErrDetail      string               // Detailed error message, if any.
	Output         map[string]cty.Value // Output data from the module execution.
	Warning        string               // Warning message, if any.
	Message        string               // Informational message, if any.
}

func NewSuccess(changed bool, output map[string]cty.Value) *Result {
	return &Result{
		Changed: changed,
		Output:  output,
	}
}

func NewSkipped() *Result {
	return &Result{
		Skipped: true,
	}
}

func NewFailure(err error, errDetail string) *Result {
	return &Result{
		Err:       err,
		ErrDetail: errDetail,
		Failed:    true,
	}
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
