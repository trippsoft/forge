// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

// Step abstracts a single Step or a procedure in a process.
type Step interface {
	ID() string                                            // ID returns the identifier of the step.
	Run(wc *WorkflowContext) (map[string]cty.Value, error) // Run executes the step using the provided WorkflowContext.
}

// StepCommonConfig holds common configuration attributes for steps and procedures.
type StepCommonConfig struct {
	id   string
	name string

	targets []*inventory.Host

	loop *StepLoopConfig

	condition *hcl.Attribute

	execTimeout *hcl.Attribute
	whatIf      *hcl.Attribute

	input map[string]*hcl.Attribute
}

// ID returns the identifier of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) ID() string {
	return s.id
}

// Name returns the name of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Name() string {
	return s.name
}

// Targets returns the slice of targets associated with the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Targets() []*inventory.Host {
	return s.targets
}

// Loop returns the loop configuration of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Loop() *StepLoopConfig {
	return s.loop
}

// Condition returns the condition attribute of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Condition() *hcl.Attribute {
	return s.condition
}

// ExecTimeout returns the execution timeout attribute of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) ExecTimeout() *hcl.Attribute {
	return s.execTimeout
}

// WhatIf returns the what-if attribute of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) WhatIf() *hcl.Attribute {
	return s.whatIf
}

// Input returns the input attributes of the step.
//
// This is used primarily for testing purposes.
func (s *StepCommonConfig) Input() map[string]*hcl.Attribute {
	return s.input
}

// Combine merges the values from another StepCommonConfig into this one, prioritizing existing values.
func (s *StepCommonConfig) Combine(other *StepCommonConfig) {
	if s == nil || other == nil {
		return
	}

	if s.targets == nil {
		s.targets = other.targets
	}

	if s.execTimeout == nil {
		s.execTimeout = other.execTimeout
	}

	if s.whatIf == nil {
		s.whatIf = other.whatIf
	}
}

// StepLoopConfig holds configuration for step looping.
//
// This represents the escalate block within a step or procedure block.
type StepLoopConfig struct {
	items     *hcl.Attribute
	label     *hcl.Attribute
	condition *hcl.Attribute
}

// Items returns the HCL attribute representing the items to loop over.
//
// This is used primarily for testing purposes.
func (s *StepLoopConfig) Items() *hcl.Attribute {
	return s.items
}

// Label returns the HCL attribute representing the label for each loop iteration.
//
// This is used primarily for testing purposes.
func (s *StepLoopConfig) Label() *hcl.Attribute {
	return s.label
}

// Condition returns the HCL attribute representing the condition for each loop iteration.
//
// This is used primarily for testing purposes.
func (s *StepLoopConfig) Condition() *hcl.Attribute {
	return s.condition
}

// StepEscalateConfig holds configuration for step escalation or impersonation.
//
// This represents the escalate block within a process, step, or procedure block.
type StepEscalateConfig struct {
	escalate        *hcl.Attribute
	impersonateUser *hcl.Attribute
}

// Escalate returns the HCL attribute representing whether escalation is enabled.
//
// This is used primarily for testing purposes.
func (s *StepEscalateConfig) Escalate() *hcl.Attribute {
	return s.escalate
}

// ImpersonateUser returns the HCL attribute representing the user to impersonate.
//
// This is used primarily for testing purposes.
func (s *StepEscalateConfig) ImpersonateUser() *hcl.Attribute {
	return s.impersonateUser
}

// Combine merges the values from another StepEscalateConfig into this one, prioritizing existing values.
func (s *StepEscalateConfig) Combine(escalate *StepEscalateConfig) {
	if s == nil || escalate == nil {
		return
	}

	if s.escalate == nil {
		s.escalate = escalate.escalate
	}

	if s.impersonateUser == nil {
		s.impersonateUser = escalate.impersonateUser
	}
}

// StepBuilder defines the interface for building Step instances.
type StepBuilder interface {
	// WithProcessCommon sets the common configuration inherited from the process.
	//
	// This will merge the provided common configuration with the step's own common configuration,
	// prioritizing the step's own values.
	WithProcessCommon(common *StepCommonConfig) StepBuilder
	// WithProcessEscalate sets the escalation configuration inherited from the process.
	//
	// This will merge the provided escalation configuration with the step's own escalation configuration,
	// prioritizing the step's own values.
	WithProcessEscalate(escalate *StepEscalateConfig) StepBuilder

	// AllTargets returns a clone of the slice of all targets associated with the step.
	AllTargets() []*inventory.Host

	// Build constructs and returns the Step instance.
	Build() (Step, error)
}

