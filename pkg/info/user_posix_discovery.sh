#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# This script is used to discover the POSIX user information on POSIX systems.
# It returns the username, user ID, user GID, home directory, shell, and
# GECOS info of the current user.

user_name=$(id -nu)
user_id=$(id -u)
user_gid=$(id -g)
user_home_dir="$HOME"
user_shell="$SHELL"
user_gecos=$(getent passwd $user_name | cut -d ':' -f 5)

printf '{"user_name": "%s", "user_id": "%s", "user_gid": "%s", "user_home_dir": "%s", "user_shell": "%s", "user_gecos": "%s"}\n' \
    "$(escape_json "$user_name")" \
    "$user_id" \
    "$user_gid" \
    "$(escape_json "$user_home_dir")" \
    "$user_shell" \
    "$(escape_json "$user_gecos")"
