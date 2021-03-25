package email

import (
	"testing"

	"github.com/fluxynet/goqa"
)

type fakemailer struct {
	subject    string
	message    string
	recipients []string
}

func (f *fakemailer) Send(subject string, message string, recipients ...string) error {
	f.subject = subject
	f.message = message
	f.recipients = recipients
	return nil
}

type fakeevent struct {
	name  string
	value string
}

func (f fakeevent) Name() string {
	return f.name
}

func (f fakeevent) String() string {
	return f.value
}

func TestEmail_Notify(t *testing.T) {
	type fields struct {
		Email string
	}

	type args struct {
		event goqa.Event
	}

	type want struct {
		subject    string
		message    string
		recipients []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr error
	}{
		{
			name: "no email",
			fields: fields{
				Email: "",
			},
			args: args{
				fakeevent{
					name:  "fakename",
					value: "fakevalue",
				},
			},
			want: want{
				subject:    "",
				message:    "",
				recipients: nil,
			},
			wantErr: ErrEmailEmpty,
		},
		{
			name: "normal email",
			fields: fields{
				Email: "john@doe.com",
			},
			args: args{
				fakeevent{
					name:  "fakename",
					value: "fakevalue",
				},
			},
			want: want{
				subject:    "fakename",
				message:    "fakevalue",
				recipients: []string{"john@doe.com"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fakemailer{}
			e := Email{
				mailer: m,
				Email:  tt.fields.Email,
			}

			err := e.Notify(tt.args.event)
			if err != tt.wantErr {
				t.Errorf("err\nwant = %v\ngot  = %v", tt.wantErr, err)
				return
			}

			if m.message != tt.want.message {
				t.Errorf("message:\nwant = %s\ngot  = %s", tt.want.message, m.message)
			}

			if m.subject != tt.want.subject {
				t.Errorf("subject:\nwant = %s\ngot  = %s", tt.want.subject, m.subject)
			}

			var (
				lg = len(m.recipients)
				lw = len(tt.want.recipients)
			)

			if lg != lw {
				t.Errorf("recipients, length\nwant = %d\ngot  = %d", lw, lg)
				return
			}

			if lg == 0 {
				return
			}

			for i := range tt.want.recipients {
				if m.recipients[i] != tt.want.recipients[i] {
					t.Errorf("recipients(%d):\nwant = %s\ngot  = %s", i, tt.want.recipients[i], m.recipients[i])
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Subscriber = New(nil, "")
	})
}
