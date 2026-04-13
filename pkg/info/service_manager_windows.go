// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

func (s *ServiceManagerInfo) discover() error {
	s.Name = "windows"
	return nil
}
