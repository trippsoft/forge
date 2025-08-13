// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type stepResultCode uint8

const (
	stepResultSkipped stepResultCode = iota
	stepResultFailure
	stepResultChanged
	stepResultNotChanged
)

var (
	stepResultText = map[stepResultCode]string{
		stepResultSkipped:    "SKIPPED",
		stepResultFailure:    "FAILED",
		stepResultChanged:    "CHANGED",
		stepResultNotChanged: "NOT CHANGED",
	}

	stepResultFormat = ui.TextFormatMap[stepResultCode]{
		stepResultSkipped:    ui.TextFormat().WithForegroundColor(ui.ForegroundCyan).WithStyle(ui.StyleItalic),
		stepResultFailure:    ui.TextFormat().WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold),
		stepResultChanged:    ui.TextFormat().WithForegroundColor(ui.ForegroundYellow).WithStyle(ui.StyleBold),
		stepResultNotChanged: ui.TextFormat().WithForegroundColor(ui.ForegroundGreen).WithStyle(ui.StyleBold),
	}

	stepErrorFormat = ui.TextFormat().WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold).WithStyle(ui.StyleItalic)
)

func (s *SingleStep) runOnHost(ctx *hostWorkflowContext) error {

	err := ctx.LoadEvalContext()
	if err != nil {
		_, _ = s.handleHostError(ctx, err) // ignore output for this error
		return err
	}

	condition := true // Default to true, in case a condition is not defined.
	if s.common.condition != nil {
		var diags hcl.Diagnostics
		condition, diags = util.ConvertHCLAttributeToBool(s.common.condition, ctx.evalContext)
		if diags.HasErrors() {
			output, _ := s.handleHostError(ctx, diags)
			ctx.host.StoreTask(s.common.id, output)
			return diags
		}
	}

	if !condition {
		hostMessage := s.formatHostResult(ctx, stepResultSkipped)
		ctx.ui.Print(hostMessage)
		output := cty.ObjectVal(map[string]cty.Value{
			"changed": cty.False,
			"failed":  cty.False,
			"skipped": cty.True,
		})
		ctx.host.StoreTask(s.common.id, output)
		return nil // Skipped
	}

	iterator, err := s.getStepIterator(ctx)
	if err != nil {
		output, _ := s.handleHostError(ctx, err)
		ctx.host.StoreTask(s.common.id, output)
		return err
	}

	resultsMap := map[cty.Value]cty.Value{}
	results := []cty.Value{}
	e := []error{}

	for iterator.Next() {
		iteration := iterator.Value()
		result, err := s.runHostIteration(ctx, iteration)
		if iteration != nil {
			resultsMap[iteration.item] = result
			results = append(results, result)
		} else {
			resultsMap[cty.DynamicVal] = result
			results = append(results, cty.DynamicVal)
		}
		e = append(e, err)
	}

	return errors.Join(e...)
}

