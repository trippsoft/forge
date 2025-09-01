package local

import (
	"context"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/module/assert"
	"github.com/trippsoft/forge/pkg/module/message"
	"github.com/trippsoft/forge/pkg/module/pkg"
	"github.com/trippsoft/forge/pkg/module/service"
	"github.com/trippsoft/forge/pkg/module/shell"
)

// Local wraps a local module implementation.
type Local struct {
	module module.Module
}

// NewLocal creates a new Local module.
func NewLocal(module module.Module) module.Module {
	return &Local{
		module: module,
	}
}

// InputSpec implements Module.
func (l *Local) InputSpec() *hclspec.Spec {
	return l.module.InputSpec()
}

// Validate implements Module.
func (l *Local) Validate(config *module.RunConfig) error {
	return l.module.Validate(config)
}

// Run implements Module.
func (l *Local) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	outputChannel := make(chan *module.Result)
	go func(ctx context.Context) {
		outputChannel <- l.module.Run(ctx, config)
	}(ctx)

	select {
	case <-ctx.Done():
		return module.NewFailure(ctx.Err(), "module run timed out")
	case result := <-outputChannel:
		return result
	}
}

func RegisterLocalModules(moduleRegistry *module.Registry) {
	moduleRegistry.Register("assert", NewLocal(&assert.Module{}))

	moduleRegistry.Register("message", NewLocal(&message.Module{}))

	moduleRegistry.Register("dnf", NewLocal(&pkg.DNFModule{}))
	moduleRegistry.Register("dnf_info", NewLocal(&pkg.DNFInfoModule{}))
	moduleRegistry.Register("pkg", NewLocal(&pkg.PkgModule{}))
	moduleRegistry.Register("pkg_info", NewLocal(&pkg.PkgInfoModule{}))

	moduleRegistry.Register("systemd_service", NewLocal(&service.SystemdServiceModule{}))
	moduleRegistry.Register("service", NewLocal(&service.ServiceModule{}))

	moduleRegistry.Register("shell", (&shell.Module{}))
}
