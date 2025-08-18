// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"os"
	"testing"

	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestOSInfo_PopulateOSInfo_Darwin(t *testing.T) {
	tests := []struct {
		name                 string
		output               string
		expectedFriendlyName string
		expectedRelease      string
		expectedMajorVersion string
		expectedVersion      string
		expectedArch         string
		expectedArchBits     int
	}{
		{
			name: "macOS 26.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "26.0.0"
				}`,
			expectedFriendlyName: "macOS 26.0.0",
			expectedRelease:      "Tahoe",
			expectedMajorVersion: "26",
			expectedVersion:      "26.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 26.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "26.0.0"
				}`,
			expectedFriendlyName: "macOS 26.0.0",
			expectedRelease:      "Tahoe",
			expectedMajorVersion: "26",
			expectedVersion:      "26.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 15.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "15.0.0"
				}`,
			expectedFriendlyName: "macOS 15.0.0",
			expectedRelease:      "Sequoia",
			expectedMajorVersion: "15",
			expectedVersion:      "15.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 15.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "15.0.0"
				}`,
			expectedFriendlyName: "macOS 15.0.0",
			expectedRelease:      "Sequoia",
			expectedMajorVersion: "15",
			expectedVersion:      "15.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 14.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "14.0.0"
				}`,
			expectedFriendlyName: "macOS 14.0.0",
			expectedRelease:      "Sonoma",
			expectedMajorVersion: "14",
			expectedVersion:      "14.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 14.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "14.0.0"
				}`,
			expectedFriendlyName: "macOS 14.0.0",
			expectedRelease:      "Sonoma",
			expectedMajorVersion: "14",
			expectedVersion:      "14.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 13.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "13.0.0"
				}`,
			expectedFriendlyName: "macOS 13.0.0",
			expectedRelease:      "Ventura",
			expectedMajorVersion: "13",
			expectedVersion:      "13.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 13.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "13.0.0"
				}`,
			expectedFriendlyName: "macOS 13.0.0",
			expectedRelease:      "Ventura",
			expectedMajorVersion: "13",
			expectedVersion:      "13.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 12.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "12.0.0"
				}`,
			expectedFriendlyName: "macOS 12.0.0",
			expectedRelease:      "Monterey",
			expectedMajorVersion: "12",
			expectedVersion:      "12.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 12.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "12.0.0"
				}`,
			expectedFriendlyName: "macOS 12.0.0",
			expectedRelease:      "Monterey",
			expectedMajorVersion: "12",
			expectedVersion:      "12.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 11.0.0 amd64",
			output: `{
				"os_arch": "amd64",
				"os_version": "11.0.0"
				}`,
			expectedFriendlyName: "macOS 11.0.0",
			expectedRelease:      "Big Sur",
			expectedMajorVersion: "11",
			expectedVersion:      "11.0.0",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "macOS 11.0.0 arm64",
			output: `{
				"os_arch": "arm64",
				"os_version": "11.0.0"
				}`,
			expectedFriendlyName: "macOS 11.0.0",
			expectedRelease:      "Big Sur",
			expectedMajorVersion: "11",
			expectedVersion:      "11.0.0",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
				Stdout: "Darwin",
			}
			mockTransport.CommandResults[osDarwinDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newOSInfo()

			diags := info.populateOSInfo(mockTransport)
			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if !info.Families().Contains("posix") {
				t.Error("expected POSIX family to be added")
			}

			if !info.Families().Contains("darwin") {
				t.Error("expected Darwin family to be added")
			}

			if !info.Families().Contains("macos") {
				t.Error("expected macOS family to be added")
			}

			if info.Families().Size() != 3 {
				t.Errorf("expected 3 families, got: %d", info.Families().Size())
			}

			expectedKernel := "darwin"
			if info.Kernel() != expectedKernel {
				t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
			}

			expectedID := "macos"
			if info.ID() != expectedID {
				t.Errorf("expected id to be %q, got: %q", expectedID, info.ID())
			}

			if info.FriendlyName() != tt.expectedFriendlyName {
				t.Errorf("expected friendlyName to be %q, got: %q", tt.expectedFriendlyName, info.FriendlyName())
			}

			if info.Release() != tt.expectedRelease {
				t.Errorf("expected release to be %q, got: %q", tt.expectedRelease, info.Release())
			}

			if info.MajorVersion() != tt.expectedMajorVersion {
				t.Errorf("expected majorVersion to be %q, got: %q", tt.expectedMajorVersion, info.MajorVersion())
			}

			if info.Version() != tt.expectedVersion {
				t.Errorf("expected version to be %q, got: %q", tt.expectedVersion, info.Version())
			}

			if info.Edition() != "" {
				t.Errorf("expected edition to be empty, got: %q", info.Edition())
			}

			if info.EditionID() != "" {
				t.Errorf("expected editionID to be empty, got: %q", info.EditionID())
			}

			if info.ProcArch() != tt.expectedArch {
				t.Errorf("expected procArch to be %q, got: %q", tt.expectedArch, info.ProcArch())
			}

			if info.OSArch() != tt.expectedArch {
				t.Errorf("expected osArch to be %q, got: %q", tt.expectedArch, info.OSArch())
			}

			if info.ProcArchBits() != tt.expectedArchBits {
				t.Errorf("expected procArchBits to be %d, got: %d", tt.expectedArchBits, info.ProcArchBits())
			}

			if info.OSArchBits() != tt.expectedArchBits {
				t.Errorf("expected osArchBits to be %d, got: %d", tt.expectedArchBits, info.OSArchBits())
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Darwin_Error(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Darwin",
	}
	mockTransport.CommandResults[osDarwinDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.Families().Contains("posix") {
		t.Error("expected POSIX family to be added")
	}

	if !info.Families().Contains("darwin") {
		t.Error("expected Darwin family to be added")
	}

	if !info.Families().Contains("macos") {
		t.Error("expected macOS family to be added")
	}

	if info.Families().Size() != 3 {
		t.Errorf("expected 3 families, got: %d", info.Families().Size())
	}

	expectedKernel := "darwin"
	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	expectedID := "macos"
	if info.ID() != expectedID {
		t.Errorf("expected id to be %q, got: %q", expectedID, info.ID())
	}

	expectedFriendlyName := "macOS"
	if info.FriendlyName() != expectedFriendlyName {
		t.Errorf("expected friendlyName to be %q, got: %q", expectedFriendlyName, info.FriendlyName())
	}

	if info.Release() != "" {
		t.Errorf("expected release to be empty, got: %q", info.Release())
	}

	if info.MajorVersion() != "" {
		t.Errorf("expected majorVersion to be empty, got: %q", info.MajorVersion())
	}

	if info.Version() != "" {
		t.Errorf("expected version to be empty, got: %q", info.Version())
	}

	if info.Edition() != "" {
		t.Errorf("expected edition to be empty, got: %q", info.Edition())
	}

	if info.EditionID() != "" {
		t.Errorf("expected editionID to be empty, got: %q", info.EditionID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected procArch to be empty, got: %q", info.ProcArch())
	}

	if info.OSArch() != "" {
		t.Errorf("expected osArch to be empty, got: %q", info.OSArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected procArchBits to be %d, got: %d", 0, info.ProcArchBits())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected osArchBits to be %d, got: %d", 0, info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to get macOS version"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing discovery command: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Darwin_NotJSON(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Darwin",
	}
	mockTransport.CommandResults[osDarwinDiscoveryScript] = &transport.MockCmd{
		Stdout: "Not a JSON output",
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if !info.Families().Contains("posix") {
		t.Error("expected POSIX family to be added")
	}

	if !info.Families().Contains("darwin") {
		t.Error("expected Darwin family to be added")
	}

	if !info.Families().Contains("macos") {
		t.Error("expected macOS family to be added")
	}

	if info.Families().Size() != 3 {
		t.Errorf("expected 3 families, got: %d", info.Families().Size())
	}

	expectedKernel := "darwin"
	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	expectedID := "macos"
	if info.ID() != expectedID {
		t.Errorf("expected id to be %q, got: %q", expectedID, info.ID())
	}

	expectedFriendlyName := "macOS"
	if info.FriendlyName() != expectedFriendlyName {
		t.Errorf("expected friendlyName to be %q, got: %q", expectedFriendlyName, info.FriendlyName())
	}

	if info.Release() != "" {
		t.Errorf("expected release to be empty, got: %q", info.Release())
	}

	if info.MajorVersion() != "" {
		t.Errorf("expected majorVersion to be empty, got: %q", info.MajorVersion())
	}

	if info.Version() != "" {
		t.Errorf("expected version to be empty, got: %q", info.Version())
	}

	if info.Edition() != "" {
		t.Errorf("expected edition to be empty, got: %q", info.Edition())
	}

	if info.EditionID() != "" {
		t.Errorf("expected editionID to be empty, got: %q", info.EditionID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected procArch to be empty, got: %q", info.ProcArch())
	}

	if info.OSArch() != "" {
		t.Errorf("expected osArch to be empty, got: %q", info.OSArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected procArchBits to be %d, got: %d", 0, info.ProcArchBits())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected osArchBits to be %d, got: %d", 0, info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse macOS discovery output"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected summary %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing JSON output: invalid character 'N' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected detail %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Darwin_UnknownArchitecture(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Darwin",
	}
	mockTransport.CommandResults[osDarwinDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			"os_arch": "newarch",
			"os_version": "26.0.0"
			}
			`,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatalf("expected warnings, got none")
	}

	if !info.Families().Contains("posix") {
		t.Error("expected POSIX family to be added")
	}

	if !info.Families().Contains("darwin") {
		t.Error("expected Darwin family to be added")
	}

	if !info.Families().Contains("macos") {
		t.Error("expected macOS family to be added")
	}

	expectedKernel := "darwin"
	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	expectedArch := "newarch"
	if info.ProcArch() != expectedArch {
		t.Errorf("expected procArch to be %q, got: %q", expectedArch, info.ProcArch())
	}

	if info.OSArch() != expectedArch {
		t.Errorf("expected osArch to be %q, got: %q", expectedArch, info.OSArch())
	}

	expectedArchBits := 0
	if info.ProcArchBits() != expectedArchBits {
		t.Errorf("expected procArchBits to be %d, got: %d", expectedArchBits, info.ProcArchBits())
	}

	if info.OSArchBits() != expectedArchBits {
		t.Errorf("expected osArchBits to be %d, got: %d", expectedArchBits, info.OSArchBits())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Unknown architecture"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected first warning summary to be %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Unknown architecture \"newarch\" detected, using it as is"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected first warning detail to be %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Darwin_UnknownVersion(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Darwin",
	}
	mockTransport.CommandResults[osDarwinDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			"os_arch": "amd64",
			"os_version": "99.0.0"
			}`,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags.Errors())
	}

	if !diags.HasWarnings() {
		t.Fatalf("expected warnings, got none")
	}

	if !info.Families().Contains("posix") {
		t.Error("expected POSIX family to be added")
	}

	if !info.Families().Contains("darwin") {
		t.Error("expected Darwin family to be added")
	}

	if !info.Families().Contains("macos") {
		t.Error("expected macOS family to be added")
	}

	if info.Families().Size() != 3 {
		t.Errorf("expected 3 families, got: %d", info.Families().Size())
	}

	if info.ID() != "macos" {
		t.Errorf("expected OS ID to be \"macos\", got: %q", info.ID())
	}

	expectedFriendlyName := "macOS 99.0.0"
	if info.FriendlyName() != expectedFriendlyName {
		t.Errorf("expected friendly name to be %q, got: %q", expectedFriendlyName, info.FriendlyName())
	}

	if info.Release() != "" {
		t.Errorf("expected release to be empty, got: %q", info.Release())
	}

	expectedMajorVersion := "99"
	if info.MajorVersion() != expectedMajorVersion {
		t.Errorf("expected major version to be %q, got: %q", expectedMajorVersion, info.MajorVersion())
	}

	if info.Version() != "99.0.0" {
		t.Errorf("expected version to be %q, got: %q", "99.0.0", info.Version())
	}

	expectedArch := "amd64"
	if info.OSArch() != expectedArch {
		t.Errorf("expected OS architecture to be %q, got: %q", expectedArch, info.OSArch())
	}

	if info.OSArchBits() != 64 {
		t.Errorf("expected OS architecture bits to be 64, got: %d", info.OSArchBits())
	}

	if info.ProcArch() != expectedArch {
		t.Errorf("expected processor architecture to be %q, got: %q", expectedArch, info.ProcArch())
	}

	if info.ProcArchBits() != 64 {
		t.Errorf("expected processor architecture bits to be 64, got: %d", info.ProcArchBits())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Unknown macOS release"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected warning summary to be %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Unknown macOS release detected for major version 99"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected warning detail to be %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Linux(t *testing.T) {
	tests := []struct {
		name                 string
		output               string
		expectedFamilies     []string
		expectedId           string
		expectedFriendlyName string
		expectedMajorVersion string
		expectedVersion      string
		expectedRelease      string
		expectedEdition      string
		expectedEditionId    string
		expectedArch         string
		expectedArchBits     int
	}{
		{
			name: "almalinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "almalinux",
				"os_friendly_name": "AlmaLinux 8.5",
				"os_release": "almalinux",
				"os_version": "8.5",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "almalinux"},
			expectedId:           "almalinux",
			expectedFriendlyName: "AlmaLinux 8.5",
			expectedMajorVersion: "8",
			expectedVersion:      "8.5",
			expectedRelease:      "almalinux",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "amazon",
			output: `{
			    "os_arch": "amd64",
				"os_id": "amzn",
				"os_friendly_name": "Amazon Linux 2",
				"os_release": "amzn",
				"os_version": "2",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "amazon"},
			expectedId:           "amazon",
			expectedFriendlyName: "Amazon Linux 2",
			expectedMajorVersion: "2",
			expectedVersion:      "2",
			expectedRelease:      "amzn",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "archlinux-arm",
			output: `{
			    "os_arch": "amd64",
				"os_id": "archlinux-arm",
				"os_friendly_name": "Arch Linux ARM",
				"os_release": "rolling",
				"os_version": "rolling",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "archlinux", "archlinux-arm"},
			expectedId:           "archlinux-arm",
			expectedFriendlyName: "Arch Linux ARM",
			expectedMajorVersion: "rolling",
			expectedVersion:      "rolling",
			expectedRelease:      "rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "arcolinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "arcolinux",
				"os_friendly_name": "ArcoLinux",
				"os_release": "rolling",
				"os_version": "rolling",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "archlinux", "arcolinux"},
			expectedId:           "arcolinux",
			expectedFriendlyName: "ArcoLinux",
			expectedMajorVersion: "rolling",
			expectedVersion:      "rolling",
			expectedRelease:      "rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "centos",
			output: `{
			    "os_arch": "amd64",
				"os_id": "centos",
				"os_friendly_name": "CentOS Linux 7",
				"os_release": "centos",
				"os_version": "7",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "centos"},
			expectedId:           "centos",
			expectedFriendlyName: "CentOS Linux 7",
			expectedMajorVersion: "7",
			expectedVersion:      "7",
			expectedRelease:      "centos",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "clearos",
			output: `{
			    "os_arch": "amd64",
				"os_id": "clearos",
				"os_friendly_name": "ClearOS",
				"os_release": "clearos",
				"os_version": "7",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "clearos"},
			expectedId:           "clearos",
			expectedFriendlyName: "ClearOS",
			expectedMajorVersion: "7",
			expectedVersion:      "7",
			expectedRelease:      "clearos",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "clearlinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "clearlinux",
				"os_friendly_name": "Clear Linux OS",
				"os_release": "clearlinux",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "clearlinux"},
			expectedId:           "clearlinux",
			expectedFriendlyName: "Clear Linux OS",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "clearlinux",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "cloudlinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "cloudlinux",
				"os_friendly_name": "CloudLinux 7",
				"os_release": "cloudlinux",
				"os_version": "7",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "cloudlinux"},
			expectedId:           "cloudlinux",
			expectedFriendlyName: "CloudLinux 7",
			expectedMajorVersion: "7",
			expectedVersion:      "7",
			expectedRelease:      "cloudlinux",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "cumuluslinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "cumuluslinux",
				"os_friendly_name": "Cumulus Linux 3.7",
				"os_release": "cumulus",
				"os_version": "3.7",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "cumuluslinux"},
			expectedId:           "cumuluslinux",
			expectedFriendlyName: "Cumulus Linux 3.7",
			expectedMajorVersion: "3",
			expectedVersion:      "3.7",
			expectedRelease:      "cumulus",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "deepin",
			output: `{
			    "os_arch": "amd64",
				"os_id": "deepin",
				"os_friendly_name": "Deepin 20.2",
				"os_release": "",
				"os_version": "20.2",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "deepin"},
			expectedId:           "deepin",
			expectedFriendlyName: "Deepin 20.2",
			expectedMajorVersion: "20",
			expectedVersion:      "20.2",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "devuan",
			output: `{
			    "os_arch": "amd64",
				"os_id": "devuan",
				"os_friendly_name": "Devuan GNU/Linux 2.1 (ASCII)",
				"os_release": "ascii",
				"os_version": "2.1",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "devuan"},
			expectedId:           "devuan",
			expectedFriendlyName: "Devuan GNU/Linux 2.1 (ASCII)",
			expectedMajorVersion: "2",
			expectedVersion:      "2.1",
			expectedRelease:      "ascii",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "elementary",
			output: `{
			    "os_arch": "amd64",
				"os_id": "elementary",
				"os_friendly_name": "elementary OS 6.1",
				"os_release": "juno",
				"os_version": "6.1",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "ubuntu", "elementary"},
			expectedId:           "elementary",
			expectedFriendlyName: "elementary OS 6.1",
			expectedMajorVersion: "6",
			expectedVersion:      "6.1",
			expectedRelease:      "juno",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "endeavouros",
			output: `{
			    "os_arch": "amd64",
				"os_id": "endeavouros",
				"os_friendly_name": "EndeavourOS",
				"os_release": "rolling",
				"os_version": "rolling",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "archlinux", "endeavouros"},
			expectedId:           "endeavouros",
			expectedFriendlyName: "EndeavourOS",
			expectedMajorVersion: "rolling",
			expectedVersion:      "rolling",
			expectedRelease:      "rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "fedora",
			output: `{
			    "os_arch": "amd64",
				"os_id": "fedora",
				"os_friendly_name": "Fedora 34 (Workstation Edition)",
				"os_release": "",
				"os_version": "34",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "fedora"},
			expectedId:           "fedora",
			expectedFriendlyName: "Fedora 34 (Workstation Edition)",
			expectedMajorVersion: "34",
			expectedVersion:      "34",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "kali",
			output: `{
			    "os_arch": "amd64",
				"os_id": "kali",
				"os_friendly_name": "Kali GNU/Linux 2021.3",
				"os_release": "kali-rolling",
				"os_version": "2021.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "kali"},
			expectedId:           "kali",
			expectedFriendlyName: "Kali GNU/Linux 2021.3",
			expectedMajorVersion: "2021",
			expectedVersion:      "2021.3",
			expectedRelease:      "kali-rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "kylin",
			output: `{
			    "os_arch": "amd64",
				"os_id": "kylin",
				"os_friendly_name": "Kylin 10",
				"os_release": "kylin",
				"os_version": "10",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "ubuntu", "kylin"},
			expectedId:           "kylin",
			expectedFriendlyName: "Kylin 10",
			expectedMajorVersion: "10",
			expectedVersion:      "10",
			expectedRelease:      "kylin",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "linuxmint",
			output: `{
			    "os_arch": "amd64",
				"os_id": "linuxmint",
				"os_friendly_name": "Linux Mint 20.2",
				"os_release": "uma",
				"os_version": "20.2",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "ubuntu", "linuxmint"},
			expectedId:           "linuxmint",
			expectedFriendlyName: "Linux Mint 20.2",
			expectedMajorVersion: "20",
			expectedVersion:      "20.2",
			expectedRelease:      "uma",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "mageia",
			output: `{
			    "os_arch": "amd64",
				"os_id": "mageia",
				"os_friendly_name": "Mageia 8",
				"os_release": "",
				"os_version": "8",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "mandrake", "mageia"},
			expectedId:           "mageia",
			expectedFriendlyName: "Mageia 8",
			expectedMajorVersion: "8",
			expectedVersion:      "8",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "manjaro",
			output: `{
			    "os_arch": "amd64",
				"os_id": "manjaro",
				"os_friendly_name": "Manjaro Linux",
				"os_release": "rolling",
				"os_version": "rolling",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "archlinux", "manjaro"},
			expectedId:           "manjaro",
			expectedFriendlyName: "Manjaro Linux",
			expectedMajorVersion: "rolling",
			expectedVersion:      "rolling",
			expectedRelease:      "rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "manjaro-arm",
			output: `{
			    "os_arch": "amd64",
				"os_id": "manjaro-arm",
				"os_friendly_name": "Manjaro ARM",
				"os_release": "rolling",
				"os_version": "rolling",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "archlinux", "manjaro", "manjaro-arm"},
			expectedId:           "manjaro-arm",
			expectedFriendlyName: "Manjaro ARM",
			expectedMajorVersion: "rolling",
			expectedVersion:      "rolling",
			expectedRelease:      "rolling",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "nobara",
			output: `{
			    "os_arch": "amd64",
				"os_id": "nobara",
				"os_friendly_name": "Nobara 38",
				"os_release": "nobara",
				"os_version": "38",
				"os_edition": "KDE Plasma Desktop Edition",
				"os_edition_id": "kde"
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "fedora", "nobara"},
			expectedId:           "nobara",
			expectedFriendlyName: "Nobara 38",
			expectedMajorVersion: "38",
			expectedVersion:      "38",
			expectedRelease:      "nobara",
			expectedEdition:      "KDE Plasma Desktop Edition",
			expectedEditionId:    "kde",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "opensuse",
			output: `{
			    "os_arch": "amd64",
				"os_id": "opensuse",
				"os_friendly_name": "openSUSE Leap 15.3",
				"os_release": "",
				"os_version": "15.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "suse", "opensuse"},
			expectedId:           "opensuse",
			expectedFriendlyName: "openSUSE Leap 15.3",
			expectedMajorVersion: "15",
			expectedVersion:      "15.3",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "oraclelinux",
			output: `{
			    "os_arch": "amd64",
				"os_id": "oraclelinux",
				"os_friendly_name": "Oracle Linux Server 8.5",
				"os_release": "ol8",
				"os_version": "8.5",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "oraclelinux"},
			expectedId:           "oraclelinux",
			expectedFriendlyName: "Oracle Linux Server 8.5",
			expectedMajorVersion: "8",
			expectedVersion:      "8.5",
			expectedRelease:      "ol8",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "pop_os",
			output: `{
			    "os_arch": "amd64",
				"os_id": "pop_os",
				"os_friendly_name": "Pop!_OS 21.04",
				"os_release": "hirsute",
				"os_version": "21.04",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "ubuntu", "pop_os"},
			expectedId:           "pop_os",
			expectedFriendlyName: "Pop!_OS 21.04",
			expectedMajorVersion: "21",
			expectedVersion:      "21.04",
			expectedRelease:      "hirsute",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "raspbian",
			output: `{
			    "os_arch": "amd64",
				"os_id": "raspbian",
				"os_friendly_name": "Raspbian GNU/Linux 10 (buster)",
				"os_release": "buster",
				"os_version": "10",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "raspbian"},
			expectedId:           "raspbian",
			expectedFriendlyName: "Raspbian GNU/Linux 10 (buster)",
			expectedMajorVersion: "10",
			expectedVersion:      "10",
			expectedRelease:      "buster",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "rhel",
			output: `{
			    "os_arch": "amd64",
				"os_id": "rhel",
				"os_friendly_name": "Red Hat Enterprise Linux 8.5",
				"os_release": "ol8",
				"os_version": "8.5",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "rhel"},
			expectedId:           "rhel",
			expectedFriendlyName: "Red Hat Enterprise Linux 8.5",
			expectedMajorVersion: "8",
			expectedVersion:      "8.5",
			expectedRelease:      "ol8",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "rocky",
			output: `{
			    "os_arch": "amd64",
				"os_id": "rocky",
				"os_friendly_name": "Rocky Linux 8.5",
				"os_release": "rocky",
				"os_version": "8.5",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "rocky"},
			expectedId:           "rocky",
			expectedFriendlyName: "Rocky Linux 8.5",
			expectedMajorVersion: "8",
			expectedVersion:      "8.5",
			expectedRelease:      "rocky",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "scientific",
			output: `{
			    "os_arch": "amd64",
				"os_id": "scientific",
				"os_friendly_name": "Scientific Linux 8",
				"os_release": "scientific",
				"os_version": "8",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "scientific"},
			expectedId:           "scientific",
			expectedFriendlyName: "Scientific Linux 8",
			expectedMajorVersion: "8",
			expectedVersion:      "8",
			expectedRelease:      "scientific",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "sled",
			output: `{
			    "os_arch": "amd64",
				"os_id": "sled",
				"os_friendly_name": "SUSE Linux Enterprise Desktop 15 SP3",
				"os_release": "",
				"os_version": "15.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "suse", "sled"},
			expectedId:           "sled",
			expectedFriendlyName: "SUSE Linux Enterprise Desktop 15 SP3",
			expectedMajorVersion: "15",
			expectedVersion:      "15.3",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "sles",
			output: `{
			    "os_arch": "amd64",
				"os_id": "sles",
				"os_friendly_name": "SUSE Linux Enterprise Server 15 SP3",
				"os_release": "",
				"os_version": "15.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "suse", "sles"},
			expectedId:           "sles",
			expectedFriendlyName: "SUSE Linux Enterprise Server 15 SP3",
			expectedMajorVersion: "15",
			expectedVersion:      "15.3",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "sles_sap",
			output: `{
			    "os_arch": "amd64",
				"os_id": "sles",
				"os_friendly_name": "SUSE Linux Enterprise Server 15 SP3 for SAP Applications",
				"os_release": "",
				"os_version": "15.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "suse", "sles"},
			expectedId:           "sles",
			expectedFriendlyName: "SUSE Linux Enterprise Server 15 SP3 for SAP Applications",
			expectedMajorVersion: "15",
			expectedVersion:      "15.3",
			expectedRelease:      "",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "ubuntu",
			output: `{
			    "os_arch": "amd64",
				"os_id": "ubuntu",
				"os_friendly_name": "Ubuntu 20.04.3 LTS",
				"os_release": "focal",
				"os_version": "20.04.3",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "debian", "ubuntu"},
			expectedId:           "ubuntu",
			expectedFriendlyName: "Ubuntu 20.04.3 LTS",
			expectedMajorVersion: "20",
			expectedVersion:      "20.04.3",
			expectedRelease:      "focal",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "virtuozzo",
			output: `{
			    "os_arch": "amd64",
				"os_id": "virtuozzo",
				"os_friendly_name": "Virtuozzo Linux 7",
				"os_release": "virtuozzo",
				"os_version": "7",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "el", "virtuozzo"},
			expectedId:           "virtuozzo",
			expectedFriendlyName: "Virtuozzo Linux 7",
			expectedMajorVersion: "7",
			expectedVersion:      "7",
			expectedRelease:      "virtuozzo",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "generic 386",
			output: `{
			    "os_arch": "386",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic i386",
			output: `{
			    "os_arch": "i386",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic i486",
			output: `{
			    "os_arch": "i486",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic i586",
			output: `{
			    "os_arch": "i586",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic i686",
			output: `{
			    "os_arch": "i686",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic x86",
			output: `{
			    "os_arch": "x86",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "386",
			expectedArchBits:     32,
		},
		{
			name: "generic x86_64",
			output: `{
			    "os_arch": "x86_64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "amd64",
			expectedArchBits:     64,
		},
		{
			name: "generic arm",
			output: `{
			    "os_arch": "arm",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "arm",
			expectedArchBits:     32,
		},
		{
			name: "generic armv6l",
			output: `{
			    "os_arch": "armv6l",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "arm",
			expectedArchBits:     32,
		},
		{
			name: "generic armv7l",
			output: `{
			    "os_arch": "armv7l",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "arm",
			expectedArchBits:     32,
		},
		{
			name: "generic aarch64",
			output: `{
			    "os_arch": "aarch64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "generic arm64",
			output: `{
			    "os_arch": "arm64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "arm64",
			expectedArchBits:     64,
		},
		{
			name: "generic mips",
			output: `{
			    "os_arch": "mips",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "mips",
			expectedArchBits:     32,
		},
		{
			name: "generic mips64",
			output: `{
			    "os_arch": "mips64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "mips64",
			expectedArchBits:     64,
		},
		{
			name: "generic ppc64",
			output: `{
			    "os_arch": "ppc64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "ppc64",
			expectedArchBits:     64,
		},
		{
			name: "generic ppc64le",
			output: `{
			    "os_arch": "ppc64le",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "ppc64le",
			expectedArchBits:     64,
		},
		{
			name: "generic riscv64",
			output: `{
			    "os_arch": "riscv64",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "riscv64",
			expectedArchBits:     64,
		},
		{
			name: "generic s390x",
			output: `{
			    "os_arch": "s390x",
				"os_id": "generic",
				"os_friendly_name": "Generic Linux 1.0",
				"os_release": "generic",
				"os_version": "1.0",
				"os_edition": "",
				"os_edition_id": ""
				}`,
			expectedFamilies:     []string{"posix", "linux", "generic"},
			expectedId:           "generic",
			expectedFriendlyName: "Generic Linux 1.0",
			expectedMajorVersion: "1",
			expectedVersion:      "1.0",
			expectedRelease:      "generic",
			expectedEdition:      "",
			expectedEditionId:    "",
			expectedArch:         "s390x",
			expectedArchBits:     64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := transport.NewMockTransport()
			mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
				Stdout: "Linux",
			}
			mockTransport.CommandResults[osLinuxDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newOSInfo()

			diags := info.populateOSInfo(mockTransport)
			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			for _, family := range tt.expectedFamilies {
				if !info.families.Contains(family) {
					t.Errorf("expected family %q to be added, but it was not", family)
				}
			}

			if len(tt.expectedFamilies) != info.families.Size() {
				t.Errorf("expected %d families, got: %d", len(tt.expectedFamilies), info.families.Size())
			}

			if info.ID() != tt.expectedId {
				t.Errorf("expected id %q, got: %q", tt.expectedId, info.ID())
			}

			if info.FriendlyName() != tt.expectedFriendlyName {
				t.Errorf("expected friendly name %q, got: %q", tt.expectedFriendlyName, info.FriendlyName())
			}

			if info.MajorVersion() != tt.expectedMajorVersion {
				t.Errorf("expected major version %q, got: %q", tt.expectedMajorVersion, info.MajorVersion())
			}

			if info.Version() != tt.expectedVersion {
				t.Errorf("expected version %q, got: %q", tt.expectedVersion, info.Version())
			}

			if info.Release() != tt.expectedRelease {
				t.Errorf("expected release %q, got: %q", tt.expectedRelease, info.Release())
			}

			if info.Edition() != tt.expectedEdition {
				t.Errorf("expected edition %q, got: %q", tt.expectedEdition, info.Edition())
			}

			if info.EditionID() != tt.expectedEditionId {
				t.Errorf("expected edition id %q, got: %q", tt.expectedEditionId, info.EditionID())
			}

			if info.ProcArch() != tt.expectedArch {
				t.Errorf("expected architecture %q, got: %q", tt.expectedArch, info.ProcArch())
			}

			if info.ProcArchBits() != tt.expectedArchBits {
				t.Errorf("expected architecture bits %d, got: %d", tt.expectedArchBits, info.ProcArchBits())
			}

			if info.OSArch() != tt.expectedArch {
				t.Errorf("expected architecture %q, got: %q", tt.expectedArch, info.OSArch())
			}

			if info.OSArchBits() != tt.expectedArchBits {
				t.Errorf("expected architecture bits %d, got: %d", tt.expectedArchBits, info.OSArchBits())
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Linux_Error(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Linux",
	}
	mockTransport.CommandResults[osLinuxDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.families.Size() != 2 {
		t.Errorf("expected 2 families (posix, linux), got: %d", info.families.Size())
	}

	if !info.families.Contains("posix") {
		t.Errorf("expected family 'posix' to be added, but it was not")
	}

	if !info.families.Contains("linux") {
		t.Errorf("expected family 'linux' to be added, but it was not")
	}

	if info.ID() != "" {
		t.Errorf("expected id to be empty, got: %q", info.ID())
	}

	if info.FriendlyName() != "" {
		t.Errorf("expected friendly name to be empty, got: %q", info.FriendlyName())
	}

	if info.MajorVersion() != "" {
		t.Errorf("expected major version to be empty, got: %q", info.MajorVersion())
	}

	if info.Version() != "" {
		t.Errorf("expected version to be empty, got: %q", info.Version())
	}

	if info.Release() != "" {
		t.Errorf("expected release to be empty, got: %q", info.Release())
	}

	if info.Edition() != "" {
		t.Errorf("expected edition to be empty, got: %q", info.Edition())
	}

	if info.EditionID() != "" {
		t.Errorf("expected edition id to be empty, got: %q", info.EditionID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected architecture to be empty, got: %q", info.ProcArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected architecture bits to be 0, got: %d", info.ProcArchBits())
	}

	if info.OSArch() != "" {
		t.Errorf("expected architecture to be empty, got: %q", info.OSArch())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected architecture bits to be 0, got: %d", info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to get Linux OS information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary to be %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing Linux discovery script: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail to be %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Linux_NotJSON(t *testing.T) {
	mockTransport := transport.NewMockTransport()
	mockTransport.CommandResults["uname -s"] = &transport.MockCmd{
		Stdout: "Linux",
	}
	mockTransport.CommandResults[osLinuxDiscoveryScript] = &transport.MockCmd{
		Stdout: "Not JSON",
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.families.Size() != 2 {
		t.Errorf("expected 2 families (posix, linux), got: %d", info.families.Size())
	}

	if !info.families.Contains("posix") {
		t.Errorf("expected family 'posix' to be added, but it was not")
	}

	if !info.families.Contains("linux") {
		t.Errorf("expected family 'linux' to be added, but it was not")
	}

	if info.ID() != "" {
		t.Errorf("expected id to be empty, got: %q", info.ID())
	}

	if info.FriendlyName() != "" {
		t.Errorf("expected friendly name to be empty, got: %q", info.FriendlyName())
	}

	if info.MajorVersion() != "" {
		t.Errorf("expected major version to be empty, got: %q", info.MajorVersion())
	}

	if info.Version() != "" {
		t.Errorf("expected version to be empty, got: %q", info.Version())
	}

	if info.Release() != "" {
		t.Errorf("expected release to be empty, got: %q", info.Release())
	}

	if info.Edition() != "" {
		t.Errorf("expected edition to be empty, got: %q", info.Edition())
	}

	if info.EditionID() != "" {
		t.Errorf("expected edition id to be empty, got: %q", info.EditionID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected architecture to be empty, got: %q", info.ProcArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected architecture bits to be 0, got: %d", info.ProcArchBits())
	}

	if info.OSArch() != "" {
		t.Errorf("expected architecture to be empty, got: %q", info.OSArch())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected architecture bits to be 0, got: %d", info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse Linux discovery output"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary to be %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing JSON output: invalid character 'N' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail to be %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Windows_Architecture(t *testing.T) {

	tests := []struct {
		name                 string
		output               string
		expectedprocArch     string
		expectedprocArchBits int
		expectedosArch       string
		expectedosArchBits   int
	}{
		{
			name: "x86",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "32-bit",
				"processor_arch": "X86"
				}
				`,
			expectedprocArch:     "386",
			expectedprocArchBits: 32,
			expectedosArch:       "386",
			expectedosArchBits:   32,
		},
		{
			name: "AMD64",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
			expectedprocArch:     "amd64",
			expectedprocArchBits: 64,
			expectedosArch:       "amd64",
			expectedosArchBits:   64,
		},
		{
			name: "32-bit OS on AMD64 processor",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "32-bit",
				"processor_arch": "AMD64"
				}
				`,
			expectedprocArch:     "amd64",
			expectedprocArchBits: 64,
			expectedosArch:       "386",
			expectedosArchBits:   32,
		},
		{
			name: "ARM64",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "64-bit",
				"processor_arch": "ARM64"
				}
				`,
			expectedprocArch:     "arm64",
			expectedprocArchBits: 64,
			expectedosArch:       "arm64",
			expectedosArchBits:   64,
		},
		{
			name: "ARM",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "32-bit",
				"processor_arch": "ARM"
				}
				`,
			expectedprocArch:     "arm",
			expectedprocArchBits: 32,
			expectedosArch:       "arm",
			expectedosArchBits:   32,
		},
		{
			name: "32-bit OS on ARM64 processor",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19041.0",
				"os_bits": "32-bit",
				"processor_arch": "ARM64"
				}
				`,
			expectedprocArch:     "arm64",
			expectedprocArchBits: 64,
			expectedosArch:       "arm",
			expectedosArchBits:   32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := transport.NewWinMockTransport()
			mockTransport.PowerShellResults["Write-Host $PSVersionTable.PSVersion"] = &transport.MockCmd{
				Stdout: "5.1.19041.1237",
			}
			mockTransport.PowerShellResults[osWindowsDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newOSInfo()

			diags := info.populateOSInfo(mockTransport)
			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			if info.Families().Size() != 2 {
				t.Errorf("expected 2 families, got: %d", info.Families().Size())
			}

			expectedKernel := "windows"
			if !info.Families().Contains(expectedKernel) {
				t.Errorf("expected family %q to be present, but it was not", expectedKernel)
			}

			expectedID := "windows-client"
			if !info.Families().Contains(expectedID) {
				t.Errorf("expected family %q to be present, but it was not", expectedID)
			}

			if info.Kernel() != expectedKernel {
				t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
			}

			if info.ID() != expectedID {
				t.Errorf("expected ID to be %q, got: %q", expectedID, info.ID())
			}

			if info.ProcArch() != tt.expectedprocArch {
				t.Errorf("expected processor architecture to be %q, got: %q", tt.expectedprocArch, info.ProcArch())
			}

			if info.ProcArchBits() != tt.expectedprocArchBits {
				t.Errorf("expected processor architecture bits to be %d, got: %d", tt.expectedprocArchBits, info.ProcArchBits())
			}

			if info.OSArch() != tt.expectedosArch {
				t.Errorf("expected OS architecture to be %q, got: %q", tt.expectedosArch, info.OSArch())
			}

			if info.OSArchBits() != tt.expectedosArchBits {
				t.Errorf("expected OS architecture bits to be %d, got: %d", tt.expectedosArchBits, info.OSArchBits())
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Windows_Architecture_Unknown(t *testing.T) {
	mockTransport := transport.NewWinMockTransport()
	mockTransport.PowerShellResults["Write-Host $PSVersionTable.PSVersion"] = &transport.MockCmd{
		Stdout: "5.1.19041.1237",
	}
	mockTransport.PowerShellResults[osWindowsDiscoveryScript] = &transport.MockCmd{
		Stdout: `{
			"os_friendly_name": "Microsoft Windows 10 Pro",
			"os_version": "10.0.19041.0",
			"os_bits": "64-bit",
			"processor_arch": "newarch"
			}
			`,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if diags.HasErrors() {
		t.Fatalf("expected no errors, got: %v", diags)
	}

	if !diags.HasWarnings() {
		t.Fatalf("expected warnings, got none")
	}

	if info.Families().Size() != 2 {
		t.Errorf("expected 2 families, got: %d", info.Families().Size())
	}

	expectedKernel := "windows"
	if !info.Families().Contains(expectedKernel) {
		t.Errorf("expected family %q to be present, but it was not", expectedKernel)
	}

	expectedID := "windows-client"
	if !info.Families().Contains(expectedID) {
		t.Errorf("expected family %q to be present, but it was not", expectedID)
	}

	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	if info.ID() != expectedID {
		t.Errorf("expected ID to be %q, got: %q", expectedID, info.ID())
	}

	expectedArch := "newarch"
	if info.ProcArch() != expectedArch {
		t.Errorf("expected processor architecture to be %q, got: %s", expectedArch, info.ProcArch())
	}

	expectedArchBits := 0
	if info.ProcArchBits() != expectedArchBits {
		t.Errorf("expected processor architecture bits to be %d, got: %d", expectedArchBits, info.ProcArchBits())
	}

	if info.OSArch() != expectedArch {
		t.Errorf("expected OS architecture to be %q, got: %s", expectedArch, info.OSArch())
	}

	if info.OSArchBits() != expectedArchBits {
		t.Errorf("expected OS architecture bits to be %d, got: %d", expectedArchBits, info.OSArchBits())
	}

	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got: %d", len(warnings))
	}

	expectedSummary := "Unknown architecture"
	if warnings[0].Summary != expectedSummary {
		t.Errorf("expected first warning summary to be %q, got: %q", expectedSummary, warnings[0].Summary)
	}

	expectedDetail := "Unknown architecture \"newarch\" detected, using it as is"
	if warnings[0].Detail != expectedDetail {
		t.Errorf("expected first warning detail to be %q, got: %q", expectedDetail, warnings[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Windows(t *testing.T) {
	tests := []struct {
		name                 string
		output               string
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
			name: "Windows Server 2008 R2 Standard",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2008 R2 Standard",
				"os_version": "6.1.7600.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 7 Professional",
			output: `{
				"os_friendly_name": "Microsoft Windows 7 Professional",
				"os_version": "6.1.7600.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2008 R2 SP1 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2008 R2 SP1 Enterprise",
				"os_version": "6.1.7601.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 7 SP1 Ultimate",
			output: `{
				"os_friendly_name": "Microsoft Windows 7 Ultimate",
				"os_version": "6.1.7601.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2012 Datacenter",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2012 Datacenter",
				"os_version": "6.2.9200.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 8 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 8 Pro",
				"os_version": "6.2.9200.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2012 R2 Standard",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2012 R2 Standard",
				"os_version": "6.3.9600.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 8.1 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows 8.1 Enterprise",
				"os_version": "6.3.9600.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1507 Home",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Home",
				"os_version": "10.0.10240.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1511 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.10586.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2016 Datacenter",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2016 Datacenter",
				"os_version": "10.0.14393.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1607 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Enterprise",
				"os_version": "10.0.14393.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1703 Education",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Education",
				"os_version": "10.0.15063.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1709 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.16299.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1803 Home",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Home",
				"os_version": "10.0.17134.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2019 Standard",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2019 Standard",
				"os_version": "10.0.17763.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1809 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Enterprise",
				"os_version": "10.0.17763.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1903 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.18362.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 1909 Core",
			output: `{
				"os_friendly_name": "Microsoft Windows Server Core",
				"os_version": "10.0.18363.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 1909 Home",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Home",
				"os_version": "10.0.18363.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2004 Datacenter",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2004 Datacenter",
				"os_version": "10.0.19041.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 2004 Education",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Education",
				"os_version": "10.0.19041.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 20H2 Standard",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 20H2 Standard",
				"os_version": "10.0.19042.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 20H2 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19042.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 21H1 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Enterprise",
				"os_version": "10.0.19043.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 21H2 Home",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Home",
				"os_version": "10.0.19044.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 10 22H2 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 10 Pro",
				"os_version": "10.0.19045.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2022 Datacenter",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2022 Datacenter",
				"os_version": "10.0.20348.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 11 21H2 Home",
			output: `{
				"os_friendly_name": "Microsoft Windows 11 Home",
				"os_version": "10.0.22000.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 11 22H2 Pro",
			output: `{
				"os_friendly_name": "Microsoft Windows 11 Pro",
				"os_version": "10.0.22621.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 11 23H2 Enterprise",
			output: `{
				"os_friendly_name": "Microsoft Windows 11 Enterprise",
				"os_version": "10.0.22631.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 23H2 Standard",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 23H2 Standard",
				"os_version": "10.0.25398.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows Server 2025 Datacenter Evaluation",
			output: `{
				"os_friendly_name": "Microsoft Windows Server 2025 Datacenter Evaluation",
				"os_version": "10.0.26100.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			name: "Windows 11 24H2 Education",
			output: `{
				"os_friendly_name": "Microsoft Windows 11 Education",
				"os_version": "10.0.26100.0",
				"os_bits": "64-bit",
				"processor_arch": "AMD64"
				}
				`,
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
			mockTransport := transport.NewWinMockTransport()
			mockTransport.PowerShellResults["Write-Host $PSVersionTable.PSVersion"] = &transport.MockCmd{
				Stdout: "5.1.19041.1237",
			}
			mockTransport.PowerShellResults[osWindowsDiscoveryScript] = &transport.MockCmd{
				Stdout: tt.output,
			}

			info := newOSInfo()

			diags := info.populateOSInfo(mockTransport)
			if diags.HasErrors() {
				t.Fatalf("expected no errors, got: %v", diags.Errors())
			}

			if diags.HasWarnings() {
				t.Fatalf("expected no warnings, got: %v", diags.Warnings())
			}

			expectedKernel := "windows"
			if !info.Families().Contains(expectedKernel) {
				t.Errorf("expected family %q to be added, but it was not", expectedKernel)
			}

			if !info.Families().Contains(tt.expectedID) {
				t.Errorf("expected family %q to be added, but it was not", tt.expectedID)
			}

			if info.Kernel() != expectedKernel {
				t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
			}

			if info.ID() != tt.expectedID {
				t.Errorf("expected OS ID to be %q, got: %q", tt.expectedID, info.ID())
			}

			if info.FriendlyName() != tt.expectedFriendlyName {
				t.Errorf("expected friendly name to be %q, got: %q", tt.expectedFriendlyName, info.FriendlyName())
			}

			if info.Release() != tt.expectedRelease {
				t.Errorf("expected release to be %q, got: %q", tt.expectedRelease, info.Release())
			}

			if info.MajorVersion() != tt.expectedMajorVersion {
				t.Errorf("expected major version to be %q, got: %q", tt.expectedMajorVersion, info.MajorVersion())
			}

			if info.Version() != tt.expectedVersion {
				t.Errorf("expected version to be %q, got: %q", tt.expectedVersion, info.Version())
			}

			if info.Edition() != tt.expectedEdition {
				t.Errorf("expected edition to be %q, got: %q", tt.expectedEdition, info.Edition())
			}

			if info.EditionID() != tt.expectedEditionId {
				t.Errorf("expected edition ID to be %q, got: %q", tt.expectedEditionId, info.EditionID())
			}
		})
	}
}

func TestOSInfo_PopulateOSInfo_Windows_Error(t *testing.T) {
	mockTransport := transport.NewWinMockTransport()
	mockTransport.PowerShellResults["Write-Host $PSVersionTable.PSVersion"] = &transport.MockCmd{
		Stdout: "5.1.19041.1237",
	}
	mockTransport.PowerShellResults[osWindowsDiscoveryScript] = &transport.MockCmd{
		Err: os.ErrPermission,
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Families().Size() != 1 {
		t.Errorf("expected 1 family, got: %d", info.Families().Size())
	}

	expectedKernel := "windows"
	if !info.Families().Contains(expectedKernel) {
		t.Errorf("expected family %q to be present, but it was not", expectedKernel)
	}

	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	if info.ID() != "" {
		t.Errorf("expected ID to be empty, got: %q", info.ID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected processor architecture to be empty, got: %q", info.ProcArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected processor architecture bits to be 0, got: %d", info.ProcArchBits())
	}

	if info.OSArch() != "" {
		t.Errorf("expected OS architecture to be empty, got: %q", info.OSArch())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected OS architecture bits to be 0, got: %d", info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to get Windows OS information"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary to be %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error executing Windows discovery script: permission denied"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail to be %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_PopulateOSInfo_Windows_NotJSON(t *testing.T) {
	mockTransport := transport.NewWinMockTransport()
	mockTransport.PowerShellResults["Write-Host $PSVersionTable.PSVersion"] = &transport.MockCmd{
		Stdout: "5.1.19041.1237",
	}
	mockTransport.PowerShellResults[osWindowsDiscoveryScript] = &transport.MockCmd{
		Stdout: "This is not JSON output",
	}

	info := newOSInfo()

	diags := info.populateOSInfo(mockTransport)
	if !diags.HasErrors() {
		t.Fatalf("expected errors, got none")
	}

	if diags.HasWarnings() {
		t.Fatalf("expected no warnings, got: %v", diags.Warnings())
	}

	if info.Families().Size() != 1 {
		t.Errorf("expected 1 family, got: %d", info.Families().Size())
	}

	expectedKernel := "windows"
	if !info.Families().Contains(expectedKernel) {
		t.Errorf("expected family %q to be present, but it was not", expectedKernel)
	}

	if info.Kernel() != expectedKernel {
		t.Errorf("expected kernel to be %q, got: %q", expectedKernel, info.Kernel())
	}

	if info.ID() != "" {
		t.Errorf("expected ID to be empty, got: %q", info.ID())
	}

	if info.ProcArch() != "" {
		t.Errorf("expected processor architecture to be empty, got: %q", info.ProcArch())
	}

	if info.ProcArchBits() != 0 {
		t.Errorf("expected processor architecture bits to be 0, got: %d", info.ProcArchBits())
	}

	if info.OSArch() != "" {
		t.Errorf("expected OS architecture to be empty, got: %q", info.OSArch())
	}

	if info.OSArchBits() != 0 {
		t.Errorf("expected OS architecture bits to be 0, got: %d", info.OSArchBits())
	}

	errors := diags.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got: %d", len(errors))
	}

	expectedSummary := "Failed to parse Windows discovery output"
	if errors[0].Summary != expectedSummary {
		t.Errorf("expected error summary to be %q, got: %q", expectedSummary, errors[0].Summary)
	}

	expectedDetail := "Error parsing JSON output: invalid character 'T' looking for beginning of value"
	if errors[0].Detail != expectedDetail {
		t.Errorf("expected error detail to be %q, got: %q", expectedDetail, errors[0].Detail)
	}
}

func TestOSInfo_ToMapOfCtyValues(t *testing.T) {
	info := newOSInfo()
	info.families.Add("linux")
	info.families.Add("debian")
	info.families.Add("ubuntu")
	info.kernel = "linux"
	info.id = "ubuntu"
	info.friendlyName = "Ubuntu 22.04.3 LTS"
	info.release = "jammy"
	info.majorVersion = "22"
	info.version = "22.04"
	info.edition = "LTS"
	info.editionID = "lts"
	info.osArch = "amd64"
	info.osArchBits = 64
	info.procArch = "amd64"
	info.procArchBits = 64

	values := info.toMapOfCtyValues()
	if values["os_families"].Type() != cty.Set(cty.String) {
		t.Errorf("expected os_families to be a set of strings, got: %s", values["os_families"].Type().GoString())
	}

	families := values["os_families"].AsValueSlice()
	if len(families) != 3 {
		t.Errorf("expected 3 families, got: %d", len(families))
	}

	for _, family := range families {
		if family.Type() != cty.String {
			t.Errorf("expected family to be a string, got: %s", family.Type().GoString())
		}

		if family.AsString() != "linux" && family.AsString() != "debian" && family.AsString() != "ubuntu" {
			t.Errorf("unexpected family value: %s", family.AsString())
		}
	}

	if values["os_kernel"].Type() != cty.String {
		t.Errorf("expected os_kernel to be a string, got: %s", values["os_kernel"].Type().GoString())
	}

	if values["os_kernel"].AsString() != "linux" {
		t.Errorf("expected os_kernel to be 'linux', got: %s", values["os_kernel"].AsString())
	}

	if values["os_id"].Type() != cty.String {
		t.Errorf("expected os_id to be a string, got: %s", values["os_id"].Type().GoString())
	}

	if values["os_id"].AsString() != "ubuntu" {
		t.Errorf("expected os_id to be 'ubuntu', got: %s", values["os_id"].AsString())
	}

	if values["os_friendly_name"].Type() != cty.String {
		t.Errorf("expected os_friendly_name to be a string, got: %s", values["os_friendly_name"].Type().GoString())
	}

	if values["os_friendly_name"].AsString() != "Ubuntu 22.04.3 LTS" {
		t.Errorf("expected os_friendly_name to be 'Ubuntu 22.04.3 LTS', got: %s", values["os_friendly_name"].AsString())
	}

	if values["os_release"].Type() != cty.String {
		t.Errorf("expected os_release to be a string, got: %s", values["os_release"].Type().GoString())
	}

	if values["os_release"].AsString() != "jammy" {
		t.Errorf("expected os_release to be 'jammy', got: %s", values["os_release"].AsString())
	}

	if values["os_major_version"].Type() != cty.String {
		t.Errorf("expected os_major_version to be a string, got: %s", values["os_major_version"].Type().GoString())
	}

	if values["os_major_version"].AsString() != "22" {
		t.Errorf("expected os_major_version to be '22', got: %s", values["os_major_version"].AsString())
	}

	if values["os_version"].Type() != cty.String {
		t.Errorf("expected os_version to be a string, got: %s", values["os_version"].Type().GoString())
	}

	if values["os_version"].AsString() != "22.04" {
		t.Errorf("expected os_version to be '22.04', got: %s", values["os_version"].AsString())
	}

	if values["os_edition"].Type() != cty.String {
		t.Errorf("expected os_edition to be a string, got: %s", values["os_edition"].Type().GoString())
	}

	if values["os_edition"].AsString() != "LTS" {
		t.Errorf("expected os_edition to be 'LTS', got: %s", values["os_edition"].AsString())
	}

	if values["os_edition_id"].Type() != cty.String {
		t.Errorf("expected os_edition_id to be a string, got: %s", values["os_edition_id"].Type().GoString())
	}

	if values["os_edition_id"].AsString() != "lts" {
		t.Errorf("expected os_edition_id to be 'lts', got: %s", values["os_edition_id"].AsString())
	}

	if values["os_architecture"].Type() != cty.String {
		t.Errorf("expected os_architecture to be a string, got: %s", values["os_architecture"].Type().GoString())
	}

	if values["os_architecture"].AsString() != "amd64" {
		t.Errorf("expected os_architecture to be 'amd64', got: %s", values["os_architecture"].AsString())
	}

	if values["processor_architecture"].Type() != cty.String {
		t.Errorf("expected processor_architecture to be a string, got: %s", values["processor_architecture"].Type().GoString())
	}

	if values["processor_architecture"].AsString() != "amd64" {
		t.Errorf("expected processor_architecture to be 'amd64', got: %s", values["processor_architecture"].AsString())
	}

	if values["os_architecture_bits"].Type() != cty.Number {
		t.Errorf("expected os_architecture_bits to be a number, got: %s", values["os_architecture_bits"].Type().GoString())
	}

	value, _ := values["os_architecture_bits"].AsBigFloat().Int64()
	if value != 64 {
		t.Errorf("expected os_architecture_bits to be 64, got: %s", values["os_architecture_bits"].AsString())
	}

	if values["processor_architecture_bits"].Type() != cty.Number {
		t.Errorf("expected processor_architecture_bits to be a number, got: %s", values["processor_architecture_bits"].Type().GoString())
	}

	value, _ = values["processor_architecture_bits"].AsBigFloat().Int64()
	if value != 64 {
		t.Errorf("expected processor_architecture_bits to be 64, got: %s", values["processor_architecture_bits"].AsString())
	}
}

func TestOSInfo_ToMapOfCtyValues_EmptyValues(t *testing.T) {
	info := newOSInfo()

	values := info.toMapOfCtyValues()

	numberKeys := []string{
		"os_architecture_bits",
		"processor_architecture_bits",
	}

	stringKeys := []string{
		"os_kernel",
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
				t.Errorf("expected %s to be null, got: %s", key, value.GoString())
			}
			if value.Type() != cty.Number {
				t.Errorf("expected %s to be of type Number, got: %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	for _, key := range stringKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected %s to be null, got: %s", key, value.GoString())
			}
			if value.Type() != cty.String {
				t.Errorf("expected %s to be of type String, got: %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}

	for _, key := range setOfStringsKeys {
		if value, exists := values[key]; exists {
			if !value.IsNull() {
				t.Errorf("expected %s to be null, got: %s", key, value.GoString())
			}
			if value.Type() != cty.Set(cty.String) {
				t.Errorf("expected %s to be of type Set(String), got %s", key, value.Type().GoString())
			}
		} else {
			t.Errorf("expected key %q to be present in values map", key)
		}
	}
}
