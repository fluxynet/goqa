package internal

import (
	"reflect"
	"sync"
	"testing"

	"github.com/fluxynet/goqa"
)

// AssertMutexUnlocked checks if a mutex is locked
func AssertMutexUnlocked(t *testing.T, m *sync.Mutex) {
	var state = reflect.ValueOf(m).Elem().FieldByName("state")
	if state.Int() == 1 {
		t.Errorf("mutex still locked")
	}
}

func AssertGithubEventsEqual(t *testing.T, got, want *goqa.GithubEvent) {
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

	AssertCoveragesEqual(t, got.Coverage, want.Coverage)
}

func AssertCoveragesEqual(t *testing.T, got, want []goqa.Coverage) {
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
