package goqa

// Emailer sends emails for us
type Emailer interface {
	// Send the mail!
	Send(subject string, message string, recipients ...string) error
}
