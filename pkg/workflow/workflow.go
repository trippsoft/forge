// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import "slices"

// Workflow represents a workflow.
type Workflow struct {
	processes []*Process // processes represents the list of processes in the workflow.
}

// NewWorkflow creates a new Workflow.
func NewWorkflow(processes ...*Process) *Workflow {
	return &Workflow{
		processes: processes,
	}
}

// Processes returns a copy of the processes in the workflow.
// This is used primarily for testing purposes.
func (w *Workflow) Processes() []*Process {
	processes := slices.Clone(w.processes)
	return processes
}

// Run executes the workflow.
func (w *Workflow) Run(ctx *workflowContext) error {

	ctx.inventory.ClearSteps() // Clear any existing steps before running the workflow.

	var err error
	for _, process := range w.processes {
		err = process.Run(ctx)
		if err != nil {
			break
		}
	}

	for _, host := range ctx.inventory.Hosts() {
		host.Transport().Close()
	}

	return err
}
