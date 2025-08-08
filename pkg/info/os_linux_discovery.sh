#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# This script is used to discover the OS and architecture on Linux systems.
# It sets the variables 'os_arch', 'os_id', 'os_friendly_name',
# 'os_release', 'os_version', 'os_edition', and 'os_edition_id'
# which can be used in other scripts or commands to determine the system's
# architecture and version.

escape_json() {
    string="$1"

    string=$(printf '%s' "$string" | sed 's/\\/\\\\/g')
    string=$(printf '%s' "$string" | sed 's/"/\\"/g') 
    string=$(printf '%s' "$string" | sed 's/\n/\\n/g')
    string=$(printf '%s' "$string" | sed 's/\r/\\r/g')
    string=$(printf '%s' "$string" | sed 's/\t/\\t/g')
    printf '%s' "$string"
}

os_arch="$(uname -m || echo \"\")"

if [ -e /etc/os-release ]; then
    source /etc/os-release
elif [ -L /etc/os-release ]; then
    source "$(readlink -f /etc/os-release || echo \"\")"
elif [ -e /usr/lib/os-release ]; then
    source /usr/lib/os-release
elif [ -L /usr/lib/os-release ]; then
    source "$(readlink -f /usr/lib/os-release || echo \"\")"
fi

if [ -n "$ID" ]; then
    os_id="$ID"
else
    os_id="$(lsb_release -si || echo \"\")"
fi

if [ -n "$PRETTY_NAME" ]; then
    os_friendly_name="$PRETTY_NAME"
else
    os_friendly_name="$(lsb_release -sd || echo \"\")"
fi

if [ -n "$VERSION_ID" ]; then
    os_version="$VERSION_ID"
else
    os_version="$(lsb_release -sr || echo \"\")"
fi

if [ -n "$VERSION_CODENAME" ]; then
    os_release="$VERSION_CODENAME"
else
    os_release="$(lsb_release -sc || echo \"\")"
fi

if [ -n "$VARIANT" ]; then
    os_edition="$VARIANT"
fi

if [ -n "$VARIANT_ID" ]; then
    os_edition_id="$VARIANT_ID"
fi

printf '{"os_arch": "%s", "os_id": "%s", "os_friendly_name": "%s", "os_release": "%s", "os_version": "%s", "os_edition": "%s", "os_edition_id": "%s"}\n' \
    "$(escape_json "$os_arch")" \
    "$(escape_json "$os_id")" \
    "$(escape_json "$os_friendly_name")" \
    "$(escape_json "$os_release")" \
    "$(escape_json "$os_version")" \
    "$(escape_json "$os_edition")" \
    "$(escape_json "$os_edition_id")"
