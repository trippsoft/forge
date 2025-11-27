// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/hclfunction"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/ui"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

var (
	cachedPrivateKeyFiles = map[string][]byte{}
)

func resolveIntermediate(intermediate *intermediateInventory) (*Inventory, hcl.Diagnostics) {
	diags := validateIntermediateInventory(intermediate)
	if diags.HasErrors() {
		return nil, diags
	}

	resolveGroupMemberships(intermediate)

	hostVars, moreDiags := resolveAllHostVars(intermediate)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	hostTransports, moreDiags := resolveAllHostTransports(intermediate, hostVars)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	hostEscalateConfigs, moreDiags := resolveAllHostEscalateConfigs(intermediate, hostVars)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	inventory, moreDiags := buildFinalInventory(intermediate, hostVars, hostTransports, hostEscalateConfigs)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	return inventory, diags
}

func validateIntermediateInventory(intermediate *intermediateInventory) hcl.Diagnostics {
	diags := hcl.Diagnostics{}

	for groupName, group := range intermediate.groups {
		if groupName == "" {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group name",
				Detail:   "The group name cannot be empty.",
				Subject:  group.hclRange,
			})
			continue
		}

		if groupName == "all" {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group name",
				Detail:   "The group name 'all' is reserved and cannot be used.",
				Subject:  group.hclRange,
			})
			continue
		}

		if group.parent == "" || group.parent == "all" {
			continue // Skip circular reference check if no parent group
		}

		if _, exists := intermediate.groups[group.parent]; !exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid parent group",
				Detail:   fmt.Sprintf("The parent group '%s' does not exist.", group.parent),
				Subject:  group.hclRange,
			})
			continue
		}

		if hasCircularReference(groupName, intermediate, util.NewSet[string]()) {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Circular group reference",
				Detail:   fmt.Sprintf("The group '%s' has a circular reference.", groupName),
				Subject:  group.hclRange,
			})
		}
	}

	moreDiags := validateNoNameConflicts(intermediate)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return diags
	}

	return diags
}

func hasCircularReference(groupName string, intermediate *intermediateInventory, visited *util.Set[string]) bool {
	if visited.Contains(groupName) {
		return true // Circular reference detected
	}

	visited.Add(groupName)

	group, exists := intermediate.groups[groupName]
	if !exists || group.parent == "" {
		return false // No parent group or group does not exist
	}

	if hasCircularReference(group.parent, intermediate, visited) {
		return true // Circular reference found in parent chain
	}

	visited.Remove(groupName)
	return false
}

func validateNoNameConflicts(intermediate *intermediateInventory) hcl.Diagnostics {
	diags := hcl.Diagnostics{}

	for groupName, group := range intermediate.groups {
		if host, exists := intermediate.hosts[groupName]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Name conflict",
				Detail:   fmt.Sprintf("The group name %q conflicts with a host name.", groupName),
				Subject:  group.hclRange,
			})

			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Name conflict",
				Detail: fmt.Sprintf(
					"The group name %q conflicts with a host name defined at %q.",
					groupName,
					host.hclRange,
				),
				Subject: host.hclRange,
			})
		}
	}

	return diags
}

func resolveGroupMemberships(intermediate *intermediateInventory) {
	for _, host := range intermediate.hosts {
		resolveHostGroupMembershipWithBFS(host, intermediate)
	}
}

func resolveHostGroupMembershipWithBFS(host *intermediateHost, intermediate *intermediateInventory) {
	queue := []string{}
	queue = append(queue, host.groups...)

	visited := util.NewSet[string]()
	for len(queue) > 0 {
		currentGroup := queue[0]
		queue = queue[1:]

		if visited.Contains(currentGroup) {
			continue // Skip if already visited
		}

		visited.Add(currentGroup)
		host.allGroups = append(host.allGroups, currentGroup)

		group, exists := intermediate.groups[currentGroup]
		if !exists {
			continue // Skip if group does not exist, this should not happen
		}

		if group.parent != "" {
			if !visited.Contains(group.parent) {
				queue = append(queue, group.parent)
			}
		}
	}
}

func resolveAllHostVars(intermediate *intermediateInventory) (map[string]map[string]cty.Value, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}
	vars := make(map[string]map[string]cty.Value)
	for hostName, host := range intermediate.hosts {
		hostVars, moreDiags := resolveHostVars(hostName, host, intermediate)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip on errors
		}

		vars[hostName] = hostVars
	}

	return vars, diags
}

