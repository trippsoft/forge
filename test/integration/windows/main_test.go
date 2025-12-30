// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/bmatcuk/go-vagrant"
	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/plugin"
)

const (
	inventoryTemplate = `
		host "cmd" {
			transport "ssh" {
				host = %q
				port = %d
				user = %q
				password = %q
				use_known_hosts = false
			}
		}

		host "windows" {
			transport "ssh" {
				host = %q
				port = %d
				user = %q
				password = %q
				use_known_hosts = false
			}
		}
		`
)

var (
	vagrantClient *vagrant.VagrantClient

	windowsHost           string
	windowsPort           uint16
	windowsUser           string
	windowsPrivateKeyPath string
	windowsPrivateKey     []byte
	windowsPassword       string

	cmdHost           string
	cmdPort           uint16
	cmdUser           string
	cmdPrivateKeyPath string
	cmdPrivateKey     []byte
	cmdPassword       string

	inventoryContent string

	inv *inventory.Inventory

	moduleRegistry *module.Registry
)

func TestMain(m *testing.M) {

	cmd := exec.Command("bash", "-c", "./build_plugins.sh")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	directory, err := os.Getwd()
	if err != nil {
		cmd := exec.Command("bash", "-c", "./stop.sh")
		cmd.Run()
		panic(err)
	}

	plugin.SharedPluginBasePath = directory + "/plugins"
	plugin.UserPluginBasePath = directory + "/plugins"

	err = setupVagrantEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set up Vagrant environment: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if vagrantClient == nil {
		os.Exit(code)
	}

	vagrantDestroy := vagrantClient.Destroy()
	_ = vagrantDestroy.Run()

	os.Exit(code)
}

func setupVagrantEnvironment() error {
	var err error
	vagrantClient, err = vagrant.NewVagrantClient(".")
	if err != nil {
		return fmt.Errorf("failed to create Vagrant client: %w", err)
	}

	vagrantUp := vagrantClient.Up()
	err = vagrantUp.Run()
	if err != nil {
		return fmt.Errorf("failed to run Vagrant up: %w", err)
	}

	vagrantSshInfo := vagrantClient.SSHConfig()
	err = vagrantSshInfo.Run()
	if err != nil {
		return fmt.Errorf("failed to get Vagrant SSH info: %w", err)
	}

	windowsHost = vagrantSshInfo.Configs["windows"].HostName
	windowsPort = uint16(vagrantSshInfo.Configs["windows"].Port)
	windowsUser = vagrantSshInfo.Configs["windows"].User
	windowsPassword = "vagrant"

	cmdHost = vagrantSshInfo.Configs["cmd"].HostName
	cmdPort = uint16(vagrantSshInfo.Configs["cmd"].Port)
	cmdUser = vagrantSshInfo.Configs["cmd"].User
	cmdPassword = "vagrant"

	windowsPrivateKeyPath = vagrantSshInfo.Configs["windows"].IdentityFile
	cmdPrivateKeyPath = vagrantSshInfo.Configs["cmd"].IdentityFile

	windowsPrivateKey, err = os.ReadFile(windowsPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read Windows private key: %w", err)
	}

	cmdPrivateKey, err = os.ReadFile(cmdPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read CMD private key: %w", err)
	}

	if inventoryContent == "" {
		inventoryContent = fmt.Sprintf(
			inventoryTemplate,
			cmdHost,
			cmdPort,
			cmdUser,
			cmdPassword,
			windowsHost,
			windowsPort,
			windowsUser,
			windowsPassword,
		)
	}

	if inv == nil {
		var diags hcl.Diagnostics
		inventoryFile := &inventory.InventoryFile{
			Content: []byte(inventoryContent),
			Path:    "test_inventory.hcl",
		}
		inv, diags = inventory.ParseInventoryFiles(inventoryFile)
		if diags.HasErrors() {
			return fmt.Errorf("failed to parse inventory files: %v", diags)
		}
	}

	if moduleRegistry == nil {
		moduleRegistry = module.NewRegistry()
		moduleRegistry.RegisterCoreModules()
		moduleRegistry.RegisterPluginModules()
	}

	return nil
}

func createTempKnownHostsFile(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_known_hosts_empty_*")
	if err != nil {
		t.Fatalf("Failed to create empty temp known hosts file: %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

func cleanupTempFile(t *testing.T, path string) {
	t.Helper()

	err := os.Remove(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOENT) {
		t.Logf("Warning: failed to cleanup temp file %s: %v", path, err)
	}
}
