// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"maps"
	"strings"

	"github.com/trippsoft/forge/pkg/result"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

// HostInfo contains information about a managed host.
type HostInfo struct {
	runtime *RuntimeInfo
}

// Populate retrieves and populates the HostInfo using the provided transport.
func (i *HostInfo) Populate(t transport.Transport, runtimeOnly bool) *result.Result {
	var err error
	i.runtime.os, i.runtime.arch, err = t.GetOSAndArch()
	if err != nil {
		return result.NewFailure(err, err.Error())
	}

	if runtimeOnly {
		return result.NewSuccess(false, nil)
	}

	// TODO - Add more info population here in the future

	return result.NewSuccess(false, nil)
}

// Runtime returns the RuntimeInfo of the managed host.
func (i *HostInfo) Runtime() *RuntimeInfo {
	return i.runtime
}

// ToMapOfCtyValues converts the HostInfo into a map of cty.Values.
func (i *HostInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	maps.Copy(values, i.runtime.ToMapOfCtyValues())
	return values
}

// String returns a string representation of the host information.
//
// This is useful for logging or debugging purposes.
func (i *HostInfo) String() string {
	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString(i.runtime.String())
	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}

// NewHostInfo creates a new HostInfo instance.
func NewHostInfo() *HostInfo {
	return &HostInfo{
		runtime: &RuntimeInfo{},
	}
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

// ToMapOfCtyValues converts the RuntimeInfo into a map of cty.Values.
func (r *RuntimeInfo) ToMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)
	if r.os == "" {
		values["runtime_os"] = cty.NullVal(cty.String)
	} else {
		values["runtime_os"] = cty.StringVal(r.os)
	}

	if r.arch == "" {
		values["runtime_arch"] = cty.NullVal(cty.String)
	} else {
		values["runtime_arch"] = cty.StringVal(r.arch)
	}

	return values
}

// String returns a string representation of the OS information.
//
// This is useful for logging or debugging purposes.
func (r *RuntimeInfo) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("runtime_os: ")
	stringBuilder.WriteString(r.os)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("runtime_arch: ")
	stringBuilder.WriteString(r.arch)
	stringBuilder.WriteString("\n")

	return stringBuilder.String()
}

// NewRuntimeInfo creates a new RuntimeInfo instance.
func NewRuntimeInfo(os, arch string) *RuntimeInfo {
	return &RuntimeInfo{
		os:   os,
		arch: arch,
	}
}
