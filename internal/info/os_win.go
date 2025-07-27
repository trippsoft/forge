package info

import (
	"context"
	"fmt"
	"strings"

	"github.com/trippsoft/forge/internal/log"
	"github.com/trippsoft/forge/internal/transport"
)

const (
	osFriendlyNamePowerShell = `$object = Get-CimInstance -ClassName Win32_OperatingSystem; Write-Host $object.Caption`
	osVersionPowerShell      = `$version = [Environment]::OSVersion.Version; Write-Host $version.ToString()`
	osArchPowerShell         = `$arch = Get-CimInstance -ClassName Win32_OperatingSystem; Write-Host $arch.OSArchitecture`
	procArchPowerShell       = `$arch = $env:PROCESSOR_ARCHITECTURE; Write-Host $arch`
)

func (o *osInfo) populateWindowsOSInfo(transport transport.Transport) error {

	err := o.populateWindowsProcArchitecture(transport)
	if err != nil {
		return fmt.Errorf("failed to populate Windows processor architecture: %w", err)
	}

	err = o.populateWindowsOSArchitecture(transport)
	if err != nil {
		return fmt.Errorf("failed to populate Windows OS architecture: %w", err)
	}

	friendlyName, err := o.getWindowsFriendlyName(transport)
	if err != nil {
		return fmt.Errorf("failed to get Windows friendly name: %w", err)
	}

	o.friendlyName = friendlyName

	isServer := false

	if strings.Contains(o.friendlyName, "Server") {
		isServer = true
		o.id = "windows-server"
	} else {
		o.id = "windows-client"
	}

	version, err := o.getWindowsVersion(transport)
	if err != nil {
		return fmt.Errorf("failed to get Windows version: %w", err)
	}

	o.version = version

	versionParts := strings.Split(version, ".")
	if len(versionParts) == 0 {
		return fmt.Errorf("failed to parse Windows version: %s", version)
	}

	o.majorVersion = versionParts[0]
	var friendlyNamePrefix string

	switch {
	case strings.HasPrefix(o.version, "6.1.7600") && isServer:
		o.release = "server-2008-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2"
	case strings.HasPrefix(o.version, "6.1.7600") && !isServer:
		o.release = "7"
		friendlyNamePrefix = "Microsoft Windows 7"
	case strings.HasPrefix(o.version, "6.1.7601") && isServer:
		o.release = "server-2008-r2-sp1"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2 SP1"
	case strings.HasPrefix(o.version, "6.1.7601") && !isServer:
		o.release = "7-sp1"
		friendlyNamePrefix = "Microsoft Windows 7 SP1"
	case strings.HasPrefix(o.version, "6.2.9200") && isServer:
		o.release = "server-2012"
		friendlyNamePrefix = "Microsoft Windows Server 2012"
	case strings.HasPrefix(o.version, "6.2.9200") && !isServer:
		o.release = "8"
		friendlyNamePrefix = "Microsoft Windows 8"
	case strings.HasPrefix(o.version, "6.3.9600") && isServer:
		o.release = "server-2012-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2012 R2"
	case strings.HasPrefix(o.version, "6.3.9600") && !isServer:
		o.release = "8.1"
		friendlyNamePrefix = "Microsoft Windows 8.1"
	case strings.HasPrefix(o.version, "10.0.10240"):
		o.release = "10-1507"
		friendlyNamePrefix = "Microsoft Windows 10 1507"
	case strings.HasPrefix(o.version, "10.0.10586"):
		o.release = "10-1511"
		friendlyNamePrefix = "Microsoft Windows 10 1511"
	case strings.HasPrefix(o.version, "10.0.14393") && isServer:
		o.release = "server-2016"
		friendlyNamePrefix = "Microsoft Windows Server 2016"
	case strings.HasPrefix(o.version, "10.0.14393") && !isServer:
		o.release = "10-1607"
		friendlyNamePrefix = "Microsoft Windows 10 1607"
	case strings.HasPrefix(o.version, "10.0.15063"):
		o.release = "10-1703"
		friendlyNamePrefix = "Microsoft Windows 10 1703"
	case strings.HasPrefix(o.version, "10.0.16299"):
		o.release = "10-1709"
		friendlyNamePrefix = "Microsoft Windows 10 1709"
	case strings.HasPrefix(o.version, "10.0.17134"):
		o.release = "10-1803"
		friendlyNamePrefix = "Microsoft Windows 10 1803"
	case strings.HasPrefix(o.version, "10.0.17763") && isServer:
		o.release = "server-2019"
		friendlyNamePrefix = "Microsoft Windows Server 2019"
	case strings.HasPrefix(o.version, "10.0.17763") && !isServer:
		o.release = "10-1809"
		friendlyNamePrefix = "Microsoft Windows 10 1809"
	case strings.HasPrefix(o.version, "10.0.18362"):
		o.release = "10-1903"
		friendlyNamePrefix = "Microsoft Windows 10 1903"
	case strings.HasPrefix(o.version, "10.0.18363") && isServer:
		o.release = "server-1909"
		friendlyNamePrefix = "Microsoft Windows Server 1909"
	case strings.HasPrefix(o.version, "10.0.18363") && !isServer:
		o.release = "10-1909"
		friendlyNamePrefix = "Microsoft Windows 10 1909"
	case strings.HasPrefix(o.version, "10.0.19041") && isServer:
		o.release = "server-2004"
		friendlyNamePrefix = "Microsoft Windows Server 2004"
	case strings.HasPrefix(o.version, "10.0.19041") && !isServer:
		o.release = "10-2004"
		friendlyNamePrefix = "Microsoft Windows 10 2004"
	case strings.HasPrefix(o.version, "10.0.19042") && isServer:
		o.release = "server-20h2"
		friendlyNamePrefix = "Microsoft Windows Server 20H2"
	case strings.HasPrefix(o.version, "10.0.19042") && !isServer:
		o.release = "10-20h2"
		friendlyNamePrefix = "Microsoft Windows 10 20H2"
	case strings.HasPrefix(o.version, "10.0.19043"):
		o.release = "10-21h1"
		friendlyNamePrefix = "Microsoft Windows 10 21H1"
	case strings.HasPrefix(o.version, "10.0.19044"):
		o.release = "10-21h2"
		friendlyNamePrefix = "Microsoft Windows 10 21H2"
	case strings.HasPrefix(o.version, "10.0.19045"):
		o.release = "10-22h2"
		friendlyNamePrefix = "Microsoft Windows 10 22H2"
	case strings.HasPrefix(o.version, "10.0.20348"):
		o.release = "server-2022"
		friendlyNamePrefix = "Microsoft Windows Server 2022"
	case strings.HasPrefix(o.version, "10.0.22000"):
		o.release = "11-21h2"
		friendlyNamePrefix = "Microsoft Windows 11 21H2"
	case strings.HasPrefix(o.version, "10.0.22621"):
		o.release = "11-22h2"
		friendlyNamePrefix = "Microsoft Windows 11 22H2"
	case strings.HasPrefix(o.version, "10.0.22631"):
		o.release = "11-23h2"
		friendlyNamePrefix = "Microsoft Windows 11 23H2"
	case strings.HasPrefix(o.version, "10.0.25398"):
		o.release = "server-23h2"
		friendlyNamePrefix = "Microsoft Windows Server 23H2"
	case strings.HasPrefix(o.version, "10.0.26100") && isServer:
		o.release = "server-2025"
		friendlyNamePrefix = "Microsoft Windows Server 2025"
	case strings.HasPrefix(o.version, "10.0.26100") && !isServer:
		o.release = "11-24h2"
		friendlyNamePrefix = "Microsoft Windows 11 24H2"
	}

	friendlyNameParts := strings.Split(o.friendlyName, " ")
	for len(friendlyNameParts) > 0 && strings.Contains(friendlyNamePrefix, friendlyNameParts[0]) {
		friendlyNameParts = friendlyNameParts[1:]
	}

	friendlyNameBuilder := strings.Builder{}
	friendlyNameBuilder.WriteString(friendlyNamePrefix)
	editionBuilder := strings.Builder{}
	editionIdBuilder := strings.Builder{}

	for _, part := range friendlyNameParts {
		if strings.Contains(part, "Edition") || strings.Contains(part, "Evaluation") {
			continue
		}

		friendlyNameBuilder.WriteString(" " + part)
		if editionBuilder.Len() == 0 {
			editionBuilder.WriteString(part)
			editionIdBuilder.WriteString(strings.ToLower(part))
		} else {
			editionBuilder.WriteString(" " + part)
			editionIdBuilder.WriteString("-" + strings.ToLower(part))
		}
	}

	o.friendlyName = strings.TrimSpace(friendlyNameBuilder.String())
	o.edition = strings.TrimSpace(editionBuilder.String())
	o.editionId = strings.TrimSpace(editionIdBuilder.String())

	o.families.Add(o.id)

	return nil
}

