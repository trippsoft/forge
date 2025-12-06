// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func discoverAppArmorInfo() (*AppArmorInfoResponse, error) {
	return &AppArmorInfoResponse{
		Supported: false,
		Enabled:   false,
	}, nil
}
