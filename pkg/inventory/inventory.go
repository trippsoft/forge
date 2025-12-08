// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import (
	"errors"
	"maps"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/result"
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

	stepContexts    []map[string]cty.Value
	procedureInputs []map[string]cty.Value
}

// Name returns the name of the host.
func (h *Host) Name() string {
	return h.name
}

// Transport returns the transport used to connect to the host.
func (h *Host) Transport() transport.Transport {
	return h.transport
}

// EscalateConfig returns the escalation configuration for the host, if any.
func (h *Host) EscalateConfig() *EscalateConfig {
	return h.escalateConfig
}

// Info returns the host's information as a HostInfo instance.
func (h *Host) Info() *info.HostInfo {
	return h.info
}

// Vars returns the variables associated with the host.
func (h *Host) Vars() map[string]cty.Value {
	return h.vars
}

// PopulateInfo populates the host's information using its transport.
func (h *Host) PopulateInfo() *result.Result {
	return h.info.Populate(h.transport)
}

// GetCurrentContextSteps retrieves the current step context for the host.
func (h *Host) GetCurrentContextSteps() (map[string]cty.Value, error) {
	if len(h.stepContexts) == 0 {
		return nil, errors.New("no step context available")
	}

	return h.stepContexts[len(h.stepContexts)-1], nil
}

// StoreStepOutput stores a step in the current step context.
//
// If the key already exists, it will overwrite the existing value.
// This is by design to allow for step updates.
func (h *Host) StoreStepOutput(key string, value cty.Value) error {
	stepContext, err := h.GetCurrentContextSteps()
	if err != nil {
		return err
	}

	stepContext[key] = value // Overwrites existing keys by design
	return nil
}

// StartProcedure initializes a new procedure context for the host.
func (h *Host) StartProcedure(inputs map[string]cty.Value) {
	h.stepContexts = append(h.stepContexts, make(map[string]cty.Value))
	h.procedureInputs = append(h.procedureInputs, inputs)
}

// EndProcedure ends the current procedure context for the host.
func (h *Host) EndProcedure() error {
	if len(h.stepContexts) < 2 {
		return errors.New("no step context to end")
	}
	h.stepContexts = h.stepContexts[:len(h.stepContexts)-1]
	if len(h.procedureInputs) < 1 {
		return errors.New("no procedure inputs to end")
	}
	h.procedureInputs = h.procedureInputs[:len(h.procedureInputs)-1]
	return nil
}

// ClearSteps clears all steps and procedure inputs for the host.
func (i *Host) ClearSteps() {
	i.stepContexts = []map[string]cty.Value{make(map[string]cty.Value)}
	i.procedureInputs = []map[string]cty.Value{}
}

// GetCurrentProcedureInputs retrieves the inputs for the current procedure context.
func (h *Host) GetCurrentProcedureInputs() map[string]cty.Value {
	if len(h.procedureInputs) == 0 {
		return nil
	}

	return h.procedureInputs[len(h.procedureInputs)-1]
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
		info:           info.NewHostInfo(),
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

// ClearSteps clears all steps and procedure inputs for the inventory.
func (i *Inventory) ClearSteps() {
	for _, host := range i.hosts {
		host.ClearSteps()
	}
}
