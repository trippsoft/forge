// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/trippsoft/forge/pkg/core"
	"github.com/zclconf/go-cty/cty"
)

type expectedHost struct {
	name          string
	transportType core.TransportType
	vars          map[string]cty.Value
}

type expectedGroup struct {
	name  string
	hosts []string // Host names in this group
}

func setupPrivateKey(t *testing.T) string {
	t.Helper()

	var dstPrivateKey string

	if runtime.GOOS == "windows" {
		os.Mkdir("C:\\temp", 0755)
		dstPrivateKey = "C:\\temp\\test_ssh_key"
	} else {
		dstPrivateKey = "/tmp/test_ssh_key"
	}

	srcPrivateKey := "test_ssh_key"

	srcFile, err := os.Open(srcPrivateKey)
	if err != nil {
		t.Fatalf("Failed to open source SSH key: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create temp SSH key: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		t.Fatalf("Failed to copy SSH key: %v", err)
	}

	err = os.Chmod(dstPrivateKey, 0600)
	if err != nil {
		t.Fatalf("Failed to set permissions on temp SSH key: %v", err)
	}

	return dstPrivateKey
}

func setupKnownHosts(t *testing.T) string {
	t.Helper()

	var knownHostsPath string
	if runtime.GOOS == "windows" {
		os.Mkdir("C:\\temp", 0755)
		knownHostsPath = "C:\\temp\\known_hosts"
	} else {
		knownHostsPath = "/tmp/known_hosts"
	}

	knownHostsFile, err := os.Create(knownHostsPath)
	if err != nil {
		t.Fatalf("Failed to create temp known_hosts file: %v", err)
	}
	defer knownHostsFile.Close()

	return knownHostsPath
}

func createExpectedTargets(t *testing.T, expectedHosts []expectedHost, expectedGroups []expectedGroup) []expectedGroup {
	t.Helper()

	expectedTargets := slices.Clone(expectedGroups)

	// Add 'all' group
	allHostNames := make([]string, len(expectedHosts))
	for i, host := range expectedHosts {
		allHostNames[i] = host.name
	}
	expectedTargets = append(expectedTargets, expectedGroup{
		name:  "all",
		hosts: allHostNames,
	})

	// Add individual host groups
	for _, host := range expectedHosts {
		expectedTargets = append(expectedTargets, expectedGroup{
			name:  host.name,
			hosts: []string{host.name},
		})
	}

	return expectedTargets
}

func verifyHosts(t *testing.T, inventory *core.Inventory, expectedHosts []expectedHost) {
	hosts := inventory.Hosts()
	if len(hosts) != len(expectedHosts) {
		t.Errorf("Expected %d hosts, got %d.", len(expectedHosts), len(hosts))
	}

	for _, expected := range expectedHosts {
		host, exists := hosts[expected.name]
		if !exists {
			t.Errorf("Expected host %q not found in inventory", expected.name)
			continue
		}

		verifyHost(t, host, expected)
	}
}

func verifyHost(t *testing.T, host *core.Host, expected expectedHost) {
	t.Helper()

	if host.Name() != expected.name {
		t.Errorf("Expected host name %q, got %q", expected.name, host.Name())
	}

	if host.Transport() == nil {
		t.Errorf("Expected host %q to have a transport", expected.name)
	} else if host.Transport().Type() != expected.transportType {
		t.Errorf(
			"Expected host %q transport type %q, got %q",
			expected.name,
			expected.transportType,
			host.Transport().Type(),
		)
	}

	vars := host.Vars()
	if vars == nil {
		t.Errorf("Expected host %q to have variables", expected.name)
		return
	}

	if len(vars) != len(expected.vars) {
		t.Errorf("Expected host %q to have %d variables, got %d", expected.name, len(expected.vars), len(vars))
		return
	}

	for name, expectedValue := range expected.vars {
		value, exists := vars[name]
		if !exists {
			t.Errorf("Expected host %q variable %q not found", expected.name, name)
			continue
		}

		if !value.Type().Equals(expectedValue.Type()) {
			t.Errorf(
				"Expected host %q variable %q to be of type %s, got %s",
				expected.name,
				name,
				expectedValue.Type().FriendlyName(),
				value.Type().FriendlyName())
			continue
		}

		if !value.Equals(expectedValue).True() {
			t.Errorf(
				"Expected host %q variable %q to have value %q, got %q",
				expected.name,
				name,
				expectedValue.AsString(),
				value.AsString(),
			)
		}
	}
}

func verifyGroups(t *testing.T, inventory *core.Inventory, expectedGroups []expectedGroup) {
	t.Helper()

	groups := inventory.Groups()
	if len(groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, got %d.", len(expectedGroups), len(groups))
	}

	for _, expected := range expectedGroups {
		hosts, exists := groups[expected.name]
		if !exists {
			t.Errorf("Expected group %q not found in inventory", expected.name)
			continue
		}

		if len(hosts) != len(expected.hosts) {
			t.Errorf("Expected group %q to have %d hosts, got %d", expected.name, len(expected.hosts), len(hosts))
			continue
		}

		for _, hostName := range expected.hosts {
			found := false
			for _, host := range hosts {
				if host.Name() == hostName {
					found = true
					inventoryHost, exists := inventory.Host(hostName)
					if !exists {
						t.Errorf("Host %q in group %q not found in inventory", hostName, expected.name)
						break
					}

					if inventoryHost != host {
						t.Errorf("Host %q in group %q does not match inventory", hostName, expected.name)
					}
					break
				}
			}
			if !found {
				t.Errorf("Host %q not found in group %q", hostName, expected.name)
			}
		}
	}
}

func verifyTargets(t *testing.T, inventory *core.Inventory, expectedTargets []expectedGroup) {
	t.Helper()

	targets := inventory.Targets()
	if len(targets) != len(expectedTargets) {
		t.Errorf("Expected %d groups, got %d.", len(expectedTargets), len(targets))
	}

	for _, expected := range expectedTargets {
		hosts, exists := targets[expected.name]
		if !exists {
			t.Errorf("Expected group %q not found in inventory", expected.name)
			continue
		}

		if len(hosts) != len(expected.hosts) {
			t.Errorf("Expected group %q to have %d hosts, got %d", expected.name, len(expected.hosts), len(hosts))
			continue
		}

		for _, hostName := range expected.hosts {
			found := false
			for _, host := range hosts {
				if host.Name() == hostName {
					found = true
					inventoryHost, exists := inventory.Host(hostName)
					if !exists {
						t.Errorf("Host %q in group %q not found in inventory", hostName, expected.name)
						break
					}

					if inventoryHost != host {
						t.Errorf("Host %q in group %q does not match inventory", hostName, expected.name)
					}
					break
				}
			}
			if !found {
				t.Errorf("Host %q not found in group %q", hostName, expected.name)
			}
		}
	}
}

func verifyDiagnostics(t *testing.T, expected hcl.Diagnostics, actual hcl.Diagnostics) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Errorf("Expected %d diagnostics, got %d", len(expected), len(actual))
		return
	}

	for _, expectedDiag := range expected {
		found := false
		for _, actualDiag := range actual {
			if expectedDiag.Summary == actualDiag.Summary &&
				expectedDiag.Detail == actualDiag.Detail &&
				expectedDiag.Severity == actualDiag.Severity {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected diagnostic not found: %+v", expectedDiag)
		}
	}
}

