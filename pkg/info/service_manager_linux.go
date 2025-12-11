// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package info

import (
	"os"
	"os/exec"
	"strings"
)

func (s *ServiceManagerInfoPB) discover() error {
	s.Name = ""

	_, err := exec.LookPath("systemctl")
	if err == nil {
		fileInfo, err := os.Stat("/run/systemd/system")
		if err == nil && fileInfo.Mode().IsDir() {
			s.Name = "systemd"
			return nil
		}

		fileInfo, err = os.Stat("/dev/.run/systemd")
		if err == nil && fileInfo.Mode().IsDir() {
			s.Name = "systemd"
			return nil
		}

		fileInfo, err = os.Stat("/dev/.systemd")
		if err == nil && fileInfo.Mode().IsDir() {
			s.Name = "systemd"
			return nil
		}
	}

	_, err = exec.LookPath("initctl")
	if err == nil {
		fileInfo, err := os.Stat("/etc/init")
		if err == nil && fileInfo.Mode().IsDir() {
			s.Name = "upstart"
			return nil
		}
	}

	fileInfo, err := os.Stat("/sbin/openrc")
	if err == nil && fileInfo.Mode().IsRegular() {
		s.Name = "openrc"
		return nil
	}

	initLinkTarget := ""
	fileInfo, err = os.Lstat("/sbin/init")
	if err == nil && fileInfo.Mode()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink("/sbin/init")
		if err == nil {
			linkTargetParts := strings.Split(linkTarget, "/")
			initLinkTarget = linkTargetParts[len(linkTargetParts)-1]
		}
	}

	if initLinkTarget == "systemd" {
		s.Name = "systemd"
		return nil
	}

	fileInfo, err = os.Stat("/etc/init.d")
	if err == nil && fileInfo.Mode().IsDir() {
		s.Name = "sysvinit"
		return nil
	}

	fileInfo, err = os.Stat("/proc/1/comm")
	if err == nil && fileInfo.Mode().IsRegular() {
		data, err := os.ReadFile("/proc/1/comm")
		if err == nil {
			comm := strings.TrimSpace(string(data))
			switch {
			case comm == "openrc-init":
				s.Name = "openrc"
				return nil
			case comm == "", comm == "COMMAND", comm == "init", strings.HasSuffix(comm, "sh"):
				// Ignore these common non-informative values
			default:
				s.Name = comm
				return nil
			}
		}
	}

	switch {
	case initLinkTarget == "openrc-init":
		s.Name = "openrc"
		return nil
	case initLinkTarget == "", initLinkTarget == "init", strings.HasSuffix(initLinkTarget, "sh"):
		// Ignore these common non-informative values
	default:
		s.Name = initLinkTarget
		return nil
	}

	return nil
}
