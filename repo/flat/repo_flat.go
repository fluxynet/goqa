package flat

import (
	"context"
	"encoding/json"
	"os"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/repo"
)

const (
	filename = "goqa.repo.json"
)

func init() {
	var _ goqa.Repo = New()
}

type Flat struct{}

func New() *Flat {
	return &Flat{}
}

func (f Flat) Save(ctx context.Context, covs ...goqa.Coverage) error {
	var b, err = json.Marshal(covs)
	if err != nil {
		return err
	}

	var cwd string
	cwd, err = os.Getwd()
	if err != nil {
		return nil
	}

	var h *os.File
	h, err = os.CreateTemp(cwd, filename+".tmp")
	if err != nil {
		return err
	}

	var t int
	t, err = h.Write(b)

	if err != nil {
		return err
	} else if t != len(b) {
		return repo.ErrFailedToWriteCompletely
	}

	os.Rename(h.Name(), filename)

	return nil
}

func (f Flat) Load(ctx context.Context) ([]goqa.Coverage, error) {
	var b, err = os.ReadFile(filename)

	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var covs []goqa.Coverage
	err = json.Unmarshal(b, &covs)

	return covs, err
}

func (f Flat) Close() error {
	return nil
}
