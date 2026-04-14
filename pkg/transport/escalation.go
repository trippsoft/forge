// Copyright (c) Forge
// SPDX-License-Identifier: MPL-2.0

package transport

// Escalation defines the privilege escalation configuration for a transport action.
type Escalation struct {
	username string
	password string
}

// User implements Escalation.
func (i *Escalation) User() string {
	return i.username
}

// Pass implements Escalation.
func (i *Escalation) Pass() string {
	return i.password
}

// NewNoPasswordEscalation creates a new escalation that escalates privileges, does not impersonate a user, and does not
// require a password.
func NewNoPasswordEscalation() *Escalation {
	return &Escalation{
		username: "",
		password: "",
	}
}

// NewEscalation creates a new escalation that just escalates privileges, does not impersonate a user, and requires a
// password.
func NewEscalation(password string) *Escalation {
	return &Escalation{
		username: "",
		password: password,
	}
}

// NewNoPasswordImpersonation creates a new escalation that impersonates the given user and does not require a password.
func NewNoPasswordImpersonation(username string) *Escalation {
	return &Escalation{
		username: username,
		password: "",
	}
}

// NewImpersonation creates a new escalation that impersonates the given user.
func NewImpersonation(username, password string) *Escalation {
	return &Escalation{
		username: username,
		password: password,
	}
}
