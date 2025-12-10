// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func discoverAppArmorInfo() (*AppArmorInfo, error) {
	return &AppArmorInfo{
		Supported: false,
		Enabled:   false,
	}, nil
}
