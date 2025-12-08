// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/trippsoft/forge/pkg/discover"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

// HostInfo contains information about a managed host.
type HostInfo struct {
	os             *OSInfo
	fips           *FIPSInfo
	appArmor       *AppArmorInfo
	seLinux        *SELinuxInfo
	packageManager *PackageManagerInfo
	serviceManager *ServiceManagerInfo
}

// Populate retrieves and populates the HostInfo using the provided transport.
func (i *HostInfo) Populate(t transport.Transport) *result.Result {
	connection, cleanup, err := t.StartPlugin("forge", "discover", nil)
	if err != nil {
		err = fmt.Errorf("failed to start discovery client: %w", err)
		return result.NewFailure(err, err.Error())
	}
	defer connection.Close()
	defer cleanup()

	client := discover.NewDiscoveryPluginClient(connection)
	response, err := client.DiscoverInfo(context.Background(), &discover.DiscoverInfoRequest{})
	if err != nil {
		err = fmt.Errorf("failed to discover host info: %w", err)
		return result.NewFailure(err, err.Error())
	}

	i.FromProtobuf(response)

	return result.NewSuccess(false, nil)
}

// OS returns the OSInfo of the managed host.
func (i *HostInfo) OS() *OSInfo {
	return i.os
}

// FIPS returns the FIPSInfo of the managed host.
func (i *HostInfo) FIPS() *FIPSInfo {
	return i.fips
}

// AppArmor returns the AppArmorInfo of the managed host.
func (i *HostInfo) AppArmor() *AppArmorInfo {
	return i.appArmor
}

// SELinux returns the SELinuxInfo of the managed host.
func (i *HostInfo) SELinux() *SELinuxInfo {
	return i.seLinux
}

// PackageManager returns the PackageManagerInfo of the managed host.
func (i *HostInfo) PackageManager() *PackageManagerInfo {
	return i.packageManager
}

// ServiceManager returns the ServiceManagerInfo of the managed host.
func (i *HostInfo) ServiceManager() *ServiceManagerInfo {
	return i.serviceManager
}

// ToMapOfCtyValues converts the HostInfo into a map of cty.Values.
func (i *HostInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	maps.Copy(values, i.os.ToMapOfCtyValues())
	maps.Copy(values, i.fips.ToMapOfCtyValues())
	maps.Copy(values, i.appArmor.ToMapOfCtyValues())
	maps.Copy(values, i.seLinux.ToMapOfCtyValues())
	maps.Copy(values, i.packageManager.ToMapOfCtyValues())
	maps.Copy(values, i.serviceManager.ToMapOfCtyValues())
	return values
}

// FromProtobuf populates the HostInfo from a protobuf representation.
func (i *HostInfo) FromProtobuf(response *discover.DiscoverInfoResponse) {
	i.os.FromProtobuf(response.Os)
	i.fips.FromProtobuf(response.Fips)
	i.appArmor.FromProtobuf(response.AppArmor)
	i.seLinux.FromProtobuf(response.Selinux)
	i.packageManager.FromProtobuf(response.PackageManager)
	i.serviceManager.FromProtobuf(response.ServiceManager)
}

// String returns a string representation of the host information.
//
// This is useful for logging or debugging purposes.
func (i *HostInfo) String() string {
	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString(i.os.String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.fips.String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.appArmor.String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.seLinux.String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.packageManager.String())
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString(i.serviceManager.String())

	return stringBuilder.String()
}

// NewHostInfo creates a new HostInfo instance.
func NewHostInfo() *HostInfo {
	return &HostInfo{
		os:             NewOSInfo(),
		fips:           &FIPSInfo{},
		appArmor:       &AppArmorInfo{},
		seLinux:        &SELinuxInfo{},
		packageManager: &PackageManagerInfo{},
		serviceManager: &ServiceManagerInfo{},
	}
}
