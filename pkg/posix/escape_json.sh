#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# Contains POSIX shell utility function to escape JSON strings

escape_json() {
    local string="$1"

    string=$(printf '%s' "$string" | sed 's/\\/\\\\/g')
    string=$(printf '%s' "$string" | sed 's/"/\\"/g') 
    string=$(printf '%s' "$string" | sed 's/\n/\\n/g')
    string=$(printf '%s' "$string" | sed 's/\r/\\r/g')
    string=$(printf '%s' "$string" | sed 's/\t/\\t/g')
    printf '%s' "$string"
}
