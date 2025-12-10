// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the SELinuxInfo into a map of cty.Values.
func (s *SELinuxInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !s.Supported {
		return map[string]cty.Value{
			"selinux_installed": cty.NullVal(cty.String),
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	if !s.Installed {
		return map[string]cty.Value{
			"selinux_installed": cty.False,
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"selinux_installed": cty.True,
		"selinux_status":    cty.StringVal(string(s.Status)),
		"selinux_type":      cty.StringVal(string(s.Type)),
	}
}

// From populates the SELinuxInfo from another SELinuxInfo.
func (s *SELinuxInfo) From(other *SELinuxInfo) {
	s.Supported = other.Supported
	s.Installed = other.Installed
	s.Status = other.Status
	s.Type = other.Type
}
