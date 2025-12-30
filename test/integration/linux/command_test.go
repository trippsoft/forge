// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package linux

import (
	"testing"

	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/workflow"
	"github.com/trippsoft/forge/test/integration"
	"github.com/zclconf/go-cty/cty"
)

func TestCommandRun_Linux(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Command"
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

			step "command_what_if" {
			    name = "Command What If"
				module = "command"
				what_if = true

				input {
				    name = "echo"
					args = ["Hello, World!"]
				}
			}

			step "command" {
			    name = "Command"
				module = "command"

				input {
				    name = "echo"
					args = ["Hello, World!"]
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
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2204": {
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
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
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

func TestCommandRun_Linux_SudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Shell"
		    targets = [
				"debian13-pw",
				"debian12-pw",
				"fedora42-pw",
				"fedora41-pw",
				"rocky10-pw",
				"rocky9-pw",
				"rocky8-pw",
				"ubuntu2404-pw",
				"ubuntu2204-pw",
			]

			step "command_what_if" {
			    name = "Command What If"
				module = "command"
				what_if = true

				escalate {
					escalate = true
				}

				input {
					name = "echo"
					args = ["Hello, World!"]
				}
			}

			step "command" {
			    name = "Command"
				module = "command"

				escalate {
					escalate = true
				}

				input {
					name = "echo"
					args = ["Hello, World!"]
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
							"debian13-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"debian12-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora42-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora41-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2404-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2204-pw": {
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
							"debian13-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204-pw": {
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

func TestCommandRun_Linux_NoSudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Shell"
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

			step "command_what_if" {
			    name = "Command What If"
				module = "command"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					name = "echo"
					args = ["Hello, World!"]
				}
			}

			step "command" {
			    name = "Command"
				module = "command"

				escalate {
					escalate = true
				}

				input {
					name = "echo"
					args = ["Hello, World!"]
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
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.NullVal(cty.String),
									"stderr": cty.NullVal(cty.String),
								},
							},
							"ubuntu2204": {
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
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal("Hello, World!"),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
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
