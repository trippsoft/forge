// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package discover

func DefaultDiscoverPluginBasePath() string {
	return "/usr/share/forge/"
}
