package info

import (
	"testing"
)

func TestUserInfo_PopulateUserInfo_Posix(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("posix")

	info := newUserInfo()

	transport := newMockTransport()
	transport.commandResponses["id -nu"] = &commandResponse{
		stdout: "testuser\n",
	}
	transport.commandResponses["id -u"] = &commandResponse{
		stdout: "1000\n",
	}
	transport.commandResponses["id -g"] = &commandResponse{
		stdout: "1000\n",
	}
	transport.commandResponses["echo $HOME"] = &commandResponse{
		stdout: "/home/testuser\n",
	}
	transport.commandResponses["echo $SHELL"] = &commandResponse{
		stdout: "/bin/bash\n",
	}
	transport.commandResponses["getent passwd testuser | cut -d ':' -f 5"] = &commandResponse{
		stdout: "Test User\n",
	}

	err := info.populateUserInfo(osInfo, transport)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if info.Name() != "testuser" {
		t.Errorf("expected name 'testuser', got '%s'", info.Name())
	}
	if info.UserId() != "1000" {
		t.Errorf("expected user ID '1000', got '%s'", info.UserId())
	}
	if info.GroupId() != "1000" {
		t.Errorf("expected group ID '1000', got '%s'", info.GroupId())
	}
	if info.HomeDir() != "/home/testuser" {
		t.Errorf("expected home directory '/home/testuser', got '%s'", info.HomeDir())
	}
	if info.Shell() != "/bin/bash" {
		t.Errorf("expected shell '/bin/bash', got '%s'", info.Shell())
	}
	if info.Gecos() != "Test User" {
		t.Errorf("expected GECOS 'Test User', got '%s'", info.Gecos())
	}
}

func TestUserInfo_PopulateUserInfo_Windows(t *testing.T) {

	osInfo := newOSInfo()
	osInfo.families.Add("windows")

	info := newUserInfo()

	transport := newMockTransport()
	transport.powerShellResponses[UserNamePowerShell] = &commandResponse{
		stdout: "testuser\n",
	}
	transport.powerShellResponses[UserHomeDirPowerShell] = &commandResponse{
		stdout: "C:\\Users\\testuser\n",
	}
	transport.powerShellResponses[UserIdPowerShell] = &commandResponse{
		stdout: "S-1-5-21-1234567890-1234567890-1234567890-1001\n",
	}

	err := info.populateUserInfo(osInfo, transport)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if info.Name() != "testuser" {
		t.Errorf("expected name 'testuser', got '%s'", info.Name())
	}
	if info.HomeDir() != "C:\\Users\\testuser" {
		t.Errorf("expected home directory 'C:\\Users\\testuser', got '%s'", info.HomeDir())
	}
	if info.UserId() != "S-1-5-21-1234567890-1234567890-1234567890-1001" {
		t.Errorf("expected user ID 'S-1-5-21-1234567890-1234567890-1234567890-1001', got '%s'", info.UserId())
	}
}
