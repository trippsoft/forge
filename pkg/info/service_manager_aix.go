// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build aix

package info

func (s *ServiceManagerInfo) discover() error {
	s.Name = "src"
	return nil
}
