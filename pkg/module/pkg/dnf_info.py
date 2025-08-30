#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.

def main():
    try:
        base = setup_base()
        installed = base.sack.query().installed()
        packages = {}

        for pkg in installed:
            packages[pkg.name] = {
                "name": pkg.name,
                "epoch": str(pkg.epoch),
                "version": pkg.version,
                "release": pkg.release,
                "architecture": pkg.arch,
                "repo": pkg.from_repo,
            }

        output = {
            "packages": packages,
        }

        dnf_success(base, output)        

    except Exception as e:
        import traceback
        dnf_fail(base, "Unknown error: " + str(e), traceback.format_exc())

main()