func TestSimpleParsing(t *testing.T) {
	path := filepath.Join("corpus", "simple")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment": cty.StringVal("test"),
				"domain":      cty.StringVal("example.com"),
				"role":        cty.StringVal("web"),
				"port":        cty.NumberIntVal(8080),
				"ip":          cty.StringVal("10.0.1.10"),
				"hostname":    cty.StringVal("web1.example.com"),
			},
		},
		{
			name:          "web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment": cty.StringVal("test"),
				"domain":      cty.StringVal("example.com"),
				"role":        cty.StringVal("web"),
				"port":        cty.NumberIntVal(8080),
				"ip":          cty.StringVal("10.0.1.11"),
				"hostname":    cty.StringVal("web2.example.com"),
			},
		},
		{
			name:          "db1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment": cty.StringVal("test"),
				"domain":      cty.StringVal("example.com"),
				"role":        cty.StringVal("database"),
				"port":        cty.NumberIntVal(5432),
				"ip":          cty.StringVal("10.0.2.10"),
				"hostname":    cty.StringVal("db1.example.com"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "webservers",
			hosts: []string{"web1", "web2"},
		},
		{
			name:  "databases",
			hosts: []string{"db1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestParentHierarchyParsing(t *testing.T) {
	knownHostsPath := setupKnownHosts(t)
	defer os.Remove(knownHostsPath)

	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join("corpus", "parent-hierarchy", "win")
	} else {
		path = filepath.Join("corpus", "parent-hierarchy", "nonwin")
	}

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("frontend"),
				"load_balanced":      cty.BoolVal(true),
				"app_port":           cty.NumberIntVal(8080),
				"ip":                 cty.StringVal("10.0.1.10"),
				"tier":               cty.StringVal("primary"),
			},
		},
		{
			name:          "web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("frontend"),
				"load_balanced":      cty.BoolVal(true),
				"app_port":           cty.NumberIntVal(8080),
				"ip":                 cty.StringVal("10.0.1.11"),
				"tier":               cty.StringVal("secondary"),
			},
		},
		{
			name:          "web3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("frontend"),
				"load_balanced":      cty.BoolVal(true),
				"app_port":           cty.NumberIntVal(8080),
				"ip":                 cty.StringVal("10.0.1.12"),
				"tier":               cty.StringVal("secondary"),
			},
		},
		{
			name:          "api1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("backend"),
				"api_version":        cty.StringVal("v2"),
				"app_port":           cty.NumberIntVal(9000),
				"ip":                 cty.StringVal("10.0.2.10"),
				"tier":               cty.StringVal("primary"),
			},
		},
		{
			name:          "api2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("backend"),
				"api_version":        cty.StringVal("v2"),
				"app_port":           cty.NumberIntVal(9000),
				"ip":                 cty.StringVal("10.0.2.11"),
				"tier":               cty.StringVal("secondary"),
			},
		},
		{
			name:          "cdn1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":        cty.StringVal("production"),
				"monitoring_enabled": cty.BoolVal(true),
				"backup_enabled":     cty.BoolVal(true),
				"log_level":          cty.StringVal("info"),
				"managed":            cty.BoolVal(true),
				"role":               cty.StringVal("frontend"),
				"load_balanced":      cty.BoolVal(true),
				"app_port":           cty.NumberIntVal(8080),
				"cache_enabled":      cty.BoolVal(true),
				"edge_locations": cty.TupleVal([]cty.Value{
					cty.StringVal("us-east"),
					cty.StringVal("us-west"),
					cty.StringVal("eu-west"),
				}),
				"ip":   cty.StringVal("10.0.3.10"),
				"tier": cty.StringVal("edge"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "infrastructure",
			hosts: []string{"web1", "web2", "web3", "api1", "api2", "cdn1"},
		},
		{
			name:  "frontend",
			hosts: []string{"web1", "web2", "web3", "cdn1"},
		},
		{
			name:  "backend",
			hosts: []string{"api1", "api2"},
		},
		{
			name:  "cdn",
			hosts: []string{"cdn1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestVariableInterpolationParsing(t *testing.T) {
	path := filepath.Join("corpus", "variable-interpolation")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":     cty.StringVal("test"),
				"domain":          cty.StringVal("example.com"),
				"datacenter":      cty.StringVal("us-east-1"),
				"internal_domain": cty.StringVal("internal.example.com"),
				"external_domain": cty.StringVal("external.example.com"),
				"network_prefix":  cty.StringVal("10.0"),
				"web_subnet":      cty.StringVal("10.0.1"),
				"api_subnet":      cty.StringVal("10.0.2"),
				"db_subnet":       cty.StringVal("10.0.3"),
				"app_name":        cty.StringVal("myapp"),
				"app_version":     cty.StringVal("1.2.3"),
				"app_image":       cty.StringVal("myapp:1.2.3"),
				"base_port":       cty.NumberIntVal(8000),
				"web_port":        cty.NumberIntVal(8080),
				"api_port":        cty.NumberIntVal(8090),
				"availability_zones": cty.TupleVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
				"log_levels": cty.TupleVal([]cty.Value{
					cty.StringVal("debug"),
					cty.StringVal("info"),
					cty.StringVal("warn"),
					cty.StringVal("error"),
				}),
				"role":         cty.StringVal("web"),
				"cluster_name": cty.StringVal("myapp-web-test"),
				"service_url":  cty.StringVal("https://web.external.example.com:8080"),
				"internal_url": cty.StringVal("http://web.internal.example.com:8080"),
				"host_id":      cty.NumberIntVal(1),
				"ip":           cty.StringVal("10.0.1.10"),
				"hostname":     cty.StringVal("web1.internal.example.com"),
				"fqdn":         cty.StringVal("web1.internal.example.com"),
				"zone":         cty.StringVal("a"),
			},
		},
		{
			name:          "web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":     cty.StringVal("test"),
				"domain":          cty.StringVal("example.com"),
				"datacenter":      cty.StringVal("us-east-1"),
				"internal_domain": cty.StringVal("internal.example.com"),
				"external_domain": cty.StringVal("external.example.com"),
				"network_prefix":  cty.StringVal("10.0"),
				"web_subnet":      cty.StringVal("10.0.1"),
				"api_subnet":      cty.StringVal("10.0.2"),
				"db_subnet":       cty.StringVal("10.0.3"),
				"app_name":        cty.StringVal("myapp"),
				"app_version":     cty.StringVal("1.2.3"),
				"app_image":       cty.StringVal("myapp:1.2.3"),
				"base_port":       cty.NumberIntVal(8000),
				"web_port":        cty.NumberIntVal(8080),
				"api_port":        cty.NumberIntVal(8090),
				"availability_zones": cty.TupleVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
				"log_levels": cty.TupleVal([]cty.Value{
					cty.StringVal("debug"),
					cty.StringVal("info"),
					cty.StringVal("warn"),
					cty.StringVal("error"),
				}),
				"role":         cty.StringVal("web"),
				"cluster_name": cty.StringVal("myapp-web-test"),
				"service_url":  cty.StringVal("https://web.external.example.com:8080"),
				"internal_url": cty.StringVal("http://web.internal.example.com:8080"),
				"host_id":      cty.NumberIntVal(2),
				"ip":           cty.StringVal("10.0.1.11"),
				"hostname":     cty.StringVal("web2.internal.example.com"),
				"fqdn":         cty.StringVal("web2.internal.example.com"),
				"zone":         cty.StringVal("b"),
			},
		},
		{
			name:          "db1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"environment":     cty.StringVal("test"),
				"domain":          cty.StringVal("example.com"),
				"datacenter":      cty.StringVal("us-east-1"),
				"internal_domain": cty.StringVal("internal.example.com"),
				"external_domain": cty.StringVal("external.example.com"),
				"network_prefix":  cty.StringVal("10.0"),
				"web_subnet":      cty.StringVal("10.0.1"),
				"api_subnet":      cty.StringVal("10.0.2"),
				"db_subnet":       cty.StringVal("10.0.3"),
				"app_name":        cty.StringVal("myapp"),
				"app_version":     cty.StringVal("1.2.3"),
				"app_image":       cty.StringVal("myapp:1.2.3"),
				"base_port":       cty.NumberIntVal(8000),
				"web_port":        cty.NumberIntVal(8080),
				"api_port":        cty.NumberIntVal(8090),
				"availability_zones": cty.TupleVal([]cty.Value{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				}),
				"log_levels": cty.TupleVal([]cty.Value{
					cty.StringVal("debug"),
					cty.StringVal("info"),
					cty.StringVal("warn"),
					cty.StringVal("error"),
				}),
				"role":         cty.StringVal("primary"),
				"cluster_name": cty.StringVal("myapp-db-test"),
				"internal_url": cty.StringVal("postgres://db.internal.example.com:5432"),
				"host_id":      cty.NumberIntVal(1),
				"ip":           cty.StringVal("10.0.3.10"),
				"hostname":     cty.StringVal("db1.internal.example.com"),
				"fqdn":         cty.StringVal("db1.internal.example.com"),
				"zone":         cty.StringVal("a"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "web",
			hosts: []string{"web1", "web2"},
		},
		{
			name:  "database",
			hosts: []string{"db1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestTransportInheritanceParsing(t *testing.T) {
	privateKeyPath := setupPrivateKey(t)
	defer os.Remove(privateKeyPath)

	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join("corpus", "transport-inheritance", "win")
	} else {
		path = filepath.Join("corpus", "transport-inheritance", "nonwin")
	}

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "secure1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"ssh_key_path":        cty.StringVal(privateKeyPath),
				"bastion_host":        cty.StringVal("bastion.example.com"),
				"admin_user":          cty.StringVal("admin"),
				"security_level":      cty.StringVal("high"),
				"compliance_required": cty.BoolVal(true),
				"ip":                  cty.StringVal("10.0.1.10"),
				"role":                cty.StringVal("security"),
			},
		},
		{
			name:          "secure2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"ssh_key_path":        cty.StringVal(privateKeyPath),
				"bastion_host":        cty.StringVal("bastion.example.com"),
				"admin_user":          cty.StringVal("admin"),
				"security_level":      cty.StringVal("high"),
				"compliance_required": cty.BoolVal(true),
				"ip":                  cty.StringVal("10.0.1.11"),
				"role":                cty.StringVal("security"),
			},
		},
		{
			name:          "web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"ssh_key_path":        cty.StringVal(privateKeyPath),
				"bastion_host":        cty.StringVal("bastion.example.com"),
				"admin_user":          cty.StringVal("admin"),
				"security_level":      cty.StringVal("standard"),
				"compliance_required": cty.BoolVal(false),
				"ip":                  cty.StringVal("10.0.2.10"),
				"role":                cty.StringVal("web"),
			},
		},
		{
			name:          "web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"ssh_key_path":        cty.StringVal(privateKeyPath),
				"bastion_host":        cty.StringVal("bastion.example.com"),
				"admin_user":          cty.StringVal("admin"),
				"security_level":      cty.StringVal("standard"),
				"compliance_required": cty.BoolVal(false),
				"ip":                  cty.StringVal("10.0.2.11"),
				"role":                cty.StringVal("web"),
			},
		},
		{
			name:          "admin1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"ssh_key_path": cty.StringVal(privateKeyPath),
				"bastion_host": cty.StringVal("bastion.example.com"),
				"admin_user":   cty.StringVal("admin"),
				"admin_access": cty.BoolVal(true),
				"ip":           cty.StringVal("10.0.3.10"),
				"role":         cty.StringVal("admin"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "secure_servers",
			hosts: []string{"secure1", "secure2"},
		},
		{
			name:  "standard_servers",
			hosts: []string{"web1", "web2"},
		},
		{
			name:  "admin_servers",
			hosts: []string{"admin1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestMultiEnvironmentParsing(t *testing.T) {
	knownHostsPath := setupKnownHosts(t)
	defer os.Remove(knownHostsPath)

	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join("corpus", "multi-environment", "win")
	} else {
		path = filepath.Join("corpus", "multi-environment", "nonwin")
	}

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("Expected 3 inventory files, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "prod-web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("web"),
				"domain":              cty.StringVal("www.acme.com"),
				"app_port":            cty.NumberIntVal(8080),
				"ip":                  cty.StringVal("10.1.1.10"),
				"tier":                cty.StringVal("primary"),
			},
		},
		{
			name:          "prod-web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("web"),
				"domain":              cty.StringVal("www.acme.com"),
				"app_port":            cty.NumberIntVal(8080),
				"ip":                  cty.StringVal("10.1.1.11"),
				"tier":                cty.StringVal("secondary"),
			},
		},
		{
			name:          "prod-web3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("web"),
				"domain":              cty.StringVal("www.acme.com"),
				"app_port":            cty.NumberIntVal(8080),
				"ip":                  cty.StringVal("10.1.1.12"),
				"tier":                cty.StringVal("secondary"),
			},
		},
		{
			name:          "prod-api1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("api"),
				"domain":              cty.StringVal("api.acme.com"),
				"app_port":            cty.NumberIntVal(9000),
				"ip":                  cty.StringVal("10.1.2.10"),
				"tier":                cty.StringVal("primary"),
			},
		},
		{
			name:          "prod-api2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("api"),
				"domain":              cty.StringVal("api.acme.com"),
				"app_port":            cty.NumberIntVal(9000),
				"ip":                  cty.StringVal("10.1.2.11"),
				"tier":                cty.StringVal("secondary"),
			},
		},
		{
			name:          "staging-web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("web"),
				"domain":              cty.StringVal("staging.acme.com"),
				"app_port":            cty.NumberIntVal(8080),
				"ip":                  cty.StringVal("10.2.1.10"),
				"tier":                cty.StringVal("single"),
			},
		},
		{
			name:          "staging-api1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"company":             cty.StringVal("acme"),
				"base_domain":         cty.StringVal("acme.com"),
				"ssh_user":            cty.StringVal("deploy"),
				"prod_environment":    cty.StringVal("production"),
				"prod_log_level":      cty.StringVal("warn"),
				"prod_replicas":       cty.NumberIntVal(3),
				"staging_environment": cty.StringVal("staging"),
				"staging_log_level":   cty.StringVal("debug"),
				"staging_replicas":    cty.NumberIntVal(1),
				"role":                cty.StringVal("api"),
				"domain":              cty.StringVal("staging-api.acme.com"),
				"app_port":            cty.NumberIntVal(9000),
				"ip":                  cty.StringVal("10.2.2.10"),
				"tier":                cty.StringVal("single"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "prod_web",
			hosts: []string{"prod-web1", "prod-web2", "prod-web3"},
		},
		{
			name:  "prod_api",
			hosts: []string{"prod-api1", "prod-api2"},
		},
		{
			name:  "staging_web",
			hosts: []string{"staging-web1"},
		},
		{
			name:  "staging_api",
			hosts: []string{"staging-api1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestMultiEnvironmentParsing_Partial(t *testing.T) {
	knownHostsPath := setupKnownHosts(t)
	defer os.Remove(knownHostsPath)

	tests := []struct {
		name           string
		paths          []string
		expectedHosts  []expectedHost
		expectedGroups []expectedGroup
	}{
		{
			name:  "Production Only",
			paths: []string{"globals.hcl", "production.hcl"},
			expectedHosts: []expectedHost{
				{
					name:          "prod-web1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":          cty.StringVal("acme"),
						"base_domain":      cty.StringVal("acme.com"),
						"ssh_user":         cty.StringVal("deploy"),
						"prod_environment": cty.StringVal("production"),
						"prod_log_level":   cty.StringVal("warn"),
						"prod_replicas":    cty.NumberIntVal(3),
						"role":             cty.StringVal("web"),
						"domain":           cty.StringVal("www.acme.com"),
						"app_port":         cty.NumberIntVal(8080),
						"ip":               cty.StringVal("10.1.1.10"),
						"tier":             cty.StringVal("primary"),
					},
				},
				{
					name:          "prod-web2",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":          cty.StringVal("acme"),
						"base_domain":      cty.StringVal("acme.com"),
						"ssh_user":         cty.StringVal("deploy"),
						"prod_environment": cty.StringVal("production"),
						"prod_log_level":   cty.StringVal("warn"),
						"prod_replicas":    cty.NumberIntVal(3),
						"role":             cty.StringVal("web"),
						"domain":           cty.StringVal("www.acme.com"),
						"app_port":         cty.NumberIntVal(8080),
						"ip":               cty.StringVal("10.1.1.11"),
						"tier":             cty.StringVal("secondary"),
					},
				},
				{
					name:          "prod-web3",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":          cty.StringVal("acme"),
						"base_domain":      cty.StringVal("acme.com"),
						"ssh_user":         cty.StringVal("deploy"),
						"prod_environment": cty.StringVal("production"),
						"prod_log_level":   cty.StringVal("warn"),
						"prod_replicas":    cty.NumberIntVal(3),
						"role":             cty.StringVal("web"),
						"domain":           cty.StringVal("www.acme.com"),
						"app_port":         cty.NumberIntVal(8080),
						"ip":               cty.StringVal("10.1.1.12"),
						"tier":             cty.StringVal("secondary"),
					},
				},
				{
					name:          "prod-api1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":          cty.StringVal("acme"),
						"base_domain":      cty.StringVal("acme.com"),
						"ssh_user":         cty.StringVal("deploy"),
						"prod_environment": cty.StringVal("production"),
						"prod_log_level":   cty.StringVal("warn"),
						"prod_replicas":    cty.NumberIntVal(3),
						"role":             cty.StringVal("api"),
						"domain":           cty.StringVal("api.acme.com"),
						"app_port":         cty.NumberIntVal(9000),
						"ip":               cty.StringVal("10.1.2.10"),
						"tier":             cty.StringVal("primary"),
					},
				},
				{
					name:          "prod-api2",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":          cty.StringVal("acme"),
						"base_domain":      cty.StringVal("acme.com"),
						"ssh_user":         cty.StringVal("deploy"),
						"prod_environment": cty.StringVal("production"),
						"prod_log_level":   cty.StringVal("warn"),
						"prod_replicas":    cty.NumberIntVal(3),
						"role":             cty.StringVal("api"),
						"domain":           cty.StringVal("api.acme.com"),
						"app_port":         cty.NumberIntVal(9000),
						"ip":               cty.StringVal("10.1.2.11"),
						"tier":             cty.StringVal("secondary"),
					},
				},
			},
			expectedGroups: []expectedGroup{
				{
					name:  "prod_web",
					hosts: []string{"prod-web1", "prod-web2", "prod-web3"},
				},
				{
					name:  "prod_api",
					hosts: []string{"prod-api1", "prod-api2"},
				},
			},
		},
		{
			name:  "Staging Only",
			paths: []string{"globals.hcl", "staging.hcl"},
			expectedHosts: []expectedHost{
				{
					name:          "staging-web1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":             cty.StringVal("acme"),
						"base_domain":         cty.StringVal("acme.com"),
						"ssh_user":            cty.StringVal("deploy"),
						"staging_environment": cty.StringVal("staging"),
						"staging_log_level":   cty.StringVal("debug"),
						"staging_replicas":    cty.NumberIntVal(1),
						"role":                cty.StringVal("web"),
						"domain":              cty.StringVal("staging.acme.com"),
						"app_port":            cty.NumberIntVal(8080),
						"ip":                  cty.StringVal("10.2.1.10"),
						"tier":                cty.StringVal("single"),
					},
				},
				{
					name:          "staging-api1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"company":             cty.StringVal("acme"),
						"base_domain":         cty.StringVal("acme.com"),
						"ssh_user":            cty.StringVal("deploy"),
						"staging_environment": cty.StringVal("staging"),
						"staging_log_level":   cty.StringVal("debug"),
						"staging_replicas":    cty.NumberIntVal(1),
						"role":                cty.StringVal("api"),
						"domain":              cty.StringVal("staging-api.acme.com"),
						"app_port":            cty.NumberIntVal(9000),
						"ip":                  cty.StringVal("10.2.2.10"),
						"tier":                cty.StringVal("single"),
					},
				},
			},
			expectedGroups: []expectedGroup{
				{
					name:  "staging_web",
					hosts: []string{"staging-web1"},
				},
				{
					name:  "staging_api",
					hosts: []string{"staging-api1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fullPaths := make([]string, len(tt.paths))
			for i, p := range tt.paths {
				if runtime.GOOS == "windows" {
					fullPaths[i] = filepath.Join("corpus", "multi-environment", "win", p)
				} else {
					fullPaths[i] = filepath.Join("corpus", "multi-environment", "nonwin", p)
				}
			}

			files, err := core.DiscoverInventoryFiles(fullPaths...)
			if err != nil {
				t.Fatalf("Failed to discover inventory files: %v", err)
			}

			if len(files) != len(tt.paths) {
				t.Fatalf("Expected %d inventory files, got %d", len(tt.paths), len(files))
			}

			inventory, diags := core.ParseInventoryFiles(files)
			if diags.HasErrors() {
				t.Fatalf("Failed to parse inventory: %s", diags.Error())
			}

			if len(diags) > 0 {
				t.Errorf("Expected no diagnostics, got: %v", diags)
			}

			if inventory == nil {
				t.Fatal("Inventory should not be nil for valid configuration")
			}

			verifyHosts(t, inventory, tt.expectedHosts)
			verifyGroups(t, inventory, tt.expectedGroups)

			expectedTargets := createExpectedTargets(t, tt.expectedHosts, tt.expectedGroups)

			verifyTargets(t, inventory, expectedTargets)
		})
	}
}

