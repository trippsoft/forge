// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"slices"
	"testing"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/zclconf/go-cty/cty"
)

type ExpectedHostOutput struct {
	Changed bool
	Failed  bool
	Skipped bool
	Output  map[string]cty.Value
}

func (e ExpectedHostOutput) Verify(
	t *testing.T,
	hostOutput map[string]cty.Value,
	processIndex int,
	stepName string,
	hostName string,
) {

	t.Helper()

	changed, ok := hostOutput["changed"]
	if !ok {
		t.Errorf(
			`Expected host %q in step %q in process %d output to contain "changed" key`,
			hostName,
			stepName,
			processIndex,
		)
	} else if changed.True() != e.Changed {
		t.Errorf(
			"Expected host %q in step %q in process %d changed value to be %t, got %t",
			hostName,
			stepName,
			processIndex,
			e.Changed,
			changed.True(),
		)
	}

	failed, ok := hostOutput["failed"]
	if !ok {
		t.Errorf(
			`Expected host %q in step %q in process %d output to contain "failed" key`,
			hostName,
			stepName,
			processIndex,
		)
	} else if failed.True() != e.Failed {
		t.Errorf(
			"Expected host %q in step %q in process %d failed value to be %t, got %t",
			hostName,
			stepName,
			processIndex,
			e.Failed,
			failed.True(),
		)
	}

	skipped, ok := hostOutput["skipped"]
	if !ok {
		t.Errorf(
			`Expected host %q in step %q in process %d output to contain "skipped" key`,
			hostName,
			stepName,
			processIndex,
		)
	} else if skipped.True() != e.Skipped {
		t.Errorf(
			"Expected host %q in step %q in process %d skipped value to be %t, got %t",
			hostName,
			stepName,
			processIndex,
			e.Skipped,
			skipped.True(),
		)
	}

	if e.Failed || e.Skipped {
		return // Skip further validation if the step failed or was skipped
	}

	output, ok := hostOutput["output"]
	if !ok {
		t.Errorf(
			`Expected host %q in step %q in process %d output to contain "output" key`,
			hostName,
			stepName,
			processIndex,
		)
		return
	}

	if !output.Type().IsObjectType() && !output.Type().IsMapType() {
		t.Errorf(
			"Expected host %q in step %q in process %d output value to be an object or map, got %s",
			hostName,
			stepName,
			processIndex,
			output.Type().FriendlyName(),
		)
		return
	}

	outputMap := output.AsValueMap()

	if len(outputMap) != len(e.Output) {
		t.Errorf(
			"Expected host %q in step %q in process %d output to contain %d keys, got %d",
			hostName,
			stepName,
			processIndex,
			len(e.Output),
			len(outputMap),
		)
	}

	for key, expected := range e.Output {
		actual, ok := outputMap[key]
		if !ok {
			t.Errorf(
				"Expected host %q in step %q in process %d output to contain key %q",
				hostName,
				stepName,
				processIndex,
				key,
			)
			continue
		}

		if !expected.IsWhollyKnown() {
			continue // Use UnknownVal to skip validation
		}

		if actual.Equals(expected).False() {
			t.Errorf(
				"Expected host %q in step %q in process %d output %q value to be %q, got %q",
				hostName,
				stepName,
				processIndex,
				key,
				expected.GoString(),
				actual.GoString(),
			)
		}
	}
}

type ExpectedStepOutput struct {
	Hosts map[string]ExpectedHostOutput
}

func (e ExpectedStepOutput) Verify(t *testing.T, stepOutput map[string]cty.Value, processIndex int, stepName string) {
	t.Helper()

	if len(stepOutput) != len(e.Hosts) {
		t.Errorf(
			"Expected step %q in process %d output to contain %d hosts, got %d",
			stepName,
			processIndex,
			len(e.Hosts),
			len(stepOutput),
		)
	}

	for hostName, expected := range e.Hosts {
		actual, ok := stepOutput[hostName]
		if !ok {
			t.Errorf("Expected step %q in process %d output to contain host %q", stepName, processIndex, hostName)
			continue
		}

		if !actual.Type().IsObjectType() && !actual.Type().IsMapType() {
			t.Errorf(
				"Expected host %q in step %q in process %d output to be an object or map, got: %v",
				hostName,
				stepName,
				processIndex,
				actual.Type().FriendlyName(),
			)
			continue
		}

		actualMap := actual.AsValueMap()

		expected.Verify(t, actualMap, processIndex, stepName, hostName)
	}
}

