// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"maps"
	"time"

	"github.com/zclconf/go-cty/cty"
)

const (
	PopulateTimeout = 30 * time.Minute
)

func (h *HostInfo) Discover() []string {
	warnings := make([]string, 0)

	w := h.Os.discover()
	warnings = append(warnings, w...)

	w = h.Fips.discover()
	warnings = append(warnings, w...)

	w = h.AppArmor.discover()
	warnings = append(warnings, w...)

	w = h.Selinux.discover()
	warnings = append(warnings, w...)

	w = h.PackageManager.discover(h.Os)
	warnings = append(warnings, w...)

	w = h.ServiceManager.discover()
	warnings = append(warnings, w...)

	if len(warnings) == 0 {
		return nil
	}

	return warnings
}

// ToMapOfCtyValues converts the HostInfo into a map of cty.Values.
func (i *HostInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	maps.Copy(values, i.Os.ToMapOfCtyValues())
	maps.Copy(values, i.Fips.ToMapOfCtyValues())
	maps.Copy(values, i.AppArmor.ToMapOfCtyValues())
	maps.Copy(values, i.Selinux.ToMapOfCtyValues())
	maps.Copy(values, i.PackageManager.ToMapOfCtyValues())
	maps.Copy(values, i.ServiceManager.ToMapOfCtyValues())
	return values
}

// CopyFrom copies the HostInfo from another HostInfo.
func (h *HostInfo) CopyFrom(other *HostInfo) {
	h.Os.CopyFrom(other.Os)
	h.Fips.CopyFrom(other.Fips)
	h.AppArmor.CopyFrom(other.AppArmor)
	h.Selinux.CopyFrom(other.Selinux)
	h.PackageManager.CopyFrom(other.PackageManager)
	h.ServiceManager.CopyFrom(other.ServiceManager)
}

// NewHostInfo creates a new HostInfo instance.
func NewHostInfo() *HostInfo {
	return &HostInfo{
		Os:             NewOSInfo(),
		Fips:           &FIPSInfo{},
		AppArmor:       &AppArmorInfo{},
		Selinux:        &SELinuxInfo{},
		PackageManager: &PackageManagerInfo{},
		ServiceManager: &ServiceManagerInfo{},
	}
}

func DiscoverHostInfo() *DiscoverResponse {
	hostInfo := NewHostInfo()
	warnings := hostInfo.Discover()

	return &DiscoverResponse{
		HostInfo: hostInfo,
		Warnings: warnings,
	}
}
