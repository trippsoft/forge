// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

// InventoryFile represents an inventory file with its path and content.
type InventoryFile struct {
	Path    string // Path returns the file path of the inventory file.
	Content []byte // Content contains the raw content of the inventory file.
}

// DiscoverInventoryFiles retrieves all inventory files from the specified paths.
//
// It walks through each path, looking for files with the ".hcl" extension that are considered inventory files.
// It returns a slice of pointers to InventoryFile structs, each containing the file path and its content.
// If an error occurs during reading the files, it returns an error.
func DiscoverInventoryFiles(paths ...string) ([]*InventoryFile, error) {
	inventoryFiles := make([]*InventoryFile, 0, len(paths))
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || filepath.Ext(path) != ".hcl" {
				return nil // Skip directories and non-HCL files
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read inventory file %q: %w", path, err)
			}

			inventoryFiles = append(inventoryFiles, &InventoryFile{
				Path:    path,
				Content: content,
			})

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return inventoryFiles, nil
}

// ParseInventoryFiles parses the content of the inventory files and returns the parsed inventory.
func ParseInventoryFiles(files []*InventoryFile) (*Inventory, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	diags := hcl.Diagnostics{}
	hclFiles := make([]*hcl.File, 0, len(files))
	for _, file := range files {
		hclFile, moreDiags := parser.ParseHCL(file.Content, file.Path)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip files with parsing errors
		}

		hclFiles = append(hclFiles, hclFile)
	}

	mergedBody := hcl.MergeFiles(hclFiles)

	inventory, moreDiags := parseHCLBody(mergedBody)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	return inventory, diags
}

func parseHCLBody(body hcl.Body) (*Inventory, hcl.Diagnostics) {
	intermediate, diags := parseHCLBodyToIntermediate(body)
	if diags.HasErrors() {
		return nil, diags
	}

	inventory, moreDiags := resolveIntermediate(intermediate)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	return inventory, diags
}

func parseHCLBodyToIntermediate(body hcl.Body) (*intermediateInventory, hcl.Diagnostics) {
	intermediate := &intermediateInventory{}
	diags := hcl.Diagnostics{}

	content, moreDiags := body.Content(inventoryBodySchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in an inventory file")
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	}

	varsBlocks := []*hcl.Block{}
	transportBlocks := []*hcl.Block{}
	escalateBlocks := []*hcl.Block{}
	groupBlocks := []*hcl.Block{}
	hostBlocks := []*hcl.Block{}

	for _, block := range content.Blocks {

		switch block.Type {
		case "vars":
			varsBlocks = append(varsBlocks, block)
		case "transport":
			transportBlocks = append(transportBlocks, block)
		case "escalate":
			escalateBlocks = append(escalateBlocks, block)
		case "group":
			groupBlocks = append(groupBlocks, block)
		case "host":
			hostBlocks = append(hostBlocks, block)
		}
	}

	vars, moreDiags := parseVarsBlocksToIntermediate(varsBlocks)
	diags = diags.Extend(moreDiags)

	if !moreDiags.HasErrors() {
		intermediate.vars = vars
	}

	transport, moreDiags := parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)

	if !moreDiags.HasErrors() {
		intermediate.transport = transport
	}

	escalate, moreDiags := parseEscalateBlocksToIntermediate(escalateBlocks)
	diags = diags.Extend(moreDiags)

	if !moreDiags.HasErrors() {
		intermediate.escalate = escalate
	}

	groups, moreDiags := parseGroupBlocksToIntermediate(groupBlocks)
	diags = diags.Extend(moreDiags)

	if !moreDiags.HasErrors() {
		intermediate.groups = groups
	}

	hosts, moreDiags := parseHostBlocksToIntermediate(hostBlocks)
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	} else {
		intermediate.hosts = hosts
	}

	for _, host := range intermediate.hosts {
		for _, groupName := range host.groups {
			if _, exists := intermediate.groups[groupName]; !exists {
				intermediate.groups[groupName] = &intermediateGroup{name: groupName} // Add groups referenced by hosts
			}
		}
	}

	return intermediate, diags
}

