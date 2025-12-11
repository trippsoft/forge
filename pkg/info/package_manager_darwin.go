// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

import "os"

func (p *PackageManagerInfoPB) discover(_ *OSInfoPB) error {
	fileInfo, err := os.Stat("/opt/homebrew/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "homebrew"
		p.Path = "/opt/homebrew/bin/brew"
		return nil
	}

	fileInfo, err = os.Stat("/usr/local/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "homebrew"
		p.Path = "/usr/local/bin/brew"
		return nil
	}

	fileInfo, err = os.Stat("/opt/local/bin/port")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "macports"
		p.Path = "/opt/local/bin/port"
		return nil
	}

	p.Name = ""
	p.Path = ""
	return nil
}
