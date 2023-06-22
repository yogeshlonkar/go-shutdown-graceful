package graceful

import (
	"context"
	"os"
	"time"

	"github.com/yogeshlonkar/go-shutdown-graceful/internal/observer"
	"github.com/yogeshlonkar/go-shutdown-graceful/internal/shutdown"
)

var (
	hook      *shutdown.Hook
	observers *observer.Pool
)

// NewObserver will add the observer to a shutdown Pool (goroutines).
// It returns a channel for listening of shutdown signal and close function to be called when routine is done.
func NewObserver() (<-chan struct{}, func()) {
	closer := observers.Add()
	return hook.Notifier(), closer
}

// Shutdown will wait for given signals once received, any of these signals will
// send shutdown signal to goroutines listening on TriggerShutdown.
// It waits for all goroutines to finish within timeout duration before exiting.
// It should be called in the main goroutine to hold the process.
//
// If timeout is 0, DefaultTimeout is used as default timeout.
// If no signals are given, syscall.SIGINT, syscall.SIGTERM are used.
//
// syscall.SIGKILL but can't be caught, so it can't be handled.
func Shutdown(timeout time.Duration, signals ...os.Signal) error {
	return ShutdownWithContext(context.Background(), timeout, signals...)
}

// ShutdownWithContext is the same as Shutdown but with context support.
func ShutdownWithContext(ctx context.Context, timeout time.Duration, signals ...os.Signal) error {
	configerLogger()
	hook = shutdown.NewHook()
	observers = observer.NewPool(logger)
	// Wait for interrupt signal to gracefully shutdown the server with
	sig, err := hook.Observe(ctx, signals...)
	if sig != nil {
		logger.Printf("shutting down: received %s", *sig)
	}
	if err != nil {
		logger.Printf("shutting down: %q", err)
	}
	if sig == nil && err == nil {
		logger.Println("shutting down: TriggerShutdown requested")
	}
	logger.Printf("waiting %d for services/ routines to finish", observers.Pending())
	if err = observers.Close(timeout); err != nil {
		return err
	}
	return err
}

// TriggerShutdown does not trigger any os signals.
// It stop waiting for os signals and send shutdown signal to observers.
func TriggerShutdown() {
	hook.Trigger()
}