func parseVarsBlocksToIntermediate(blocks []*hcl.Block) (map[string]*hcl.Attribute, hcl.Diagnostics) {
	if len(blocks) == 0 {
		return map[string]*hcl.Attribute{}, hcl.Diagnostics{}
	}

	vars := make(map[string]*hcl.Attribute)
	diags := hcl.Diagnostics{}
	for _, block := range blocks {
		blockVars, moreDiags := parseVarsBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue
		}

		for name, attr := range blockVars {
			if _, exists := vars[name]; exists {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate variable name",
					Detail: fmt.Sprintf(
						"Variable '%s' is defined multiple times in the inventory file. "+
							"Each variable must have a unique name.",
						name,
					),
					Subject: &attr.Range,
				})
				continue
			}

			vars[name] = attr
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return vars, diags
}

func parseVarsBlockToIntermediate(block *hcl.Block) (map[string]*hcl.Attribute, hcl.Diagnostics) {
	vars := make(map[string]*hcl.Attribute)
	diags := hcl.Diagnostics{}

	attributes, moreDiags := block.Body.JustAttributes()
	util.ModifyUnexpectedElementDiags(moreDiags, "in a vars block")
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	for name, attr := range attributes {
		if _, exists := vars[name]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable name",
				Detail: fmt.Sprintf(
					"Variable '%s' is defined multiple times in the vars block. Each variable must have a unique name.",
					attr.Name,
				),
				Subject: &attr.Range,
			})
			continue
		}

		vars[attr.Name] = attr
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return vars, diags
}

func parseTransportBlocksToIntermediate(blocks []*hcl.Block) (*intermediateTransport, hcl.Diagnostics) {
	if len(blocks) == 0 {
		return nil, hcl.Diagnostics{}
	}

	diags := hcl.Diagnostics{}
	var transport *intermediateTransport
	for _, block := range blocks {
		t, moreDiags := parseTransportBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip blocks with errors
		}

		if transport == nil {
			transport = t
		} else {
			if transport.name != t.name {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Multiple transport blocks with different types",
					Detail: fmt.Sprintf(
						"Found multiple 'transport' blocks with different types: '%s' and '%s'. "+
							"Only one transport type is allowed.",
						transport.name,
						t.name,
					),
					Subject: t.hclRange,
				})
				continue
			}

			// Merge the configuration attributes
			for name, attr := range t.config {
				if _, exists := transport.config[name]; exists {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Duplicate transport configuration",
						Detail: fmt.Sprintf(
							"The attribute '%s' is defined multiple times in the 'transport' block. "+
								"Each attribute must have a unique name.",
							name,
						),
						Subject: &attr.Range,
					})
					continue
				}

				transport.config[name] = attr
			}
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return transport, diags
}

func parseTransportBlockToIntermediate(block *hcl.Block) (*intermediateTransport, hcl.Diagnostics) {
	transportType := block.Labels[0]

	diags := hcl.Diagnostics{}
	var body *hcl.BodyContent

	switch transportType {
	case string(transport.TransportTypeLocal):
		var moreDiags hcl.Diagnostics
		body, moreDiags = block.Body.Content(transportLocalSchema)
		util.ModifyUnexpectedElementDiags(moreDiags, "in a transport \"none\" block")
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			return nil, diags
		}

	case string(transport.TransportTypeSSH):
		var moreDiags hcl.Diagnostics
		body, moreDiags = block.Body.Content(transportSSHSchema)
		util.ModifyUnexpectedElementDiags(moreDiags, "in a transport \"ssh\" block")
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			return nil, diags
		}

	default:
		return nil, diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid transport type",
			Detail: fmt.Sprintf(
				"The transport type %q is not supported. Allowed types are: %q, %q",
				transportType,
				transport.TransportTypeLocal,
				transport.TransportTypeSSH),
			Subject: &block.DefRange,
		})
	}

	transport := &intermediateTransport{
		name:     transportType,
		config:   make(map[string]*hcl.Attribute),
		hclRange: &block.DefRange,
	}

	maps.Copy(transport.config, body.Attributes)

	if diags.HasErrors() {
		return nil, diags
	}

	return transport, diags
}

