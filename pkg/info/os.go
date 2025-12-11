// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

// OSInfo represents information about the operating system.
type OSInfo struct {
	kernel       string
	id           string
	friendlyName string
	release      string
	releaseId    string
	majorVersion string
	version      string
	edition      string
	editionId    string
	arch         string
	families     *util.Set[string]
}

// Kernel returns the OS kernel.
func (o *OSInfo) Kernel() string {
	return o.kernel
}

// Id returns the OS ID.
func (o *OSInfo) Id() string {
	return o.id
}

// FriendlyName returns the OS friendly name.
func (o *OSInfo) FriendlyName() string {
	return o.friendlyName
}

// Release returns the OS release.
func (o *OSInfo) Release() string {
	return o.release
}

// ReleaseId returns the OS release ID.
func (o *OSInfo) ReleaseId() string {
	return o.releaseId
}

// MajorVersion returns the OS major version.
func (o *OSInfo) MajorVersion() string {
	return o.majorVersion
}

// Version returns the OS version.
func (o *OSInfo) Version() string {
	return o.version
}

// Edition returns the OS edition.
func (o *OSInfo) Edition() string {
	return o.edition
}

// EditionId returns the OS edition ID.
func (o *OSInfo) EditionId() string {
	return o.editionId
}

// Arch returns the OS architecture.
func (o *OSInfo) Arch() string {
	return o.arch
}

// Families returns the OS families.
func (o *OSInfo) Families() util.ReadOnlySet[string] {
	return o.families
}

// ToMapOfCtyValues converts the OSInfo into a map of cty.Values.
func (o *OSInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	if o.kernel == "" {
		values["os_kernel"] = cty.NullVal(cty.String)
	} else {
		values["os_kernel"] = cty.StringVal(o.kernel)
	}

	if o.id == "" {
		values["os_id"] = cty.NullVal(cty.String)
	} else {
		values["os_id"] = cty.StringVal(o.id)
	}

	if o.friendlyName == "" {
		values["os_friendly_name"] = cty.NullVal(cty.String)
	} else {
		values["os_friendly_name"] = cty.StringVal(o.friendlyName)
	}

	if o.release == "" {
		values["os_release"] = cty.NullVal(cty.String)
	} else {
		values["os_release"] = cty.StringVal(o.release)
	}

	if o.releaseId == "" {
		values["os_release_id"] = cty.NullVal(cty.String)
	} else {
		values["os_release_id"] = cty.StringVal(o.releaseId)
	}

	if o.majorVersion == "" {
		values["os_major_version"] = cty.NullVal(cty.String)
	} else {
		values["os_major_version"] = cty.StringVal(o.majorVersion)
	}

	if o.version == "" {
		values["os_version"] = cty.NullVal(cty.String)
	} else {
		values["os_version"] = cty.StringVal(o.version)
	}

	if o.edition == "" {
		values["os_edition"] = cty.NullVal(cty.String)
	} else {
		values["os_edition"] = cty.StringVal(o.edition)
	}

	if o.editionId == "" {
		values["os_edition_id"] = cty.NullVal(cty.String)
	} else {
		values["os_edition_id"] = cty.StringVal(o.editionId)
	}

	if o.arch == "" {
		values["os_arch"] = cty.NullVal(cty.String)
	} else {
		values["os_arch"] = cty.StringVal(o.arch)
	}

	if o.families.Size() > 0 {
		families := make([]cty.Value, 0, o.families.Size())
		for _, family := range o.families.Items() {
			families = append(families, cty.StringVal(family))
		}

		values["os_families"] = cty.SetVal(families)
	} else {
		values["os_families"] = cty.NullVal(cty.Set(cty.String))
	}

	return values
}

// FromProtobuf populates the OSInfo fields from a protobuf representation.
func (o *OSInfo) FromProtobuf(other *OSInfoPB) {
	o.kernel = other.Kernel
	o.id = other.Id
	o.friendlyName = other.FriendlyName
	o.release = other.Release
	o.releaseId = other.ReleaseId
	o.majorVersion = other.MajorVersion
	o.version = other.Version
	o.edition = other.Edition
	o.editionId = other.EditionId
	o.arch = other.Arch
	o.families = util.NewSet(other.Families...)
}

// NewOSInfo creates a new OSInfo instance.
func NewOSInfo() *OSInfo {
	return &OSInfo{
		families: util.NewSet[string](),
	}
}
