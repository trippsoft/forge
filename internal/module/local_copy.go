// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/plugin"
	pluginv1 "github.com/trippsoft/forge/pkg/plugin/v1"
	"github.com/trippsoft/forge/pkg/result"
	"github.com/zclconf/go-cty/cty"
)

var (
	localCopyInputSpec = hclspec.NewSpec(
		hclspec.Object(
			hclspec.RequiredField("source", hclspec.String).WithAliases("src"),
			hclspec.RequiredField("destination", hclspec.String).WithAliases("dest", "dst"),
		),
	)

	LocalCopy pluginv1.PluginModule = &LocalCopyModule{}
)

// LocalCopyModule is a module that copies a source file from the local filesystem to a destination file.
type LocalCopyModule struct{}

// Name implements pluginv1.PluginModule.
func (f *LocalCopyModule) Name() string {
	return "local_copy"
}

// Type implements pluginv1.PluginModule.
func (f *LocalCopyModule) Type() plugin.ModuleType {
	return plugin.ModuleType_REMOTE
}

// InputSpec implements pluginv1.PluginModule.
func (f *LocalCopyModule) InputSpec() *hclspec.Spec {
	return localCopyInputSpec
}

// RunModule implements pluginv1.PluginModule.
func (f *LocalCopyModule) RunModule(
	hostInfo *info.HostInfo,
	input map[string]cty.Value,
	whatIf bool,
) *result.ModuleResult {

	sourcePath := input["source"].AsString()
	sourceHash, err := hashFile(sourcePath)
	if err != nil {
		return pluginv1.NewFailure(fmt.Errorf("failed to hash source file from path %q: %w", sourcePath, err), "")
	}

	destinationPath := input["destination"].AsString()
	destinationHash, err := hashFile(destinationPath)
	if err == nil && bytes.Equal(sourceHash, destinationHash) {
		// Files are identical, no need to copy
		output := cty.ObjectVal(map[string]cty.Value{
			"sha256_hash": cty.StringVal(fmt.Sprintf("%x", sourceHash)),
		})

		success, err := pluginv1.NewNotChanged(output)
		if err != nil {
			return pluginv1.NewFailure(fmt.Errorf("failed to create module success result: %w", err), "")
		}

		return success
	}

	output := cty.ObjectVal(map[string]cty.Value{
		"sha256_hash": cty.StringVal(fmt.Sprintf("%x", sourceHash)),
	})

	success, err := pluginv1.NewChanged(output)
	if err != nil {
		return pluginv1.NewFailure(fmt.Errorf("failed to create module success result: %w", err), "")
	}

	if whatIf {
		// In what-if mode, do not perform the copy
		return success
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return pluginv1.NewFailure(fmt.Errorf("failed to open source file from path %q: %w", sourcePath, err), "")
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return pluginv1.NewFailure(
			fmt.Errorf("failed to create destination file at path %q: %w", destinationPath, err),
			"",
		)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return pluginv1.NewFailure(
			fmt.Errorf("failed to copy content from %q to %q: %w", sourcePath, destinationPath, err),
			"",
		)
	}

	return success
}

func hashFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
