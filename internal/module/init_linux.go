// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package module

import (
	"github.com/trippsoft/forge/pkg/python"
)

func init() {
	prunedDnfInfoScript = python.RemoveEmptyLinesAndComments(dnfInfoScript)
	prunedDnfAbsentScript = python.RemoveEmptyLinesAndComments(dnfAbsentScript)
	prunedDnfLatestScript = python.RemoveEmptyLinesAndComments(dnfLatestScript)
	prunedDnfPresentScript = python.RemoveEmptyLinesAndComments(dnfPresentScript)
}
