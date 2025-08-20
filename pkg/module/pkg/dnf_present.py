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

def package_should_be_processed(base, package):
    subject = dnf.subject.Subject(package)
    solution = subject.get_best_solution(base.sack)

    if not solution:
        return True

    if solution["nevra"] is None:
        return True

    if solution["nevra"].has_just_name():
        return True

    nevra = solution["nevra"]
    kargs = {"name": nevra.name}

    if nevra.epoch is not None:
        kargs["epoch"] = nevra.epoch

    if nevra.version is not None:
        kargs["version"] = nevra.version

    if nevra.release is not None:
        kargs["release"] = nevra.release

    if nevra.arch is not None:
        kargs["architecture"] = nevra.arch

    installed = base.sack.query().installed().filter(**kargs)

    return len(installed) == 0

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

        packages_to_be_processed = []
        for name in package_names:
            name = name.strip()
            if package_should_be_processed(base, name):
                packages_to_be_processed.append(name)

        if len(packages_to_be_processed) == 0:
            succeed(base, [], [])

        for name in packages_to_be_processed:
            try:
                _ = base.install(name, strict=base.conf.strict)
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
            base.download_packages(base.transaction.install_set)
        except dnf.exceptions.DownloadError as e:
            fail(base, f"Download error: {str(e)}", traceback.format_exc())

        # TODO - Configure disable GPG check
        for pkg in base.transaction.install_set:
            gpg_response, gpg_error = base._sig_check_pkg(pkg)

            if gpg_response == 0:
                continue

            if gpg_response == 1:
                try:
                    base._get_key_for_package(pkg)
                    continue
                except dnf.exceptions.Error as e:
                    fail(base, f"GPG key error for package {pkg.name}: {str(e)}", traceback.format_exc())

            fail(base, f"GPG check failed for package {pkg.name}: {gpg_error}")

        try:
            tid = base.do_transaction()
        except dnf.exceptions.Error as e:
            fail(base, f"Transaction failed: {str(e)}", traceback.format_exc())

        if tid is None:
            succeed(base, [], [])

        if ARGS["autoremove"]:
            try:
                base.autoremove()
            except dnf.exceptions.Error as e:
                fail(base, f"Autoremove failed: {str(e)}", traceback.format_exc())

        succeed(base, installed, removed)

    except Exception as e:
        fail(base, f"Unknown error: {str(e)}", traceback.format_exc())

if not HAS_DNF:
    fail("dnf module is not available", DNF_TRACEBACK)

if not dnf.util.am_i_root():
    fail("This module must be run as root")

main()
