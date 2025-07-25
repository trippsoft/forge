package info

import (
	"testing"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestHostInfo_Populate_NilTransport(t *testing.T) {
	hostInfo := NewHostInfo()
	err := hostInfo.Populate(nil)

	if err == nil {
		t.Fatal("expected error for nil transport, got nil")
	}

	expectedError := "transport cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestHostInfo_Populate_NullFileSystem(t *testing.T) {
	// Create a mock transport with null filesystem
	mockTransport := transport.NewNoneTransport()

	hostInfo := NewHostInfo()
	err := hostInfo.Populate(mockTransport)

	if err == nil {
		t.Fatal("expected error for null filesystem, got nil")
	}

	expectedError := "file system is null or not supported"
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
	err = hostInfo.Populate(localTransport)
	if err != nil {
		t.Fatalf("failed to populate host info: %v", err)
	}

	// Verify that all components were populated
	if hostInfo.osInfo.families == nil || hostInfo.osInfo.families.Size() == 0 {
		t.Error("expected OS families to be populated")
	}

	// Basic validation of OS info
	if hostInfo.osInfo.osArch == "" {
		t.Error("expected OS architecture to be populated")
	}

	if hostInfo.osInfo.procArch == "" {
		t.Error("expected processor architecture to be populated")
	}
}

func TestHostInfo_Populate_ErrorHandling(t *testing.T) {
	// Test with transport that returns errors
	errorTransport := &errorTransport{shouldError: true}

	hostInfo := NewHostInfo()
	err := hostInfo.Populate(errorTransport)
	if err == nil {
		t.Error("expected error when transport operations fail")
	}
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
	err = hostInfo.Populate(localTransport)
	if err != nil {
		t.Fatalf("failed to populate host info: %v", err)
	}

	values := hostInfo.ToMapOfCtyValues()

	// Verify that we have expected keys from all components
	expectedKeys := []string{
		"os_families",
		"os_id",
		"os_friendly_name",
		"os_release",
		"os_major_version",
		"os_version",
		"os_edition",
		"os_edition_id",
		"os_architecture",
		"os_architecture_bits",
		"processor_architecture",
		"processor_architecture_bits",
		"selinux_status",
		"selinux_type",
		"apparmor_enabled",
		"fips_enabled",
	}

	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	// Verify that os_families is a set
	if osFamilies, exists := values["os_families"]; exists {
		if osFamilies.Type() != cty.Set(cty.String) && !osFamilies.IsNull() {
			t.Errorf("expected os_families to be a set of strings or null, got %s", osFamilies.Type().GoString())
		}
	}
}

func TestHostInfo_ToMapOfCtyValues_EmptyHostInfo(t *testing.T) {
	hostInfo := NewHostInfo()
	values := hostInfo.ToMapOfCtyValues()

	// For an unpopulated HostInfo, most values should be null
	expectedNullKeys := []string{
		"os_families",
		"os_id",
		"os_friendly_name",
		"os_release",
		"os_major_version",
		"os_version",
		"os_edition",
		"os_edition_id",
		"os_architecture",
		"os_architecture_bits",
		"processor_architecture",
		"processor_architecture_bits",
		"selinux_status",
		"selinux_type",
		"apparmor_enabled",
		"fips_enabled",
	}

	for _, key := range expectedNullKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected key %q to be null for empty HostInfo, got %s", key, value.GoString())
			}
		}
	}
}
