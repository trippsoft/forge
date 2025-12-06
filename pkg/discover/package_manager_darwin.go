// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package discover

import "os"

func discoverPackageManagerInfo(_ *OSInfoResponse) (*PackageManagerInfoResponse, error) {
	packageManagerInfo := &PackageManagerInfoResponse{}

	fileInfo, err := os.Stat("/opt/homebrew/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "homebrew"
		packageManagerInfo.Path = "/opt/homebrew/bin/brew"
		return packageManagerInfo, nil
	}

	fileInfo, err = os.Stat("/usr/local/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "homebrew"
		packageManagerInfo.Path = "/usr/local/bin/brew"
		return packageManagerInfo, nil
	}

	fileInfo, err = os.Stat("/opt/local/bin/port")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "macports"
		packageManagerInfo.Path = "/opt/local/bin/port"
		return packageManagerInfo, nil
	}

	packageManagerInfo.Name = ""
	packageManagerInfo.Path = ""
	return packageManagerInfo, nil
}
