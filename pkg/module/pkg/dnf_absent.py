#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.
# The module arguments will be available in the ARGS variable.

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
            dnf_success(base, {"installed": [], "removed": []})

        for name in to_be_removed:
            try:
                _ = base.remove(name)
            except dnf.exceptions.MarkingError as e:
                import traceback
                dnf_fail(base, f"Package {name} not found: {str(e)}", traceback.format_exc())
            except dnf.exceptions.DepsolveError as e:
                import traceback
                dnf_fail(base, f"Dependency error for package {name}: {str(e)}", traceback.format_exc())
            except dnf.exceptions.Error as e:
                import traceback
                dnf_fail(base, f"Error processing package {name}: {str(e)}", traceback.format_exc())

        # TODO - Add allow erasing as input
        try:
            if not base.resolve(allow_erasing=False):
                dnf_success(base, {"installed": [], "removed": []})
        except dnf.exceptions.DepsolveError as e:
            import traceback
            dnf_fail(base, f"Dependency resolution failed: {str(e)}", traceback.format_exc())
        except dnf.exceptions.Error as e:
            import traceback
            dnf_fail(base, f"Error creating transaction: {str(e)}", traceback.format_exc())

        transaction = base.transaction
        if transaction is None:
            dnf_success(base, {"installed": [], "removed": []})

        installed, removed = get_output_from_transaction(transaction)
        if ARGS["what_if"] or (len(installed) == 0 and len(removed) == 0):
            dnf_success(base, {"installed": installed, "removed": removed})

        try:
            tid = base.do_transaction()
        except dnf.exceptions.Error as e:
            import traceback
            dnf_fail(base, f"Transaction failed: {str(e)}", traceback.format_exc())

        if tid is None:
            dnf_success(base, {"installed": [], "removed": []})

        dnf_success(base, {"installed": installed, "removed": removed})

    except Exception as e:
        import traceback
        dnf_fail(base, f"Unknown error: {str(e)}", traceback.format_exc())

if not dnf.util.am_i_root():
    fail("This module must be run as root")

main()
