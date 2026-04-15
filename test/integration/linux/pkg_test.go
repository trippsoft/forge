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

func TestPackageRun_Dnf_Linux(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "package_what_if" {
			    name = "Package What If"
				module = "package"
				what_if = true

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package" {
			    name = "Package"
				module = "package"

				input {
					names = ["sos"]
					state = "present"
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
					"package_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky9": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky8": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
						},
					},
					"package": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky9": {
								Changed: false,
								Failed:  true,
								Skipped: false,
							},
							"rocky8": {
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

func TestPackageRun_Dnf_Linux_SudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package"
		    targets = [
				"rocky10-pw",
				"rocky9-pw",
				"rocky8-pw",
			]

			step "package_present_what_if" {
			    name = "Package Present What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present" {
			    name = "Package Present"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present_repeat_what_if" {
			    name = "Package Present Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present_repeat" {
			    name = "Package Present Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_absent_what_if" {
			    name = "Package Absent What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent" {
			    name = "Package Absent"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent_repeat_what_if" {
			    name = "Package Absent Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent_repeat" {
			    name = "Package Absent Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_latest_what_if" {
			    name = "Package Latest What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest" {
			    name = "Package Latest"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest_repeat_what_if" {
			    name = "Package Latest Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest_repeat" {
			    name = "Package Latest Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_cleanup" {
				name = "Package Cleanup"
				module = "package"

				escalate {
					escalate = true
				}
				
				input {
					names = ["sos"]
					state = "absent"
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
					"package_present_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_absent_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
					"package_absent": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
					"package_absent_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_absent_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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

func TestPackageRun_Dnf_Linux_NoSudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Package Info"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "package_present_what_if" {
			    name = "Package Present What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present" {
			    name = "Package Present"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present_repeat_what_if" {
			    name = "Package Present Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_present_repeat" {
			    name = "Package Present Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "package_absent_what_if" {
			    name = "Package Absent What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent" {
			    name = "Package Absent"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent_repeat_what_if" {
			    name = "Package Absent Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_absent_repeat" {
			    name = "Package Absent Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "package_latest_what_if" {
			    name = "Package Latest What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest" {
			    name = "Package Latest"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest_repeat_what_if" {
			    name = "Package Latest Repeat What If"
				module = "package"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_latest_repeat" {
			    name = "Package Latest Repeat"
				module = "package"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "package_cleanup" {
				name = "Package Cleanup"
				module = "package"

				escalate {
					escalate = true
				}
				
				input {
					names = ["sos"]
					state = "absent"
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
					"package_present_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_present_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_absent_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
					"package_absent": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
					"package_absent_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_absent_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.UnknownVal(
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
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_repeat_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_latest_repeat": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
								},
							},
						},
					},
					"package_cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"installed": cty.ListValEmpty(
										cty.Object(map[string]cty.Type{
											"name":         cty.String,
											"epoch":        cty.String,
											"version":      cty.String,
											"release":      cty.String,
											"architecture": cty.String,
											"repo":         cty.String,
										}),
									),
									"removed": cty.UnknownVal(
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
