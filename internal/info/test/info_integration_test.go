package test

import (
	"strings"
	"testing"

	"github.com/trippsoft/forge/internal/info"
	"github.com/trippsoft/forge/internal/transport"
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

	osInfo := hostInfo.GetOSInfo()
	if osInfo == nil {
		t.Error("expected OS info to be populated")
	} else {
		families := osInfo.Families()
		if !families.Contains("posix") {
			t.Error("expected OS families to contain 'posix'")
		}
		if !families.Contains("linux") {
			t.Error("expected OS families to contain 'linux'")
		}
		if !families.Contains("el") {
			t.Error("expected OS families to contain 'el'")
		}
		if !families.Contains("rocky") {
			t.Error("expected OS families to contain 'rocky'")
		}
		if families.Size() != 4 {
			t.Errorf("expected OS families to have size 4, got %d", families.Size())
		}

		if osInfo.Id() != "rocky" {
			t.Errorf("expected OS ID to be 'rocky', got '%s'", osInfo.Id())
		}

		if !strings.Contains(osInfo.FriendlyName(), "Rocky Linux 9") {
			t.Errorf("expected OS friendly name to contain 'Rocky Linux 9', got '%s'", osInfo.FriendlyName())
		}

		if osInfo.Release() != "" {
			t.Errorf("expected OS release to be empty, got '%s'", osInfo.Release())
		}

		if osInfo.MajorVersion() != "9" {
			t.Errorf("expected OS major version to be '9', got '%s'", osInfo.MajorVersion())
		}

		if !strings.Contains(osInfo.Version(), "9") || osInfo.Version() == "9" {
			t.Errorf("expected OS version to contain '9', got '%s'", osInfo.Version())
		}

		if osInfo.Edition() != "" {
			t.Errorf("expected OS edition to be empty, got '%s'", osInfo.Edition())
		}

		if osInfo.EditionId() != "" {
			t.Errorf("expected OS edition ID to be empty, got '%s'", osInfo.EditionId())
		}

		if osInfo.OsArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OsArch())
		}

		if osInfo.OsArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OsArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.GetSELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if !selinuxInfo.Supported() {
			t.Error("expected SELinux to be supported on Rocky Linux")
		}
		if !selinuxInfo.Installed() {
			t.Error("expected SELinux to be installed on Rocky Linux")
		}
		if selinuxInfo.Status() != info.SelinuxEnforcing {
			t.Errorf("expected SELinux status to be 'enforcing', got '%s'", selinuxInfo.Status())
		}
		if selinuxInfo.SelinuxType() != info.SelinuxTypeTargeted {
			t.Errorf("expected SELinux type to be 'targeted', got '%s'", selinuxInfo.SelinuxType())
		}
	}

	appArmorInfo := hostInfo.GetAppArmorInfo()
	if appArmorInfo == nil {
		t.Error("expected AppArmor info to be populated")
	} else {
		if !appArmorInfo.Supported() {
			t.Error("expected AppArmor to be supported on Rocky Linux")
		}
		if appArmorInfo.Enabled() {
			t.Error("expected AppArmor to be disabled on Rocky Linux")
		}
	}

	fipsInfo := hostInfo.GetFipsInfo()
	if fipsInfo == nil {
		t.Error("expected FIPS info to be populated")
	} else {
		if !fipsInfo.Known() {
			t.Error("expected FIPS info to be known on Rocky Linux")
		}
		if fipsInfo.Enabled() {
			t.Error("expected FIPS to be disabled on Rocky Linux")
		}
	}

	packageManagerInfo := hostInfo.GetPackageManagerInfo()
	if packageManagerInfo == nil {
		t.Error("expected Package Manager info to be populated")
	} else {
		if packageManagerInfo.Name() != "dnf" {
			t.Errorf("expected Package Manager name to be 'dnf', got '%s'", packageManagerInfo.Name())
		}
		if packageManagerInfo.Path() != "/usr/bin/dnf-3" {
			t.Errorf("expected Package Manager path to be '/usr/bin/dnf-3', got '%s'", packageManagerInfo.Path())
		}
	}

	serviceManagerInfo := hostInfo.GetServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "systemd" {
			t.Errorf("expected Service Manager name to be 'systemd', got '%s'", serviceManagerInfo.Name())
		}
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

	osInfo := hostInfo.GetOSInfo()
	if osInfo == nil {
		t.Error("expected OS info to be populated")
	} else {
		families := osInfo.Families()
		if !families.Contains("windows") {
			t.Error("expected OS families to contain 'windows'")
		}
		if !families.Contains("windows-server") {
			t.Error("expected OS families to contain 'windows-server'")
		}
		if families.Size() != 2 {
			t.Errorf("expected OS families to have size 2, got %d", families.Size())
		}

		if osInfo.Id() != "windows-server" {
			t.Errorf("expected OS ID to be 'windows-server', got '%s'", osInfo.Id())
		}

		if !strings.Contains(osInfo.FriendlyName(), "Microsoft Windows Server 2025 Datacenter") {
			t.Errorf("expected OS friendly name to contain 'Microsoft Windows Server 2025 Datacenter', got '%s'", osInfo.FriendlyName())
		}

		if osInfo.Release() != "server-2025" {
			t.Errorf("expected OS release to be 'server-2025', got '%s'", osInfo.Release())
		}

		if osInfo.MajorVersion() != "10" {
			t.Errorf("expected OS major version to be '10', got '%s'", osInfo.MajorVersion())
		}

		if osInfo.Version() != "10.0.26100.0" {
			t.Errorf("expected OS version to be '10.0.26100.0', got '%s'", osInfo.Version())
		}

		if osInfo.Edition() != "Datacenter" {
			t.Errorf("expected OS edition to be 'Datacenter', got '%s'", osInfo.Edition())
		}

		if osInfo.EditionId() != "datacenter" {
			t.Errorf("expected OS edition ID to be 'datacenter', got '%s'", osInfo.EditionId())
		}

		if osInfo.OsArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OsArch())
		}

		if osInfo.OsArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OsArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.GetSELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if selinuxInfo.Supported() {
			t.Error("expected SELinux to be unsupported on Windows")
		}
	}

	appArmorInfo := hostInfo.GetAppArmorInfo()
	if appArmorInfo == nil {
		t.Error("expected AppArmor info to be populated")
	} else {
		if appArmorInfo.Supported() {
			t.Error("expected AppArmor to be unsupported on Windows")
		}
	}

	fipsInfo := hostInfo.GetFipsInfo()
	if fipsInfo == nil {
		t.Error("expected FIPS info to be populated")
	} else {
		if !fipsInfo.Known() {
			t.Error("expected FIPS info to be known on Windows")
		}
		if fipsInfo.Enabled() {
			t.Error("expected FIPS to be enabled on Windows")
		}
	}

	packageManagerInfo := hostInfo.GetPackageManagerInfo()
	if packageManagerInfo == nil {
		t.Error("expected Package Manager info to be populated")
	} else {
		// Blank for now, Windows is not supported yet
	}

	serviceManagerInfo := hostInfo.GetServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "windows-service-manager" {
			t.Errorf("expected Service Manager name to be 'windows-service-manager', got '%s'", serviceManagerInfo.Name())
		}
	}
}

