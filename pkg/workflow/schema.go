// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package workflow

import "github.com/hashicorp/hcl/v2"

var (
	workflowBodySchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "process",
				LabelNames: []string{},
			},
		},
	}
	processBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "step",
				LabelNames: []string{"id"},
			},
			{
				Type:       "procedure",
				LabelNames: []string{"id"},
			},
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
		},
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "name",
				Required: true,
			},
			{
				Name:     "targets",
				Required: true,
			},
			{
				Name:     "exec_timeout",
				Required: false,
			},
			{
				Name:     "what_if",
				Required: false,
			},
			{
				Name:     "discover_info",
				Required: false,
			},
		},
	}
	stepBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
			{
				Type:       "loop",
				LabelNames: []string{},
			},
			{
				Type:       "input",
				LabelNames: []string{},
			},
			{
				Type:       "output",
				LabelNames: []string{},
			},
		},
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "name",
				Required: true,
			},
			{
				Name:     "module",
				Required: true,
			},
			{
				Name:     "condition",
				Required: false,
			},
			{
				Name:     "run_once",
				Required: false,
			},
			{
				Name:     "targets",
				Required: false,
			},
			{
				Name:     "exec_timeout",
				Required: false,
			},
			{
				Name:     "what_if",
				Required: false,
			},
		},
	}
	escalateBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "escalate",
				Required: false,
			},
			{
				Name:     "impersonate_user",
				Required: false,
			},
		},
	}
	loopBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "items",
				Required: true,
			},
			{
				Name:     "label",
				Required: false,
			},
			{
				Name:     "condition",
				Required: false,
			},
		},
	}
	outputBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "continue_on_fail",
				Required: false,
			},
			{
				Name:     "changed_condition",
				Required: false,
			},
			{
				Name:     "failed_condition",
				Required: false,
			},
		},
	}
)
