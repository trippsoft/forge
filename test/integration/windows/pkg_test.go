// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
)

func TestPackageRun_Windows(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package"
		    targets = [
				"cmd",
				"windows",
			]

			step "package_what_if" {
			    name = "Package What If"
				module = "package"
				what_if = true

				input {
					names = ["dummy-package"]
					state = "present"
				}
			}

			step "package" {
			    name = "Package"
				module = "package"

				input {
					names = ["dummy-package"]
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
					"package_what_if": {
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
					"package": {
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
