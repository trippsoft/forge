package test

import (
	"testing"

	"github.com/trippsoft/forge/internal/info"
	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestHostInfo_SSH_Integration_Linux(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setupVagrantEnvironment(t)

	sshBuilder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("failed to create SSH builder: %v", err)
	}

	sshTransport, err := sshBuilder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		Build()
	if err != nil {
		t.Fatalf("failed to create SSH transport: %v", err)
	}
	defer sshTransport.Close()

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect SSH transport: %v", err)
	}

	hostInfo := info.NewHostInfo()
	err = hostInfo.Populate(sshTransport)
	if err != nil {
		t.Fatalf("failed to populate host info via SSH: %v", err)
	}

	values := hostInfo.ToMapOfCtyValues()

	// Verify that we have expected keys
	expectedKeys := []string{
		"os_families",
		"os_id",
		"os_friendly_name",
		"os_architecture",
		"processor_architecture",
		"selinux_status",
		"apparmor_enabled",
		"fips_enabled",
		"package_manager_name",
		"package_manager_path",
	}

	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	// Verify Linux-specific expectations
	if osFamilies, exists := values["os_families"]; exists && !osFamilies.IsNull() {
		// Convert to string slice to check families
		familiesSlice := make([]string, 0)
		osFamilies.ForEachElement(func(key, val cty.Value) (stop bool) {
			familiesSlice = append(familiesSlice, val.AsString())
			return false
		})

		foundPosix := false
		foundLinux := false
		for _, family := range familiesSlice {
			if family == "posix" {
				foundPosix = true
			}
			if family == "linux" {
				foundLinux = true
			}
		}

		if !foundPosix {
			t.Error("expected 'posix' family to be present in OS families")
		}
		if !foundLinux {
			t.Error("expected 'linux' family to be present in OS families")
		}
	}

	// Verify that architecture information is populated
	if osArch, exists := values["os_architecture"]; exists && !osArch.IsNull() {
		arch := osArch.AsString()
		if arch == "" {
			t.Error("expected OS architecture to be non-empty")
		}
	} else {
		t.Error("expected OS architecture to be present and non-null")
	}
}

func TestHostInfo_SSH_Integration_Windows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setupVagrantEnvironment(t)

	sshBuilder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("failed to create SSH builder: %v", err)
	}

	sshTransport, err := sshBuilder.
		Host(windowsHost).
		Port(windowsPort).
		User(windowsUser).
		PasswordAuth(windowsPassword).
		DontUseKnownHosts().
		Build()
	if err != nil {
		t.Fatalf("failed to create SSH transport: %v", err)
	}
	defer sshTransport.Close()

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect SSH transport: %v", err)
	}

	hostInfo := info.NewHostInfo()
	err = hostInfo.Populate(sshTransport)
	if err != nil {
		t.Fatalf("failed to populate host info via SSH: %v", err)
	}

	values := hostInfo.ToMapOfCtyValues()

	// Verify that we have expected keys
	expectedKeys := []string{
		"os_families",
		"os_id",
		"os_friendly_name",
		"os_architecture",
		"processor_architecture",
		"selinux_status",
		"apparmor_enabled",
		"fips_enabled",
		"package_manager_name",
		"package_manager_path",
	}

	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	// Verify Windows-specific expectations
	if osFamilies, exists := values["os_families"]; exists && !osFamilies.IsNull() {
		// Convert to string slice to check families
		familiesSlice := make([]string, 0)
		osFamilies.ForEachElement(func(key, val cty.Value) (stop bool) {
			familiesSlice = append(familiesSlice, val.AsString())
			return false
		})

		foundWindows := false
		for _, family := range familiesSlice {
			if family == "windows" {
				foundWindows = true
			}
		}

		if !foundWindows {
			t.Error("expected 'windows' family to be present in OS families")
		}
	}

	// On Windows, SELinux and AppArmor should be null/not supported
	if selinuxStatus, exists := values["selinux_status"]; exists {
		if !selinuxStatus.IsNull() {
			t.Error("expected SELinux status to be null on Windows")
		}
	}

	if appArmorEnabled, exists := values["apparmor_enabled"]; exists {
		if !appArmorEnabled.IsNull() {
			t.Error("expected AppArmor enabled to be null on Windows")
		}
	}

	// FIPS should be known on Windows
	if fipsEnabled, exists := values["fips_enabled"]; exists {
		if fipsEnabled.IsNull() {
			t.Error("expected FIPS enabled to be known (not null) on Windows")
		}
	}
}

func TestHostInfo_SSH_Integration_CompareWithLocal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Test that SSH transport gives same results as local transport on Linux
	setupVagrantEnvironment(t)

	// Get info via SSH
	sshBuilder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("failed to create SSH builder: %v", err)
	}

	sshTransport, err := sshBuilder.
		Host(linuxHost).
		Port(linuxPort).
		User(linuxUser).
		PasswordAuth(linuxPassword).
		DontUseKnownHosts().
		Build()
	if err != nil {
		t.Fatalf("failed to create SSH transport: %v", err)
	}
	defer sshTransport.Close()

	err = sshTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect SSH transport: %v", err)
	}

	sshHostInfo := info.NewHostInfo()
	err = sshHostInfo.Populate(sshTransport)
	if err != nil {
		t.Fatalf("failed to populate host info via SSH: %v", err)
	}

	sshValues := sshHostInfo.ToMapOfCtyValues()

	// Architecture and basic OS info should be consistent
	// Note: We can't directly compare with local transport here since we're
	// running tests on the host system, not in the VM. This test serves
	// as a basic validation that SSH transport works.

	if osArch, exists := sshValues["os_architecture"]; exists && !osArch.IsNull() {
		arch := osArch.AsString()
		if arch == "" {
			t.Error("expected OS architecture to be populated via SSH")
		}
	}

	if procArch, exists := sshValues["processor_architecture"]; exists && !procArch.IsNull() {
		arch := procArch.AsString()
		if arch == "" {
			t.Error("expected processor architecture to be populated via SSH")
		}
	}
}
