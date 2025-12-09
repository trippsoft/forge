// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package plugin

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

func init() {
	SharedPluginBasePath = "/usr/share/forge/plugins"
	home, _ := homedir.Dir()
	UserPluginBasePath = home + `/.local/share/forge/plugins`
	os.MkdirAll(UserPluginBasePath, 0777)
}

func FindPluginPath(namespace, pluginName, osName, arch string) (string, error) {
	var extension string
	if osName == "windows" {
		extension = ".exe"
	}

	pluginPathSuffix := fmt.Sprintf("/%s/%s/%s-%s_%s_%s%s",
		namespace,
		pluginName,
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
