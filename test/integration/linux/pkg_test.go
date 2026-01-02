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
			name = "Test DNF"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "dnf_what_if" {
			    name = "DNF What If"
				module = "dnf"
				what_if = true

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf" {
			    name = "DNF"
				module = "dnf"

				input {
					names = ["sos"]
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
					"dnf": {
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
			name = "Test DNF"
		    targets = [
				"rocky10-pw",
				"rocky9-pw",
				"rocky8-pw",
			]

			step "dnf_present_what_if" {
			    name = "DNF Present What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present" {
			    name = "DNF Present"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present_repeat_what_if" {
			    name = "DNF Present Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present_repeat" {
			    name = "DNF Present Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_absent_what_if" {
			    name = "DNF Absent What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent" {
			    name = "DNF Absent"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent_repeat_what_if" {
			    name = "DNF Absent Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent_repeat" {
			    name = "DNF Absent Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_latest_what_if" {
			    name = "DNF Latest What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest" {
			    name = "DNF Latest"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest_repeat_what_if" {
			    name = "DNF Latest Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest_repeat" {
			    name = "DNF Latest Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_cleanup" {
				name = "DNF Cleanup"
				module = "dnf"

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

	actual, err := w.Run(workflow.NewWorkflowContext(ui.MockUI, inv, false))
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	expected := integration.ExpectedWorkflowOutput{
		Processes: []integration.ExpectedProcessOutput{
			{
				Steps: map[string]integration.ExpectedStepOutput{
					"dnf_present_what_if": {
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
					"dnf_present": {
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
					"dnf_present_repeat_what_if": {
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
					"dnf_present_repeat": {
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
					"dnf_absent_what_if": {
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
					"dnf_absent": {
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
					"dnf_absent_repeat_what_if": {
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
					"dnf_absent_repeat": {
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
					"dnf_latest_what_if": {
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
					"dnf_latest": {
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
					"dnf_latest_repeat_what_if": {
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
					"dnf_latest_repeat": {
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
					"dnf_cleanup": {
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
			name = "Test DNF Info"
		    targets = [
				"rocky10",
				"rocky9",
				"rocky8",
			]

			step "dnf_present_what_if" {
			    name = "DNF Present What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present" {
			    name = "DNF Present"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present_repeat_what_if" {
			    name = "DNF Present Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_present_repeat" {
			    name = "DNF Present Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "present"
				}
			}

			step "dnf_absent_what_if" {
			    name = "DNF Absent What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent" {
			    name = "DNF Absent"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent_repeat_what_if" {
			    name = "DNF Absent Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_absent_repeat" {
			    name = "DNF Absent Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "absent"
				}
			}

			step "dnf_latest_what_if" {
			    name = "DNF Latest What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest" {
			    name = "DNF Latest"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest_repeat_what_if" {
			    name = "DNF Latest Repeat What If"
				module = "dnf"
				what_if = true

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_latest_repeat" {
			    name = "DNF Latest Repeat"
				module = "dnf"

				escalate {
				    escalate = true
				}

				input {
					names = ["sos"]
					state = "latest"
				}
			}

			step "dnf_cleanup" {
				name = "DNF Cleanup"
				module = "dnf"

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

	actual, err := w.Run(workflow.NewWorkflowContext(ui.MockUI, inv, false))
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	expected := integration.ExpectedWorkflowOutput{
		Processes: []integration.ExpectedProcessOutput{
			{
				Steps: map[string]integration.ExpectedStepOutput{
					"dnf_present_what_if": {
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
					"dnf_present": {
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
					"dnf_present_repeat_what_if": {
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
					"dnf_present_repeat": {
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
					"dnf_absent_what_if": {
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
					"dnf_absent": {
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
					"dnf_absent_repeat_what_if": {
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
					"dnf_absent_repeat": {
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
					"dnf_latest_what_if": {
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
					"dnf_latest": {
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
					"dnf_latest_repeat_what_if": {
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
					"dnf_latest_repeat": {
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
					"dnf_cleanup": {
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
