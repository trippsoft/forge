// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux && !windows

package info

func (f *FIPSInfo) discover() error {
	f.Known = false
	f.Enabled = false
	return nil
}
