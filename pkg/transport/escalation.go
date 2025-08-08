package transport

// Escalation defines the interface for privilege escalation configuration for a command.
type Escalation interface {
	User() string // The username to impersonate, if applicable.
	Pass() string // The password for the user, if applicable.
}

type escalation struct {
	username string
	password string
}

// NewNoPasswordEscalation creates a new escalation that escalates privileges,
// does not impersonate a user, and does not require a password.
func NewNoPasswordEscalation() Escalation {
	return &escalation{
		username: "",
		password: "",
	}
}

// NewEscalation creates a new escalation that just escalates privileges,
// does not impersonate a user, and requires a password.
func NewEscalation(password string) Escalation {
	return &escalation{
		username: "",
		password: password,
	}
}

// NewNoPasswordImpersonation creates a new escalation that impersonates the given user
// and does not require a password.
func NewNoPasswordImpersonation(username string) Escalation {
	return &escalation{
		username: username,
		password: "",
	}
}

// NewImpersonation creates a new escalation that impersonates the given user.
func NewImpersonation(username, password string) Escalation {
	return &escalation{
		username: username,
		password: password,
	}
}

// User implements Escalation.
func (i *escalation) User() string {
	return i.username
}

// Pass implements Escalation.
func (i *escalation) Pass() string {
	return i.password
}
