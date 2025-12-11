// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

func (s *ServiceManagerInfoPB) discover() error {
	s.Name = "launchd"
	return nil
}
