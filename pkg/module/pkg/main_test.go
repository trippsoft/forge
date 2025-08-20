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

func TestScratch(t *testing.T) {
	builder, err := transport.NewSSHBuilder()
	if err != nil {
		t.Fatalf("failed to create SSH builder: %v", err)
	}

	tp, err := builder.Host("192.168.121.198").
		Port(22).
		User("vagrant").
		PasswordAuth("vagrant").
		Build()

	if err != nil {
		t.Fatalf("failed to build transport: %v", err)
	}

	config := &module.RunConfig{
		Transport: tp,
		HostInfo:  &info.HostInfo{},
		Input:     map[string]cty.Value{},
	}

	m := &DNFInfoModule{}

	m.Run(context.Background(), config)
}
