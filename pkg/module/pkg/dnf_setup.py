#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# dnf_setup.py provides utility functions for setting up DNF operations.
# It expects the fail.py and success.py utility functions to be available.

try:
    import dnf
    import dnf.base
    import dnf.exceptions
    import dnf.subject
    import dnf.util
except ImportError:
    import traceback
    fail("dnf module is not available", traceback.format_exc())

def dnf_fail(base, error, detail=""):
    if base is not None:
        try:
            base.close()
        except Exception:
            pass

    fail(error, detail)

def dnf_success(base, output):
    if base is not None:
        try:
            base.close()
        except Exception:
            pass

    success(output)

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
        if ARGS.get("update_cache", False):
            _ = base.update_cache()

        _ = base.fill_sack()
    except dnf.exceptions.RepoError as e:
        dnf_fail(base, f"Repository error: {str(e)}", traceback.format_exc())

    return base
