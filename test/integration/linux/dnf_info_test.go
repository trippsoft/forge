// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package linux

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
)

func TestShellRun_Linux(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test DNF Info"
		    targets = [
				"debian13",
				"debian12",
				"fedora42",
				"fedora41",
				"rocky10",
				"rocky9",
				"rocky8",
				"ubuntu2404",
				"ubuntu2204",
			]

			step "dnf_info_what_if" {
			    name = "DNF Info What-If"
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

	_ = actual
}
