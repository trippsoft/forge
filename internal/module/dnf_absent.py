#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.
# The module arguments will be available in the ARGS variable.

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
        if INPUT.get("update_cache", False):
            _ = base.update_cache()

        _ = base.fill_sack()
    except dnf.exceptions.RepoError as e:
        fail(base, f"Repository error: {str(e)}", traceback.format_exc())

    return base

def packages_to_be_processed(base, packages):
    to_remove_packages = []
    for package in packages:
        package = package.strip()

        subject = dnf.subject.Subject(package)
        query = subject.get_best_query(base.sack)

        if not query:
            to_remove_packages.append(package)
            continue

        installed = query.installed()
        if len(installed) > 0:
            for pkg in installed:
                to_remove_packages.append(pkg.name)

    return to_remove_packages

def get_output_from_transaction(transaction):
    installed = []
    removed = []

    for pkg in transaction.install_set:
        installed.append({
            "name": pkg.name,
            "epoch": str(pkg.epoch),
            "version": pkg.version,
            "release": pkg.release,
            "architecture": pkg.arch,
            "repo": pkg.repoid,
        })

    for pkg in transaction.remove_set:
        removed.append({
            "name": pkg.name,
            "epoch": str(pkg.epoch),
            "version": pkg.version,
            "release": pkg.release,
            "architecture": pkg.arch,
            "repo": pkg.from_repo,
        })

    return installed, removed

def main():
    try:
        base = setup_base()

        package_names = INPUT["names"]

        to_be_removed = packages_to_be_processed(base, package_names)

        if len(to_be_removed) == 0:
            success(base, {"installed": [], "removed": []})

        for name in to_be_removed:
            try:
                _ = base.remove(name)
            except dnf.exceptions.MarkingError as e:
                fail(base, f"Package {name} not found: {str(e)}", traceback.format_exc())
            except dnf.exceptions.DepsolveError as e:
                fail(base, f"Dependency error for package {name}: {str(e)}", traceback.format_exc())
            except dnf.exceptions.Error as e:
                fail(base, f"Error processing package {name}: {str(e)}", traceback.format_exc())

        # TODO - Add allow erasing as input
        try:
            if not base.resolve(allow_erasing=False):
                success(base, {"installed": [], "removed": []})
        except dnf.exceptions.DepsolveError as e:
            fail(base, f"Dependency resolution failed: {str(e)}", traceback.format_exc())
        except dnf.exceptions.Error as e:
            fail(base, f"Error creating transaction: {str(e)}", traceback.format_exc())

        transaction = base.transaction
        if transaction is None:
            success(base, {"installed": [], "removed": []})

        installed, removed = get_output_from_transaction(transaction)
        if WHAT_IF or (len(installed) == 0 and len(removed) == 0):
            success(base, {"installed": installed, "removed": removed})

        try:
            tid = base.do_transaction()
        except dnf.exceptions.Error as e:
            fail(base, f"Transaction failed: {str(e)}", traceback.format_exc())

        if tid is None:
            success(base, {"installed": [], "removed": []})

        success(base, {"installed": installed, "removed": removed})

    except Exception as e:
        fail(base, f"Unknown error: {str(e)}", traceback.format_exc())

if not dnf.util.am_i_root():
    fail(None, "This module must be run as root")

main()
