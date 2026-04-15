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

func TestSlurpRun_Linux(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Slurp"
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

			step "create_file" {
			    name = "Create File"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "echo 'Hello, World!' > /tmp/hello.txt"]
				}
			}

			step "slurp_what_if" {
			    name = "Slurp What If"
				module = "slurp"
				what_if = true

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "slurp" {
			    name = "Slurp"
				module = "slurp"

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "cleanup" {
				name = "Cleanup"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "rm -rf /tmp/hello.txt"]
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
					"create_file": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
						},
					},
					"slurp_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"slurp": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
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

func TestSlurpRun_Linux_SudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Slurp"
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

			step "create_file" {
			    name = "Create File"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "echo 'Hello, World!' > /tmp/hello.txt"]
				}
			}

			step "slurp_what_if" {
			    name = "Slurp What If"
				module = "slurp"
				what_if = true

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "slurp" {
			    name = "Slurp"
				module = "slurp"

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "cleanup" {
				name = "Cleanup"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "rm -rf /tmp/hello.txt"]
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
					"create_file": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
						},
					},
					"slurp_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"slurp": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204-pw": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204-pw": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
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

func TestSlurpRun_Linux_NoSudoPassword(t *testing.T) {
	workflowContent := `
	    process {
			name = "Test Slurp"
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

			step "create_file" {
			    name = "Create File"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "echo 'Hello, World!' > /tmp/hello.txt"]
				}
			}

			step "slurp_what_if" {
			    name = "Slurp What If"
				module = "slurp"
				what_if = true

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "slurp" {
			    name = "Slurp"
				module = "slurp"

				input {
				    path = "/tmp/hello.txt"
				}
			}

			step "cleanup" {
				name = "Cleanup"
				module = "command"

				input {
				    name = "/bin/sh"
					args = ["-c", "rm -rf /tmp/hello.txt"]
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
					"create_file": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
						},
					},
					"slurp_what_if": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"slurp": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"debian12": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora42": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"fedora41": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky10": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky9": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"rocky8": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2404": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
							"ubuntu2204": {
								Changed: false,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"content":     cty.StringVal("SGVsbG8sIFdvcmxkIQo="),
									"sha256_hash": cty.StringVal("c98c24b677eff44860afea6f493bbaec5bb1c4cbb209c6fc2bbb47f66ff2ad31"),
								},
							},
						},
					},
					"cleanup": {
						Hosts: map[string]integration.ExpectedHostOutput{
							"debian13": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"debian12": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora42": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"fedora41": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky10": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky9": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"rocky8": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2404": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
									"stderr": cty.StringVal(""),
								},
							},
							"ubuntu2204": {
								Changed: true,
								Failed:  false,
								Skipped: false,
								Output: map[string]cty.Value{
									"stdout": cty.StringVal(""),
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
