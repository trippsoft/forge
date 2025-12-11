// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pluginv1

import "github.com/trippsoft/forge/pkg/hclspec"

type PluginModule interface {
	Name() string
	InputSpec() *hclspec.Spec
}
