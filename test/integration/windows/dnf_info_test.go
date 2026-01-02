// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
)

func TestDnfInfoRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test DNF Info"
		    targets = [
				"cmd",
				"windows",
			]

			step "dnf_info_what_if" {
			    name = "DNF Info What If"
				module = "dnf_info"
				what_if = true

				input {
				}
			}

			step "dnf_info" {
			    name = "DNF Info"
				module = "dnf_info"

				input {
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
					"dnf_info_what_if": {
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
					"dnf_info": {
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
