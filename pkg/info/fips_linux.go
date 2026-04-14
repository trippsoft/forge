// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package info

import (
	"bytes"
	"os"
)

func (f *FIPSInfo) discover() []string {
	f.Known = true

	fileInfo, err := os.Stat("/proc/sys/crypto/fips_enabled")
	if os.IsNotExist(err) {
		f.Enabled = false
		return nil
	}

	if err != nil {
		return []string{"failed to stat /proc/sys/crypto/fips_enabled: " + err.Error()}
	}

	if !fileInfo.Mode().IsRegular() {
		f.Enabled = false
		return []string{"/proc/sys/crypto/fips_enabled is not a regular file"}
	}

	data, err := os.ReadFile("/proc/sys/crypto/fips_enabled")
	if err != nil {
		return []string{"failed to read /proc/sys/crypto/fips_enabled: " + err.Error()}
	}

	content := string(bytes.TrimSpace(data))
	if content == "1" {
		f.Enabled = true
	} else {
		f.Enabled = false
	}

	return nil
}
