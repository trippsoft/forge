// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type intermediateProcess struct {
	common     *StepCommonConfig
	escalate   *StepEscalationConfig
	gatherInfo bool
}

type parser struct {
	inventory      *inventory.Inventory
	parser         *hclparse.Parser
	moduleRegistry *module.Registry
}

func NewParser(inventory *inventory.Inventory, moduleRegistry *module.Registry) *parser {
	return &parser{
		inventory:      inventory,
		parser:         hclparse.NewParser(),
		moduleRegistry: moduleRegistry,
	}
}

func (p *parser) ParseWorkflowFile(path string, content []byte) (*Workflow, hcl.Diagnostics) {

	file, diags := p.parser.ParseHCL(content, path)
	if diags.HasErrors() {
		return nil, diags
	}

	bodyContent, moreDiags := file.Body.Content(workflowFileSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a workflow file")
	diags = diags.Extend(moreDiags)

	if diags.HasErrors() {
		return nil, diags
	}

	processes, moreDiags := p.parseProcessBlocks(bodyContent)
	diags = diags.Extend(moreDiags)

	if diags.HasErrors() {
		return nil, diags
	}

	return &Workflow{processes: processes}, diags
}

func (p *parser) parseProcessBlocks(content *hcl.BodyContent) ([]*Process, hcl.Diagnostics) {

	processes := make([]*Process, 0, len(content.Blocks))
	diags := hcl.Diagnostics{}

	for _, block := range content.Blocks {

		process, moreDiags := p.parseProcessBlock(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue
		}

		processes = append(processes, process)
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return processes, diags
}

func (p *parser) parseProcessBlock(block *hcl.Block) (*Process, hcl.Diagnostics) {

	if block.Type != "process" {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block type",
			Detail:   "Expected 'process' block type.",
			Subject:  &block.TypeRange,
		}}
	}

	if len(block.Labels) != 0 {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block labels",
			Detail:   "Expected no labels for 'process' block.",
			Subject:  &block.TypeRange,
		}}
	}

	diags := hcl.Diagnostics{}

	content, moreDiags := block.Body.Content(processBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a process block")
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	}

	intermediate := &intermediateProcess{gatherInfo: true}

	foundEscalate := false
	for _, block := range content.Blocks {

		if block.Type != "escalate" {
			continue
		}

		if foundEscalate {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate escalate block",
				Detail:   "Only one escalate block is allowed.",
				Subject:  &block.TypeRange,
			})
			continue
		}

		foundEscalate = true

		escalate, moreDiags := parseEscalateBlock(block)
		diags = diags.Extend(moreDiags)

		if moreDiags.HasErrors() {
			continue
		}

		intermediate.escalate = escalate
	}

	common, moreDiags := p.parseCommonElements(content)
	diags = diags.Extend(moreDiags)

	if !moreDiags.HasErrors() {
		intermediate.common = common
	}

	steps := []Step{}

	for _, block := range content.Blocks {

		switch block.Type {
		case "step":

			step, moreDiags := p.parseStepBlock(block, intermediate)
			diags = diags.Extend(moreDiags)
			if moreDiags.HasErrors() {
				continue
			}
			steps = append(steps, step)

		case "procedure":

			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Procedure blocks are not yet supported",
				Detail:   "Procedure blocks are currently not implemented. Skipping this block.",
				Subject:  &block.TypeRange,
			})
		}
	}

	for name, attr := range content.Attributes {
		switch name {
		case "gather_info":
			gatherInfo, moreDiags := util.ConvertHCLAttributeToBool(attr, nil)
			diags = diags.Extend(moreDiags)
			if !moreDiags.HasErrors() {
				intermediate.gatherInfo = gatherInfo
			}
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	if common == nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing common attributes",
			Detail:   "The 'process' block is missing required common attributes. This is likely a parser error.",
			Subject:  &block.DefRange,
		})
		return nil, diags
	}

	process := NewProcess(common.name, intermediate.gatherInfo, steps...)

	return process, diags
}