type ExpectedProcessOutput struct {
	Steps map[string]ExpectedStepOutput
}

func (e ExpectedProcessOutput) Verify(t *testing.T, processOutput map[string]map[string]cty.Value, index int) {
	t.Helper()

	if len(processOutput) != len(e.Steps) {
		t.Errorf("Expected process %d output to contain %d steps, got %d", index, len(e.Steps), len(processOutput))
	}

	for stepName, expected := range e.Steps {
		actual, ok := processOutput[stepName]
		if !ok {
			t.Errorf("Expected process %d output to contain step %q", index, stepName)
			continue
		}

		expected.Verify(t, actual, index, stepName)
	}
}

type ExpectedWorkflowOutput struct {
	Processes []ExpectedProcessOutput
}

func (e ExpectedWorkflowOutput) Verify(t *testing.T, workflowOutput []map[string]map[string]cty.Value) {
	t.Helper()

	if len(workflowOutput) != len(e.Processes) {
		t.Errorf("Expected workflow output to contain %d processes, got %d", len(e.Processes), len(workflowOutput))
	}

	for i, expected := range e.Processes {
		actual := workflowOutput[i]
		expected.Verify(t, actual, i)
	}
}

type ExpectedOSInfo struct {
	Kernel       string
	Id           string
	FriendlyName string
	Release      string
	ReleaseId    string
	MajorVersion string
	Version      string
	Edition      string
	EditionID    string
	Arch         string

	Families []string
}

func (e *ExpectedOSInfo) Verify(t *testing.T, actual *info.OSInfo) {
	if actual == nil {
		t.Fatal("actual OSInfo is nil")
		return
	}

	for _, family := range e.Families {
		if !slices.Contains(actual.Families, family) {
			t.Errorf("expected OS families to contain %q, got %v", family, actual.Families)
		}
	}

	for _, family := range actual.Families {
		if !slices.Contains(e.Families, family) {
			t.Errorf("unexpected OS family %q found, expected only %v", family, e.Families)
		}
	}

	if actual.Kernel != e.Kernel {
		t.Errorf("expected OS kernel to be %q, got %q", e.Kernel, actual.Kernel)
	}

	if actual.Id != e.Id {
		t.Errorf("expected OS ID to be %q, got %q", e.Id, actual.Id)
	}

	if actual.FriendlyName != e.FriendlyName {
		t.Errorf("expected OS friendly name to be %q, got %q", e.FriendlyName, actual.FriendlyName)
	}

	if actual.Release != e.Release {
		t.Errorf("expected OS release to be %q, got %q", e.Release, actual.Release)
	}

	if actual.ReleaseId != e.ReleaseId {
		t.Errorf("expected OS release ID to be %q, got %q", e.ReleaseId, actual.ReleaseId)
	}

	if actual.MajorVersion != e.MajorVersion {
		t.Errorf("expected OS major version to be %q, got %q", e.MajorVersion, actual.MajorVersion)
	}

	if actual.Version != e.Version {
		t.Errorf("expected OS version to be %q, got %q", e.Version, actual.Version)
	}

	if actual.Edition != e.Edition {
		t.Errorf("expected OS edition to be %q, got %q", e.Edition, actual.Edition)
	}

	if actual.EditionId != e.EditionID {
		t.Errorf("expected OS edition ID to be %q, got %q", e.EditionID, actual.EditionId)
	}

	if actual.Arch != e.Arch {
		t.Errorf("expected OS architecture to be %q, got %q", e.Arch, actual.Arch)
	}
}

type ExpectedSELinuxInfo struct {
	Supported bool
	Installed bool
	Status    string
	Type      string
}

