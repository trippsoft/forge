// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// PackageManagerInfo represents information about the package manager on the host.
type PackageManagerInfo struct {
	name string
	path string
}

// Name returns the name of the package manager.
func (p *PackageManagerInfo) Name() string {
	return p.name
}

// Path returns the path to the package manager executable.
func (p *PackageManagerInfo) Path() string {
	return p.path
}

// ToMapOfCtyValues converts the PackageManagerInfo into a map of cty.Values.
func (p *PackageManagerInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	if p.name != "" {
		values["package_manager_name"] = cty.StringVal(p.name)
	} else {
		values["package_manager_name"] = cty.NullVal(cty.String)
	}

	if p.path != "" {
		values["package_manager_path"] = cty.StringVal(p.path)
	} else {
		values["package_manager_path"] = cty.NullVal(cty.String)
	}

	return values
}

// FromProtobuf populates the PackageManagerInfo from a protobuf representation.
func (p *PackageManagerInfo) FromProtobuf(other *PackageManagerInfoPB) {
	p.name = other.Name
	p.path = other.Path
}