func resolveHostVars(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
) (map[string]cty.Value, hcl.Diagnostics) {

	inheritanceChain, diags := buildVarInheritanceChain(hostName, host, intermediate)
	if diags.HasErrors() {
		return nil, diags // Return on errors
	}

	combinedVars := combineVarsFromChain(inheritanceChain)

	hostVars, moreDiags := evaluateVarsIteratively(combinedVars)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags // Return on errors
	}

	return hostVars, diags
}

func buildVarInheritanceChain(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
) ([]map[string]*hcl.Attribute, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	chain := []map[string]*hcl.Attribute{}

	if len(host.vars) > 0 {
		chain = append(chain, host.vars)
	}

	for _, groupName := range host.allGroups {
		group, exists := intermediate.groups[groupName]
		if !exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group reference",
				Detail:   fmt.Sprintf("The group '%s' referenced by host '%s' does not exist.", groupName, hostName),
				Subject:  host.hclRange,
			})
			continue
		}

		if len(group.vars) > 0 {
			chain = append(chain, group.vars)
		}
	}

	if len(intermediate.vars) > 0 {
		chain = append(chain, intermediate.vars)
	}

	return chain, diags
}

func combineVarsFromChain(chain []map[string]*hcl.Attribute) map[string]*hcl.Attribute {
	combined := make(map[string]*hcl.Attribute)
	for _, attrs := range chain {
		for name, attr := range attrs {
			if _, exists := combined[name]; !exists {
				combined[name] = attr
			}
		}
	}

	return combined
}

func evaluateVarsIteratively(vars map[string]*hcl.Attribute) (map[string]cty.Value, hcl.Diagnostics) {
	diags := hcl.Diagnostics{}
	evaluatedVars := make(map[string]cty.Value)
	pendingVars := maps.Clone(vars)
	newlyEvaluated := util.NewSet[string]()

	max := len(vars) + 10 // Arbitrary buffer to ensure all vars can be evaluated
	i := 0

	for len(pendingVars) > 0 && i < max {
		i++
		progressMade := false
		for name, attr := range pendingVars {
			if attr == nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing expression",
					Detail:   fmt.Sprintf("The variable '%s' has no expression to evaluate.", name),
				})
				continue
			}

			if attr.Expr == nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Missing expression",
					Detail:   fmt.Sprintf("The variable '%s' has no expression to evaluate.", name),
					Subject:  &attr.Range,
				})
				continue
			}

			evalCtx := &hcl.EvalContext{
				Functions: hclfunction.HCLFunctions(),
			}
			evalCtx.Variables = map[string]cty.Value{
				"var": cty.ObjectVal(evaluatedVars),
			}

			value, moreDiags := attr.Expr.Value(evalCtx)
			if !moreDiags.HasErrors() {
				diags = diags.Extend(moreDiags)
				evaluatedVars[name] = value
				newlyEvaluated.Add(name)
				progressMade = true
				continue
			}

			if !containsDependencyError(moreDiags) || i >= max-1 {
				diags = diags.Extend(moreDiags)
			}
		}

		for _, varName := range newlyEvaluated.Items() {
			delete(pendingVars, varName)
		}

		newlyEvaluated.Clear()

		if !progressMade && len(pendingVars) > 0 {
			for name, attr := range pendingVars {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Unresolvable variable",
					Detail: fmt.Sprintf(
						"The variable '%s' could not be resolved due to missing or circular dependencies.",
						name,
					),
					Subject: &attr.Range,
				})
			}
			break
		}
	}

	if len(pendingVars) > 0 && i >= max {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Variable resolution limit reached",
			Detail: "The variable resolution process reached the maximum number of iterations. " +
				"This may indicate circular dependencies.",
			Subject: nil,
		})
	}

	return evaluatedVars, diags
}

func containsDependencyError(diags hcl.Diagnostics) bool {
	for _, diag := range diags {
		if strings.Contains(diag.Detail, "Unknown variable") ||
			strings.Contains(diag.Detail, "There is no variable named") ||
			strings.Contains(diag.Summary, "Unknown variable") ||
			strings.Contains(diag.Summary, "Unsupported attribute") {
			return true
		}
	}

	return false
}