func TestComplexVariableParsing(t *testing.T) {
	path := filepath.Join("corpus", "complex-variables")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "web1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"app_version":        cty.StringVal("2.1.0"),
				"health_check_port":  cty.NumberIntVal(8080),
				"service_url":        cty.StringVal("https://app.internal.company.com"),
				"role":               cty.StringVal("frontend"),
				"load_balancer_pool": cty.StringVal("frontend-staging"),
				"replicas":           cty.NumberIntVal(3),
				"ip":                 cty.StringVal("10.0.1.10"),
				"hostname":           cty.StringVal("web1.internal.company.com"),
				"server_id":          cty.NumberIntVal(1),
				"memory_limit":       cty.StringVal("2GB"),
			},
		},
		{
			name:          "web2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"app_version":        cty.StringVal("2.1.0"),
				"health_check_port":  cty.NumberIntVal(8080),
				"service_url":        cty.StringVal("https://app.internal.company.com"),
				"role":               cty.StringVal("frontend"),
				"load_balancer_pool": cty.StringVal("frontend-staging"),
				"replicas":           cty.NumberIntVal(3),
				"ip":                 cty.StringVal("10.0.1.11"),
				"hostname":           cty.StringVal("web2.internal.company.com"),
				"server_id":          cty.NumberIntVal(2),
				"memory_limit":       cty.StringVal("2GB"),
			},
		},
		{
			name:          "web3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"app_version":        cty.StringVal("2.1.0"),
				"health_check_port":  cty.NumberIntVal(8080),
				"service_url":        cty.StringVal("https://app.internal.company.com"),
				"role":               cty.StringVal("frontend"),
				"load_balancer_pool": cty.StringVal("frontend-staging"),
				"replicas":           cty.NumberIntVal(3),
				"ip":                 cty.StringVal("10.0.1.12"),
				"hostname":           cty.StringVal("web3.internal.company.com"),
				"server_id":          cty.NumberIntVal(3),
				"memory_limit":       cty.StringVal("4GB"),
			},
		},
		{
			name:          "api1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"app_version":        cty.StringVal("2.1.0"),
				"health_check_port":  cty.NumberIntVal(8080),
				"service_url":        cty.StringVal("https://app.internal.company.com"),
				"role":               cty.StringVal("backend"),
				"database_url":       cty.StringVal("postgres://db.internal.company.com:5432/app"),
				"api_prefix":         cty.StringVal("/api/v2"),
				"ip":                 cty.StringVal("10.0.2.10"),
				"hostname":           cty.StringVal("api1.internal.company.com"),
				"instance_type":      cty.StringVal("primary"),
			},
		},
		{
			name:          "api2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"app_version":        cty.StringVal("2.1.0"),
				"health_check_port":  cty.NumberIntVal(8080),
				"service_url":        cty.StringVal("https://app.internal.company.com"),
				"role":               cty.StringVal("backend"),
				"database_url":       cty.StringVal("postgres://db.internal.company.com:5432/app"),
				"api_prefix":         cty.StringVal("/api/v2"),
				"ip":                 cty.StringVal("10.0.2.11"),
				"hostname":           cty.StringVal("api2.internal.company.com"),
				"instance_type":      cty.StringVal("secondary"),
			},
		},
		{
			name:          "db1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"role":               cty.StringVal("database"),
				"backup_enabled":     cty.BoolVal(true),
				"replication_factor": cty.NumberIntVal(2),
				"engine":             cty.StringVal("postgresql"),
				"version":            cty.StringVal("14.9"),
				"port":               cty.NumberIntVal(5432),
				"ip":                 cty.StringVal("10.0.3.10"),
				"hostname":           cty.StringVal("db1.internal.company.com"),
				"is_primary":         cty.BoolVal(true),
				"storage_size":       cty.StringVal("100GB"),
			},
		},
		{
			name:          "db2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"global_env":         cty.StringVal("staging"),
				"base_domain":        cty.StringVal("internal.company.com"),
				"backup_schedule":    cty.StringVal("0 2 * * *"),
				"monitoring":         cty.BoolVal(true),
				"log_retention_days": cty.NumberIntVal(30),
				"environment":        cty.StringVal("staging"),
				"fqdn_suffix":        cty.StringVal("internal.company.com"),
				"role":               cty.StringVal("database"),
				"backup_enabled":     cty.BoolVal(true),
				"replication_factor": cty.NumberIntVal(2),
				"engine":             cty.StringVal("postgresql"),
				"version":            cty.StringVal("14.9"),
				"port":               cty.NumberIntVal(5432),
				"ip":                 cty.StringVal("10.0.3.11"),
				"hostname":           cty.StringVal("db2.internal.company.com"),
				"is_primary":         cty.BoolVal(false),
				"storage_size":       cty.StringVal("100GB"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "base",
			hosts: []string{"web1", "web2", "web3", "api1", "api2", "db1", "db2"},
		},
		{
			name:  "app_tier",
			hosts: []string{"web1", "web2", "web3", "api1", "api2"},
		},
		{
			name:  "frontend",
			hosts: []string{"web1", "web2", "web3"},
		},
		{
			name:  "backend",
			hosts: []string{"api1", "api2"},
		},
		{
			name:  "data_tier",
			hosts: []string{"db1", "db2"},
		},
		{
			name:  "databases",
			hosts: []string{"db1", "db2"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestDeepHierarchyParsing(t *testing.T) {
	path := filepath.Join("corpus", "deep-hierarchy")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "user-svc-1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// platform
				"container_runtime": cty.StringVal("docker"),
				"orchestrator":      cty.StringVal("k8s"),
				// app_services
				"service_mesh":    cty.StringVal("istio"),
				"tracing_enabled": cty.BoolVal(true),
				// user_service
				"tier":         cty.StringVal("microservice"),
				"service_name": cty.StringVal("user-service"),
				"version":      cty.StringVal("v2.3.1"),
				"replicas":     cty.NumberIntVal(3),
				// user-svc-1
				"ip":          cty.StringVal("10.1.1.10"),
				"instance_id": cty.StringVal("i-user-001"),
				"cpu_cores":   cty.NumberIntVal(2),
			},
		},
		{
			name:          "user-svc-2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// platform
				"container_runtime": cty.StringVal("docker"),
				"orchestrator":      cty.StringVal("k8s"),
				// app_services
				"service_mesh":    cty.StringVal("istio"),
				"tracing_enabled": cty.BoolVal(true),
				// user_service
				"tier":         cty.StringVal("microservice"),
				"service_name": cty.StringVal("user-service"),
				"version":      cty.StringVal("v2.3.1"),
				"replicas":     cty.NumberIntVal(3),
				// user-svc-2
				"ip":          cty.StringVal("10.1.1.11"),
				"instance_id": cty.StringVal("i-user-002"),
				"cpu_cores":   cty.NumberIntVal(2),
			},
		},
		{
			name:          "user-svc-3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// platform
				"container_runtime": cty.StringVal("docker"),
				"orchestrator":      cty.StringVal("k8s"),
				// app_services
				"service_mesh":    cty.StringVal("istio"),
				"tracing_enabled": cty.BoolVal(true),
				// user_service
				"tier":         cty.StringVal("microservice"),
				"service_name": cty.StringVal("user-service"),
				"version":      cty.StringVal("v2.3.1"),
				"replicas":     cty.NumberIntVal(3),
				// user-svc-3
				"ip":          cty.StringVal("10.1.1.12"),
				"instance_id": cty.StringVal("i-user-003"),
				"cpu_cores":   cty.NumberIntVal(4),
			},
		},
		{
			name:          "payment-svc-1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// platform
				"container_runtime": cty.StringVal("docker"),
				"orchestrator":      cty.StringVal("k8s"),
				// app_services
				"service_mesh":    cty.StringVal("istio"),
				"tracing_enabled": cty.BoolVal(true),
				// payment_service
				"tier":          cty.StringVal("microservice"),
				"service_name":  cty.StringVal("payment-service"),
				"version":       cty.StringVal("v1.8.2"),
				"replicas":      cty.NumberIntVal(2),
				"pci_compliant": cty.BoolVal(true),
				// payment-svc-1
				"ip":             cty.StringVal("10.1.2.10"),
				"instance_id":    cty.StringVal("i-payment-001"),
				"cpu_cores":      cty.NumberIntVal(4),
				"secure_enclave": cty.BoolVal(true),
			},
		},
		{
			name:          "payment-svc-2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// platform
				"container_runtime": cty.StringVal("docker"),
				"orchestrator":      cty.StringVal("k8s"),
				// app_services
				"service_mesh":    cty.StringVal("istio"),
				"tracing_enabled": cty.BoolVal(true),
				// payment_service
				"tier":          cty.StringVal("microservice"),
				"service_name":  cty.StringVal("payment-service"),
				"version":       cty.StringVal("v1.8.2"),
				"replicas":      cty.NumberIntVal(2),
				"pci_compliant": cty.BoolVal(true),
				// payment-svc-2
				"ip":             cty.StringVal("10.1.2.11"),
				"instance_id":    cty.StringVal("i-payment-002"),
				"cpu_cores":      cty.NumberIntVal(4),
				"secure_enclave": cty.BoolVal(true),
			},
		},
		{
			name:          "analytics-db-1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// data_platform
				"encryption_at_rest":  cty.BoolVal(true),
				"data_classification": cty.StringVal("sensitive"),
				// analytics_db
				"tier":         cty.StringVal("database"),
				"engine":       cty.StringVal("postgresql"),
				"cluster_mode": cty.BoolVal(true),
				// analytics-db-1
				"ip":           cty.StringVal("10.2.1.10"),
				"instance_id":  cty.StringVal("i-analytics-db-001"),
				"storage_type": cty.StringVal("ssd"),
				"storage_size": cty.StringVal("500GB"),
				"is_primary":   cty.BoolVal(true),
			},
		},
		{
			name:          "analytics-db-2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// data_platform
				"encryption_at_rest":  cty.BoolVal(true),
				"data_classification": cty.StringVal("sensitive"),
				// analytics_db
				"tier":         cty.StringVal("database"),
				"engine":       cty.StringVal("postgresql"),
				"cluster_mode": cty.BoolVal(true),
				// analytics-db-2
				"ip":           cty.StringVal("10.2.1.11"),
				"instance_id":  cty.StringVal("i-analytics-db-002"),
				"storage_type": cty.StringVal("ssd"),
				"storage_size": cty.StringVal("500GB"),
				"is_primary":   cty.BoolVal(false),
			},
		},
		{
			name:          "analytics-db-3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				// global
				"organization":     cty.StringVal("acme-corp"),
				"compliance_level": cty.StringVal("high"),
				// foundation
				"security_baseline": cty.StringVal("cis-level-2"),
				"audit_enabled":     cty.BoolVal(true),
				// infrastructure
				"monitoring_enabled": cty.BoolVal(true),
				"backup_retention":   cty.StringVal("30d"),
				// data_platform
				"encryption_at_rest":  cty.BoolVal(true),
				"data_classification": cty.StringVal("sensitive"),
				// analytics_db
				"tier":         cty.StringVal("database"),
				"engine":       cty.StringVal("postgresql"),
				"cluster_mode": cty.BoolVal(true),
				// analytics-db-3
				"ip":           cty.StringVal("10.2.1.12"),
				"instance_id":  cty.StringVal("i-analytics-db-003"),
				"storage_type": cty.StringVal("ssd"),
				"storage_size": cty.StringVal("500GB"),
				"is_primary":   cty.BoolVal(false),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name: "foundation",
			hosts: []string{
				"user-svc-1",
				"user-svc-2",
				"user-svc-3",
				"payment-svc-1",
				"payment-svc-2",
				"analytics-db-1",
				"analytics-db-2",
				"analytics-db-3",
			},
		},
		{
			name: "infrastructure",
			hosts: []string{
				"user-svc-1",
				"user-svc-2",
				"user-svc-3",
				"payment-svc-1",
				"payment-svc-2",
				"analytics-db-1",
				"analytics-db-2",
				"analytics-db-3",
			},
		},
		{
			name: "platform",
			hosts: []string{
				"user-svc-1",
				"user-svc-2",
				"user-svc-3",
				"payment-svc-1",
				"payment-svc-2",
			},
		},
		{
			name: "app_services",
			hosts: []string{
				"user-svc-1",
				"user-svc-2",
				"user-svc-3",
				"payment-svc-1",
				"payment-svc-2",
			},
		},
		{
			name: "user_service",
			hosts: []string{
				"user-svc-1",
				"user-svc-2",
				"user-svc-3",
			},
		},
		{
			name: "payment_service",
			hosts: []string{
				"payment-svc-1",
				"payment-svc-2",
			},
		},
		{
			name: "data_platform",
			hosts: []string{
				"analytics-db-1",
				"analytics-db-2",
				"analytics-db-3",
			},
		},
		{
			name: "analytics_db",
			hosts: []string{
				"analytics-db-1",
				"analytics-db-2",
				"analytics-db-3",
			},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestMultiPathCombinedParsing(t *testing.T) {
	path1 := filepath.Join("corpus", "multi-path", "global")
	path2 := filepath.Join("corpus", "multi-path", "production")
	path3 := filepath.Join("corpus", "multi-path", "staging")

	files, err := core.DiscoverInventoryFiles(path1, path2, path3)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("Expected 3 inventory files, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse inventory: %s", diags.Error())
	}

	if len(diags) > 0 {
		t.Errorf("Expected no diagnostics, got: %v", diags)
	}

	if inventory == nil {
		t.Fatal("Inventory should not be nil for valid configuration")
	}

	expectedHosts := []expectedHost{
		{
			name:          "prod-web-1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"organization":   cty.StringVal("test-corp"),
				"domain":         cty.StringVal("test.local"),
				"ssh_user":       cty.StringVal("deploy"),
				"prod_env":       cty.StringVal("production"),
				"prod_count":     cty.NumberIntVal(3),
				"staging_env":    cty.StringVal("staging"),
				"staging_count":  cty.NumberIntVal(1),
				"tier":           cty.StringVal("infrastructure"),
				"backup_enabled": cty.BoolVal(true),
				"role":           cty.StringVal("web"),
				"environment":    cty.StringVal("production"),
				"instance_count": cty.NumberIntVal(3),
				"ip":             cty.StringVal("10.1.1.10"),
				"hostname":       cty.StringVal("prod-web-1.test.local"),
			},
		},
		{
			name:          "prod-web-2",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"organization":   cty.StringVal("test-corp"),
				"domain":         cty.StringVal("test.local"),
				"ssh_user":       cty.StringVal("deploy"),
				"prod_env":       cty.StringVal("production"),
				"prod_count":     cty.NumberIntVal(3),
				"staging_env":    cty.StringVal("staging"),
				"staging_count":  cty.NumberIntVal(1),
				"tier":           cty.StringVal("infrastructure"),
				"backup_enabled": cty.BoolVal(true),
				"role":           cty.StringVal("web"),
				"environment":    cty.StringVal("production"),
				"instance_count": cty.NumberIntVal(3),
				"ip":             cty.StringVal("10.1.1.11"),
				"hostname":       cty.StringVal("prod-web-2.test.local"),
			},
		},
		{
			name:          "prod-web-3",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"organization":   cty.StringVal("test-corp"),
				"domain":         cty.StringVal("test.local"),
				"ssh_user":       cty.StringVal("deploy"),
				"prod_env":       cty.StringVal("production"),
				"prod_count":     cty.NumberIntVal(3),
				"staging_env":    cty.StringVal("staging"),
				"staging_count":  cty.NumberIntVal(1),
				"tier":           cty.StringVal("infrastructure"),
				"backup_enabled": cty.BoolVal(true),
				"role":           cty.StringVal("web"),
				"environment":    cty.StringVal("production"),
				"instance_count": cty.NumberIntVal(3),
				"ip":             cty.StringVal("10.1.1.12"),
				"hostname":       cty.StringVal("prod-web-3.test.local"),
			},
		},
		{
			name:          "staging-web-1",
			transportType: "ssh",
			vars: map[string]cty.Value{
				"organization":   cty.StringVal("test-corp"),
				"domain":         cty.StringVal("test.local"),
				"ssh_user":       cty.StringVal("deploy"),
				"prod_env":       cty.StringVal("production"),
				"prod_count":     cty.NumberIntVal(3),
				"staging_env":    cty.StringVal("staging"),
				"staging_count":  cty.NumberIntVal(1),
				"tier":           cty.StringVal("infrastructure"),
				"backup_enabled": cty.BoolVal(true),
				"role":           cty.StringVal("web"),
				"environment":    cty.StringVal("staging"),
				"instance_count": cty.NumberIntVal(1),
				"ip":             cty.StringVal("10.2.1.10"),
				"hostname":       cty.StringVal("staging-web-1.test.local"),
			},
		},
	}

	verifyHosts(t, inventory, expectedHosts)

	expectedGroups := []expectedGroup{
		{
			name:  "infrastructure",
			hosts: []string{"prod-web-1", "prod-web-2", "prod-web-3", "staging-web-1"},
		},
		{
			name:  "prod_web",
			hosts: []string{"prod-web-1", "prod-web-2", "prod-web-3"},
		},
		{
			name:  "staging_web",
			hosts: []string{"staging-web-1"},
		},
	}

	verifyGroups(t, inventory, expectedGroups)

	expectedTargets := createExpectedTargets(t, expectedHosts, expectedGroups)

	verifyTargets(t, inventory, expectedTargets)
}

