// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package plugin

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

func init() {
	SharedPluginBasePath = `C:\ProgramData\Forge\plugins`
	home, _ := homedir.Dir()
	UserPluginBasePath = home + `\AppData\Local\Forge\plugins`
	os.MkdirAll(UserPluginBasePath, 0777)
}

func FindPluginPath(basePath, namespace, pluginName, osName, arch string) (string, error) {
	var extension string
	if osName == "windows" {
		extension = ".exe"
	}

	pluginPath := fmt.Sprintf(`%s\%s\%s\%s-%s_%s_%s%s`,
		basePath,
		namespace,
		pluginName,
		namespace,
		pluginName,
		osName,
		arch,
		extension,
	)

	fileInfo, err := os.Stat(pluginPath)
	if err == nil && !fileInfo.IsDir() {
		return pluginPath, nil
	}

	return "", fmt.Errorf(
		`plugin "%s/%s" does not exist for OS %q and architecture %q`,
		namespace,
		pluginName,
		osName,
		arch,
	)
}
