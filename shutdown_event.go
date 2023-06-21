package graceful

import (
	"context"
	"os"
	"os/signal"
)

type shutdownEvent struct {
	Fired      os.Signal // <- add this
	ContextErr error
}

// observe copy of signal.NotifyContext with added Fired field to return the signal that caused the context to be canceled.
func observe(parent context.Context, signals ...os.Signal) *shutdownEvent {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	event := &shutdownEvent{}
	if ctx.Err() == nil {
		select {
		case fired := <-ch:
			event.Fired = fired // <- add this
		case <-ctx.Done():
			event.ContextErr = ctx.Err()
		}
	}
	return event
}
