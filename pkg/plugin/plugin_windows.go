// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package plugin

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

func init() {
	SharedPluginBasePath = `C:\Program Files\Forge\plugins`
	home, _ := homedir.Dir()
	UserPluginBasePath = home + `\AppData\Local\Forge\plugins`
	os.MkdirAll(UserPluginBasePath, 0777)
}
