// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"os"
	"testing"

	"github.com/bmatcuk/go-vagrant"
)

var (
	vagrantClient *vagrant.VagrantClient

	linuxHost       string
	linuxPort       uint16
	linuxUser       string
	linuxPrivateKey []byte
	linuxPassword   string

	linuxPWHost       string
	linuxPWPort       uint16
	linuxPWUser       string
	linuxPWPrivateKey []byte
	linuxPWPassword   string

	cmdHost       string
	cmdPort       uint16
	cmdUser       string
	cmdPrivateKey []byte
	cmdPassword   string

	windowsHost       string
	windowsPort       uint16
	windowsUser       string
	windowsPrivateKey []byte
	windowsPassword   string
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

	linuxPrivateKeyPath := vagrantSshInfo.Configs["linux"].IdentityFile
	linuxPWPrivateKeyPath := vagrantSshInfo.Configs["linuxpw"].IdentityFile
	cmdPrivateKeyPath := vagrantSshInfo.Configs["cmd"].IdentityFile
	windowsPrivateKeyPath := vagrantSshInfo.Configs["windows"].IdentityFile

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
}
