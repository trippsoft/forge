#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_stop.sh is used to stop a systemd service.

if [ $PREVIOUS_IS_ACTIVE != "inactive" ]; then
    systemctl stop "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf "{\"error\": \"Failed to stop service %s\"}\n" "$FORGE_NAME"
        exit 0
    fi
fi
