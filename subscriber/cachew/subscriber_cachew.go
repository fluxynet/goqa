package cachew

import (
	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

func New(c goqa.Cache) *Cache {
	return &Cache{cache: c}
}

// Cache is a subscriber that listens to goqa.GithubEvent, updates a goqa.Cache and then emits a goqa.CoverageEvent
type Cache struct {
	subscriber.Identifiable
	cache goqa.Cache
}

func (c *Cache) Notify(event goqa.Event) error {
	var e *goqa.GithubEvent

	switch v := event.(type) {
	default:
		return subscriber.ErrUnsupportedEvent
	case *goqa.GithubEvent:
		e = v
	case goqa.GithubEvent:
		e = &v
	}

	var err = c.cache.Reset(e.Coverage...)
	return err
}
