// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package discover

import (
	"errors"
	"runtime"
	"strings"

	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
)

func discoverOSInfo() (*OSInfoResponse, error) {
	b := &backend.Local{}
	shell, err := ps.New(b)
	if err != nil {
		return nil, err
	}
	defer shell.Exit()

	osInfo := &OSInfoResponse{
		Kernel: "windows",
		Arch:   runtime.GOARCH,
	}

	stdout, _, err := shell.Execute("[System.Environment]::OSVersion.Version.ToString()")
	if err != nil {
		return nil, err
	}

	osInfo.Version = strings.TrimSpace(stdout)
	versionParts := strings.Split(osInfo.Version, ".")
	if len(versionParts) < 3 {
		return nil, errors.New("unexpected Windows version format")
	}

	osInfo.MajorVersion = versionParts[0]

	stdout, _, err = shell.Execute("(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
	if err != nil {
		return nil, err
	}

	caption := strings.TrimSpace(stdout)
	isServer := strings.Contains(strings.ToLower(caption), "server")

	if isServer {
		osInfo.Id = "windows-server"
		osInfo.Families = []string{"windows", "windows-server"}
	} else {
		osInfo.Id = "windows-client"
		osInfo.Families = []string{"windows", "windows-client"}
	}

	var friendlyNamePrefix string
	switch {
	case strings.HasPrefix(osInfo.Version, "6.1.7600") && isServer:
		osInfo.Release = "Server 2008 R2"
		osInfo.ReleaseId = "server-2008-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2"
	case strings.HasPrefix(osInfo.Version, "6.1.7600") && !isServer:
		osInfo.Release = "7"
		osInfo.ReleaseId = "7"
		friendlyNamePrefix = "Microsoft Windows 7"
	case strings.HasPrefix(osInfo.Version, "6.1.7601") && isServer:
		osInfo.Release = "Server 2008 R2 SP1"
		osInfo.ReleaseId = "server-2008-r2-sp1"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2 SP1"
	case strings.HasPrefix(osInfo.Version, "6.1.7601") && !isServer:
		osInfo.Release = "7 SP1"
		osInfo.ReleaseId = "7-sp1"
		friendlyNamePrefix = "Microsoft Windows 7 SP1"
	case strings.HasPrefix(osInfo.Version, "6.2.9200") && isServer:
		osInfo.Release = "Server 2012"
		osInfo.ReleaseId = "server-2012"
		friendlyNamePrefix = "Microsoft Windows Server 2012"
	case strings.HasPrefix(osInfo.Version, "6.2.9200") && !isServer:
		osInfo.Release = "8"
		osInfo.ReleaseId = "8"
		friendlyNamePrefix = "Microsoft Windows 8"
	case strings.HasPrefix(osInfo.Version, "6.3.9600") && isServer:
		osInfo.Release = "Server 2012 R2"
		osInfo.ReleaseId = "server-2012-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2012 R2"
	case strings.HasPrefix(osInfo.Version, "6.3.9600") && !isServer:
		osInfo.Release = "8.1"
		osInfo.ReleaseId = "8.1"
		friendlyNamePrefix = "Microsoft Windows 8.1"
	case strings.HasPrefix(osInfo.Version, "10.0.10240"):
		osInfo.Release = "10 1507"
		osInfo.ReleaseId = "10-1507"
		friendlyNamePrefix = "Microsoft Windows 10 1507"
	case strings.HasPrefix(osInfo.Version, "10.0.10586"):
		osInfo.Release = "10 1511"
		osInfo.ReleaseId = "10-1511"
		friendlyNamePrefix = "Microsoft Windows 10 1511"
	case strings.HasPrefix(osInfo.Version, "10.0.14393") && isServer:
		osInfo.Release = "Server 2016"
		osInfo.ReleaseId = "server-2016"
		friendlyNamePrefix = "Microsoft Windows Server 2016"
	case strings.HasPrefix(osInfo.Version, "10.0.14393") && !isServer:
		osInfo.Release = "10 1607"
		osInfo.ReleaseId = "10-1607"
		friendlyNamePrefix = "Microsoft Windows 10 1607"
	case strings.HasPrefix(osInfo.Version, "10.0.15063"):
		osInfo.Release = "10 1703"
		osInfo.ReleaseId = "10-1703"
		friendlyNamePrefix = "Microsoft Windows 10 1703"
	case strings.HasPrefix(osInfo.Version, "10.0.16299"):
		osInfo.Release = "10 1709"
		osInfo.ReleaseId = "10-1709"
		friendlyNamePrefix = "Microsoft Windows 10 1709"
	case strings.HasPrefix(osInfo.Version, "10.0.17134"):
		osInfo.Release = "10 1803"
		osInfo.ReleaseId = "10-1803"
		friendlyNamePrefix = "Microsoft Windows 10 1803"
	case strings.HasPrefix(osInfo.Version, "10.0.17763") && isServer:
		osInfo.Release = "Server 2019"
		osInfo.ReleaseId = "server-2019"
		friendlyNamePrefix = "Microsoft Windows Server 2019"
	case strings.HasPrefix(osInfo.Version, "10.0.17763") && !isServer:
		osInfo.Release = "10 1809"
		osInfo.ReleaseId = "10-1809"
		friendlyNamePrefix = "Microsoft Windows 10 1809"
	case strings.HasPrefix(osInfo.Version, "10.0.18362"):
		osInfo.Release = "10 1903"
		osInfo.ReleaseId = "10-1903"
		friendlyNamePrefix = "Microsoft Windows 10 1903"
	case strings.HasPrefix(osInfo.Version, "10.0.18363") && isServer:
		osInfo.Release = "Server 1909"
		osInfo.ReleaseId = "server-1909"
		friendlyNamePrefix = "Microsoft Windows Server 1909"
	case strings.HasPrefix(osInfo.Version, "10.0.18363") && !isServer:
		osInfo.Release = "10 1909"
		osInfo.ReleaseId = "10-1909"
		friendlyNamePrefix = "Microsoft Windows 10 1909"
	case strings.HasPrefix(osInfo.Version, "10.0.19041") && isServer:
		osInfo.Release = "Server 2004"
		osInfo.ReleaseId = "server-2004"
		friendlyNamePrefix = "Microsoft Windows Server 2004"
	case strings.HasPrefix(osInfo.Version, "10.0.19041") && !isServer:
		osInfo.Release = "10 2004"
		osInfo.ReleaseId = "10-2004"
		friendlyNamePrefix = "Microsoft Windows 10 2004"
	case strings.HasPrefix(osInfo.Version, "10.0.19042") && isServer:
		osInfo.Release = "Server 20H2"
		osInfo.ReleaseId = "server-20h2"
		friendlyNamePrefix = "Microsoft Windows Server 20H2"
	case strings.HasPrefix(osInfo.Version, "10.0.19042") && !isServer:
		osInfo.Release = "10 20H2"
		osInfo.ReleaseId = "10-20h2"
		friendlyNamePrefix = "Microsoft Windows 10 20H2"
	case strings.HasPrefix(osInfo.Version, "10.0.19043"):
		osInfo.Release = "10 21H1"
		osInfo.ReleaseId = "10-21h1"
		friendlyNamePrefix = "Microsoft Windows 10 21H1"
	case strings.HasPrefix(osInfo.Version, "10.0.19044"):
		osInfo.Release = "10 21H2"
		osInfo.ReleaseId = "10-21h2"
		friendlyNamePrefix = "Microsoft Windows 10 21H2"
	case strings.HasPrefix(osInfo.Version, "10.0.19045"):
		osInfo.Release = "10 22H2"
		osInfo.ReleaseId = "10-22h2"
		friendlyNamePrefix = "Microsoft Windows 10 22H2"
	case strings.HasPrefix(osInfo.Version, "10.0.20348"):
		osInfo.Release = "Server 2022"
		osInfo.ReleaseId = "server-2022"
		friendlyNamePrefix = "Microsoft Windows Server 2022"
	case strings.HasPrefix(osInfo.Version, "10.0.22000"):
		osInfo.Release = "11"
		osInfo.ReleaseId = "11-21h2"
		friendlyNamePrefix = "Microsoft Windows 11 21H2"
	case strings.HasPrefix(osInfo.Version, "10.0.22621"):
		osInfo.Release = "11 22H2"
		osInfo.ReleaseId = "11-22h2"
		friendlyNamePrefix = "Microsoft Windows 11 22H2"
	case strings.HasPrefix(osInfo.Version, "10.0.22631"):
		osInfo.Release = "11 23H2"
		osInfo.ReleaseId = "11-23h2"
		friendlyNamePrefix = "Microsoft Windows 11 23H2"
	case strings.HasPrefix(osInfo.Version, "10.0.25398"):
		osInfo.Release = "Server 23H2"
		osInfo.ReleaseId = "server-23h2"
		friendlyNamePrefix = "Microsoft Windows Server 23H2"
	case strings.HasPrefix(osInfo.Version, "10.0.26100") && isServer:
		osInfo.Release = "Server 2025"
		osInfo.ReleaseId = "server-2025"
		friendlyNamePrefix = "Microsoft Windows Server 2025"
	case strings.HasPrefix(osInfo.Version, "10.0.26100") && !isServer:
		osInfo.Release = "11 24H2"
		osInfo.ReleaseId = "11-24h2"
		friendlyNamePrefix = "Microsoft Windows 11 24H2"
	case strings.HasPrefix(osInfo.Version, "10.0.26200"):
		osInfo.Release = "11 25H2"
		osInfo.ReleaseId = "11-25h2"
		friendlyNamePrefix = "Microsoft Windows 11 25H2"
	default:
		osInfo.Release = ""
		osInfo.ReleaseId = ""
		osInfo.FriendlyName = caption
		osInfo.Edition = ""
		osInfo.EditionId = ""
		return osInfo, nil
	}

	captionParts := strings.Split(caption, " ")
	for len(captionParts) > 0 && strings.Contains(friendlyNamePrefix, captionParts[0]) {
		captionParts = captionParts[1:]
	}

	friendlyNameBuilder := strings.Builder{}
	friendlyNameBuilder.WriteString(friendlyNamePrefix)

	editionBuilder := strings.Builder{}
	editionIdBuilder := strings.Builder{}

	for _, part := range captionParts {
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

	osInfo.FriendlyName = strings.TrimSpace(friendlyNameBuilder.String())
	osInfo.Edition = strings.TrimSpace(editionBuilder.String())
	osInfo.EditionId = strings.TrimSpace(editionIdBuilder.String())

	return osInfo, nil
}
