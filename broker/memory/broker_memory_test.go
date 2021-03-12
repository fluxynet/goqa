package memory

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
)

const (
	// EventDummy is a dummy event
	EventDummy = "DUMMY_EVENT"

	// ContentDummy is the content of the dummy event
	ContentDummy = "<dummy>"
)

// DummyEvent is a not for real
type DummyEvent struct{}

func (d DummyEvent) Name() string {
	return EventDummy
}

func (d DummyEvent) String() string {
	return ContentDummy
}

func makeListeners(t int) []chan goqa.Event {
	var c = make([]chan goqa.Event, t)
	for i := range c {
		c[i] = make(chan goqa.Event)
	}

	return c
}

func TestMemory_Close(t *testing.T) {
	type fields struct {
		listeners []chan goqa.Event
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				listeners: nil,
			},
			wantErr: false,
		},
		{
			name: "non-empty",
			fields: fields{
				listeners: makeListeners(10),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				listeners: tt.fields.listeners,
				mutex:     sync.Mutex{},
			}

			l := append([]chan goqa.Event{}, tt.fields.listeners...)

			if err := m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i := range l {
				select {
				case <-l[i]:
				case <-time.After(time.Millisecond):
					t.Errorf("channel %d not closed", i)
				}
			}

			internal.AssertMutexUnlocked(t, &m.mutex)
		})
	}
}

func TestMemory_Listen(t *testing.T) {
	type fields struct {
		listeners []chan goqa.Event
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
		{
			name: "non empty",
			fields: fields{
				listeners: makeListeners(50),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				listeners: tt.fields.listeners,
				mutex:     sync.Mutex{},
			}

			before := len(m.listeners)

			got, err := m.Listen(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			after := len(m.listeners)

			if after-before != 1 {
				t.Errorf("length did not increase by 1, diff = %d", after-before)
				return
			}

			if m.listeners[before] != got {
				t.Errorf("channel returned is different from the one created")
				return
			}

			internal.AssertMutexUnlocked(t, &m.mutex)
		})
	}
}

func TestMemory_Publish(t *testing.T) {
	type fields struct {
		listeners []chan goqa.Event
	}

	type args struct {
		ctx   context.Context
		event goqa.Event
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "no listener",
			fields: fields{},
			args: args{
				ctx:   context.Background(),
				event: DummyEvent{},
			},
			wantErr: false,
		},
		{
			name: "2 listeners",
			fields: fields{
				listeners: makeListeners(2),
			},
			args: args{
				ctx:   context.Background(),
				event: DummyEvent{},
			},
			wantErr: false,
		},
		{
			name: "99999 listeners",
			fields: fields{
				listeners: makeListeners(99999),
			},
			args: args{
				ctx:   context.Background(),
				event: DummyEvent{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				listeners: tt.fields.listeners,
				mutex:     sync.Mutex{},
			}

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				if err := m.Publish(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
					t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				}

				internal.AssertMutexUnlocked(t, &m.mutex)
			}()

			for i := range m.listeners {
				select {
				case ev := <-m.listeners[i]:
					if !reflect.DeepEqual(tt.args.event, ev) {
						t.Errorf(
							"event different at %d\nwant: [%s] %s\ngot:  [%s] [%s]\n",
							i,
							tt.args.event.Name(), tt.args.event.String(),
							ev.Name(), ev.String(),
						)
						return
					}
				case <-time.After(time.Millisecond):
					t.Errorf("event not received!")
					return
				}
			}

			wg.Wait()
		})
	}
}

func TestNew(t *testing.T) {
	var _ goqa.Broker = New()
}
