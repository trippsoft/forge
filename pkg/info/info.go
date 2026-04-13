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

func (h *HostInfo) Discover() error {
	err := h.Os.discover()
	if err != nil {
		return err
	}

	err = h.Fips.discover()
	if err != nil {
		return err
	}

	err = h.AppArmor.discover()
	if err != nil {
		return err
	}

	err = h.Selinux.discover()
	if err != nil {
		return err
	}

	err = h.PackageManager.discover(h.Os)
	if err != nil {
		return err
	}

	err = h.ServiceManager.discover()
	if err != nil {
		return err
	}

	return nil
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
