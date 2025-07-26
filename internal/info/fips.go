package info

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

const (
	fipsPowerShellCommand = `$value = Get-ItemPropertyValue -LiteralPath 'HKLM:\SYSTEM\CurrentControlSet\Control\Lsa\FipsAlgorithm' -Name 'Enabled' -ErrorAction SilentlyContinue; Write-Host $value`
)

type fipsInfo struct {
	known   bool
	enabled bool
}

func newFipsInfo() *fipsInfo {
	return &fipsInfo{
		known:   false,
		enabled: false,
	}
}

func (f *fipsInfo) Known() bool {
	return f.known
}

func (f *fipsInfo) Enabled() bool {
	return f.enabled
}

func (f *fipsInfo) populateFipsInfo(osInfo *osInfo, transport transport.Transport) error {

	if osInfo.families.Contains("linux") {
		f.known = true
		return f.populateLinuxFipsInfo(transport)
	}

	if osInfo.families.Contains("windows") {
		f.known = true
		return f.populateWindowsFipsInfo(transport)
	}

	f.known = false
	f.enabled = false
	return nil
}

func (f *fipsInfo) populateLinuxFipsInfo(transport transport.Transport) error {

	stdout, _, err := transport.ExecuteCommand(context.Background(), "cat /proc/sys/crypto/fips_enabled || echo 0")
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	value, err := strconv.Atoi(strings.TrimSpace(stdout))
	if err != nil {
		return fmt.Errorf("failed to convert FIPS status: %w", err)
	}

	f.enabled = value == 1
	return nil
}

func (f *fipsInfo) populateWindowsFipsInfo(transport transport.Transport) error {

	stdout, _, err := transport.ExecutePowerShell(context.Background(), fipsPowerShellCommand)
	if err != nil {
		return fmt.Errorf("failed to execute PowerShell command: %w", err)
	}

	stdout = strings.TrimSpace(stdout)

	f.known = true
	f.enabled = stdout == "1"
	return nil
}

func (f *fipsInfo) toMapOfCtyValues() map[string]cty.Value {

	if !f.known {
		return map[string]cty.Value{
			"fips_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"fips_enabled": cty.BoolVal(f.enabled),
	}
}

// String returns a string representation of the FIPS information.
// This is useful for logging or debugging purposes.
func (f *fipsInfo) String() string {

	if !f.known {
		return "fips_enabled: unknown on this OS"
	}

	return fmt.Sprintf("fips_enabled: %t", f.enabled)
}
