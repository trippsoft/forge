// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package version provides version information for the Forge application.
package version

import (
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed VERSION
	RawVersion string

	MajorVersion  uint
	MinorVersion  uint
	PatchVersion  uint
	VersionSuffix string
)

func init() {
	parts := strings.SplitN(strings.TrimSpace(RawVersion), "-", 2)
	if len(parts) == 2 {
		VersionSuffix = parts[1]
	}

	parts = strings.SplitN(parts[0], ".", 3)
	switch len(parts) {
	case 3:
		fmt.Sscanf(parts[2], "%d", &PatchVersion)
		fallthrough
	case 2:
		fmt.Sscanf(parts[1], "%d", &MinorVersion)
		fallthrough
	case 1:
		fmt.Sscanf(parts[0], "%d", &MajorVersion)
	}
}

// Version returns the version number of Forge as a string.
func Version() string {
	if VersionSuffix != "" {
		return fmt.Sprintf("%d.%d.%d-%s", MajorVersion, MinorVersion, PatchVersion, VersionSuffix)
	}

	return fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)
}
