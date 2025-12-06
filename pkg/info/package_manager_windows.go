// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

func discoverPackageManagerInfo(_ *OSInfoResponse) (*PackageManagerInfoResponse, error) {
	// Windows package managers are not handled in this implementation.
	// Winget and Chocolatey implementations are to be separate.
	return &PackageManagerInfoResponse{
		Name: "",
		Path: "",
	}, nil
}
