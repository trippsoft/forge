// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package plugin

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

func init() {
	SharedPluginBasePath = "/usr/share/forge/plugins"
	home, _ := homedir.Dir()
	UserPluginBasePath = home + `/.local/share/forge/plugins`
	os.MkdirAll(UserPluginBasePath, 0777)
}
