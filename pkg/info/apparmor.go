// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the AppArmorInfo into a map of cty.Values.
func (a *AppArmorInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !a.Supported {
		return map[string]cty.Value{
			"apparmor_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"apparmor_enabled": cty.BoolVal(a.Enabled),
	}
}

// From populates the AppArmorInfo from another AppArmorInfo.
func (a *AppArmorInfo) From(other *AppArmorInfo) {
	a.Supported = other.Supported
	a.Enabled = other.Enabled
}
