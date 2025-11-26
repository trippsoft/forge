// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

// HostInfo contains information about a managed host.
type HostInfo struct {
	runtime *RuntimeInfo
}

// RuntimeInfo contains the OS and architecture information of a managed host.
type RuntimeInfo struct {
	os   string
	arch string
}

// OS returns the operating system of the managed host.
func (r *RuntimeInfo) OS() string {
	return r.os
}

// Arch returns the architecture of the managed host.
func (r *RuntimeInfo) Arch() string {
	return r.arch
}

// NewRuntimeInfo creates a new RuntimeInfo instance.
func NewRuntimeInfo(os, arch string) *RuntimeInfo {
	return &RuntimeInfo{
		os:   os,
		arch: arch,
	}
}
