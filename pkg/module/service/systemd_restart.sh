#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_restart.sh is used to restart a systemd service.

systemctl restart "$FORGE_NAME" > /dev/null

if [ "$?" -ne 0 ]; then
    printf "{\"error\": \"Failed to restart service %s\"}\n" "$FORGE_NAME"
    exit 0
fi
