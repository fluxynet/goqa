package memory

import (
	"context"
	"sync"

	"github.com/fluxynet/goqa"
)

func init() {
	var _ goqa.Broker = New()
}

type Memory struct {
	listeners []chan goqa.Event
	mutex     sync.RWMutex
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Listen(ctx context.Context) (<-chan goqa.Event, error) {
	var c = make(chan goqa.Event)

	defer m.mutex.Unlock()

	m.listeners = append(m.listeners, c)

	m.mutex.Lock()

	return c, nil
}

func (m *Memory) Publish(ctx context.Context, event goqa.Event) error {
	defer m.mutex.RUnlock()
	m.mutex.RLock()

	var wg sync.WaitGroup
	wg.Add(len(m.listeners))

	for i := range m.listeners {
		go func(l chan goqa.Event) {
			defer wg.Done()
			l <- event
		}(m.listeners[i])
	}

	return nil
}

func (m *Memory) Close() error {
	defer m.mutex.RUnlock()
	m.mutex.RLock()

	var wg sync.WaitGroup
	wg.Add(len(m.listeners))

	for i := range m.listeners {
		go func(l chan goqa.Event) {
			defer wg.Done()
			close(l)
		}(m.listeners[i])
	}

	m.listeners = nil

	return nil
}
