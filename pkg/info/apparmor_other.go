// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func (a *AppArmorInfoPB) discover() error {
	a.Supported = false
	a.Enabled = false
	return nil
}
