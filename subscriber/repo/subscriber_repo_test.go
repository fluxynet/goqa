package repo

import (
	"context"
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

type fakerepo struct {
	covs []goqa.Coverage
}

func (f *fakerepo) Save(ctx context.Context, covs ...goqa.Coverage) error {
	f.covs = covs
	return nil
}

func (f *fakerepo) Load(ctx context.Context) ([]goqa.Coverage, error) {
	return f.covs, nil
}

func (f *fakerepo) Close() error {
	return nil
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Subscriber = New(nil)
	})
}

func TestRepo_Notify(t *testing.T) {
	type args struct {
		event goqa.Event
	}

	tests := []struct {
		name    string
		args    args
		repo    goqa.Repo
		want    []goqa.Coverage
		wantErr error
	}{
		{
			name: "empty and nil event",
			args: args{
				event: nil,
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "non empty and nil event",
			args: args{
				event: nil,
			},
			repo: &fakerepo{
				covs: []goqa.Coverage{
					{
						Pkg:        "pkg a",
						Percentage: 10,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
					{
						Pkg:        "pkg b",
						Percentage: 20,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "pkg a",
					Percentage: 10,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg b",
					Percentage: 20,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "empty and coverage event",
			args: args{
				event: goqa.CoverageEvent{
					Pkg:        "foo",
					Percentage: 70,
					Time:       "2006-01-02-T15:04:05Z07:00",
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "non empty and coverage event",
			args: args{
				event: goqa.CoverageEvent{
					Pkg:        "foo",
					Percentage: 70,
					Time:       "2006-01-02-T15:04:05Z07:00",
				},
			},
			repo: &fakerepo{
				covs: []goqa.Coverage{
					{
						Pkg:        "pkg 1",
						Percentage: 50,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
					{
						Pkg:        "pkg 2",
						Percentage: 30,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
				},
			},
			wantErr: subscriber.ErrUnsupportedEvent,
		},
		{
			name: "empty and github event",
			repo: &fakerepo{},
			args: args{
				event: goqa.GithubEvent{
					Event:      "Event A",
					Repository: "Repository A",
					Commit:     "Commit A",
					Ref:        "Ref A",
					Head:       "Head A",
					Workflow:   "Workflow A",
					Coverage: []goqa.Coverage{
						{
							Pkg:        "pkg foo",
							Percentage: 55,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg bar",
							Percentage: 35,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "pkg foo",
					Percentage: 55,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg bar",
					Percentage: 35,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: nil,
		},
		{
			name: "non empty and github event",
			repo: &fakerepo{
				covs: []goqa.Coverage{
					{
						Pkg:        "foobar",
						Percentage: 15,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
				},
			},
			args: args{
				event: goqa.GithubEvent{
					Event:      "Event B",
					Repository: "Repository B",
					Commit:     "Commit B",
					Ref:        "Ref B",
					Head:       "Head B",
					Workflow:   "Workflow B",
					Coverage: []goqa.Coverage{
						{
							Pkg:        "pkg abc",
							Percentage: 95,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
						{
							Pkg:        "pkg def",
							Percentage: 85,
							Time:       "2006-01-02T15:04:05Z07:00",
						},
					},
				},
			},
			want: []goqa.Coverage{
				{
					Pkg:        "pkg abc",
					Percentage: 95,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
				{
					Pkg:        "pkg def",
					Percentage: 85,
					Time:       "2006-01-02T15:04:05Z07:00",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repo{
				repo: tt.repo,
			}

			err := r.Notify(tt.args.event)
			if err != tt.wantErr {
				t.Errorf("Notify() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
