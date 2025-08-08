#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used to discover FIPS mode status on Linux systems.
# It returns 1 if FIPS mode is enabled, 0 otherwise.

if [ -f /proc/sys/crypto/fips_enabled ]; then
    fips_enabled=$(cat /proc/sys/crypto/fips_enabled)
else
    fips_enabled=0
fi

echo "$fips_enabled"
