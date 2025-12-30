// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"errors"
	"slices"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

// Process represents a process within a workflow.
//
// This represents a parsed process block from a workflow file.
type Process struct {
	name         string
	discoverInfo bool
	allTargets   []*inventory.Host
	steps        []Step
}

// Name returns the name of the process.
func (p *Process) Name() string {
	return p.name
}

// DiscoverInfo indicates whether the process is set to discover information.
func (p *Process) DiscoverInfo() bool {
	return p.discoverInfo
}

// AllTargets returns a clone of the slice of all targets associated with the process.
//
// This is done to prevent external modification of the internal state.
// This method is primarily intended for integration tests and debugging purposes.
func (p *Process) AllTargets() []*inventory.Host {
	allTargets := slices.Clone(p.allTargets)
	return allTargets
}

// Steps returns a clone of the slice of all steps in the process.
//
// This is done to prevent external modification of the internal state.
// This method is primarily intended for integration tests and debugging purposes.
func (p *Process) Steps() []Step {
	steps := slices.Clone(p.steps)
	return steps
}

// Run executes the process using the provided WorkflowContext.
func (p *Process) Run(wc *WorkflowContext) (map[string]map[string]cty.Value, error) {
	wc.ui.PrintHeader(ui.HeaderLevel1, "PROCESS - ", p.name)

	err := p.discoverInfoForTargets(wc)

	outputs := make(map[string]map[string]cty.Value)

	for _, step := range p.steps {
		output, stepErr := step.Run(wc)
		outputs[step.ID()] = output
		err = errors.Join(err, stepErr)
	}

	return outputs, err
}

func (p *Process) discoverInfoForTargets(wc *WorkflowContext) error {
	if !p.discoverInfo {
		return nil
	}

	wc.ui.PrintHeader(ui.HeaderLevel2, "Discovering Host Info", "")

	errChannel := make(chan error)
	var err error
	for _, target := range p.allTargets {
		go func(host *inventory.Host) {
			result := host.PopulateInfo()

			var e error
			if result.Err != nil {
				e = result.Err
			}

			wc.ui.PrintHostResult(host.Name(), result)
			errChannel <- e
		}(target)
	}

	for range p.allTargets {
		e := <-errChannel
		err = errors.Join(err, e)
	}

	return err
}

// ProcessBuilder is used to build a Process instance during parsing.
type ProcessBuilder struct {
	common       *StepCommonConfig
	escalate     *StepEscalateConfig
	discoverInfo bool

	steps []StepBuilder
}

// WithCommon sets the common configuration for the process.
func (pb *ProcessBuilder) WithCommon(common *StepCommonConfig) *ProcessBuilder {
	pb.common = common
	return pb
}

// WithEscalate sets the escalation configuration for the process.
func (pb *ProcessBuilder) WithEscalate(escalate *StepEscalateConfig) *ProcessBuilder {
	pb.escalate = escalate
	return pb
}

// WithDiscoverInfo sets the discover_info flag for the process.
func (pb *ProcessBuilder) WithDiscoverInfo(discoverInfo bool) *ProcessBuilder {
	pb.discoverInfo = discoverInfo
	return pb
}

// AddStep adds a StepBuilder to the ProcessBuilder.
func (pb *ProcessBuilder) AddStep(sb StepBuilder) *ProcessBuilder {
	pb.steps = append(pb.steps, sb)
	sb.WithProcessCommon(pb.common).WithProcessEscalate(pb.escalate)
	return pb
}

// Build constructs and returns the Process instance.
func (pb *ProcessBuilder) Build() (*Process, error) {
	steps := make([]Step, 0, len(pb.steps))
	allTargetsSet := util.NewSet[*inventory.Host]()
	var err error
	for _, sb := range pb.steps {
		step, stepErr := sb.Build()
		if stepErr != nil || err != nil {
			err = errors.Join(err, stepErr)
			continue
		}

		steps = append(steps, step)
		for _, target := range sb.AllTargets() {
			allTargetsSet.Add(target)
		}
	}

	if err != nil {
		return nil, err
	}

	return &Process{
		name:         pb.common.name,
		discoverInfo: pb.discoverInfo,
		allTargets:   allTargetsSet.Items(),
		steps:        steps,
	}, nil
}

// NewProcessBuilder creates a new instance of ProcessBuilder.
func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{
		discoverInfo: true,
		steps:        []StepBuilder{},
	}
}
