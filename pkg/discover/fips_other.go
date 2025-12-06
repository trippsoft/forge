// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux && !windows

package discover

func discoverFIPSInfo() (*FIPSInfoResponse, error) {
	return &FIPSInfoResponse{
		Known:   false,
		Enabled: false,
	}, nil
}