func TestMultiPathCombinedParsing_Partial(t *testing.T) {
	tests := []struct {
		name           string
		paths          []string
		expectedHosts  []expectedHost
		expectedGroups []expectedGroup
	}{
		{
			name: "Production only",
			paths: []string{
				filepath.Join("corpus", "multi-path", "global"),
				filepath.Join("corpus", "multi-path", "production"),
			},
			expectedHosts: []expectedHost{
				{
					name:          "prod-web-1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"organization":   cty.StringVal("test-corp"),
						"domain":         cty.StringVal("test.local"),
						"ssh_user":       cty.StringVal("deploy"),
						"prod_env":       cty.StringVal("production"),
						"prod_count":     cty.NumberIntVal(3),
						"tier":           cty.StringVal("infrastructure"),
						"backup_enabled": cty.BoolVal(true),
						"role":           cty.StringVal("web"),
						"environment":    cty.StringVal("production"),
						"instance_count": cty.NumberIntVal(3),
						"ip":             cty.StringVal("10.1.1.10"),
						"hostname":       cty.StringVal("prod-web-1.test.local"),
					},
				},
				{
					name:          "prod-web-2",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"organization":   cty.StringVal("test-corp"),
						"domain":         cty.StringVal("test.local"),
						"ssh_user":       cty.StringVal("deploy"),
						"prod_env":       cty.StringVal("production"),
						"prod_count":     cty.NumberIntVal(3),
						"tier":           cty.StringVal("infrastructure"),
						"backup_enabled": cty.BoolVal(true),
						"role":           cty.StringVal("web"),
						"environment":    cty.StringVal("production"),
						"instance_count": cty.NumberIntVal(3),
						"ip":             cty.StringVal("10.1.1.11"),
						"hostname":       cty.StringVal("prod-web-2.test.local"),
					},
				},
				{
					name:          "prod-web-3",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"organization":   cty.StringVal("test-corp"),
						"domain":         cty.StringVal("test.local"),
						"ssh_user":       cty.StringVal("deploy"),
						"prod_env":       cty.StringVal("production"),
						"prod_count":     cty.NumberIntVal(3),
						"tier":           cty.StringVal("infrastructure"),
						"backup_enabled": cty.BoolVal(true),
						"role":           cty.StringVal("web"),
						"environment":    cty.StringVal("production"),
						"instance_count": cty.NumberIntVal(3),
						"ip":             cty.StringVal("10.1.1.12"),
						"hostname":       cty.StringVal("prod-web-3.test.local"),
					},
				},
			},
			expectedGroups: []expectedGroup{
				{
					name:  "infrastructure",
					hosts: []string{"prod-web-1", "prod-web-2", "prod-web-3"},
				},
				{
					name:  "prod_web",
					hosts: []string{"prod-web-1", "prod-web-2", "prod-web-3"},
				},
			},
		},
		{
			name: "Staging only",
			paths: []string{
				filepath.Join("corpus", "multi-path", "global"),
				filepath.Join("corpus", "multi-path", "staging"),
			},
			expectedHosts: []expectedHost{
				{
					name:          "staging-web-1",
					transportType: "ssh",
					vars: map[string]cty.Value{
						"organization":   cty.StringVal("test-corp"),
						"domain":         cty.StringVal("test.local"),
						"ssh_user":       cty.StringVal("deploy"),
						"staging_env":    cty.StringVal("staging"),
						"staging_count":  cty.NumberIntVal(1),
						"tier":           cty.StringVal("infrastructure"),
						"backup_enabled": cty.BoolVal(true),
						"role":           cty.StringVal("web"),
						"environment":    cty.StringVal("staging"),
						"instance_count": cty.NumberIntVal(1),
						"ip":             cty.StringVal("10.2.1.10"),
						"hostname":       cty.StringVal("staging-web-1.test.local"),
					},
				},
			},
			expectedGroups: []expectedGroup{
				{
					name:  "infrastructure",
					hosts: []string{"staging-web-1"},
				},
				{
					name:  "staging_web",
					hosts: []string{"staging-web-1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			files, err := core.DiscoverInventoryFiles(tt.paths...)
			if err != nil {
				t.Fatalf("Failed to discover inventory files: %v", err)
			}

			if len(files) != len(tt.paths) {
				t.Fatalf("Expected %d inventory files, got %d", len(tt.paths), len(files))
			}

			inventory, diags := core.ParseInventoryFiles(files)
			if diags.HasErrors() {
				t.Fatalf("Failed to parse inventory: %s", diags.Error())
			}

			if len(diags) > 0 {
				t.Errorf("Expected no diagnostics, got: %v", diags)
			}

			if inventory == nil {
				t.Fatal("Inventory should not be nil for valid configuration")
			}

			verifyHosts(t, inventory, tt.expectedHosts)
			verifyGroups(t, inventory, tt.expectedGroups)

			expectedTargets := createExpectedTargets(t, tt.expectedHosts, tt.expectedGroups)

			verifyTargets(t, inventory, expectedTargets)
		})
	}
}

