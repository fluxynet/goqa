package flat

import (
	"context"
	"encoding/json"
	"os"

	"github.com/fluxynet/goqa"
)

const (
	filename = "goqa.repo.json"
)

type filewriterFunc func(name string, data []byte, perm os.FileMode) error
type filereaderFunc func(name string) ([]byte, error)

var (
	filewriter filewriterFunc = os.WriteFile
	filereader filereaderFunc = os.ReadFile
)

type Flat struct{}

func New() *Flat {
	return &Flat{}
}

func (f Flat) Save(ctx context.Context, covs ...goqa.Coverage) error {
	var (
		b   []byte
		err error
	)

	if len(covs) == 0 {
		b = []byte("[]")
	} else {
		b, err = json.Marshal(covs)
	}

	if err != nil {
		return err
	}

	return filewriter(filename, b, 0644)
}

func (f Flat) Load(ctx context.Context) ([]goqa.Coverage, error) {
	var b, err = filereader(filename)

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
