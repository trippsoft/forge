// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import (
	"errors"
	"maps"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

// EscalateConfig defines the configuration for privilege escalation.
type EscalateConfig struct {
	password string
}

// NewEscalateConfig creates a new EscalateConfig with the given password.
func NewEscalateConfig(password string) *EscalateConfig {
	return &EscalateConfig{
		password: password,
	}
}

// Pass returns the password used for privilege escalation on the host, if any.
func (e *EscalateConfig) Pass() string {
	return e.password
}

// Host represents a single host in the inventory.
type Host struct {
	name string

	transport      transport.Transport
	escalateConfig *EscalateConfig

	info *info.HostInfo
	vars map[string]cty.Value
}

// Name returns the name of the host.
func (h *Host) Name() string {
	return h.name
}

// Transport returns the transport used to connect to the host.
func (h *Host) Transport() transport.Transport {
	return h.transport
}

// Vars returns the variables associated with the host.
func (h *Host) Vars() map[string]cty.Value {
	return h.vars
}

// HostBuilder is used to build Host instances.
type HostBuilder struct {
	name string

	transport      transport.Transport
	escalateConfig *EscalateConfig

	vars map[string]cty.Value
}

// NewHostBuilder creates a new HostBuilder instance.
func NewHostBuilder() *HostBuilder {
	return &HostBuilder{}
}

// WithName sets the name of the host.
func (b *HostBuilder) WithName(name string) *HostBuilder {
	b.name = name
	return b
}

// WithTransport sets the transport for the host.
func (b *HostBuilder) WithTransport(transport transport.Transport) *HostBuilder {
	b.transport = transport
	return b
}

// WithEscalateConfig sets the escalate configuration for the host.
func (b *HostBuilder) WithEscalateConfig(escalateConfig *EscalateConfig) *HostBuilder {
	b.escalateConfig = escalateConfig
	return b
}

// WithVars sets the variables for the host.
func (b *HostBuilder) WithVars(vars map[string]cty.Value) *HostBuilder {
	b.vars = vars
	return b
}

// Build constructs the Host instance.
func (b *HostBuilder) Build() (*Host, error) {
	if b.name == "" {
		return nil, errors.New("host name is required")
	}

	if b.transport == nil {
		return nil, errors.New("host transport is required")
	}

	if b.vars == nil {
		b.vars = make(map[string]cty.Value, 0)
	}

	return &Host{
		name:           b.name,
		transport:      b.transport,
		escalateConfig: b.escalateConfig,
		info:           &info.HostInfo{},
		vars:           b.vars,
	}, nil
}

// Inventory represents a collection of hosts, groups, and targets.
type Inventory struct {
	hosts map[string]*Host

	groups  map[string][]*Host
	targets map[string][]*Host
}

func NewInventory(hosts map[string]*Host, groups map[string][]*Host, targets map[string][]*Host) *Inventory {
	return &Inventory{
		hosts:   hosts,
		groups:  groups,
		targets: targets,
	}
}

// Host retrieves a host by name from the inventory.
//
// It returns the host and a boolean indicating if the host exists.
func (i *Inventory) Host(name string) (*Host, bool) {
	host, exists := i.hosts[name]
	return host, exists
}

// Hosts returns a copy of all hosts in the inventory.
func (i *Inventory) Hosts() map[string]*Host {
	hosts := maps.Clone(i.hosts)
	return hosts
}

// Group retrieves a group of hosts by name from the inventory.
//
// It returns the hosts in the group and a boolean indicating if the group exists.
// It does not include pseudo-groups like 'all' or hostnames.
func (i *Inventory) Group(name string) ([]*Host, bool) {
	hosts, exists := i.groups[name]
	return hosts, exists
}

// Groups returns a copy of all groups in the inventory.
//
// This does not include pseudo-groups like 'all' or hostnames.
func (i *Inventory) Groups() map[string][]*Host {
	groups := maps.Clone(i.groups)
	return groups
}

// Target retrieves a target group of hosts by name from the inventory.
//
// It returns the hosts in the target group and a boolean indicating if the target exists.
// It includes the pseudo-group 'all' and hostnames as targets.
func (i *Inventory) Target(name string) ([]*Host, bool) {
	hosts, exists := i.targets[name]
	return hosts, exists
}

// Targets returns a copy of all target groups in the inventory.
//
// This includes the pseudo-group 'all' and hostnames as targets.
func (i *Inventory) Targets() map[string][]*Host {
	targets := maps.Clone(i.targets)
	return targets
}