func (p *parser) parseStepBlock(block *hcl.Block, intermediate *intermediateProcess) (Step, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	if block == nil {
		return nil, diags
	}

	if block.Type != "step" {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block type",
			Detail:   "Expected 'step' block type.",
			Subject:  &block.TypeRange,
		}}
	}

	if len(block.Labels) != 1 {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block labels",
			Detail:   "Expected exactly one label for 'step' block.",
			Subject:  &block.TypeRange,
		}}
	}

	content, moreDiags := block.Body.Content(stepBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a step block")
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	}

	step := &SingleStep{}

	common, moreDiags := p.parseCommonElements(content)
	diags = diags.Extend(moreDiags)

	if common != nil {
		if intermediate != nil && intermediate.common != nil {
			if common.targets == nil {
				common.targets = intermediate.common.targets
			}

			if common.execTimeout == nil {
				common.execTimeout = intermediate.common.execTimeout
			}

			if common.whatIf == nil {
				common.whatIf = intermediate.common.whatIf
			}
		}

		common.id = block.Labels[0]
	} else if intermediate != nil && intermediate.common != nil {
		common = &StepCommonConfig{
			id:          block.Labels[0],
			name:        intermediate.common.name,
			targets:     intermediate.common.targets,
			execTimeout: intermediate.common.execTimeout,
			input:       intermediate.common.input,
			whatIf:      intermediate.common.whatIf,
		}
	}

	if !moreDiags.HasErrors() {
		step.common = common
	}

	foundEscalate := false
	foundOutput := false

	for _, block := range content.Blocks {

		switch block.Type {
		case "escalate":

			if foundEscalate {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate escalate block",
					Detail:   "Only one escalate block is allowed.",
					Subject:  &block.TypeRange,
				})
				continue
			}

			foundEscalate = true

			escalate, moreDiags := parseEscalateBlock(block)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			step.escalation = escalate

		case "output":

			if foundOutput {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate output block",
					Detail:   "Only one output block is allowed.",
					Subject:  &block.TypeRange,
				})
				continue
			}

			foundOutput = true

			output, moreDiags := parseOutputBlock(block)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			step.output = output
		}
	}

	if step.escalation != nil {
		if intermediate != nil && intermediate.escalate != nil {
			if step.escalation.escalate == nil {
				step.escalation.escalate = intermediate.escalate.escalate
			}

			if step.escalation.impersonateUser == nil {
				step.escalation.impersonateUser = intermediate.escalate.impersonateUser
			}
		}
	} else if intermediate != nil && intermediate.escalate != nil {
		step.escalation = &StepEscalationConfig{
			escalate:        intermediate.escalate.escalate,
			impersonateUser: intermediate.escalate.impersonateUser,
		}
	}

	for name, attr := range content.Attributes {

		if name != "module" {
			continue
		}

		moduleName, moreDiags := util.ConvertHCLAttributeToString(attr, nil)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue
		}

		var exists bool
		step.module, exists = p.moduleRegistry.Lookup(moduleName)
		if !exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Module not found",
				Detail:   fmt.Sprintf("Module %q not found", moduleName),
				Subject:  attr.NameRange.Ptr(),
			})
		}
	}

	if step.module == nil && !diags.HasErrors() {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Step missing module",
			Detail:   "Step is missing required 'module' attribute. This is likely a parser error.",
			Subject:  &block.DefRange,
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return step, diags
}

func (p *parser) parseCommonElements(content *hcl.BodyContent) (*StepCommonConfig, hcl.Diagnostics) {

	config := &StepCommonConfig{}
	diags := hcl.Diagnostics{}

	foundLoop := false
	foundInput := false

	for _, block := range content.Blocks {
		switch block.Type {
		case "loop":
			if foundLoop {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate loop block",
					Detail:   "Only one loop block is allowed.",
					Subject:  &block.TypeRange,
				})
				continue
			}

			foundLoop = true

			loop, moreDiags := parseLoopBlock(block)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			config.loop = loop

		case "input":
			if foundInput {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate input block",
					Detail:   "Only one input block is allowed.",
					Subject:  &block.TypeRange,
				})
				continue
			}

			foundInput = true

			input, moreDiags := parseInputBlock(block)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			config.input = input
		}
	}

	for name, attr := range content.Attributes {
		switch name {
		case "name":

			name, moreDiags := util.ConvertHCLAttributeToString(attr, nil)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			config.name = name

		case "targets":

			targets, moreDiags := p.parseTargetsAttribute(attr)
			diags = diags.Extend(moreDiags)

			if moreDiags.HasErrors() {
				continue
			}

			config.targets = targets

		case "condition":
			config.condition = attr
		case "exec_timeout":
			config.execTimeout = attr
		case "what_if":
			config.whatIf = attr
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return config, diags
}

func (p *parser) parseTargetsAttribute(attr *hcl.Attribute) ([]*inventory.Host, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	targetsValue, moreDiags := attr.Expr.Value(nil)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	switch {
	case targetsValue.IsNull():
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid targets",
			Detail:   "The 'targets' attribute cannot be null.",
			Subject:  attr.Expr.Range().Ptr(),
		})
		return nil, diags
	case !targetsValue.IsWhollyKnown():
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unknown targets",
			Detail:   "The 'targets' attribute must be a known value.",
			Subject:  attr.Expr.Range().Ptr(),
		})
		return nil, diags
	case targetsValue.Type().Equals(cty.String):
		target, exists := p.inventory.Host(targetsValue.AsString())
		if exists {
			return []*inventory.Host{target}, diags
		}

		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Target not found",
			Detail:   fmt.Sprintf("The target %q does not exist in the inventory", targetsValue.AsString()),
			Subject:  attr.Expr.Range().Ptr(),
		})

		return nil, diags

	case targetsValue.Type().IsListType() || targetsValue.Type().IsSetType() || targetsValue.Type().IsTupleType():

		it := targetsValue.ElementIterator()
		seenTargets := util.NewSet[string]()
		targetHosts := util.NewSet[*inventory.Host]()

		for it.Next() {
			_, elem := it.Element()

			if elem.IsNull() {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid target",
					Detail:   "The 'targets' attribute cannot contain null values.",
					Subject:  attr.Expr.Range().Ptr(),
				})
				continue
			}

			if !elem.Type().Equals(cty.String) {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid target type",
					Detail:   "The 'targets' attribute must be a list of strings.",
					Subject:  attr.Expr.Range().Ptr(),
				})
				continue
			}

			targetName := elem.AsString()

			if seenTargets.Contains(targetName) {
				continue
			}

			seenTargets.Add(targetName)

			target, exists := p.inventory.Host(elem.AsString())
			if exists {
				targetHosts.Add(target)
			} else {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Target not found",
					Detail:   fmt.Sprintf("The target %q does not exist in the inventory", elem.AsString()),
					Subject:  attr.Expr.Range().Ptr(),
				})
			}
		}

		if diags.HasErrors() {
			return nil, diags
		}

		return targetHosts.Items(), diags

	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid targets type",
			Detail:   "The 'targets' attribute must be a string or a list of strings.",
			Subject:  attr.Expr.Range().Ptr(),
		})
		return nil, diags
	}
}

