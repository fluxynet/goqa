package memory

import (
	"context"
	"strconv"
	"sync"

	"github.com/fluxynet/goqa"
)

func init() {
	var _ goqa.Roster = New()
}

type Memory struct {
	subs map[string]goqa.Subscriber
	id   int
	mut  sync.RWMutex
}

func New() *Memory {
	return &Memory{subs: make(map[string]goqa.Subscriber)}
}

func (s *Memory) Subscribe(ctx context.Context, name string, sub goqa.Subscriber) error {
	defer s.mut.Unlock()
	s.mut.Lock()

	s.id++

	var id = name + "-" + strconv.Itoa(s.id) // poor man's uuid
	s.subs[id] = sub

	sub.SetID(id)

	return nil
}

func (s *Memory) Unsubscribe(ctx context.Context, id string) error {
	defer s.mut.Unlock()
	s.mut.Lock()

	delete(s.subs, id)
	return nil
}

func (s *Memory) Subscribers(ctx context.Context, name string) ([]goqa.Subscriber, error) {
	if s.subs == nil {
		return nil, nil
	}

	defer s.mut.RUnlock()
	s.mut.RLock()

	var (
		subs = make([]goqa.Subscriber, len(s.subs))
		i    int
	)

	for k := range s.subs {
		subs[i] = s.subs[k]
		i++
	}

	return subs, nil
}

func (s *Memory) Close() error {
	return nil
}
