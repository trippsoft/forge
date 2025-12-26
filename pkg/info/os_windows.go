// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package info

import (
	"errors"
	"runtime"
	"strings"

	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
)

func (o *OSInfo) discover() error {
	b := &backend.Local{}
	shell, err := ps.New(b)
	if err != nil {
		return err
	}
	defer shell.Exit()

	o.Kernel = "windows"
	o.Arch = runtime.GOARCH

	stdout, _, err := shell.Execute("[System.Environment]::OSVersion.Version.ToString()")
	if err != nil {
		return err
	}

	o.Version = strings.TrimSpace(stdout)
	versionParts := strings.Split(o.Version, ".")
	if len(versionParts) < 3 {
		return errors.New("unexpected Windows version format")
	}

	o.MajorVersion = versionParts[0]

	stdout, _, err = shell.Execute("(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
	if err != nil {
		return err
	}

	caption := strings.TrimSpace(stdout)
	isServer := strings.Contains(strings.ToLower(caption), "server")

	if isServer {
		o.Id = "windows-server"
		o.Families = []string{"windows", "windows-server"}
	} else {
		o.Id = "windows-client"
		o.Families = []string{"windows", "windows-client"}
	}

	var friendlyNamePrefix string
	switch {
	case strings.HasPrefix(o.Version, "6.1.7600") && isServer:
		o.Release = "Server 2008 R2"
		o.ReleaseId = "server-2008-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2"
	case strings.HasPrefix(o.Version, "6.1.7600") && !isServer:
		o.Release = "7"
		o.ReleaseId = "7"
		friendlyNamePrefix = "Microsoft Windows 7"
	case strings.HasPrefix(o.Version, "6.1.7601") && isServer:
		o.Release = "Server 2008 R2 SP1"
		o.ReleaseId = "server-2008-r2-sp1"
		friendlyNamePrefix = "Microsoft Windows Server 2008 R2 SP1"
	case strings.HasPrefix(o.Version, "6.1.7601") && !isServer:
		o.Release = "7 SP1"
		o.ReleaseId = "7-sp1"
		friendlyNamePrefix = "Microsoft Windows 7 SP1"
	case strings.HasPrefix(o.Version, "6.2.9200") && isServer:
		o.Release = "Server 2012"
		o.ReleaseId = "server-2012"
		friendlyNamePrefix = "Microsoft Windows Server 2012"
	case strings.HasPrefix(o.Version, "6.2.9200") && !isServer:
		o.Release = "8"
		o.ReleaseId = "8"
		friendlyNamePrefix = "Microsoft Windows 8"
	case strings.HasPrefix(o.Version, "6.3.9600") && isServer:
		o.Release = "Server 2012 R2"
		o.ReleaseId = "server-2012-r2"
		friendlyNamePrefix = "Microsoft Windows Server 2012 R2"
	case strings.HasPrefix(o.Version, "6.3.9600") && !isServer:
		o.Release = "8.1"
		o.ReleaseId = "8.1"
		friendlyNamePrefix = "Microsoft Windows 8.1"
	case strings.HasPrefix(o.Version, "10.0.10240"):
		o.Release = "10 1507"
		o.ReleaseId = "10-1507"
		friendlyNamePrefix = "Microsoft Windows 10 1507"
	case strings.HasPrefix(o.Version, "10.0.10586"):
		o.Release = "10 1511"
		o.ReleaseId = "10-1511"
		friendlyNamePrefix = "Microsoft Windows 10 1511"
	case strings.HasPrefix(o.Version, "10.0.14393") && isServer:
		o.Release = "Server 2016"
		o.ReleaseId = "server-2016"
		friendlyNamePrefix = "Microsoft Windows Server 2016"
	case strings.HasPrefix(o.Version, "10.0.14393") && !isServer:
		o.Release = "10 1607"
		o.ReleaseId = "10-1607"
		friendlyNamePrefix = "Microsoft Windows 10 1607"
	case strings.HasPrefix(o.Version, "10.0.15063"):
		o.Release = "10 1703"
		o.ReleaseId = "10-1703"
		friendlyNamePrefix = "Microsoft Windows 10 1703"
	case strings.HasPrefix(o.Version, "10.0.16299"):
		o.Release = "10 1709"
		o.ReleaseId = "10-1709"
		friendlyNamePrefix = "Microsoft Windows 10 1709"
	case strings.HasPrefix(o.Version, "10.0.17134"):
		o.Release = "10 1803"
		o.ReleaseId = "10-1803"
		friendlyNamePrefix = "Microsoft Windows 10 1803"
	case strings.HasPrefix(o.Version, "10.0.17763") && isServer:
		o.Release = "Server 2019"
		o.ReleaseId = "server-2019"
		friendlyNamePrefix = "Microsoft Windows Server 2019"
	case strings.HasPrefix(o.Version, "10.0.17763") && !isServer:
		o.Release = "10 1809"
		o.ReleaseId = "10-1809"
		friendlyNamePrefix = "Microsoft Windows 10 1809"
	case strings.HasPrefix(o.Version, "10.0.18362"):
		o.Release = "10 1903"
		o.ReleaseId = "10-1903"
		friendlyNamePrefix = "Microsoft Windows 10 1903"
	case strings.HasPrefix(o.Version, "10.0.18363") && isServer:
		o.Release = "Server 1909"
		o.ReleaseId = "server-1909"
		friendlyNamePrefix = "Microsoft Windows Server 1909"
	case strings.HasPrefix(o.Version, "10.0.18363") && !isServer:
		o.Release = "10 1909"
		o.ReleaseId = "10-1909"
		friendlyNamePrefix = "Microsoft Windows 10 1909"
	case strings.HasPrefix(o.Version, "10.0.19041") && isServer:
		o.Release = "Server 2004"
		o.ReleaseId = "server-2004"
		friendlyNamePrefix = "Microsoft Windows Server 2004"
	case strings.HasPrefix(o.Version, "10.0.19041") && !isServer:
		o.Release = "10 2004"
		o.ReleaseId = "10-2004"
		friendlyNamePrefix = "Microsoft Windows 10 2004"
	case strings.HasPrefix(o.Version, "10.0.19042") && isServer:
		o.Release = "Server 20H2"
		o.ReleaseId = "server-20h2"
		friendlyNamePrefix = "Microsoft Windows Server 20H2"
	case strings.HasPrefix(o.Version, "10.0.19042") && !isServer:
		o.Release = "10 20H2"
		o.ReleaseId = "10-20h2"
		friendlyNamePrefix = "Microsoft Windows 10 20H2"
	case strings.HasPrefix(o.Version, "10.0.19043"):
		o.Release = "10 21H1"
		o.ReleaseId = "10-21h1"
		friendlyNamePrefix = "Microsoft Windows 10 21H1"
	case strings.HasPrefix(o.Version, "10.0.19044"):
		o.Release = "10 21H2"
		o.ReleaseId = "10-21h2"
		friendlyNamePrefix = "Microsoft Windows 10 21H2"
	case strings.HasPrefix(o.Version, "10.0.19045"):
		o.Release = "10 22H2"
		o.ReleaseId = "10-22h2"
		friendlyNamePrefix = "Microsoft Windows 10 22H2"
	case strings.HasPrefix(o.Version, "10.0.20348"):
		o.Release = "Server 2022"
		o.ReleaseId = "server-2022"
		friendlyNamePrefix = "Microsoft Windows Server 2022"
	case strings.HasPrefix(o.Version, "10.0.22000"):
		o.Release = "11"
		o.ReleaseId = "11-21h2"
		friendlyNamePrefix = "Microsoft Windows 11 21H2"
	case strings.HasPrefix(o.Version, "10.0.22621"):
		o.Release = "11 22H2"
		o.ReleaseId = "11-22h2"
		friendlyNamePrefix = "Microsoft Windows 11 22H2"
	case strings.HasPrefix(o.Version, "10.0.22631"):
		o.Release = "11 23H2"
		o.ReleaseId = "11-23h2"
		friendlyNamePrefix = "Microsoft Windows 11 23H2"
	case strings.HasPrefix(o.Version, "10.0.25398"):
		o.Release = "Server 23H2"
		o.ReleaseId = "server-23h2"
		friendlyNamePrefix = "Microsoft Windows Server 23H2"
	case strings.HasPrefix(o.Version, "10.0.26100") && isServer:
		o.Release = "Server 2025"
		o.ReleaseId = "server-2025"
		friendlyNamePrefix = "Microsoft Windows Server 2025"
	case strings.HasPrefix(o.Version, "10.0.26100") && !isServer:
		o.Release = "11 24H2"
		o.ReleaseId = "11-24h2"
		friendlyNamePrefix = "Microsoft Windows 11 24H2"
	case strings.HasPrefix(o.Version, "10.0.26200"):
		o.Release = "11 25H2"
		o.ReleaseId = "11-25h2"
		friendlyNamePrefix = "Microsoft Windows 11 25H2"
	default:
		o.Release = ""
		o.ReleaseId = ""
		o.FriendlyName = caption
		o.Edition = ""
		o.EditionId = ""
		return nil
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

	o.FriendlyName = strings.TrimSpace(friendlyNameBuilder.String())
	o.Edition = strings.TrimSpace(editionBuilder.String())
	o.EditionId = strings.TrimSpace(editionIdBuilder.String())

	return nil
}
