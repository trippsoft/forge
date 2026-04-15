// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"errors"
	"slices"

	"github.com/zclconf/go-cty/cty"
)

// Workflow represents a workflow that contains one or more processes.
//
// This is the parsed representation of a workflow file.
type Workflow struct {
	processes []*Process
}

// Run executes the workflow using the provided WorkflowContext.
func (w *Workflow) Run(wc *WorkflowContext) ([]map[string]map[string]cty.Value, error) {
	wc.inventory.ClearSteps() // Clear any previous step contexts and procedure inputs

	outputs := make([]map[string]map[string]cty.Value, 0, len(w.processes))
	var err error
	for _, process := range w.processes {
		var output map[string]map[string]cty.Value
		output, err = process.Run(wc)
		outputs = append(outputs, output)
		if err != nil {
			break
		}
	}

	for _, host := range wc.inventory.Hosts() {
		host.Transport().Close()
	}

	return outputs, err
}

// Processes returns a clone of the slice of all processes in the workflow.
//
// This is done to prevent external modification of the internal state.
// This method is primarily intended for integration tests and debugging purposes.
func (w *Workflow) Processes() []*Process {
	processes := slices.Clone(w.processes)
	return processes
}

// WorkflowBuilder is used to build a Workflow instance during parsing.
type WorkflowBuilder struct {
	processes []*ProcessBuilder
}

// AddProcess adds a ProcessBuilder to the WorkflowBuilder.
func (wb *WorkflowBuilder) AddProcess(pb ...*ProcessBuilder) *WorkflowBuilder {
	wb.processes = append(wb.processes, pb...)
	return wb
}

// Build constructs and returns the Workflow instance.
func (wb *WorkflowBuilder) Build() (*Workflow, error) {
	processes := make([]*Process, 0, len(wb.processes))
	var err error
	for _, pb := range wb.processes {
		process, processErr := pb.Build()
		if processErr != nil || err != nil {
			err = errors.Join(err, processErr)
			continue
		}

		processes = append(processes, process)
	}

	if err != nil {
		return nil, err
	}

	return &Workflow{
		processes: processes,
	}, nil
}

// NewWorkflowBuilder creates a new instance of WorkflowBuilder.
func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		processes: []*ProcessBuilder{},
	}
}
