package info

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/trippsoft/forge/pkg/diag"
	"github.com/zclconf/go-cty/cty"
)

const (
	userPosixDiscoveryScript = `user_name=$(id -nu); ` +
		`user_id=$(id -u); ` +
		`user_gid=$(id -g); ` +
		`user_home_dir="$HOME"; ` +
		`user_shell="$SHELL"; ` +
		`user_gecos=$(getent passwd $user_name | cut -d ':' -f 5); ` +
		`output=$(jq -n ` +
		`--arg user_name "$user_name" ` +
		`--arg user_id "$user_id" ` +
		`--arg user_gid "$user_gid" ` +
		`--arg user_home_dir "$user_home_dir" ` +
		`--arg user_shell "$user_shell" ` +
		`--arg user_gecos "$user_gecos" ` +
		`'{user_name: $user_name, user_id: $user_id, user_gid: $user_gid, user_home_dir: $user_home_dir, user_shell: $user_shell, user_gecos: $user_gecos}'); ` +
		`echo "$output"`
	userWindowsDiscoveryScript = `$userName = $env:USERNAME; ` +
		`$userId = $userId = [Security.Principal.WindowsIdentity]::GetCurrent().User.Value; ` +
		`$userHomeDir = $env:USERPROFILE; ` +
		`$output = @{user_name = $userName; user_id = $userId; user_home_dir = $userHomeDir}; ` +
		`$json = $output | ConvertTo-Json -Depth 3; ` +
		`Write-Host $json`
	UserHomeDirPowerShell = `Write-Host $env:USERPROFILE`
	UserNamePowerShell    = `Write-Host $env:USERNAME`
	UserIdPowerShell      = `$obj = [Security.Principal.WindowsIdentity]::GetCurrent(); Write-Host $obj.User.Value`
)

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

func (u *UserInfo) populatePosixUserInfo(transport transport.Transport) diag.Diags {

	cmd := transport.NewCommand(userPosixDiscoveryScript)

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on POSIX host: %v", err),
		}}
	}

	stdout := strings.TrimSpace(string(stdoutBytes))

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

func (u *UserInfo) populateWindowsUserInfo(transport transport.Transport) diag.Diags {

	cmd, err := transport.NewPowerShellCommand(userWindowsDiscoveryScript)
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on Windows host: %v", err),
		}}
	}

	stdoutBytes, err := cmd.Output(context.Background())
	if err != nil {
		return diag.Diags{&diag.Diag{
			Severity: diag.DiagError,
			Summary:  "Failed to get user information",
			Detail:   fmt.Sprintf("Error getting user information on Windows host: %v", err),
		}}
	}

	stdout := strings.TrimSpace(string(stdoutBytes))

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
