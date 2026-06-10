// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:generate protoc --proto_path=../.. --go_out=../.. --go_opt=paths=source_relative pkg/plugin/plugin.proto

package plugin
