// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package discover

func discoverSELinuxInfo() (*SELinuxInfoResponse, error) {
	return &SELinuxInfoResponse{
		Supported: false,
		Installed: false,
		Status:    "",
		Type:      "",
	}, nil
}
