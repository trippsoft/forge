// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

func discoverPackageManagerInfo(_ *OSInfo) (*PackageManagerInfo, error) {
	// Windows package managers are not handled in this implementation.
	// Winget and Chocolatey implementations are to be separate.
	return &PackageManagerInfo{
		Name: "",
		Path: "",
	}, nil
}
