#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
# 
# This script is used by the dnf module to apply the package changes.
# ARGS will be replaced by the module arguments.

def packages_to_be_processed(base, packages):
    to_be_installed = []
    to_be_updated = []
    update_all = False

    for package in packages:
        if package == "*":
            update_all = True
            continue

        subject = dnf.subject.Subject(package)
        solution = subject.get_best_solution(base.sack)

        if not solution:
            to_be_installed.append(package)
            continue

        if solution["nevra"] is None:
            to_be_installed.append(package)
            continue

        if solution["nevra"].has_just_name():
            if not update_all:
                to_be_updated.append(package)

            continue

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

        if len(installed) == 0:
            to_be_installed.append(package)

    if update_all:
        return to_be_installed, [], update_all

    return to_be_installed, to_be_updated, update_all

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

        to_be_installed, to_be_updated, update_all = packages_to_be_processed(base, package_names)
        if len(to_be_installed) == 0 and len(to_be_updated) == 0 and not update_all:
            dnf_success(base, {"installed": [], "removed": []})

        if update_all:
            try:
                _ = base.upgrade_all()
            except dnf.exceptions.DepsolveError as e:
                import traceback
                dnf_fail(base, f"Dependency error for package upgrades: {str(e)}", traceback.format_exc())
            except dnf.exceptions.Error as e:
                import traceback
                dnf_fail(base, f"Error processing package upgrades: {str(e)}", traceback.format_exc())

        for name in to_be_installed:
            try:
                _ = base.install(name, strict=base.conf.strict)
            except dnf.exceptions.MarkingError as e:
                import traceback
                dnf_fail(base, f"Package {name} not found: {str(e)}", traceback.format_exc())
            except dnf.exceptions.DepsolveError as e:
                import traceback
                dnf_fail(base, f"Dependency error for package {name}: {str(e)}", traceback.format_exc())
            except dnf.exceptions.Error as e:
                import traceback
                dnf_fail(base, f"Error processing package {name}: {str(e)}", traceback.format_exc())

        for name in to_be_updated:
            try:
                _ = base.upgrade(name)
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
            base.download_packages(base.transaction.install_set)
        except dnf.exceptions.DownloadError as e:
            import traceback
            dnf_fail(base, f"Download error: {str(e)}", traceback.format_exc())

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
                    import traceback
                    dnf_fail(base, f"GPG key error for package {pkg.name}: {str(e)}", traceback.format_exc())

            dnf_fail(base, f"GPG check failed for package {pkg.name}: {gpg_error}")

        try:
            tid = base.do_transaction()
        except dnf.exceptions.Error as e:
            import traceback
            dnf_fail(base, f"Transaction failed: {str(e)}", traceback.format_exc())

        if tid is None:
            dnf_success(base, {"installed": [], "removed": []})

        if ARGS["autoremove"]:
            try:
                base.autoremove()
            except dnf.exceptions.Error as e:
                import traceback
                dnf_fail(base, f"Autoremove failed: {str(e)}", traceback.format_exc())

        dnf_success(base, {"installed": installed, "removed": removed})

    except Exception as e:
        import traceback
        dnf_fail(base, f"Unknown error: {str(e)}", traceback.format_exc())

if not dnf.util.am_i_root():
    fail("This module must be run as root")

main()
