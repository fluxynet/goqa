package email

import (
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
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

func TestEmail_Serialize(t *testing.T) {
	type fields struct {
		Email string
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name:    "empty email",
			fields:  fields{},
			want:    "",
			wantErr: subscriber.ErrSerializeNotSupported,
		},
		{
			name: "non-empty email",
			fields: fields{
				Email: "john@doe.com",
			},
			want:    "john@doe.com",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Email{
				Email: tt.fields.Email,
			}

			got, err := e.Serialize()
			if err != tt.wantErr {
				t.Errorf("err\nwant = %v\ngot  = %v", tt.wantErr, err)
				return
			}

			if got != tt.want {
				t.Errorf("Serialize() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestEmail_Unserialize(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: subscriber.ErrSerializeNotSupported,
		},
		{
			name: "non-empty",
			args: args{
				s: "foo@bar.com",
			},
			want: "foo@bar.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{}

			err := e.Unserialize(tt.args.s)

			if err != tt.wantErr {
				t.Errorf("err\nwant = %v\ngot  = %v", tt.wantErr, err)
				return
			}

			if e.Email != tt.want {
				t.Errorf("Serialize() got = %s, want %s", e.Email, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Subscriber = New(nil, "")
	})
}
