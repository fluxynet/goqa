package coverage

import (
	"context"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/subscriber"
)

type Coverage struct {
	subscriber.Identifiable
	broker goqa.Broker
}

func New(broker goqa.Broker) *Coverage {
	return &Coverage{
		broker: broker,
	}
}

func (c Coverage) Notify(event goqa.Event) error {
	var e *goqa.GithubEvent

	switch v := event.(type) {
	default:
		return subscriber.ErrUnsupportedEvent
	case *goqa.GithubEvent:
		e = v
	case goqa.GithubEvent:
		e = &v
	}

	for i := range e.Coverage {
		var ev = goqa.CoverageEvent{
			Pkg:        e.Coverage[i].Pkg,
			Percentage: e.Coverage[i].Percentage,
			Time:       e.Coverage[i].Time,
		}

		err := c.broker.Publish(context.Background(), ev)
		if err != nil {
			return err
		}
	}

	return nil
}
