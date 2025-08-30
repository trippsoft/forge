#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# systemd_shared.sh provides utility functions shared between systemd service management scripts.

get_is_active() {
    local service_name="$1"
    local is_active

    is_active=$(systemctl is-active "$service_name")

    if [ "$?" -eq 4 ]; then
        is_active="unknown"
    fi

    echo "$is_active"
}

get_is_enabled() {
    local service_name="$1"
    local is_enabled

    is_enabled=$(systemctl is-enabled "$service_name")

    if [ "$?" -eq 4 ]; then
        is_enabled="unknown"
    fi

    echo "$is_enabled"
}

print_json_result() {
    local previous_is_active="$1"
    local previous_is_enabled="$2"

    printf "{\"previous_is_active\": \"%s\", \"previous_is_enabled\": \"%s\"}\n" \
        "$previous_is_active" \
        "$previous_is_enabled"
}

CURRENT_USER=$(whoami)

if [ "$CURRENT_USER" != "root" ]; then
    printf "{\"error\": \"Service %s can only be managed by root\"}\n" "$FORGE_NAME"
    exit 0
fi

PREVIOUS_IS_ACTIVE=$(get_is_active "$FORGE_NAME")

if [ "$PREVIOUS_IS_ACTIVE" = "unknown" ]; then
    printf "{\"error\": \"Service %s does not appear to exist\"}\n" "$FORGE_NAME"
    exit 0
fi

PREVIOUS_IS_ENABLED=$(get_is_enabled "$FORGE_NAME")

if [ "$PREVIOUS_IS_ENABLED" = "unknown" ]; then
    printf '{"error": "Service %s does not appear to exist"}\n' "$FORGE_NAME"
    exit 0
fi
