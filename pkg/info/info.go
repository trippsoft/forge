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
	ctx, cancel := context.WithTimeout(context.Background(), PopulateTimeout)
	defer cancel()
	connection, cleanup, err := t.StartPlugin(
		ctx,
		plugin.SharedPluginBasePath,
		"forge",
		"discover",
		nil,
	)

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

	i.CopyFrom(response.HostInfo)

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

// CopyFrom copies the HostInfo from another HostInfo.
func (h *HostInfo) CopyFrom(other *HostInfo) {
	h.Os.CopyFrom(other.Os)
	h.Fips.CopyFrom(other.Fips)
	h.AppArmor.CopyFrom(other.AppArmor)
	h.Selinux.CopyFrom(other.Selinux)
	h.PackageManager.CopyFrom(other.PackageManager)
	h.ServiceManager.CopyFrom(other.ServiceManager)
}

func (h *HostInfo) discover() error {
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

func discoverHostInfo() (*HostInfo, error) {
	hostInfo := NewHostInfo()

	err := hostInfo.discover()
	if err != nil {
		return nil, err
	}

	return hostInfo, nil
}
