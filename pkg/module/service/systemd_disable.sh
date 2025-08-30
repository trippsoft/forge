#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_disable.sh is used to disable a systemd service.

if [ "$PREVIOUS_IS_ENABLED" != "disabled" ]; then
    systemctl disable "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf "{\"error\": \"Failed to disable service %s\"}\n" "$FORGE_NAME"
        exit 0
    fi
fi
