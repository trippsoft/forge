#!/bin/sh
# This script is used to discover the OS and architecture on Darwin systems.
# It sets the variables 'os_arch' and 'os_version' which can be used in
# other scripts or commands to determine the system's architecture and version.

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
os_version="$(sw_vers -productVersion || echo \"\")"

printf '{"os_arch": "%s", "os_version": "%s"}\n' \
    "$(escape_json "$os_arch")" \
    "$(escape_json "$os_version")"
