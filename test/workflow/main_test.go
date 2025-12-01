// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/zclconf/go-cty/cty"
)

type mockModule struct {
	inputSpec    *hclspec.Spec
	validateFunc func(config *module.RunConfig) error
	Result       *result.Result
}

func newMockModule(spec *hclspec.Spec, validate func(config *module.RunConfig) error) *mockModule {
	return &mockModule{
		inputSpec:    spec,
		validateFunc: validate,
	}
}

func (m *mockModule) InputSpec() *hclspec.Spec {
	return m.inputSpec
}

func (m *mockModule) Validate(config *module.RunConfig) error {
	if m.validateFunc == nil {
		return nil
	}

	return m.validateFunc(config)
}

func (m *mockModule) Run(ctx context.Context, config *module.RunConfig) *result.Result {
	return m.Result
}

func createMockHost(name string) *inventory.Host {
	builder := inventory.NewHostBuilder()
	host, _ := builder.WithName(name).
		WithTransport(transport.NewMockTransport()).
		WithEscalateConfig(inventory.NewEscalateConfig("")).
		WithVars(map[string]cty.Value{}).
		Build()

	return host
}

func createMockInventory(h ...*inventory.Host) *inventory.Inventory {
	hosts := map[string]*inventory.Host{}
	groups := map[string][]*inventory.Host{}
	targets := map[string][]*inventory.Host{
		"all": {},
	}

	for _, host := range h {
		hosts[host.Name()] = host
		targets["all"] = append(targets["all"], host)
		targets[host.Name()] = []*inventory.Host{host}
	}

	return inventory.NewInventory(hosts, groups, targets)
}

func createMockModule() *mockModule {
	return newMockModule(&hclspec.Spec{}, nil)
}

type expectedLoop struct {
	items     bool
	label     bool
	condition bool
}

func (e *expectedLoop) verify(t *testing.T, actual *workflow.StepLoopConfig) {

	if actual == nil {
		t.Fatalf("expected loop config to be non-nil, got nil")
	}

	if e.items {
		if actual.Items() == nil {
			t.Errorf("expected loop items to be present, got nil")
		}
	} else {
		if actual.Items() != nil {
			t.Error("expected loop items to be nil, got non-nil")
		}
	}

	if e.label {
		if actual.Label() == nil {
			t.Errorf("expected loop label to be present, got nil")
		}
	} else {
		if actual.Label() != nil {
			t.Error("expected loop label to be nil, got non-nil")
		}
	}

	if e.condition {
		if actual.Condition() == nil {
			t.Errorf("expected loop condition to be present, got nil")
		}
	} else {
		if actual.Condition() != nil {
			t.Errorf("expected loop condition to be nil, got non-nil")
		}
	}
}

type expectedCommon struct {
	loop *expectedLoop

	id   string
	name string

	targets []*inventory.Host

	condition bool

	execTimeout bool

	input map[string]struct{}
}

func (e *expectedCommon) verify(t *testing.T, actual *workflow.StepCommonConfig) {

	if actual == nil {
		t.Fatalf("expected common config to be non-nil, got nil")
	}

	if e.loop != nil {
		e.loop.verify(t, actual.Loop())
	} else if actual.Loop() != nil {
		t.Fatalf("expected loop config to be non-nil, got nil")
	}

	if e.id != actual.ID() {
		t.Errorf("expected step ID %q, got %q", e.id, actual.ID())
	}

	if e.name != actual.Name() {
		t.Errorf("expected step name %q, got %q", e.name, actual.Name())
	}

	if len(e.targets) != len(actual.Targets()) {
		t.Errorf("expected %d targets, got %d", len(e.targets), len(actual.Targets()))
	}

	for _, target := range e.targets {
		if !slices.Contains(actual.Targets(), target) {
			t.Errorf("unexpected target %q", target.Name())
		}
	}

	for _, actualTarget := range actual.Targets() {
		if !slices.Contains(e.targets, actualTarget) {
			t.Errorf("missing target %q", actualTarget.Name())
		}
	}

	if e.condition {
		if actual.Condition() == nil {
			t.Errorf("expected step condition to be present, got nil")
		}
	} else {
		if actual.Condition() != nil {
			t.Error("expected step condition to be nil, got non-nil")
		}
	}

	if e.execTimeout {
		if actual.ExecTimeout() == nil {
			t.Errorf("expected step exec_timeout to be present, got nil")
		}
	} else {
		if actual.ExecTimeout() != nil {
			t.Error("expected step exec_timeout to be nil, got non-nil")
		}
	}

	if e.input == nil {
		e.input = map[string]struct{}{}
	}

	if len(e.input) != len(actual.Input()) {
		t.Errorf("expected step input to have %d items, got %d", len(e.input), len(actual.Input()))
	}

	for key := range e.input {
		if _, ok := actual.Input()[key]; !ok {
			t.Errorf("missing step input %q", key)
		}
	}

	for key := range actual.Input() {
		if _, ok := e.input[key]; !ok {
			t.Errorf("unexpected step input %q", key)
		}
	}
}

type expectedEscalation struct {
	escalate        bool
	impersonateUser bool
}

