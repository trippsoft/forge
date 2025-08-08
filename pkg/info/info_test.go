// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"testing"
)

func TestHostInfo_Populate_NilTransport(t *testing.T) {

	hostInfo := NewHostInfo()
	err := hostInfo.Populate(nil)

	if err == nil {
		t.Fatal("expected error for nil transport, got nil")
	}

	expectedError := "Invalid transport; Transport cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestHostInfo_ToMapOfCtyValues(t *testing.T) {

	hostInfo := NewHostInfo()

	osValues := hostInfo.OSInfo().toMapOfCtyValues()
	selinuxValues := hostInfo.SELinuxInfo().toMapOfCtyValues()
	appArmorValues := hostInfo.AppArmorInfo().toMapOfCtyValues()
	fipsValues := hostInfo.FipsInfo().toMapOfCtyValues()
	packageManagerValues := hostInfo.PackageManagerInfo().toMapOfCtyValues()
	serviceManagerValues := hostInfo.ServiceManagerInfo().toMapOfCtyValues()
	userValues := hostInfo.UserInfo().toMapOfCtyValues()

	totalLength := len(osValues) + len(selinuxValues) + len(appArmorValues) +
		len(fipsValues) + len(packageManagerValues) + len(serviceManagerValues) +
		len(userValues)

	values := hostInfo.ToMapOfCtyValues()

	if len(values) != totalLength {
		t.Errorf("expected %d values, got %d", totalLength, len(values))
	}
}
