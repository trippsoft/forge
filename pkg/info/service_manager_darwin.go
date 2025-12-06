// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

func discoverServiceManagerInfo() (*ServiceManagerInfoResponse, error) {
	return &ServiceManagerInfoResponse{
		Name: "launchd",
	}, nil
}
