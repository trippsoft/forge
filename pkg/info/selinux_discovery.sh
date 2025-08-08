#!/bin/sh
# This script is used to discover SELinux status on Linux systems.
# It returns whether SELinux is installed, its status, and type.

if [ ! -f /etc/selinux/config ]; then
	selinux_installed=0
else
	selinux_installed=1
	selinux_status=$(grep -E '^SELINUX=' /etc/selinux/config | cut -d '=' -f 2)
	selinux_type=$(grep -E '^SELINUXTYPE=' /etc/selinux/config | cut -d '=' -f 2)
fi

printf '{"selinux_installed": "%s", "selinux_status": "%s", "selinux_type": "%s"}\n' \
    "$selinux_installed" \
    "$selinux_status" \
    "$selinux_type"
