package transport

type EscalateConfig interface {
	User() string // The username to impersonate, if applicable.
	Pass() string // The password for the user, if applicable.
}

type escalateConfig struct {
	username string
	password string
}

// NewNoPasswordEscalation creates a new escalateConfig that escalates privileges,
// does not impersonate a user, and does not require a password.
func NewNoPasswordEscalation() EscalateConfig {
	return &escalateConfig{
		username: "",
		password: "",
	}
}

// NewEscalation creates a new escalateConfig that just escalates privileges,
// does not impersonate a user, and requires a password.
func NewEscalation(password string) EscalateConfig {
	return &escalateConfig{
		username: "",
		password: password,
	}
}

// NewNoPasswordImpersonation creates a new escalateConfig that impersonates the given user
// and does not require a password.
func NewNoPasswordImpersonation(username string) EscalateConfig {
	return &escalateConfig{
		username: username,
		password: "",
	}
}

// NewImpersonation creates a new escalateConfig that impersonates the given user.
func NewImpersonation(username, password string) EscalateConfig {
	return &escalateConfig{
		username: username,
		password: password,
	}
}

// User implements EscalateConfig.
func (i *escalateConfig) User() string {
	return i.username
}

// Pass implements EscalateConfig.
func (i *escalateConfig) Pass() string {
	return i.password
}
