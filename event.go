package graceful

import (
	"context"
	"os"
	"os/signal"
)

type event struct {
	Fired      os.Signal // <- add this
	ContextErr error
}

// observe copy of signal.NotifyContext with added Fired field to return the signal that caused the context to be canceled.
func observe(parent context.Context, signals ...os.Signal) *event {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	_event := &event{}
	if ctx.Err() == nil {
		select {
		case fired := <-ch:
			_event.Fired = fired // <- add this
		case <-ctx.Done():
			_event.ContextErr = ctx.Err()
		}
	}
	return _event
}
