package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// DefaultSignals for graceful shutdown.
// kill (no param) default sends syscall.SIGTERM
// kill -2 is syscall.SIGINT
// kill -15 is syscall.SIGTERM.
var DefaultSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

type Hook struct {
	shutdown chan struct{}
	signals  []os.Signal
	trigger  chan struct{}
}

func NewHook() *Hook {
	return &Hook{
		shutdown: make(chan struct{}),
		signals:  DefaultSignals,
		trigger:  make(chan struct{}),
	}
}

func (h *Hook) Notifier() <-chan struct{} {
	return h.shutdown
}

// Observe is copy of signal.NotifyContext modified to capture signal.
// It returns the signal that caused shutdown or error in case of timeout/ context cancellation.
func (h *Hook) Observe(parent context.Context, signals ...os.Signal) (*string, error) {
	if len(signals) > 0 {
		h.signals = signals
	}
	defer close(h.shutdown)
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	if err := ctx.Err(); err != nil {
		return nil, err //nolint:wrapcheck // handled in caller
	}
	select {
	case sig := <-ch:
		sigStr := fmt.Sprintf("%s(%#v)!", sig.String(), sig)
		return &sigStr, nil
	case <-ctx.Done():
		return nil, ctx.Err() //nolint:wrapcheck // handled in caller
	case <-h.trigger:
		return nil, nil
	}
}

// Trigger will send a syscall.SIGTERM signal to given process id.
func (h *Hook) Trigger() {
	close(h.trigger)
}
