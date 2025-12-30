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

func TestDnfInfoRun_Linux(t *testing.T) {
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
							"debian13": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
					"dnf_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204": {
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

func TestDnfInfoRun_Linux_SudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test DNF Info"
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

			step "dnf_info_what_if" {
			    name = "DNF Info What If"
				module = "dnf_info"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
				}
			}

			step "dnf_info" {
			    name = "DNF Info"
				module = "dnf_info"

				escalate {
				    escalate = true
				}

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
							"debian13-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
					"dnf_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404-pw": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204-pw": {
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

func TestDnfInfoRun_Linux_NoSudoPassword(t *testing.T) {
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
			    name = "DNF Info What If"
				module = "dnf_info"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
				}
			}

			step "dnf_info" {
			    name = "DNF Info"
				module = "dnf_info"

				escalate {
				    escalate = true
				}

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
							"debian13": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
					"dnf_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"debian12": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora42": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"fedora41": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"packages": cty.UnknownVal(
										cty.List(
											cty.Object(map[string]cty.Type{
												"name":         cty.String,
												"epoch":        cty.String,
												"version":      cty.String,
												"release":      cty.String,
												"architecture": cty.String,
												"repo":         cty.String,
											}),
										),
									),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"ubuntu2204": {
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
