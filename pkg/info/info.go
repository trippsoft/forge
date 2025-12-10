// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"fmt"
	"maps"

	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

// Populate retrieves and populates the HostInfo using the provided transport.
func (i *HostInfo) Populate(t transport.Transport) *result.Result {
	connection, cleanup, err := t.StartPlugin("forge", "discover", nil)
	if err != nil {
		err = fmt.Errorf("failed to start discovery client: %w", err)
		return result.NewFailure(err, err.Error())
	}
	defer connection.Close()
	defer cleanup()

	client := NewDiscoveryPluginClient(connection)
	response, err := client.DiscoverInfo(context.Background(), &DiscoverInfoRequest{})
	_, _ = client.Shutdown(context.Background(), &ShutdownRequest{})
	if err != nil {
		err = fmt.Errorf("failed to discover host info: %w", err)
		return result.NewFailure(err, err.Error())
	}

	i.From(response.HostInfo)

	return result.NewSuccess(false, nil)
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

// From populates the HostInfo from another HostInfo.
func (i *HostInfo) From(other *HostInfo) {
	i.Os.From(other.Os)
	i.Fips.From(other.Fips)
	i.AppArmor.From(other.AppArmor)
	i.Selinux.From(other.Selinux)
	i.PackageManager.From(other.PackageManager)
	i.ServiceManager.From(other.ServiceManager)
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
