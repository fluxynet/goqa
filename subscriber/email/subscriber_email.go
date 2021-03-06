package email

import (
	"errors"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

var (
	// ErrEmailEmpty email is not good
	ErrEmailEmpty = errors.New("email cannot be empty")
)

func New(mailer goqa.Emailer, email string) *Email {
	return &Email{mailer: mailer, Email: email}
}

type Email struct {
	subscriber.Identifiable
	mailer goqa.Emailer
	Email  string
}

func (e Email) Notify(event goqa.Event) error {
	if e.Email == "" {
		return ErrEmailEmpty
	}

	return e.mailer.Send(event.Name(), event.String(), e.Email)
}
