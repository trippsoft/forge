// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/workflow"
)

func TestCommentsWhitespace(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "comments_whitespace.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule("shell")
	copyModule := createMockModule("copy")

	moduleRegistry.Register(shellModule)
	moduleRegistry.Register(copyModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Comment Test Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:        "test",
							name:      "Test Step",
							targets:   []*inventory.Host{host1},
							condition: true,
						},
						module: shellModule,
					},
					{
						common: &expectedCommon{
							id:      "another",
							name:    "Another Step",
							targets: []*inventory.Host{host1},
						},
						module: copyModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestDataTypesInput(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "data_types_input.hcl")

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Data Types Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:      "test",
							name:    "Test Step",
							targets: []*inventory.Host{host1},
							input: map[string]struct{}{
								"string_val":   {},
								"number_val":   {},
								"float_val":    {},
								"bool_val":     {},
								"array_val":    {},
								"null_val":     {},
								"empty_string": {},
								"zero_val":     {},
								"false_val":    {},
							},
						},
						module: shellModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestEmptyEscalateBlock(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "empty_escalate_block.hcl")

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Empty Escalate Process",
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
}

func TestEmptyFile(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "empty_file.hcl")

	i := createMockInventory() // No hosts in inventory

	moduleRegistry := module.NewRegistry() // Empty module registry

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{},
	}

	expected.verify(t, w)
}

func TestEmptyInputBlock(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "empty_input_block.hcl")

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Empty Input Process",
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
}

func TestEmptyOutputBlock(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "empty_output_block.hcl")

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Empty Output Process",
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
}

func TestEmptyProcess(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "empty_process.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name:  "Empty Process",
				steps: []*expectedStep{},
			},
		},
	}

	expected.verify(t, w)
}

func TestLongNames(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "long_names.hcl")

	host := createMockHost("host-with-very-long-name-that-should-still-be-parsed-correctly")

	i := createMockInventory(host)

	moduleRegistry := module.NewRegistry()

	module := createMockModule("module-with-long-name")

	moduleRegistry.Register(module)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "A Very Long Process Name That Tests The Parser's Ability To Handle Extended Names Without Issues",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:      "step_with_very_long_identifier_name_that_tests_parser_limits",
							name:    "A Very Long Step Name That Also Tests The Parser's String Handling Capabilities",
							targets: []*inventory.Host{host},
						},
						module: module,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestOnlyComments(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "only_comments.hcl")

	i := createMockInventory() // No hosts in inventory

	moduleRegistry := module.NewRegistry() // Empty module registry

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{},
	}

	expected.verify(t, w)
}

func TestSingleTargetString(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "single_target_string.hcl")

	singleHost := createMockHost("single-host")
	differentHost := createMockHost("different-host")

	i := createMockInventory(singleHost, differentHost)

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Single Target Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:      "test",
							name:    "Test Step",
							targets: []*inventory.Host{differentHost},
						},
						module: shellModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestTargetInheritance(t *testing.T) {

	path := filepath.Join("corpus", "edge_cases", "target_inheritance.hcl")

	parent1 := createMockHost("parent1")
	parent2 := createMockHost("parent2")
	override1 := createMockHost("override1")

	i := createMockInventory(parent1, parent2, override1)

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

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Target Inheritance Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:          "inherits_targets",
							name:        "Step That Inherits Targets",
							targets:     []*inventory.Host{parent1, parent2},
							execTimeout: true,
						},
						module: shellModule,
					},
					{
						common: &expectedCommon{
							id:          "overrides_targets",
							name:        "Step That Overrides Targets",
							targets:     []*inventory.Host{override1},
							execTimeout: true,
						},
						module: shellModule,
					},
					{
						common: &expectedCommon{
							id:          "inherits_timeout",
							name:        "Step That Inherits Timeout",
							targets:     []*inventory.Host{parent1, parent2},
							execTimeout: true,
						},
						module: shellModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}
