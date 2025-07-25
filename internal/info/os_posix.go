package info

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/trippsoft/forge/internal/log"
	"github.com/trippsoft/forge/internal/transport"
)

func (o *osInfo) populatePosixOSInfo(transport transport.Transport, fileSystem transport.FileSystem) error {

	err := o.populatePosixArchitecture(transport)
	if err != nil {
		return fmt.Errorf("failed to populate POSIX architecture: %w", err)
	}

	system, err := o.getPosixKernelFromUname(transport)
	if err != nil {
		return fmt.Errorf("failed to get POSIX kernel: %w", err)
	}

	o.families.Add(system)

	switch system {
	case "darwin":
		return o.populateDarwinOSInfo(transport)
	case "linux":
		return o.populateLinuxOSInfo(transport, fileSystem)
	default:
		log.Warnf("unknown POSIX OS family %q detected", system)
	}

	return nil
}

func (o *osInfo) populatePosixArchitecture(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "uname -m")
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	archString := strings.ToLower(strings.TrimSpace(stdout))

	arch, exists := architectureMap[archString]
	if !exists {
		log.Warnf("unknown architecture %q detected", archString)
		o.procArch = archString
		o.procArchBits = 0
		o.osArch = archString
		o.osArchBits = 0
		return nil
	}

	o.procArch = arch
	o.osArch = arch

	archBits, exists := architectureBitsMap[arch]
	if !exists {
		log.Warnf("unknown architecture bits for %q detected", arch)
		o.procArchBits = 0
		o.osArchBits = 0
		return nil
	}

	o.procArchBits = archBits
	o.osArchBits = archBits

	return nil
}

func (o *osInfo) getPosixKernelFromUname(transport transport.Transport) (string, error) {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "uname -s")
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return strings.ToLower(strings.TrimSpace(stdout)), nil
}

func (o *osInfo) populateDarwinOSInfo(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "/usr/bin/sw_vers -productVersion")
	if err != nil {
		return fmt.Errorf("failed to execute command for OS version: %w", err)
	}

	o.id = "darwin"
	o.friendlyName = "macOS"
	o.version = strings.TrimSpace(stdout)

	versionParts := strings.Split(o.version, ".")
	o.majorVersion = versionParts[0]

	switch {
	case o.majorVersion == "26":
		o.release = "Tahoe"
	case o.majorVersion == "15":
		o.release = "Sequoia"
	case o.majorVersion == "14":
		o.release = "Sonoma"
	case o.majorVersion == "13":
		o.release = "Ventura"
	case o.majorVersion == "12":
		o.release = "Monterey"
	case o.majorVersion == "11":
		o.release = "Big Sur"
	default:
		log.Warnf("possibly unsupported macOS version %s detected", o.majorVersion)
	}

	return nil
}

func (o *osInfo) populateLinuxOSInfo(transport transport.Transport, fileSystem transport.FileSystem) error {

	osReleaseData, err := o.parseOSReleaseFile(fileSystem)
	if err != nil && !errors.Is(err, osNoReleaseFileError) {
		log.Warnf("failed to parse os-release file(s): %v", err)
	}

	osID, exists := osReleaseData["ID"]
	if exists {
		o.id = strings.ToLower(osID)
	}

	osPrettyName, exists := osReleaseData["PRETTY_NAME"]
	if exists {
		o.friendlyName = osPrettyName
	}

	osVersionID, exists := osReleaseData["VERSION_ID"]
	if exists {
		o.version = osVersionID
	}

	osVersionCodename, exists := osReleaseData["VERSION_CODENAME"]
	if exists {
		o.release = osVersionCodename
	}

	osVariant, exists := osReleaseData["VARIANT"]
	if exists {
		o.edition = osVariant
	}

	osVariantID, exists := osReleaseData["VARIANT_ID"]
	if exists {
		o.editionId = strings.ToLower(osVariantID)
	}

	lsbReleaseData, err := o.parseLSBReleaseOutput(transport)

	if o.id == "" {
		lsbDistributorID, exists := lsbReleaseData["Distributor ID"]
		if exists && lsbDistributorID != "n/a" {
			o.id = strings.ToLower(lsbDistributorID)
		}
	}

	if o.friendlyName == "" {
		lsbDescription, exists := lsbReleaseData["Description"]
		if exists && lsbDescription != "n/a" {
			o.friendlyName = lsbDescription
		}
	}

	if o.version == "" {
		lsbRelease, exists := lsbReleaseData["Release"]
		if exists && lsbRelease != "n/a" {
			o.version = lsbRelease
		}
	}

	if o.release == "" {
		lsbCodename, exists := lsbReleaseData["Codename"]
		if exists && lsbCodename != "n/a" {
			o.release = lsbCodename
		}
	}

	osID, exists = osIDCorrectionMap[o.id]
	if exists {
		o.id = osID // Correct the ID if it exists in the correction map to make it consistent and identifiable
	}

	osFamilies, exists := osFamiliesMap[o.id]
	if exists {
		for _, family := range osFamilies {
			o.families.Add(family)
		}
	}

	o.families.Add(o.id)

	return nil
}

func (o *osInfo) parseOSReleaseFile(fileSystem transport.FileSystem) (map[string]string, error) {

	osReleaseFile, err := fileSystem.Open("/etc/os-release")
	if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		return map[string]string{}, fmt.Errorf("failed to open /etc/os-release: %w", err)
	}

	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) && osReleaseFile != nil {
		defer osReleaseFile.Close()
		return o.parseOSReleaseFileContent(osReleaseFile)
	}

	osReleaseFile, err = fileSystem.Open("/usr/lib/os-release")
	if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		return map[string]string{}, fmt.Errorf("failed to open /usr/lib/os-release: %w", err)
	}

	if !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) && osReleaseFile != nil {
		defer osReleaseFile.Close()
		return o.parseOSReleaseFileContent(osReleaseFile)
	}

	return map[string]string{}, osNoReleaseFileError
}

func (o *osInfo) parseOSReleaseFileContent(file transport.File) (map[string]string, error) {

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read os-release file content: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	osReleaseData := make(map[string]string)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"") // Remove quotes around value
		osReleaseData[key] = value
	}

	return osReleaseData, nil
}

func (o *osInfo) parseLSBReleaseOutput(transport transport.Transport) (map[string]string, error) {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "lsb_release -a")
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to execute lsb_release command: %w", err)
	}

	content := make(map[string]string)

	lines := strings.SplitSeq(strings.TrimSpace(stdout), "\n")
	for line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		content[key] = value
	}

	return content, nil
}
