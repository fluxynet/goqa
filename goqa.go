package goqa

import (
	"context"
	"strconv"
)

const (
	// EventGithub means we got something from github
	EventGithub = "EVENT_GITHUB"

	// EventCoverage denotes new coverage data being available
	EventCoverage = "EVENT_COVERAGE"
)

// Coverage represents actual coverage of a package
type Coverage struct {
	// Pkg name
	Pkg string `json:"pkg"`

	// Percentage expressed to the nearest integer (on a scale of 0-100)
	Percentage int `json:"percentage"`

	// Time the measurement was done
	Time string `json:"time"`
}

// String representation of coverage information
func (c Coverage) String() string {
	return `[` + c.Time + `] pkg = "` + c.Pkg + `" %` + strconv.Itoa(c.Percentage)
}

// Repo allows persistance of data to permanent storage
type Repo interface {
	// Save coverage data on permanent storage
	Save(ctx context.Context, covs ...Coverage) error

	// Load latest coverage data available
	Load(ctx context.Context) ([]Coverage, error)

	// Close the repo
	Close() error
}

// Cache mechanism to store coverage data in memory
type Cache interface {
	// Reinit purges existing coverage info and stores new info
	Reset(covs ...Coverage) error

	// Get coverage data for a package
	Get(pkg string) (*Coverage, bool)

	// Keys get all keys therein stored
	Keys() ([]string, error)

	// Close the cache
	Close() error
}

// Event is emitted and listened to
type Event interface {
	// Name of the event
	Name() string

	// String representation of the event
	// not the best, a simple string cannot satisfy all use cases
	// maybe should not use string but map[string]interface{} for key => values map
	String() string
}

// Subscriber is someone who listens to events
type Subscriber interface {
	// ID
	ID() string

	// SetID
	SetID(id string)

	// Notify sends the event to the subscriber
	Notify(event Event) error
}

// Roster keeps track of subscribers
type Roster interface {
	// Subscribe to an event source for an event of a given name and get a subscription id
	Subscribe(ctx context.Context, name string, sub Subscriber) error

	// Unsubscribe from an event source by the subscription id
	Unsubscribe(ctx context.Context, id string) error

	// Roster returns active subscribers for an event of a given name
	Subscribers(ctx context.Context, name string) ([]Subscriber, error)

	// Close the source
	Close() error
}

// Broker is an event broker like SQS or RabbitMQ
type Broker interface {
	// Listen to events from a broker
	Listen(ctx context.Context) (<-chan Event, error)

	// Publish event to a broker
	Publish(ctx context.Context, event Event) error

	// Close the broker
	Close() error
}
