package plugin

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/zclconf/go-cty/cty"
)

// Plugin defines the interface for a plugin in the system.
// This is being implemented behind an interface to allow for remote plugins eventually.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string
	// FullName returns the full name of the plugin, which may include a namespace or other identifiers.
	FullName() string
	// Version returns the version of the plugin.
	Version() string
	// Description returns a brief description of the plugin.
	Description() string
	// Author returns the author of the plugin.
	Author() string

	// InputSpec returns the specification for the plugin's input.
	InputSpec() hcldec.ObjectSpec

	// Validate checks if the plugin input is valid.
	Validate(input map[string]cty.Value) error

	// Run executes the plugin with the provided host and input.
	Run(host *inventory.Host, input map[string]cty.Value) (map[string]cty.Value, error)
}

type Registry struct {
	plugins map[string]Plugin
}
