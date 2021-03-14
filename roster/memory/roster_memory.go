package memory

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/fluxynet/goqa"
)

type Memory struct {
	// subs event => [subscription_id => [subscriber]]
	subs map[string]map[string]goqa.Subscriber
	id   int
	mut  sync.Mutex
}

func New() *Memory {
	return &Memory{
		subs: make(map[string]map[string]goqa.Subscriber),
	}
}

func (s *Memory) Subscribe(ctx context.Context, name string, sub goqa.Subscriber) error {
	if name == "" || sub == nil {
		return nil
	}

	defer s.mut.Unlock()
	s.mut.Lock()

	s.id++

	var id = name + "-" + strconv.Itoa(s.id) // poor man's uuid

	if s.subs == nil {
		s.subs = make(map[string]map[string]goqa.Subscriber)
	}

	if _, ok := s.subs[name]; !ok {
		s.subs[name] = map[string]goqa.Subscriber{id: sub}
	} else {
		s.subs[name][id] = sub
	}

	sub.SetID(id)

	return nil
}

func (s *Memory) Unsubscribe(ctx context.Context, id string) error {
	defer s.mut.Unlock()
	s.mut.Lock()

	var i = strings.LastIndex(id, "-")
	if i == -1 {
		return nil
	}

	var name = id[:i]
	if _, ok := s.subs[name]; !ok {
		return nil
	}

	if _, ok := s.subs[name][id]; ok {
		delete(s.subs[name], id)
	}

	return nil
}

func (s *Memory) Subscribers(ctx context.Context, name string) ([]goqa.Subscriber, error) {
	if s.subs == nil {
		return nil, nil
	}

	defer s.mut.Unlock()
	s.mut.Lock()

	if _, ok := s.subs[name]; !ok {
		return nil, nil
	}

	var (
		subs = make([]goqa.Subscriber, len(s.subs[name]))
		i    int
	)

	for k := range s.subs[name] {
		subs[i] = s.subs[name][k]
		i++
	}

	return subs, nil
}

func (s *Memory) Close() error {
	return nil
}
