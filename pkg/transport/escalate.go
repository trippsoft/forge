package transport

var (
	_ EscalateConfig = &SimpleEscalate{}
	_ EscalateConfig = &Impersonate{}
)

type EscalateConfig interface {
	User() string // The username to impersonate, if applicable.
	Pass() string // The password for the user, if applicable.
}

type SimpleEscalate struct {
	Password string
}

// User implements EscalateConfig.
func (s *SimpleEscalate) User() string {
	return "" // Use default user (root on most POSIX or SYSTEM on Windows).
}

// Pass implements EscalateConfig.
func (s *SimpleEscalate) Pass() string {
	return s.Password
}

type Impersonate struct {
	Username string
	Password string
}

// User implements EscalateConfig.
func (i *Impersonate) User() string {
	return i.Username
}

// Pass implements EscalateConfig.
func (i *Impersonate) Pass() string {
	return i.Password
}
