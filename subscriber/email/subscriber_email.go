package email

import (
	"github.com/fluxynet/goqa"
)

func init() {
	var _ goqa.Subscriber = New(nil, "")
}

func New(mailer goqa.Emailer, email string) *Email {
	return &Email{mailer: mailer, Email: email}
}

type Email struct {
	goqa.Identifiable
	mailer goqa.Emailer
	Email  string
}

func (e Email) Notify(event goqa.Event) error {
	return e.mailer.Send(event.Name(), event.String(), e.Email)
}

func (e Email) Serialize() (string, error) {
	return e.Email, nil
}

func (e *Email) Unserialize(s string) error {
	e.Email = s
	return nil
}
