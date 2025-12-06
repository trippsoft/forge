// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux && !windows

package info

func discoverFIPSInfo() (*FIPSInfoResponse, error) {
	return &FIPSInfoResponse{
		Known:   false,
		Enabled: false,
	}, nil
}
