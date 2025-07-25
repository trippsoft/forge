package info

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

var (
	packageManagerMap = map[string]string{
		"/QOpenSys/pkgs/bin/yum": "yum",
		"/usr/bin/installp":      "installp",
		"/usr/sbin/sorcery":      "sorcery",
		"/usr/bin/swupd":         "swupd",
		"/usr/local/sbin/pkg":    "pkgng",
		"/usr/bin/xbps-install":  "xbps",
		"/usr/bin/pkg":           "pkg5",
		"/usr/sbin/pkgadd":       "svr4pkg",
		"/usr/bin/emerge":        "portage",
		"/usr/sbin/swlist":       "swdepot",
		"/usr/sbin/pkg":          "pkgng",
		"/sbin/apk":              "apk",
		"/opt/homebrew/bin/brew": "homebrew",
		"/usr/local/bin/brew":    "homebrew",
		"/opt/local/bin/port":    "macports",
		"/opt/tools/bin/pkgin":   "pkgin",
		"/opt/local/bin/pkgin":   "pkgin",
		"/usr/pkg/bin/pkgin":     "pkgin",
		"/bin/opkg":              "opkg",
		"/usr/bin/pacman":        "pacman",
		"/usr/sbin/urpmi":        "urpmi",
		"/usr/bin/zypper":        "zypper",
		"/usr/bin/apt-get":       "apt",
		"/usr/bin/dnf5":          "dnf5",
		"/usr/bin/dnf-3":         "dnf",
		"/usr/bin/dnf":           "dnf",
		"/usr/bin/microdnf":      "dnf",
		"/usr/bin/yum":           "yum",
	}

	packageManagers = []string{
		"/QOpenSys/pkgs/bin/yum",
		"/usr/bin/installp",
		"/usr/sbin/sorcery",
		"/usr/bin/swupd",
		"/usr/local/sbin/pkg",
		"/usr/bin/xbps-install",
		"/usr/bin/pkg",
		"/usr/sbin/pkgadd",
		"/usr/bin/emerge",
		"/usr/sbin/swlist",
		"/usr/sbin/pkg",
		"/sbin/apk",
		"/opt/homebrew/bin/brew",
		"/usr/local/bin/brew",
		"/opt/local/bin/port",
		"/opt/tools/bin/pkgin",
		"/opt/local/bin/pkgin",
		"/usr/pkg/bin/pkgin",
		"/bin/opkg",
		"/usr/bin/pacman",
		"/usr/sbin/urpmi",
		"/usr/bin/zypper",
		"/usr/bin/apt-get",
		"/usr/bin/dnf5",
		"/usr/bin/dnf-3",
		"/usr/bin/dnf",
		"/usr/bin/microdnf",
		"/usr/bin/yum",
	}

	elPackageManagers = []string{
		"/usr/bin/dnf5",
		"/usr/bin/dnf-3",
		"/usr/bin/dnf",
		"/usr/bin/microdnf",
		"/usr/bin/yum",
	}

	debianPackageManagers = []string{
		"/usr/bin/apt-get",
	}
)

type packageManagerInfo struct {
	supported bool
	name      string
	path      string
}

func newPackageManagerInfo() *packageManagerInfo {
	return &packageManagerInfo{
		supported: false,
		name:      "",
		path:      "",
	}
}
func (p *packageManagerInfo) populatePackageManagerInfo(osInfo *osInfo, transport transport.Transport, fileSystem transport.FileSystem) error {

	if osInfo.families.Contains("windows") {
		p.supported = false // Windows does not have a standard package manager like Linux
		p.name = ""
		p.path = ""
		return nil
	}

	if osInfo.families.Contains("el") {
		return p.populateELPackageManagerInfo(fileSystem)
	}

	if osInfo.families.Contains("debian") {
		return p.populateDebianPackageManagerInfo(fileSystem)
	}

	if osInfo.families.Contains("altlinux") {
		err := p.populateDebianPackageManagerInfo(fileSystem)
		if err != nil {
			return err
		}

		if p.supported {
			p.name = "apt-rpm"
		}

		return nil
	}

	matchingPackageManager, err := p.getFirstMatchingPackageManager(fileSystem, packageManagers)
	if err != nil {
		return err
	}

	if matchingPackageManager == "" {
		p.supported = false
		p.name = ""
		p.path = ""
		return nil
	}

	p.supported = true
	p.name = packageManagerMap[matchingPackageManager]
	p.path = matchingPackageManager

	if p.name == "apt" {
		err = p.determineIfAptIsRpmBacked(transport, fileSystem)
		if err != nil {
			return fmt.Errorf("failed to determine if apt is RPM-backed: %w", err)
		}
	}

	return nil
}

