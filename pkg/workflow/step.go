// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"fmt"
	"maps"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/zclconf/go-cty/cty"
)

// Step abstracts a single Step or a procedure in a process.
type Step interface {
	Run(ctx *workflowContext) // Run executes the step with the given workflow context.
}

type StepCommonConfig struct {
	loop *StepLoopConfig

	id   string
	name string

	targets []*inventory.Host

	condition *hcl.Attribute

	execTimeout *hcl.Attribute
	whatIf      *hcl.Attribute

	input map[string]*hcl.Attribute
}

// Loop returns the loop configuration for the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Loop() *StepLoopConfig {
	return s.loop
}

// ID returns the ID of the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) ID() string {
	return s.id
}

// Name returns the name of the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Name() string {
	return s.name
}

// Targets returns the targets of the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Targets() []*inventory.Host {
	return s.targets
}

// Condition returns the condition of the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Condition() *hcl.Attribute {
	return s.condition
}

// ExecTimeout returns the execTimeout of the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) ExecTimeout() *hcl.Attribute {
	return s.execTimeout
}

// Input returns the input variables for the step.
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Input() map[string]*hcl.Attribute {
	input := maps.Clone(s.input)
	return input
}

type StepLoopConfig struct {
	items     *hcl.Attribute
	label     *hcl.Attribute
	condition *hcl.Attribute
}

// Items returns the items of the loop.
func (s *StepLoopConfig) Items() *hcl.Attribute {
	return s.items
}

// Label returns the label of the loop.
func (s *StepLoopConfig) Label() *hcl.Attribute {
	return s.label
}

// Condition returns the condition of the loop.
func (s *StepLoopConfig) Condition() *hcl.Attribute {
	return s.condition
}

type StepEscalationConfig struct {
	escalate        *hcl.Attribute
	impersonateUser *hcl.Attribute
}

// Escalate returns the escalation attribute for the step.
func (s *StepEscalationConfig) Escalate() *hcl.Attribute {
	return s.escalate
}

// ImpersonateUser returns the impersonate user attribute for the step.
func (s *StepEscalationConfig) ImpersonateUser() *hcl.Attribute {
	return s.impersonateUser
}

type StepOutputConfig struct {
	continueOnFail   *hcl.Attribute
	changedCondition *hcl.Attribute
	failedCondition  *hcl.Attribute
}

// ContinueOnFail returns the continue on fail attribute for the step.
func (s *StepOutputConfig) ContinueOnFail() *hcl.Attribute {
	return s.continueOnFail
}

// ChangedCondition returns the changed condition attribute for the step.
func (s *StepOutputConfig) ChangedCondition() *hcl.Attribute {
	return s.changedCondition
}

// FailedCondition returns the failed condition attribute for the step.
func (s *StepOutputConfig) FailedCondition() *hcl.Attribute {
	return s.failedCondition
}

type SingleStep struct {
	common     *StepCommonConfig
	escalation *StepEscalationConfig
	output     *StepOutputConfig

	module module.Module
}

// Common returns the common configuration for the step.
// This is used primarily for testing.
func (s *SingleStep) Common() *StepCommonConfig {
	return s.common
}

// Escalation returns the escalation configuration for the step.
// This is used primarily for testing.
func (s *SingleStep) Escalation() *StepEscalationConfig {
	return s.escalation
}

// Output returns the output configuration for the step.
// This is used primarily for testing.
func (s *SingleStep) Output() *StepOutputConfig {
	return s.output
}

// Module returns the module configuration for the step.
// This is used primarily for testing.
func (s *SingleStep) Module() module.Module {
	return s.module
}

// Run implements Step.
func (s *SingleStep) Run(ctx *workflowContext) {

	nameText := ui.Text(s.common.name).WithStyle(ui.StyleBold)
	name := ctx.ui.Format(nameText)
	line := ctx.ui.FormatLine('=', nil)

	message := fmt.Sprintf("\nSTEP - %s\n%s", name, line)
	ctx.ui.Print(message)

	ctx.LoadHostVars()

	wg := sync.WaitGroup{}

	for _, host := range s.common.targets {
		wg.Add(1)
		go func(h *inventory.Host) {
			defer wg.Done()
			if ctx.IsFailed(h) {
				return
			}
			err := s.runOnHost(HostWorkflowContext(ctx, h))
			if err != nil {
				ctx.MarkFailed(h)
			}
		}(host)
	}

	wg.Wait()
}

type stepIteration struct {
	label string
	index cty.Value
	item  cty.Value
}

// StepIterator handles any loop behavior for a step.
type StepIterator interface {
	Next() bool
	Value() *stepIteration
}

type singleIterator struct {
	completed bool
}

func (s *singleIterator) Next() bool {
	if s.completed {
		return false
	}
	s.completed = true
	return true
}

func (s *singleIterator) Value() *stepIteration {
	return nil // No iteration values to return
}

type multiIterator struct {
	index      int
	iterations []*stepIteration
}

func (m *multiIterator) Next() bool {

	if m.index >= len(m.iterations) {
		return false
	}

	m.index++
	return true
}

func (m *multiIterator) Value() *stepIteration {

	if m.index == 0 || m.index > len(m.iterations) {
		return nil
	}

	return m.iterations[m.index-1]
}
