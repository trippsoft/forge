package info

import (
	"context"
	"errors"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	UserHomeDirPowerShell = `Write-Host $env:USERPROFILE`
	UserNamePowerShell    = `Write-Host $env:USERNAME`
	UserIdPowerShell      = `$obj = [Security.Principal.WindowsIdentity]::GetCurrent(); Write-Host $obj.User.Value`
)

type userInfo struct {
	name    string
	userId  string // UID on POSIX systems, SID on Windows
	groupId string // GID on POSIX systems, not applicable on Windows
	homeDir string
	shell   string // Login shell on POSIX systems, not applicable on Windows
	gecos   string // User information on POSIX systems, not applicable on Windows
}

func newUserInfo() *userInfo {
	return &userInfo{}
}

func (u *userInfo) Name() string {
	return u.name
}

func (u *userInfo) UserId() string {
	return u.userId
}

func (u *userInfo) GroupId() string {
	return u.groupId
}

func (u *userInfo) HomeDir() string {
	return u.homeDir
}

func (u *userInfo) Shell() string {
	return u.shell
}

func (u *userInfo) Gecos() string {
	return u.gecos
}

func (u *userInfo) populateUserInfo(osInfo *osInfo, transport transport.Transport) error {

	if osInfo.families.Contains("posix") {
		return u.populatePosixUserInfo(transport)
	}

	if osInfo.families.Contains("windows") {
		return u.populateWindowsUserInfo(transport)
	}

	return errors.New("unsupported OS family for user info population")
}

func (u *userInfo) populatePosixUserInfo(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "id -nu")
	if err != nil {
		return err
	}
	u.name = strings.TrimSpace(stdout)

	stdout, _, err = transport.ExecuteCommand(context.Background(), "id -u")
	if err != nil {
		return err
	}
	u.userId = strings.TrimSpace(stdout)

	stdout, _, err = transport.ExecuteCommand(context.Background(), "id -g")
	if err != nil {
		return err
	}
	u.groupId = strings.TrimSpace(stdout)

	stdout, _, err = transport.ExecuteCommand(context.Background(), "echo $HOME")
	if err != nil {
		return err
	}
	u.homeDir = strings.TrimSpace(stdout)

	stdout, _, err = transport.ExecuteCommand(context.Background(), "echo $SHELL")
	if err != nil {
		return err
	}
	u.shell = strings.TrimSpace(stdout)

	stdout, _, err = transport.ExecuteCommand(context.Background(), "getent passwd "+u.name+" | cut -d ':' -f 5")
	if err != nil {
		return err
	}
	u.gecos = strings.TrimSpace(stdout)

	return nil
}

func (u *userInfo) populateWindowsUserInfo(transport transport.Transport) error {

	stdout, err := transport.ExecutePowerShell(context.Background(), UserNamePowerShell)
	if err != nil {
		return err
	}

	u.name = strings.TrimSpace(stdout)

	stdout, err = transport.ExecutePowerShell(context.Background(), UserIdPowerShell)
	if err != nil {
		return err
	}
	u.userId = strings.TrimSpace(stdout)

	stdout, err = transport.ExecutePowerShell(context.Background(), UserHomeDirPowerShell)
	if err != nil {
		return err
	}
	u.homeDir = strings.TrimSpace(stdout)

	return nil
}

func (u *userInfo) toMapOfCtyValues() map[string]cty.Value {
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

func (u *userInfo) String() string {
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
