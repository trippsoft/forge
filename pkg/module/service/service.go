// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"errors"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

var (
	inputSpec = hclspec.NewSpec(
		hclspec.Object(
			hclspec.RequiredField("name", hclspec.String),
			hclspec.OptionalField("state", hclspec.String).
				WithDefaultValue(cty.StringVal("started")).
				WithConstraints(hclspec.AllowedValues(
					cty.StringVal("started"),
					cty.StringVal("stopped"),
					cty.StringVal("restarted"),
				)),
			hclspec.OptionalField("masked", hclspec.Bool).
				WithDefaultValue(cty.BoolVal(false)),
			hclspec.OptionalField("enabled", hclspec.Bool),
		).WithConstraints(
			hclspec.RequiredOneOf("state", "masked", "enabled"),
			hclspec.ConditionalConstraint(
				hclspec.FieldEquals("masked", cty.True),
				hclspec.AllowedFieldValues("state", cty.StringVal("stopped")),
			),
			hclspec.ConditionalConstraint(
				hclspec.FieldEquals("masked", cty.True),
				hclspec.AllowedFieldValues("enabled", cty.False),
			),
			hclspec.ConditionalConstraint(
				hclspec.FieldEquals("enabled", cty.True),
				hclspec.AllowedFieldValues("masked", cty.False),
			),
			hclspec.ConditionalConstraint(
				hclspec.FieldEquals("state", cty.StringVal("started")),
				hclspec.AllowedFieldValues("masked", cty.False),
			),
			hclspec.ConditionalConstraint(
				hclspec.FieldEquals("state", cty.StringVal("restarted")),
				hclspec.AllowedFieldValues("masked", cty.False),
			),
		),
	)

	serviceManagerModules = map[string]module.Module{
		"systemd": &SystemdServiceModule{},
	}

	_ module.Module = (*ServiceModule)(nil)
)

// ServiceModule is a module for managing services, selecting the
// implementation based on discovered service manager.
type ServiceModule struct{}

// InputSpec implements module.Module.
func (s *ServiceModule) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (s *ServiceModule) Validate(config *module.RunConfig) error {
	serviceManagerName := config.HostInfo.ServiceManagerInfo().Name()
	mod, exists := serviceManagerModules[serviceManagerName]
	if !exists {
		return errors.New("unsupported service manager")
	}

	return mod.Validate(config)
}

// Run implements module.Module.
func (s *ServiceModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	serviceManagerName := config.HostInfo.ServiceManagerInfo().Name()
	mod, exists := serviceManagerModules[serviceManagerName]
	if !exists {
		return module.NewFailure(nil, "unsupported service manager")
	}

	return mod.Run(ctx, config)
}
