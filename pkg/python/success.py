#!/usr/bin/env python3
# Copyright (c) Forge
# SPDX-License-Identifier: MPL-2.0
#
# success.py provides a utility function to succeed gracefully and
# return the result as JSON.

def success(output):
    import json
    print(json.dumps(output), flush=True)
    exit(0)