// StepOutputConfig holds configuration for step output handling.
//
// This represents the output block within a step block.
type StepOutputConfig struct {
	continueOnFail   *hcl.Attribute
	changedCondition *hcl.Attribute
	failedCondition  *hcl.Attribute
}

// ContinueOnFail returns the HCL attribute representing whether to continue on failure.
//
// This is used primarily for testing purposes.
func (s *StepOutputConfig) ContinueOnFail() *hcl.Attribute {
	return s.continueOnFail
}

// ChangedCondition returns the HCL attribute representing the condition for marking output as changed.
//
// This is used primarily for testing purposes.
func (s *StepOutputConfig) ChangedCondition() *hcl.Attribute {
	return s.changedCondition
}

// FailedCondition returns the HCL attribute representing the condition for marking output as failed.
//
// This is used primarily for testing purposes.
func (s *StepOutputConfig) FailedCondition() *hcl.Attribute {
	return s.failedCondition
}

// SingleStep represents a single step within a process.
type SingleStep struct {
	common   *StepCommonConfig
	escalate *StepEscalateConfig
	output   *StepOutputConfig

	module module.Module
}

// ID implements Step.
func (s *SingleStep) ID() string {
	return s.common.id
}

// Common returns the common configuration of the step.
//
// This is used primarily for testing purposes.
func (s *SingleStep) Common() *StepCommonConfig {
	return s.common
}

// Escalate returns the escalation configuration of the step.
//
// This is used primarily for testing purposes.
func (s *SingleStep) Escalate() *StepEscalateConfig {
	return s.escalate
}

// Output returns the output configuration of the step.
//
// This is used primarily for testing purposes.
func (s *SingleStep) Output() *StepOutputConfig {
	return s.output
}

// Module returns the module associated with the step.
//
// This is used primarily for testing purposes.
func (s *SingleStep) Module() module.Module {
	return s.module
}

// Run implements Step.
func (s *SingleStep) Run(wc *WorkflowContext) (map[string]cty.Value, error) {
	header := "STEP - " + s.common.name
	wc.ui.PrintHeader(ui.HeaderLevel2, header)

	wc.LoadHostVars()

	var err error
	mutex := sync.Mutex{}
	errChannel := make(chan error)
	outputs := make(map[string]cty.Value)

	for _, host := range s.common.targets {
		go func(h *inventory.Host) {
			if wc.IsFailed(h) {
				errChannel <- nil
				return
			}

			output, hostErr := s.runOnHost(NewHostWorkflowContext(wc, h))

			mutex.Lock()
			outputs[h.Name()] = output
			mutex.Unlock()

			errChannel <- hostErr
		}(host)
	}

	for range s.common.targets {
		hostErr := <-errChannel
		err = errors.Join(err, hostErr)
	}

	return outputs, err
}

func (s *SingleStep) runOnHost(hwc *HostWorkflowContext) (cty.Value, error) {
	err := hwc.LoadEvalContext()
	if err != nil {
		result := result.NewFailure(err, "failed to load evaluation context")
		output := s.handleHostResult(hwc, result)
		return output, err
	}

	condition := true
	if s.common.condition != nil {
		var diags hcl.Diagnostics
		condition, diags = util.ConvertHCLAttributeToBool(s.common.condition, hwc.evalContext)
		if diags.HasErrors() {
			result := result.NewFailure(diags, diags.Error())
			output := s.handleHostResult(hwc, result)
			return output, diags
		}
	}

	if !condition {
		result := result.NewSkipped()
		output := s.handleHostResult(hwc, result)
		return output, nil
	}

	iterator, err := s.getStepIterator(hwc)
	if err != nil {
		result := result.NewFailure(err, err.Error())
		output := s.handleHostResult(hwc, result)
		return output, err
	}

	results := make([]*StepIterationResult, 0)
	for {
		iteration, ok := iterator.Next()

		if iteration == nil {
			break
		}

		var result cty.Value
		result, err = s.runHostIteration(hwc, iteration)
		results = append(results, &StepIterationResult{
			iteration: iteration,
			result:    result,
		})

		if err != nil {
			hwc.MarkFailed(hwc.host)
			break
		}

		if !ok {
			break
		}
	}

	var output cty.Value
	switch iterator.Type() {
	case StepIteratorTypeSingle:
		output = results[0].result
	case StepIteratorTypeList:
		outputList := make([]cty.Value, len(results))
		for i, r := range results {
			outputList[i] = r.result
		}

		if len(outputList) == 0 {
			output = cty.EmptyTupleVal
		} else {
			output = cty.TupleVal(outputList)
		}

	case StepIteratorTypeMap:
		outputMap := make(map[string]cty.Value, len(results))
		for _, r := range results {
			outputMap[r.iteration.index.AsString()] = r.result
		}

		if len(outputMap) == 0 {
			output = cty.EmptyObjectVal
		} else {
			output = cty.ObjectVal(outputMap)
		}

	default:
		output = cty.EmptyObjectVal
	}

	hwc.host.StoreStepOutput(s.common.id, output)

	return output, err
}