func TestHostInfo_SSH_Integration_Cmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setupVagrantEnvironment(t)

	sshBuilder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("failed to create SSH builder: %v", err)
	}

	sshTransport, err := sshBuilder.
		Host(cmdHost).
		Port(cmdPort).
		User(cmdUser).
		PasswordAuth(cmdPassword).
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

	osInfo := hostInfo.GetOSInfo()
	if osInfo == nil {
		t.Error("expected OS info to be populated")
	} else {
		families := osInfo.Families()
		if !families.Contains("windows") {
			t.Error("expected OS families to contain 'windows'")
		}
		if !families.Contains("windows-server") {
			t.Error("expected OS families to contain 'windows-server'")
		}
		if families.Size() != 2 {
			t.Errorf("expected OS families to have size 2, got %d", families.Size())
		}

		if osInfo.Id() != "windows-server" {
			t.Errorf("expected OS ID to be 'windows-server', got '%s'", osInfo.Id())
		}

		if !strings.Contains(osInfo.FriendlyName(), "Microsoft Windows Server 2025 Datacenter") {
			t.Errorf("expected OS friendly name to contain 'Microsoft Windows Server 2025 Datacenter', got '%s'", osInfo.FriendlyName())
		}

		if osInfo.Release() != "server-2025" {
			t.Errorf("expected OS release to be 'server-2025', got '%s'", osInfo.Release())
		}

		if osInfo.MajorVersion() != "10" {
			t.Errorf("expected OS major version to be '10', got '%s'", osInfo.MajorVersion())
		}

		if osInfo.Version() != "10.0.26100.0" {
			t.Errorf("expected OS version to be '10.0.26100.0', got '%s'", osInfo.Version())
		}

		if osInfo.Edition() != "Datacenter" {
			t.Errorf("expected OS edition to be 'Datacenter', got '%s'", osInfo.Edition())
		}

		if osInfo.EditionId() != "datacenter" {
			t.Errorf("expected OS edition ID to be 'datacenter', got '%s'", osInfo.EditionId())
		}

		if osInfo.OsArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OsArch())
		}

		if osInfo.OsArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OsArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.GetSELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if selinuxInfo.Supported() {
			t.Error("expected SELinux to be unsupported on Windows")
		}
	}

	appArmorInfo := hostInfo.GetAppArmorInfo()
	if appArmorInfo == nil {
		t.Error("expected AppArmor info to be populated")
	} else {
		if appArmorInfo.Supported() {
			t.Error("expected AppArmor to be unsupported on Windows")
		}
	}

	fipsInfo := hostInfo.GetFipsInfo()
	if fipsInfo == nil {
		t.Error("expected FIPS info to be populated")
	} else {
		if !fipsInfo.Known() {
			t.Error("expected FIPS info to be known on Windows")
		}
		if fipsInfo.Enabled() {
			t.Error("expected FIPS to be enabled on Windows")
		}
	}

	packageManagerInfo := hostInfo.GetPackageManagerInfo()
	if packageManagerInfo == nil {
		t.Error("expected Package Manager info to be populated")
	} else {
		// Blank for now, Windows is not supported yet
	}

	serviceManagerInfo := hostInfo.GetServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "windows-service-manager" {
			t.Errorf("expected Service Manager name to be 'windows-service-manager', got '%s'", serviceManagerInfo.Name())
		}
	}
}
