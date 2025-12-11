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

// HostInfo represents comprehensive information about the host system.
type HostInfo struct {
	os             *OSInfo
	fips           *FIPSInfo
	appArmor       *AppArmorInfo
	seLinux        *SELinuxInfo
	packageManager *PackageManagerInfo
	serviceManager *ServiceManagerInfo
}

// OS returns the OSInfo of the host.
func (i *HostInfo) OS() *OSInfo {
	return i.os
}

// SELinux returns the SELinuxInfo of the host.
func (i *HostInfo) SELinux() *SELinuxInfo {
	return i.seLinux
}

// AppArmor returns the AppArmorInfo of the host.
func (i *HostInfo) AppArmor() *AppArmorInfo {
	return i.appArmor
}

// FIPS returns the FIPSInfo of the host.
func (i *HostInfo) FIPS() *FIPSInfo {
	return i.fips
}

// PackageManager returns the PackageManagerInfo of the host.
func (i *HostInfo) PackageManager() *PackageManagerInfo {
	return i.packageManager
}

// ServiceManager returns the ServiceManagerInfo of the host.
func (i *HostInfo) ServiceManager() *ServiceManagerInfo {
	return i.serviceManager
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

	client := NewDiscoveryPluginClient(connection)
	response, err := client.DiscoverInfo(context.Background(), &DiscoverInfoRequest{})
	_, _ = client.Shutdown(context.Background(), &ShutdownRequest{})
	if err != nil {
		err = fmt.Errorf("failed to discover host info: %w", err)
		return result.NewFailure(err, err.Error())
	}

	i.FromProtobuf(response.HostInfo)

	return result.NewSuccess(false, nil)
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
func (i *HostInfo) FromProtobuf(other *HostInfoPB) {
	i.os.FromProtobuf(other.Os)
	i.fips.FromProtobuf(other.Fips)
	i.appArmor.FromProtobuf(other.AppArmor)
	i.seLinux.FromProtobuf(other.Selinux)
	i.packageManager.FromProtobuf(other.PackageManager)
	i.serviceManager.FromProtobuf(other.ServiceManager)
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

func discoverHostInfo() (*HostInfoPB, error) {
	hostInfo := &HostInfoPB{
		Os:             &OSInfoPB{},
		Fips:           &FIPSInfoPB{},
		AppArmor:       &AppArmorInfoPB{},
		Selinux:        &SELinuxInfoPB{},
		PackageManager: &PackageManagerInfoPB{},
		ServiceManager: &ServiceManagerInfoPB{},
	}

	err := hostInfo.discover()
	if err != nil {
		return nil, err
	}

	return hostInfo, nil
}

func (h *HostInfoPB) discover() error {
	err := h.Os.discover()
	if err != nil {
		return err
	}

	err = h.Fips.discover()
	if err != nil {
		return err
	}

	err = h.AppArmor.discover()
	if err != nil {
		return err
	}

	err = h.Selinux.discover()
	if err != nil {
		return err
	}

	err = h.PackageManager.discover(h.Os)
	if err != nil {
		return err
	}

	err = h.ServiceManager.discover()
	if err != nil {
		return err
	}

	return nil
}
