package inventory

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type inventoryFile struct {
	path    string
	content []byte
}

// DiscoverInventoryFiles retrieves all inventory files from the specified paths.
// It walks through each path, looking for files with the ".hcl" extension
// that are considered inventory files. It returns a slice of pointers to inventoryFile
// structs, each containing the file path and its content. If an error occurs during
// reading the files, it returns an error.
func DiscoverInventoryFiles(paths ...string) ([]*inventoryFile, error) {

	inventoryFiles := make([]*inventoryFile, 0, len(paths))

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
				return err
			}

			inventoryFiles = append(inventoryFiles, &inventoryFile{
				path:    path,
				content: content,
			})

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return inventoryFiles, nil
}

// ParseInventoryFiles parses the content of the inventory files and returns the
// parsed inventory.
func ParseInventoryFiles(files []*inventoryFile) (*Inventory, hcl.Diagnostics) {

	parser := hclparse.NewParser()
	diags := hcl.Diagnostics{}
	hclFiles := make([]*hcl.File, 0, len(files))

	for _, file := range files {

		hclFile, moreDiags := parser.ParseHCL(file.content, file.path)
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

	diags := hcl.Diagnostics{}
	intermediate := &intermediateInventory{
		allHosts: make(map[string]*intermediateHost),
	}

	varBlocks, body, moreDiags := body.PartialContent(varsBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	intermediate.vars, moreDiags = parseVarsBlocksToIntermediate(varBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	transportBlocks, body, moreDiags := body.PartialContent(transportBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	intermediate.transport, moreDiags = parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	groupBlocks, body, moreDiags := body.PartialContent(groupBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	intermediate.groups, moreDiags = parseGroupBlocksToIntermediate(groupBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	hostBlocks, body, moreDiags := body.PartialContent(hostBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	intermediate.hosts, moreDiags = parseHostBlocksToIntermediate(hostBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	if !body.MissingItemRange().Empty() {
		missingItemRange := body.MissingItemRange()
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unexpected content in inventory file",
			Detail:   "The inventory file contains unexpected content that was not parsed. This may indicate a syntax error or unsupported features.",
			Subject:  &missingItemRange,
		})
	}

	return intermediate, diags
}

func parseVarsBlocksToIntermediate(varBlocks *hcl.BodyContent) (map[string]*hcl.Attribute, hcl.Diagnostics) {

	vars := make(map[string]*hcl.Attribute)
	diags := hcl.Diagnostics{}

	for _, block := range varBlocks.Blocks {

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
					Detail:   fmt.Sprintf("Variable '%s' is defined multiple times in the inventory file. Each variable must have a unique name.", name),
					Subject:  &attr.Range,
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

	if block.Type != "vars" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid block type",
				Detail:   fmt.Sprintf("Expected 'vars' block, but found '%s'.", block.Type),
				Subject:  &block.DefRange,
			},
		}
	}

	vars := make(map[string]*hcl.Attribute)
	diags := hcl.Diagnostics{}

	attributes, moreDiags := block.Body.JustAttributes()
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	for _, attr := range attributes {

		if _, exists := vars[attr.Name]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate variable name",
				Detail:   fmt.Sprintf("Variable '%s' is defined multiple times in the vars block. Each variable must have a unique name.", attr.Name),
				Subject:  &attr.Range,
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

func parseTransportBlocksToIntermediate(transportBlocks *hcl.BodyContent) (*intermediateTransport, hcl.Diagnostics) {

	if len(transportBlocks.Blocks) == 0 {
		return nil, hcl.Diagnostics{}
	}

	diags := hcl.Diagnostics{}
	var transport *intermediateTransport

	for _, block := range transportBlocks.Blocks {

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
					Detail:   fmt.Sprintf("Found multiple 'transport' blocks with different types: '%s' and '%s'. Only one transport type is allowed.", transport.name, t.name),
					Subject:  t.hclRange,
				})
				continue
			}

			// Merge the configuration attributes
			for name, attr := range t.config {
				if _, exists := transport.config[name]; exists {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Duplicate transport configuration",
						Detail:   fmt.Sprintf("The attribute '%s' is defined multiple times in the 'transport' block. Each attribute must have a unique name.", name),
						Subject:  &attr.Range,
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

	if block.Type != "transport" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid block type",
				Detail:   fmt.Sprintf("Expected 'transport' block, but found '%s'.", block.Type),
				Subject:  &block.DefRange,
			},
		}
	}

	if len(block.Labels) != 1 {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid transport block",
				Detail:   "The 'transport' block must have exactly one label specifying the transport type.",
				Subject:  &block.DefRange,
			},
		}
	}

	transportType := block.Labels[0]

	diags := hcl.Diagnostics{}
	var body *hcl.BodyContent
	var remaining hcl.Body

	switch transportType {
	case string(transport.TransportTypeNone):
		body, remaining, _ = block.Body.PartialContent(transportNoneSchema)
	case string(transport.TransportTypeSSH):
		var moreDiags hcl.Diagnostics
		body, remaining, moreDiags = block.Body.PartialContent(transportSSHSchema)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			return nil, diags
		}
	case string(transport.TransportTypeWinRM):
		body, remaining, _ = block.Body.PartialContent(transportWinRMSchema) // TODO - implement WinRM schema
	default:
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid transport type",
				Detail: fmt.Sprintf(
					"The transport type '%s' is not supported. Allowed types are: %s, %s, %s",
					transportType,
					transport.TransportTypeNone,
					transport.TransportTypeSSH,
					transport.TransportTypeWinRM),
				Subject: &block.DefRange,
			},
		}
	}

	transport := &intermediateTransport{
		name:     transportType,
		config:   make(map[string]*hcl.Attribute),
		hclRange: &block.DefRange,
	}

	for _, attr := range body.Attributes {
		if _, exists := transport.config[attr.Name]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate transport configuration",
				Detail:   fmt.Sprintf("The attribute '%s' is defined multiple times in the 'transport' block. Each attribute must have a unique name.", attr.Name),
				Subject:  &attr.Range,
			})
			continue
		}

		transport.config[attr.Name] = attr
	}

	if !remaining.MissingItemRange().Empty() {
		missingItemRange := remaining.MissingItemRange()
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unexpected content in transport block",
			Detail:   "The 'transport' block contains unexpected content that was not parsed. This may indicate a syntax error or unsupported features.",
			Subject:  &missingItemRange,
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return transport, diags
}

