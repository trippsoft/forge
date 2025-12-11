// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !aix && !darwin && !linux && !windows

package info

import "os"

func (s *ServiceManagerInfoPB) discover() error {
	s.Name = ""

	fileInfo, err := os.Lstat("/sbin/init")
	if err != nil {
		return err
	}

	if (fileInfo.Mode() & os.ModeSymlink) == 0 {
		// Assume BSD-style init system for other Unix-like OSes
		s.Name = "bsdinit"
		return nil
	}

	return nil
}
