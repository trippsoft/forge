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
)

var (
	varsBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "vars",
				LabelNames: []string{},
			},
		},
	}
	transportBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "transport",
				LabelNames: []string{"type"},
			},
		},
	}
	transportNoneSchema = &hcl.BodySchema{}
	transportSSHSchema  = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: sshAttrHost, Required: false},
			{Name: sshAttrPort, Required: false},
			{Name: sshAttrUser, Required: false},
			{Name: sshAttrPassword, Required: false},
			{Name: sshAttrPrivateKeyPath, Required: false},
			{Name: sshAttrPrivateKeyPass, Required: false},
			{Name: sshAttrUseKnownHosts, Required: false},
			{Name: sshAttrKnownHostsPath, Required: false},
			{Name: sshAttrAddUnknownHosts, Required: false},
			{Name: sshAttrConnectionTimeout, Required: false},
		},
		Blocks: []hcl.BlockHeaderSchema{},
	}
	escalateBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "escalate",
				LabelNames: []string{},
			},
		},
	}
	escalateAttributesSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "password", Required: false},
		},
	}
	groupBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "group",
				LabelNames: []string{"name"},
			},
		},
	}
	groupAttributesSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "parent", Required: false},
			{Name: "hosts", Required: false},
		},
	}
	hostBlockSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "host",
				LabelNames: []string{"name"},
			},
		},
	}
)
