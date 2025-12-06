// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package inventory

import "github.com/hashicorp/hcl/v2"

const (
	sshAttrHost              = "host"
	sshAttrPort              = "port"
	sshAttrUser              = "user"
	sshAttrPassword          = "password"
	sshAttrPrivateKeyPath    = "private_key_path"
	sshAttrPrivateKeyPass    = "private_key_pass"
	sshAttrUseKnownHosts     = "use_known_hosts"
	sshAttrKnownHostsPath    = "known_hosts_path"
	sshAttrAddUnknownHosts   = "add_unknown_hosts"
	sshAttrConnectionTimeout = "connection_timeout"
	sshAttrMinPluginPort     = "min_plugin_port"
	sshAttrMaxPluginPort     = "max_plugin_port"
	sshAttrTempPath          = "temp_path"
)

var (
	inventoryBodySchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "vars",
				LabelNames: []string{},
			},
			{
				Type:       "transport",
				LabelNames: []string{"type"},
			},
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
			{
				Type:       "group",
				LabelNames: []string{"name"},
			},
			{
				Type:       "host",
				LabelNames: []string{"name"},
			},
		},
		Attributes: []hcl.AttributeSchema{},
	}
	transportLocalSchema = &hcl.BodySchema{}
	transportSSHSchema   = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     sshAttrHost,
				Required: false,
			},
			{
				Name:     sshAttrPort,
				Required: false,
			},
			{
				Name:     sshAttrUser,
				Required: false,
			},
			{
				Name:     sshAttrPassword,
				Required: false,
			},
			{
				Name:     sshAttrPrivateKeyPath,
				Required: false,
			},
			{
				Name:     sshAttrPrivateKeyPass,
				Required: false,
			},
			{
				Name:     sshAttrUseKnownHosts,
				Required: false,
			},
			{
				Name:     sshAttrKnownHostsPath,
				Required: false,
			},
			{
				Name:     sshAttrAddUnknownHosts,
				Required: false,
			},
			{
				Name:     sshAttrConnectionTimeout,
				Required: false,
			},
			{
				Name:     sshAttrMinPluginPort,
				Required: false,
			},
			{
				Name:     sshAttrMaxPluginPort,
				Required: false,
			},
			{
				Name:     sshAttrTempPath,
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{},
	}
	escalateBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "password",
				Required: false,
			},
		},
	}
	groupBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "vars",
				LabelNames: []string{},
			},
			{
				Type:       "transport",
				LabelNames: []string{"type"},
			},
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
		},
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "parent",
				Required: false,
			},
		},
	}
	hostBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "vars",
				LabelNames: []string{},
			},
			{
				Type:       "transport",
				LabelNames: []string{"type"},
			},
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
		},
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "groups",
				Required: false,
			},
		},
	}
)
