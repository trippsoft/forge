// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

//go:build linux

package module

import "github.com/trippsoft/forge/pkg/util"

func init() {
	prunedDnfInfoScript = util.RemoveEmptyLinesAndComments(dnfInfoScript)
	prunedDnfAbsentScript = util.RemoveEmptyLinesAndComments(dnfAbsentScript)
	prunedDnfLatestScript = util.RemoveEmptyLinesAndComments(dnfLatestScript)
	prunedDnfPresentScript = util.RemoveEmptyLinesAndComments(dnfPresentScript)
}
