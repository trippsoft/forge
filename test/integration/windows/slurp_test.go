// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
	"github.com/zclconf/go-cty/cty"
)

func TestSlurpRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Slurp"
		    targets = [
				"cmd",
				"windows",
			]

			step "create_file" {
			    name = "Create File"
				module = "command"

				input {
				    name = "powershell.exe"
					args = ["-Command", "[byte[]]$content = @(0x48,0x65,0x6c,0x6c,0x6f,0x2c,0x20,0x57,0x6f,0x72,0x6c,0x64,0x21,0x0a); [System.IO.File]::WriteAllBytes('C:\\Windows\\Temp\\hello.txt', $content)"]
				}
			}

			step "slurp_what_if" {
			    name = "Slurp What If"
				module = "slurp"
				what_if = true

				input {
				    path = "C:\\Windows\\Temp\\hello.txt"
				}
			}

			step "slurp" {
			    name = "Slurp"
				module = "slurp"

				input {
				    path = "C:\\Windows\\Temp\\hello.txt"
				}
			}

			step "cleanup" {
			    name = "Cleanup"
				module = "command"

				input {
				    name = "powershell.exe"
					args = ["-Command", "Remove-Item 'C:\\Windows\\Temp\\hello.txt' -Force"]
				}
			}
		}
		`

	w := integration.ParseWorkflow(t, inv, moduleRegistry, workflowContent)

	workflowContext, err := workflow.NewWorkflowContext(ui.MockUI, inv, false)
	if err != nil {
		t.Fatalf("Failed to create workflow context: %v", err)
	}

	actual, err := w.Run(workflowContext)
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	expected := integration.ExpectedWorkflowOutput{
		Processes: []integration.ExpectedProcessOutput{
			{
				Steps: map[string]integration.ExpectedStepOutput{
					"create_file": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
						},
					},
					"slurp_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"slurp": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
						},
					},
				},
			},
		},
	}

	expected.Verify(t, actual)
}