func resolveAllHostTransports(
	intermediate *intermediateInventory,
	hostVars map[string]map[string]cty.Value,
) (map[string]transport.Transport, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	transports := make(map[string]transport.Transport)
	for hostName, host := range intermediate.hosts {
		vars, exists := hostVars[hostName]
		if !exists {
			vars = make(map[string]cty.Value)
		}

		hostTransport, moreDiags := resolveHostTransport(hostName, host, intermediate, vars)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip on errors
		}

		transports[hostName] = hostTransport
	}

	return transports, diags
}

func resolveHostTransport(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
	vars map[string]cty.Value,
) (transport.Transport, hcl.Diagnostics) {

	inheritanceChain, diags := buildTransportInheritanceChain(hostName, host, intermediate)
	if diags.HasErrors() {
		return nil, diags // Return on errors
	}

	combinedTransport := combineTransportsFromChain(inheritanceChain)
	if combinedTransport == nil {
		return transport.LocalTransport, diags
	}

	transport, moreDiags := createTransportFromConfig(combinedTransport, vars)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags // Return on errors
	}

	return transport, diags
}

func buildTransportInheritanceChain(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
) ([]*intermediateTransport, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	chain := []*intermediateTransport{}
	if host.transport != nil {
		chain = append(chain, host.transport)
	}

	for _, groupName := range host.allGroups {
		group, exists := intermediate.groups[groupName]
		if !exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group reference",
				Detail:   fmt.Sprintf("The group '%s' referenced by host '%s' does not exist.", groupName, hostName),
				Subject:  host.hclRange,
			})
			continue
		}

		if group.transport != nil {
			chain = append(chain, group.transport)
		}
	}

	if intermediate.transport != nil {
		chain = append(chain, intermediate.transport)
	}

	return chain, diags
}

func combineTransportsFromChain(chain []*intermediateTransport) *intermediateTransport {
	if len(chain) == 0 {
		return nil
	}

	combined := &intermediateTransport{
		name:     "",
		config:   make(map[string]*hcl.Attribute),
		hclRange: nil,
	}

	for _, transport := range chain {
		if combined.name == "" {
			combined.name = transport.name
		} else if combined.name != transport.name {
			continue // Skip conflicting transport types
		}

		for key, attr := range transport.config {
			if _, exists := combined.config[key]; !exists {
				combined.config[key] = attr
			}
		}
	}

	return combined
}

func createTransportFromConfig(
	intermediate *intermediateTransport,
	vars map[string]cty.Value,
) (transport.Transport, hcl.Diagnostics) {

	switch intermediate.name {
	case string(transport.TransportTypeLocal):
		return transport.LocalTransport, hcl.Diagnostics{}
	case string(transport.TransportTypeSSH):
		return createSSHTransport(intermediate.config, vars)
	default:
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unknown transport type",
				Detail:   fmt.Sprintf("The transport type '%s' is not recognized.", intermediate.name),
			},
		}
	}
}

