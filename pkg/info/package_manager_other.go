// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !darwin && !windows

package info

import (
	"bytes"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func discoverPackageManagerInfo(osInfo *OSInfoResponse) (*PackageManagerInfoResponse, error) {
	packageManagerInfo := &PackageManagerInfoResponse{}

	if slices.Contains(osInfo.Families, "archlinux") {
		err := populateArchLinuxPackageManagerInfo(packageManagerInfo)
		if err != nil {
			return nil, err
		}
	}

	if slices.Contains(osInfo.Families, "debian") || slices.Contains(osInfo.Families, "altlinux") {
		err := populateDebianPackageManagerInfo(packageManagerInfo)
		if err != nil {
			return nil, err
		}

		return packageManagerInfo, nil
	}

	if slices.Contains(osInfo.Families, "el") {
		err := populateEnterpriseLinuxPackageManagerInfo(packageManagerInfo)
		if err != nil {
			return nil, err
		}

		return packageManagerInfo, nil
	}

	if slices.Contains(osInfo.Families, "gentoo") {
		err := populateGentooPackageManagerInfo(packageManagerInfo)
		if err != nil {
			return nil, err
		}

		return packageManagerInfo, nil
	}

	if slices.Contains(osInfo.Families, "suse") {
		err := populateSUSEPackageManagerInfo(packageManagerInfo)
		if err != nil {
			return nil, err
		}

		return packageManagerInfo, nil
	}

	err := populateOtherLinuxPackageManagerInfo(packageManagerInfo)
	if err != nil {
		return nil, err
	}

	return packageManagerInfo, nil
}

func populateArchLinuxPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/usr/bin/pacman")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pacman"
		packageManagerInfo.Path = "/usr/bin/pacman"
		return nil
	}

	return populateOtherLinuxPackageManagerInfo(packageManagerInfo)
}

func populateDebianPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/usr/bin/apt-get")
	if err == nil && fileInfo.Mode().IsRegular() {
		fileInfo, err := os.Stat("/usr/bin/rpm")
		if err == nil && fileInfo.Mode().IsRegular() {
			cmd := exec.Command("/usr/bin/rpm", "-q", "--whatprovides", "/usr/bin/apt-get")
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			err := cmd.Run()
			if err == nil {
				output := strings.TrimSpace(stdout.String())
				if output != "" {
					packageManagerInfo.Name = "apt-rpm"
					packageManagerInfo.Path = "/usr/bin/apt-get"
					return nil
				}
			}
		}

		packageManagerInfo.Name = "apt"
		packageManagerInfo.Path = "/usr/bin/apt-get"
		return nil
	}

	return populateOtherLinuxPackageManagerInfo(packageManagerInfo)
}

func populateEnterpriseLinuxPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/usr/bin/dnf5")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf5"
		packageManagerInfo.Path = "/usr/bin/dnf5"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf-3")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf"
		packageManagerInfo.Path = "/usr/bin/dnf-3"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf"
		packageManagerInfo.Path = "/usr/bin/dnf"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "yum"
		packageManagerInfo.Path = "/usr/bin/yum"
		return nil
	}

	return populateOtherLinuxPackageManagerInfo(packageManagerInfo)
}

func populateGentooPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/usr/bin/emerge")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "emerge"
		packageManagerInfo.Path = "/usr/bin/emerge"
		return nil
	}

	return populateOtherLinuxPackageManagerInfo(packageManagerInfo)
}

func populateSUSEPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/usr/bin/zypper")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "zypper"
		packageManagerInfo.Path = "/usr/bin/zypper"
		return nil
	}

	return populateOtherLinuxPackageManagerInfo(packageManagerInfo)
}

func populateOtherLinuxPackageManagerInfo(packageManagerInfo *PackageManagerInfoResponse) error {
	fileInfo, err := os.Stat("/QOpenSys/pkgs/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "yum"
		packageManagerInfo.Path = "/QOpenSys/pkgs/bin/yum"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/installp")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "installp"
		packageManagerInfo.Path = "/usr/bin/installp"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/sorcery")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "sorcery"
		packageManagerInfo.Path = "/usr/sbin/sorcery"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/swupd")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "swupd"
		packageManagerInfo.Path = "/usr/bin/swupd"
		return nil
	}

	fileInfo, err = os.Stat("/usr/local/sbin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkgng"
		packageManagerInfo.Path = "/usr/local/sbin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/xbps-install")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "xbps"
		packageManagerInfo.Path = "/usr/bin/xbps-install"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkg5"
		packageManagerInfo.Path = "/usr/bin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/pkgadd")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "svr4pkg"
		packageManagerInfo.Path = "/usr/sbin/pkgadd"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/emerge")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "portage"
		packageManagerInfo.Path = "/usr/bin/emerge"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/swlist")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "swdepot"
		packageManagerInfo.Path = "/usr/sbin/swlist"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkgng"
		packageManagerInfo.Path = "/usr/sbin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/sbin/apk")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "apk"
		packageManagerInfo.Path = "/sbin/apk"
		return nil
	}

	fileInfo, err = os.Stat("/opt/homebrew/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "homebrew"
		packageManagerInfo.Path = "/opt/homebrew/bin/brew"
		return nil
	}

	fileInfo, err = os.Stat("/usr/local/bin/brew")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "homebrew"
		packageManagerInfo.Path = "/usr/local/bin/brew"
		return nil
	}

	fileInfo, err = os.Stat("/opt/local/bin/port")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "macports"
		packageManagerInfo.Path = "/opt/local/bin/port"
		return nil
	}

	fileInfo, err = os.Stat("/opt/tools/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkgin"
		packageManagerInfo.Path = "/opt/tools/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/opt/local/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkgin"
		packageManagerInfo.Path = "/opt/local/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/usr/pkg/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pkgin"
		packageManagerInfo.Path = "/usr/pkg/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/bin/opkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "opkg"
		packageManagerInfo.Path = "/bin/opkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/pacman")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "pacman"
		packageManagerInfo.Path = "/usr/bin/pacman"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/urpmi")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "urpmi"
		packageManagerInfo.Path = "/usr/sbin/urpmi"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/zypper")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "zypper"
		packageManagerInfo.Path = "/usr/bin/zypper"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/apt-get")
	if err == nil && fileInfo.Mode().IsRegular() {
		fileInfo, err := os.Stat("/usr/bin/rpm")
		if err == nil && fileInfo.Mode().IsRegular() {
			cmd := exec.Command("/usr/bin/rpm", "-q", "--whatprovides", "/usr/bin/apt-get")
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			err := cmd.Run()
			if err == nil {
				output := strings.TrimSpace(stdout.String())
				if output != "" {
					packageManagerInfo.Name = "apt-rpm"
					packageManagerInfo.Path = "/usr/bin/apt-get"
					return nil
				}
			}
		}

		packageManagerInfo.Name = "apt"
		packageManagerInfo.Path = "/usr/bin/apt-get"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf5")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf5"
		packageManagerInfo.Path = "/usr/bin/dnf5"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf-3")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf"
		packageManagerInfo.Path = "/usr/bin/dnf-3"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "dnf"
		packageManagerInfo.Path = "/usr/bin/dnf"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		packageManagerInfo.Name = "yum"
		packageManagerInfo.Path = "/usr/bin/yum"
		return nil
	}

	packageManagerInfo.Name = ""
	packageManagerInfo.Path = ""
	return nil
}
