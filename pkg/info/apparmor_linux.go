// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package info

import (
	"errors"
	"os"
)

func (a *AppArmorInfoPB) discover() error {
	a.Supported = true

	fileInfo, err := os.Stat("/sys/kernel/security/apparmor")
	if errors.Is(err, os.ErrNotExist) {
		a.Enabled = false
		return nil
	}

	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		a.Enabled = false
		return nil
	}

	a.Enabled = true
	return nil
}
