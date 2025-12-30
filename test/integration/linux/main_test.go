// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package linux

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/inventory"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/plugin"
)

const (
	inventoryContent = `
		host "debian13" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2201
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "debian13-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2201
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "debian12" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2202
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "debian12-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2202
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "fedora42" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2211
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "fedora42-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2211
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "fedora41" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2212
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "fedora41-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2212
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "rocky10" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2221
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "rocky10-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2221
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "rocky9" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2222
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "rocky9-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2222
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "rocky8" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2223
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "rocky8-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2223
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "ubuntu2404" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2231
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "ubuntu2404-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2231
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}

		host "ubuntu2204" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2232
				user = "forge"
				password = "forge"
				use_known_hosts = false
			}
		}

		host "ubuntu2204-pw" {
			transport "ssh" {
				host = "127.0.0.1"
				port = 2232
				user = "forge-pw"
				password = "forge"
				use_known_hosts = false
			}

			escalate {
				password = "forge"
			}
		}
		`
)

var (
	inv            *inventory.Inventory
	moduleRegistry *module.Registry

	privateKeyContent []byte
)

func TestMain(m *testing.M) {
	cmd := exec.Command("bash", "-c", "./stop.sh")
	cmd.Run()

	cmd = exec.Command("bash", "-c", "./start.sh")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("bash", "-c", "./build_plugins.sh")
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

	var diags hcl.Diagnostics
	inventoryFile := &inventory.InventoryFile{
		Path:    "test_inventory.hcl",
		Content: []byte(inventoryContent),
	}
	inv, diags = inventory.ParseInventoryFiles(inventoryFile)
	if diags.HasErrors() {
		panic(diags.Error())
	}

	moduleRegistry = module.NewRegistry()
	moduleRegistry.RegisterCoreModules()
	moduleRegistry.RegisterPluginModules()

	privateKeyContent, _ = os.ReadFile("id_rsa")

	code := m.Run()

	cmd = exec.Command("bash", "-c", "./stop.sh")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(code)
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
