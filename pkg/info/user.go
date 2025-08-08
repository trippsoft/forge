package info

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/pkg/diag"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run ../../cmd/scriptimport/main.go info user_posix_discovery.sh
//go:generate go run ../../cmd/scriptimport/main.go info user_windows_discovery.ps1

type UserInfo struct {
	name    string
	userId  string // UID on POSIX systems, SID on Windows
	groupId string // GID on POSIX systems, not applicable on Windows
	homeDir string
	shell   string // Login shell on POSIX systems, not applicable on Windows
	gecos   string // User information on POSIX systems, not applicable on Windows
}

func newUserInfo() *UserInfo {
	return &UserInfo{}
}

func (u *UserInfo) Name() string {
	return u.name
}

func (u *UserInfo) UserId() string {
	return u.userId
}

func (u *UserInfo) GroupId() string {
	return u.groupId
}

func (u *UserInfo) HomeDir() string {
	return u.homeDir
}

func (u *UserInfo) Shell() string {
	return u.shell
}

func (u *UserInfo) Gecos() string {
	return u.gecos
}

func (u *UserInfo) populateUserInfo(osInfo *OSInfo, transport transport.Transport) diag.Diags {

	if osInfo == nil || osInfo.id == "" {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagWarning,
			Summary:  "Invalid OS information",
			Detail:   "Skipping user information collection due to missing or invalid OS info",
		}}
	}

	if osInfo.families.Contains("posix") {
		return u.populatePosixUserInfo(transport)
	}

	if osInfo.families.Contains("windows") {
		return u.populateWindowsUserInfo(transport)
	}

	return diag.Diags{&diag.Diag{
		Severity: diag.DiagError,
		Summary:  "Unsupported OS family",
		Detail:   "User information collection is not supported on this OS",
	}}
}

func (u *UserInfo) populatePosixUserInfo(t transport.Transport) diag.Diags {

	cmd, err := t.NewCommand(userPosixDiscoveryScript, nil)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on POSIX host: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on POSIX host: %v", err),
		}, &diag.Diag{
			Severity: diag.DiagDebug,
			Summary:  "Discovery command stderr",
			Detail:   fmt.Sprintf("stderr: %s", stderr),
		}}
	}

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse user information",
			Detail:   fmt.Sprintf("Error parsing user information on POSIX host: %v", err),
		}}
	}

	u.name, _ = discoveredData["user_name"]
	u.userId, _ = discoveredData["user_id"]
	u.groupId, _ = discoveredData["user_gid"]
	u.homeDir, _ = discoveredData["user_home_dir"]
	u.shell, _ = discoveredData["user_shell"]
	u.gecos, _ = discoveredData["user_gecos"]

	return diag.Diags{}
}

func (u *UserInfo) populateWindowsUserInfo(t transport.Transport) diag.Diags {

	cmd, err := t.NewPowerShellCommand(userWindowsDiscoveryScript, nil)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on Windows host: %v", err),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on Windows host: %v", err),
		}, &diag.Diag{
			Severity: diag.DiagDebug,
			Summary:  "Discovery command stderr",
			Detail:   fmt.Sprintf("stderr: %s", stderr),
		}}
	}

	discoveredData := make(map[string]string)
	err = json.Unmarshal([]byte(stdout), &discoveredData)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to parse user information",
			Detail:   fmt.Sprintf("Error parsing user information on Windows host: %v", err),
		}}
	}

	u.name, _ = discoveredData["user_name"]
	u.userId, _ = discoveredData["user_id"]
	u.homeDir, _ = discoveredData["user_home_dir"]

	return diag.Diags{}
}

func (u *UserInfo) toMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	if u.name != "" {
		values["user_name"] = cty.StringVal(u.name)
	} else {
		values["user_name"] = cty.NullVal(cty.String)
	}

	if u.userId != "" {
		values["user_id"] = cty.StringVal(u.userId)
	} else {
		values["user_id"] = cty.NullVal(cty.String)
	}

	if u.groupId != "" {
		values["user_gid"] = cty.StringVal(u.groupId)
	} else {
		values["user_gid"] = cty.NullVal(cty.String)
	}

	if u.homeDir != "" {
		values["user_home_dir"] = cty.StringVal(u.homeDir)
	} else {
		values["user_home_dir"] = cty.NullVal(cty.String)
	}

	if u.shell != "" {
		values["user_shell"] = cty.StringVal(u.shell)
	} else {
		values["user_shell"] = cty.NullVal(cty.String)
	}

	if u.gecos != "" {
		values["user_gecos"] = cty.StringVal(u.gecos)
	} else {
		values["user_gecos"] = cty.NullVal(cty.String)
	}

	return values
}

func (u *UserInfo) String() string {
	stringBuilder := &strings.Builder{}

	stringBuilder.WriteString("user_name: ")
	stringBuilder.WriteString(u.name)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("user_id: ")
	stringBuilder.WriteString(u.userId)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("user_gid: ")
	stringBuilder.WriteString(u.groupId)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("user_home_dir: ")
	stringBuilder.WriteString(u.homeDir)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("user_shell: ")
	stringBuilder.WriteString(u.shell)
	stringBuilder.WriteString("\n")

	stringBuilder.WriteString("user_gecos: ")
	stringBuilder.WriteString(u.gecos)

	return stringBuilder.String()
}
