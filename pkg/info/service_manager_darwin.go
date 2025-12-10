// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

func discoverServiceManagerInfo() (*ServiceManagerInfo, error) {
	return &ServiceManagerInfo{
		Name: "launchd",
	}, nil
}
