// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
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
	runtime        *RuntimeInfo
	os             *OSInfo
	fips           *FIPSInfo
	appArmor       *AppArmorInfo
	seLinux        *SELinuxInfo
	packageManager *PackageManagerInfo
	serviceManager *ServiceManagerInfo
}

// Populate retrieves and populates the HostInfo using the provided transport.
func (i *HostInfo) Populate(t transport.Transport, runtimeOnly bool) *result.Result {
	var err error
	i.runtime.os, err = t.OS()
	if err != nil {
		err = fmt.Errorf("failed to get OS from transport: %w", err)
		return result.NewFailure(err, err.Error())
	}

	i.runtime.arch, err = t.Arch()
	if err != nil {
		err = fmt.Errorf("failed to get architecture from transport: %w", err)
		return result.NewFailure(err, err.Error())
	}

	if runtimeOnly {
		return result.NewSuccess(false, nil)
	}

	discoveryClient, err := t.StartDiscovery()
	if err != nil {
		err = fmt.Errorf("failed to start discovery client: %w", err)
		return result.NewFailure(err, err.Error())
	}

	response, err := discoveryClient.Discover()
	if err != nil {
		err = fmt.Errorf("failed to discover host info: %w", err)
		return result.NewFailure(err, err.Error())
	}

	i.FromProtobuf(response)

	return result.NewSuccess(false, nil)
}

// Runtime returns the RuntimeInfo of the managed host.
func (i *HostInfo) Runtime() *RuntimeInfo {
	return i.runtime
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
	maps.Copy(values, i.runtime.ToMapOfCtyValues())
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

	stringBuilder.WriteString(i.runtime.String())
	stringBuilder.WriteString("\n")

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
		runtime:        &RuntimeInfo{},
		os:             NewOSInfo(),
		fips:           &FIPSInfo{},
		appArmor:       &AppArmorInfo{},
		seLinux:        &SELinuxInfo{},
		packageManager: &PackageManagerInfo{},
		serviceManager: &ServiceManagerInfo{},
	}
}

// RuntimeInfo contains the OS and architecture information of a managed host.
type RuntimeInfo struct {
	os   string
	arch string
}

// OS returns the operating system of the managed host.
func (r *RuntimeInfo) OS() string {
	return r.os
}

// Arch returns the architecture of the managed host.
func (r *RuntimeInfo) Arch() string {
	return r.arch
}

// ToMapOfCtyValues converts the RuntimeInfo into a map of cty.Values.
func (r *RuntimeInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	if r.os == "" {
		values["runtime_os"] = cty.NullVal(cty.String)
	} else {
		values["runtime_os"] = cty.StringVal(r.os)
	}

	if r.arch == "" {
		values["runtime_arch"] = cty.NullVal(cty.String)
	} else {
		values["runtime_arch"] = cty.StringVal(r.arch)
	}

	return values
}

// String returns a string representation of the OS information.
//
// This is useful for logging or debugging purposes.
func (r *RuntimeInfo) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("runtime_os: ")
	stringBuilder.WriteString(r.os)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("runtime_arch: ")
	stringBuilder.WriteString(r.arch)
	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}

// NewRuntimeInfo creates a new RuntimeInfo instance.
func NewRuntimeInfo(os, arch string) *RuntimeInfo {
	return &RuntimeInfo{
		os:   os,
		arch: arch,
	}
}
