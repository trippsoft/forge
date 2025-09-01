// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/hclutil"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/ui"
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

	stepErrorFormat       = ui.TextFormat().WithForegroundColor(ui.ForegroundRed).WithStyle(ui.StyleBold).WithStyle(ui.StyleItalic)
	stepErrorDetailFormat = ui.TextFormat().WithForegroundColor(ui.ForegroundBlack).WithStyle(ui.StyleItalic)
	stepWarningFormat     = ui.TextFormat().WithForegroundColor(ui.ForegroundYellow).WithStyle(ui.StyleBold).WithStyle(ui.StyleItalic)
	stepMessageFormat     = ui.TextFormat().WithForegroundColor(ui.ForegroundGreen).WithStyle(ui.StyleBold).WithStyle(ui.StyleItalic)
)

type iterationResult struct {
	iteration *stepIteration
	result    cty.Value
}

func (s *SingleStep) runOnHost(ctx *hostWorkflowContext) (cty.Value, error) {
	err := ctx.LoadEvalContext()
	if err != nil {
		result := module.NewFailure(err, "failed to load evaluation context")
		output := s.handleHostResult(ctx, result)

		return output, err
	}

	condition := true // Default to true, in case a condition is not defined.
	if s.common.condition != nil {
		var diags hcl.Diagnostics
		condition, diags = hclutil.ConvertHCLAttributeToBool(s.common.condition, ctx.evalContext)
		if diags.HasErrors() {
			result := module.NewFailure(diags, diags.Error())
			output := s.handleHostResult(ctx, result)

			return output, diags
		}
	}

	if !condition {
		result := module.NewSkipped()
		output := s.handleHostResult(ctx, result)

		return output, nil // Skipped
	}

	iterator, err := s.getStepIterator(ctx)
	if err != nil {
		result := module.NewFailure(err, err.Error())
		output := s.handleHostResult(ctx, result)

		return output, err
	}

	results := []*iterationResult{}
	for iterator.Next() {
		iteration := iterator.Value()
		var result cty.Value
		result, err = s.runHostIteration(ctx, iteration)
		results = append(results, &iterationResult{iteration: iteration, result: result})
		if err != nil {
			ctx.MarkFailed(ctx.host)
			break
		}
	}

	var output cty.Value
	switch iterator.Type() {
	case stepIteratorSingle:
		output = results[0].result
	case stepIteratorMap:
		outputMap := make(map[string]cty.Value, len(results))
		for _, r := range results {
			outputMap[r.iteration.index.AsString()] = r.result
		}

		if len(outputMap) == 0 {
			output = cty.EmptyObjectVal
		} else {
			output = cty.ObjectVal(outputMap)
		}

	case stepIteratorList:
		outputList := make([]cty.Value, len(results))
		for i, r := range results {
			outputList[i] = r.result
		}

		if len(outputList) == 0 {
			output = cty.EmptyTupleVal
		} else {
			output = cty.TupleVal(outputList)
		}

	default:
		output = cty.EmptyObjectVal
	}

	ctx.host.StoreStepOutput(s.common.id, output)

	return output, err
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
			iteration.label, diags = hclutil.ConvertHCLAttributeToString(s.common.loop.label, ctx.evalContext)
			if diags.HasErrors() {
				output := module.NewFailure(diags, diags.Error())

				return s.handleHostIterationResult(ctx, iteration, output)
			}
		}

		if iteration.label == "" {
			iteration.label = getIndexAsString(iteration.index)
		}

		if s.common != nil && s.common.loop != nil && s.common.loop.condition != nil {
			condition, diags = hclutil.ConvertHCLAttributeToBool(s.common.loop.condition, ctx.evalContext)
		}

		if !condition {
			result := module.NewSkipped()

			return s.handleHostIterationResult(ctx, iteration, result)
		}
	}

	escalation, err := s.getEscalation(ctx)
	if err != nil {
		result := module.NewFailure(err, err.Error())

		return s.handleHostIterationResult(ctx, iteration, result)
	}

	timeout := module.DefaultTimeout
	if s.common != nil && s.common.execTimeout != nil {
		var diags hcl.Diagnostics
		timeout, diags = hclutil.ConvertHCLAttributeToDuration(s.common.execTimeout, ctx.evalContext)
		if diags.HasErrors() {
			result := module.NewFailure(diags, diags.Error())

			return s.handleHostIterationResult(ctx, iteration, result)
		}
	}

	whatIf := false
	if s.common != nil && s.common.whatIf != nil {
		var diags hcl.Diagnostics
		whatIf, diags = hclutil.ConvertHCLAttributeToBool(s.common.whatIf, ctx.evalContext)
		if diags.HasErrors() {
			result := module.NewFailure(diags, diags.Error())

			return s.handleHostIterationResult(ctx, iteration, result)
		}
	}

	input := make(map[string]cty.Value, len(s.common.input))
	if s.common != nil && s.common.input != nil {
		for k, attr := range s.common.input {
			var diags hcl.Diagnostics
			input[k], diags = attr.Expr.Value(ctx.evalContext)
			if diags.HasErrors() {
				result := module.NewFailure(diags, diags.Error())

				return s.handleHostIterationResult(ctx, iteration, result)
			}
		}
	}

	input, err = s.module.InputSpec().Convert(input)
	if err != nil {
		result := module.NewFailure(err, err.Error())

		return s.handleHostIterationResult(ctx, iteration, result)
	}

	err = s.module.InputSpec().Validate(input)
	if err != nil {
		result := module.NewFailure(err, err.Error())

		return s.handleHostIterationResult(ctx, iteration, result)
	}

	config := &module.RunConfig{
		Transport:  ctx.host.Transport(),
		HostInfo:   ctx.host.Info(),
		Escalation: escalation,
		WhatIf:     whatIf,
		Input:      input,
	}

	err = s.module.Validate(config)
	if err != nil {
		result := module.NewFailure(err, err.Error())

		return s.handleHostIterationResult(ctx, iteration, result)
	}

	runCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := s.module.Run(runCtx, config)
	if result == nil || s.output == nil {
		return s.handleHostIterationResult(ctx, iteration, result)
	}

	output := formatResultOutput(result)

	ctx.evalContext.Variables["result"] = output
	defer delete(ctx.evalContext.Variables, "result")

	if s.output.failedCondition != nil {
		var diags hcl.Diagnostics
		result.Failed, diags = hclutil.ConvertHCLAttributeToBool(s.output.failedCondition, ctx.evalContext)
		if diags.HasErrors() {
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	if s.output.changedCondition != nil {
		var diags hcl.Diagnostics
		result.Changed, diags = hclutil.ConvertHCLAttributeToBool(s.output.changedCondition, ctx.evalContext)
		if diags.HasErrors() {
			result.Changed = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	continueOnFail := false
	if s.output.continueOnFail != nil {
		var diags hcl.Diagnostics
		continueOnFail, diags = hclutil.ConvertHCLAttributeToBool(s.output.continueOnFail, ctx.evalContext)
		if diags.HasErrors() {
			continueOnFail = false
			result.Failed = true
			result.Err = errors.Join(result.Err, diags)
		}
	}

	output, err = s.handleHostIterationResult(ctx, iteration, result)
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

	var iteratorType stepIteratorType
	switch {
	case itemsType.IsListType() || itemsType.IsTupleType():
		iteratorType = stepIteratorList
	case itemsType.IsMapType() || itemsType.IsObjectType():
		iteratorType = stepIteratorMap
	default:
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

	return &multiIterator{iteratorType: iteratorType, iterations: iterations}, nil
}

func (s *SingleStep) getEscalation(ctx *hostWorkflowContext) (transport.Escalation, error) {
	if s.escalation == nil || s.escalation.escalate == nil {
		return nil, nil // No escalation configured
	}

	escalate, diags := hclutil.ConvertHCLAttributeToBool(s.escalation.escalate, ctx.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if !escalate {
		return nil, nil // No escalation needed
	}

	if s.escalation.impersonateUser == nil {
		return transport.NewEscalation(ctx.host.EscalateConfig().Pass()), nil
	}

	impersonateUser, diags := hclutil.ConvertHCLAttributeToString(s.escalation.impersonateUser, ctx.evalContext)
	if diags.HasErrors() {
		return nil, diags
	}

	if impersonateUser == "" {
		return transport.NewEscalation(ctx.host.EscalateConfig().Pass()), nil
	}

	return transport.NewImpersonation(impersonateUser, ctx.host.EscalateConfig().Pass()), nil
}

func (s *SingleStep) handleHostResult(ctx *hostWorkflowContext, result *module.Result) cty.Value {
	if result == nil {
		result = module.NewFailure(errors.New("no result returned from module"), "")
	}

	errMessage := s.formatHostError(ctx, result.Err)
	errMessage += s.formatHostErrorDetail(ctx, result.ErrDetail)
	outMessage := s.formatHostResult(ctx, result)

	ctx.ui.Print(outMessage)
	if errMessage != "" {
		ctx.ui.Error(errMessage)
	}

	output := formatResultOutput(result)
	ctx.host.StoreStepOutput(s.common.id, output)

	return output
}

func (s *SingleStep) handleHostIterationResult(ctx *hostWorkflowContext, iteration *stepIteration, result *module.Result) (cty.Value, error) {
	if result == nil {
		result = module.NewFailure(errors.New("no result returned from module"), "no result returned from module")
	}

	errMessage := s.formatHostError(ctx, result.Err)
	errMessage += s.formatHostErrorDetail(ctx, result.ErrDetail)
	outMessage := s.formatHostIterationResult(ctx, iteration, result)

	ctx.ui.Print(outMessage)
	if errMessage != "" {
		ctx.ui.Error(errMessage)
	}

	return formatResultOutput(result), result.Err
}

func (s *SingleStep) formatHostError(ctx *hostWorkflowContext, err error) string {
	if err == nil {
		return ""
	}

	errMessage := strings.ReplaceAll(err.Error(), "\n", "\n     ")

	message := fmt.Sprintf("ERROR:   %s\n", errMessage)
	messageText := ui.Text(message).WithFormat(stepErrorFormat).WithLeftMargin(4)

	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostErrorDetail(ctx *hostWorkflowContext, errDetail string) string {
	if errDetail == "" || !ctx.debug {
		return ""
	}

	errDetail = strings.ReplaceAll(errDetail, "\n", "\n       ")

	message := fmt.Sprintf("DETAIL:  %s\n", errDetail)
	messageText := ui.Text(message).WithFormat(stepErrorDetailFormat).WithLeftMargin(6)

	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostWarning(ctx *hostWorkflowContext, warning string) string {
	if warning == "" {
		return ""
	}

	warning = strings.ReplaceAll(warning, "\n", "\n     ")

	message := fmt.Sprintf("WARNING: %s\n", warning)
	messageText := ui.Text(message).WithFormat(stepWarningFormat).WithLeftMargin(4)

	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostMessage(ctx *hostWorkflowContext, message string) string {
	if message == "" {
		return ""
	}

	message = strings.ReplaceAll(message, "\n", "\n     ")

	message = fmt.Sprintf("MESSAGE: %s\n", message)
	messageText := ui.Text(message).WithFormat(stepMessageFormat).WithLeftMargin(4)

	return ctx.ui.Format(messageText)
}

func (s *SingleStep) formatHostResult(ctx *hostWorkflowContext, result *module.Result) string {
	hostName := ctx.host.Name()
	hostMessage := fmt.Sprintf("%s:", hostName)
	runeCount := utf8.RuneCountInString(hostName)
	hostText := ui.Text(hostMessage).WithLeftMargin(2).WithRightMargin(65 - runeCount)

	resultCode := getResultCode(result)

	statusMessage := stepResultText[resultCode]
	runeCount = utf8.RuneCountInString(statusMessage)
	statusText := ui.Text(statusMessage).WithFormat(stepResultFormat[resultCode]).WithLeftMargin(12 - runeCount)

	warning := s.formatHostWarning(ctx, result.Warning)
	message := s.formatHostMessage(ctx, result.Message)

	return fmt.Sprintf("%s%s\n%s%s", ctx.ui.Format(hostText), ctx.ui.Format(statusText), warning, message)
}

func (s *SingleStep) formatHostIterationResult(ctx *hostWorkflowContext, iteration *stepIteration, result *module.Result) string {
	if iteration == nil {
		return s.formatHostResult(ctx, result)
	}

	label := iteration.label

	if label == "" {
		label = getIndexAsString(iteration.index)
	}

	hostName := ctx.host.Name()
	hostRuneCount := utf8.RuneCountInString(hostName)
	labelRuneCount := utf8.RuneCountInString(label)
	iterationMessage := ctx.ui.Format(ui.Text(label).WithStyle(ui.StyleItalic))
	hostMessage := fmt.Sprintf("%s->%s:", hostName, iterationMessage)
	hostText := ui.Text(hostMessage).WithLeftMargin(2).WithRightMargin(63 - hostRuneCount - labelRuneCount)

	resultCode := getResultCode(result)

	statusMessage := stepResultText[resultCode]
	statusRuneCount := utf8.RuneCountInString(statusMessage)
	statusText := ui.Text(statusMessage).WithFormat(stepResultFormat[resultCode]).WithLeftMargin(12 - statusRuneCount)

	warning := s.formatHostWarning(ctx, result.Warning)
	message := s.formatHostMessage(ctx, result.Message)

	return fmt.Sprintf("%s%s\n%s%s", ctx.ui.Format(hostText), ctx.ui.Format(statusText), warning, message)
}

func formatResultOutput(result *module.Result) cty.Value {
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

func getResultCode(result *module.Result) stepResultCode {
	if result.Skipped {
		return stepResultSkipped
	}

	if result.Failed {
		return stepResultFailure
	}

	if result.Changed {
		return stepResultChanged
	}

	return stepResultNotChanged
}
