package test

import (
	"strings"
	"testing"

	"github.com/trippsoft/forge/internal/info"
	"github.com/trippsoft/forge/internal/transport"
)

func TestHostInfo_Local(t *testing.T) {

	transport, err := transport.NewLocalTransport()
	if err != nil {
		t.Fatalf("failed to create local transport: %v", err)
	}
	transport.Connect()
	defer transport.Close()

	hostInfo := info.NewHostInfo()
	diags := hostInfo.Populate(transport)

	if diags.HasErrors() {
		t.Fatalf("failed to populate host info: %v", diags)
	}

	t.Log(hostInfo.String()) // This test is used primarily for manual review.
}

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
	diags := hostInfo.Populate(sshTransport)
	if diags.HasErrors() {
		t.Fatalf("failed to populate host info via SSH: %v", diags.Errors())
	}

	osInfo := hostInfo.OSInfo()
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

		if osInfo.ID() != "rocky" {
			t.Errorf("expected OS ID to be 'rocky', got '%s'", osInfo.ID())
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

		if osInfo.EditionID() != "" {
			t.Errorf("expected OS edition ID to be empty, got '%s'", osInfo.EditionID())
		}

		if osInfo.OSArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OSArch())
		}

		if osInfo.OSArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OSArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.SELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if !selinuxInfo.Supported() {
			t.Error("expected SELinux to be supported on Rocky Linux")
		}
		if !selinuxInfo.Installed() {
			t.Error("expected SELinux to be installed on Rocky Linux")
		}
		if selinuxInfo.Status() != info.SELinuxEnforcing {
			t.Errorf("expected SELinux status to be 'enforcing', got '%s'", selinuxInfo.Status())
		}
		if selinuxInfo.SelinuxType() != info.SELinuxTypeTargeted {
			t.Errorf("expected SELinux type to be 'targeted', got '%s'", selinuxInfo.SelinuxType())
		}
	}

	appArmorInfo := hostInfo.AppArmorInfo()
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

	fipsInfo := hostInfo.FipsInfo()
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

	// packageManagerInfo := hostInfo.PackageManagerInfo()
	// if packageManagerInfo == nil {
	// 	t.Error("expected Package Manager info to be populated")
	// } else {
	// 	if packageManagerInfo.Name() != "dnf" {
	// 		t.Errorf("expected Package Manager name to be 'dnf', got '%s'", packageManagerInfo.Name())
	// 	}
	// 	if packageManagerInfo.Path() != "/usr/bin/dnf-3" {
	// 		t.Errorf("expected Package Manager path to be '/usr/bin/dnf-3', got '%s'", packageManagerInfo.Path())
	// 	}
	// }

	serviceManagerInfo := hostInfo.ServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "systemd" {
			t.Errorf("expected Service Manager name to be 'systemd', got '%s'", serviceManagerInfo.Name())
		}
	}

	userInfo := hostInfo.UserInfo()
	if userInfo == nil {
		t.Error("expected User info to be populated")
	} else {
		if userInfo.Name() != "vagrant" {
			t.Errorf("expected User name to be 'vagrant', got '%s'", userInfo.Name())
		}
		if userInfo.UserId() != "1000" {
			t.Errorf("expected User ID to be '1000', got '%s'", userInfo.UserId())
		}
		if userInfo.GroupId() != "1000" {
			t.Errorf("expected Group ID to be '1000', got '%s'", userInfo.GroupId())
		}
		if userInfo.HomeDir() != "/home/vagrant" {
			t.Errorf("expected Home Directory to be '/home/vagrant', got '%s'", userInfo.HomeDir())
		}
		if userInfo.Shell() != "/bin/bash" {
			t.Errorf("expected Shell to be '/bin/bash', got '%s'", userInfo.Shell())
		}
		if userInfo.Gecos() != "vagrant" {
			t.Errorf("expected GECOS to be 'vagrant', got '%s'", userInfo.Gecos())
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
	diags := hostInfo.Populate(sshTransport)
	if diags.HasErrors() {
		t.Fatalf("failed to populate host info via SSH: %v", diags)
	}

	osInfo := hostInfo.OSInfo()
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

		if osInfo.ID() != "windows-server" {
			t.Errorf("expected OS ID to be 'windows-server', got '%s'", osInfo.ID())
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

		if osInfo.EditionID() != "datacenter" {
			t.Errorf("expected OS edition ID to be 'datacenter', got '%s'", osInfo.EditionID())
		}

		if osInfo.OSArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OSArch())
		}

		if osInfo.OSArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OSArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.SELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if selinuxInfo.Supported() {
			t.Error("expected SELinux to be unsupported on Windows")
		}
	}

	appArmorInfo := hostInfo.AppArmorInfo()
	if appArmorInfo == nil {
		t.Error("expected AppArmor info to be populated")
	} else {
		if appArmorInfo.Supported() {
			t.Error("expected AppArmor to be unsupported on Windows")
		}
	}

	fipsInfo := hostInfo.FipsInfo()
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

	// packageManagerInfo := hostInfo.PackageManagerInfo()
	// if packageManagerInfo == nil {
	// 	t.Error("expected Package Manager info to be populated")
	// } else {
	// 	// Blank for now, Windows is not supported yet
	// }

	serviceManagerInfo := hostInfo.ServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "windows-service-manager" {
			t.Errorf("expected Service Manager name to be 'windows-service-manager', got '%s'", serviceManagerInfo.Name())
		}
	}

	userInfo := hostInfo.UserInfo()
	if userInfo == nil {
		t.Error("expected User info to be populated")
	} else {
		if userInfo.Name() != "vagrant" {
			t.Errorf("expected User name to be 'vagrant', got '%s'", userInfo.Name())
		}
		if !strings.HasPrefix(userInfo.UserId(), "S-1-5-21-") || !strings.HasSuffix(userInfo.UserId(), "-1000") {
			t.Errorf("expected User ID to start with 'S-1-5-21-' and end with '-1000', got '%s'", userInfo.UserId())
		}
		if userInfo.GroupId() != "" {
			t.Errorf("expected Group ID to be empty on Windows, got '%s'", userInfo.GroupId())
		}
		if userInfo.HomeDir() != "C:\\Users\\vagrant" {
			t.Errorf("expected Home Directory to be 'C:\\Users\\vagrant', got '%s'", userInfo.HomeDir())
		}
		if userInfo.Shell() != "" {
			t.Errorf("expected Shell to be empty on Windows, got '%s'", userInfo.Shell())
		}
		if userInfo.Gecos() != "" {
			t.Errorf("expected GECOS to be empty on Windows, got '%s'", userInfo.Gecos())
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
	diags := hostInfo.Populate(sshTransport)
	if diags.HasErrors() {
		t.Fatalf("failed to populate host info via SSH: %v", diags)
	}

	osInfo := hostInfo.OSInfo()
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

		if osInfo.ID() != "windows-server" {
			t.Errorf("expected OS ID to be 'windows-server', got '%s'", osInfo.ID())
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

		if osInfo.EditionID() != "datacenter" {
			t.Errorf("expected OS edition ID to be 'datacenter', got '%s'", osInfo.EditionID())
		}

		if osInfo.OSArch() != "amd64" {
			t.Errorf("expected OS architecture to be 'amd64', got '%s'", osInfo.OSArch())
		}

		if osInfo.OSArchBits() != 64 {
			t.Errorf("expected OS architecture bits to be 64, got %d", osInfo.OSArchBits())
		}

		if osInfo.ProcArch() != "amd64" {
			t.Errorf("expected processor architecture to be 'amd64', got '%s'", osInfo.ProcArch())
		}

		if osInfo.ProcArchBits() != 64 {
			t.Errorf("expected processor architecture bits to be 64, got %d", osInfo.ProcArchBits())
		}
	}

	selinuxInfo := hostInfo.SELinuxInfo()
	if selinuxInfo == nil {
		t.Error("expected SELinux info to be populated")
	} else {
		if selinuxInfo.Supported() {
			t.Error("expected SELinux to be unsupported on Windows")
		}
	}

	appArmorInfo := hostInfo.AppArmorInfo()
	if appArmorInfo == nil {
		t.Error("expected AppArmor info to be populated")
	} else {
		if appArmorInfo.Supported() {
			t.Error("expected AppArmor to be unsupported on Windows")
		}
	}

	fipsInfo := hostInfo.FipsInfo()
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

	// packageManagerInfo := hostInfo.PackageManagerInfo()
	// if packageManagerInfo == nil {
	// 	t.Error("expected Package Manager info to be populated")
	// } else {
	// 	// Blank for now, Windows is not supported yet
	// }

	serviceManagerInfo := hostInfo.ServiceManagerInfo()
	if serviceManagerInfo == nil {
		t.Error("expected Service Manager info to be populated")
	} else {
		if serviceManagerInfo.Name() != "windows-service-manager" {
			t.Errorf("expected Service Manager name to be 'windows-service-manager', got '%s'", serviceManagerInfo.Name())
		}
	}

	userInfo := hostInfo.UserInfo()
	if userInfo == nil {
		t.Error("expected User info to be populated")
	} else {
		if userInfo.Name() != "vagrant" {
			t.Errorf("expected User name to be 'vagrant', got '%s'", userInfo.Name())
		}
		if !strings.HasPrefix(userInfo.UserId(), "S-1-5-21-") || !strings.HasSuffix(userInfo.UserId(), "-1000") {
			t.Errorf("expected User ID to start with 'S-1-5-21-' and end with '-1000', got '%s'", userInfo.UserId())
		}
		if userInfo.GroupId() != "" {
			t.Errorf("expected Group ID to be empty on Windows, got '%s'", userInfo.GroupId())
		}
		if userInfo.HomeDir() != "C:\\Users\\vagrant" {
			t.Errorf("expected Home Directory to be 'C:\\Users\\vagrant', got '%s'", userInfo.HomeDir())
		}
		if userInfo.Shell() != "" {
			t.Errorf("expected Shell to be empty on Windows, got '%s'", userInfo.Shell())
		}
		if userInfo.Gecos() != "" {
			t.Errorf("expected GECOS to be empty on Windows, got '%s'", userInfo.Gecos())
		}
	}
}
