package info

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestOSInfo_PopulateOSInfo_Darwin(t *testing.T) {

	tests := []struct {
		name                 string
		unameArchOutput      string
		swVersOutput         string
		expectedFriendlyName string
		expectedMajorVersion string
		expectedVersion      string
		expectedRelease      string
		expectedArch         string
	}{
		{
			name:                 "macOS Tahoe x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "26.0.0\n",
			expectedFriendlyName: "macOS 26.0.0",
			expectedMajorVersion: "26",
			expectedVersion:      "26.0.0",
			expectedRelease:      "Tahoe",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Sequoia x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "15.0.0\n",
			expectedFriendlyName: "macOS 15.0.0",
			expectedMajorVersion: "15",
			expectedVersion:      "15.0.0",
			expectedRelease:      "Sequoia",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Sonoma x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "14.0.0\n",
			expectedFriendlyName: "macOS 14.0.0",
			expectedMajorVersion: "14",
			expectedVersion:      "14.0.0",
			expectedRelease:      "Sonoma",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Ventura x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "13.0.0\n",
			expectedFriendlyName: "macOS 13.0.0",
			expectedMajorVersion: "13",
			expectedVersion:      "13.0.0",
			expectedRelease:      "Ventura",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Monterey x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "12.0.0\n",
			expectedFriendlyName: "macOS 12.0.0",
			expectedMajorVersion: "12",
			expectedVersion:      "12.0.0",
			expectedRelease:      "Monterey",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Big Sur x64",
			unameArchOutput:      "x86_64\n",
			swVersOutput:         "11.0.0\n",
			expectedFriendlyName: "macOS 11.0.0",
			expectedMajorVersion: "11",
			expectedVersion:      "11.0.0",
			expectedRelease:      "Big Sur",
			expectedArch:         "amd64",
		},
		{
			name:                 "macOS Tahoe arm64",
			unameArchOutput:      "arm64\n",
			swVersOutput:         "26.0.0\n",
			expectedFriendlyName: "macOS 26.0.0",
			expectedMajorVersion: "26",
			expectedVersion:      "26.0.0",
			expectedRelease:      "Tahoe",
			expectedArch:         "arm64",
		},
		{
			name:                 "macOS Ventura arm64",
			unameArchOutput:      "arm64\n",
			swVersOutput:         "13.0.0\n",
			expectedFriendlyName: "macOS 13.0.0",
			expectedMajorVersion: "13",
			expectedVersion:      "13.0.0",
			expectedRelease:      "Ventura",
			expectedArch:         "arm64",
		},
		{
			name:                 "macOS Monterey arm64",
			unameArchOutput:      "arm64\n",
			swVersOutput:         "12.0.0\n",
			expectedFriendlyName: "macOS 12.0.0",
			expectedMajorVersion: "12",
			expectedVersion:      "12.0.0",
			expectedRelease:      "Monterey",
			expectedArch:         "arm64",
		},
		{
			name:                 "macOS Big Sur arm64",
			unameArchOutput:      "arm64\n",
			swVersOutput:         "11.0.0\n",
			expectedFriendlyName: "macOS 11.0.0",
			expectedMajorVersion: "11",
			expectedVersion:      "11.0.0",
			expectedRelease:      "Big Sur",
			expectedArch:         "arm64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Darwin\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: tt.unameArchOutput,
			}
			transport.commandResponses["/usr/bin/sw_vers -productVersion"] = &commandResponse{
				stdout: tt.swVersOutput,
			}

			fileSystem := transport.fileSystem

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Darwin family, got: %v", err)
			}

			if !osInfo.families.Contains("posix") {
				t.Error("expected POSIX family to be added")
			}

			if !osInfo.families.Contains("darwin") {
				t.Error("expected Darwin family to be added")
			}

			if osInfo.id != "darwin" {
				t.Errorf("expected OS ID to be 'darwin', got: %s", osInfo.id)
			}

			if osInfo.friendlyName != tt.expectedFriendlyName {
				t.Errorf("expected friendly name to be '%s', got: %s", tt.expectedFriendlyName, osInfo.friendlyName)
			}

			if osInfo.release != tt.expectedRelease {
				t.Errorf("expected release to be '%s', got: %s", tt.expectedRelease, osInfo.release)
			}

			if osInfo.majorVersion != tt.expectedMajorVersion {
				t.Errorf("expected major version to be '%s', got: %s", tt.expectedMajorVersion, osInfo.majorVersion)
			}

			if osInfo.version != tt.expectedVersion {
				t.Errorf("expected version to be '%s', got: %s", tt.expectedVersion, osInfo.version)
			}

			if osInfo.osArch != tt.expectedArch {
				t.Errorf("expected OS architecture to be '%s', got: %s", tt.expectedArch, osInfo.osArch)
			}

			if osInfo.osArchBits != 64 {
				t.Errorf("expected OS architecture bits to be 64, got: %d", osInfo.osArchBits)
			}

			if osInfo.procArch != tt.expectedArch {
				t.Errorf("expected processor architecture to be '%s', got: %s", tt.expectedArch, osInfo.procArch)
			}

			if osInfo.procArchBits != 64 {
				t.Errorf("expected processor architecture bits to be 64, got: %d", osInfo.procArchBits)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Families(t *testing.T) {

	tests := []struct {
		name              string
		osReleaseContent  string
		lsbReleaseContent string
		expectedFamilies  []string
	}{
		{
			name: "almalinux from os-release",
			osReleaseContent: `PRETTY_NAME="AlmaLinux 8.5"
							NAME="AlmaLinux"
							VERSION_ID="8.5"
							VERSION="8.5"
							VERSION_CODENAME="almalinux"
							ID="almalinux"
							ID_LIKE="centos rhel"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "almalinux"},
		},
		{
			name: "almalinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: AlmaLinux
							Description:    AlmaLinux 8.5
							Release:        8.5
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "almalinux"},
		},
		{
			name: "amazon from os-release",
			osReleaseContent: `PRETTY_NAME="Amazon Linux 2"
							NAME="Amazon Linux"
							VERSION_ID="2"
							VERSION="2"
							VERSION_CODENAME="amzn"
							ID="amzn"
							ID_LIKE="centos rhel"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "amazon"},
		},
		{
			name: "amazon from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: amzn
							Description:    Amazon Linux 2
							Release:        2
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "amazon"},
		},
		{
			name: "archlinux-arm from os-release",
			osReleaseContent: `PRETTY_NAME="Arch Linux ARM"
							NAME="Arch Linux ARM"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="archarm"
							ID_LIKE="arch"
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "archlinux-arm"},
		},
		{
			name: "archlinux-arm from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: archarm
							Description:    Arch Linux ARM
							Release:        rolling
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "archlinux-arm"},
		},
		{
			name: "arcolinux from os-release",
			osReleaseContent: `PRETTY_NAME="ArcoLinux"
							NAME="ArcoLinux"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="arcolinux"
							ID_LIKE="arch"
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "arcolinux"},
		},
		{
			name: "arcolinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: arcolinux
							Description:    ArcoLinux
							Release:        rolling
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "arcolinux"},
		},
		{
			name: "centos from os-release",
			osReleaseContent: `PRETTY_NAME="CentOS Linux 7"
							NAME="CentOS Linux"
							VERSION_ID="7"
							VERSION="7"
							VERSION_CODENAME="centos"
							ID="centos"
							ID_LIKE="rhel fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "centos"},
		},
		{
			name: "centos from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: centos
							Description:    CentOS Linux 7
							Release:        7
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "centos"},
		},
		{
			name: "clearos from os-release",
			osReleaseContent: `PRETTY_NAME="ClearOS"
							NAME="ClearOS"
							VERSION_ID="7"
							VERSION="7"
							VERSION_CODENAME="clearos"
							ID="clearos"
							ID_LIKE="rhel fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "clearos"},
		},
		{
			name: "clearos from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: clearos
							Description:    ClearOS
							Release:        7
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "clearos"},
		},
		{
			name: "cloudlinux from os-release",
			osReleaseContent: `PRETTY_NAME="CloudLinux 7"
							NAME="CloudLinux"
							VERSION_ID="7"
							VERSION="7"
							VERSION_CODENAME="cloudlinux"
							ID="cloudlinux"
							ID_LIKE="rhel fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "cloudlinux"},
		},
		{
			name: "cloudlinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: cloudlinux
							Description:    CloudLinux 7
							Release:        7
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "cloudlinux"},
		},
		{
			name: "deepin from os-release",
			osReleaseContent: `PRETTY_NAME="Deepin 20.2"
							NAME="Deepin"
							VERSION_ID="20.2"
							VERSION="20.2"
							VERSION_CODENAME="n/a"
							ID="deepin"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "deepin"},
		},
		{
			name: "deepin from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: deepin
							Description:    Deepin
							Release:        20.2
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "deepin"},
		},
		{
			name: "devuan from os-release",
			osReleaseContent: `PRETTY_NAME="Devuan GNU/Linux 2.1 (ASCII)"
							NAME="Devuan"
							VERSION_ID="2.1"
							VERSION="2.1 (ASCII)"
							VERSION_CODENAME="ascii"
							ID="devuan"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "devuan"},
		},
		{
			name: "devuan from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: devuan
							Description:    Devuan GNU/Linux 2.1 (ASCII)
							Release:        2.1 (ASCII)
							Codename:       ascii
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "devuan"},
		},
		{
			name: "elementary from os-release",
			osReleaseContent: `PRETTY_NAME="elementary OS 6.1"
							NAME="elementary OS"
							VERSION_ID="6.1"
							VERSION="6.1"
							VERSION_CODENAME="juno"
							ID="elementary"
							ID_LIKE="ubuntu"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "elementary"},
		},
		{
			name: "elementary from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: elementary
							Description:    elementary OS 6.1
							Release:        6.1
							Codename:       juno
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "elementary"},
		},
		{
			name: "endeavouros from os-release",
			osReleaseContent: `PRETTY_NAME="EndeavourOS"
							NAME="EndeavourOS"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="endeavouros"
							ID_LIKE="arch"
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "endeavouros"},
		},
		{
			name: "endeavouros from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: endeavouros
							Description:    EndeavourOS
							Release:        rolling
							Codename:       rolling
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "endeavouros"},
		},
		{
			name: "fedora from os-release",
			osReleaseContent: `PRETTY_NAME="Fedora 34 (Workstation Edition)"
							NAME="Fedora"
							VERSION_ID="34"
							VERSION="34 (Workstation Edition)"
							VERSION_CODENAME="n/a"
							ID="fedora"
							ID_LIKE="rhel"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "fedora"},
		},
		{
			name: "fedora from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Fedora
							Description:    Fedora 34 (Workstation Edition)
							Release:        34
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "el", "fedora"},
		},
		{
			name: "kali from os-release",
			osReleaseContent: `PRETTY_NAME="Kali GNU/Linux 2021.3"
							NAME="Kali GNU/Linux"
							VERSION_ID="2021.3"
							VERSION="2021.3"
							VERSION_CODENAME="kali-rolling"
							ID="kali"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "kali"},
		},
		{
			name: "kali from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Kali
							Description:    Kali GNU/Linux 2021.3
							Release:        2021.3
							Codename:       kali-rolling
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "kali"},
		},
		{
			name: "kylin from os-release",
			osReleaseContent: `PRETTY_NAME="Kylin 10"
							NAME="Kylin"
							VERSION_ID="10"
							VERSION="10"
							VERSION_CODENAME="kylin"
							ID="kylin"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "kylin"},
		},
		{
			name: "kylin from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Kylin
							Description:    Kylin 10
							Release:        10
							Codename:       kylin
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "kylin"},
		},
		{
			name: "linuxmint from os-release",
			osReleaseContent: `PRETTY_NAME="Linux Mint 20.2"
							NAME="Linux Mint"
							VERSION_ID="20.2"
							VERSION="20.2"
							VERSION_CODENAME="uma"
							ID="linuxmint"
							ID_LIKE="ubuntu"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "linuxmint"},
		},
		{
			name: "linuxmint from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: LinuxMint
							Description:    Linux Mint 20.2
							Release:        20.2
							Codename:       uma
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "linuxmint"},
		},
		{
			name: "mageia from os-release",
			osReleaseContent: `PRETTY_NAME="Mageia 8"
							NAME="Mageia"
							VERSION_ID="8"
							VERSION="8"
							VERSION_CODENAME="n/a"
							ID="mageia"
							ID_LIKE="mandriva mandrake"
							`,
			expectedFamilies: []string{"posix", "linux", "mandrake", "mageia"},
		},
		{
			name: "mageia from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Mageia
							Description:    Mageia 8
							Release:        8
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "mandrake", "mageia"},
		},
		{
			name: "manjaro from os-release",
			osReleaseContent: `PRETTY_NAME="Manjaro Linux"
							NAME="Manjaro Linux"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="manjaro"
							ID_LIKE="arch"
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "manjaro"},
		},
		{
			name: "manjaro from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Manjaro
							Description:    Manjaro Linux
							Release:        rolling
							Codename:       rolling
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "manjaro"},
		},
		{
			name: "manjaro-arm from os-release",
			osReleaseContent: `PRETTY_NAME="Manjaro ARM"
							NAME="Manjaro ARM"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="manjaro-arm"
							ID_LIKE="arch"
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "manjaro", "manjaro-arm"},
		},
		{
			name: "manjaro-arm from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Manjaro-ARM
							Description:    Manjaro ARM
							Release:        rolling
							Codename:       rolling
							`,
			expectedFamilies: []string{"posix", "linux", "archlinux", "manjaro", "manjaro-arm"},
		},
		{
			name: "nobara from os-release",
			osReleaseContent: `PRETTY_NAME="Nobara 38"
							NAME="Nobara"
							VERSION_ID="38"
							VERSION="38"
							VERSION_CODENAME="nobara"
							ID="nobara"
							ID_LIKE="fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "fedora", "nobara"},
		},
		{
			name: "nobara from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Nobara
							Description:    Nobara 38
							Release:        38
							Codename:       nobara
							`,
			expectedFamilies: []string{"posix", "linux", "el", "fedora", "nobara"},
		},
		{
			name: "opensuse from os-release",
			osReleaseContent: `PRETTY_NAME="openSUSE Leap 15.3"
							NAME="openSUSE Leap"
							VERSION_ID="15.3"
							VERSION="15.3"
							VERSION_CODENAME="n/a"
							ID="opensuse-leap"
							ID_LIKE="suse"
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "opensuse"},
		},
		{
			name: "opensuse from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: openSUSE-leap
							Description:    openSUSE Leap 15.3
							Release:        15.3
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "opensuse"},
		},
		{
			name: "oraclelinux from os-release",
			osReleaseContent: `PRETTY_NAME="Oracle Linux Server 8.5"
							NAME="Oracle Linux Server"
							VERSION_ID="8.5"
							VERSION="8.5"
							VERSION_CODENAME="ol8"
							ID="ol"
							ID_LIKE="fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "oraclelinux"},
		},
		{
			name: "oraclelinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: ol
							Description:    Oracle Linux Server 8.5
							Release:        8.5
							Codename:       ol8
							`,
			expectedFamilies: []string{"posix", "linux", "el", "oraclelinux"},
		},
		{
			name: "pop_os from os-release",
			osReleaseContent: `PRETTY_NAME="Pop!_OS 21.04"
							NAME="Pop!_OS"
							VERSION_ID="21.04"
							VERSION="21.04"
							VERSION_CODENAME="hirsute"
							ID="pop"
							ID_LIKE="ubuntu"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "pop_os"},
		},
		{
			name: "pop_os from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Pop
							Description:    Pop!_OS 21.04
							Release:        21.04
							Codename:       hirsute
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu", "pop_os"},
		},
		{
			name: "raspbian from os-release",
			osReleaseContent: `PRETTY_NAME="Raspbian GNU/Linux 10 (buster)"
							NAME="Raspbian GNU/Linux"
							VERSION_ID="10"
							VERSION="10 (buster)"
							VERSION_CODENAME="buster"
							ID="raspbian"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "raspbian"},
		},
		{
			name: "raspbian from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Raspbian
							Description:    Raspbian GNU/Linux 10 (buster)
							Release:        10
							Codename:       buster
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "raspbian"},
		},
		{
			name: "rhel from os-release",
			osReleaseContent: `PRETTY_NAME="Red Hat Enterprise Linux 8.5"
							NAME="Red Hat Enterprise Linux"
							VERSION_ID="8.5"
							VERSION="8.5"
							VERSION_CODENAME="ol8"
							ID="rhel"
							ID_LIKE="fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "rhel"},
		},
		{
			name: "rhel from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: rhel
							Description:    Red Hat Enterprise Linux 8.5
							Release:        8.5
							Codename:       ol8
							`,
			expectedFamilies: []string{"posix", "linux", "el", "rhel"},
		},
		{
			name: "rocky from os-release",
			osReleaseContent: `PRETTY_NAME="Rocky Linux 8.5"
							NAME="Rocky Linux"
							VERSION_ID="8.5"
							VERSION="8.5"
							VERSION_CODENAME="rocky"
							ID="rocky"
							ID_LIKE="centos rhel"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "rocky"},
		},
		{
			name: "rocky from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: rocky
							Description:    Rocky Linux 8.5
							Release:        8.5
							Codename:       rocky
							`,
			expectedFamilies: []string{"posix", "linux", "el", "rocky"},
		},
		{
			name: "scientific from os-release",
			osReleaseContent: `PRETTY_NAME="Scientific Linux 8"
							NAME="Scientific Linux"
							VERSION_ID="8"
							VERSION="8"
							VERSION_CODENAME="scientific"
							ID="scientific"
							ID_LIKE="centos rhel"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "scientific"},
		},
		{
			name: "scientific from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: scientific
							Description:    Scientific Linux 8
							Release:        8
							Codename:       scientific
							`,
			expectedFamilies: []string{"posix", "linux", "el", "scientific"},
		},
		{
			name: "sled from os-release",
			osReleaseContent: `PRETTY_NAME="SUSE Linux Enterprise Desktop 15 SP3"
							NAME="SUSE Linux Enterprise Desktop"
							VERSION_ID="15.3"
							VERSION="15.3"
							VERSION_CODENAME="n/a"
							ID="sled"
							ID_LIKE="suse"
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "sled"},
		},
		{
			name: "sled from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: SLED
							Description:    SUSE Linux Enterprise Desktop 15 SP3
							Release:        15.3
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "sled"},
		},
		{
			name: "sles from os-release",
			osReleaseContent: `PRETTY_NAME="SUSE Linux Enterprise Server 15 SP3"
							NAME="SUSE Linux Enterprise Server"
							VERSION_ID="15.3"
							VERSION="15.3"
							VERSION_CODENAME="n/a"
							ID="sles"
							ID_LIKE="suse"
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "sles"},
		},
		{
			name: "sles from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: SLES
							Description:    SUSE Linux Enterprise Server 15 SP3
							Release:        15.3
							Codename:       n/a
							`,
			expectedFamilies: []string{"posix", "linux", "suse", "sles"},
		},
		{
			name: "ubuntu from os-release",
			osReleaseContent: `PRETTY_NAME="Ubuntu 20.04.3 LTS"
							NAME="Ubuntu"
							VERSION_ID="20.04"
							VERSION="20.04.3 LTS (Focal Fossa)"
							VERSION_CODENAME="focal"
							ID="ubuntu"
							ID_LIKE="debian"
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu"},
		},
		{
			name: "ubuntu from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Ubuntu
							Description:    Ubuntu 20.04.3 LTS
							Release:        20.04
							Codename:       focal
							`,
			expectedFamilies: []string{"posix", "linux", "debian", "ubuntu"},
		},
		{
			name: "virtuozzo from os-release",
			osReleaseContent: `PRETTY_NAME="Virtuozzo Linux 7"
							NAME="Virtuozzo Linux"
							VERSION_ID="7"
							VERSION="7"
							VERSION_CODENAME="virtuozzo"
							ID="virtuozzo"
							ID_LIKE="rhel fedora"
							`,
			expectedFamilies: []string{"posix", "linux", "el", "virtuozzo"},
		},
		{
			name: "virtuozzo from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Virtuozzo
							Description:    Virtuozzo Linux 7
							Release:        7
							Codename:       virtuozzo
							`,
			expectedFamilies: []string{"posix", "linux", "el", "virtuozzo"},
		},
		{
			name: "generic from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedFamilies: []string{"posix", "linux", "generic"},
		},
		{
			name: "generic from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Generic
							Description:    Generic Linux 1.0
							Release:        1.0
							Codename:       generic
							`,
			expectedFamilies: []string{"posix", "linux", "generic"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}
			transport.commandResponses["/usr/bin/lsb_release -a"] = &commandResponse{
				stdout: tt.lsbReleaseContent,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux family, got: %v", err)
			}

			for _, family := range tt.expectedFamilies {
				if !osInfo.families.Contains(family) {
					t.Errorf("expected family %q to be added, but it was not", family)
				}
			}

			if len(tt.expectedFamilies) != osInfo.families.Size() {
				t.Errorf("expected %d families, got: %d", len(tt.expectedFamilies), osInfo.families.Size())
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_ID(t *testing.T) {

	tests := []struct {
		name              string
		osReleaseContent  string
		lsbReleaseContent string
		expectedID        string
	}{
		{
			name: "amazon from os-release",
			osReleaseContent: `PRETTY_NAME="Amazon Linux 2"
							NAME="Amazon Linux"
							VERSION_ID="2"
							VERSION="2"
							VERSION_CODENAME="amzn"
							ID="amzn"
							ID_LIKE="centos rhel"
							`,
			expectedID: "amazon",
		},
		{
			name: "amazon from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: amzn
							Description:    Amazon Linux 2
							Release:        2
							Codename:       n/a
							`,
			expectedID: "amazon",
		},
		{
			name: "archlinux from os-release",
			osReleaseContent: `PRETTY_NAME="Arch Linux"
							NAME="Arch Linux"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="arch"
							ID_LIKE="arch"
							`,
			expectedID: "archlinux",
		},
		{
			name: "archlinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: arch
							Description:    Arch Linux
							Release:        rolling
							Codename:       rolling
							`,
			expectedID: "archlinux",
		},
		{
			name: "archlinux-arm from os-release",
			osReleaseContent: `PRETTY_NAME="Arch Linux ARM"
							NAME="Arch Linux ARM"
							VERSION_ID="rolling"
							VERSION="rolling"
							VERSION_CODENAME="rolling"
							ID="archarm"
							ID_LIKE="arch"
							`,
			expectedID: "archlinux-arm",
		},
		{
			name: "archlinux-arm from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: archarm
							Description:    Arch Linux ARM
							Release:        rolling
							Codename:       rolling
							`,
			expectedID: "archlinux-arm",
		},
		{
			name: "clearlinux from os-release",
			osReleaseContent: `PRETTY_NAME="Clear Linux OS"
							NAME="Clear Linux OS"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="clearlinux"
							ID="clear-linux-os"
							ID_LIKE="linux"
							`,
			expectedID: "clearlinux",
		},
		{
			name: "clearlinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: clear-linux-os
							Description:    Clear Linux OS 1.0
							Release:        1.0
							Codename:       clearlinux
							`,
			expectedID: "clearlinux",
		},
		{
			name: "cumuluslinux from os-release",
			osReleaseContent: `PRETTY_NAME="Cumulus Linux 3.7"
							NAME="Cumulus Linux"
							VERSION_ID="3.7"
							VERSION="3.7"
							VERSION_CODENAME="cumulus"
							ID="cumulus-linux"
							ID_LIKE="debian"
							`,
			expectedID: "cumuluslinux",
		},
		{
			name: "cumuluslinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: cumulus-linux
							Description:    Cumulus Linux 3.7
							Release:        3.7
							Codename:       cumulus
							`,
			expectedID: "cumuluslinux",
		},
		{
			name: "pop_os from os-release",
			osReleaseContent: `PRETTY_NAME="Pop!_OS 21.04"
							NAME="Pop!_OS"
							VERSION_ID="21.04"
							VERSION="21.04"
							VERSION_CODENAME="hirsute"
							ID="pop"
							ID_LIKE="ubuntu"
							`,
			expectedID: "pop_os",
		},
		{
			name: "pop_os from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Pop
							Description:    Pop!_OS 21.04
							Release:        21.04
							Codename:       hirsute
							`,
			expectedID: "pop_os",
		},
		{
			name: "oraclelinux from os-release",
			osReleaseContent: `PRETTY_NAME="Oracle Linux Server 8.5"
							NAME="Oracle Linux Server"
							VERSION_ID="8.5"
							VERSION="8.5"
							VERSION_CODENAME="ol8"
							ID="ol"
							ID_LIKE="fedora"
							`,
			expectedID: "oraclelinux",
		},
		{
			name: "oraclelinux from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: ol
							Description:    Oracle Linux Server 8.5
							Release:        8.5
							Codename:       ol8
							`,
			expectedID: "oraclelinux",
		},
		{
			name: "opensuse from os-release",
			osReleaseContent: `PRETTY_NAME="openSUSE Leap 15.3"
							NAME="openSUSE Leap"
							VERSION_ID="15.3"
							VERSION="15.3"
							VERSION_CODENAME="n/a"
							ID="opensuse-leap"
							ID_LIKE="suse"
							`,
			expectedID: "opensuse",
		},
		{
			name: "opensuse from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: openSUSE-leap
							Description:    openSUSE Leap 15.3
							Release:        15.3
							Codename:       n/a
							`,
			expectedID: "opensuse",
		},
		{
			name: "sles from os-release",
			osReleaseContent: `PRETTY_NAME="SUSE Linux Enterprise Server 15 SP3 for SAP Applications"
							NAME="SUSE Linux Enterprise Server for SAP Applications"
							VERSION_ID="15.3"
							VERSION="15.3"
							VERSION_CODENAME="n/a"
							ID="sles_sap"
							ID_LIKE="suse"
							`,
			expectedID: "sles",
		},
		{
			name: "sles from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: SLES_SAP
							Description:    SUSE Linux Enterprise Server 15 SP3 for SAP Applications
							Release:        15.3
							Codename:       n/a
							`,
			expectedID: "sles",
		},
		{
			name: "generic from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedID: "generic",
		},
		{
			name: "generic from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Generic
							Description:    Generic Linux 1.0
							Release:        1.0
							Codename:       generic
							`,
			expectedID: "generic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}
			transport.commandResponses["/usr/bin/lsb_release -a"] = &commandResponse{
				stdout: tt.lsbReleaseContent,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux OS, got: %v", err)
			}

			if osInfo.id != tt.expectedID {
				t.Errorf("expected ID %q, got: %q", tt.expectedID, osInfo.id)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_FriendlyName(t *testing.T) {

	tests := []struct {
		name                 string
		osReleaseContent     string
		lsbReleaseContent    string
		expectedFriendlyName string
	}{
		{
			name: "from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedFriendlyName: "Generic Linux 1.0",
		},
		{
			name: "from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Generic
							Description:    Generic Linux 1.0
							Release:        1.0
							Codename:       generic
							`,
			expectedFriendlyName: "Generic Linux 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}
			transport.commandResponses["/usr/bin/lsb_release -a"] = &commandResponse{
				stdout: tt.lsbReleaseContent,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux OS, got: %v", err)
			}

			if osInfo.friendlyName != tt.expectedFriendlyName {
				t.Errorf("expected friendly name %q, got: %q", tt.expectedFriendlyName, osInfo.friendlyName)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Release(t *testing.T) {

	tests := []struct {
		name              string
		osReleaseContent  string
		lsbReleaseContent string
		expectedRelease   string
	}{
		{
			name: "from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedRelease: "generic",
		},
		{
			name: "from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Generic
							Description:    Generic Linux 1.0
							Release:        1.0
							Codename:       generic
							`,
			expectedRelease: "generic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}
			transport.commandResponses["/usr/bin/lsb_release -a"] = &commandResponse{
				stdout: tt.lsbReleaseContent,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux OS, got: %v", err)
			}

			if osInfo.release != tt.expectedRelease {
				t.Errorf("expected release %q, got: %q", tt.expectedRelease, osInfo.release)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Version(t *testing.T) {

	tests := []struct {
		name                 string
		osReleaseContent     string
		lsbReleaseContent    string
		expectedMajorVersion string
		expectedVersion      string
	}{
		{
			name: "from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
		},
		{
			name: "from lsb-release",
			lsbReleaseContent: `LSB Version:    n/a
							Distributor ID: Generic
							Description:    Generic Linux 1.0
							Release:        1.0
							Codename:       generic
							`,
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}
			transport.commandResponses["/usr/bin/lsb_release -a"] = &commandResponse{
				stdout: tt.lsbReleaseContent,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux OS, got: %v", err)
			}

			if osInfo.majorVersion != tt.expectedMajorVersion {
				t.Errorf("expected major version %q, got: %q", tt.expectedMajorVersion, osInfo.majorVersion)
			}

			if osInfo.version != tt.expectedVersion {
				t.Errorf("expected version %q, got: %q", tt.expectedVersion, osInfo.version)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Edition(t *testing.T) {

	tests := []struct {
		name              string
		osReleaseContent  string
		expectedEdition   string
		expectedEditionId string
	}{
		{
			name: "from os-release",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							VARIANT="Generic"
							VARIANT_ID="generic"
							`,
			expectedEdition:   "Generic",
			expectedEditionId: "generic",
		},
		{
			name: "not specified",
			osReleaseContent: `PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME="generic"
							ID="generic"
							ID_LIKE="linux"
							`,
			expectedEdition:   "",
			expectedEditionId: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: "x86_64\n",
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(tt.osReleaseContent)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux OS, got: %v", err)
			}

			if osInfo.edition != tt.expectedEdition {
				t.Errorf("expected edition %q, got: %q", tt.expectedEdition, osInfo.edition)
			}

			if osInfo.editionId != tt.expectedEditionId {
				t.Errorf("expected edition ID %q, got: %q", tt.expectedEditionId, osInfo.editionId)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Architecture(t *testing.T) {

	tests := []struct {
		name             string
		unameArchOutput  string
		expectedArch     string
		expectedArchBits int
	}{
		{
			name:             "i386",
			unameArchOutput:  "i386\n",
			expectedArch:     "386",
			expectedArchBits: 32,
		},
		{
			name:             "i486",
			unameArchOutput:  "i486\n",
			expectedArch:     "386",
			expectedArchBits: 32,
		},
		{
			name:             "i586",
			unameArchOutput:  "i586\n",
			expectedArch:     "386",
			expectedArchBits: 32,
		},
		{
			name:             "i686",
			unameArchOutput:  "i686\n",
			expectedArch:     "386",
			expectedArchBits: 32,
		},
		{
			name:             "x86_64",
			unameArchOutput:  "x86_64\n",
			expectedArch:     "amd64",
			expectedArchBits: 64,
		},
		{
			name:             "armv6l",
			unameArchOutput:  "armv6l\n",
			expectedArch:     "arm",
			expectedArchBits: 32,
		},
		{
			name:             "armv7l",
			unameArchOutput:  "armv7l\n",
			expectedArch:     "arm",
			expectedArchBits: 32,
		},
		{
			name:             "aarch64",
			unameArchOutput:  "aarch64\n",
			expectedArch:     "arm64",
			expectedArchBits: 64,
		},
		{
			name:             "mips",
			unameArchOutput:  "mips\n",
			expectedArch:     "mips",
			expectedArchBits: 32,
		},
		{
			name:             "mips64",
			unameArchOutput:  "mips64\n",
			expectedArch:     "mips64",
			expectedArchBits: 64,
		},
		{
			name:             "ppc64",
			unameArchOutput:  "ppc64\n",
			expectedArch:     "ppc64",
			expectedArchBits: 64,
		},
		{
			name:             "ppc64le",
			unameArchOutput:  "ppc64le\n",
			expectedArch:     "ppc64le",
			expectedArchBits: 64,
		},
		{
			name:             "riscv64",
			unameArchOutput:  "riscv64\n",
			expectedArch:     "riscv64",
			expectedArchBits: 64,
		},
		{
			name:             "s390x",
			unameArchOutput:  "s390x\n",
			expectedArch:     "s390x",
			expectedArchBits: 64,
		},
		{
			name:             "newarch",
			unameArchOutput:  "newarch\n",
			expectedArch:     "newarch",
			expectedArchBits: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				stdout: "Linux\n",
			}
			transport.commandResponses["uname -m"] = &commandResponse{
				stdout: tt.unameArchOutput,
			}

			fileSystem := transport.fileSystem
			fileSystem.files["/etc/os-release"] = &mockFile{
				content: io.NopCloser(bytes.NewBufferString(`PRETTY_NAME="Generic Linux 1.0"
							NAME="Generic Linux"
							VERSION_ID="1.0"
							VERSION="1.0"
							VERSION_CODENAME=generic
							ID=generic
							ID_LIKE=debian
							`)),
			}

			err := osInfo.populateOSInfo(transport, fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Linux family, got: %v", err)
			}

			if osInfo.osArch != tt.expectedArch {
				t.Errorf("expected OS architecture to be %q, got: %s", tt.expectedArch, osInfo.osArch)
			}

			if osInfo.osArchBits != tt.expectedArchBits {
				t.Errorf("expected OS architecture bits to be %d, got: %d", tt.expectedArchBits, osInfo.osArchBits)
			}

			if osInfo.procArch != tt.expectedArch {
				t.Errorf("expected processor architecture to be %q, got: %s", tt.expectedArch, osInfo.procArch)
			}

			if osInfo.procArchBits != tt.expectedArchBits {
				t.Errorf("expected processor architecture bits to be %d, got: %d", tt.expectedArchBits, osInfo.procArchBits)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Windows_Architecture(t *testing.T) {

	tests := []struct {
		name string

		procArchPowerShell string
		osArchPowerShell   string

		expectedprocArch     string
		expectedprocArchBits int

		expectedosArch     string
		expectedosArchBits int
	}{
		{
			name:                 "x86",
			procArchPowerShell:   "x86",
			osArchPowerShell:     "32-bit",
			expectedprocArch:     "386",
			expectedprocArchBits: 32,
			expectedosArch:       "386",
			expectedosArchBits:   32,
		},
		{
			name:                 "x64",
			procArchPowerShell:   "AMD64",
			osArchPowerShell:     "64-bit",
			expectedprocArch:     "amd64",
			expectedprocArchBits: 64,
			expectedosArch:       "amd64",
			expectedosArchBits:   64,
		},
		{
			name:                 "32-bit OS on x64 processor",
			procArchPowerShell:   "AMD64",
			osArchPowerShell:     "32-bit",
			expectedprocArch:     "amd64",
			expectedprocArchBits: 64,
			expectedosArch:       "386",
			expectedosArchBits:   32,
		},
		{
			name:                 "arm64",
			procArchPowerShell:   "ARM64",
			osArchPowerShell:     "64-bit",
			expectedprocArch:     "arm64",
			expectedprocArchBits: 64,
			expectedosArch:       "arm64",
			expectedosArchBits:   64,
		},
		{
			name:                 "arm",
			procArchPowerShell:   "ARM",
			osArchPowerShell:     "32-bit",
			expectedprocArch:     "arm",
			expectedprocArchBits: 32,
			expectedosArch:       "arm",
			expectedosArchBits:   32,
		},
		{
			name:                 "32-bit OS on arm64 processor",
			procArchPowerShell:   "ARM64",
			osArchPowerShell:     "32-bit",
			expectedprocArch:     "arm64",
			expectedprocArchBits: 64,
			expectedosArch:       "arm",
			expectedosArchBits:   32,
		},
		{
			name:                 "unknown architecture",
			procArchPowerShell:   "newarch",
			osArchPowerShell:     "64-bit",
			expectedprocArch:     "newarch",
			expectedprocArchBits: 0,
			expectedosArch:       "newarch",
			expectedosArchBits:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				err: errors.New("command not found"),
			}

			transport.powerShellResponses[procArchPowerShell] = &commandResponse{
				stdout: tt.procArchPowerShell,
			}
			transport.powerShellResponses[osArchPowerShell] = &commandResponse{
				stdout: tt.osArchPowerShell,
			}

			err := osInfo.populateOSInfo(transport, transport.fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Windows family, got: %v", err)
			}

			if osInfo.procArch != tt.expectedprocArch {
				t.Errorf("expected processor architecture to be %q, got: %s", tt.expectedprocArch, osInfo.procArch)
			}

			if osInfo.procArchBits != tt.expectedprocArchBits {
				t.Errorf("expected processor architecture bits to be %d, got: %d", tt.expectedprocArchBits, osInfo.procArchBits)
			}

			if osInfo.osArch != tt.expectedosArch {
				t.Errorf("expected OS architecture to be %q, got: %s", tt.expectedosArch, osInfo.osArch)
			}

			if osInfo.osArchBits != tt.expectedosArchBits {
				t.Errorf("expected OS architecture bits to be %d, got: %d", tt.expectedosArchBits, osInfo.osArchBits)
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Windows(t *testing.T) {

	tests := []struct {
		name                 string
		friendNamePowerShell string
		versionPowerShell    string
		expectedID           string
		expectedFriendlyName string
		expectedRelease      string
		expectedMajorVersion string
		expectedVersion      string
		expectedEdition      string
		expectedEditionId    string
	}{
		// Windows Server 2008 R2 (6.1.7600)
		{
			name:                 "Windows Server 2008 R2 Standard",
			friendNamePowerShell: "Microsoft Windows Server 2008 R2 Standard",
			versionPowerShell:    "6.1.7600.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2008 R2 Standard",
			expectedRelease:      "server-2008-r2",
			expectedMajorVersion: "6",
			expectedVersion:      "6.1.7600.0",
			expectedEdition:      "Standard",
			expectedEditionId:    "standard",
		},
		// Windows 7 (6.1.7600)
		{
			name:                 "Windows 7 Professional",
			friendNamePowerShell: "Microsoft Windows 7 Professional",
			versionPowerShell:    "6.1.7600.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 7 Professional",
			expectedRelease:      "7",
			expectedMajorVersion: "6",
			expectedVersion:      "6.1.7600.0",
			expectedEdition:      "Professional",
			expectedEditionId:    "professional",
		},
		// Windows Server 2008 R2 SP1 (6.1.7601)
		{
			name:                 "Windows Server 2008 R2 SP1 Enterprise",
			friendNamePowerShell: "Microsoft Windows Server 2008 R2 Enterprise",
			versionPowerShell:    "6.1.7601.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2008 R2 SP1 Enterprise",
			expectedRelease:      "server-2008-r2-sp1",
			expectedMajorVersion: "6",
			expectedVersion:      "6.1.7601.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows 7 SP1 (6.1.7601)
		{
			name:                 "Windows 7 SP1 Ultimate",
			friendNamePowerShell: "Microsoft Windows 7 Ultimate",
			versionPowerShell:    "6.1.7601.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 7 SP1 Ultimate",
			expectedRelease:      "7-sp1",
			expectedMajorVersion: "6",
			expectedVersion:      "6.1.7601.0",
			expectedEdition:      "Ultimate",
			expectedEditionId:    "ultimate",
		},
		// Windows Server 2012 (6.2.9200)
		{
			name:                 "Windows Server 2012 Datacenter",
			friendNamePowerShell: "Microsoft Windows Server 2012 Datacenter",
			versionPowerShell:    "6.2.9200.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2012 Datacenter",
			expectedRelease:      "server-2012",
			expectedMajorVersion: "6",
			expectedVersion:      "6.2.9200.0",
			expectedEdition:      "Datacenter",
			expectedEditionId:    "datacenter",
		},
		// Windows 8 (6.2.9200)
		{
			name:                 "Windows 8 Pro",
			friendNamePowerShell: "Microsoft Windows 8 Pro",
			versionPowerShell:    "6.2.9200.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 8 Pro",
			expectedRelease:      "8",
			expectedMajorVersion: "6",
			expectedVersion:      "6.2.9200.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows Server 2012 R2 (6.3.9600)
		{
			name:                 "Windows Server 2012 R2 Standard",
			friendNamePowerShell: "Microsoft Windows Server 2012 R2 Standard",
			versionPowerShell:    "6.3.9600.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2012 R2 Standard",
			expectedRelease:      "server-2012-r2",
			expectedMajorVersion: "6",
			expectedVersion:      "6.3.9600.0",
			expectedEdition:      "Standard",
			expectedEditionId:    "standard",
		},
		// Windows 8.1 (6.3.9600)
		{
			name:                 "Windows 8.1 Enterprise",
			friendNamePowerShell: "Microsoft Windows 8.1 Enterprise",
			versionPowerShell:    "6.3.9600.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 8.1 Enterprise",
			expectedRelease:      "8.1",
			expectedMajorVersion: "6",
			expectedVersion:      "6.3.9600.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows 10 1507 (10.0.10240)
		{
			name:                 "Windows 10 1507 Home",
			friendNamePowerShell: "Microsoft Windows 10 Home",
			versionPowerShell:    "10.0.10240.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1507 Home",
			expectedRelease:      "10-1507",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.10240.0",
			expectedEdition:      "Home",
			expectedEditionId:    "home",
		},
		// Windows 10 1511 (10.0.10586)
		{
			name:                 "Windows 10 1511 Pro",
			friendNamePowerShell: "Microsoft Windows 10 Pro",
			versionPowerShell:    "10.0.10586.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1511 Pro",
			expectedRelease:      "10-1511",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.10586.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows Server 2016 (10.0.14393)
		{
			name:                 "Windows Server 2016 Datacenter",
			friendNamePowerShell: "Microsoft Windows Server 2016 Datacenter",
			versionPowerShell:    "10.0.14393.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2016 Datacenter",
			expectedRelease:      "server-2016",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.14393.0",
			expectedEdition:      "Datacenter",
			expectedEditionId:    "datacenter",
		},
		// Windows 10 1607 (10.0.14393)
		{
			name:                 "Windows 10 1607 Enterprise",
			friendNamePowerShell: "Microsoft Windows 10 Enterprise",
			versionPowerShell:    "10.0.14393.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1607 Enterprise",
			expectedRelease:      "10-1607",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.14393.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows 10 1703 (10.0.15063)
		{
			name:                 "Windows 10 1703 Education",
			friendNamePowerShell: "Microsoft Windows 10 Education",
			versionPowerShell:    "10.0.15063.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1703 Education",
			expectedRelease:      "10-1703",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.15063.0",
			expectedEdition:      "Education",
			expectedEditionId:    "education",
		},
		// Windows 10 1709 (10.0.16299)
		{
			name:                 "Windows 10 1709 Pro",
			friendNamePowerShell: "Microsoft Windows 10 Pro",
			versionPowerShell:    "10.0.16299.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1709 Pro",
			expectedRelease:      "10-1709",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.16299.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows 10 1803 (10.0.17134)
		{
			name:                 "Windows 10 1803 Home",
			friendNamePowerShell: "Microsoft Windows 10 Home",
			versionPowerShell:    "10.0.17134.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1803 Home",
			expectedRelease:      "10-1803",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.17134.0",
			expectedEdition:      "Home",
			expectedEditionId:    "home",
		},
		// Windows Server 2019 (10.0.17763)
		{
			name:                 "Windows Server 2019 Standard",
			friendNamePowerShell: "Microsoft Windows Server 2019 Standard",
			versionPowerShell:    "10.0.17763.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2019 Standard",
			expectedRelease:      "server-2019",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.17763.0",
			expectedEdition:      "Standard",
			expectedEditionId:    "standard",
		},
		// Windows 10 1809 (10.0.17763)
		{
			name:                 "Windows 10 1809 Enterprise",
			friendNamePowerShell: "Microsoft Windows 10 Enterprise",
			versionPowerShell:    "10.0.17763.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1809 Enterprise",
			expectedRelease:      "10-1809",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.17763.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows 10 1903 (10.0.18362)
		{
			name:                 "Windows 10 1903 Pro",
			friendNamePowerShell: "Microsoft Windows 10 Pro",
			versionPowerShell:    "10.0.18362.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1903 Pro",
			expectedRelease:      "10-1903",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.18362.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows Server 1909 (10.0.18363)
		{
			name:                 "Windows Server 1909 Core",
			friendNamePowerShell: "Microsoft Windows Server Core",
			versionPowerShell:    "10.0.18363.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 1909 Core",
			expectedRelease:      "server-1909",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.18363.0",
			expectedEdition:      "Core",
			expectedEditionId:    "core",
		},
		// Windows 10 1909 (10.0.18363)
		{
			name:                 "Windows 10 1909 Home",
			friendNamePowerShell: "Microsoft Windows 10 Home",
			versionPowerShell:    "10.0.18363.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 1909 Home",
			expectedRelease:      "10-1909",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.18363.0",
			expectedEdition:      "Home",
			expectedEditionId:    "home",
		},
		// Windows Server 2004 (10.0.19041)
		{
			name:                 "Windows Server 2004 Datacenter",
			friendNamePowerShell: "Microsoft Windows Server 2004 Datacenter",
			versionPowerShell:    "10.0.19041.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2004 Datacenter",
			expectedRelease:      "server-2004",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19041.0",
			expectedEdition:      "Datacenter",
			expectedEditionId:    "datacenter",
		},
		// Windows 10 2004 (10.0.19041)
		{
			name:                 "Windows 10 2004 Education",
			friendNamePowerShell: "Microsoft Windows 10 Education",
			versionPowerShell:    "10.0.19041.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 2004 Education",
			expectedRelease:      "10-2004",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19041.0",
			expectedEdition:      "Education",
			expectedEditionId:    "education",
		},
		// Windows Server 20H2 (10.0.19042)
		{
			name:                 "Windows Server 20H2 Standard",
			friendNamePowerShell: "Microsoft Windows Server 20H2 Standard",
			versionPowerShell:    "10.0.19042.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 20H2 Standard",
			expectedRelease:      "server-20h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19042.0",
			expectedEdition:      "Standard",
			expectedEditionId:    "standard",
		},
		// Windows 10 20H2 (10.0.19042)
		{
			name:                 "Windows 10 20H2 Pro",
			friendNamePowerShell: "Microsoft Windows 10 Pro",
			versionPowerShell:    "10.0.19042.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 20H2 Pro",
			expectedRelease:      "10-20h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19042.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows 10 21H1 (10.0.19043)
		{
			name:                 "Windows 10 21H1 Enterprise",
			friendNamePowerShell: "Microsoft Windows 10 Enterprise",
			versionPowerShell:    "10.0.19043.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 21H1 Enterprise",
			expectedRelease:      "10-21h1",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19043.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows 10 21H2 (10.0.19044)
		{
			name:                 "Windows 10 21H2 Home",
			friendNamePowerShell: "Microsoft Windows 10 Home",
			versionPowerShell:    "10.0.19044.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 21H2 Home",
			expectedRelease:      "10-21h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19044.0",
			expectedEdition:      "Home",
			expectedEditionId:    "home",
		},
		// Windows 10 22H2 (10.0.19045)
		{
			name:                 "Windows 10 22H2 Pro",
			friendNamePowerShell: "Microsoft Windows 10 Pro",
			versionPowerShell:    "10.0.19045.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 10 22H2 Pro",
			expectedRelease:      "10-22h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.19045.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows Server 2022 (10.0.20348)
		{
			name:                 "Windows Server 2022 Datacenter",
			friendNamePowerShell: "Microsoft Windows Server 2022 Datacenter",
			versionPowerShell:    "10.0.20348.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2022 Datacenter",
			expectedRelease:      "server-2022",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.20348.0",
			expectedEdition:      "Datacenter",
			expectedEditionId:    "datacenter",
		},
		// Windows 11 21H2 (10.0.22000)
		{
			name:                 "Windows 11 21H2 Home",
			friendNamePowerShell: "Microsoft Windows 11 Home",
			versionPowerShell:    "10.0.22000.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 11 21H2 Home",
			expectedRelease:      "11-21h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.22000.0",
			expectedEdition:      "Home",
			expectedEditionId:    "home",
		},
		// Windows 11 22H2 (10.0.22621)
		{
			name:                 "Windows 11 22H2 Pro",
			friendNamePowerShell: "Microsoft Windows 11 Pro",
			versionPowerShell:    "10.0.22621.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 11 22H2 Pro",
			expectedRelease:      "11-22h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.22621.0",
			expectedEdition:      "Pro",
			expectedEditionId:    "pro",
		},
		// Windows 11 23H2 (10.0.22631)
		{
			name:                 "Windows 11 23H2 Enterprise",
			friendNamePowerShell: "Microsoft Windows 11 Enterprise",
			versionPowerShell:    "10.0.22631.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 11 23H2 Enterprise",
			expectedRelease:      "11-23h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.22631.0",
			expectedEdition:      "Enterprise",
			expectedEditionId:    "enterprise",
		},
		// Windows Server 23H2 (10.0.25398)
		{
			name:                 "Windows Server 23H2 Standard",
			friendNamePowerShell: "Microsoft Windows Server 23H2 Standard",
			versionPowerShell:    "10.0.25398.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 23H2 Standard",
			expectedRelease:      "server-23h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.25398.0",
			expectedEdition:      "Standard",
			expectedEditionId:    "standard",
		},
		// Windows Server 2025 (10.0.26100)
		{
			name:                 "Windows Server 2025 Datacenter Evaluation",
			friendNamePowerShell: "Microsoft Windows Server 2025 Datacenter Evaluation",
			versionPowerShell:    "10.0.26100.0",
			expectedID:           "windows-server",
			expectedFriendlyName: "Microsoft Windows Server 2025 Datacenter",
			expectedRelease:      "server-2025",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.26100.0",
			expectedEdition:      "Datacenter",
			expectedEditionId:    "datacenter",
		},
		// Windows 11 24H2 (10.0.26100)
		{
			name:                 "Windows 11 24H2 Education",
			friendNamePowerShell: "Microsoft Windows 11 Education",
			versionPowerShell:    "10.0.26100.0",
			expectedID:           "windows-client",
			expectedFriendlyName: "Microsoft Windows 11 24H2 Education",
			expectedRelease:      "11-24h2",
			expectedMajorVersion: "10",
			expectedVersion:      "10.0.26100.0",
			expectedEdition:      "Education",
			expectedEditionId:    "education",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			osInfo := newOSInfo()

			transport := newMockTransport()
			transport.commandResponses["uname -s"] = &commandResponse{
				err: errors.New("command not found"),
			}

			transport.powerShellResponses[osFriendlyNamePowerShell] = &commandResponse{
				stdout: tt.friendNamePowerShell + "\n",
			}
			transport.powerShellResponses[osVersionPowerShell] = &commandResponse{
				stdout: tt.versionPowerShell + "\n",
			}
			transport.powerShellResponses[procArchPowerShell] = &commandResponse{
				stdout: "AMD64\n",
			}
			transport.powerShellResponses[osArchPowerShell] = &commandResponse{
				stdout: "64-bit\n",
			}

			err := osInfo.populateOSInfo(transport, transport.fileSystem)
			if err != nil {
				t.Fatalf("expected no error for Windows family, got: %v", err)
			}

			if !osInfo.families.Contains("windows") {
				t.Error("expected 'windows' family to be added, but it was not")
			}

			if !osInfo.families.Contains(tt.expectedID) {
				t.Errorf("expected family %q to be added, but it was not", tt.expectedID)
			}

			if osInfo.id != tt.expectedID {
				t.Errorf("expected OS ID to be %q, got: %s", tt.expectedID, osInfo.id)
			}

			if osInfo.friendlyName != tt.expectedFriendlyName {
				t.Errorf("expected friendly name to be %q, got: %s", tt.expectedFriendlyName, osInfo.friendlyName)
			}

			if osInfo.release != tt.expectedRelease {
				t.Errorf("expected release to be %q, got: %s", tt.expectedRelease, osInfo.release)
			}

			if osInfo.majorVersion != tt.expectedMajorVersion {
				t.Errorf("expected major version to be %q, got: %s", tt.expectedMajorVersion, osInfo.majorVersion)
			}

			if osInfo.version != tt.expectedVersion {
				t.Errorf("expected version to be %q, got: %s", tt.expectedVersion, osInfo.version)
			}

			if osInfo.edition != tt.expectedEdition {
				t.Errorf("expected edition to be %q, got: %s", tt.expectedEdition, osInfo.edition)
			}

			if osInfo.editionId != tt.expectedEditionId {
				t.Errorf("expected edition ID to be %q, got: %s", tt.expectedEditionId, osInfo.editionId)
			}
		})
	}
}

func TestOSInfo_ToMapOfCtyValues(t *testing.T) {
	osInfo := newOSInfo()
	osInfo.families.Add("linux")
	osInfo.families.Add("debian")
	osInfo.families.Add("ubuntu")
	osInfo.id = "ubuntu"
	osInfo.friendlyName = "Ubuntu 22.04.3 LTS"
	osInfo.release = "jammy"
	osInfo.majorVersion = "22"
	osInfo.version = "22.04"
	osInfo.edition = "LTS"
	osInfo.editionId = "lts"
	osInfo.osArch = "amd64"
	osInfo.osArchBits = 64
	osInfo.procArch = "amd64"
	osInfo.procArchBits = 64

	values := osInfo.toMapOfCtyValues()

	if values["os_families"].Type() != cty.Set(cty.String) {
		t.Errorf("expected os_families to be a set of strings, got %s", values["os_families"].Type().GoString())
	}

	families := values["os_families"].AsValueSlice()

	if len(families) != 3 {
		t.Errorf("expected 3 families, got %d", len(families))
	}

	for _, family := range families {
		if family.Type() != cty.String {
			t.Errorf("expected family to be a string, got %s", family.Type().GoString())
		}

		if family.AsString() != "linux" && family.AsString() != "debian" && family.AsString() != "ubuntu" {
			t.Errorf("unexpected family value: %s", family.AsString())
		}
	}

	if values["os_id"].Type() != cty.String {
		t.Errorf("expected os_id to be a string, got %s", values["os_id"].Type().GoString())
	}
	if values["os_id"].AsString() != "ubuntu" {
		t.Errorf("expected os_id to be 'ubuntu', got %s", values["os_id"].AsString())
	}

	if values["os_friendly_name"].Type() != cty.String {
		t.Errorf("expected os_friendly_name to be a string, got %s", values["os_friendly_name"].Type().GoString())
	}
	if values["os_friendly_name"].AsString() != "Ubuntu 22.04.3 LTS" {
		t.Errorf("expected os_friendly_name to be 'Ubuntu 22.04.3 LTS', got %s", values["os_friendly_name"].AsString())
	}

	if values["os_release"].Type() != cty.String {
		t.Errorf("expected os_release to be a string, got %s", values["os_release"].Type().GoString())
	}
	if values["os_release"].AsString() != "jammy" {
		t.Errorf("expected os_release to be 'jammy', got %s", values["os_release"].AsString())
	}

	if values["os_major_version"].Type() != cty.String {
		t.Errorf("expected os_major_version to be a string, got %s", values["os_major_version"].Type().GoString())
	}
	if values["os_major_version"].AsString() != "22" {
		t.Errorf("expected os_major_version to be '22', got %s", values["os_major_version"].AsString())
	}

	if values["os_version"].Type() != cty.String {
		t.Errorf("expected os_version to be a string, got %s", values["os_version"].Type().GoString())
	}
	if values["os_version"].AsString() != "22.04" {
		t.Errorf("expected os_version to be '22.04', got %s", values["os_version"].AsString())
	}

	if values["os_edition"].Type() != cty.String {
		t.Errorf("expected os_edition to be a string, got %s", values["os_edition"].Type().GoString())
	}
	if values["os_edition"].AsString() != "LTS" {
		t.Errorf("expected os_edition to be 'LTS', got %s", values["os_edition"].AsString())
	}

	if values["os_edition_id"].Type() != cty.String {
		t.Errorf("expected os_edition_id to be a string, got %s", values["os_edition_id"].Type().GoString())
	}
	if values["os_edition_id"].AsString() != "lts" {
		t.Errorf("expected os_edition_id to be 'lts', got %s", values["os_edition_id"].AsString())
	}

	if values["os_architecture"].Type() != cty.String {
		t.Errorf("expected os_architecture to be a string, got %s", values["os_architecture"].Type().GoString())
	}
	if values["os_architecture"].AsString() != "amd64" {
		t.Errorf("expected os_architecture to be 'amd64', got %s", values["os_architecture"].AsString())
	}

	if values["processor_architecture"].Type() != cty.String {
		t.Errorf("expected processor_architecture to be a string, got %s", values["processor_architecture"].Type().GoString())
	}
	if values["processor_architecture"].AsString() != "amd64" {
		t.Errorf("expected processor_architecture to be 'amd64', got %s", values["processor_architecture"].AsString())
	}

	if values["os_architecture_bits"].Type() != cty.Number {
		t.Errorf("expected os_architecture_bits to be a number, got %s", values["os_architecture_bits"].Type().GoString())
	}
	value, _ := values["os_architecture_bits"].AsBigFloat().Int64()
	if value != 64 {
		t.Errorf("expected os_architecture_bits to be 64, got %s", values["os_architecture_bits"].AsString())
	}

	if values["processor_architecture_bits"].Type() != cty.Number {
		t.Errorf("expected processor_architecture_bits to be a number, got %s", values["processor_architecture_bits"].Type().GoString())
	}
	value, _ = values["processor_architecture_bits"].AsBigFloat().Int64()
	if value != 64 {
		t.Errorf("expected processor_architecture_bits to be 64, got %s", values["processor_architecture_bits"].AsString())
	}
}

func TestOSInfo_ToMapOfCtyValues_EmptyValues(t *testing.T) {
	osInfo := newOSInfo()

	values := osInfo.toMapOfCtyValues()

	numberKeys := []string{
		"os_architecture_bits",
		"processor_architecture_bits",
	}

	stringKeys := []string{
		"os_id",
		"os_friendly_name",
		"os_release",
		"os_major_version",
		"os_version",
		"os_edition",
		"os_edition_id",
		"os_architecture",
		"processor_architecture",
	}

	setOfStringsKeys := []string{
		"os_families",
	}

	for _, key := range numberKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected %s to be null, got %s", key, value.GoString())
			}
			if value.Type() != cty.Number {
				t.Errorf("expected %s to be of type Number, got %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	for _, key := range stringKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected %s to be null, got %s", key, value.GoString())
			}
			if value.Type() != cty.String {
				t.Errorf("expected %s to be of type String, got %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	for _, key := range setOfStringsKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected %s to be null, got %s", key, value.GoString())
			}
			if value.Type() != cty.Set(cty.String) {
				t.Errorf("expected %s to be of type Set(String), got %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}
}
