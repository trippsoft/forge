// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/trippsoft/forge/pkg/plugin"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	PopulateTimeout = 30 * time.Minute
)

// Populate retrieves and populates the HostInfo using the provided transport.
func (i *HostInfo) Populate(t transport.Transport) *result.Result {
	session, err := t.StartPluginSession(
		context.Background(),
		plugin.SharedPluginBasePath,
		"forge",
		"discover",
		nil,
	)
	if err != nil {
		err = fmt.Errorf("failed to start plugin session: %w", err)
		return result.NewFailure(err, "")
	}
	defer session.Close()

	request := &DiscoverRequest{}
	err = plugin.Write(session.Stdin(), request)
	if err != nil {
		err = fmt.Errorf("failed to write discover request: %w", err)
		return result.NewFailure(err, "")
	}

	response := &DiscoverResponse{}
	err = plugin.Read(session.Stdout(), response)
	if err != nil {
		err = fmt.Errorf("failed to read discover response: %w", err)
		return result.NewFailure(err, "")
	}

	i.CopyFrom(response.HostInfo)

	r := result.NewNotChanged(cty.EmptyObjectVal)
	if len(response.Warnings) > 0 {
		r.Messages = response.Warnings
	}

	return r
}

// Discover performs discovery of the host information and returns any warnings encountered during discovery.
func (h *HostInfo) Discover() []string {
	warnings := make([]string, 0)

	w := h.Os.discover()
	warnings = append(warnings, w...)

	w = h.Fips.discover()
	warnings = append(warnings, w...)

	w = h.AppArmor.discover()
	warnings = append(warnings, w...)

	w = h.Selinux.discover()
	warnings = append(warnings, w...)

	w = h.PackageManager.discover(h.Os)
	warnings = append(warnings, w...)

	w = h.ServiceManager.discover()
	warnings = append(warnings, w...)

	if len(warnings) == 0 {
		return nil
	}

	return warnings
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

// CopyFrom copies the HostInfo from another HostInfo.
func (h *HostInfo) CopyFrom(other *HostInfo) {
	h.Os.CopyFrom(other.Os)
	h.Fips.CopyFrom(other.Fips)
	h.AppArmor.CopyFrom(other.AppArmor)
	h.Selinux.CopyFrom(other.Selinux)
	h.PackageManager.CopyFrom(other.PackageManager)
	h.ServiceManager.CopyFrom(other.ServiceManager)
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

func DiscoverHostInfo() *DiscoverResponse {
	hostInfo := NewHostInfo()
	warnings := hostInfo.Discover()

	return &DiscoverResponse{
		HostInfo: hostInfo,
		Warnings: warnings,
	}
}
