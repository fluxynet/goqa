package goqa

import (
	"reflect"
	"testing"
)

type fakesubscriber struct {
	event Event
	id    string
}

func (f fakesubscriber) String() string {
	return "(" + f.id + ")"
}

func (f fakesubscriber) ID() string {
	return f.id
}

func (f *fakesubscriber) SetID(id string) {
	f.id = id
}

func (f *fakesubscriber) Notify(event Event) error {
	f.event = event
	return nil
}

func (f *fakesubscriber) Serialize() (string, error) {
	panic("not implemented")
}

func (f *fakesubscriber) Unserialize(s string) error {
	panic("not implemented")
}

type fakeevent struct {
	name string
}

func (f fakeevent) Name() string {
	return f.name
}

func (f fakeevent) String() string {
	return f.name + ".data"
}

func TestPublish(t *testing.T) {
	type args struct {
		event Event
		subs  []Subscriber
	}

	tests := []struct {
		name string
		args args
		want []Subscriber
	}{
		{
			name: "no subscriber",
			args: args{
				event: fakeevent{
					name: "event 1",
				},
				subs: nil,
			},
		},
		{
			name: "nil event",
			args: args{
				event: nil,
				subs: []Subscriber{
					&fakesubscriber{id: "foo"},
				},
			},
		},
		{
			name: "1 subscriber",
			args: args{
				event: fakeevent{
					name: "event a",
				},
				subs: []Subscriber{
					&fakesubscriber{id: "foo"},
				},
			},
		},
		{
			name: "many subscribers",
			args: args{
				event: fakeevent{
					name: "event bar",
				},
				subs: []Subscriber{
					&fakesubscriber{id: "foo1"},
					&fakesubscriber{id: "foo2"},
					&fakesubscriber{id: "foo3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Publish(tt.args.event, tt.args.subs...)

			for i := range tt.args.subs {
				s, ok := tt.args.subs[i].(*fakesubscriber)
				if !ok {
					t.Errorf("%d is not a fakesubscriber:[ %T] %v", i, tt.args.subs[i], tt.args.subs[i])
					continue
				}

				if !reflect.DeepEqual(s.event, tt.args.event) {
					t.Errorf("%d event\nwant = %v\ngot  = %v", i, tt.args.event, s.event)
				}
			}
		})
	}
}
