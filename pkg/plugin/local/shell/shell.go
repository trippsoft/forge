package shell

import (
	"context"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hcldec.ObjectSpec{
		"command": &hcldec.AttrSpec{
			Type:     cty.String,
			Required: true,
		},
	}

	_ plugin.LocalPlugin = &Plugin{} // Ensure Plugin implements the plugin.LocalPlugin interface.
)

type Plugin struct{}

// InputSpec implements plugin.Plugin.
func (s *Plugin) InputSpec() hcldec.ObjectSpec {
	return inputSpec
}

// Validate implements plugin.Plugin.
func (s *Plugin) Validate(host *inventory.Host, input map[string]cty.Value) error {
	return nil // No specific validation needed for this plugin.
}

// Run implements plugin.Plugin.
func (s *Plugin) Run(host *inventory.Host, common *plugin.CommonConfig, input map[string]cty.Value) *plugin.Result {

	t := host.Transport()

	command := input["command"].AsString()
	cmd, err := t.NewCommand(command, common.Escalation)
	if err != nil {
		return plugin.NewFailure(err)
	}

	timeout := common.Timeout

	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
	}

	stdout, stderr, err := cmd.OutputWithError(ctx)
	if err != nil {
		return plugin.NewFailure(err)
	}

	output := map[string]cty.Value{
		"stdout": cty.StringVal(stdout),
		"stderr": cty.StringVal(stderr),
	}

	return plugin.NewSuccess(true, output)
}
