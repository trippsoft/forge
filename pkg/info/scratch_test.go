// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"runtime"
	"testing"
)

func TestScratch(t *testing.T) {
	os, err := discoverOSInfo()
	if err != nil {
		t.Fatalf("discoverOSInfo() returned error: %v", err)
	}
	if os == nil {
		t.Fatalf("discoverOSInfo() returned nil info")
	}

	appArmor, err := discoverAppArmorInfo()
	if err != nil {
		t.Fatalf("discoverAppArmorInfo() returned error: %v", err)
	}
	if appArmor == nil {
		t.Fatalf("discoverAppArmorInfo() returned nil info")
	}

	seLinux, err := discoverSELinuxInfo()
	if err != nil {
		t.Fatalf("discoverSELinuxInfo() returned error: %v", err)
	}
	if seLinux == nil {
		t.Fatalf("discoverSELinuxInfo() returned nil info")
	}

	fips, err := discoverFIPSInfo()
	if err != nil {
		t.Fatalf("discoverFIPSInfo() returned error: %v", err)
	}
	if fips == nil {
		t.Fatalf("discoverFIPSInfo() returned nil info")
	}

	packageManager, err := discoverPackageManagerInfo(os)
	if err != nil {
		t.Fatalf("discoverPackageManagerInfo() returned error: %v", err)
	}
	if packageManager == nil {
		t.Fatalf("discoverPackageManagerInfo() returned nil info")
	}

	serviceManager, err := discoverServiceManagerInfo()
	if err != nil {
		t.Fatalf("discoverServiceManagerInfo() returned error: %v", err)
	}
	if serviceManager == nil {
		t.Fatalf("discoverServiceManagerInfo() returned nil info")
	}

	response := &DiscoverInfoResponse{
		Os:             os,
		Fips:           fips,
		AppArmor:       appArmor,
		Selinux:        seLinux,
		PackageManager: packageManager,
		ServiceManager: serviceManager,
	}

	hostInfo := NewHostInfo()
	hostInfo.runtime.os = runtime.GOOS
	hostInfo.runtime.arch = runtime.GOARCH
	hostInfo.FromProtobuf(response)
}
