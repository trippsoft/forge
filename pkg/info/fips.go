// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// FIPSInfo represents information about FIPS status on the host.
type FIPSInfo struct {
	known   bool
	enabled bool
}

// Known returns whether FIPS status is known on the host.
func (f *FIPSInfo) Known() bool {
	return f.known
}

// Enabled returns whether FIPS is enabled on the host.
func (f *FIPSInfo) Enabled() bool {
	return f.enabled
}

// ToMapOfCtyValues converts the FIPSInfo into a map of cty.Values.
func (f *FIPSInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !f.Known() {
		return map[string]cty.Value{
			"fips_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"fips_enabled": cty.BoolVal(f.Enabled()),
	}
}

// FromProtobuf populates the FIPSInfo from a protobuf representation.
func (f *FIPSInfo) FromProtobuf(other *FIPSInfoPB) {
	f.known = other.Known
	f.enabled = other.Enabled
}
