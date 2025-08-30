#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_start.sh is used to start a systemd service.

if [ $PREVIOUS_IS_ACTIVE != "active" ]; then
    systemctl start "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf '{"error": "Failed to start service %s"}\n' "$FORGE_NAME"
        exit 0
    fi
fi
