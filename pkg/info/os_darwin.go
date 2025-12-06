// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package info

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func discoverOSInfo() (*OSInfoResponse, error) {
	osInfo := &OSInfoResponse{
		Id:        "macos",
		Edition:   "",
		EditionId: "",
		Families:  []string{"posix", "darwin", "macos"},
	}

	cmd := exec.Command("sw_vers", "-productVersion")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	osInfo.Version = string(bytes.TrimSpace(stdout.Bytes()))
	osInfo.FriendlyName = "macOS " + osInfo.Version
	versionParts := strings.Split(osInfo.Version, ".")
	if len(versionParts) < 2 {
		return nil, errors.New("invalid version format")
	}

	osInfo.MajorVersion = versionParts[0]
	minorVersion := versionParts[1]

	switch osInfo.MajorVersion {
	case "10":
		switch minorVersion {
		case "6":
			osInfo.Release = "Snow Leopard"
			osInfo.ReleaseId = "snow-leopard"
		case "7":
			osInfo.Release = "Lion"
			osInfo.ReleaseId = "lion"
		case "8":
			osInfo.Release = "Mountain Lion"
			osInfo.ReleaseId = "mountain-lion"
		case "9":
			osInfo.Release = "Mavericks"
			osInfo.ReleaseId = "mavericks"
		case "10":
			osInfo.Release = "Yosemite"
			osInfo.ReleaseId = "yosemite"
		case "11":
			osInfo.Release = "El Capitan"
			osInfo.ReleaseId = "el-capitan"
		case "12":
			osInfo.Release = "Sierra"
			osInfo.ReleaseId = "sierra"
		case "13":
			osInfo.Release = "High Sierra"
			osInfo.ReleaseId = "high-sierra"
		case "14":
			osInfo.Release = "Mojave"
			osInfo.ReleaseId = "mojave"
		case "15":
			osInfo.Release = "Catalina"
			osInfo.ReleaseId = "catalina"
		default:
			return nil, errors.New("unknown macOS release for version 10." + minorVersion)
		}
	case "11":
		osInfo.Release = "Big Sur"
		osInfo.ReleaseId = "big-sur"
	case "12":
		osInfo.Release = "Monterey"
		osInfo.ReleaseId = "monterey"
	case "13":
		osInfo.Release = "Ventura"
		osInfo.ReleaseId = "ventura"
	case "14":
		osInfo.Release = "Sonoma"
		osInfo.ReleaseId = "sonoma"
	case "15":
		osInfo.Release = "Sequoia"
		osInfo.ReleaseId = "sequoia"
	case "26":
		osInfo.Release = "Tahoe"
		osInfo.ReleaseId = "tahoe"
	default:
		return nil, errors.New("unknown macOS major version: " + osInfo.MajorVersion)
	}

	return osInfo, nil
}
