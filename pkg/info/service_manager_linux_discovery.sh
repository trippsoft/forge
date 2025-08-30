#!/bin/sh
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# This script is used to discover init system information on Linux systems.
# It returns whether systemd, initctl, OpenRC, and other init systems exist,
# along with the target of the /sbin/init symlink and the command running as PID 1.

systemctl_exists=0
run_systemd_system_exists=0
dev_run_systemd_exists=0
dev_systemd_exists=0
initctl_exists=0
etc_init_exists=0
openrc_exists=0
init_link_target=""
etc_init_d_exists=0
proc1_comm=""

systemctl_path=$(which systemctl 2>/dev/null || echo "")

if [ -z "$systemctl_path" ]; then
	systemctl_exists=0
elif [ -x "$systemctl_path" ]; then
	systemctl_exists=1
fi
if [ -d /run/systemd/system ]; then
	run_systemd_system_exists=1
fi
if [ -d /dev/.run/systemd ]; then
	dev_run_systemd_exists=1
fi
if [ -d /dev/.systemd ]; then
	dev_systemd_exists=1
fi

initctl_path=$(which initctl 2>/dev/null || echo "")

if [ -z "$initctl_path" ]; then
	initctl_exists=0
elif [ -f "$initctl_path" ]; then
	initctl_exists=1
fi
if [ -d /etc/init ]; then
	etc_init_exists=1
fi
if [ -f /sbin/openrc ]; then
	openrc_exists=1
fi
if [ -L /sbin/init ]; then
	init_link_target=$(readlink /sbin/init)
fi
if [ -d /etc/init.d ]; then
	etc_init_d_exists=1
fi
if [ -f /proc/1/comm ]; then
	proc1_comm=$(cat /proc/1/comm)
fi

printf '{"systemctl_exists": "%s", "run_systemd_system_exists": "%s", "dev_run_systemd_exists": "%s", ' \
    "$systemctl_exists" \
    "$run_systemd_system_exists" \
    "$dev_run_systemd_exists"
printf '"dev_systemd_exists": "%s", "initctl_exists": "%s", "etc_init_exists": "%s", ' \
    "$dev_systemd_exists" \
    "$initctl_exists" \
    "$etc_init_exists"
printf '"openrc_exists": "%s", "init_link_target": "%s", ' \
    "$openrc_exists" \
    "$init_link_target"
printf '"etc_init_d_exists": "%s", "proc1_comm": "%s"}\n' \
    "$etc_init_d_exists" \
    "$proc1_comm"
