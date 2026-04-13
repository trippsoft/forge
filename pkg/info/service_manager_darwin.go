// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

func (s *ServiceManagerInfo) discover() error {
	s.Name = "launchd"
	return nil
}
