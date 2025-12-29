// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"encoding/base64"
	"errors"

	"golang.org/x/text/encoding/unicode"
)

func encodePowerShellAsUTF16LEBase64(powershell string) (string, error) {
	if powershell == "" {
		return "", errors.New("input PowerShell command cannot be empty")
	}

	utf16Encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	utf16Encoded, err := utf16Encoder.String(powershell)
	if err != nil {
		return "", err
	}

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(utf16Encoded))
	if base64Encoded == "" {
		return "", errors.New("failed to encode PowerShell command to base64")
	}

	return base64Encoded, nil
}
