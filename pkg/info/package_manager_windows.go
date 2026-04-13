// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

func (p *PackageManagerInfo) discover(_ *OSInfo) error {
	// Windows package managers are not handled in this implementation.
	// Winget and Chocolatey implementations are to be separate.
	p.Name = ""
	p.Path = ""
	return nil
}
