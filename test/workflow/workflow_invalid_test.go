// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/workflow"
)

func TestInvalidBlockType(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "invalid_block_type.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Test Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:      "test",
							name:    "Test Step",
							targets: []*inventory.Host{host1},
						},
						module: shellModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagWarning,
			summary:  "Unsupported block type",
			detail:   "Blocks of type \"invalid_block\" are not expected in a process block.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestInvalidStepAttribute(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "invalid_step_attribute.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Test Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:      "test",
							name:    "Test Step",
							targets: []*inventory.Host{host1},
						},
						module: shellModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagWarning,
			summary:  "Unsupported argument",
			detail:   "An argument named \"invalid_attribute\" is not expected in a step block.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestLoopMissingItems(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "loop_missing_items.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing required argument",
			detail:   "The argument \"items\" is required, but no definition was found.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMissingProcessName(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "missing_process_name.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing required argument",
			detail:   "The argument \"name\" is required, but no definition was found.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMissingProcessTargets(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "missing_process_targets.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing required argument",
			detail:   "The argument \"targets\" is required, but no definition was found.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMissingStepModule(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "missing_step_module.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing required argument",
			detail:   "The argument \"module\" is required, but no definition was found.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMissingStepName(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "missing_step_name.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing required argument",
			detail:   "The argument \"name\" is required, but no definition was found.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMultipleEscalateBlocks(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "multiple_escalate_blocks.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Duplicate escalate block",
			detail:   "Only one escalate block is allowed per process.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMultipleInputBlocks(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "multiple_input_blocks.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Duplicate input block",
			detail:   "Only one input block is allowed.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMulitpleLoopBlocks(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "multiple_loop_blocks.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Duplicate loop block",
			detail:   "Only one loop block is allowed.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestProcessWithLabels(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "process_with_labels.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Extraneous label for process",
			detail:   "No labels are expected for process blocks.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestStepMissingLabel(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "step_missing_label.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing id for step",
			detail:   "All step blocks must have 1 labels (id).",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestStepTooManyLabels(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "step_too_many_labels.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Extraneous label for step",
			detail:   "Only 1 labels (id) are expected for step blocks.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestSyntaxErrorMissingBrace(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "syntax_error_missing_brace.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Unclosed configuration block",
			detail:   "There is no closing brace for this block before the end of the file. This may be caused by incorrect brace nesting elsewhere in this file.",
		},
	}

	expectedDiags.verify(t, diags)
}

func TestSyntaxErrorMissingQuotes(t *testing.T) {

	path := filepath.Join("corpus", "invalid", "syntax_error_missing_quotes.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")

	moduleRegistry.Register(shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected error, got none")
	}

	if w != nil {
		t.Fatalf("expected nil workflow, got: %v", w)
	}

	expectedDiags := &expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Missing newline after argument",
			detail:   "An argument definition must end with a newline.",
		},
	}

	expectedDiags.verify(t, diags)
}
