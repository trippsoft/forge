package plugin

import (
	"fmt"
	"maps"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

// CommonConfig holds configuration options common to all plugins.
type CommonConfig struct {
	transport.Escalation               // Configuration for privilege escalation.
	Timeout              time.Duration // Maximum duration to wait for a command to complete.
}

// Plugin defines the interface for a plugin in the system.
// This is being implemented behind an interface to allow for remote plugins eventually.
type Plugin interface {
	// InputSpec returns the specification for the plugin's input.
	InputSpec() hcldec.ObjectSpec

	// Validate checks if the plugin input is valid.
	// This validation is done after ensuring the input matches the InputSpec.
	Validate(input map[string]cty.Value) error
}

// Result holds the result of a plugin execution.
// It includes whether the plugin made any changes, any error encountered, and the output data.
type Result struct {
	Changed bool                 // Indicates if the plugin made any changes.
	Err     error                // Error encountered during plugin execution, if any.
	Output  map[string]cty.Value // Output data from the plugin execution.
}

func NewSuccess(changed bool, output map[string]cty.Value) *Result {
	return &Result{
		Changed: changed,
		Output:  output,
	}
}

func NewFailure(err error) *Result {
	return &Result{
		Err: err,
	}
}

// LocalPlugin extends the Plugin interface with a Run method for local execution.
type LocalPlugin interface {
	Plugin

	// Run executes the plugin with the provided host and input.
	Run(host *inventory.Host, common *CommonConfig, input map[string]cty.Value) *Result
}

// Registry manages a collection of plugins.
// It allows for registering new plugins and looking them up by name.
type Registry struct {
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a new plugin to the registry.
func (r *Registry) Register(name string, plugin Plugin) error {

	if r.plugins == nil {
		r.plugins = make(map[string]Plugin)
	}

	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %q is already registered", name)
	}

	r.plugins[name] = plugin
	return nil
}

// Lookup retrieves a plugin by its name from the registry.
func (r *Registry) Lookup(name string) (Plugin, bool) {
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// Plugins returns a copy of the registered plugins in the registry.
// This is used mostly for testing purposes to avoid modifying the original map.
func (r *Registry) Plugins() map[string]Plugin {

	if r.plugins == nil {
		return make(map[string]Plugin)
	}

	return maps.Clone(r.plugins)
}
