// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"fmt"

	"github.com/trippsoft/forge/pkg/discover"
	"github.com/zclconf/go-cty/cty"
)

// AppArmorInfo contains information about the AppArmor status of a managed host.
type AppArmorInfo struct {
	supported bool
	enabled   bool
}

// Supported indicates whether AppArmor is supported on the host.
func (a *AppArmorInfo) Supported() bool {
	return a.supported
}

// Enabled indicates whether AppArmor is enabled on the host.
func (a *AppArmorInfo) Enabled() bool {
	return a.enabled
}

// ToMapOfCtyValues converts the AppArmorInfo into a map of cty.Values.
func (a *AppArmorInfo) ToMapOfCtyValues() map[string]cty.Value {
	if !a.supported {
		return map[string]cty.Value{
			"apparmor_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"apparmor_enabled": cty.BoolVal(a.enabled),
	}
}

// FromProtobuf populates the AppArmorInfo from a protobuf representation.
func (a *AppArmorInfo) FromProtobuf(response *discover.AppArmorInfoResponse) {
	a.supported = response.Supported
	a.enabled = response.Enabled
}

// String returns a string representation of the AppArmor information.
//
// This is useful for logging or debugging purposes.
func (a *AppArmorInfo) String() string {
	if !a.supported {
		return "apparmor_enabled: not supported\n"
	}

	return fmt.Sprintf("apparmor_enabled: %t\n", a.enabled)
}
