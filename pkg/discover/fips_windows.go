// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package discover

import "golang.org/x/sys/windows/registry"

func discoverFIPSInfo() (*FIPSInfoResponse, error) {
	fipsInfo := &FIPSInfoResponse{
		Known: true,
	}

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Lsa\FipsAlgorithm\Enabled`,
		registry.QUERY_VALUE,
	)

	// If the key doesn't exist, FIPS is not enabled.
	if registry.ErrNotExist.Is(err) {
		fipsInfo.Enabled = false
		return fipsInfo, nil
	}

	if err != nil {
		return nil, err
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("")
	if err != nil {
		return nil, err
	}

	fipsInfo.Enabled = val != 0

	return fipsInfo, nil
}
