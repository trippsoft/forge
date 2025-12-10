// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

func discoverServiceManagerInfo() (*ServiceManagerInfo, error) {
	return &ServiceManagerInfo{
		Name: "windows",
	}, nil
}
