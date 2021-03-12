package memory

import (
	"sort"
	"strconv"
	"sync"
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
)

func makeMapCoverages(t, skip int) map[string]goqa.Coverage {
	var m = make(map[string]goqa.Coverage, t)
	t = t + skip

	for i := skip; i < t; i++ {
		var pkg = "pkg." + strconv.Itoa(i)
		m[pkg] = goqa.Coverage{
			Pkg:        pkg,
			Percentage: i,
			Time:       "2006-01-02T15:04:05Z07:00",
		}
	}

	return m
}

func makeCoverages(t, skip int) []goqa.Coverage {
	var m = make([]goqa.Coverage, t)
	t = t + skip

	if t == 0 {
		return nil
	}

	var n int
	for i := skip; i < t; i++ {
		var pkg = "pkg." + strconv.Itoa(i)
		m[n] = goqa.Coverage{
			Pkg:        pkg,
			Percentage: i,
			Time:       "2006-01-02T15:04:05Z07:00",
		}

		n++
	}

	return m
}

func assertStringSlicesEqual(t *testing.T, got, want []string) {
	var lg, lw int

	if got != nil {
		lg = len(got)
	}

	if want != nil {
		lw = len(want)
	}

	if lg != lw {
		t.Errorf("Length nok:\ngot  = %d\nwant = %d", lg, lw)
		return
	}

	if lg == 0 {
		return
	}

	for k, w := range want {
		var g = got[k]
		if w != g {
			t.Errorf("not equal [%d]:\ngot  = %s\nwant = %s", k, w, g)
			return
		}
	}
}

func assertCoveragesEqual(t *testing.T, got map[string]goqa.Coverage, want []goqa.Coverage) {
	var lg, lw int

	if got != nil {
		lg = len(got)
	}

	if want != nil {
		lw = len(want)
	}

	if lg != lw {
		t.Errorf("Length nok:\ngot  = %d\nwant = %d", lg, lw)
		return
	}

	if lg == 0 {
		return
	}

	for k, w := range want {
		var (
			g    = got[w.Pkg]
			same = g.Percentage == w.Percentage && g.Pkg == w.Pkg && g.Time == w.Time
		)

		if !same {
			t.Errorf(
				"not equal [%d]:\ngot\ntime = %s pkg = %s perc = %d\nwant\ntime = %s pkg = %s perc = %d",
				k,
				g.Time, g.Pkg, g.Percentage,
				w.Time, w.Pkg, w.Percentage,
			)
			return
		}
	}
}

func TestMemory_Close(t *testing.T) {
	type fields struct {
		items map[string]goqa.Coverage
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "not empty",
			fields: fields{
				items: makeMapCoverages(100, 0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				items: tt.fields.items,
				mut:   sync.Mutex{},
			}

			if err := m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if l := len(m.items); l != 0 {
				t.Errorf("expected length = 0; got = %d", l)
				return
			}

			internal.AssertMutexUnlocked(t, &m.mut)
		})
	}
}

func TestMemory_Get(t *testing.T) {
	type fields struct {
		items map[string]goqa.Coverage
	}

	tests := []struct {
		name   string
		fields fields
		total  int
	}{
		{
			name:   "empty",
			fields: fields{},
			total:  0,
		},
		{
			name: "non-empty",
			fields: fields{
				items: makeMapCoverages(1000, 0),
			},
			total: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				items: tt.fields.items,
				mut:   sync.Mutex{},
			}

			for i := 0; i < tt.total; i++ {
				var pkg = "pkg." + strconv.Itoa(i)

				got, ok := m.Get(pkg)
				if !ok {
					t.Errorf("could not get pkg: %s", pkg)
					return
				}

				want := m.items[pkg]
				same := got.Time == want.Time && got.Pkg == want.Pkg && got.Percentage == want.Percentage
				if !same {
					t.Errorf(
						"Get(%d)\ngot:\ntime = %s pkg = %s perc = %d\nwant:\ntime = %s pkg = %s perc = %d",
						i,
						got.Time, got.Pkg, got.Percentage,
						want.Time, want.Pkg, want.Percentage,
					)
					return
				}
			}

			for i := -100; i < 0; i++ {
				var pkg = "pkg." + strconv.Itoa(i)

				_, ok := m.Get(pkg)
				if ok {
					t.Errorf("should not get pkg: %s", pkg)
					return
				}
			}

			internal.AssertMutexUnlocked(t, &m.mut)
		})
	}
}

func TestMemory_Keys(t *testing.T) {
	type fields struct {
		items map[string]goqa.Coverage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			want:    nil,
			wantErr: false,
		},
		{
			name: "non-empty",
			fields: fields{
				items: makeMapCoverages(40, 0),
			},
			want: []string{
				"pkg.0", "pkg.1", "pkg.2", "pkg.3", "pkg.4", "pkg.5", "pkg.6", "pkg.7", "pkg.8", "pkg.9",
				"pkg.10", "pkg.11", "pkg.12", "pkg.13", "pkg.14", "pkg.15", "pkg.16", "pkg.17", "pkg.18", "pkg.19",
				"pkg.20", "pkg.21", "pkg.22", "pkg.23", "pkg.24", "pkg.25", "pkg.26", "pkg.27", "pkg.28", "pkg.29",
				"pkg.30", "pkg.31", "pkg.32", "pkg.33", "pkg.34", "pkg.35", "pkg.36", "pkg.37", "pkg.38", "pkg.39",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				items: tt.fields.items,
				mut:   sync.Mutex{},
			}

			got, err := m.Keys()
			if (err != nil) != tt.wantErr {
				t.Errorf("Keys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// order is not guaranteed and also not important
			sort.Strings(got)
			sort.Strings(tt.want)

			assertStringSlicesEqual(t, got, tt.want)
			internal.AssertMutexUnlocked(t, &m.mut)
		})
	}
}

func TestNew(t *testing.T) {
	var _ goqa.Cache = New()
}

func TestMemory_Reset(t *testing.T) {
	type fields struct {
		items map[string]goqa.Coverage
	}

	type args struct {
		covs []goqa.Coverage
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "empty -> empty",
			fields:  fields{},
			args:    args{},
			wantErr: false,
		},
		{
			name:   "empty -> non-empty",
			fields: fields{},
			args: args{
				covs: makeCoverages(1000, 0),
			},
			wantErr: false,
		},
		{
			name: "non-empty -> empty",
			fields: fields{
				items: makeMapCoverages(1000, 0),
			},
			args:    args{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				items: tt.fields.items,
				mut:   sync.Mutex{},
			}

			if err := m.Reset(tt.args.covs...); (err != nil) != tt.wantErr {
				t.Errorf("Reset() error = %v, wantErr %v", err, tt.wantErr)
			}

			assertCoveragesEqual(t, m.items, tt.args.covs)

			internal.AssertMutexUnlocked(t, &m.mut)
		})
	}
}