func (s *SingleStep) formatResultOutput(result *result.Result) cty.Value {
	outputMap := map[string]cty.Value{
		"failed":  cty.BoolVal(result.Failed),
		"skipped": cty.BoolVal(result.Skipped),
		"changed": cty.BoolVal(result.Changed),
	}

	if result.Err != nil {
		outputMap["error"] = cty.StringVal(result.Err.Error())
		outputMap["error_detail"] = cty.StringVal(result.ErrDetail)
	}

	if result.Output != nil {
		outputMap["output"] = cty.ObjectVal(result.Output)
	}

	return cty.ObjectVal(outputMap)
}

func (s *SingleStep) handleHostResult(hwc *HostWorkflowContext, r *result.Result) cty.Value {
	if r == nil {
		r = result.NewFailure(errors.New("no result returned from module"), "")
	}

	hwc.ui.PrintHostResult(hwc.host.Name(), r)
	output := s.formatResultOutput(r)
	hwc.host.StoreStepOutput(s.common.id, output)

	return output
}

func (s *SingleStep) handleHostIterationResult(
	hwc *HostWorkflowContext,
	iteration *StepIteration,
	r *result.Result,
) cty.Value {

	if r == nil {
		r = result.NewFailure(errors.New("no result returned from module"), "")
	}

	hwc.ui.PrintIterationResult(hwc.host.Name(), iteration.label, r)
	return s.formatResultOutput(r)
}

func (s *SingleStep) getStepIterator(hwc *HostWorkflowContext) (StepIterator, error) {
	if s.common == nil || s.common.loop == nil || s.common.loop.items == nil {
		return StepIteratorSingle, nil
	}

	diags := hcl.Diagnostics{}

	itemsValue, moreDiags := s.common.loop.items.Expr.Value(hwc.evalContext)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	if !itemsValue.IsWhollyKnown() || itemsValue.IsNull() {
		return nil, fmt.Errorf("items must be a wholly known and non-null value, got %q", itemsValue.GoString())
	}

	itemsType := itemsValue.Type()

	var iteratorType StepIteratorType
	switch {
	case itemsType.IsListType() || itemsType.IsTupleType():
		iteratorType = StepIteratorTypeList
	case itemsType.IsMapType() || itemsType.IsObjectType():
		iteratorType = StepIteratorTypeMap
	default:
		return nil, fmt.Errorf("items must be a list, tuple, map, or object, got %q", itemsType.FriendlyName())
	}

	iterations := make([]*StepIteration, 0, itemsValue.LengthInt())
	it := itemsValue.ElementIterator()
	for it.Next() {
		key, value := it.Element()
		iterations = append(iterations, &StepIteration{
			index: key,
			value: value,
		})
	}

	return &MultipleStepIterator{
		iteratorType: iteratorType,
		items:        iterations,
	}, nil
}

