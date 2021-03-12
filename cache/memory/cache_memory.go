package memory

import (
	"sync"

	"github.com/fluxynet/goqa"
)

type Memory struct {
	items map[string]goqa.Coverage
	mut   sync.Mutex
}

func New() *Memory {
	return &Memory{items: make(map[string]goqa.Coverage)}
}

func (m *Memory) Reset(covs ...goqa.Coverage) error {
	var d = make(map[string]goqa.Coverage, len(covs))
	for i := range covs {
		d[covs[i].Pkg] = covs[i]
	}

	defer m.mut.Unlock()
	m.mut.Lock()

	m.items = d

	return nil
}

func (m *Memory) Get(pkg string) (*goqa.Coverage, bool) {
	if m.items == nil {
		return nil, false
	}

	var v, ok = m.items[pkg]
	return &v, ok
}

func (m *Memory) Keys() ([]string, error) {
	if m.items == nil {
		return nil, nil
	}

	defer m.mut.Unlock()
	m.mut.Lock()

	var (
		keys = make([]string, len(m.items))
		i    int
	)

	for k := range m.items {
		keys[i] = k
		i++
	}

	return keys, nil
}

func (m *Memory) Close() error {
	defer m.mut.Unlock()
	m.mut.Lock()

	m.items = nil
	return nil
}
