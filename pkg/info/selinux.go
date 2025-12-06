// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// SELinuxInfo contains information about the SELinux status of a managed host.
type SELinuxInfo struct {
	supported   bool
	installed   bool
	status      string
	selinuxType string
}

// Supported indicates whether SELinux is supported on the host.
func (s *SELinuxInfo) Supported() bool {
	return s.supported
}

// Installed indicates whether SELinux is installed on the host.
func (s *SELinuxInfo) Installed() bool {
	return s.installed
}

// Status returns the current SELinux status.
func (s *SELinuxInfo) Status() string {
	return s.status
}

// Type returns the SELinux type.
func (s *SELinuxInfo) Type() string {
	return s.selinuxType
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
		"selinux_type":      cty.StringVal(string(s.selinuxType)),
	}
}

// FromProtobuf populates the SELinuxInfo from a protobuf representation.
func (s *SELinuxInfo) FromProtobuf(response *SELinuxInfoResponse) {
	s.supported = response.Supported
	s.installed = response.Installed
	s.status = response.Status
	s.selinuxType = response.Type
}

// String returns a string representation of the SELinux information.
//
// This is useful for logging or debugging purposes.
func (s *SELinuxInfo) String() string {
	stringBuilder := &strings.Builder{}
	if !s.supported {
		stringBuilder.WriteString("selinux_installed: not supported on this OS\n")
		stringBuilder.WriteString("selinux_status: not supported on this OS\n")
		stringBuilder.WriteString("selinux_type: not supported on this OS\n")

		return stringBuilder.String()
	}

	if !s.installed {
		stringBuilder.WriteString("selinux_installed: false\n")
		stringBuilder.WriteString("selinux_status: not installed\n")
		stringBuilder.WriteString("selinux_type: not installed\n")

		return stringBuilder.String()
	}

	stringBuilder.WriteString("selinux_installed: true\n")
	stringBuilder.WriteString("selinux_status: ")
	stringBuilder.WriteString(string(s.status))
	stringBuilder.WriteString("\n")
	stringBuilder.WriteString("selinux_type: ")
	stringBuilder.WriteString(string(s.selinuxType))
	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}
