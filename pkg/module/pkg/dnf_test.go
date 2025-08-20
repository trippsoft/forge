// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package pkg

import (
	"context"
	"testing"

	"github.com/trippsoft/forge/pkg/info"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/trippsoft/forge/pkg/transport"
	"github.com/zclconf/go-cty/cty"
)

func TestDNF(t *testing.T) {
	tp := transport.NewMockTransport()

	config := &module.RunConfig{
		Transport: tp,
		HostInfo:  &info.HostInfo{},
		Input: map[string]cty.Value{
			"names": cty.ListVal([]cty.Value{cty.StringVal("test")}),
			"state": cty.StringVal("present"),
		},
	}

	m := &DNFModule{}

	m.Run(context.Background(), config)
}
