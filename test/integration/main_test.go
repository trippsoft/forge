// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/bmatcuk/go-vagrant"
	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
)

const (
	inventoryTemplate = `
		host "linux" {
			transport "ssh" {
				host = %q
				port = %d
				user = %q
				private_key_path = %q
				use_known_hosts = false
			}

			escalate {
				password = %q
			}
		}

		host "linuxpw" {
			transport "ssh" {
				host = %q
				port = %d
				user = %q
				private_key_path = %q
				use_known_hosts = false
			}

			escalate {
				password = %q
			}
		}

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

	linuxHost           string
	linuxPort           uint16
	linuxUser           string
	linuxPrivateKeyPath string
	linuxPrivateKey     []byte
	linuxPassword       string

	linuxPWHost           string
	linuxPWPort           uint16
	linuxPWUser           string
	linuxPWPrivateKeyPath string
	linuxPWPrivateKey     []byte
	linuxPWPassword       string

	cmdHost           string
	cmdPort           uint16
	cmdUser           string
	cmdPrivateKeyPath string
	cmdPrivateKey     []byte
	cmdPassword       string

	windowsHost           string
	windowsPort           uint16
	windowsUser           string
	windowsPrivateKeyPath string
	windowsPrivateKey     []byte
	windowsPassword       string

	inventoryContent string

	inv *inventory.Inventory
)

func TestMain(m *testing.M) {

	code := m.Run()

	if vagrantClient == nil {
		os.Exit(code)
	}

	vagrantDestroy := vagrantClient.Destroy()
	_ = vagrantDestroy.Run()

	os.Exit(code)
}

func setupVagrantEnvironment(t testing.TB) {

	t.Helper()

	var err error
	vagrantClient, err = vagrant.NewVagrantClient(".")
	if err != nil {
		t.Fatalf("Failed to create Vagrant client: %v", err)
	}

	vagrantUp := vagrantClient.Up()
	err = vagrantUp.Run()
	if err != nil {
		t.Fatalf("Failed to run Vagrant up: %v", err)
	}

	vagrantSshInfo := vagrantClient.SSHConfig()
	err = vagrantSshInfo.Run()
	if err != nil {
		t.Fatalf("Failed to get Vagrant SSH info: %v", err)
	}

	linuxHost = vagrantSshInfo.Configs["linux"].HostName
	linuxPort = uint16(vagrantSshInfo.Configs["linux"].Port)
	linuxUser = vagrantSshInfo.Configs["linux"].User
	linuxPassword = "vagrant"

	linuxPWHost = vagrantSshInfo.Configs["linuxpw"].HostName
	linuxPWPort = uint16(vagrantSshInfo.Configs["linuxpw"].Port)
	linuxPWUser = vagrantSshInfo.Configs["linuxpw"].User
	linuxPWPassword = "vagrant"

	cmdHost = vagrantSshInfo.Configs["cmd"].HostName
	cmdPort = uint16(vagrantSshInfo.Configs["cmd"].Port)
	cmdUser = vagrantSshInfo.Configs["cmd"].User
	cmdPassword = "vagrant"

	windowsHost = vagrantSshInfo.Configs["windows"].HostName
	windowsPort = uint16(vagrantSshInfo.Configs["windows"].Port)
	windowsUser = vagrantSshInfo.Configs["windows"].User
	windowsPassword = "vagrant"

	linuxPrivateKeyPath = vagrantSshInfo.Configs["linux"].IdentityFile
	linuxPWPrivateKeyPath = vagrantSshInfo.Configs["linuxpw"].IdentityFile
	cmdPrivateKeyPath = vagrantSshInfo.Configs["cmd"].IdentityFile
	windowsPrivateKeyPath = vagrantSshInfo.Configs["windows"].IdentityFile

	linuxPrivateKey, err = os.ReadFile(linuxPrivateKeyPath)
	if err != nil {
		t.Fatalf("Failed to read Linux private key: %v", err)
	}

	linuxPWPrivateKey, err = os.ReadFile(linuxPWPrivateKeyPath)
	if err != nil {
		t.Fatalf("Failed to read Linux PW private key: %v", err)
	}

	cmdPrivateKey, err = os.ReadFile(cmdPrivateKeyPath)
	if err != nil {
		t.Fatalf("Failed to read CMD private key: %v", err)
	}

	windowsPrivateKey, err = os.ReadFile(windowsPrivateKeyPath)
	if err != nil {
		t.Fatalf("Failed to read Windows private key: %v", err)
	}

	if inventoryContent == "" {
		inventoryContent = fmt.Sprintf(
			inventoryTemplate,
			linuxHost,
			linuxPort,
			linuxUser,
			linuxPrivateKeyPath,
			linuxPassword,
			linuxPWHost,
			linuxPWPort,
			linuxPWUser,
			linuxPWPrivateKeyPath,
			linuxPWPassword,
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
		inv, diags = inventory.ParseTestInventoryFile(inventoryContent)

		if diags.HasErrors() {
			t.Fatalf("Failed to parse inventory files: %v", diags)
		}
	}
}
