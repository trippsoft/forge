// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// SELinuxInfo represents information about SELinux status on the host.
type SELinuxInfo struct {
	supported bool
	installed bool
	status    string
	typ       string
}

// Supported returns whether SELinux is supported on the host.
func (s *SELinuxInfo) Supported() bool {
	return s.supported
}

// Installed returns whether SELinux is installed on the host.
func (s *SELinuxInfo) Installed() bool {
	return s.installed
}

// Status returns the current SELinux status on the host.
func (s *SELinuxInfo) Status() string {
	return s.status
}

// Type returns the SELinux type on the host.
func (s *SELinuxInfo) Type() string {
	return s.typ
}

// ToMapOfCtyValues converts the SELinuxInfo into a map of cty.Values.
func (s *SELinuxInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !s.supported {
		return map[string]cty.Value{
			"selinux_installed": cty.NullVal(cty.String),
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	if !s.installed {
		return map[string]cty.Value{
			"selinux_installed": cty.False,
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"selinux_installed": cty.True,
		"selinux_status":    cty.StringVal(string(s.status)),
		"selinux_type":      cty.StringVal(string(s.typ)),
	}
}

// FromProtobuf populates the SELinuxInfo from a protobuf representation.
func (s *SELinuxInfo) FromProtobuf(other *SELinuxInfoPB) {
	s.supported = other.Supported
	s.installed = other.Installed
	s.status = other.Status
	s.typ = other.Type
}
