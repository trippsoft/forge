// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !aix && !darwin && !linux && !windows

package info

import "os"

func discoverServiceManagerInfo() (*ServiceManagerInfoResponse, error) {
	serviceManagerInfo := &ServiceManagerInfoResponse{}

	fileInfo, err := os.Lstat("/sbin/init")
	if err != nil {
		return nil, err
	}

	if (fileInfo.Mode() & os.ModeSymlink) == 0 {
		// Assume BSD-style init system for other Unix-like OSes
		serviceManagerInfo.Name = "bsdinit"
		return serviceManagerInfo, nil
	}

	return serviceManagerInfo, nil
}
