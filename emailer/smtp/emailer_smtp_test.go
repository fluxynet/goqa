package smtp

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"testing"

	"github.com/fluxynet/goqa"
)

var (
	errSendMail = errors.New("mail send error")
)

func assertMessagesEqual(t *testing.T, want, got []fakemsg) {
	lw := len(want)
	lg := len(got)

	if lw != lg {
		t.Errorf("length not same, want = %d, got = %d", lw, lg)
		return
	}

	var replacer = strings.NewReplacer("\r", "[R]", "\n", "[N]", " ", "[S]")

	for i := range want {
		if len(want[i].from) != len(got[i].from) {
			t.Errorf("(%d) len from: want = %d, got = %d", i, len(want[i].from), len(got[i].from))
		}

		if len(want[i].to) != len(got[i].to) {
			t.Errorf("(%d) len to: want = %d, got = %d", i, len(want[i].to), len(got[i].to))
		}

		if bytes.Compare(want[i].msg, got[i].msg) != 0 {
			t.Errorf(
				"(%d) msg not same\nwant = %s\ngot  = %s",
				i,
				replacer.Replace(string(want[i].msg)),
				replacer.Replace(string(got[i].msg)),
			)
		}
	}
}

type fakesendmail struct {
	msgs []fakemsg
}

type fakemsg struct {
	from string
	to   []string
	msg  []byte
}

func (f *fakesendmail) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	f.msgs = append(f.msgs, fakemsg{from: from, to: to, msg: msg})
	return nil
}

func SendMailError(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return errSendMail
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Emailer = New("", "", "", "")
	})
}

func TestSmtp_Send(t *testing.T) {
	var oldsendmail = sendmail

	defer func() {
		sendmail = oldsendmail
	}()

	type fields struct {
		host string
		usr  string
		pass string
		from string
	}

	type args struct {
		subject    string
		message    string
		recipients []string
	}

	type want struct {
		msgs []fakemsg
	}

	tests := []struct {
		name       string
		sendmail   sendmailFunc
		fields     fields
		args       args
		want       want
		mustNotErr bool
	}{
		{
			name:       "no recipients",
			mustNotErr: true,
			fields: fields{
				host: "smtp.foobar.com",
				usr:  "john",
				pass: "doe",
				from: "foo@bar.com",
			},
			args: args{
				subject:    "test",
				message:    "sample message",
				recipients: nil,
			},
			want: want{msgs: nil},
		},
		{
			name: "one recipient",
			fields: fields{
				host: "smtp.foobar.com",
				usr:  "john",
				pass: "doe",
				from: "foo@bar.com",
			},
			args: args{
				subject:    "test",
				message:    "sample message",
				recipients: []string{"john@doe.com"},
			},
			want: want{
				[]fakemsg{
					{
						from: "foo@bar.com",
						to:   []string{"john@doe.com"},
						msg: []byte(
							fmt.Sprintf(
								"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
								"john@doe.com",
								"test",
								"sample message",
							),
						),
					},
				},
			},
		},
		{
			name: "3 recipients",
			fields: fields{
				host: "smtp.foobar.com",
				usr:  "john",
				pass: "doe",
				from: "foo@bar.com",
			},
			args: args{
				subject: "testing",
				message: `
This is 
a longer 
message.
`,
				recipients: []string{"abc@def.com", "hello@world.com", "what@ever.com"},
			},
			want: want{
				[]fakemsg{
					{
						from: "foo@bar.com",
						to:   []string{"abc@def.com"},
						msg: []byte(
							fmt.Sprintf(
								"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
								"abc@def.com",
								"testing",
								`
This is 
a longer 
message.
`,
							),
						),
					},
					{
						from: "foo@bar.com",
						to:   []string{"hello@world.com"},
						msg: []byte(
							fmt.Sprintf(
								"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
								"hello@world.com",
								"testing",
								`
This is 
a longer 
message.
`,
							),
						),
					},
					{
						from: "foo@bar.com",
						to:   []string{"what@ever.com"},
						msg: []byte(
							fmt.Sprintf(
								"To: %s\r\nSubject: %s\r\n\r\n%s\r\n",
								"what@ever.com",
								"testing",
								`
This is 
a longer 
message.
`,
							),
						),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Smtp{
				host: tt.fields.host,
				usr:  tt.fields.usr,
				pass: tt.fields.pass,
				from: tt.fields.from,
			}

			var faker = fakesendmail{}
			sendmail = faker.SendMail

			_ = s.Send(tt.args.subject, tt.args.message, tt.args.recipients...)

			assertMessagesEqual(t, tt.want.msgs, faker.msgs)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name+" error", func(t *testing.T) {
			s := Smtp{
				host: tt.fields.host,
				usr:  tt.fields.usr,
				pass: tt.fields.pass,
				from: tt.fields.from,
			}

			sendmail = SendMailError

			var err = s.Send(tt.args.subject, tt.args.message, tt.args.recipients...)
			if tt.mustNotErr && err != nil {
				t.Errorf("%s must not err, but gave an error anyway", tt.name)
				return
			}

			if !tt.mustNotErr && err == nil {
				t.Errorf("%s did not give an error as expected", tt.name)
				return
			}

			if !errors.Is(err, errSendMail) {
				t.Errorf("got an unknown error: %s", err.Error())
			}
		})
	}

}
