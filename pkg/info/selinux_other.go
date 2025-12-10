// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func discoverSELinuxInfo() (*SELinuxInfo, error) {
	return &SELinuxInfo{
		Supported: false,
		Installed: false,
		Status:    "",
		Type:      "",
	}, nil
}
