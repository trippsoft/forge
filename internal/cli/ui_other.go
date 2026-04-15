// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package cli

import (
	"os"
)

func InitUI(debug bool) {
	UI = &CLI{
		color:  true,
		stdout: os.Stdout,
		stderr: os.Stderr,
		debug:  debug,
	}
}
