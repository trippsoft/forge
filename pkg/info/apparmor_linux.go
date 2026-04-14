// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package info

import (
	"errors"
	"os"
)

func (a *AppArmorInfo) discover() []string {
	a.Supported = true

	fileInfo, err := os.Stat("/sys/kernel/security/apparmor")
	if errors.Is(err, os.ErrNotExist) {
		a.Enabled = false
		return nil
	}

	if err != nil {
		return []string{"failed to stat /sys/kernel/security/apparmor: " + err.Error()}
	}

	if !fileInfo.IsDir() {
		a.Enabled = false
		return nil
	}

	a.Enabled = true

	return nil
}
