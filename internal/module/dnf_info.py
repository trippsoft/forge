#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.

import json
import sys
import traceback

try:
    import dnf
    import dnf.base
    import dnf.exceptions
    DNF_AVAILABLE = True
except ImportError:
    DNF_AVAILABLE = False

def success(base, output):
    if base is not None:
        try:
            base.close()
        except:
            pass

    print(json.dumps(output), flush=True)
    exit(0)

def fail(base, message, details=""):
    if base is not None:
        try:
            base.close()
        except:
            pass

    output = {
        "error": message,
        "error_details": details,
    }
    print(json.dumps(output), flush=True, file=sys.stderr)
    exit(1)

def setup_base():
    base = dnf.base.Base()
    base.conf.read()
    base.conf.debuglevel = 0
    base.conf.assumeyes = True
    base.conf.substitutions.update_from_etc(base.conf.installroot)

    try:
        base.setup_loggers()
    except AttributeError:
        pass

    try:
        # TODO - Add/remove plugins
        base.init_plugins()
        base.pre_configure_plugins()
    except AttributeError:
        pass

    base.read_all_repos()

    try:
        base.configure_plugins()
    except AttributeError:
        pass

    try:
        _ = base.fill_sack()
    except dnf.exceptions.RepoError as e:
        fail(base, f"Repository error: {str(e)}", traceback.format_exc())

    return base

def main():
    if not DNF_AVAILABLE:
        fail(None, "DNF Python bindings are not available.")

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

        success(base, output)        

    except Exception as e:
        fail(base, "Unknown error: " + str(e), traceback.format_exc())

main()
