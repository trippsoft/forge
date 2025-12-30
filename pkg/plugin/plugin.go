// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"os"
	"strconv"
)

var (
	SharedPluginBasePath string
	UserPluginBasePath   string

	LocalPluginMinPort uint16 = 25000
	LocalPluginMaxPort uint16 = 40000

	DefaultRemotePluginMinPort uint16 = 25000
	DefaultRemotePluginMaxPort uint16 = 40000
)

func GetMinimumPort() uint16 {
	env := os.Getenv("FORGE_PLUGIN_MIN_PORT")
	if env != "" {
		minPort, err := strconv.ParseUint(env, 10, 16)
		if err == nil {
			return uint16(minPort)
		}
	}

	return 25000
}

func GetMaximumPort() uint16 {
	env := os.Getenv("FORGE_PLUGIN_MAX_PORT")
	if env != "" {
		maxPort, err := strconv.ParseUint(env, 10, 16)
		if err == nil {
			return uint16(maxPort)
		}
	}

	return 40000
}