func TestCircularParentReferenceParsing(t *testing.T) {
	path := filepath.Join("corpus", "error-cases", "circular-parent.hcl")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if !diags.HasErrors() {
		t.Fatal("Expected parsing to fail due to circular parent reference")
	}

	expectedDiags := hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular group reference",
			Detail:   "The group 'group_a' has a circular reference.",
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular group reference",
			Detail:   "The group 'group_b' has a circular reference.",
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular group reference",
			Detail:   "The group 'group_c' has a circular reference.",
		},
	}

	verifyDiagnostics(t, expectedDiags, diags)

	if inventory != nil {
		t.Fatal("Inventory should be nil for invalid configuration")
	}
}

func TestInvalidParentReferenceParsing(t *testing.T) {
	path := filepath.Join("corpus", "error-cases", "invalid-parent.hcl")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if !diags.HasErrors() {
		t.Fatal("Expected parsing to fail due to invalid parent reference")
	}

	expectedDiags := hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid parent group",
			Detail:   "The parent group 'nonexistent_group' does not exist.",
		},
	}

	verifyDiagnostics(t, expectedDiags, diags)

	if inventory != nil {
		t.Fatal("Inventory should be nil for invalid configuration")
	}
}

func TestReservedGroupNameParsing(t *testing.T) {
	path := filepath.Join("corpus", "error-cases", "reserved-name.hcl")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if !diags.HasErrors() {
		t.Fatal("Expected parsing to fail due to reserved group name")
	}

	expectedDiags := hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid group name",
			Detail:   "The group name 'all' is reserved and cannot be used.",
		},
	}

	verifyDiagnostics(t, expectedDiags, diags)

	if inventory != nil {
		t.Fatal("Inventory should be nil for invalid configuration")
	}
}

