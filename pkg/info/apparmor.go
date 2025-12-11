// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// AppArmorInfo represents information about AppArmor status on the host.
type AppArmorInfo struct {
	supported bool
	enabled   bool
}

// Supported returns whether AppArmor is supported on the host.
func (a *AppArmorInfo) Supported() bool {
	return a.supported
}

// Enabled returns whether AppArmor is enabled on the host.
func (a *AppArmorInfo) Enabled() bool {
	return a.enabled
}

// ToMapOfCtyValues converts the AppArmorInfo into a map of cty.Values.
func (a *AppArmorInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !a.Supported() {
		return map[string]cty.Value{
			"apparmor_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"apparmor_enabled": cty.BoolVal(a.Enabled()),
	}
}

// FromProtobuf populates the AppArmorInfo from a protobuf representation.
func (a *AppArmorInfo) FromProtobuf(other *AppArmorInfoPB) {
	a.supported = other.Supported
	a.enabled = other.Enabled
}