func (p *packageManagerInfo) populateELPackageManagerInfo(fileSystem transport.FileSystem) error {

	osTreeBooted, err := isOSTreeBooted(fileSystem)
	if err != nil {
		return fmt.Errorf("failed to check if OSTree is booted: %w", err)
	}

	if osTreeBooted {
		p.supported = false // Atomic containers do not use traditional package managers
		p.name = ""
		p.path = ""
		return nil
	}

	matchingPackageManager, err := p.getFirstMatchingPackageManager(fileSystem, elPackageManagers)
	if err != nil {
		return fmt.Errorf("failed to get matching package managers: %w", err)
	}

	if matchingPackageManager == "" {
		p.supported = false
		p.name = ""
		p.path = ""
		return nil
	}

	p.supported = true
	p.name = packageManagerMap[matchingPackageManager]
	p.path = matchingPackageManager

	return nil
}

func isOSTreeBooted(fileSystem transport.FileSystem) (bool, error) {
	_, err := fileSystem.Stat("/run/ostree-booted")
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return false, nil // Not booted with OSTree
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat /run/ostree-booted: %w", err)
	}
	return true, nil // Booted with OSTree
}

func (p *packageManagerInfo) populateDebianPackageManagerInfo(fileSystem transport.FileSystem) error {

	matchingPackageManager, err := p.getFirstMatchingPackageManager(fileSystem, debianPackageManagers)
	if err != nil {
		return fmt.Errorf("failed to get matching package managers: %w", err)
	}

	if matchingPackageManager == "" {
		p.supported = false
		p.name = ""
		p.path = ""
		return nil
	}

	p.supported = true
	p.name = packageManagerMap[matchingPackageManager]
	p.path = matchingPackageManager

	return nil
}

func (p *packageManagerInfo) getFirstMatchingPackageManager(fileSystem transport.FileSystem, possiblePackageManagers []string) (string, error) {

	for _, path := range possiblePackageManagers {
		stat, err := fileSystem.Stat(path)
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
			continue // Skip if the package manager path does not exist
		}
		if err != nil {
			return "", fmt.Errorf("failed to stat package manager path %s: %w", path, err)
		}
		if stat == nil {
			return "", fmt.Errorf("stat returned nil for package manager path %s", path) // This should not happen, but handle it gracefully
		}
		if stat.IsDir() {
			continue
		}
		return path, nil
	}

	return "", nil
}

func (p *packageManagerInfo) determineIfAptIsRpmBacked(transport transport.Transport, fileSystem transport.FileSystem) error {

	_, err := fileSystem.Stat("/usr/bin/rpm")
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		return nil // RPM executable is not present, so not RPM-backed
	}

	if err != nil {
		return fmt.Errorf("failed to stat /usr/bin/rpm: %w", err)
	}

	_, _, err = transport.ExecuteCommand(context.Background(), "/usr/bin/rpm -q --whatprovides /usr/bin/apt-get")

	if err != nil {
		return nil // If the command fails, we assume it's not RPM-backed
	}

	p.name = "apt-rpm" // If the command succeeds, apt-get is provided by an RPM package

	return nil
}

func (p *packageManagerInfo) toMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	if !p.supported {
		values["package_manager_name"] = cty.NullVal(cty.String)
		values["package_manager_path"] = cty.NullVal(cty.String)
		return values
	}

	values["package_manager_name"] = cty.StringVal(p.name)
	values["package_manager_path"] = cty.StringVal(p.path)

	return values
}
