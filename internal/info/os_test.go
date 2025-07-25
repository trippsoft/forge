package info

import (
	"runtime"
	"testing"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestOSInfo_PopulateOSInfo_WithLocalTransport(t *testing.T) {
	localTransport, err := transport.NewLocalTransport()
	if err != nil {
		t.Fatalf("failed to create local transport: %v", err)
	}
	defer localTransport.Close()

	err = localTransport.Connect()
	if err != nil {
		t.Fatalf("failed to connect local transport: %v", err)
	}

	osInfo := newOSInfo()
	err = osInfo.populateOSInfo(localTransport, localTransport.FileSystem())
	if err != nil {
		t.Fatalf("failed to populate OS info: %v", err)
	}

	// Verify families are populated
	if osInfo.families == nil || osInfo.families.Size() == 0 {
		t.Error("expected OS families to be populated")
	}

	// Verify architecture information
	if osInfo.osArch == "" {
		t.Error("expected OS architecture to be populated")
	}

	if osInfo.procArch == "" {
		t.Error("expected processor architecture to be populated")
	}

	if osInfo.osArchBits == 0 {
		t.Error("expected OS architecture bits to be populated")
	}

	if osInfo.procArchBits == 0 {
		t.Error("expected processor architecture bits to be populated")
	}

	// Verify OS family is appropriate for current system
	if runtime.GOOS == "windows" {
		if !osInfo.families.Contains("windows") {
			t.Error("expected Windows family to be present on Windows system")
		}
	} else {
		if !osInfo.families.Contains("posix") {
			t.Error("expected POSIX family to be present on non-Windows system")
		}
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

	// Verify all expected keys are present
	expectedKeys := []string{
		"os_families", "os_id", "os_friendly_name", "os_release", "os_major_version",
		"os_version", "os_edition", "os_edition_id", "os_architecture", "os_architecture_bits",
		"processor_architecture", "processor_architecture_bits",
	}

	for _, key := range expectedKeys {
		if _, exists := values[key]; !exists {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	// Verify specific values
	if !values["os_families"].Type().Equals(cty.Set(cty.String)) {
		t.Errorf("expected os_families to be a set of strings, got %s", values["os_families"].Type().GoString())
	}

	if values["os_id"].AsString() != "ubuntu" {
		t.Errorf("expected os_id to be 'ubuntu', got %q", values["os_id"].AsString())
	}

	archBits, _ := values["os_architecture_bits"].AsBigFloat().Int64()
	if archBits != 64 {
		t.Errorf("expected os_architecture_bits to be 64, got %d", archBits)
	}
}

func TestOSInfo_ToMapOfCtyValues_EmptyValues(t *testing.T) {
	osInfo := newOSInfo()

	values := osInfo.toMapOfCtyValues()

	// All values should be null except os_families which should be null set
	nullKeys := []string{
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
	}

	for _, key := range nullKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected key %q to be null for empty osInfo, got %s", key, value.GoString())
			}
		}
	}

	// os_families should be null set
	if osFamilies, exists := values["os_families"]; exists {
		if !osFamilies.IsNull() || !osFamilies.Type().Equals(cty.Set(cty.String)) {
			t.Errorf("expected os_families to be null set of strings, got %s", osFamilies.GoString())
		}
	}
}

func TestArchitectureMap(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"i386", "386"},
		{"i486", "386"},
		{"i586", "386"},
		{"i686", "386"},
		{"x86", "386"},
		{"amd64", "amd64"},
		{"x86_64", "amd64"},
		{"armv6l", "arm"},
		{"armv7l", "arm"},
		{"aarch64", "arm64"},
		{"arm64", "arm64"},
		{"mips", "mips"},
		{"mips64", "mips64"},
		{"ppc64", "ppc64"},
		{"ppc64le", "ppc64le"},
		{"riscv64", "riscv64"},
		{"s390x", "s390x"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if mapped, exists := architectureMap[tc.input]; exists {
				if mapped != tc.expected {
					t.Errorf("expected architecture %q to map to %q, got %q", tc.input, tc.expected, mapped)
				}
			} else {
				t.Errorf("expected architecture %q to be in architectureMap", tc.input)
			}
		})
	}
}

func TestArchitectureBitsMap(t *testing.T) {
	testCases := []struct {
		arch string
		bits int
	}{
		{"386", 32},
		{"amd64", 64},
		{"arm", 32},
		{"arm64", 64},
		{"mips", 32},
		{"mips64", 64},
		{"ppc64", 64},
		{"ppc64le", 64},
		{"riscv64", 64},
		{"s390x", 64},
	}

	for _, tc := range testCases {
		t.Run(tc.arch, func(t *testing.T) {
			if bits, exists := architectureBitsMap[tc.arch]; exists {
				if bits != tc.bits {
					t.Errorf("expected architecture %q to have %d bits, got %d", tc.arch, tc.bits, bits)
				}
			} else {
				t.Errorf("expected architecture %q to be in architectureBitsMap", tc.arch)
			}
		})
	}
}

func TestOSFamiliesMap(t *testing.T) {
	testCases := []struct {
		osID     string
		families []string
	}{
		{"ubuntu", []string{"debian"}},
		{"centos", []string{"el"}},
		{"fedora", []string{"el"}},
		{"archlinux", []string{}}, // Not in the map, should have no additional families
		{"linuxmint", []string{"debian", "ubuntu"}},
		{"rhel", []string{"el"}},
		{"opensuse", []string{"suse"}},
	}

	for _, tc := range testCases {
		t.Run(tc.osID, func(t *testing.T) {
			if families, exists := osFamiliesMap[tc.osID]; exists {
				if len(families) != len(tc.families) {
					t.Errorf("expected %d families for %q, got %d", len(tc.families), tc.osID, len(families))
					return
				}
				for i, family := range tc.families {
					if families[i] != family {
						t.Errorf("expected family %q at index %d for %q, got %q", family, i, tc.osID, families[i])
					}
				}
			} else if len(tc.families) > 0 {
				t.Errorf("expected OS ID %q to be in osFamiliesMap", tc.osID)
			}
		})
	}
}

func TestOSIDCorrectionMap(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"amzn", "amazon"},
		{"arch", "archlinux"},
		{"archarm", "archlinux-arm"},
		{"clear-linux-os", "clearlinux"},
		{"cumulus-linux", "cumuluslinux"},
		{"pop", "pop_os"},
		{"ol", "oraclelinux"},
		{"opensuse-leap", "opensuse"},
		{"sles_sap", "sles"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if corrected, exists := osIDCorrectionMap[tc.input]; exists {
				if corrected != tc.expected {
					t.Errorf("expected OS ID %q to be corrected to %q, got %q", tc.input, tc.expected, corrected)
				}
			} else {
				t.Errorf("expected OS ID %q to be in osIDCorrectionMap", tc.input)
			}
		})
	}
}
