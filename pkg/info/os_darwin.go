// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func (o *OSInfo) discover() error {
	o.Kernel = "darwin"
	o.Id = "macos"
	o.Edition = ""
	o.EditionId = ""
	o.Arch = runtime.GOARCH
	o.Families = []string{"posix", "darwin", "macos"}

	cmd := exec.Command("sw_vers", "-productVersion")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	o.Version = string(bytes.TrimSpace(stdout.Bytes()))
	o.FriendlyName = "macOS " + o.Version
	versionParts := strings.Split(o.Version, ".")
	if len(versionParts) < 2 {
		return errors.New("invalid version format")
	}

	o.MajorVersion = versionParts[0]
	minorVersion := versionParts[1]

	switch o.MajorVersion {
	case "10":
		switch minorVersion {
		case "6":
			o.Release = "Snow Leopard"
			o.ReleaseId = "snow-leopard"
		case "7":
			o.Release = "Lion"
			o.ReleaseId = "lion"
		case "8":
			o.Release = "Mountain Lion"
			o.ReleaseId = "mountain-lion"
		case "9":
			o.Release = "Mavericks"
			o.ReleaseId = "mavericks"
		case "10":
			o.Release = "Yosemite"
			o.ReleaseId = "yosemite"
		case "11":
			o.Release = "El Capitan"
			o.ReleaseId = "el-capitan"
		case "12":
			o.Release = "Sierra"
			o.ReleaseId = "sierra"
		case "13":
			o.Release = "High Sierra"
			o.ReleaseId = "high-sierra"
		case "14":
			o.Release = "Mojave"
			o.ReleaseId = "mojave"
		case "15":
			o.Release = "Catalina"
			o.ReleaseId = "catalina"
		default:
			return errors.New("unknown macOS release for version 10." + minorVersion)
		}
	case "11":
		o.Release = "Big Sur"
		o.ReleaseId = "big-sur"
	case "12":
		o.Release = "Monterey"
		o.ReleaseId = "monterey"
	case "13":
		o.Release = "Ventura"
		o.ReleaseId = "ventura"
	case "14":
		o.Release = "Sonoma"
		o.ReleaseId = "sonoma"
	case "15":
		o.Release = "Sequoia"
		o.ReleaseId = "sequoia"
	case "26":
		o.Release = "Tahoe"
		o.ReleaseId = "tahoe"
	default:
		return errors.New("unknown macOS major version: " + o.MajorVersion)
	}

	return nil
}
