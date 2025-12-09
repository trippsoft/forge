// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package linux

import (
	"testing"

	"github.com/trippsoft/forge/test/integration"
)

func TestHostInfo(t *testing.T) {
	tests := []struct {
		name     string
		expected integration.ExpectedHostInfo
	}{
		{
			name: "debian13",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "debian"},
					Kernel:       "linux",
					ID:           "debian",
					FriendlyName: "Debian GNU/Linux 13 (trixie)",
					Release:      "trixie",
					ReleaseId:    "trixie",
					MajorVersion: "13",
					Version:      "13",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "apt",
					Path: "/usr/bin/apt-get",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "debian12",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "debian"},
					Kernel:       "linux",
					ID:           "debian",
					FriendlyName: "Debian GNU/Linux 12 (bookworm)",
					Release:      "bookworm",
					ReleaseId:    "bookworm",
					MajorVersion: "12",
					Version:      "12",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "apt",
					Path: "/usr/bin/apt-get",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "fedora42",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "el", "fedora"},
					Kernel:       "linux",
					ID:           "fedora",
					FriendlyName: "Fedora Linux 42 (Container Image)",
					Release:      "",
					ReleaseId:    "",
					MajorVersion: "42",
					Version:      "42",
					Edition:      "Container Image",
					EditionID:    "container",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "dnf5",
					Path: "/usr/bin/dnf5",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "fedora41",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "el", "fedora"},
					Kernel:       "linux",
					ID:           "fedora",
					FriendlyName: "Fedora Linux 41 (Container Image)",
					Release:      "",
					ReleaseId:    "",
					MajorVersion: "41",
					Version:      "41",
					Edition:      "Container Image",
					EditionID:    "container",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "dnf5",
					Path: "/usr/bin/dnf5",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "rocky10",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "el", "rocky"},
					Kernel:       "linux",
					ID:           "rocky",
					FriendlyName: "Rocky Linux 10.0 (Red Quartz)",
					Release:      "",
					ReleaseId:    "",
					MajorVersion: "10",
					Version:      "10.0",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "dnf",
					Path: "/usr/bin/dnf-3",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "rocky9",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "el", "rocky"},
					Kernel:       "linux",
					ID:           "rocky",
					FriendlyName: "Rocky Linux 9.6 (Blue Onyx)",
					Release:      "",
					ReleaseId:    "",
					MajorVersion: "9",
					Version:      "9.6",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "dnf",
					Path: "/usr/bin/dnf-3",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "rocky8",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "el", "rocky"},
					Kernel:       "linux",
					ID:           "rocky",
					FriendlyName: "Rocky Linux 8.10 (Green Obsidian)",
					Release:      "",
					ReleaseId:    "",
					MajorVersion: "8",
					Version:      "8.10",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   true,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "dnf",
					Path: "/usr/bin/dnf-3",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "ubuntu2404",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "debian", "ubuntu"},
					Kernel:       "linux",
					ID:           "ubuntu",
					FriendlyName: "Ubuntu 24.04.3 LTS",
					Release:      "noble",
					ReleaseId:    "noble",
					MajorVersion: "24",
					Version:      "24.04",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "apt",
					Path: "/usr/bin/apt-get",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
		{
			name: "ubuntu2204",
			expected: integration.ExpectedHostInfo{
				OS: integration.ExpectedOSInfo{
					Families:     []string{"posix", "linux", "debian", "ubuntu"},
					Kernel:       "linux",
					ID:           "ubuntu",
					FriendlyName: "Ubuntu 22.04.5 LTS",
					Release:      "jammy",
					ReleaseId:    "jammy",
					MajorVersion: "22",
					Version:      "22.04",
					Edition:      "",
					EditionID:    "",
					Arch:         "amd64",
				},
				SELinux: integration.ExpectedSELinuxInfo{
					Supported: true,
					Installed: false,
					Status:    "",
					Type:      "",
				},
				AppArmor: integration.ExpectedAppArmorInfo{
					Supported: true,
					Enabled:   false,
				},
				FIPS: integration.ExpectedFIPSInfo{
					Known:   true,
					Enabled: false,
				},
				PackageManager: integration.ExpectedPackageManagerInfo{
					Name: "apt",
					Path: "/usr/bin/apt-get",
				},
				ServiceManager: integration.ExpectedServiceManagerInfo{
					Name: "systemd",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, ok := inv.Host(tt.name)
			if !ok {
				t.Fatalf("Host %q not found in inventory", tt.name)
			}

			hostInfo := host.Info()
			result := hostInfo.Populate(host.Transport())
			if result.Err != nil {
				t.Fatalf("failed to populate host info via SSH: %v", result.Err)
			}

			tt.expected.Verify(t, hostInfo)
		})
	}
}
