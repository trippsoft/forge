// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the PackageManagerInfo into a map of cty.Values.
func (p *PackageManagerInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	if p.Name != "" {
		values["package_manager_name"] = cty.StringVal(p.Name)
	} else {
		values["package_manager_name"] = cty.NullVal(cty.String)
	}

	if p.Path != "" {
		values["package_manager_path"] = cty.StringVal(p.Path)
	} else {
		values["package_manager_path"] = cty.NullVal(cty.String)
	}

	return values
}

// From populates the PackageManagerInfo from another PackageManagerInfo.
func (p *PackageManagerInfo) From(other *PackageManagerInfo) {
	p.Name = other.Name
	p.Path = other.Path
}
