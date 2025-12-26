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

func (p *PackageManagerInfo) discover(osInfo *OSInfo) error {
	if slices.Contains(osInfo.Families, "archlinux") {
		err := p.populateArchLinuxPackageManagerInfo()
		if err != nil {
			return err
		}
	}

	if slices.Contains(osInfo.Families, "debian") || slices.Contains(osInfo.Families, "altlinux") {
		err := p.populateDebianPackageManagerInfo()
		if err != nil {
			return err
		}

		return nil
	}

	if slices.Contains(osInfo.Families, "el") {
		err := p.populateEnterpriseLinuxPackageManagerInfo()
		if err != nil {
			return err
		}

		return nil
	}

	if slices.Contains(osInfo.Families, "gentoo") {
		err := p.populateGentooPackageManagerInfo()
		if err != nil {
			return err
		}

		return nil
	}

	if slices.Contains(osInfo.Families, "suse") {
		err := p.populateSUSEPackageManagerInfo()
		if err != nil {
			return err
		}

		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateArchLinuxPackageManagerInfo() error {
	fileInfo, err := os.Stat("/usr/bin/pacman")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pacman"
		p.Path = "/usr/bin/pacman"
		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateDebianPackageManagerInfo() error {
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
					p.Name = "apt-rpm"
					p.Path = "/usr/bin/apt-get"
					return nil
				}
			}
		}

		p.Name = "apt"
		p.Path = "/usr/bin/apt-get"
		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateEnterpriseLinuxPackageManagerInfo() error {
	fileInfo, err := os.Stat("/usr/bin/dnf5")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf5"
		p.Path = "/usr/bin/dnf5"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf-3")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf"
		p.Path = "/usr/bin/dnf-3"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf"
		p.Path = "/usr/bin/dnf"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "yum"
		p.Path = "/usr/bin/yum"
		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateGentooPackageManagerInfo() error {
	fileInfo, err := os.Stat("/usr/bin/emerge")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "emerge"
		p.Path = "/usr/bin/emerge"
		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateSUSEPackageManagerInfo() error {
	fileInfo, err := os.Stat("/usr/bin/zypper")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "zypper"
		p.Path = "/usr/bin/zypper"
		return nil
	}

	return p.populateOtherLinuxPackageManagerInfo()
}

func (p *PackageManagerInfo) populateOtherLinuxPackageManagerInfo() error {
	fileInfo, err := os.Stat("/QOpenSys/pkgs/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "yum"
		p.Path = "/QOpenSys/pkgs/bin/yum"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/installp")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "installp"
		p.Path = "/usr/bin/installp"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/sorcery")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "sorcery"
		p.Path = "/usr/sbin/sorcery"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/swupd")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "swupd"
		p.Path = "/usr/bin/swupd"
		return nil
	}

	fileInfo, err = os.Stat("/usr/local/sbin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkgng"
		p.Path = "/usr/local/sbin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/xbps-install")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "xbps"
		p.Path = "/usr/bin/xbps-install"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkg5"
		p.Path = "/usr/bin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/pkgadd")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "svr4pkg"
		p.Path = "/usr/sbin/pkgadd"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/emerge")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "portage"
		p.Path = "/usr/bin/emerge"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/swlist")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "swdepot"
		p.Path = "/usr/sbin/swlist"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/pkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkgng"
		p.Path = "/usr/sbin/pkg"
		return nil
	}

	fileInfo, err = os.Stat("/sbin/apk")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "apk"
		p.Path = "/sbin/apk"
		return nil
	}

	fileInfo, err = os.Stat("/opt/homebrew/bin/brew")
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

	fileInfo, err = os.Stat("/opt/tools/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkgin"
		p.Path = "/opt/tools/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/opt/local/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkgin"
		p.Path = "/opt/local/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/usr/pkg/bin/pkgin")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pkgin"
		p.Path = "/usr/pkg/bin/pkgin"
		return nil
	}

	fileInfo, err = os.Stat("/bin/opkg")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "opkg"
		p.Path = "/bin/opkg"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/pacman")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "pacman"
		p.Path = "/usr/bin/pacman"
		return nil
	}

	fileInfo, err = os.Stat("/usr/sbin/urpmi")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "urpmi"
		p.Path = "/usr/sbin/urpmi"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/zypper")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "zypper"
		p.Path = "/usr/bin/zypper"
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
					p.Name = "apt-rpm"
					p.Path = "/usr/bin/apt-get"
					return nil
				}
			}
		}

		p.Name = "apt"
		p.Path = "/usr/bin/apt-get"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf5")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf5"
		p.Path = "/usr/bin/dnf5"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf-3")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf"
		p.Path = "/usr/bin/dnf-3"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/dnf")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "dnf"
		p.Path = "/usr/bin/dnf"
		return nil
	}

	fileInfo, err = os.Stat("/usr/bin/yum")
	if err == nil && fileInfo.Mode().IsRegular() {
		p.Name = "yum"
		p.Path = "/usr/bin/yum"
		return nil
	}

	p.Name = ""
	p.Path = ""
	return nil
}
