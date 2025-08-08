#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used to discover AppArmor status on Linux systems.
# It returns 1 if AppArmor is enabled, 0 otherwise.
if [ -d /sys/kernel/security/apparmor ]; then
    apparmor_enabled=1
else
    apparmor_enabled=0
fi
echo "$apparmor_enabled"
