// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !darwin && !windows

package discover

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	osReleaseKeys = map[string]any{
		"ID":               nil,
		"VERSION_ID":       nil,
		"VERSION_CODENAME": nil,
		"PRETTY_NAME":      nil,
		"VARIANT":          nil,
		"VARIANT_ID":       nil,
	}
	osIDCorrectionMap = map[string]string{
		"amzn":           "amazon",
		"arch":           "archlinux",
		"archarm":        "archlinux-arm",
		"clear-linux-os": "clearlinux",
		"cumulus-linux":  "cumuluslinux",
		"pop":            "pop_os",
		"ol":             "oraclelinux",
		"opensuse-leap":  "opensuse",
		"sles_sap":       "sles",
	}
	osFamiliesMap = map[string][]string{
		"almalinux":     {"el"},
		"amazon":        {"el"},
		"archlinux-arm": {"archlinux"},
		"arcolinux":     {"archlinux"},
		"centos":        {"el"},
		"clearos":       {"el"},
		"cloudlinux":    {"el"},
		"deepin":        {"debian"},
		"devuan":        {"debian"},
		"elementary":    {"debian", "ubuntu"},
		"endeavouros":   {"archlinux"},
		"fedora":        {"el"},
		"kali":          {"debian"},
		"kylin":         {"debian", "ubuntu"},
		"linuxmint":     {"debian", "ubuntu"}, // We will treat Linux Mint as always being Ubuntu-based, despite the existence of Debian-based versions.
		"mageia":        {"mandrake"},
		"manjaro":       {"archlinux"},
		"manjaro-arm":   {"archlinux", "manjaro"},
		"nobara":        {"el", "fedora"},
		"opensuse":      {"suse"},
		"oraclelinux":   {"el"},
		"pop_os":        {"debian", "ubuntu"},
		"raspbian":      {"debian"},
		"rhel":          {"el"},
		"rocky":         {"el"},
		"scientific":    {"el"},
		"sled":          {"suse"},
		"sles":          {"suse"},
		"ubuntu":        {"debian"},
		"virtuozzo":     {"el"},
	}
)

func discoverOSInfo() (*OSInfoResponse, error) {
	osInfo := &OSInfoResponse{}
	osReleaseErr := populateFromOsReleaseFile(osInfo)
	lsbReleaseErr := populateFromLsbRelease(osInfo)

	if osReleaseErr != nil && lsbReleaseErr != nil {
		return nil, errors.Join(osReleaseErr, lsbReleaseErr)
	}

	if osInfo.Release != "" {
		osInfo.ReleaseId = strings.ReplaceAll(strings.ToLower(osInfo.Release), " ", "-")
	}

	if osInfo.Version != "" {
		versionParts := strings.SplitN(osInfo.Version, ".", 2)
		if len(versionParts) >= 1 {
			osInfo.MajorVersion = versionParts[0]
		}
	}

	correctedID, exists := osIDCorrectionMap[osInfo.Id]
	if exists {
		osInfo.Id = correctedID
	}

	families := []string{runtime.GOOS, "posix"}
	matchingFamilies, exists := osFamiliesMap[osInfo.Id]
	if exists {
		families = append(families, matchingFamilies...)
	}

	if osInfo.Id != "" {
		families = append(families, osInfo.Id)
	}
	osInfo.Families = families

	return osInfo, nil
}

func populateFromOsReleaseFile(osInfo *OSInfoResponse) error {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		file, err = os.Open("/usr/lib/os-release")
		if err != nil {
			return err
		}
	}
	defer file.Close()

	contents := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // Skip empty lines
		}

		if line[0] == '#' {
			continue // Skip comments
		}

		stringParts := strings.SplitN(line, "=", 2)
		if len(stringParts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(stringParts[0])
		_, exists := osReleaseKeys[key]
		if !exists {
			continue // Skip keys we're not interested in
		}

		value := strings.Trim(strings.TrimSpace(stringParts[1]), `"`)
		value = strings.ReplaceAll(value, `\"`, `"`)
		value = strings.ReplaceAll(value, `\$`, `$`)
		value = strings.ReplaceAll(value, `\\`, `\`)
		value = strings.ReplaceAll(value, "\\`", "`")

		contents[key] = value
	}

	id, ok := contents["ID"]
	if ok && id != "" && id != "n/a" {
		osInfo.Id = strings.ToLower(id)
	}

	friendlyName, ok := contents["PRETTY_NAME"]
	if ok && friendlyName != "" && friendlyName != "n/a" {
		osInfo.FriendlyName = friendlyName
	}

	release, ok := contents["VERSION_CODENAME"]
	if ok && release != "" && release != "n/a" {
		osInfo.Release = release
	}

	version, ok := contents["VERSION_ID"]
	if ok && version != "" && version != "n/a" {
		osInfo.Version = version
	}

	edition, ok := contents["VARIANT"]
	if ok && edition != "" && edition != "n/a" {
		osInfo.Edition = edition
	}

	editionId, ok := contents["VARIANT_ID"]
	if ok && editionId != "" && editionId != "n/a" {
		osInfo.EditionId = editionId
	}

	return nil
}

func populateFromLsbRelease(osInfo *OSInfoResponse) error {
	if osInfo.Id == "" {
		cmd := exec.Command("lsb_release", "-si")
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return err
		}

		id := strings.ToLower(strings.TrimSpace(stdout.String()))
		if id != "" && id != "n/a" {
			osInfo.Id = id
		}
	}

	if osInfo.FriendlyName == "" {
		cmd := exec.Command("lsb_release", "-sd")
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return err
		}

		friendlyName := strings.TrimSpace(stdout.String())
		if friendlyName != "" && friendlyName != "n/a" {
			osInfo.FriendlyName = friendlyName
		}
	}

	if osInfo.Version == "" {
		cmd := exec.Command("lsb_release", "-sr")
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return err
		}

		version := strings.TrimSpace(stdout.String())
		if version != "" && version != "n/a" {
			osInfo.Version = version
		}
	}

	if osInfo.Release == "" {
		cmd := exec.Command("lsb_release", "-sc")
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err := cmd.Run()
		if err != nil {
			return err
		}

		release := strings.TrimSpace(stdout.String())
		if release != "" && release != "n/a" {
			osInfo.Release = release
		}
	}

	return nil
}
