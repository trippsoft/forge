// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package info

import (
	"context"
	"slices"
	"strings"

	"github.com/trippsoft/forge/pkg/transport"
	"github.com/trippsoft/forge/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

// LocaleInfo holds information about a POSIX system's locale settings.
type LocaleInfo struct {
	locales []string
}

func newLocaleInfo() *LocaleInfo {
	return &LocaleInfo{
		locales: []string{},
	}
}

// Locales returns a copy of the locales slice.
func (l *LocaleInfo) Locales() []string {
	locales := slices.Clone(l.locales)
	return locales
}

func (l *LocaleInfo) populateLocaleInfo(osInfo *OSInfo, transport transport.Transport) util.Diags {
	if osInfo == nil || osInfo.id == "" {
		return util.Diags{&util.Diag{
			Severity: util.DiagWarning,
			Summary:  "Missing OS information",
			Detail:   "Skipping locale information collection due to missing or invalid OS info",
		}}
	}

	if !osInfo.families.Contains("posix") {
		return nil
	}

	cmd, err := transport.NewCommand("locale -a", nil)
	if err != nil {
		return util.Diags{&util.Diag{
			Severity: util.DiagError,
			Summary:  "Failed to create command",
			Detail:   err.Error(),
		}}
	}

	stdout, stderr, err := cmd.OutputWithError(context.Background())
	if err != nil {
		return util.Diags{
			&util.Diag{
				Severity: util.DiagError,
				Summary:  "Failed to execute command: %v",
				Detail:   err.Error(),
			},
			&util.Diag{
				Severity: util.DiagDebug,
				Summary:  "Command stderr",
				Detail:   stderr,
			},
		}
	}

	locales := strings.Split(strings.TrimSpace(stdout), "\n")
	l.locales = locales

	return nil
}

func (l *LocaleInfo) toMapOfCtyValues() map[string]cty.Value {
	values := make(map[string]cty.Value)

	if len(l.locales) == 0 {
		values["locales"] = cty.NullVal(cty.List(cty.String))
		return values
	}

	localeValues := make([]cty.Value, 0, len(l.locales))

	for _, locale := range l.locales {
		localeValues = append(localeValues, cty.StringVal(locale))
	}

	values["locales"] = cty.ListVal(localeValues)

	return values
}

func (l *LocaleInfo) String() string {
	stringBuilder := &strings.Builder{}
	stringBuilder.WriteString("locales: ")

	if len(l.locales) == 0 {
		stringBuilder.WriteString("none\n")
	} else {
		for i, locale := range l.locales {
			if i > 0 {
				stringBuilder.WriteString(", ")
			}

			stringBuilder.WriteString(locale)
		}

		stringBuilder.WriteString("\n")
	}

	return stringBuilder.String()
}
