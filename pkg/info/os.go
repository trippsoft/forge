// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"slices"

	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the OSInfo into a map of cty.Values.
func (o *OSInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	if o.Kernel == "" {
		values["os_kernel"] = cty.NullVal(cty.String)
	} else {
		values["os_kernel"] = cty.StringVal(o.Kernel)
	}

	if o.Id == "" {
		values["os_id"] = cty.NullVal(cty.String)
	} else {
		values["os_id"] = cty.StringVal(o.Id)
	}

	if o.FriendlyName == "" {
		values["os_friendly_name"] = cty.NullVal(cty.String)
	} else {
		values["os_friendly_name"] = cty.StringVal(o.FriendlyName)
	}

	if o.Release == "" {
		values["os_release"] = cty.NullVal(cty.String)
	} else {
		values["os_release"] = cty.StringVal(o.Release)
	}

	if o.ReleaseId == "" {
		values["os_release_id"] = cty.NullVal(cty.String)
	} else {
		values["os_release_id"] = cty.StringVal(o.ReleaseId)
	}

	if o.MajorVersion == "" {
		values["os_major_version"] = cty.NullVal(cty.String)
	} else {
		values["os_major_version"] = cty.StringVal(o.MajorVersion)
	}

	if o.Version == "" {
		values["os_version"] = cty.NullVal(cty.String)
	} else {
		values["os_version"] = cty.StringVal(o.Version)
	}

	if o.Edition == "" {
		values["os_edition"] = cty.NullVal(cty.String)
	} else {
		values["os_edition"] = cty.StringVal(o.Edition)
	}

	if o.EditionId == "" {
		values["os_edition_id"] = cty.NullVal(cty.String)
	} else {
		values["os_edition_id"] = cty.StringVal(o.EditionId)
	}

	if o.Arch == "" {
		values["os_arch"] = cty.NullVal(cty.String)
	} else {
		values["os_arch"] = cty.StringVal(o.Arch)
	}

	if len(o.Families) > 0 {
		families := make([]cty.Value, 0, len(o.Families))
		for _, family := range o.Families {
			families = append(families, cty.StringVal(family))
		}

		values["os_families"] = cty.SetVal(families)
	} else {
		values["os_families"] = cty.NullVal(cty.Set(cty.String))
	}

	return values
}

// CopyFrom copies the OSInfo from another OSInfo.
func (o *OSInfo) CopyFrom(other *OSInfo) {
	o.Kernel = other.Kernel
	o.Id = other.Id
	o.FriendlyName = other.FriendlyName
	o.Release = other.Release
	o.ReleaseId = other.ReleaseId
	o.MajorVersion = other.MajorVersion
	o.Version = other.Version
	o.Edition = other.Edition
	o.EditionId = other.EditionId
	o.Arch = other.Arch
	o.Families = slices.Clone(other.Families)
}

// NewOSInfo creates a new OSInfo instance.
func NewOSInfo() *OSInfo {
	return &OSInfo{
		Families: make([]string, 0),
	}
}