func parseEscalateBlock(block *hcl.Block) (*StepEscalationConfig, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	if block == nil {
		return nil, diags
	}

	if block.Type != "escalate" {
		return nil, diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block type",
			Detail:   "Expected 'escalate' block.",
			Subject:  &block.TypeRange,
		})
	}

	if len(block.Labels) != 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block labels",
			Detail:   "Expected no labels for 'escalate' block.",
			Subject:  &block.TypeRange,
		})
	}

	content, moreDiags := block.Body.Content(escalateBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in an escalate block")
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	}

	config := &StepEscalationConfig{}

	for name, attr := range content.Attributes {
		switch name {
		case "escalate":
			config.escalate = attr
		case "impersonate_user":
			config.impersonateUser = attr
		}
	}

	if config.escalate == nil && config.impersonateUser == nil {
		return nil, diags
	}

	return config, diags
}

func parseLoopBlock(block *hcl.Block) (*StepLoopConfig, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	if block == nil {
		return nil, diags
	}

	if block.Type != "loop" {
		return nil, diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block type",
			Detail:   "Expected 'loop' block.",
			Subject:  &block.TypeRange,
		})
	}

	if len(block.Labels) != 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block labels",
			Detail:   "Expected no labels for 'loop' block.",
			Subject:  &block.TypeRange,
		})
	}

	content, moreDiags := block.Body.Content(loopBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a loop block")
	diags = diags.Extend(moreDiags)

	if diags.HasErrors() {
		return nil, diags
	}

	config := &StepLoopConfig{}

	for name, attr := range content.Attributes {
		switch name {
		case "items":
			config.items = attr
		case "label":
			config.label = attr
		case "condition":
			config.condition = attr
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	if config.items == nil && config.label == nil && config.condition == nil {
		return nil, diags
	}

	return config, diags
}

func parseInputBlock(block *hcl.Block) (map[string]*hcl.Attribute, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	if block == nil {
		return nil, diags
	}

	attributes, moreDiags := block.Body.JustAttributes()
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	input := make(map[string]*hcl.Attribute, len(attributes))
	for name, attr := range attributes {
		if _, exists := input[name]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate input",
				Detail:   fmt.Sprintf("The input variable '%s' is defined multiple times.", name),
				Subject:  attr.NameRange.Ptr(),
			})
			continue
		}
		input[name] = attr
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return input, diags
}

func parseOutputBlock(block *hcl.Block) (*StepOutputConfig, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}

	if block == nil {
		return nil, diags
	}

	if block.Type != "output" {
		return nil, diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block type",
			Detail:   "Expected 'output' block.",
			Subject:  &block.TypeRange,
		})
	}

	if len(block.Labels) != 0 {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid block labels",
			Detail:   "Expected no labels for 'output' block.",
			Subject:  &block.TypeRange,
		})
	}

	content, moreDiags := block.Body.Content(outputBlockSchema)
	diags = diags.Extend(moreDiags)

	if diags.HasErrors() {
		return nil, diags
	}

	config := &StepOutputConfig{}

	for name, attr := range content.Attributes {
		switch name {
		case "continue_on_fail":
			config.continueOnFail = attr
		case "changed_condition":
			config.changedCondition = attr
		case "failed_condition":
			config.failedCondition = attr
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	if config.continueOnFail == nil && config.changedCondition == nil && config.failedCondition == nil {
		return nil, diags
	}

	return config, diags
}