func (s *SingleStep) runHostIteration(hwc *HostWorkflowContext, iteration *StepIteration) (cty.Value, error) {
	if iteration == nil {
		return cty.NilVal, errors.New("iteration cannot be nil")
	}

	if !iteration.index.IsNull() && !iteration.value.IsNull() {
		hwc.evalContext.Variables["item"] = iteration.value
		hwc.evalContext.Variables["index"] = iteration.index
		defer delete(hwc.evalContext.Variables, "item")
		defer delete(hwc.evalContext.Variables, "index")

		condition := true

		var diags hcl.Diagnostics

		if s.common != nil && s.common.loop != nil && s.common.loop.label != nil {
			iteration.label, diags = util.ConvertHCLAttributeToString(s.common.loop.label, hwc.evalContext)
			if diags.HasErrors() {
				output := result.NewFailure(diags, diags.Error())
				return s.handleHostIterationResult(hwc, iteration, output), diags
			}
		}

		if iteration.label == "" {
			switch {
			case !iteration.index.IsWhollyKnown() || iteration.index.IsNull():
				iteration.label = iteration.index.GoString()
			case iteration.index.Type().Equals(cty.Number):
				number, _ := iteration.index.AsBigFloat().Int64()
				iteration.label = fmt.Sprintf("%d", number)
			case iteration.index.Type().Equals(cty.String):
				iteration.label = iteration.index.AsString()
			default:
				iteration.label = iteration.index.GoString()
			}
		}

		if s.common != nil && s.common.loop != nil && s.common.loop.condition != nil {
			condition, diags = util.ConvertHCLAttributeToBool(s.common.loop.condition, hwc.evalContext)
			if diags.HasErrors() {
				output := result.NewFailure(diags, diags.Error())
				return s.handleHostIterationResult(hwc, iteration, output), diags
			}
		}

		if !condition {
			result := result.NewSkipped()
			return s.handleHostIterationResult(hwc, iteration, result), nil
		}
	}

	escalation, err := s.getEscalation(hwc)
	if err != nil {
		result := result.NewFailure(err, err.Error())
		return s.handleHostIterationResult(hwc, iteration, result), err
	}
	timeout := module.DefaultTimeout
	if s.common != nil && s.common.execTimeout != nil {
		var diags hcl.Diagnostics
		timeout, diags = util.ConvertHCLAttributeToDuration(s.common.execTimeout, hwc.evalContext)
		if diags.HasErrors() {
			result := result.NewFailure(diags, diags.Error())
			return s.handleHostIterationResult(hwc, iteration, result), diags
		}
	}

	whatIf := false
	if s.common != nil && s.common.whatIf != nil {
		var diags hcl.Diagnostics
		whatIf, diags = util.ConvertHCLAttributeToBool(s.common.whatIf, hwc.evalContext)
		if diags.HasErrors() {
			result := result.NewFailure(diags, diags.Error())
			return s.handleHostIterationResult(hwc, iteration, result), diags
		}
	}

	input := make(map[string]cty.Value, len(s.common.input))
	if s.common != nil && s.common.input != nil {
		for k, attr := range s.common.input {
			var diags hcl.Diagnostics
			input[k], diags = attr.Expr.Value(hwc.evalContext)
			if diags.HasErrors() {
				result := result.NewFailure(diags, diags.Error())
				return s.handleHostIterationResult(hwc, iteration, result), diags
			}
		}
	}

	input, err = s.module.InputSpec().Convert(input)
	if err != nil {
		result := result.NewFailure(err, err.Error())
		return s.handleHostIterationResult(hwc, iteration, result), err
	}

	err = s.module.InputSpec().Validate(input)
	if err != nil {
		result := result.NewFailure(err, err.Error())
		return s.handleHostIterationResult(hwc, iteration, result), err
	}

	config := &module.RunConfig{
		Transport:  hwc.host.Transport(),
		HostInfo:   hwc.host.Info(),
		Escalation: escalation,
		WhatIf:     whatIf,
		Input:      input,
	}

	err = s.module.Validate(config)
	if err != nil {
		result := result.NewFailure(err, err.Error())
		return s.handleHostIterationResult(hwc, iteration, result), err
	}

	runCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := s.module.Run(runCtx, config)
	if result == nil || s.output == nil {
		return s.handleHostIterationResult(hwc, iteration, result), nil
	}

	output := s.formatResultOutput(result)

	hwc.evalContext.Variables["result"] = output
	defer delete(hwc.evalContext.Variables, "result")

	if s.output.failedCondition != nil {
		var diags hcl.Diagnostics
		result.Failed, diags = util.ConvertHCLAttributeToBool(s.output.failedCondition, hwc.evalContext)
		if diags.HasErrors() {
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	if s.output.changedCondition != nil {
		var diags hcl.Diagnostics
		result.Changed, diags = util.ConvertHCLAttributeToBool(s.output.changedCondition, hwc.evalContext)
		if diags.HasErrors() {
			result.Changed = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	continueOnFail := false
	if s.output.continueOnFail != nil {
		var diags hcl.Diagnostics
		continueOnFail, diags = util.ConvertHCLAttributeToBool(s.output.continueOnFail, hwc.evalContext)
		if diags.HasErrors() {
			continueOnFail = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	s.handleHostIterationResult(hwc, iteration, result)

	if continueOnFail {
		return output, nil
	}

	return output, result.Err
}

func (s *SingleStep) getEscalation(hwc *HostWorkflowContext) (*transport.Escalation, error) {
	if s.escalate == nil || s.escalate.escalate == nil {
		return nil, nil // No escalation configured
	}

	escalate, diags := util.ConvertHCLAttributeToBool(s.escalate.escalate, hwc.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if !escalate {
		return nil, nil // No escalation needed
	}

	if s.escalate.impersonateUser == nil {
		return transport.NewEscalation(hwc.host.EscalateConfig().Pass()), nil
	}

	impersonateUser, diags := util.ConvertHCLAttributeToString(s.escalate.impersonateUser, hwc.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if impersonateUser == "" {
		return transport.NewEscalation(hwc.host.EscalateConfig().Pass()), nil
	}

	return transport.NewImpersonation(impersonateUser, hwc.host.EscalateConfig().Pass()), nil
}

// SingleStepBuilder is used to build a SingleStep instance during parsing.
type SingleStepBuilder struct {
	common   *StepCommonConfig
	escalate *StepEscalateConfig
	output   *StepOutputConfig

	module module.Module
}

// WithCommon sets the common configuration for the single step.
func (s *SingleStepBuilder) WithCommon(common *StepCommonConfig) *SingleStepBuilder {
	s.common = common
	return s
}

// WithEscalate sets the escalation configuration for the single step.
func (s *SingleStepBuilder) WithEscalate(escalate *StepEscalateConfig) *SingleStepBuilder {
	s.escalate = escalate
	return s
}

// WithOutput sets the output configuration for the single step.
func (s *SingleStepBuilder) WithOutput(output *StepOutputConfig) *SingleStepBuilder {
	s.output = output
	return s
}

// WithModule sets the module for the single step.
func (s *SingleStepBuilder) WithModule(module module.Module) *SingleStepBuilder {
	s.module = module
	return s
}

// WithProcessCommon implements StepBuilder.
func (s *SingleStepBuilder) WithProcessCommon(common *StepCommonConfig) StepBuilder {
	s.common.Combine(common)
	return s
}

// WithProcessEscalate implements StepBuilder.
func (s *SingleStepBuilder) WithProcessEscalate(escalate *StepEscalateConfig) StepBuilder {
	if s.escalate == nil {
		s.escalate = escalate
		return s
	}

	s.escalate.Combine(escalate)
	return s
}

// AllTargets implements StepBuilder.
func (s *SingleStepBuilder) AllTargets() []*inventory.Host {
	targets := slices.Clone(s.common.targets)
	return targets
}

// Build implements StepBuilder.
func (s *SingleStepBuilder) Build() (Step, error) {
	if s.common == nil {
		return nil, errors.New("Step failed to build: common configuration is missing")
	}

	if s.common.id == "" {
		return nil, errors.New("Step failed to build: id is missing")
	}

	if s.common.name == "" {
		return nil, errors.New("Step failed to build: name is missing")
	}

	if s.common.targets == nil {
		return nil, errors.New("Step failed to build: targets are missing")
	}

	if s.module == nil {
		return nil, errors.New("Step failed to build: module is missing")
	}

	return &SingleStep{
		common:   s.common,
		escalate: s.escalate,
		output:   s.output,
		module:   s.module,
	}, nil
}

// NewSingleStepBuilder creates a new instance of SingleStepBuilder.
func NewSingleStepBuilder() StepBuilder {
	return &SingleStepBuilder{}
}
