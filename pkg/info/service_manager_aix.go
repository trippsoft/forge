// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build aix

package info

func (s *ServiceManagerInfoPB) discover() error {
	s.Name = "src"
	return nil
}
