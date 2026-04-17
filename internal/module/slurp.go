// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
)

var (
	slurpInputSpec = hclspec.NewSpec(hclspec.Object(
		hclspec.RequiredField("path", hclspec.String).WithAliases("source", "src"),
	))

	Slurp pluginv1.PluginModule = &SlurpModule{}
)

// SlurpModule is a module that gets files contents from a given path and returns it in a base64 encoded string.
type SlurpModule struct{}

// Name implements [pluginv1.PluginModule].
func (f *SlurpModule) Name() string {
	return "slurp"
}

// Type implements [pluginv1.PluginModule].
func (f *SlurpModule) Type() plugin.ModuleType {
	return plugin.ModuleType_REMOTE
}

// InputSpec implements [pluginv1.PluginModule].
func (f *SlurpModule) InputSpec() *hclspec.Spec {
	return slurpInputSpec
}

// RunModule implements [pluginv1.PluginModule].
func (f *SlurpModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *result.ModuleResult {
	path := input["path"].AsString()

	content, err := os.ReadFile(path)
	if err != nil {
		return pluginv1.NewFailure(fmt.Errorf("failed to read file from path %q: %w", path, err), "")
	}

	stringBuilder := &strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, stringBuilder)
	_, err = encoder.Write(content)
	if err != nil {
		return pluginv1.NewFailure(
			fmt.Errorf("failed to encode file content from path %q to base64: %w", path, err),
			"",
		)
	}

	err = encoder.Close()
	if err != nil {
		return pluginv1.NewFailure(
			fmt.Errorf("failed to finalize base64 encoding for file content from path %q: %w", path, err),
			"",
		)
	}

	output := cty.ObjectVal(map[string]cty.Value{
		"content":     cty.StringVal(stringBuilder.String()),
		"sha256_hash": cty.StringVal(fmt.Sprintf("%x", sha256.Sum256(content))),
	})

	success, err := pluginv1.NewNotChanged(output)
	if err != nil {
		return pluginv1.NewFailure(
			fmt.Errorf("failed to create module success result: %w", err),
			"",
		)
	}

	return success
}
