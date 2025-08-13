package workflow

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/hclfunction"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/zclconf/go-cty/cty"
)

// workflowContext holds the context for a workflow.
type workflowContext struct {
	ui        ui.UI                // ui represents the user interface used by the workflow.
	inventory *inventory.Inventory // inventory represents the inventory used by the workflow.
	hostVars  map[string]cty.Value // hostVars holds the variables for each host.
}

// WorkflowContext creates a new workflowContext.
func WorkflowContext(ui ui.UI, inventory *inventory.Inventory) *workflowContext {
	return &workflowContext{
		ui:        ui,
		inventory: inventory,
	}
}

func (ctx *workflowContext) LoadHostVars() {

	ctx.hostVars = make(map[string]cty.Value)

	for _, host := range ctx.inventory.Hosts() {
		vars := host.Vars()
		if len(vars) > 0 {
			ctx.hostVars[host.Name()] = cty.ObjectVal(vars)
		}
	}
}

// hostWorkflowContext holds the context for a host workflow.
type hostWorkflowContext struct {
	*workflowContext
	host        *inventory.Host  // host represents the current host being processed in the workflow.
	evalContext *hcl.EvalContext // evalContext is the evaluation context for the workflow.
}

// HostWorkflowContext creates a new hostWorkflowContext.
func HostWorkflowContext(ctx *workflowContext, host *inventory.Host) *hostWorkflowContext {
	return &hostWorkflowContext{
		workflowContext: ctx,
		host:            host,
	}
}

func (ctx *hostWorkflowContext) LoadEvalContext() error {

	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"hostvars": cty.ObjectVal(ctx.hostVars),
			"info":     cty.ObjectVal(ctx.host.Info().ToMapOfCtyValues()),
		},
		Functions: hclfunction.HCLFunctions(),
	}

	vars, exists := ctx.hostVars[ctx.host.Name()]
	if exists {
		evalCtx.Variables["var"] = vars
	}

	tasks, err := ctx.host.GetCurrentContextTasks()
	if err != nil {
		return err
	}

	if len(tasks) > 0 {
		evalCtx.Variables["tasks"] = cty.ObjectVal(tasks)
	}

	procedureInputs := ctx.host.GetCurrentProcedureInputs()
	if procedureInputs != nil {
		evalCtx.Variables["input"] = cty.ObjectVal(procedureInputs)
	}

	ctx.evalContext = evalCtx

	return nil
}
