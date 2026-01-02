// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
)

func TestDnfRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test DNF"
		    targets = [
				"cmd",
				"windows",
			]

			step "dnf_what_if" {
			    name = "DNF What If"
				module = "dnf"
				what_if = true

				input {
					names = ["vim"]
					state = "present"
				}
			}

			step "dnf" {
			    name = "DNF"
				module = "dnf"

				input {
					names = ["vim"]
					state = "present"
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
					"dnf_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"windows": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
					"dnf": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"cmd": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"windows": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
				},
			},
		},
	}

	expected.Verify(t, actual)
}
