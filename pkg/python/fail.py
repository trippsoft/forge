#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# fail.py provides a utility function to fail gracefully and
# return a meaningful error message.

def fail(error, detail = ""):
    output = {
        "error": error,
        "error_detail": detail
    }

    import json
    print(json.dumps(output), flush=True)
    exit(0)