func parseGroupBlocksToIntermediate(groupBlocks *hcl.BodyContent) (map[string]*intermediateGroup, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	groups := make(map[string]*intermediateGroup)

	for _, block := range groupBlocks.Blocks {

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
			// TODO - merge groups if they have the same name
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate group name",
				Detail:   fmt.Sprintf("Group '%s' is defined multiple times in the inventory file. Each group must have a unique name.", group.name),
				Subject:  group.hclRange,
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

	if block.Type != "group" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid block type",
				Detail:   fmt.Sprintf("Expected 'group' block, but found '%s'.", block.Type),
				Subject:  &block.DefRange,
			},
		}
	}

	if len(block.Labels) != 1 {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group block",
				Detail:   "The 'group' block must have exactly one label specifying the group name.",
				Subject:  &block.DefRange,
			},
		}
	}

	groupName := block.Labels[0]

	if groupName == "" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Empty group name",
				Detail:   "The group name cannot be empty.",
				Subject:  &block.DefRange,
			},
		}
	}

	group := &intermediateGroup{
		name:           groupName,
		childHosts:     make(map[string]*intermediateHost),
		hostReferences: []string{},
		vars:           make(map[string]*hcl.Attribute),
	}

	varsBlocks, body, moreDiags := block.Body.PartialContent(varsBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	group.vars, moreDiags = parseVarsBlocksToIntermediate(varsBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	transportBlocks, body, moreDiags := body.PartialContent(transportBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	group.transport, moreDiags = parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	childHostBlocks, body, moreDiags := body.PartialContent(hostBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	group.childHosts, moreDiags = parseHostBlocksToIntermediate(childHostBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	for _, host := range group.childHosts {
		host.parentGroup = group.name
	}

	attributes, body, moreDiags := body.PartialContent(groupAttributesSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	for _, attr := range attributes.Attributes {
		switch attr.Name {
		case "parent":
			if group.parent != "" {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate parent attribute",
					Detail:   fmt.Sprintf("The 'parent' attribute is defined multiple times in the 'group' block for group '%s'. Each group can have at most one parent.", group.name),
					Subject:  &attr.Range,
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
					Detail:   fmt.Sprintf("The 'parent' attribute in the 'group' block for group '%s' must be a string, but got '%s'.", group.name, value.Type().FriendlyName()),
					Subject:  &attr.Range,
				})
				continue // Skip this attribute if the type is incorrect
			}

			group.parent = value.AsString()

		case "hosts":
			if len(group.hostReferences) > 0 {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate hosts attribute",
					Detail:   fmt.Sprintf("The 'hosts' attribute is defined multiple times in the 'group' block for group '%s'. Each group can have at most one 'hosts' attribute.", group.name),
					Subject:  &attr.Range,
				})
				continue // Skip this attribute if it is duplicated
			}

			value, moreDiags := attr.Expr.Value(nil)
			diags = diags.Extend(moreDiags)
			if moreDiags.HasErrors() {
				continue // Skip this attribute if there are errors
			}

			if value.Type() != cty.List(cty.String) && value.Type() != cty.Set(cty.String) && !value.Type().IsTupleType() {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid hosts attribute type",
					Detail:   fmt.Sprintf("The 'hosts' attribute in the 'group' block for group '%s' must be a list of strings, but got '%s'.", group.name, value.Type().FriendlyName()),
					Subject:  &attr.Range,
				})
				continue // Skip this attribute if the type is incorrect
			}

			for _, elem := range value.AsValueSlice() {
				if elem.Type() != cty.String {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Invalid host reference type",
						Detail:   fmt.Sprintf("Each element in the 'hosts' list for group '%s' must be a string, but got '%s'.", group.name, elem.Type().FriendlyName()),
						Subject:  &attr.Range,
					})
					continue // Skip this element if the type is incorrect
				}

				hostName := elem.AsString()
				if hostName == "" {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Empty host reference",
						Detail:   fmt.Sprintf("The 'hosts' list for group '%s' contains an empty string. Each host reference must be a non-empty string.", group.name),
						Subject:  &attr.Range,
					})
					continue // Skip empty host references
				}

				group.hostReferences = append(group.hostReferences, hostName)
			}
		}
	}

	if !body.MissingItemRange().Empty() {
		missingItemRange := body.MissingItemRange()
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Unexpected content in group block",
			Detail:   "The 'group' block contains unexpected content that was not parsed. This may indicate a syntax error or unsupported features.",
			Subject:  &missingItemRange,
		})
	}

	if diags.HasErrors() {
		return nil, diags
	}

	return group, diags
}

