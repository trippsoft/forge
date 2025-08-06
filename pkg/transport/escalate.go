package transport

var (
	_ EscalateConfig = &NoEscalate{}
	_ EscalateConfig = &SimpleEscalate{}
	_ EscalateConfig = &Impersonate{}
)

type EscalateConfig interface {
	Enabled() bool // Whether escalation is enabled.
	User() string  // The username to impersonate, if applicable.
	Pass() string  // The password for the user, if applicable.
}

type NoEscalate struct{}

// Enabled implements EscalateConfig.
func (n *NoEscalate) Enabled() bool {
	return false
}

// User implements EscalateConfig.
func (n *NoEscalate) User() string {
	return ""
}

// Pass implements EscalateConfig.
func (n *NoEscalate) Pass() string {
	return ""
}

type SimpleEscalate struct {
	Password string
}

// Enabled implements EscalateConfig.
func (s *SimpleEscalate) Enabled() bool {
	return true
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

// Enabled implements EscalateConfig.
func (i *Impersonate) Enabled() bool {
	return true
}

// User implements EscalateConfig.
func (i *Impersonate) User() string {
	return i.Username
}

// Pass implements EscalateConfig.
func (i *Impersonate) Pass() string {
	return i.Password
}
