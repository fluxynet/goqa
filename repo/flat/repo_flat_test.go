package flat

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/fluxynet/goqa"
)

var (
	errWrite = errors.New("write err")
	errRead  = errors.New("read err")
)

func assertCoveragesEqual(t *testing.T, got, want []goqa.Coverage) {
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

	mgot := make(map[string]goqa.Coverage, lg)
	for i := range got {
		mgot[got[i].Pkg] = got[i]
	}

	mwant := make(map[string]goqa.Coverage, lw)
	for i := range want {
		mwant[want[i].Pkg] = want[i]
	}

	for k, w := range mwant {
		var (
			g    = mgot[k]
			same = g.Percentage == w.Percentage && g.Pkg == w.Pkg && g.Time == w.Time
		)

		if !same {
			t.Errorf(
				"not equal [%s]:\ngot\ntime = %s pkg = %s perc = %d\nwant\ntime = %s pkg = %s perc = %d",
				k,
				g.Time, g.Pkg, g.Percentage,
				w.Time, w.Pkg, w.Percentage,
			)
			return
		}
	}
}

type fakereadwriter struct {
	filename string
	data     []byte
	perm     os.FileMode
}

func (f *fakereadwriter) write(name string, data []byte, perm os.FileMode) error {
	f.filename = name
	f.data = data
	f.perm = perm
	return nil
}

func (f *fakereadwriter) read(name string) ([]byte, error) {
	return f.data, nil
}

func (f *fakereadwriter) String() string {
	return fmt.Sprintf("filename = %s perm = %O data = \n%s", f.filename, f.perm, f.data)
}

type badReadWriter struct {
	err error
}

func (f badReadWriter) write(name string, data []byte, perm os.FileMode) error {
	return f.err
}

func (f badReadWriter) read(name string) ([]byte, error) {
	return nil, f.err
}

func assertFakereadwritersEqual(t *testing.T, got, want fakereadwriter) {
	if got.filename != want.filename {
		t.Errorf("filename; got = %s, want = %s", got.filename, want.filename)
	}

	if got.perm != want.perm {
		t.Errorf("perm; got = %O, want = %O", got.perm, want.perm)
	}

	lg := len(got.data)
	lw := len(want.data)

	if lg != lw {
		t.Errorf("len data; got = %d, want = %d", lg, lw)
		return
	}

	if lg == 0 {
		return
	}

	if bytes.Compare(got.data, want.data) == 0 {
		return
	}

	var replacer = strings.NewReplacer("\r", "[R]", "\n", "[N]", " ", "[S]")
	t.Errorf("data;\nwant = %s\ngot  = %s", replacer.Replace(string(want.data)), replacer.Replace(string(got.data)))
}

func writeErr(name string, data []byte, perm os.FileMode) error {
	return errors.New("error writing")
}

func readErr(name string) ([]byte, error) {
	return nil, errors.New("error reading")
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

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var _ goqa.Repo = New()
	})
}

func TestFlat_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "close",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flat{}
			if err := f.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlat_LoadSave(t *testing.T) {
	var (
		oldfilewriter = filewriter
		oldfilereader = filereader
	)

	defer func() {
		filewriter = oldfilewriter
		filereader = oldfilereader
	}()

	type args struct {
		covs []goqa.Coverage
	}

	tests := []struct {
		name       string
		args       args
		want       fakereadwriter
		mustNotErr bool
	}{
		{
			name: "empty",
			args: args{
				covs: makeCoverages(0, 0),
			},
		},
		{
			name: "1 coverage",
			args: args{
				covs: makeCoverages(1, 0),
			},
		},
		{
			name: "many coverages",
			args: args{
				covs: makeCoverages(100, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flat{}

			faker := &fakereadwriter{}
			filewriter = faker.write
			filereader = faker.read

			err := f.Save(context.Background(), tt.args.covs...)

			if err != nil {
				t.Errorf("writer err is not nil, %s\n%s", err.Error(), faker.String())
				return
			}

			covs, err := f.Load(context.Background())
			if err != nil {
				t.Errorf("reader err is not nil, %s\n%s", err.Error(), faker.String())
				return
			}

			assertCoveragesEqual(t, covs, tt.args.covs)
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flat{}
			faker := badReadWriter{err: errRead}
			gfaker := fakereadwriter{}

			filewriter = gfaker.write
			filereader = faker.read

			err := f.Save(context.Background(), tt.args.covs...)

			if err != nil {
				t.Errorf("writer err is not nil, %s\n", err.Error())
				return
			}

			covs, err := f.Load(context.Background())

			if covs != nil {
				t.Errorf("covs is not nil, len = %d", len(covs))
			}

			if tt.mustNotErr && err != nil {
				t.Errorf("%s must not err, but gave an error anyway", tt.name)
				return
			}

			if !tt.mustNotErr && err == nil {
				t.Errorf("%s did not give an error as expected", tt.name)
				return
			}

			if !errors.Is(err, errRead) {
				t.Errorf("got an unknown error: %s", err.Error())
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flat{}
			faker := badReadWriter{err: errWrite}
			filewriter = faker.write

			err := f.Save(context.Background(), tt.args.covs...)

			if tt.mustNotErr && err != nil {
				t.Errorf("%s must not err, but gave an error anyway", tt.name)
				return
			}

			if !tt.mustNotErr && err == nil {
				t.Errorf("%s did not give an error as expected", tt.name)
				return
			}

			if !errors.Is(err, errWrite) {
				t.Errorf("got an unknown error: %s", err.Error())
			}
		})
	}
}
