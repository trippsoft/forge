// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/hclfunction"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type WorkflowContext struct {
	ui          ui.UI
	inventory   *inventory.Inventory
	debug       bool
	hostVars    map[string]cty.Value
	failedHosts *util.Set[*inventory.Host]
}

// LoadHostVars loads the variables for each host in the inventory into the WorkflowContext.
func (wc *WorkflowContext) LoadHostVars() {
	wc.hostVars = make(map[string]cty.Value)
	for _, host := range wc.inventory.Hosts() {
		vars := host.Vars()
		if len(vars) > 0 {
			wc.hostVars[host.Name()] = cty.ObjectVal(vars)
		}
	}
}

// IsFailed checks if the given host has been marked as failed in the workflow context.
func (wc *WorkflowContext) IsFailed(host *inventory.Host) bool {
	return wc.failedHosts.Contains(host)
}

// MarkFailed marks the given host as failed in the workflow context.
func (wc *WorkflowContext) MarkFailed(host *inventory.Host) {
	wc.failedHosts.Add(host)
}

// NewWorkflowContext creates a new WorkflowContext with the provided parameters.
func NewWorkflowContext(ui ui.UI, i *inventory.Inventory, debug bool) *WorkflowContext {
	return &WorkflowContext{
		ui:          ui,
		inventory:   i,
		debug:       debug,
		failedHosts: util.NewSet[*inventory.Host](),
	}
}

// HostWorkflowContext represents the context for a specific host within a workflow execution.
type HostWorkflowContext struct {
	*WorkflowContext
	host        *inventory.Host
	evalContext *hcl.EvalContext
}

// LoadEvalContext initializes the HCL evaluation context for the host workflow context.
func (hwc *HostWorkflowContext) LoadEvalContext() error {
	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"hostvars": cty.ObjectVal(hwc.hostVars),
			"info":     cty.ObjectVal(hwc.host.Info().ToMapOfCtyValues()),
		},
		Functions: hclfunction.HCLFunctions(),
	}

	vars, exists := hwc.hostVars[hwc.host.Name()]
	if exists {
		evalCtx.Variables["var"] = vars
	}

	steps, err := hwc.host.GetCurrentContextSteps()
	if err != nil {
		return err
	}

	if len(steps) > 0 {
		evalCtx.Variables["steps"] = cty.ObjectVal(steps)
	}

	procedureInputs := hwc.host.GetCurrentProcedureInputs()
	if procedureInputs != nil {
		evalCtx.Variables["input"] = cty.ObjectVal(procedureInputs)
	}

	hwc.evalContext = evalCtx

	return nil
}

// HostWorkflowContext creates a new HostWorkflowContext for the given host.
func NewHostWorkflowContext(wc *WorkflowContext, host *inventory.Host) *HostWorkflowContext {
	return &HostWorkflowContext{
		WorkflowContext: wc,
		host:            host,
	}
}
