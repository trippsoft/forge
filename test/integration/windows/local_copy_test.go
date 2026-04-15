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

func TestLocalCopyRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Local Copy"
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

			step "local_copy_what_if" {
			    name = "Local Copy What If"
				module = "local_copy"
				what_if = true

				input {
				    source = "C:\\Windows\\Temp\\hello.txt"
					destination = "C:\\Windows\\Temp\\hello_copy.txt"
				}
			}

			step "local_copy" {
			    name = "Local Copy"
				module = "local_copy"

				input {
				    source = "C:\\Windows\\Temp\\hello.txt"
					destination = "C:\\Windows\\Temp\\hello_copy.txt"
				}
			}

			step "local_copy_repeat_what_if" {
			    name = "Local Copy Repeat What If"
				module = "local_copy"
				what_if = true

				input {
				    source = "C:\\Windows\\Temp\\hello.txt"
					destination = "C:\\Windows\\Temp\\hello_copy.txt"
				}
			}

			step "local_copy_repeat" {
			    name = "Local Copy Repeat"
				module = "local_copy"

				input {
				    source = "C:\\Windows\\Temp\\hello.txt"
					destination = "C:\\Windows\\Temp\\hello_copy.txt"
				}
			}

			step "cleanup" {
			    name = "Cleanup"
				module = "command"

				input {
				    name = "powershell.exe"
					args = ["-Command", "Remove-Item 'C:\\Windows\\Temp\\hello.txt' -Force; Remove-Item 'C:\\Windows\\Temp\\hello_copy.txt' -Force"]
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
					"local_copy_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"local_copy": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"local_copy_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"local_copy_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"windows": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
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
