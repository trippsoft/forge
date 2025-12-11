// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !darwin && !windows

package info

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

func (o *OSInfoPB) discover() error {
	o.Kernel = runtime.GOOS
	o.Arch = runtime.GOARCH

	osReleaseErr := o.populateFromOsReleaseFile()
	lsbReleaseErr := o.populateFromLsbRelease()

	if osReleaseErr != nil && lsbReleaseErr != nil {
		return errors.Join(osReleaseErr, lsbReleaseErr)
	}

	if o.Release != "" {
		o.ReleaseId = strings.ReplaceAll(strings.ToLower(o.Release), " ", "-")
	}

	if o.Version != "" {
		versionParts := strings.SplitN(o.Version, ".", 2)
		if len(versionParts) >= 1 {
			o.MajorVersion = versionParts[0]
		}
	}

	correctedID, exists := osIDCorrectionMap[o.Id]
	if exists {
		o.Id = correctedID
	}

	families := []string{runtime.GOOS, "posix"}
	matchingFamilies, exists := osFamiliesMap[o.Id]
	if exists {
		families = append(families, matchingFamilies...)
	}

	if o.Id != "" {
		families = append(families, o.Id)
	}
	o.Families = families

	return nil
}

func (o *OSInfoPB) populateFromOsReleaseFile() error {
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
		o.Id = strings.ToLower(id)
	}

	friendlyName, ok := contents["PRETTY_NAME"]
	if ok && friendlyName != "" && friendlyName != "n/a" {
		o.FriendlyName = friendlyName
	}

	release, ok := contents["VERSION_CODENAME"]
	if ok && release != "" && release != "n/a" {
		o.Release = release
	}

	version, ok := contents["VERSION_ID"]
	if ok && version != "" && version != "n/a" {
		o.Version = version
	}

	edition, ok := contents["VARIANT"]
	if ok && edition != "" && edition != "n/a" {
		o.Edition = edition
	}

	editionId, ok := contents["VARIANT_ID"]
	if ok && editionId != "" && editionId != "n/a" {
		o.EditionId = editionId
	}

	return nil
}

func (o *OSInfoPB) populateFromLsbRelease() error {
	if o.Id == "" {
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
			o.Id = id
		}
	}

	if o.FriendlyName == "" {
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
			o.FriendlyName = friendlyName
		}
	}

	if o.Version == "" {
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
			o.Version = version
		}
	}

	if o.Release == "" {
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
			o.Release = release
		}
	}

	return nil
}