func createSSHTransport(
	transportSSH map[string]*hcl.Attribute,
	vars map[string]cty.Value,
) (transport.Transport, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	var moreDiags hcl.Diagnostics

	var host string
	port := transport.DefaultSSHPort // Default SSH port
	var user string
	var privateKeyPath string
	var privateKeyPass string
	var password string
	useKnownHosts := transport.DefaultUseKnownHostsFile

	knownHostsPath, err := transport.DefaultKnownHostsPath()
	if err != nil {
		knownHostsPath = "" // Fallback to empty if default path cannot be determined
	}

	addUnknownHosts := transport.DefaultAddUnknownHostsToFile
	connectionTimeout := transport.DefaultSSHConnectionTimeout

	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(vars),
		},
		Functions: hclfunction.HCLFunctions(),
	}

	if attr, exists := transportSSH["host"]; exists && attr != nil {
		host, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the host
		}
	}

	if host == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing SSH host",
			Detail:   "The 'host' attribute is required for SSH transport.",
		})

		return nil, diags // Return early if the host is missing
	}

	if attr, exists := transportSSH["port"]; exists && attr != nil {
		port, moreDiags = util.ConvertHCLAttributeToUint16(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the port
		}
	}

	if attr, exists := transportSSH["user"]; exists && attr != nil {
		user, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the user
		}
	}

	if user == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing SSH user",
			Detail:   "The 'user' attribute is required for SSH transport.",
		})

		return nil, diags // Return early if the user is missing
	}

	if attr, exists := transportSSH["private_key_path"]; exists && attr != nil {
		privateKeyPath, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the private key path
		}
	}

	if attr, exists := transportSSH["private_key_pass"]; exists && attr != nil {
		privateKeyPass, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the private key pass
		}
	}

	if attr, exists := transportSSH["password"]; exists && attr != nil {
		password, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the password
		}
	}

	if attr, exists := transportSSH["use_known_hosts"]; exists && attr != nil {
		useKnownHosts, moreDiags = util.ConvertHCLAttributeToBool(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the use_known_hosts
		}
	}

	if attr, exists := transportSSH["known_hosts_path"]; exists && attr != nil {
		knownHostsPath, moreDiags = util.ConvertHCLAttributeToString(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the known hosts path
		}
	}

	if attr, exists := transportSSH["add_unknown_hosts"]; exists && attr != nil {
		addUnknownHosts, moreDiags = util.ConvertHCLAttributeToBool(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the add_unknown_hosts
		}
	}

	if attr, exists := transportSSH["connection_timeout"]; exists && attr != nil {
		connectionTimeout, moreDiags = util.ConvertHCLAttributeToDuration(attr, evalCtx)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			return nil, diags // Return early if there are errors in converting the connection_timeout
		}
	}

	builder, err := transport.NewSSHBuilder()
	if err != nil {
		return nil, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to create SSH transport builder",
			Detail:   fmt.Sprintf("An error occurred while creating the SSH transport builder: %v", err),
		})
	}

	builder = builder.WithHost(host).
		WithPort(port).
		WithUser(user).
		WithConnectionTimeout(connectionTimeout)

	if privateKeyPath != "" {
		privateKey, exists := cachedPrivateKeyFiles[privateKeyPath]
		if !exists {
			privateKey, err = os.ReadFile(privateKeyPath)
			if err != nil {
				return nil, append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Failed to read private key file",
					Detail: fmt.Sprintf(
						"An error occurred while reading the private key file '%s': %v",
						privateKeyPath,
						err,
					),
				})
			}
			cachedPrivateKeyFiles[privateKeyPath] = privateKey // Cache the private key
		}
		if privateKeyPass != "" {
			ui.SecretFilter.AddSecret(privateKeyPass)
			builder = builder.WithPublicKeyAuthWithPass(privateKey, privateKeyPass)
		} else {
			builder = builder.WithPublicKeyAuth(privateKey)
		}
	}

	if password != "" {
		ui.SecretFilter.AddSecret(password)
		builder = builder.WithPasswordAuth(password)
	}

	if useKnownHosts && addUnknownHosts {
		builder = builder.UseKnownHosts(knownHostsPath)
	} else if useKnownHosts {
		builder = builder.UseStrictKnownHosts(knownHostsPath)
	} else {
		builder = builder.DontUseKnownHosts()
	}

	sshTransport, err := builder.Build()
	if err != nil {
		return nil, append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to build SSH transport",
			Detail:   fmt.Sprintf("An error occurred while building the SSH transport: %v", err),
		})
	}

	return sshTransport, diags
}

func resolveAllHostEscalateConfigs(
	intermediate *intermediateInventory,
	hostVars map[string]map[string]cty.Value,
) (map[string]*EscalateConfig, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	escalateConfigs := make(map[string]*EscalateConfig)
	for hostName, host := range intermediate.hosts {
		vars, exists := hostVars[hostName]
		if !exists {
			vars = make(map[string]cty.Value)
		}

		hostEscalateConfig, moreDiags := resolveHostEscalateConfig(hostName, host, intermediate, vars)
		diags = diags.Extend(moreDiags)
		if moreDiags.HasErrors() {
			continue // Skip on errors
		}

		escalateConfigs[hostName] = hostEscalateConfig
	}

	return escalateConfigs, diags
}

func resolveHostEscalateConfig(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
	vars map[string]cty.Value,
) (*EscalateConfig, hcl.Diagnostics) {

	inheritanceChain, diags := buildEscalateConfigInheritanceChain(hostName, host, intermediate)
	if diags.HasErrors() {
		return nil, diags // Return on errors
	}

	combinedEscalate := combineEscalateConfigsFromChain(inheritanceChain)
	if combinedEscalate == nil {
		return &EscalateConfig{}, diags
	}

	escalateConfig, moreDiags := createEscalateConfigFromCombined(combinedEscalate, vars)
	diags = diags.Extend(moreDiags)
	if moreDiags.HasErrors() {
		return nil, diags // Return on errors
	}

	return escalateConfig, diags
}

