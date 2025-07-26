package info

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

var (
	serviceManagerCorrectionMap = map[string]string{
		"runit-init":  "runit",
		"openrc-init": "openrc",
	}
)

type serviceManagerInfo struct {
	name string
}

func newServiceManagerInfo() *serviceManagerInfo {
	return &serviceManagerInfo{}
}

func (s *serviceManagerInfo) Name() string {
	return s.name
}

func (s *serviceManagerInfo) populateServiceManagerInfo(osInfo *osInfo, transport transport.Transport, fileSystem transport.FileSystem) error {

	if osInfo.families.Contains("windows") {
		s.name = "windows-service-manager"
		return nil
	}

	err := s.getProcess1FromVirtualFile(transport)
	if err == nil {
		return nil
	}

	err = s.getProcess1FromInitLink(transport)
	if err == nil {
		return nil
	}

	if osInfo.families.Contains("darwin") {
		return s.getDarwinServiceManager(osInfo)
	}

	if osInfo.families.Contains("linux") {
		return s.getLinuxServiceManager(transport, fileSystem)
	}

	// TODO - Support for other OS families

	return errors.New("could not determine service manager")
}

func (s *serviceManagerInfo) getProcess1FromVirtualFile(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "cat /proc/1/comm")
	if err != nil {
		return fmt.Errorf("failed to read /proc/1/comm: %w", err)
	}

	proc1 := strings.TrimSpace(stdout)

	if proc1 == "COMMAND" || proc1 == "init" || strings.HasSuffix(proc1, "sh") {
		return fmt.Errorf("got imprecise or unexpected value in /proc/1/comm: %s", proc1)
	}

	pathParts := strings.Split(proc1, "/") // Not using path.Base here for cross-platform compatibility
	if len(pathParts) > 0 {
		proc1 = pathParts[len(pathParts)-1]
	}

	if correctedName, ok := serviceManagerCorrectionMap[proc1]; ok {
		proc1 = correctedName
	}

	s.name = proc1

	return nil
}

func (s *serviceManagerInfo) getProcess1FromInitLink(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "realpath /sbin/init")
	if err != nil {
		return fmt.Errorf("failed to read /sbin/init link: %w", err)
	}

	proc1 := strings.TrimSpace(stdout)

	if proc1 == "init" || strings.HasSuffix(proc1, "sh") {
		return fmt.Errorf("got imprecise or unexpected value in /sbin/init link: %s", proc1)
	}

	pathParts := strings.Split(proc1, "/") // Not using path.Base here for cross-platform compatibility
	if len(pathParts) > 0 {
		proc1 = pathParts[len(pathParts)-1]
	}

	if correctedName, ok := serviceManagerCorrectionMap[proc1]; ok {
		proc1 = correctedName
	}

	s.name = proc1

	return nil
}

func (s *serviceManagerInfo) getDarwinServiceManager(osInfo *osInfo) error {

	majorVersion, err := strconv.Atoi(osInfo.majorVersion)
	if err != nil {
		return fmt.Errorf("failed to parse macOS major version: %w", err)
	}

	if majorVersion > 10 {
		s.name = "launchd"
	} else if majorVersion == 10 {

		versionParts := strings.Split(osInfo.version, ".")
		if len(versionParts) < 2 {
			return fmt.Errorf("failed to parse macOS version: %s", osInfo.version)
		}

		minorVersion, err := strconv.Atoi(versionParts[1])
		if err != nil {
			return fmt.Errorf("failed to parse macOS minor version: %w", err)
		}

		if minorVersion >= 4 {
			s.name = "launchd"
		} else {
			s.name = "systemstarter"
		}
	} else {
		s.name = "systemstarter"
	}

	return nil
}

func (s *serviceManagerInfo) getLinuxServiceManager(transport transport.Transport, fileSystem transport.FileSystem) error {

	_, _, err := transport.ExecuteCommand(context.Background(), "realpath systemctl")
	if err == nil {
		for _, dir := range []string{"/run/systemd/system", "/dev/.run/systemd", "/dev/.systemd"} {
			_, err := fileSystem.Stat(dir)
			if err == nil {
				s.name = "systemd"
				return nil
			}
		}
	}

	_, _, err = transport.ExecuteCommand(context.Background(), "realpath initctl")
	if err == nil {
		_, err := fileSystem.Stat("/etc/init")
		if err == nil {
			s.name = "upstart"
			return nil
		}
	}

	_, err = fileSystem.Stat("/sbin/openrc")
	if err == nil {
		s.name = "openrc"
		return nil
	}

	_, err = fileSystem.Stat("/etc/init.d")
	if err == nil {
		s.name = "sysvinit"
		return nil
	}

	_, err = fileSystem.Stat("/etc/dinit.d")
	if err == nil {
		s.name = "dinit"
		return nil
	}

	return errors.New("could not determine Linux service manager")
}

func (s *serviceManagerInfo) toMapOfCtyValues() map[string]cty.Value {

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
func (s *serviceManagerInfo) String() string {

	if s.name == "" {
		return "service_manager: unknown"
	}

	return fmt.Sprintf("service_manager: %s", s.name)
}