func (s *SingleStep) runHostIteration(ctx *hostWorkflowContext, iteration *stepIteration) (cty.Value, error) {

	if iteration != nil {
		ctx.evalContext.Variables["item"] = iteration.item
		ctx.evalContext.Variables["index"] = iteration.index
		defer delete(ctx.evalContext.Variables, "item")
		defer delete(ctx.evalContext.Variables, "index")

		condition := true

		var diags hcl.Diagnostics

		if s.common != nil && s.common.loop != nil && s.common.loop.label != nil {
			iteration.label, diags = util.ConvertHCLAttributeToString(s.common.loop.label, ctx.evalContext)
			if diags.HasErrors() {
				return s.handleHostIterationError(ctx, iteration, diags, diags.Error())
			}
		}

		if iteration.label == "" {
			iteration.label = getIndexAsString(iteration.index)
		}

		if s.common != nil && s.common.loop != nil && s.common.loop.condition != nil {
			condition, diags = util.ConvertHCLAttributeToBool(s.common.loop.condition, ctx.evalContext)
		}

		if !condition {
			message := s.formatHostIterationResult(ctx, iteration, stepResultSkipped)
			ctx.ui.Print(message)
			output := cty.ObjectVal(map[string]cty.Value{
				"changed": cty.False,
				"failed":  cty.False,
				"skipped": cty.True,
			})
			return output, nil // Skipped iteration
		}
	}

	escalation, err := s.getEscalation(ctx)
	if err != nil {
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	timeout := module.DefaultTimeout

	if s.common != nil && s.common.execTimeout != nil {
		var diags hcl.Diagnostics
		timeout, diags = util.ConvertHCLAttributeToDuration(s.common.execTimeout, ctx.evalContext)
		if diags.HasErrors() {
			return s.handleHostIterationError(ctx, iteration, diags, diags.Error())
		}
	}

	commonConfig := &module.CommonConfig{
		Escalation: escalation,
		Timeout:    timeout,
	}

	input := make(map[string]cty.Value, len(s.common.input))
	if s.common != nil && s.common.input != nil {
		for k, attr := range s.common.input {
			var diags hcl.Diagnostics
			input[k], diags = attr.Expr.Value(ctx.evalContext)
			if diags.HasErrors() {
				return s.handleHostIterationError(ctx, iteration, diags, diags.Error())
			}
		}
	}

	input, err = s.module.InputSpec().Convert(input)
	if err != nil {
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	err = s.module.InputSpec().Validate(input)
	if err != nil {
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	err = s.module.Validate(ctx.host, input)
	if err != nil {
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	result := s.module.Run(ctx.host, commonConfig, input)
	if result == nil {
		err := errors.New("no result returned from module")
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	outputMap := map[string]cty.Value{}

	outputMap["failed"] = cty.BoolVal(result.Failed)
	outputMap["skipped"] = cty.BoolVal(result.Skipped)
	outputMap["changed"] = cty.BoolVal(result.Changed)

	if result.Err != nil {
		outputMap["error"] = cty.StringVal(result.Err.Error())
		outputMap["error_detail"] = cty.StringVal(result.ErrDetail)
	}

	if len(result.Output) > 0 {
		outputMap["output"] = cty.ObjectVal(result.Output)
	}

	ctx.evalContext.Variables["result"] = cty.ObjectVal(outputMap)
	defer delete(ctx.evalContext.Variables, "result")

	if s.output == nil {
		return s.handleHostIterationResult(ctx, iteration, result)
	}

	if s.output.failedCondition != nil {
		var diags hcl.Diagnostics
		result.Failed, diags = util.ConvertHCLAttributeToBool(s.output.failedCondition, ctx.evalContext)
		if diags.HasErrors() {
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	if s.output.changedCondition != nil {
		var diags hcl.Diagnostics
		result.Changed, diags = util.ConvertHCLAttributeToBool(s.output.changedCondition, ctx.evalContext)
		if diags.HasErrors() {
			result.Changed = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	continueOnFail := false
	if s.output.continueOnFail != nil {
		var diags hcl.Diagnostics
		continueOnFail, diags = util.ConvertHCLAttributeToBool(s.output.continueOnFail, ctx.evalContext)
		if diags.HasErrors() {
			continueOnFail = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	output, err := s.handleHostIterationResult(ctx, iteration, result)

	if err != nil && !continueOnFail {
		return output, err
	}

	return output, nil
}

func (s *SingleStep) getStepIterator(ctx *hostWorkflowContext) (StepIterator, error) {

	if s.common == nil || s.common.loop == nil || s.common.loop.items == nil {
		return &singleIterator{}, nil
	}

	diags := hcl.Diagnostics{}

	itemsValue, moreDiags := s.common.loop.items.Expr.Value(ctx.evalContext)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	if !itemsValue.IsWhollyKnown() || itemsValue.IsNull() {
		return nil, fmt.Errorf("items must be a wholly known and non-null value, got %q", itemsValue.GoString())
	}

	itemsType := itemsValue.Type()

	if !itemsType.IsListType() && !itemsType.IsTupleType() && !itemsType.IsMapType() && !itemsType.IsObjectType() {
		return nil, fmt.Errorf("items must be a list, tuple, map, or object, got %q", itemsType.FriendlyName())
	}

	iterations := make([]*stepIteration, 0, itemsValue.LengthInt())

	it := itemsValue.ElementIterator()
	for it.Next() {
		key, value := it.Element()
		iterations = append(iterations, &stepIteration{
			index: key,
			item:  value,
		})
	}

	return &multiIterator{iterations: iterations}, nil
}

func (s *SingleStep) getEscalation(ctx *hostWorkflowContext) (transport.Escalation, error) {

	if s.escalation == nil || s.escalation.escalate == nil {
		return nil, nil // No escalation configured
	}

	escalate, diags := util.ConvertHCLAttributeToBool(s.escalation.escalate, ctx.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if !escalate {
		return nil, nil // No escalation needed
	}

	if s.escalation.impersonateUser == nil {
		return transport.NewEscalation(ctx.host.EscalateConfig().Pass()), nil
	}

	impersonateUser, diags := util.ConvertHCLAttributeToString(s.escalation.impersonateUser, ctx.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if impersonateUser == "" {
		return transport.NewEscalation(ctx.host.EscalateConfig().Pass()), nil
	}

	return transport.NewImpersonation(impersonateUser, ctx.host.EscalateConfig().Pass()), nil
}

func (s *SingleStep) handleHostIterationResult(ctx *hostWorkflowContext, iteration *stepIteration, result *module.Result) (cty.Value, error) {

	if result == nil {
		err := errors.New("no result returned from module")
		return s.handleHostIterationError(ctx, iteration, err, err.Error())
	}

	if result.Err != nil {
		return s.handleHostIterationError(ctx, iteration, result.Err, result.ErrDetail)
	}

	var hostMessage string
	if result.Changed {
		hostMessage = s.formatHostIterationResult(ctx, iteration, stepResultChanged)
	} else {
		hostMessage = s.formatHostIterationResult(ctx, iteration, stepResultNotChanged)
	}

	ctx.ui.Print(hostMessage)

	outputMap := map[string]cty.Value{}

	outputMap["failed"] = cty.BoolVal(result.Failed)
	outputMap["skipped"] = cty.BoolVal(result.Skipped)
	outputMap["changed"] = cty.BoolVal(result.Changed)

	if result.Err != nil {
		outputMap["error"] = cty.StringVal(result.Err.Error())
		outputMap["error_detail"] = cty.StringVal(result.ErrDetail)
	}

	if len(result.Output) > 0 {
		outputMap["output"] = cty.ObjectVal(result.Output)
	}

	return cty.ObjectVal(outputMap), nil
}

func (s *SingleStep) handleHostError(ctx *hostWorkflowContext, err error) (cty.Value, error) {
	hostMessage := s.formatHostResult(ctx, stepResultFailure)
	errMessage := s.formatHostError(ctx, err)
	outMessage := fmt.Sprintf("%s\n%s\n", hostMessage, errMessage)
	ctx.ui.Print(outMessage)
	ctx.ui.Error(errMessage)

	return cty.ObjectVal(map[string]cty.Value{
		"changed":      cty.BoolVal(false),
		"failed":       cty.BoolVal(true),
		"skipped":      cty.BoolVal(false),
		"error":        cty.StringVal(err.Error()),
		"error_detail": cty.StringVal(err.Error()),
	}), err
}

func (s *SingleStep) handleHostIterationError(ctx *hostWorkflowContext, iteration *stepIteration, err error, detail string) (cty.Value, error) {

	hostMessage := s.formatHostIterationResult(ctx, iteration, stepResultFailure)
	errMessage := s.formatHostIterationError(ctx, iteration, err)
	outMessage := fmt.Sprintf("%s\n%s\n", hostMessage, errMessage)
	ctx.ui.Print(outMessage)
	ctx.ui.Error(errMessage)

	return cty.ObjectVal(map[string]cty.Value{
		"changed":      cty.BoolVal(false),
		"failed":       cty.BoolVal(true),
		"skipped":      cty.BoolVal(false),
		"error":        cty.StringVal(err.Error()),
		"error_detail": cty.StringVal(detail),
	}), err
}

func (s *SingleStep) formatHostError(ctx *hostWorkflowContext, err error) string {
	message := fmt.Sprintf("host %s: step %s: ERROR: %v\n", ctx.host.Name(), s.common.id, err)
	messageText := ui.Text(message).WithFormat(stepErrorFormat).WithLeftMargin(4)
	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostIterationError(ctx *hostWorkflowContext, iteration *stepIteration, err error) string {
	if iteration == nil {
		return s.formatHostError(ctx, err)
	}

	label := iteration.label

	if label == "" {
		label = getIndexAsString(iteration.index)
	}

	message := fmt.Sprintf("host %s: step %s: item %s: ERROR: %v\n", ctx.host.Name(), s.common.id, label, err)
	messageText := ui.Text(message).WithFormat(stepErrorFormat).WithLeftMargin(4)
	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostResult(ctx *hostWorkflowContext, result stepResultCode) string {

	hostMessage := fmt.Sprintf("%s:", ctx.host.Name())
	hostText := ui.Text(hostMessage).WithFormat(stepResultFormat[result]).WithLeftMargin(2)

	statusText := ui.Text(stepResultText[result]).WithFormat(stepResultFormat[result]).WithRightMargin(2)

	return ctx.ui.FormatColumns(hostText, statusText)
}

func (s *SingleStep) formatHostIterationResult(ctx *hostWorkflowContext, iteration *stepIteration, result stepResultCode) string {

	if iteration == nil {
		return s.formatHostResult(ctx, result)
	}

	label := iteration.label

	if label == "" {
		label = getIndexAsString(iteration.index)
	}

	iterationMessage := ctx.ui.Format(ui.Text(label).WithStyle(ui.StyleItalic))
	hostMessage := fmt.Sprintf("%s->%s:", ctx.host.Name(), iterationMessage)
	hostText := ui.Text(hostMessage).WithFormat(stepResultFormat[result]).WithLeftMargin(2)

	statusText := ui.Text(stepResultText[result]).WithFormat(stepResultFormat[result]).WithRightMargin(2)

	return ctx.ui.FormatColumns(hostText, statusText)
}

func getIndexAsString(value cty.Value) string {

	switch {
	case !value.IsWhollyKnown() || value.IsNull():
		return value.GoString()
	case value.Type().Equals(cty.Number):
		num, _ := value.AsBigFloat().Int64()
		return fmt.Sprintf("%d", num)
	case value.Type().Equals(cty.String):
		return value.AsString()
	}

	return value.GoString()
}