func buildEscalateConfigInheritanceChain(
	hostName string,
	host *intermediateHost,
	intermediate *intermediateInventory,
) ([]*intermediateEscalate, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	chain := []*intermediateEscalate{}
	if host.escalate != nil {
		chain = append(chain, host.escalate)
	}

	for _, groupName := range host.allGroups {
		group, exists := intermediate.groups[groupName]
		if !exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid group reference",
				Detail:   fmt.Sprintf("The group '%s' referenced by host '%s' does not exist.", groupName, hostName),
				Subject:  host.hclRange,
			})
			continue
		}

		if group.escalate != nil {
			chain = append(chain, group.escalate)
		}
	}

	if intermediate.escalate != nil {
		chain = append(chain, intermediate.escalate)
	}

	return chain, diags
}

func combineEscalateConfigsFromChain(inheritanceChain []*intermediateEscalate) *intermediateEscalate {
	if len(inheritanceChain) == 0 {
		return nil
	}

	combined := &intermediateEscalate{}
	for _, escalate := range inheritanceChain {
		if combined.password == nil {
			combined.password = escalate.password
			break
		}
	}

	return combined
}

func createEscalateConfigFromCombined(
	combinedEscalate *intermediateEscalate,
	vars map[string]cty.Value,
) (*EscalateConfig, hcl.Diagnostics) {

	if combinedEscalate.password == nil {
		return nil, hcl.Diagnostics{}
	}

	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(vars),
		},
		Functions: hclfunction.HCLFunctions(),
	}

	password := combinedEscalate.password
	value, diags := password.Expr.Value(evalCtx)
	if diags.HasErrors() {
		return nil, diags // Return on errors
	}

	if !value.IsKnown() || value.IsNull() {
		return nil, diags // Return if value is unknown or null
	}

	if value.Type() != cty.String {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid escalate password type",
			Detail:   "The escalate password must be a string.",
			Subject:  &password.Range,
		})

		return nil, diags // Return on type error
	}

	escalateConfig := NewEscalateConfig(value.AsString())

	return escalateConfig, diags
}

func buildFinalInventory(
	intermediate *intermediateInventory,
	hostVars map[string]map[string]cty.Value,
	hostTransports map[string]transport.Transport,
	hostEscalateConfigs map[string]*EscalateConfig,
) (*Inventory, hcl.Diagnostics) {

	diags := hcl.Diagnostics{}
	inventory := NewInventory(map[string]*Host{}, map[string][]*Host{}, map[string][]*Host{})
	inventory.targets["all"] = make([]*Host, 0, len(intermediate.hosts))
	for hostName, intermediateHost := range intermediate.hosts {
		vars, exists := hostVars[hostName]
		if !exists {
			vars = make(map[string]cty.Value)
		}

		t, exists := hostTransports[hostName]
		if !exists {
			t = transport.LocalTransport
		}

		escalateConfig, exists := hostEscalateConfigs[hostName]
		if !exists {
			escalateConfig = NewEscalateConfig("")
		}

		if exists && escalateConfig.Pass() != "" {
			ui.SecretFilter.AddSecret(escalateConfig.Pass())
		}

		builder := NewHostBuilder().
			WithName(hostName).
			WithTransport(t).
			WithEscalateConfig(escalateConfig).
			WithVars(vars)

		host, err := builder.Build()
		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to build host",
				Detail:   fmt.Sprintf("An error occurred while building the host '%s': %v", hostName, err),
			})
			continue
		}

		inventory.hosts[hostName] = host
		inventory.targets[hostName] = []*Host{host}
		inventory.targets["all"] = append(inventory.targets["all"], host)

		for _, groupName := range intermediateHost.allGroups {
			if _, exists := inventory.groups[groupName]; !exists {
				inventory.groups[groupName] = make([]*Host, 0)
			}

			if _, exists := inventory.targets[groupName]; !exists {
				inventory.targets[groupName] = make([]*Host, 0)
			}

			inventory.groups[groupName] = append(inventory.groups[groupName], host)
			inventory.targets[groupName] = append(inventory.targets[groupName], host)
		}
	}

	return inventory, diags
}
