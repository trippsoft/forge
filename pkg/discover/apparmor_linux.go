// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package discover

import (
	"errors"
	"os"
)

func discoverAppArmorInfo() (*AppArmorInfoResponse, error) {
	appArmorInfo := &AppArmorInfoResponse{
		Supported: true,
	}

	fileInfo, err := os.Stat("/sys/kernel/security/apparmor")
	if errors.Is(err, os.ErrNotExist) {
		appArmorInfo.Enabled = false
		return appArmorInfo, nil
	}

	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		appArmorInfo.Enabled = false
		return appArmorInfo, nil
	}

	appArmorInfo.Enabled = true
	return appArmorInfo, nil
}
