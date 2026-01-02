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

func TestCommandRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Command"
		    targets = [
				"cmd",
				"windows",
			]

			step "command_what_if" {
			    name = "Command What If"
				module = "command"
				what_if = true

				input {
				    name = "cmd.exe"
					args = ["/C", "echo Hello, World!"]
				}
			}

			step "command" {
			    name = "Command"
				module = "command"

				input {
				    name = "cmd.exe"
					args = ["/C", "echo Hello, World!"]
				}
			}
		}
		`

	w := integration.ParseWorkflow(t, inv, moduleRegistry, workflowContent)

	actual, err := w.Run(workflow.NewWorkflowContext(ui.MockUI, inv, false))
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	expected := integration.ExpectedWorkflowOutput{
		Processes: []integration.ExpectedProcessOutput{
			{
				Steps: map[string]integration.ExpectedStepOutput{
					"command_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
						},
					},
					"command": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"windows": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
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
