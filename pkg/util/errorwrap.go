// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package util

// UnwrapErrors unwraps any error messages recursively that were joined using errors.Join(...)
func UnwrapErrors(e error) []error {
	if e == nil {
		return nil
	}

	wrapped, ok := e.(interface{ Unwrap() []error })
	if !ok {
		return []error{e}
	}

	var errs []error
	for _, err := range wrapped.Unwrap() {
		if err == nil {
			continue
		}

		if errs == nil {
			errs = UnwrapErrors(err)
			continue
		}

		errs = append(errs, UnwrapErrors(err)...)
	}

	return errs
}
