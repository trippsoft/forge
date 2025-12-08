// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"fmt"
	"os"
)

var (
	SharedPluginBasePath string
	UserPluginBasePath   string

	LocalPluginMinPort uint16 = 25000
	LocalPluginMaxPort uint16 = 40000

	DefaultRemotePluginMinPort uint16 = 25000
	DefaultRemotePluginMaxPort uint16 = 40000
)

func FindPluginPath(namespace, pluginName, osName, arch string) (string, error) {
	var pathSeparator string
	var extension string
	if osName == "windows" {
		pathSeparator = `\`
		extension = ".exe"
	} else {
		pathSeparator = "/"
		extension = ""
	}

	pluginPathSuffix := fmt.Sprintf("%s%s%s%s%s%s-%s_%s_%s%s",
		pathSeparator,
		namespace,
		pathSeparator,
		pluginName,
		pathSeparator,
		namespace,
		pluginName,
		osName,
		arch,
		extension,
	)

	userPluginPath := UserPluginBasePath + pluginPathSuffix
	fileInfo, err := os.Stat(userPluginPath)
	if err == nil && !fileInfo.IsDir() {
		return userPluginPath, nil
	}

	sharedPluginPath := SharedPluginBasePath + pluginPathSuffix
	fileInfo, err = os.Stat(sharedPluginPath)
	if err == nil && !fileInfo.IsDir() {
		return sharedPluginPath, nil
	}

	return "", fmt.Errorf(
		`plugin "%s/%s" does not exist for OS %q and architecture %q`,
		namespace,
		pluginName,
		osName,
		arch,
	)
}
