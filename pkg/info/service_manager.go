// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ToMapOfCtyValues converts the ServiceManagerInfo into a map of cty.Values.
func (s *ServiceManagerInfo) ToMapOfCtyValues() map[string]cty.Value {
	if s.Name == "" {
		return map[string]cty.Value{
			"service_manager": cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"service_manager": cty.StringVal(s.Name),
	}
}

// From populates the ServiceManagerInfo from another ServiceManagerInfo.
func (s *ServiceManagerInfo) From(other *ServiceManagerInfo) {
	s.Name = other.Name
}
