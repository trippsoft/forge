// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the FIPSInfo into a map of cty.Values.
func (f *FIPSInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !f.Known {
		return map[string]cty.Value{
			"fips_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"fips_enabled": cty.BoolVal(f.Enabled),
	}
}

// CopyFrom copies the FIPSInfo from another FIPSInfo.
func (f *FIPSInfo) CopyFrom(other *FIPSInfo) {
	f.Known = other.Known
	f.Enabled = other.Enabled
}
