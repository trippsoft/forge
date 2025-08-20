#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.

import json
import traceback

try:
    import dnf
    import dnf.base
    import dnf.exceptions
    HAS_DNF = True
    DNF_TRACEBACK = None
except ImportError:
    HAS_DNF = False
    DNF_TRACEBACK = traceback.format_exc()

def fail(base, error, detail = ""):
    if base is not None:
        try:
            base.close()
        except Exception:
            pass

    output = {
        "error": error,
        "error_detail": detail,
        "packages": {},
    }
    print(json.dumps(output))
    exit(0)

def setup_base():
    base = dnf.base.Base()

    try:
        base.setup_loggers()
    except AttributeError:
        pass

    try:
        base.init_plugins()
        base.pre_configure_plugins()
    except AttributeError:
        pass

    try:
        _ = base.fill_sack()
    except dnf.exceptions.RepoError as e:
        fail(base, "Repository error: " + str(e), traceback.format_exc())

    return base

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
            "error": None,
            "error_detail": None,
            "packages": packages,
        }

        try:
            base.close()
        except Exception:
            pass

        print(json.dumps(output))

    except Exception as e:
        fail(base, "Unknown error: " + str(e), traceback.format_exc())

if not HAS_DNF:
    fail(None, "dnf module is not available", DNF_TRACEBACK)

main()