func parseEscalateBlocksToIntermediate(blocks []*hcl.Block) (*intermediateEscalate, hcl.Diagnostics) {
	if len(blocks) == 0 {
		return nil, hcl.Diagnostics{}
	}

	diags := hcl.Diagnostics{}
	var escalate *intermediateEscalate

	for _, block := range blocks {
		e, moreDiags := parseEscalateBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip blocks with errors
		}

		if escalate == nil {
			escalate = e
		} else {
			// Merge the configuration attributes
			if escalate.password != nil && e.password != nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate escalate configuration",
					Detail:   "The attribute 'password' is defined multiple times in the 'escalate' block.",
					Subject:  &e.password.Range,
				})
				continue
			}

			if escalate.password == nil && e.password != nil {
				escalate.password = e.password
			}
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return escalate, diags
}

func parseEscalateBlockToIntermediate(block *hcl.Block) (*intermediateEscalate, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}

	body, moreDiags := block.Body.Content(escalateBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in an escalate block")
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	escalate := &intermediateEscalate{}
	for _, attr := range body.Attributes {
		switch attr.Name {
		case "password":
			escalate.password = attr
		}
	}

	return escalate, diags
}

func parseGroupBlocksToIntermediate(blocks []*hcl.Block) (map[string]*intermediateGroup, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}
	groups := make(map[string]*intermediateGroup)
	for _, block := range blocks {
		group, moreDiags := parseGroupBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue
		}

		if group.name == "" {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group name",
				Detail:   "The group name cannot be empty.",
				Subject:  &block.DefRange,
			})
		}

		if _, exists := groups[group.name]; exists {
			// TODO - merge groups if they have the same name/context?
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate group name",
				Detail: fmt.Sprintf(
					"Group '%s' is defined multiple times in the inventory file. Each group must have a unique name.",
					group.name,
				),
				Subject: group.hclRange,
			})
			continue
		}

		groups[group.name] = group
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return groups, diags
}

func parseGroupBlockToIntermediate(block *hcl.Block) (*intermediateGroup, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}

	groupName := block.Labels[0]
	if groupName == "" {
		return nil, hcl.Diagnostics{&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Empty group name",
			Detail:   "The group name cannot be empty.",
			Subject:  &block.DefRange,
		}}
	}

	group := &intermediateGroup{
		name: groupName,
		vars: make(map[string]*hcl.Attribute),
	}

	content, moreDiags := block.Body.Content(groupBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a group block")
	diags = diags.Extend(moreDiags)

	if moreDiags.HasErrors() {
		return nil, diags
	}

	varsBlocks := []*hcl.Block{}
	transportBlocks := []*hcl.Block{}
	escalateBlocks := []*hcl.Block{}

	for _, block := range content.Blocks {
		switch block.Type {
		case "vars":
			varsBlocks = append(varsBlocks, block)
		case "transport":
			transportBlocks = append(transportBlocks, block)
		case "escalate":
			escalateBlocks = append(escalateBlocks, block)
		}
	}

	vars, moreDiags := parseVarsBlocksToIntermediate(varsBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		group.vars = vars
	}

	transport, moreDiags := parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		group.transport = transport
	}

	escalate, moreDiags := parseEscalateBlocksToIntermediate(escalateBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		group.escalate = escalate
	}

	for _, attr := range content.Attributes {
		switch attr.Name {
		case "parent":
			if group.parent != "" {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate parent attribute",
					Detail: fmt.Sprintf(
						"The 'parent' attribute is defined multiple times in the 'group' block for group '%s'. "+
							"Each group can have at most one parent.",
						group.name,
					),
					Subject: &attr.Range,
				})
				continue // Skip this attribute if it is duplicated
			}

			value, moreDiags := attr.Expr.Value(nil)
			diags = diags.Extend(moreDiags)
			if moreDiags.HasErrors() {
				continue // Skip this attribute if there are errors
			}

			if value.Type() != cty.String {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid parent attribute type",
					Detail: fmt.Sprintf(
						"The 'parent' attribute in the 'group' block for group '%s' must be a string, but got '%s'.",
						group.name,
						value.Type().FriendlyName(),
					),
					Subject: &attr.Range,
				})
				continue // Skip this attribute if the type is incorrect
			}

			group.parent = value.AsString()
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return group, diags
}

