// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import "github.com/hashicorp/hcl/v2"

type intermediateTransport struct {
	name   string
	config map[string]*hcl.Attribute

	hclRange *hcl.Range
}

type intermediateEscalate struct {
	password *hcl.Attribute
}

type intermediateHost struct {
	name string

	vars      map[string]*hcl.Attribute
	transport *intermediateTransport
	escalate  *intermediateEscalate

	groups    []string
	allGroups []string

	hclRange *hcl.Range
}

type intermediateGroup struct {
	name string

	parent string

	vars      map[string]*hcl.Attribute
	transport *intermediateTransport
	escalate  *intermediateEscalate

	hclRange *hcl.Range
}

type intermediateInventory struct {
	vars      map[string]*hcl.Attribute
	transport *intermediateTransport
	escalate  *intermediateEscalate

	groups map[string]*intermediateGroup
	hosts  map[string]*intermediateHost
}
