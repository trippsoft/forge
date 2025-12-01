// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/workflow"
)

func TestBasicProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "basic_process.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule()

	moduleRegistry.Register("shell", shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Basic Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "setup",
							name: "Setup Environment",
							targets: []*inventory.Host{
								host1,
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

func TestComplexProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "complex_process.hcl")

	db1 := createMockHost("db1")
	db2 := createMockHost("db2")

	i := createMockInventory(db1, db2)

	moduleRegistry := module.NewRegistry()

	backupModule := createMockModule()
	maintenanceModule := createMockModule()

	moduleRegistry.Register("backup", backupModule)
	moduleRegistry.Register("maintenance", maintenanceModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Complex Workflow",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "backup",
							name: "Backup Database",
							targets: []*inventory.Host{
								db1,
								db2,
							},
							input: map[string]struct{}{
								"backup_path": {},
								"compress":    {},
							},
							condition:   true,
							execTimeout: true,
						},
						escalation: &expectedEscalation{
							escalate:        true,
							impersonateUser: true,
						},
						output: &expectedOutput{
							continueOnFail:  true,
							failedCondition: true,
						},
						module: backupModule,
					},
					{
						common: &expectedCommon{
							id:   "maintenance",
							name: "Database Maintenance",
							targets: []*inventory.Host{
								db1,
							},
							loop: &expectedLoop{
								items: true,
								label: true,
							},
							execTimeout: true,
						},
						escalation: &expectedEscalation{
							escalate:        true,
							impersonateUser: true,
						},
						module: maintenanceModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestEscalationProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "escalation_process.hcl")

	privilegedHost := createMockHost("privileged-host")

	i := createMockInventory(privilegedHost)

	moduleRegistry := module.NewRegistry()

	packageModule := createMockModule()
	serviceModule := createMockModule()

	moduleRegistry.Register("package", packageModule)
	moduleRegistry.Register("service", serviceModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Escalated Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "install",
							name: "Install Package",
							targets: []*inventory.Host{
								privilegedHost,
							},
						},
						escalation: &expectedEscalation{
							escalate:        true,
							impersonateUser: true,
						},
						module: packageModule,
					},
					{
						common: &expectedCommon{
							id:   "start_service",
							name: "Start Service",
							targets: []*inventory.Host{
								privilegedHost,
							},
						},
						escalation: &expectedEscalation{
							escalate:        true,
							impersonateUser: true,
						},
						module: serviceModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestInputOutputProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "input_output_process.hcl")

	appServerHost := createMockHost("app-server")

	i := createMockInventory(appServerHost)

	moduleRegistry := module.NewRegistry()

	configModule := createMockModule()

	moduleRegistry.Register("config", configModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Input/Output Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "configure",
							name: "Configure Application",
							targets: []*inventory.Host{
								appServerHost,
							},
							input: map[string]struct{}{
								"app_name":    {},
								"version":     {},
								"config_file": {},
								"debug":       {},
								"port":        {},
							},
						},
						output: &expectedOutput{
							continueOnFail:   true,
							changedCondition: true,
							failedCondition:  true,
						},
						module: configModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestLoopProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "loop_process.hcl")

	web1Host := createMockHost("web1")
	web2Host := createMockHost("web2")
	web3Host := createMockHost("web3")

	i := createMockInventory(web1Host, web2Host, web3Host)

	moduleRegistry := module.NewRegistry()

	deployModule := createMockModule()

	moduleRegistry.Register("deploy", deployModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Looping Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "deploy_services",
							name: "Deploy Services",
							targets: []*inventory.Host{
								web1Host,
								web2Host,
								web3Host,
							},
							loop: &expectedLoop{
								items:     true,
								label:     true,
								condition: true,
							},
						},
						module: deployModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestMultiStepProcess(t *testing.T) {

	path := filepath.Join("corpus", "valid", "multi_step_process.hcl")

	host1 := createMockHost("host1")
	host2 := createMockHost("host2")
	host3 := createMockHost("host3")

	i := createMockInventory(host1, host2, host3)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule()
	configModule := createMockModule()
	deployModule := createMockModule()

	moduleRegistry.Register("shell", shellModule)
	moduleRegistry.Register("config", configModule)
	moduleRegistry.Register("deploy", deployModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "Multi-Step Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "setup",
							name: "Setup Environment",
							targets: []*inventory.Host{
								host1,
								host2,
								host3,
							},
							execTimeout: true,
						},
						module: shellModule,
					},
					{
						common: &expectedCommon{
							id:   "configure",
							name: "Configure System",
							targets: []*inventory.Host{
								host1,
								host2,
								host3,
							},
							execTimeout: true,
							condition:   true,
						},
						module: configModule,
					},
					{
						common: &expectedCommon{
							id:   "deploy",
							name: "Deploy Application",
							targets: []*inventory.Host{
								host1,
							},
							execTimeout: true,
						},
						module: deployModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestMultipleProcesses(t *testing.T) {

	path := filepath.Join("corpus", "valid", "multiple_processes.hcl")

	host1 := createMockHost("host1")
	host2 := createMockHost("host2")
	host3 := createMockHost("host3")

	i := createMockInventory(host1, host2, host3)

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule()
	copyModule := createMockModule()
	serviceModule := createMockModule()

	moduleRegistry.Register("shell", shellModule)
	moduleRegistry.Register("copy", copyModule)
	moduleRegistry.Register("service", serviceModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if diags.HasErrors() {
		t.Fatalf("failed to parse workflow file: %v", diags)
	}

	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	expected := &expected{
		processes: []*expectedProcess{
			{
				name: "First Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "step1",
							name: "First Step",
							targets: []*inventory.Host{
								host1,
							},
						},
						module: shellModule,
					},
				},
			},
			{
				name: "Second Process",
				steps: []*expectedStep{
					{
						common: &expectedCommon{
							id:   "step1",
							name: "Another Step",
							targets: []*inventory.Host{
								host2,
								host3,
							},
						},
						module: copyModule,
					},
					{
						common: &expectedCommon{
							id:   "step2",
							name: "Final Step",
							targets: []*inventory.Host{
								host2,
								host3,
							},
						},
						module: serviceModule,
					},
				},
			},
		},
	}

	expected.verify(t, w)
}

