// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build aix

package info

func discoverServiceManagerInfo() (*ServiceManagerInfoResponse, error) {
	return &ServiceManagerInfoResponse{
		Name: "src",
	}, nil
}
