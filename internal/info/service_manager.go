package info

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/pkg/diag"
	"github.com/zclconf/go-cty/cty"
)

const (
	linuxServiceManagerDiscoveryScript = `systemctl_exists=0; ` +
		`run_systemd_system_exists=0; ` +
		`dev_run_systemd_exists=0; ` +
		`dev_systemd_exists=0; ` +
		`initctl_exists=0; ` +
		`etc_init_exists=0; ` +
		`openrc_exists=0; ` +
		`init_link_target=""; ` +
		`etc_init_d_exists=0; ` +
		`proc1_comm=""; ` +
		`systemctl_path=$(which systemctl 2>/dev/null || echo ""); ` +
		`if [ -z "$systemctl_path" ]; ` +
		`then systemctl_exists=0; ` +
		`elif [ -x "$systemctl_path" ]; ` +
		`then systemctl_exists=1; ` +
		`fi; ` +
		`if [ -d /run/systemd/system ]; ` +
		`then run_systemd_system_exists=1; ` +
		`fi; ` +
		`if [ -d /dev/.run/systemd ]; ` +
		`then dev_run_systemd_exists=1; ` +
		`fi; ` +
		`if [ -d /dev/.systemd ]; ` +
		`then dev_systemd_exists=1; ` +
		`fi; ` +
		`initctl_path=$(which initctl 2>/dev/null || echo ""); ` +
		`if [ -z "$initctl_path" ]; ` +
		`then initctl_exists=0; ` +
		`elif [ -f "$initctl_path" ]; ` +
		`then initctl_exists=1; ` +
		`fi; ` +
		`if [ -d /etc/init ]; ` +
		`then etc_init_exists=1; ` +
		`fi; ` +
		`if [ -f /sbin/openrc ]; ` +
		`then openrc_exists=1; ` +
		`fi; ` +
		`if [ -L /sbin/init ]; ` +
		`then init_link_target=$(readlink /sbin/init); ` +
		`fi; ` +
		`if [ -d /etc/init.d ]; ` +
		`then etc_init_d_exists=1; ` +
		`fi; ` +
		`if [ -f /proc/1/comm ]; ` +
		`then proc1_comm=$(cat /proc/1/comm); ` +
		`fi; ` +
		`output=$(jq -n ` +
		`--arg systemctl_exists "$systemctl_exists" ` +
		`--arg run_systemd_system_exists "$run_systemd_system_exists" ` +
		`--arg dev_run_systemd_exists "$dev_run_systemd_exists" ` +
		`--arg dev_systemd_exists "$dev_systemd_exists" ` +
		`--arg initctl_exists "$initctl_exists" ` +
		`--arg etc_init_exists "$etc_init_exists" ` +
		`--arg openrc_exists "$openrc_exists" ` +
		`--arg init_link_target "$init_link_target" ` +
		`--arg etc_init_d_exists "$etc_init_d_exists" ` +
		`--arg proc1_comm "$proc1_comm" ` +
		`'{` +
		`systemctl_exists: $systemctl_exists, ` +
		`run_systemd_system_exists: $run_systemd_system_exists, ` +
		`dev_run_systemd_exists: $dev_run_systemd_exists, ` +
		`dev_systemd_exists: $dev_systemd_exists, ` +
		`initctl_exists: $initctl_exists, ` +
		`etc_init_exists: $etc_init_exists, ` +
		`openrc_exists: $openrc_exists, ` +
		`init_link_target: $init_link_target, ` +
		`etc_init_d_exists: $etc_init_d_exists, ` +
		`proc1_comm: $proc1_comm}'); ` +
		`echo "$output"`
)

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

func (s *ServiceManagerInfo) populateServiceManagerInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.ID() == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping service manager information collection due to missing or invalid OS info",
		}}
	}

	if osInfo.Families().Contains("darwin") {
		return s.populateDarwinServiceManagerInfo(osInfo)
	}

	if osInfo.Families().Contains("windows") {
		s.name = "windows-service-manager"
		return diag.Diags{}
	}

	if osInfo.Families().Contains("linux") {
		return s.populateLinuxServiceManagerInfo(transport)
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagError,
		Summary:  "Unsupported OS family",
		Detail:   "Service manager information collection is not supported for this OS family",
	}}
}

