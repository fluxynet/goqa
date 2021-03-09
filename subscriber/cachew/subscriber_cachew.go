package cachew

import (
	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

func init() {
	var _ goqa.Subscriber = New(nil)
}

func New(c goqa.Cache) *Cache {
	return &Cache{cache: c}
}

// Cache is a subscriber that listens to goqa.GithubEvent, updates a goqa.Cache and then emits a goqa.CoverageEvent
type Cache struct {
	subscriber.Identifiable
	cache goqa.Cache
}

func (c *Cache) Notify(event goqa.Event) error {
	var e, ok = event.(goqa.GithubEvent)
	if !ok {
		return subscriber.ErrUnsupportedEvent
	}

	var err = c.cache.Reinit(e.Coverage...)
	return err
}

func (c *Cache) Serialize() (string, error) {
	return "", subscriber.ErrSerializeNotSupported
}

func (c *Cache) Unserialize(string) error {
	return subscriber.ErrSerializeNotSupported
}
