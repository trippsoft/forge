// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/zclconf/go-cty/cty"
)

// ServiceManagerInfo represents information about the service manager on the host.
type ServiceManagerInfo struct {
	name string
}

// Name returns the name of the service manager.
func (s *ServiceManagerInfo) Name() string {
	return s.name
}

// ToMapOfCtyValues converts the ServiceManagerInfo into a map of cty.Values.
func (s *ServiceManagerInfo) ToMapOfCtyValues() map[string]cty.Value {
	if s.name == "" {
		return map[string]cty.Value{
			"service_manager": cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"service_manager": cty.StringVal(s.name),
	}
}

// FromProtobuf populates the ServiceManagerInfo from a protobuf representation.
func (s *ServiceManagerInfo) FromProtobuf(other *ServiceManagerInfoPB) {
	s.name = other.Name
}
