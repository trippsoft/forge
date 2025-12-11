// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

import (
	"errors"
	"os"

	"golang.org/x/sys/windows/registry"
)

func (f *FIPSInfoPB) discover() error {
	f.Known = true

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Lsa\FipsAlgorithm\Enabled`,
		registry.QUERY_VALUE,
	)

	// If the key doesn't exist, FIPS is not enabled.
	if errors.Is(err, os.ErrNotExist) {
		f.Enabled = false
		return nil
	}

	if err != nil {
		return err
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("")
	if err != nil {
		return err
	}

	f.Enabled = val != 0

	return nil
}
