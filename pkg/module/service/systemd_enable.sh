#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_enable.sh is used to enable a systemd service.

if [ "$PREVIOUS_IS_ENABLED" != "enabled" ]; then
    systemctl enable "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf '{"error": "Failed to enable service %s"}\n' "$FORGE_NAME"
        exit 0
    fi
fi
