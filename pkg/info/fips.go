// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// FIPSInfo contains information about the FIPS status of a managed host.
type FIPSInfo struct {
	known   bool
	enabled bool
}

// Known indicates whether the FIPS status is known.
func (f *FIPSInfo) Known() bool {
	return f.known
}

// Enabled indicates whether FIPS mode is enabled.
func (f *FIPSInfo) Enabled() bool {
	return f.enabled
}

// ToMapOfCtyValues converts the FIPSInfo into a map of cty.Values.
func (f *FIPSInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !f.known {
		return map[string]cty.Value{
			"fips_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"fips_enabled": cty.BoolVal(f.enabled),
	}
}

// FromProtobuf populates the FIPSInfo from a protobuf representation.
func (f *FIPSInfo) FromProtobuf(response *FIPSInfoResponse) {
	f.known = response.Known
	f.enabled = response.Enabled
}

// String returns a string representation of the FIPS information.
//
// This is useful for logging or debugging purposes.
func (f *FIPSInfo) String() string {
	if !f.known {
		return "fips_enabled: unknown on this OS"
	}

	return fmt.Sprintf("fips_enabled: %t\n", f.enabled)
}
