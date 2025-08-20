// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"strings"
	"unicode"
)

// RemoveEmptyLinesAndComments removes empty lines and comments from a script.
func RemoveEmptyLinesAndComments(script string) string {
	lines := strings.Split(script, "\n")
	filteredLines := []string{}
	for _, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		filteredLines = append(filteredLines, line)
	}

	return strings.Join(filteredLines, "\n")
}
