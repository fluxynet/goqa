package roster

import (
	"context"

	"github.com/fluxynet/goqa"
)

// WatchCtx a context and unsubscribe when done
func WatchCtx(ctx context.Context, roster goqa.Roster, sub goqa.Subscriber) {
	defer roster.Unsubscribe(context.Background(), sub.ID())
	<-ctx.Done()
}