func (s *ServiceManagerInfo) populateDarwinServiceManagerInfo(osInfo *OSInfo) diag.Diags {

	s.name = "launchd" // Default to launchd for macOS

	majorVersion, err := strconv.Atoi(osInfo.MajorVersion())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Invalid OS major version",
			Detail:   fmt.Sprintf("Error parsing OS major version: %v", err),
		}}
	}

	if majorVersion < 10 {
		s.name = "systemstarter"
		return diag.Diags{}
	}

	if majorVersion > 10 {
		s.name = "launchd"
		return diag.Diags{}
	}

	versionParts := strings.Split(osInfo.Version(), ".")
	if len(versionParts) < 2 {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Invalid OS version format",
			Detail:   "OS version does not contain enough parts to determine service manager",
		}}
	}

	minorVersion, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Invalid OS minor version",
			Detail:   fmt.Sprintf("Error parsing OS minor version: %v", err),
		}}
	}

	if minorVersion < 4 {
		s.name = "systemstarter"
	} else {
		s.name = "launchd"
	}

	return diag.Diags{}
}

func (s *ServiceManagerInfo) populateLinuxServiceManagerInfo(transport transport.Transport) diag.Diags {

	cmd := transport.NewCommand(linuxServiceManagerDiscoveryScript)
	var outBuf bytes.Buffer
	cmd.SetStdout(&outBuf)

	err := cmd.Run(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to execute service manager discovery script",
			Detail:   fmt.Sprintf("Error executing service manager discovery script: %v", err),
		}}
	}

	stdout := strings.TrimSpace(outBuf.String())

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
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
			return diag.Diags{}
		}

		devRunSystemdExists, _ := discoveredData["dev_run_systemd_exists"]

		if devRunSystemdExists == "1" {
			s.name = "systemd"
			return diag.Diags{}
		}

		devSystemdExists, _ := discoveredData["dev_systemd_exists"]

		if devSystemdExists == "1" {
			s.name = "systemd"
			return diag.Diags{}
		}
	}

	initctlExists, _ := discoveredData["initctl_exists"]

	if initctlExists == "1" {
		etcInitExists, _ := discoveredData["etc_init_exists"]

		if etcInitExists == "1" {
			s.name = "upstart"
			return diag.Diags{}
		}
	}

	openrcExists, _ := discoveredData["openrc_exists"]

	if openrcExists == "1" {
		s.name = "openrc"
		return diag.Diags{}
	}

	initLinkTarget, _ := discoveredData["init_link_target"]

	initLinkTargetParts := strings.Split(initLinkTarget, "/")
	initLinkTarget = initLinkTargetParts[len(initLinkTargetParts)-1] // Get the last part of the path
	if initLinkTarget == "systemd" {
		s.name = "systemd"
		return diag.Diags{}
	}

	etcInitDExists, _ := discoveredData["etc_init_d_exists"]

	if etcInitDExists == "1" {
		s.name = "sysvinit"
		return diag.Diags{}
	}

	proc1Comm, _ := discoveredData["proc1_comm"]

	if proc1Comm != "" && proc1Comm != "COMMAND" && proc1Comm != "init" && !strings.HasSuffix(proc1Comm, "sh") {
		if serviceManager, ok := proc1CommMap[proc1Comm]; ok {
			s.name = serviceManager
			return diag.Diags{}
		}

		s.name = proc1Comm
		return diag.Diags{}
	}

	if initLinkTarget != "" && initLinkTarget != "init" && !strings.HasSuffix(initLinkTarget, "sh") {
		if serviceManager, ok := proc1CommMap[initLinkTarget]; ok {
			s.name = serviceManager
			return diag.Diags{}
		}

		s.name = initLinkTarget
		return diag.Diags{}
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagError,
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
