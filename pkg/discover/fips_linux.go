// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package discover

import (
	"bytes"
	"os"
)

func discoverFIPSInfo() (*FIPSInfoResponse, error) {
	fipsInfo := &FIPSInfoResponse{
		Known: true,
	}

	fileInfo, err := os.Stat("/proc/sys/crypto/fips_enabled")
	if os.IsNotExist(err) {
		fipsInfo.Enabled = false
		return fipsInfo, nil
	}

	if err != nil {
		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		fipsInfo.Enabled = false
		return fipsInfo, nil
	}

	data, err := os.ReadFile("/proc/sys/crypto/fips_enabled")
	if err != nil {
		return nil, err
	}

	content := string(bytes.TrimSpace(data))
	if content == "1" {
		fipsInfo.Enabled = true
	} else {
		fipsInfo.Enabled = false
	}

	return fipsInfo, nil
}
