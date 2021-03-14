package cachew

import (
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
	"github.com/fluxynet/goqa/subscriber"
)

type fakecache struct {
	covs []goqa.Coverage
}

func (f *fakecache) Reset(covs ...goqa.Coverage) error {
	f.covs = covs
	return nil
}

func (f *fakecache) Get(pkg string) (*goqa.Coverage, bool) {
	panic("not implemented")
	return nil, false
}

func (f *fakecache) Keys() ([]string, error) {
	panic("not implemented")
	return nil, nil
}

func (f *fakecache) Close() error {
	return nil
}

func TestCache_Notify(t *testing.T) {
	type fields struct {
		cache *fakecache
	}

	type args struct {
		event goqa.Event
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []goqa.Coverage
		wantErr error
	}{
		{
			name: "empty and nil event",
			fields: fields{
				cache: &fakecache{},
			},
			args: args{
				event: nil,
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "empty and coverage event",
			fields: fields{
				cache: &fakecache{},
			},
			args: args{
				event: goqa.CoverageEvent{
					Pkg:        "foo",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "non-empty and nil event",
			fields: fields{
				cache: &fakecache{
					covs: []goqa.Coverage{
						{
							Pkg:        "pkg1",
							Percentage: 10,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg2",
							Percentage: 20,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg3",
							Percentage: 30,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			args: args{
				event: nil,
			},
			want: []goqa.Coverage{
				{
					Pkg:        "pkg1",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg2",
					Percentage: 20,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg3",
					Percentage: 30,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "non-empty and coverage event",
			fields: fields{
				cache: &fakecache{
					covs: []goqa.Coverage{
						goqa.Coverage{
							Pkg:        "pkg1",
							Percentage: 10,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						goqa.Coverage{
							Pkg:        "pkg2",
							Percentage: 20,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						goqa.Coverage{
							Pkg:        "pkg3",
							Percentage: 30,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			args: args{
				event: goqa.CoverageEvent{
					Pkg:        "foobar",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "pkg1",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg2",
					Percentage: 20,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg3",
					Percentage: 30,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "empty and github event",
			fields: fields{
				cache: &fakecache{},
			},
			args: args{
				event: goqa.GithubEvent{
					Event:      "Event foo",
					Repository: "Repository foo",
					Commit:     "Commit foo",
					Ref:        "Ref foo",
					Head:       "Head foo",
					Workflow:   "Workflow foo",
					Coverage: []goqa.Coverage{
						{
							Pkg:        "foobar",
							Percentage: 10,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "barbaz",
							Percentage: 20,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "foobar",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "barbaz",
					Percentage: 20,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
		},
		{
			name: "non-empty and github event",
			fields: fields{
				cache: &fakecache{
					covs: []goqa.Coverage{
						{
							Pkg:        "pkg1",
							Percentage: 10,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg2",
							Percentage: 20,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg3",
							Percentage: 30,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			args: args{
				event: goqa.GithubEvent{
					Event:      "Event foo",
					Repository: "Repository foo",
					Commit:     "Commit foo",
					Ref:        "Ref foo",
					Head:       "Head foo",
					Workflow:   "Workflow foo",
					Coverage: []goqa.Coverage{
						{
							Pkg:        "foobar",
							Percentage: 10,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "barbaz",
							Percentage: 20,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "foobar",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "barbaz",
					Percentage: 20,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.fields.cache

			c := &Cache{
				cache: cache,
			}

			if err := c.Notify(tt.args.event); err != tt.wantErr {
				t.Errorf("Notify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			internal.AssertCoveragesEqual(t, cache.covs, tt.want)
		})
	}
}

func TestCache_Serialize(t *testing.T) {
	t.Run("serialize", func(t *testing.T) {
		c := &Cache{}
		got, err := c.Serialize()
		if err != subscriber.ErrSerializeNotSupported {
			t.Errorf("Serialize() error = %v, wantErr %v", err, subscriber.ErrSerializeNotSupported)
			return
		}

		if got != "" {
			t.Errorf("Serialize() got = %v, want %v", got, "")
		}
	})
}

func TestCache_Unserialize(t *testing.T) {
	t.Run("unserialize", func(t *testing.T) {
		c := &Cache{}
		err := c.Unserialize("")
		if err != subscriber.ErrSerializeNotSupported {
			t.Errorf("Serialize() error = %v, wantErr %v", err, subscriber.ErrSerializeNotSupported)
			return
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Subscriber = New(nil)
	})
}
