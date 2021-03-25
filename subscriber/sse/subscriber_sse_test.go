package sse

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fluxynet/goqa"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Subscriber = New(nil, nil)
	})
}

type fakeevent struct {
	name string
}

func (f fakeevent) Name() string {
	return f.name
}

func (f fakeevent) String() string {
	return f.name + "::data"
}

func TestSSE_Notify(t *testing.T) {
	tests := []struct {
		name   string
		events []goqa.Event
		want   string
	}{
		{
			name: "1 event",
			events: []goqa.Event{
				fakeevent{
					name: "foo",
				},
			},
			want: "event: foo\ndata: foo::data\n\n",
		},
		{
			name: "3 events",
			events: []goqa.Event{
				fakeevent{
					name: "event 1",
				},
				fakeevent{
					name: "event 2",
				},
				fakeevent{
					name: "event 3",
				},
			},
			want: "event: event 1\ndata: event 1::data\n\nevent: event 2\ndata: event 2::data\n\nevent: event 3\ndata: event 3::data\n\n",
		},
		{
			name: "event without name",
			events: []goqa.Event{
				fakeevent{
					name: "",
				},
			},
			want: "data: ::data\n\n",
		},
		{
			name: "event with new line",
			events: []goqa.Event{
				fakeevent{
					name: "foo\nbar",
				},
			},
			want: "event: foo_bar\ndata: foo_bar::data\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			s := &SSE{
				ctx:     context.Background(),
				flusher: w,
				writer:  w,
			}

			for i := range tt.events {
				err := s.Notify(tt.events[i])
				if err != nil {
					t.Errorf("error not nil = %v", err)
					return
				}
			}

			b := w.Body.String()
			if b != tt.want {
				r := strings.NewReplacer("\r", "[R]", "\n", "[N]")
				t.Errorf("body not same\nwant = %s\ngot  = %s", r.Replace(tt.want), r.Replace(b))
			}

		})
	}
}
