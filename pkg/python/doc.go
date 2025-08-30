// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

// Package python provides Python utility functions.
package python

import (
	_ "embed"
)

//go:embed fail.py
var FailFunction string

//go:embed success.py
var SuccessFunction string
