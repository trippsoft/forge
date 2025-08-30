// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/trippsoft/forge/pkg/hclspec"
	"github.com/trippsoft/forge/pkg/module"
	"github.com/zclconf/go-cty/cty"
)

//go:generate go run ../../../cmd/scriptimport/main.go service systemd_disable.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_enable.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_mask.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_print.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_restart.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_shared.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_start.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_stop.sh
//go:generate go run ../../../cmd/scriptimport/main.go service systemd_unmask.sh

var (
	_ module.Module = (*SystemdServiceModule)(nil)
)

// SystemdServiceModule is a module for managing systemd services.
type SystemdServiceModule struct{}

// InputSpec implements module.Module.
func (s *SystemdServiceModule) InputSpec() *hclspec.Spec {
	return inputSpec
}

// Validate implements module.Module.
func (s *SystemdServiceModule) Validate(config *module.RunConfig) error {
	return nil
}

// Run implements module.Module.
func (s *SystemdServiceModule) Run(ctx context.Context, config *module.RunConfig) *module.Result {
	if config == nil {
		return module.NewFailure(errors.New("config cannot be nil"), "")
	}

	if config.Input == nil {
		return module.NewFailure(errors.New("input cannot be nil"), "")
	}

	if config.Transport == nil {
		return module.NewFailure(errors.New("transport cannot be nil"), "")
	}

	name, _ := config.Input["name"]
	masked, _ := config.Input["masked"]
	enabled, _ := config.Input["enabled"]
	state, _ := config.Input["state"]

	if masked.IsNull() {
		if !enabled.IsNull() && enabled.True() {
			masked = cty.False
		}

		if !state.IsNull() && state.AsString() != "stopped" {
			masked = cty.False
		}
	}

	if enabled.IsNull() && !masked.IsNull() && masked.True() {
		enabled = cty.False
	}

	if state.IsNull() && !masked.IsNull() && masked.True() {
		state = cty.StringVal("stopped")
	}

	command := fmt.Sprintf("export FORGE_NAME=%q; %s", name.AsString(), SystemdSharedScript)

	if !config.WhatIf {
		if !masked.IsNull() && masked.True() {
			command = fmt.Sprintf("%s; %s", command, SystemdMaskScript)
		} else {
			if !masked.IsNull() && !masked.True() {
				command = fmt.Sprintf("%s; %s", command, SystemdUnmaskScript)
			}

			if !state.IsNull() {
				switch state.AsString() {
				case "started":
					command = fmt.Sprintf("%s; %s", command, SystemdStartScript)
				case "stopped":
					command = fmt.Sprintf("%s; %s", command, SystemdStopScript)
				case "restarted":
					command = fmt.Sprintf("%s; %s", command, SystemdRestartScript)
				}
			}

			if !enabled.IsNull() {
				if enabled.True() {
					command = fmt.Sprintf("%s; %s", command, SystemdEnableScript)
				} else {
					command = fmt.Sprintf("%s; %s", command, SystemdDisableScript)
				}
			}
		}
	}

	command = fmt.Sprintf("%s; %s", command, SystemdPrintScript)

	cmd, err := config.Transport.NewCommand(command, config.Escalation)
	if err != nil {
		return module.NewFailure(err, "")
	}
	stdout, stderr, err := cmd.OutputWithError(ctx)
	if err != nil {
		return module.NewFailure(err, stderr)
	}

	var result systemdServiceResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		return module.NewFailure(err, stderr)
	}

	if result.Error != "" {
		return module.NewFailure(errors.New(result.Error), stderr)
	}

	changed := false

	previousMasked := result.masked()
	previousEnabled := result.enabled()
	previousState := result.state()

	previous := map[string]cty.Value{
		"masked":  previousMasked,
		"enabled": previousEnabled,
		"state":   previousState,
	}

	if masked.IsNull() {
		masked = previousMasked
	}

	if enabled.IsNull() {
		enabled = previousEnabled
	}

	if state.IsNull() {
		state = previousState
	} else if state.AsString() == "restarted" {
		changed = true
		state = cty.StringVal("started")
	}

	if masked.Equals(previousMasked).False() {
		changed = true
	}

	if enabled.Equals(previousEnabled).False() {
		changed = true
	}

	if state.Equals(previousState).False() {
		changed = true
	}

	current := map[string]cty.Value{
		"masked":  masked,
		"enabled": enabled,
		"state":   state,
	}

	output := map[string]cty.Value{
		"previous": cty.ObjectVal(previous),
		"current":  cty.ObjectVal(current),
	}

	return module.NewSuccess(changed, output)
}

type systemdServiceResult struct {
	Error             string `json:"error,omitempty"`
	PreviousIsActive  string `json:"previous_is_active,omitempty"`
	PreviousIsEnabled string `json:"previous_is_enabled,omitempty"`
}

func (r systemdServiceResult) masked() cty.Value {
	return cty.BoolVal(r.PreviousIsEnabled == "masked")
}

func (r systemdServiceResult) enabled() cty.Value {
	return cty.BoolVal(r.PreviousIsEnabled == "enabled")
}

func (r systemdServiceResult) state() cty.Value {
	switch r.PreviousIsActive {
	case "active":
		return cty.StringVal("started")
	default:
		return cty.StringVal("stopped")
	}
}
