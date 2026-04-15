// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package cli

import (
	"os"

	"golang.org/x/sys/windows"
)

func InitUI(debug bool) {
	var outMode uint32
	out := windows.Handle(os.Stdout.Fd())
	err := windows.GetConsoleMode(out, &outMode)
	if err != nil {
		UI = &CLI{
			color:  false,
			stdout: os.Stdout,
			stderr: os.Stderr,
			debug:  debug,
		}
		return
	}

	outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(out, outMode)
	if err != nil {
		UI = &CLI{
			color:  false,
			stdout: os.Stdout,
			stderr: os.Stderr,
			debug:  debug,
		}
		return
	}

	UI = &CLI{
		color:  true,
		stdout: os.Stdout,
		stderr: os.Stderr,
		debug:  debug,
	}
}
