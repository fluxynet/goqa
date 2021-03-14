package hook

import (
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
)

func TestCreateGithubEvent(t *testing.T) {
	type args struct {
		p *Payload
	}

	tests := []struct {
		name string
		args args
		want *goqa.GithubEvent
	}{
		{
			name: "nil",
			args: args{
				p: nil,
			},
			want: nil,
		},
		{
			name: "no datum",
			args: args{
				p: &Payload{
					Event:      "Event A",
					Repository: "Repo A",
					Commit:     "Commit A",
					Ref:        "Ref A",
					Head:       "Head A",
					Workflow:   "Workflow A",
					Data:       nil,
				},
			},
			want: &goqa.GithubEvent{
				Event:      "Event A",
				Repository: "Repo A",
				Commit:     "Commit A",
				Ref:        "Ref A",
				Head:       "Head A",
				Workflow:   "Workflow A",
				Coverage:   nil,
			},
		},
		{
			name: "1 datum no coverage",
			args: args{
				p: &Payload{
					Event:      "Event A",
					Repository: "Repo A",
					Commit:     "Commit A",
					Ref:        "Ref A",
					Head:       "Head A",
					Workflow:   "Workflow A",
					Data: []Datum{
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 1",
							Package: "Package 1",
							Test:    "Test 1",
							Output:  "Output 1",
						},
					},
				},
			},
			want: &goqa.GithubEvent{
				Event:      "Event A",
				Repository: "Repo A",
				Commit:     "Commit A",
				Ref:        "Ref A",
				Head:       "Head A",
				Workflow:   "Workflow A",
				Coverage:   nil,
			},
		},
		{
			name: "1 datum with coverage",
			args: args{
				p: &Payload{
					Event:      "Event A",
					Repository: "Repository A",
					Commit:     "Commit A",
					Ref:        "Ref A",
					Head:       "Head A",
					Workflow:   "Workflow A",
					Data: []Datum{
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 1",
							Package: "Package 1",
							Test:    "Test 1",
							Output:  "coverage: 10.5% of statements\n",
						},
					},
				},
			},
			want: &goqa.GithubEvent{
				Event:      "Event A",
				Repository: "Repository A",
				Commit:     "Commit A",
				Ref:        "Ref A",
				Head:       "Head A",
				Workflow:   "Workflow A",
				Coverage: []goqa.Coverage{
					{Pkg: "Package 1", Percentage: 10, Time: "2006-01-02T15:04:05Z07:00"},
				},
			},
		},
		{
			name: "3 data, 2 coverages",
			args: args{
				p: &Payload{
					Event:      "Event Foo",
					Repository: "Repository Foo",
					Commit:     "Commit Foo",
					Ref:        "Ref Foo",
					Head:       "Head Foo",
					Workflow:   "Workflow Foo",
					Data: []Datum{
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 1",
							Package: "Package 1",
							Test:    "Test 1",
							Output:  "coverage: 5% of statements\n",
						},
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 2",
							Package: "Package 2",
							Test:    "Test 2",
							Output:  "Output 2",
						},
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 3",
							Package: "Package 3",
							Test:    "Test 3",
							Output:  "coverage: 7.4% of statements\n",
						},
					},
				},
			},
			want: &goqa.GithubEvent{
				Event:      "Event Foo",
				Repository: "Repository Foo",
				Commit:     "Commit Foo",
				Ref:        "Ref Foo",
				Head:       "Head Foo",
				Workflow:   "Workflow Foo",
				Coverage: []goqa.Coverage{
					{
						Pkg:        "Package 1",
						Percentage: 5,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
					{
						Pkg:        "Package 3",
						Percentage: 7,
						Time:       "2006-01-02T15:04:05Z07:00",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateGithubEvent(tt.args.p)

			internal.AssertGithubEventsEqual(t, got, tt.want)
		})
	}
}
