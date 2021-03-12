package memory

import (
	"context"
	"sync"

	"github.com/fluxynet/goqa"
)

type Memory struct {
	listeners []chan goqa.Event
	mutex     sync.Mutex
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Listen(ctx context.Context) (<-chan goqa.Event, error) {
	var c = make(chan goqa.Event)

	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.listeners = append(m.listeners, c)

	return c, nil
}

func (m *Memory) Publish(ctx context.Context, event goqa.Event) error {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	var wg sync.WaitGroup
	wg.Add(len(m.listeners))

	for i := range m.listeners {
		go func(l chan goqa.Event) {
			defer wg.Done()
			l <- event
		}(m.listeners[i])
	}

	wg.Wait()

	return nil
}

func (m *Memory) Close() error {
	if m.listeners == nil {
		return nil
	}

	defer m.mutex.Unlock()
	m.mutex.Lock()

	var wg sync.WaitGroup
	wg.Add(len(m.listeners))

	for i := range m.listeners {
		go func(l chan goqa.Event) {
			defer wg.Done()
			close(l)
		}(m.listeners[i])
	}

	wg.Wait()

	m.listeners = nil

	return nil
}
