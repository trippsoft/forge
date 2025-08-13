// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run ../../cmd/scriptimport/main.go info service_manager_linux_discovery.sh

var (
	proc1CommMap = map[string]string{
		"openrc-init": "openrc",
	}
)

type ServiceManagerInfo struct {
	name string
}

func newServiceManagerInfo() *ServiceManagerInfo {
	return &ServiceManagerInfo{}
}

func (s *ServiceManagerInfo) Name() string {
	return s.name
}

func (s *ServiceManagerInfo) populateServiceManagerInfo(osInfo *OSInfo, transport transport.Transport) util.Diags {

	if osInfo == nil || osInfo.ID() == "" {
		return util.Diags{&util.Diag{
			Severity: util.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping service manager information collection due to missing or invalid OS info",
		}}
	}

	if osInfo.Families().Contains("darwin") {
		return s.populateDarwinServiceManagerInfo(osInfo)
	}

	if osInfo.Families().Contains("windows") {
		s.name = "windows-service-manager"
		return util.Diags{}
	}

	if osInfo.Families().Contains("linux") {
		return s.populateLinuxServiceManagerInfo(transport)
	}

	return util.Diags{&util.Diag{
		Severity: util.DiagError,
		Summary:  "Unsupported OS family",
		Detail:   "Service manager information collection is not supported for this OS family",
	}}
}

func (s *ServiceManagerInfo) populateDarwinServiceManagerInfo(osInfo *OSInfo) util.Diags {

	s.name = "launchd" // Default to launchd for macOS

	majorVersion, err := strconv.Atoi(osInfo.MajorVersion())
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Invalid OS major version",
			Detail:   fmt.Sprintf("Error parsing OS major version: %v", err),
		}}
	}

	if majorVersion < 10 {
		s.name = "systemstarter"
		return util.Diags{}
	}

	if majorVersion > 10 {
		s.name = "launchd"
		return util.Diags{}
	}

	versionParts := strings.Split(osInfo.Version(), ".")
	if len(versionParts) < 2 {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Invalid OS version format",
			Detail:   "OS version does not contain enough parts to determine service manager",
		}}
	}

	minorVersion, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Invalid OS minor version",
			Detail:   fmt.Sprintf("Error parsing OS minor version: %v", err),
		}}
	}

	if minorVersion < 4 {
		s.name = "systemstarter"
	} else {
		s.name = "launchd"
	}

	return util.Diags{}
}

func (s *ServiceManagerInfo) populateLinuxServiceManagerInfo(t transport.Transport) util.Diags {

	cmd, err := t.NewCommand(serviceManagerLinuxDiscoveryScript, nil)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to create service manager discovery command",
			Detail:   fmt.Sprintf("Error creating command: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to execute service manager discovery script",
			Detail:   fmt.Sprintf("Error executing service manager discovery script: %v", err),
		}, &util.Diag{
			Severity: util.DiagDebug,
			Summary:  "Discovery command stderr",
			Detail:   fmt.Sprintf("stderr: %s", stderr),
		}}
	}

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to parse service manager discovery script output",
			Detail:   fmt.Sprintf("Error parsing service manager discovery script output: %v", err),
		}}
	}

	systemctlExists, _ := discoveredData["systemctl_exists"]

	hasSystemctl := systemctlExists == "1"

	if hasSystemctl {

		runSystemdSystemExists, _ := discoveredData["run_systemd_system_exists"]

		if runSystemdSystemExists == "1" {
			s.name = "systemd"
			return util.Diags{}
		}

		devRunSystemdExists, _ := discoveredData["dev_run_systemd_exists"]

		if devRunSystemdExists == "1" {
			s.name = "systemd"
			return util.Diags{}
		}

		devSystemdExists, _ := discoveredData["dev_systemd_exists"]

		if devSystemdExists == "1" {
			s.name = "systemd"
			return util.Diags{}
		}
	}

	initctlExists, _ := discoveredData["initctl_exists"]

	if initctlExists == "1" {
		etcInitExists, _ := discoveredData["etc_init_exists"]

		if etcInitExists == "1" {
			s.name = "upstart"
			return util.Diags{}
		}
	}

	openrcExists, _ := discoveredData["openrc_exists"]

	if openrcExists == "1" {
		s.name = "openrc"
		return util.Diags{}
	}

	initLinkTarget, _ := discoveredData["init_link_target"]

	initLinkTargetParts := strings.Split(initLinkTarget, "/")
	initLinkTarget = initLinkTargetParts[len(initLinkTargetParts)-1] // Get the last part of the path
	if initLinkTarget == "systemd" {
		s.name = "systemd"
		return util.Diags{}
	}

	etcInitDExists, _ := discoveredData["etc_init_d_exists"]

	if etcInitDExists == "1" {
		s.name = "sysvinit"
		return util.Diags{}
	}

	proc1Comm, _ := discoveredData["proc1_comm"]

	if proc1Comm != "" && proc1Comm != "COMMAND" && proc1Comm != "init" && !strings.HasSuffix(proc1Comm, "sh") {
		if serviceManager, ok := proc1CommMap[proc1Comm]; ok {
			s.name = serviceManager
			return util.Diags{}
		}

		s.name = proc1Comm
		return util.Diags{}
	}

	if initLinkTarget != "" && initLinkTarget != "init" && !strings.HasSuffix(initLinkTarget, "sh") {
		if serviceManager, ok := proc1CommMap[initLinkTarget]; ok {
			s.name = serviceManager
			return util.Diags{}
		}

		s.name = initLinkTarget
		return util.Diags{}
	}

	return util.Diags{&util.Diag{
		Severity: util.DiagError,
		Summary:  "Failed to determine service manager",
		Detail:   "Could not identify the service manager for the current Linux system",
	}}
}

func (s *ServiceManagerInfo) toMapOfCtyValues() map[string]cty.Value {

	if s.name == "" {
		return map[string]cty.Value{
			"service_manager": cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"service_manager": cty.StringVal(s.name),
	}
}

// String returns a string representation of the service manager information.
// This is useful for logging or debugging purposes.
func (s *ServiceManagerInfo) String() string {

	if s.name == "" {
		return "service_manager: unknown"
	}

	return "service_manager: " + s.name
}
