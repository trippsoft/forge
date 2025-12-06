// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package discover

import (
	"os"
	"os/exec"
	"strings"
)

func discoverServiceManagerInfo() (*ServiceManagerInfoResponse, error) {
	serviceManagerInfo := &ServiceManagerInfoResponse{}

	_, err := exec.LookPath("systemctl")
	if err == nil {
		fileInfo, err := os.Stat("/run/systemd/system")
		if err == nil && fileInfo.Mode().IsDir() {
			serviceManagerInfo.Name = "systemd"
			return serviceManagerInfo, nil
		}

		fileInfo, err = os.Stat("/dev/.run/systemd")
		if err == nil && fileInfo.Mode().IsDir() {
			serviceManagerInfo.Name = "systemd"
			return serviceManagerInfo, nil
		}

		fileInfo, err = os.Stat("/dev/.systemd")
		if err == nil && fileInfo.Mode().IsDir() {
			serviceManagerInfo.Name = "systemd"
			return serviceManagerInfo, nil
		}
	}

	_, err = exec.LookPath("initctl")
	if err == nil {
		fileInfo, err := os.Stat("/etc/init")
		if err == nil && fileInfo.Mode().IsDir() {
			serviceManagerInfo.Name = "upstart"
			return serviceManagerInfo, nil
		}
	}

	fileInfo, err := os.Stat("/sbin/openrc")
	if err == nil && fileInfo.Mode().IsRegular() {
		serviceManagerInfo.Name = "openrc"
		return serviceManagerInfo, nil
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
		serviceManagerInfo.Name = "systemd"
		return serviceManagerInfo, nil
	}

	fileInfo, err = os.Stat("/etc/init.d")
	if err == nil && fileInfo.Mode().IsDir() {
		serviceManagerInfo.Name = "sysvinit"
		return serviceManagerInfo, nil
	}

	fileInfo, err = os.Stat("/proc/1/comm")
	if err == nil && fileInfo.Mode().IsRegular() {
		data, err := os.ReadFile("/proc/1/comm")
		if err == nil {
			comm := strings.TrimSpace(string(data))
			switch {
			case comm == "openrc-init":
				serviceManagerInfo.Name = "openrc"
				return serviceManagerInfo, nil
			case comm == "", comm == "COMMAND", comm == "init", strings.HasSuffix(comm, "sh"):
				// Ignore these common non-informative values
			default:
				serviceManagerInfo.Name = comm
				return serviceManagerInfo, nil
			}
		}
	}

	switch {
	case initLinkTarget == "openrc-init":
		serviceManagerInfo.Name = "openrc"
		return serviceManagerInfo, nil
	case initLinkTarget == "", initLinkTarget == "init", strings.HasSuffix(initLinkTarget, "sh"):
		// Ignore these common non-informative values
	default:
		serviceManagerInfo.Name = initLinkTarget
		return serviceManagerInfo, nil
	}

	return serviceManagerInfo, nil
}
