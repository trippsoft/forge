package info

import (
	"testing"

	"github.com/trippsoft/forge/internal/transport"
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

func TestHostInfo_Populate_WithLocalTransport(t *testing.T) {
	localTransport, err := transport.NewLocalTransport()
	if err != nil {
		t.Fatalf("failed to create local transport: %v", err)
	}
	defer localTransport.Close()

	err = localTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect local transport: %v", err)
	}

	hostInfo := NewHostInfo()
	diags := hostInfo.Populate(localTransport)
	if diags.HasErrors() {
		t.Fatalf("failed to populate host info: %v", diags)
	}

	t.Logf("%s", hostInfo.String()) // This test is used primarily for manual review.
}

func TestHostInfo_ToMapOfCtyValues(t *testing.T) {
	localTransport, err := transport.NewLocalTransport()
	if err != nil {
		t.Fatalf("failed to create local transport: %v", err)
	}
	defer localTransport.Close()

	err = localTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect local transport: %v", err)
	}

	hostInfo := NewHostInfo()
	diags := hostInfo.Populate(localTransport)
	if diags.HasErrors() {
		t.Fatalf("failed to populate host info: %v", diags)
	}

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
