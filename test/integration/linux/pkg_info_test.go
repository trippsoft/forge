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

func TestPackageInfoRun_Dnf_Linux(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package Info"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "package_info_what_if" {
			    name = "Package Info What If"
				module = "package_info"
				what_if = true

				input {
				}
			}

			step "package_info" {
			    name = "Package Info"
				module = "package_info"

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
					"package_info_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
					"package_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
				},
			},
		},
	}

	expected.Verify(t, actual)
}

func TestPackageInfoRun_Dnf_Linux_SudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package Info"
		    targets = [
				"rocky10-pw",
				"rocky9-pw",
				"rocky8-pw",
			]

			step "package_info_what_if" {
			    name = "Package Info What If"
				module = "package_info"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
				}
			}

			step "package_info" {
			    name = "Package Info"
				module = "package_info"

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
					"package_info_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
					"package_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
				},
			},
		},
	}

	expected.Verify(t, actual)
}

func TestPackageInfoRun_Dnf_Linux_NoSudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package Info"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "package_info_what_if" {
			    name = "Package Info What If"
				module = "package_info"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
				}
			}

			step "package_info" {
			    name = "Package Info"
				module = "package_info"

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
					"package_info_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
					"package_info": {
						Hosts: map[string]integration.ExpectedHostOutput{
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
						},
					},
				},
			},
		},
	}

	expected.Verify(t, actual)
}
