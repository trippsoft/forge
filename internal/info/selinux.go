package info

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type selinuxStatus string
type selinuxType string

const (
	SelinuxEnforcing    selinuxStatus = "enforcing"
	SelinuxDisabled     selinuxStatus = "disabled"
	SelinuxPermissive   selinuxStatus = "permissive"
	SelinuxNotSupported selinuxStatus = ""
)

const (
	SelinuxTypeTargeted     selinuxType = "targeted"
	SelinuxTypeMinimum      selinuxType = "minimum"
	SelinuxTypeMLS          selinuxType = "mls"
	SelinuxTypeNotSupported selinuxType = ""
)

type selinuxInfo struct {
	supported   bool
	installed   bool
	status      selinuxStatus
	selinuxType selinuxType
}

func newSELinuxInfo() *selinuxInfo {
	return &selinuxInfo{
		supported:   false,
		installed:   false,
		status:      SelinuxNotSupported,
		selinuxType: SelinuxTypeNotSupported,
	}
}

func (s *selinuxInfo) Supported() bool {
	return s.supported
}

func (s *selinuxInfo) Installed() bool {
	return s.installed
}

func (s *selinuxInfo) Status() selinuxStatus {
	return s.status
}

func (s *selinuxInfo) SelinuxType() selinuxType {
	return s.selinuxType
}

func (s *selinuxInfo) populateSelinuxInfo(osInfo *osInfo, fileSystem transport.FileSystem) error {

	if !osInfo.families.Contains("linux") {
		s.supported = false
		s.installed = false
		s.status = SelinuxNotSupported
		s.selinuxType = SelinuxTypeNotSupported
		return nil
	}

	s.supported = true

	selinuxConfigFile, err := fileSystem.Open("/etc/selinux/config")
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		s.installed = false
		s.status = SelinuxNotSupported
		s.selinuxType = SelinuxTypeNotSupported
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to open SELinux config file: %w", err)
	}
	if selinuxConfigFile == nil {
		return fmt.Errorf("SELinux config file is nil, expected a valid file handle") // This should not happen, but handle it gracefully
	}

	defer selinuxConfigFile.Close()

	s.installed = true

	scanner := bufio.NewScanner(selinuxConfigFile)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue // Skip empty lines and comments
		}

		line = strings.TrimSpace(line)

		switch line {
		case "SELINUX=disabled":
			s.status = SelinuxDisabled
			s.selinuxType = SelinuxTypeNotSupported
			return nil
		case "SELINUX=enforcing":
			s.status = SelinuxEnforcing
		case "SELINUX=permissive":
			s.status = SelinuxPermissive
		case "SELINUXTYPE=targeted":
			s.selinuxType = SelinuxTypeTargeted
		case "SELINUXTYPE=minimum":
			s.selinuxType = SelinuxTypeMinimum
		case "SELINUXTYPE=mls":
			s.selinuxType = SelinuxTypeMLS
		}

		if s.status != "" && s.selinuxType != "" {
			break // Stop scanning once we find SELINUX status and type
		}
	}

	return nil
}

func (s *selinuxInfo) toMapOfCtyValues() map[string]cty.Value {

	if !s.supported {
		return map[string]cty.Value{
			"selinux_installed": cty.NullVal(cty.String),
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	if !s.installed {
		return map[string]cty.Value{
			"selinux_installed": cty.False,
			"selinux_status":    cty.NullVal(cty.String),
			"selinux_type":      cty.NullVal(cty.String),
		}
	}

	return map[string]cty.Value{
		"selinux_installed": cty.True,
		"selinux_status":    cty.StringVal(string(s.status)),
		"selinux_type":      cty.StringVal(string(s.selinuxType)),
	}
}

// String returns a string representation of the SELinux information.
// This is useful for logging or debugging purposes.
func (s *selinuxInfo) String() string {

	stringBuilder := &strings.Builder{}
	if !s.supported {
		stringBuilder.WriteString("selinux_installed: not supported on this OS\n")
		stringBuilder.WriteString("selinux_status: not supported on this OS\n")
		stringBuilder.WriteString("selinux_type: not supported on this OS")

		return stringBuilder.String()
	}

	if !s.installed {
		stringBuilder.WriteString("selinux_installed: false\n")
		stringBuilder.WriteString("selinux_status: not installed\n")
		stringBuilder.WriteString("selinux_type: not installed\n")

		return stringBuilder.String()
	}

	stringBuilder.WriteString("selinux_installed: true\n")
	stringBuilder.WriteString("selinux_status: ")
	stringBuilder.WriteString(string(s.status))
	stringBuilder.WriteString("\n")
	stringBuilder.WriteString("selinux_type: ")
	stringBuilder.WriteString(string(s.selinuxType))

	return stringBuilder.String()
}
