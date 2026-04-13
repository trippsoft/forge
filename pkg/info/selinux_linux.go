// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package info

import (
	"bufio"
	"os"
	"strings"
)

func (s *SELinuxInfo) discover() error {
	s.Supported = true

	fileInfo, err := os.Stat("/etc/selinux/config")
	if os.IsNotExist(err) {
		s.Installed = false
		return nil
	}

	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		s.Installed = false
		return nil
	}

	s.Installed = true
	file, err := os.Open("/etc/selinux/config")
	if err != nil {
		return err
	}
	defer file.Close()

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
		value := strings.Trim(strings.TrimSpace(stringParts[1]), `"`)

		switch key {
		case "SELINUX":
			s.Status = value
		case "SELINUXTYPE":
			s.Type = value
		}
	}

	return nil
}
