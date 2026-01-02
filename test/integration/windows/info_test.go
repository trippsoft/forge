// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package windows

import (
	"testing"

	"github.com/trippsoft/forge/test/integration"
)

func TestHostInfo(t *testing.T) {
	hostnames := []string{"windows", "cmd"}

	expected := integration.ExpectedHostInfo{
		OS: integration.ExpectedOSInfo{
			Families:     []string{"windows", "windows-server"},
			Kernel:       "windows",
			Id:           "windows-server",
			FriendlyName: "Microsoft Windows Server 2025 Datacenter",
			Release:      "Server 2025",
			ReleaseId:    "server-2025",
			MajorVersion: "10",
			Version:      "10.0.26100.0",
			Edition:      "Datacenter",
			EditionID:    "datacenter",
			Arch:         "amd64",
		},
		SELinux: integration.ExpectedSELinuxInfo{
			Supported: false,
			Installed: false,
			Status:    "",
			Type:      "",
		},
		AppArmor: integration.ExpectedAppArmorInfo{
			Supported: false,
			Enabled:   false,
		},
		FIPS: integration.ExpectedFIPSInfo{
			Known:   true,
			Enabled: false,
		},
		PackageManager: integration.ExpectedPackageManagerInfo{
			Name: "",
			Path: "",
		},
		ServiceManager: integration.ExpectedServiceManagerInfo{
			Name: "windows",
		},
	}

	for _, hostname := range hostnames {
		t.Run(hostname, func(t *testing.T) {
			host, ok := inv.Host(hostname)
			if !ok {
				t.Fatalf("failed to get %s host", hostname)
			}

			hostInfo := host.Info()
			result := hostInfo.Populate(host.Transport())
			if result.Err != nil {
				t.Fatalf("failed to populate host info via SSH: %v", result.Err)
			}

			expected.Verify(t, hostInfo)
		})
	}
}
