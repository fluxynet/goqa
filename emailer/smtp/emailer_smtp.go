package smtp

import (
	"bytes"
	"net/smtp"
)

type sendmailFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

var (
	sendmail sendmailFunc = smtp.SendMail
)

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
		auth     = smtp.PlainAuth("", s.usr, s.pass, s.host)
		err      error
		contents bytes.Buffer
	)

	contents.WriteString("Subject: ")
	contents.WriteString(subject)
	contents.WriteString("\r\n")

	contents.WriteString("\r\n")
	contents.WriteString(message)
	contents.WriteString("\r\n")

	var body = contents.Bytes()

	// todo make this parallel maybe?
	for i := range recipients {
		var msg = append([]byte("To: "+recipients[i]+"\r\n"), body...)
		err = sendmail(s.host, auth, s.from, []string{recipients[i]}, msg)

		if err != nil {
			return err
		}
	}

	return nil
}
