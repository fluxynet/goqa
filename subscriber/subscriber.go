package subscriber

import (
	"errors"
)

var (
	// ErrUnsupportedEvent when a subscriber gets an event it cannot handle
	ErrUnsupportedEvent = errors.New("event is not supported by this subscriber")
)

// Identifiable intended for composition to satisfy ID segment of subscriber interface
type Identifiable struct {
	id string
}

func (s Identifiable) ID() string {
	return s.id
}

func (s *Identifiable) SetID(id string) {
	s.id = id
}