func (e *ExpectedSELinuxInfo) Verify(t *testing.T, actual *info.SELinuxInfo) {
	if actual == nil {
		t.Fatal("actual SELinuxInfo is nil")
		return
	}

	if actual.Supported != e.Supported {
		t.Errorf("expected SELinux supported to be %t, got %t", e.Supported, actual.Supported)
	}

	if actual.Installed != e.Installed {
		t.Errorf("expected SELinux installed to be %t, got %t", e.Installed, actual.Installed)
	}

	if actual.Status != e.Status {
		t.Errorf("expected SELinux status to be %q, got %q", e.Status, actual.Status)
	}

	if actual.Type != e.Type {
		t.Errorf("expected SELinux type to be %q, got %q", e.Type, actual.Type)
	}
}

type ExpectedAppArmorInfo struct {
	Supported bool
	Enabled   bool
}

func (e *ExpectedAppArmorInfo) Verify(t *testing.T, actual *info.AppArmorInfo) {
	if actual == nil {
		t.Fatal("actual AppArmorInfo is nil")
		return
	}

	if actual.Supported != e.Supported {
		t.Errorf("expected AppArmor supported to be %t, got %t", e.Supported, actual.Supported)
	}

	if actual.Enabled != e.Enabled {
		t.Errorf("expected AppArmor enabled to be %t, got %t", e.Enabled, actual.Enabled)
	}
}

type ExpectedFIPSInfo struct {
	Known   bool
	Enabled bool
}

func (e *ExpectedFIPSInfo) Verify(t *testing.T, actual *info.FIPSInfo) {
	if actual == nil {
		t.Fatal("actual FIPSInfo is nil")
		return
	}

	if actual.Known != e.Known {
		t.Errorf("expected FIPS known to be %t, got %t", e.Known, actual.Known)
	}

	if actual.Enabled != e.Enabled {
		t.Errorf("expected FIPS enabled to be %t, got %t", e.Enabled, actual.Enabled)
	}
}

type ExpectedPackageManagerInfo struct {
	Name string
	Path string
}

func (e *ExpectedPackageManagerInfo) Verify(t *testing.T, actual *info.PackageManagerInfo) {
	if actual == nil {
		t.Fatal("actual PackageManagerInfo is nil")
		return
	}

	if actual.Name != e.Name {
		t.Errorf("expected Package Manager name to be %q, got %q", e.Name, actual.Name)
	}

	if actual.Path != e.Path {
		t.Errorf("expected Package Manager path to be %q, got %q", e.Path, actual.Path)
	}
}

type ExpectedServiceManagerInfo struct {
	Name string
}

func (e *ExpectedServiceManagerInfo) Verify(t *testing.T, actual *info.ServiceManagerInfo) {
	if actual == nil {
		t.Fatal("actual ServiceManagerInfo is nil")
		return
	}

	if actual.Name != e.Name {
		t.Errorf("expected Service Manager name to be %q, got %q", e.Name, actual.Name)
	}
}

type ExpectedHostInfo struct {
	OS             ExpectedOSInfo
	SELinux        ExpectedSELinuxInfo
	AppArmor       ExpectedAppArmorInfo
	FIPS           ExpectedFIPSInfo
	PackageManager ExpectedPackageManagerInfo
	ServiceManager ExpectedServiceManagerInfo
}

func (e *ExpectedHostInfo) Verify(t *testing.T, actual *info.HostInfo) {
	if actual == nil {
		t.Fatal("actual HostInfo is nil")
		return
	}

	e.OS.Verify(t, actual.Os)
	e.SELinux.Verify(t, actual.Selinux)
	e.AppArmor.Verify(t, actual.AppArmor)
	e.FIPS.Verify(t, actual.Fips)
	e.PackageManager.Verify(t, actual.PackageManager)
	e.ServiceManager.Verify(t, actual.ServiceManager)
}

func ParseWorkflow(
	t *testing.T,
	inv *inventory.Inventory,
	moduleRegistry *module.Registry,
	content string,
) *workflow.Workflow {

	t.Helper()

	parser := workflow.NewParser(inv, moduleRegistry)

	w, diags := parser.ParseWorkflowFile("test_workflow.hcl", []byte(content))
	if diags.HasErrors() {
		t.Fatalf("Failed to parse workflow file: %v", diags)
	}

	return w
}
