#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_mask.sh is used to mask a systemd service.

if [ "$PREVIOUS_IS_ENABLED" != "masked" ]; then
    systemctl stop "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf "{\"error\": \"Failed to stop service %s\"}\n" "$FORGE_NAME"
        exit 0
    fi

    systemctl mask "$FORGE_NAME" > /dev/null

    if [ "$?" -ne 0 ]; then
        printf "{\"error\": \"Failed to mask service %s\"}\n" "$FORGE_NAME"
        exit 0
    fi
fi