func TestCircularVariablesParsing(t *testing.T) {
	path := filepath.Join("corpus", "error-cases", "circular-variables.hcl")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if !diags.HasErrors() {
		t.Fatal("Expected parsing to fail due to circular variable references")
	}

	expectedDiags := hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unresolvable variable",
			Detail:   "The variable 'hostname' could not be resolved due to missing or circular dependencies.",
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unresolvable variable",
			Detail:   "The variable 'var_c' could not be resolved due to missing or circular dependencies.",
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unresolvable variable",
			Detail:   "The variable 'var_a' could not be resolved due to missing or circular dependencies.",
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Unresolvable variable",
			Detail:   "The variable 'var_b' could not be resolved due to missing or circular dependencies.",
		},
	}

	verifyDiagnostics(t, expectedDiags, diags)

	if inventory != nil {
		t.Fatal("Inventory should be nil for invalid configuration")
	}
}

func TestNameConflictsParsing(t *testing.T) {
	path := filepath.Join("corpus", "error-cases", "name-conflicts.hcl")

	files, err := core.DiscoverInventoryFiles(path)
	if err != nil {
		t.Fatalf("Failed to discover inventory files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 inventory file, got %d", len(files))
	}

	inventory, diags := core.ParseInventoryFiles(files)
	if !diags.HasErrors() {
		t.Fatal("Expected parsing to fail due to name conflicts")
	}

	expectedDiags := hcl.Diagnostics{
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Name conflict",
			Detail:   `The group name "server1" conflicts with a host name.`,
		},
		&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Name conflict",
			Detail:   fmt.Sprintf(`The group name "server1" conflicts with a host name defined at "%s:21,1-15".`, path),
		},
	}

	verifyDiagnostics(t, expectedDiags, diags)

	if inventory != nil {
		t.Fatal("Inventory should be nil for invalid configuration")
	}
}