func TestMissingModule(t *testing.T) {

	path := filepath.Join("corpus", "valid", "basic_process.hcl")

	host1 := createMockHost("host1")

	i := createMockInventory(host1)

	moduleRegistry := module.NewRegistry()

	// Leave moduleRegistry empty to simulate missing module

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected diagnostics for missing module, got none")
	}

	if w != nil {
		t.Fatalf("expected workflow to be nil due to missing module, got non-nil")
	}

	expectedDiags := expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Module not found",
			detail:   fmt.Sprintf("Module %q not found", "shell"),
		},
	}

	expectedDiags.verify(t, diags)
}

func TestMissingTarget(t *testing.T) {

	path := filepath.Join("corpus", "valid", "basic_process.hcl")

	i := createMockInventory() // No hosts in inventory

	moduleRegistry := module.NewRegistry()

	shellModule := createMockModule()

	moduleRegistry.Register("shell", shellModule)

	parser := workflow.NewParser(i, moduleRegistry)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read workflow file: %v", err)
	}

	w, diags := parser.ParseWorkflowFile(path, content)
	if !diags.HasErrors() {
		t.Fatalf("expected diagnostics for missing target, got none")
	}

	if w != nil {
		t.Fatalf("expected workflow to be nil due to missing target, got non-nil")
	}

	expectedDiags := expectedDiagnostics{
		{
			severity: hcl.DiagError,
			summary:  "Target not found",
			detail:   fmt.Sprintf("The target %q does not exist in the inventory", "host1"),
		},
	}

	expectedDiags.verify(t, diags)
}