func parseHostBlocksToIntermediate(hostBlocks *hcl.BodyContent) (map[string]*intermediateHost, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	hosts := make(map[string]*intermediateHost)

	for _, block := range hostBlocks.Blocks {

		host, moreDiags := parseHostBlockToIntermediate(block)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue
		}

		if _, exists := hosts[host.name]; exists {
			// TODO - merge hosts if they have the same name
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate host name",
				Detail:   fmt.Sprintf("Host '%s' is defined multiple times in the inventory file. Each host must have a unique name.", host.name),
				Subject:  host.hclRange,
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

	if block.Type != "host" {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid block type",
				Detail:   fmt.Sprintf("Expected 'host' block, but found '%s'.", block.Type),
				Subject:  &block.DefRange,
			},
		}
	}

	if len(block.Labels) != 1 {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid host block",
				Detail:   "The 'host' block must have exactly one label specifying the host name.",
				Subject:  &block.DefRange,
			},
		}
	}

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

	varsBlocks, body, moreDiags := block.Body.PartialContent(varsBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	host.vars, moreDiags = parseVarsBlocksToIntermediate(varsBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	transportBlocks, body, moreDiags := body.PartialContent(transportBlockSchema)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	host.transport, moreDiags = parseTransportBlocksToIntermediate(transportBlocks)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	return host, diags
}
