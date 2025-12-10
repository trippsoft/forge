// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux && !windows

package info

func discoverFIPSInfo() (*FIPSInfo, error) {
	return &FIPSInfo{
		Known:   false,
		Enabled: false,
	}, nil
}
