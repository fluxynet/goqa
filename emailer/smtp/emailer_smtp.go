package smtp

import (
	"bytes"
	"net/smtp"

	"github.com/fluxynet/goqa"
)

func init() {
	var _ goqa.Emailer = New("", "", "", "")
}

func New(host, usr, pass, from string) Smtp {
	return Smtp{
		host: host,
		usr:  usr,
		pass: pass,
		from: from,
	}
}

type Smtp struct {
	host string
	usr  string
	pass string
	from string
}

func (s Smtp) Send(subject string, message string, recipients ...string) error {
	var (
		auth = smtp.PlainAuth("", s.usr, s.pass, s.host)
		msg  bytes.Buffer
		err  error
	)

	// todo make this parallel maybe?
	for i := range recipients {
		msg.WriteString("To: ")
		msg.WriteString(recipients[i])
		msg.WriteString("\r\n")

		msg.WriteString("Subject: ")
		msg.WriteString(subject)
		msg.WriteString("\r\n")

		msg.WriteString("\r\n")
		msg.WriteString(message)
		msg.WriteString("\r\n")

		err = smtp.SendMail(s.host, auth, s.from, []string{recipients[i]}, msg.Bytes())

		if err != nil {
			return err
		}
	}

	return nil
}