func (e *expectedEscalation) verify(t *testing.T, actual *workflow.StepEscalateConfig) {

	if actual == nil {
		t.Fatal("expected step escalate config to be non-nil, got nil")
	}

	if e.escalate {
		if actual.Escalate() == nil {
			t.Errorf("expected step escalate to be present, got nil")
		}
	} else {
		if actual.Escalate() != nil {
			t.Error("expected step escalate to be nil, got non-nil")
		}
	}

	if e.impersonateUser {
		if actual.ImpersonateUser() == nil {
			t.Errorf("expected step impersonate_user to be present, got nil")
		}
	} else {
		if actual.ImpersonateUser() != nil {
			t.Error("expected step impersonate_user to be nil, got non-nil")
		}
	}
}

type expectedOutput struct {
	continueOnFail   bool
	changedCondition bool
	failedCondition  bool
}

func (e *expectedOutput) verify(t *testing.T, actual *workflow.StepOutputConfig) {

	if actual == nil {
		t.Fatalf("expected step output config to be non-nil, got nil")
	}

	if e.continueOnFail {
		if actual.ContinueOnFail() == nil {
			t.Errorf("expected step continue_on_fail to be present, got nil")
		}
	} else {
		if actual.ContinueOnFail() != nil {
			t.Error("expected step continue_on_fail to be nil, got non-nil")
		}
	}

	if e.changedCondition {
		if actual.ChangedCondition() == nil {
			t.Errorf("expected step changed_condition to be present, got nil")
		}
	} else {
		if actual.ChangedCondition() != nil {
			t.Error("expected step changed_condition to be nil, got non-nil")
		}
	}

	if e.failedCondition {
		if actual.FailedCondition() == nil {
			t.Errorf("expected step failed_condition to be present, got nil")
		}
	} else {
		if actual.FailedCondition() != nil {
			t.Error("expected step failed_condition to be nil, got non-nil")
		}
	}
}

type expectedStep struct {
	common     *expectedCommon
	escalation *expectedEscalation
	output     *expectedOutput

	module *mockModule
}

func (e *expectedStep) verify(t *testing.T, a workflow.Step) {

	if a == nil {
		t.Fatalf("expected step to be non-nil, got nil")
	}

	actual, ok := a.(*workflow.SingleStep)
	if !ok {
		t.Fatalf("expected step to be of type *workflow.SingleStep")
	}

	actualCommon := actual.Common()
	if e.common != nil {
		e.common.verify(t, actualCommon)
	} else if actualCommon != nil {
		t.Fatal("expected common config to be nil, got non-nil")
	}

	actualEscalation := actual.Escalate()
	if e.escalation != nil {
		e.escalation.verify(t, actualEscalation)
	} else if actualEscalation != nil {
		t.Fatal("expected escalation config to be nil, got non-nil")
	}

	actualOutput := actual.Output()
	if e.output != nil {
		e.output.verify(t, actualOutput)
	} else if actualOutput != nil {
		t.Fatal("expected output config to be nil, got non-nil")
	}

	if e.module != actual.Module() {
		t.Errorf("expected step module to be %v, got %v", e.module, actual.Module())
	}
}

type expectedProcess struct {
	name  string
	steps []*expectedStep
}

func (e *expectedProcess) verify(t *testing.T, actual *workflow.Process) {

	if actual == nil {
		t.Fatalf("expected process to be non-nil, got nil")
	}

	if e.name != actual.Name() {
		t.Errorf("expected process name %q, got %q", e.name, actual.Name())
	}

	actualSteps := actual.Steps()
	if len(e.steps) != len(actualSteps) {
		t.Fatalf("expected %d steps, got %d", len(e.steps), len(actualSteps))
	}

	for i := range e.steps {
		e.steps[i].verify(t, actualSteps[i])
	}
}

type expected struct {
	processes []*expectedProcess
}

func (e *expected) verify(t *testing.T, actual *workflow.Workflow) {

	actualProcesses := actual.Processes()
	if len(e.processes) != len(actualProcesses) {
		t.Fatalf("expected %d processes, got %d", len(e.processes), len(actualProcesses))
	}

	for i := range e.processes {
		e.processes[i].verify(t, actualProcesses[i])
	}
}

type expectedDiagnostic struct {
	severity hcl.DiagnosticSeverity
	summary  string
	detail   string
}

func (e *expectedDiagnostic) verify(t *testing.T, actual *hcl.Diagnostic) {

	if actual == nil {
		t.Fatalf("expected diagnostic to be non-nil, got nil")
	}

	if e.severity != actual.Severity {
		t.Errorf("expected diagnostic severity %q, got %q", e.severity, actual.Severity)
	}

	if e.summary != actual.Summary {
		t.Errorf("expected diagnostic summary %q, got %q", e.summary, actual.Summary)
	}

	if e.detail != actual.Detail {
		t.Errorf("expected diagnostic detail %q, got %q", e.detail, actual.Detail)
	}
}

type expectedDiagnostics []*expectedDiagnostic

func (e expectedDiagnostics) verify(t *testing.T, actual hcl.Diagnostics) {

	if len(e) != len(actual) {
		t.Fatalf("expected %d diagnostics, got %d", len(e), len(actual))
	}

	for i := range e {
		e[i].verify(t, actual[i])
	}
}
