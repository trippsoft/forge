package inventory

import (
	"errors"

	"github.com/trippsoft/forge/internal/info"
	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type Host struct {
	name string

	transport transport.Transport

	taskContexts []map[string]cty.Value
	info         *info.HostInfo
	vars         map[string]cty.Value
}

func newHost(name string, transport transport.Transport, vars map[string]cty.Value) *Host {
	return &Host{
		name:         name,
		transport:    transport,
		taskContexts: []map[string]cty.Value{make(map[string]cty.Value)},
		info:         info.NewHostInfo(),
		vars:         vars,
	}
}

func (h *Host) Name() string {
	return h.name
}

func (h *Host) Transport() transport.Transport {
	return h.transport
}

func (h *Host) Info() *info.HostInfo {
	return h.info
}

func (h *Host) Vars() map[string]cty.Value {
	return h.vars
}

func (h *Host) StoreTask(key string, value cty.Value) error {

	taskContext, err := h.getCurrentContextTasks()
	if err != nil {
		return err
	}

	taskContext[key] = value // Overwrites existing keys by design
	return nil
}

func (h *Host) StartProcedure() {
	h.taskContexts = append(h.taskContexts, make(map[string]cty.Value))
}

func (h *Host) EndProcedure() error {
	if len(h.taskContexts) < 2 {
		return errors.New("no task context to end")
	}
	h.taskContexts = h.taskContexts[:len(h.taskContexts)-1]
	return nil
}

func (h *Host) getCurrentContextTasks() (map[string]cty.Value, error) {

	if len(h.taskContexts) == 0 {
		return nil, errors.New("no task context available")
	}

	return h.taskContexts[len(h.taskContexts)-1], nil
}

type Inventory struct {
	hosts map[string]*Host

	groups  map[string][]*Host
	targets map[string][]*Host
}

// Host retrieves a host by name from the inventory.
// It returns the host and a boolean indicating if the host exists.
func (i *Inventory) Host(name string) (*Host, bool) {

	host, exists := i.hosts[name]
	return host, exists
}

// Hosts returns all hosts in the inventory.
func (i *Inventory) Hosts() map[string]*Host {
	return i.hosts
}

// Group retrieves a group of hosts by name from the inventory.
// It returns the hosts in the group and a boolean indicating if the group exists.
// It does not include pseudo-groups like 'all' or hostnames.
func (i *Inventory) Group(name string) ([]*Host, bool) {
	hosts, exists := i.groups[name]
	return hosts, exists
}

// Groups returns all groups in the inventory.
// This does not include pseudo-groups like 'all' or hostnames.
func (i *Inventory) Groups() map[string][]*Host {
	return i.groups
}

// Target retrieves a target group of hosts by name from the inventory.
// It returns the hosts in the target group and a boolean indicating if the target exists.
// It includes the pseudo-group 'all' and hostnames as targets.
func (i *Inventory) Target(name string) ([]*Host, bool) {
	hosts, exists := i.targets[name]
	return hosts, exists
}

// Targets returns all target groups in the inventory.
// This includes the pseudo-group 'all' and hostnames as targets.
func (i *Inventory) Targets() map[string][]*Host {
	return i.targets
}