func parseHostBlocksToIntermediate(blocks []*hcl.Block) (map[string]*intermediateHost, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}
	hosts := make(map[string]*intermediateHost)

	for _, block := range blocks {
		host, moreDiags := parseHostBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)

		if moreDiags.HasErrors() {
			continue
		}

		if _, exists := hosts[host.name]; exists {
			// TODO - merge hosts if they have the same name/context?
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate host name",
				Detail: fmt.Sprintf(
					"Host '%s' is defined multiple times in the inventory file. Each host must have a unique name.",
					host.name,
				),
				Subject: host.hclRange,
			})
			continue
		}

		hosts[host.name] = host
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return hosts, diags
}

func parseHostBlockToIntermediate(block *hcl.Block) (*intermediateHost, hcl.Diagnostics) {
	hostName := block.Labels[0]
	if hostName == "" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Empty host name",
				Detail:   "The host name cannot be empty.",
				Subject:  &block.DefRange,
			},
		}
	}

	diags := hcl.Diagnostics{}
	host := &intermediateHost{
		name:      hostName,
		vars:      make(map[string]*hcl.Attribute),
		allGroups: []string{},
		hclRange:  &block.DefRange,
	}

	content, moreDiags := block.Body.Content(hostBlockSchema)
	util.ModifyUnexpectedElementDiags(moreDiags, "in a host block")
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	varsBlocks := []*hcl.Block{}
	transportBlocks := []*hcl.Block{}
	escalateBlocks := []*hcl.Block{}

	for _, block := range content.Blocks {
		switch block.Type {
		case "vars":
			varsBlocks = append(varsBlocks, block)
		case "transport":
			transportBlocks = append(transportBlocks, block)
		case "escalate":
			escalateBlocks = append(escalateBlocks, block)
		}
	}

	vars, moreDiags := parseVarsBlocksToIntermediate(varsBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		host.vars = vars
	}

	transport, moreDiags := parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		host.transport = transport
	}

	escalate, moreDiags := parseEscalateBlocksToIntermediate(escalateBlocks)
	diags = diags.Extend(moreDiags)
	if !moreDiags.HasErrors() {
		host.escalate = escalate
	}

	for name, attr := range content.Attributes {
		switch name {
		case "groups":
			value, moreDiags := attr.Expr.Value(nil)
			diags = diags.Extend(moreDiags)
			if moreDiags.HasErrors() {
				continue // Skip this attribute if there are errors
			}

			if value.Type() != cty.List(cty.String) &&
				value.Type() != cty.Set(cty.String) &&
				!value.Type().IsTupleType() {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid groups attribute type",
					Detail: fmt.Sprintf(
						`The "groups" attribute in the "host" block for host %q must be a list of strings, but got %q.`,
						host.name,
						value.Type().FriendlyName(),
					),
					Subject: &attr.Range,
				})
				continue // Skip this attribute if the type is incorrect
			}

			for _, elem := range value.AsValueSlice() {
				if elem.Type() != cty.String {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid group reference type",
						Detail: fmt.Sprintf(
							`Each element in the "groups" list for host %q must be a string, but got %q.`,
							host.name,
							elem.Type().FriendlyName(),
						),
						Subject: &attr.Range,
					})
					continue // Skip this element if the type is incorrect
				}

				hostName := elem.AsString()
				if hostName == "" {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Empty group reference",
						Detail: fmt.Sprintf(
							`The "groups" list for host %q contains an empty string. `+
								`Each group reference must be a non-empty string.`,
							host.name,
						),
						Subject: &attr.Range,
					})
					continue // Skip empty group references
				}

				host.groups = append(host.groups, hostName)
			}
		}
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return host, diags
}
