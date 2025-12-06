// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"strings"

	"github.com/trippsoft/forge/pkg/discover"
	"github.com/zclconf/go-cty/cty"
)

// PackageManagerInfo contains information about the package manager of a managed host.
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
func (p *PackageManagerInfo) FromProtobuf(response *discover.PackageManagerInfoResponse) {
	p.name = response.Name
	p.path = response.Path
}

// String returns a string representation of the package manager information.
//
// This is useful for logging or debugging purposes.
func (p *PackageManagerInfo) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("package_manager_name: ")
	if p.name != "" {
		stringBuilder.WriteString(p.name)
	} else {
		stringBuilder.WriteString("unknown")
	}

	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("package_manager_path: ")
	if p.path != "" {
		stringBuilder.WriteString(p.path)
	} else {
		stringBuilder.WriteString("unknown")
	}

	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}
