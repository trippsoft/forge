// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !linux

package info

func (a *AppArmorInfo) discover() []string {
	a.Supported = false
	a.Enabled = false
	return nil
}
