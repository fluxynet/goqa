package hook

import (
	"testing"

	"github.com/fluxynet/goqa"
)

func assertGithubEventsEqual(t *testing.T, got, want *goqa.GithubEvent) {
	if (got == nil) != (want == nil) {
		t.Errorf("Nil, want = %t got = %t", want == nil, got == nil)
		return
	}

	if got == nil {
		return
	}

	if got.Event != want.Event {
		t.Errorf("Event\nwant = %s\ngot  = %s", got.Event, want.Event)
	}

	if got.Repository != want.Repository {
		t.Errorf("Repository\nwant = %s\ngot  = %s", got.Repository, want.Repository)
	}

	if got.Commit != want.Commit {
		t.Errorf("Commit\nwant = %s\ngot  = %s", got.Commit, want.Commit)
	}

	if got.Ref != want.Ref {
		t.Errorf("Ref\nwant = %s\ngot  = %s", got.Ref, want.Ref)
	}

	if got.Head != want.Head {
		t.Errorf("Head\nwant = %s\ngot  = %s", got.Head, want.Head)
	}

	if got.Workflow != want.Workflow {
		t.Errorf("Workflow\nwant = %s\ngot  = %s", got.Workflow, want.Workflow)
	}

	assertCoveragesEqual(t, got.Coverage, want.Coverage)
}

func assertCoveragesEqual(t *testing.T, got, want []goqa.Coverage) {
	var lg, lw = len(got), len(want)
	if lg != lw {
		t.Errorf("coverage length got = %d, want = %d", lg, lw)
		return
	}

	for i := range want {
		if got[i].Pkg != want[i].Pkg {
			t.Errorf("coverage(%d) pkg\nwant = %s\ngot  = %s", i, want[i].Pkg, got[i].Pkg)
		}

		if got[i].Percentage != want[i].Percentage {
			t.Errorf("coverage(%d) percentage\nwant = %d\ngot  = %d", i, want[i].Percentage, got[i].Percentage)
		}

		if got[i].Time != want[i].Time {
			t.Errorf("coverage(%d) time\nwant = %s\ngot  = %s", i, want[i].Time, got[i].Time)
		}
	}
}

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
							Elapsed: "Elapsed 1",
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
							Elapsed: "Elapsed 1",
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
							Elapsed: "Elapsed 1",
						},
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 2",
							Package: "Package 2",
							Test:    "Test 2",
							Output:  "Output 2",
							Elapsed: "Elapsed 2",
						},
						{
							Time:    "2006-01-02T15:04:05Z07:00",
							Action:  "Action 3",
							Package: "Package 3",
							Test:    "Test 3",
							Output:  "coverage: 7.4% of statements\n",
							Elapsed: "Elapsed 3",
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

			assertGithubEventsEqual(t, got, tt.want)
		})
	}
}
