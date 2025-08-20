#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.
# ARGS will be replaced by the module arguments.

import json
import traceback

try:
    import dnf
    import dnf.base
    import dnf.exceptions
    import dnf.subject
    import dnf.util
    HAS_DNF = True
    DNF_TRACEBACK = None
except ImportError:
    HAS_DNF = False
    DNF_TRACEBACK = traceback.format_exc()

ARGS = {}

def fail(base, error, detail = ""):
    if base is not None:
        try:
            base.close()
        except Exception:
            pass

    output = {
        "error": error,
        "error_detail": detail,
        "changed": False,
        "installed_packages": [],
        "removed_packages": [],
    }
    print(json.dumps(output), flush=True)
    exit(0)

def succeed(base, installed, removed):
    if base is not None:
        try:
            base.close()
        except Exception:
            pass

    output = {
        "error": None,
        "error_detail": None,
        "changed": len(installed) > 0 or len(removed) > 0,
        "installed_packages": installed,
        "removed_packages": removed,
    }
    print(json.dumps(output), flush=True)
    exit(0)

def setup_base():
    # TODO - Configure DNF (GPG checking, disable/enable repo, SSL verification)
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
        if ARGS["update_cache"]:
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

        package_names = ARGS["names"]

        to_be_removed = packages_to_be_processed(base, package_names)

        if len(to_be_removed) == 0:
            succeed(base, [], [])

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
                succeed(base, [], [])
        except dnf.exceptions.DepsolveError as e:
            fail(base, f"Dependency resolution failed: {str(e)}", traceback.format_exc())
        except dnf.exceptions.Error as e:
            fail(base, f"Error creating transaction: {str(e)}", traceback.format_exc())

        transaction = base.transaction
        if transaction is None:
            succeed(base, [], [])

        installed, removed = get_output_from_transaction(transaction)
        if ARGS["what_if"] or (len(installed) == 0 and len(removed) == 0):
            succeed(base, installed, removed)

        try:
            tid = base.do_transaction()
        except dnf.exceptions.Error as e:
            fail(base, f"Transaction failed: {str(e)}", traceback.format_exc())

        if tid is None:
            succeed(base, [], [])

        succeed(base, installed, removed)

    except Exception as e:
        fail(base, f"Unknown error: {str(e)}", traceback.format_exc())

if not HAS_DNF:
    fail("dnf module is not available", DNF_TRACEBACK)

if not dnf.util.am_i_root():
    fail("This module must be run as root")

main()
