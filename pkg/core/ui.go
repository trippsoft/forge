// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package core

// UI represents a user interface for text output.
//
// Each implementation will be specific to the type of UI (e.g. CLI, Packer plugin).
// The implementation should handle secret filtering and text formatting.
type UI interface {
	AddSecret(secret string) // AddSecret adds a secret to be filtered from output.
	ClearSecrets()           // ClearSecrets clears all secrets from the secret filter.
}
