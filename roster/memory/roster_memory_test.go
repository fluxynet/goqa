package memory

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
)

type fakesubscriber struct {
	id string
}

func (f fakesubscriber) String() string {
	return "(" + f.id + ")"
}

func (f *fakesubscriber) ID() string {
	return f.id
}

func (f *fakesubscriber) SetID(id string) {
	f.id = id
}

func (f *fakesubscriber) Notify(event goqa.Event) error {
	panic("not supported")
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Roster = New()
	})
}

func TestMemory_Close(t *testing.T) {
	t.Run("close", func(t *testing.T) {
		s := &Memory{
			subs: map[string]map[string]goqa.Subscriber{},
			mut:  sync.Mutex{},
		}
		err := s.Close()

		if err != nil {
			t.Errorf("Close() error not nil = %v", err)
		}
	})
}

func TestMemory_Subscribers(t *testing.T) {
	type fields struct {
		subs map[string]map[string]goqa.Subscriber
	}

	type args struct {
		name string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []goqa.Subscriber
	}{
		{
			name: "nil",
			fields: fields{
				subs: nil,
			},
			args: args{
				name: "foo",
			},
			want: nil,
		},
		{
			name: "1 event",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
				},
			},
			args: args{
				name: "foo",
			},
			want: []goqa.Subscriber{
				&fakesubscriber{id: "foo-1"},
				&fakesubscriber{id: "foo-2"},
				&fakesubscriber{id: "foo-3"},
			},
		},
		{
			name: "2 events",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
					"bar": {
						"bar-2": &fakesubscriber{id: "bar-2"},
						"bar-4": &fakesubscriber{id: "bar-4"},
					},
				},
			},
			args: args{
				name: "bar",
			},
			want: []goqa.Subscriber{
				&fakesubscriber{id: "bar-2"},
				&fakesubscriber{id: "bar-4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Memory{
				subs: tt.fields.subs,
				mut:  sync.Mutex{},
			}

			got, err := s.Subscribers(context.Background(), tt.args.name)
			if err != nil {
				t.Errorf("Subscribers() error not nil = %v", err)
				return
			}

			sort.SliceStable(got, func(i, j int) bool {
				return strings.Compare(got[i].ID(), got[j].ID()) == 1
			})

			sort.SliceStable(tt.want, func(i, j int) bool {
				return strings.Compare(tt.want[i].ID(), tt.want[j].ID()) == 1
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscribers() want = %v\ngot  = %v", tt.want, got)
			}

			internal.AssertMutexUnlocked(t, &s.mut)
		})
	}
}

func TestMemory_Unsubscribe(t *testing.T) {
	type fields struct {
		subs map[string]map[string]goqa.Subscriber
		id   int
	}

	type args struct {
		id string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]map[string]goqa.Subscriber
		wantID int
	}{
		{
			name: "nil",
			fields: fields{
				subs: nil,
				id:   0,
			},
			args: args{
				id: "foo-1",
			},
			want:   nil,
			wantID: 0,
		},
		{
			name: "1 event, bad id",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
				},
				id: 3,
			},
			args: args{
				id: "foo",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-2": &fakesubscriber{id: "foo-2"},
					"foo-3": &fakesubscriber{id: "foo-3"},
				},
			},
			wantID: 3,
		},
		{
			name: "1 event, good id",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
				},
				id: 3,
			},
			args: args{
				id: "foo-2",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-3": &fakesubscriber{id: "foo-3"},
				},
			},
			wantID: 3,
		},
		{
			name: "2 events",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
					"bar": {
						"bar-2": &fakesubscriber{id: "bar-2"},
						"bar-4": &fakesubscriber{id: "bar-4"},
					},
				},
				id: 4,
			},
			args: args{
				id: "bar-4",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-3": &fakesubscriber{id: "foo-3"},
				},
				"bar": {
					"bar-2": &fakesubscriber{id: "bar-2"},
				},
			},
			wantID: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Memory{
				subs: tt.fields.subs,
				id:   tt.fields.id,
				mut:  sync.Mutex{},
			}

			err := s.Unsubscribe(context.Background(), tt.args.id)
			if err != nil {
				t.Errorf("Subscribers() error not nil = %v", err)
				return
			}

			if tt.wantID != s.id {
				t.Errorf("id\nwant = %d\ngot  = %d", tt.wantID, s.id)
			}

			if !reflect.DeepEqual(s.subs, tt.want) {
				t.Errorf("Subscribers() got = %v\nwant %v", s.subs, tt.want)
			}

			internal.AssertMutexUnlocked(t, &s.mut)
		})
	}
}

func TestMemory_Subscribe(t *testing.T) {
	type fields struct {
		subs map[string]map[string]goqa.Subscriber
		id   int
	}

	type args struct {
		name string
		sub  goqa.Subscriber
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]map[string]goqa.Subscriber
		wantID int
	}{
		{
			name: "nil",
			fields: fields{
				subs: nil,
				id:   0,
			},
			args: args{
				sub:  &fakesubscriber{},
				name: "foo",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
				},
			},
			wantID: 1,
		},
		{
			name: "1 event, nil subscriber",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
					},
				},
				id: 2,
			},
			args: args{
				sub:  nil,
				name: "foo",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-2": &fakesubscriber{id: "foo-2"},
				},
			},
			wantID: 2,
		},
		{
			name: "1 event, empty event",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
					},
				},
				id: 2,
			},
			args: args{
				sub:  &fakesubscriber{id: "????"},
				name: "",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-2": &fakesubscriber{id: "foo-2"},
				},
			},
			wantID: 2,
		},
		{
			name: "1 event, good subscriber",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-2": &fakesubscriber{id: "foo-2"},
					},
				},
				id: 2,
			},
			args: args{
				sub:  &fakesubscriber{id: "???"},
				name: "foo",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-2": &fakesubscriber{id: "foo-2"},
					"foo-3": &fakesubscriber{id: "foo-3"},
				},
			},
			wantID: 3,
		},
		{
			name: "2 events",
			fields: fields{
				subs: map[string]map[string]goqa.Subscriber{
					"foo": {
						"foo-1": &fakesubscriber{id: "foo-1"},
						"foo-3": &fakesubscriber{id: "foo-3"},
					},
					"bar": {
						"bar-2": &fakesubscriber{id: "bar-2"},
					},
				},
				id: 3,
			},
			args: args{
				sub:  &fakesubscriber{id: "???"},
				name: "bar",
			},
			want: map[string]map[string]goqa.Subscriber{
				"foo": {
					"foo-1": &fakesubscriber{id: "foo-1"},
					"foo-3": &fakesubscriber{id: "foo-3"},
				},
				"bar": {
					"bar-2": &fakesubscriber{id: "bar-2"},
					"bar-4": &fakesubscriber{id: "bar-4"},
				},
			},
			wantID: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Memory{
				subs: tt.fields.subs,
				id:   tt.fields.id,
				mut:  sync.Mutex{},
			}

			err := s.Subscribe(context.Background(), tt.args.name, tt.args.sub)
			if err != nil {
				t.Errorf("Subscribers() error not nil = %v", err)
				return
			}

			if tt.wantID != s.id {
				t.Errorf("id\nwant = %d\ngot  = %d", tt.wantID, s.id)
			}

			if !reflect.DeepEqual(s.subs, tt.want) {
				t.Errorf("Subscribers() got = %v\nwant %v", s.subs, tt.want)
			}

			internal.AssertMutexUnlocked(t, &s.mut)
		})
	}
}
