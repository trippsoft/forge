package info

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/trippsoft/forge/internal/transport"
	"github.com/zclconf/go-cty/cty"
)

type appArmorInfo struct {
	supported bool
	enabled   bool
}

func newAppArmorInfo() *appArmorInfo {
	return &appArmorInfo{
		supported: false,
		enabled:   false,
	}
}

func (a *appArmorInfo) Supported() bool {
	return a.supported
}

func (a *appArmorInfo) Enabled() bool {
	return a.enabled
}

func (a *appArmorInfo) populateAppArmorInfo(osInfo *osInfo, fileSystem transport.FileSystem) error {

	if !osInfo.families.Contains("linux") {
		a.supported = false
		a.enabled = false
		return nil
	}

	a.supported = true

	_, err := fileSystem.Stat("/sys/kernel/security/apparmor")
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
		a.enabled = false
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to stat AppArmor directory: %w", err)
	}

	a.enabled = true
	return nil
}

func (a *appArmorInfo) toMapOfCtyValues() map[string]cty.Value {

	if !a.supported {
		return map[string]cty.Value{
			"apparmor_enabled": cty.NullVal(cty.Bool),
		}
	}

	return map[string]cty.Value{
		"apparmor_enabled": cty.BoolVal(a.enabled),
	}
}

// String returns a string representation of the AppArmor information.
// This is useful for logging or debugging purposes.
func (a *appArmorInfo) String() string {

	if !a.supported {
		return "apparmor_enabled: not supported"
	}

	return fmt.Sprintf("apparmor_enabled: %t", a.enabled)
}