func (o *osInfo) populateWindowsProcArchitecture(transport transport.Transport) error {

	stdout, err := transport.ExecutePowerShell(context.Background(), procArchPowerShell)
	if err != nil {
		return fmt.Errorf("failed to exec PowerShell command: %w", err)
	}

	procArchString := strings.TrimSpace(strings.ToLower(stdout))

	procArch, exists := architectureMap[procArchString]
	if !exists {
		log.Warnf("unknown architecture %s detected", procArchString)
		o.procArch = procArchString
		o.procArchBits = 0
		return nil
	}

	o.procArch = procArch

	procArchBits, exists := architectureBitsMap[o.procArch]
	if !exists {
		log.Warnf("unknown architecture bits for %s detected", o.procArch)
		o.procArchBits = 0
		return nil
	}

	o.procArchBits = procArchBits

	return nil
}

func (o *osInfo) populateWindowsOSArchitecture(transport transport.Transport) error {

	stdout, err := transport.ExecutePowerShell(context.Background(), osArchPowerShell)
	if err != nil {
		return fmt.Errorf("failed to execute command for OS architecture: %w", err)
	}

	osBits := strings.TrimSpace(stdout)

	switch {
	case osBits == "64-bit" && o.procArch == "amd64":
		o.osArch = "amd64"
		o.osArchBits = 64
	case osBits == "64-bit" && o.procArch == "arm64":
		o.osArch = "arm64"
		o.osArchBits = 64
	case osBits == "32-bit" && o.procArch == "amd64":
		o.osArch = "386"
		o.osArchBits = 32
	case osBits == "32-bit" && o.procArch == "arm64":
		o.osArch = "arm"
		o.osArchBits = 32
	case osBits == "32-bit" && o.procArch == "386":
		o.osArch = "386"
		o.osArchBits = 32
	case osBits == "32-bit" && o.procArch == "arm":
		o.osArch = "arm"
		o.osArchBits = 32
	default:
		log.Warnf("unknown OS architecture detected: %s on %s machine", osBits, o.procArch)
		o.osArch = o.procArch
		o.osArchBits = o.procArchBits
	}

	return nil
}

func (o *osInfo) getWindowsFriendlyName(transport transport.Transport) (string, error) {

	stdout, err := transport.ExecutePowerShell(context.Background(), osFriendlyNamePowerShell)
	if err != nil {
		return "", fmt.Errorf("failed to execute command for OS friendly name: %w", err)
	}

	return strings.TrimSpace(stdout), nil
}

func (o *osInfo) getWindowsVersion(transport transport.Transport) (string, error) {

	stdout, err := transport.ExecutePowerShell(context.Background(), osVersionPowerShell)
	if err != nil {
		return "", fmt.Errorf("failed to execute command for OS version: %w", err)
	}

	return strings.TrimSpace(stdout), nil
}
