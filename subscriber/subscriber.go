package subscriber

import (
	"errors"
)

var (
	// ErrSerializeNotSupported for subscribers that are ephemeral
	ErrSerializeNotSupported = errors.New("serialize is not supported")

	// ErrUnsupportedEvent when a subscriber gets an event it cannot handle
	ErrUnsupportedEvent = errors.New("event is not supported by this subscriber")
)

// Identifiable composition
type Identifiable struct {
	id string
}

func (s Identifiable) ID() string {
	return s.id
}

func (s *Identifiable) SetID(id string) {
	s.id = id
}
