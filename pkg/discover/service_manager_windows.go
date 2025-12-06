// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package discover

func discoverServiceManagerInfo() (*ServiceManagerInfoResponse, error) {
	return &ServiceManagerInfoResponse{
		Name: "windows",
	}, nil
}
