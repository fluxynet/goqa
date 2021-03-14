package goqa

import (
	"context"
	"log"
)

// Attach a roster to a broker; stops if the context stops
func Attach(ctx context.Context, b Broker, r Roster) {
	var c, err = b.Listen(context.Background())
	if err != nil {
		return
	}

	defer Closed(b)

	for {
		select {
		case event := <-c:
			var subs, err = r.Subscribers(context.Background(), event.Name())
			if err == nil {
				go Publish(event, subs...)
			} else {
				log.Printf("failed to get subscribers\n%s\n%s\n", err.Error(), event.Name())
			}
		case <-ctx.Done():
			return
		}
	}
}

// Publish event to many subscribers
func Publish(event Event, subs ...Subscriber) {
	for i := range subs {
		var err = subs[i].Notify(event)
		if err != nil {
			log.Printf("failed to notify subscriber\n%s\n%s\n", err.Error(), event.Name())
		}
	}
}
