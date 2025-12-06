// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"github.com/trippsoft/forge/pkg/discover"
	"github.com/zclconf/go-cty/cty"
)

// ServiceManagerInfo contains information about the service manager of a managed host.
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
func (s *ServiceManagerInfo) FromProtobuf(response *discover.ServiceManagerInfoResponse) {
	s.name = response.Name
}

// String returns a string representation of the service manager information.
//
// This is useful for logging or debugging purposes.
func (s *ServiceManagerInfo) String() string {
	if s.name == "" {
		return "service_manager: unknown\n"
	}

	return "service_manager: " + s.name + "\n"
}
