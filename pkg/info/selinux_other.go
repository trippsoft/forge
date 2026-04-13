// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func (s *SELinuxInfo) discover() error {
	s.Supported = false
	s.Installed = false
	s.Status = ""
	s.Type = ""
	return nil
}
