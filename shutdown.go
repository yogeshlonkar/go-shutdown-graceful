package graceful

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// DefaultTimeout is the default timeout for shutdown.
const DefaultTimeout = 30 * time.Second

var (
	shutdown     chan struct{}
	routinesDone chan struct{}
	observers    *observerPool
	ErrTimeout   = errors.New("timeout waiting")
)

// kill (no param) default sends syscall.SIGTERM
// kill -2 is syscall.SIGINT
// kill -15 is syscall.SIGTERM.
var defaultOsSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

func init() {
	initGrace()
}

func initGrace() {
	shutdown = make(chan struct{})
	routinesDone = make(chan struct{})
	observers = &observerPool{wg: &sync.WaitGroup{}, count: &atomic.Int64{}}
}

// NewShutdownObserver will add a shutdown observerPool (goroutine) to the wait list.
// It returns a channel for listening of shutdown signal and close function to be called when routine is done.
func NewShutdownObserver() (<-chan struct{}, func()) {
	closer := observers.Add()
	return shutdown, closer
}

// HandleSignals will wait for given signals once received, any of these signals will
// send shutdown signal to goroutines listening on Shutdown.
// It waits for all goroutines to finish within timeout duration before exiting.
// It should be called in the main goroutine to hold the process.
//
// If timeout is 0, DefaultTimeout is used as default timeout.
// If no signals are given, syscall.SIGINT, syscall.SIGTERM are used.
//
// syscall.SIGKILL but can't be caught, so it can't be handled.
func HandleSignals(timeout time.Duration, signals ...os.Signal) error {
	return HandleSignalsWithContext(context.Background(), timeout, signals...)
}

// HandleSignalsWithContext is the same as HandleSignals but with context support.
func HandleSignalsWithContext(ctx context.Context, timeout time.Duration, signals ...os.Signal) error {
	configerLogger()
	if len(signals) == 0 {
		signals = defaultOsSignals
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	// Wait for interrupt signal to gracefully shutdown the server with
	event := observe(ctx, signals...)
	if event.Fired != nil {
		logger.Printf("received %s(%#v)! shutting down", event.Fired.String(), event.Fired)
	}
	go func() {
		close(shutdown)
		observers.Wait()
		close(routinesDone)
	}()
	logger.Printf("waiting %d for services/ routines to finish", observers.Pending())
	select {
	case <-time.After(timeout):
		if observers.Pending() > 0 {
			return fmt.Errorf("graceful shutodwn: %d observers not closed: %w", observers.Pending(), ErrTimeout)
		}
	case <-routinesDone:
	}
	logger.Println("all observers closed")
	if event.Fired == nil {
		return fmt.Errorf("graceful shutodwn: %w", event.ContextErr)
	}
	return nil
}

// Shutdown will send shutdown signal to goroutines listening on Shutdown.
// Goroutine suspended by calling HandleSignals or HandleSignalsWithContext will resume.
func Shutdown() error {
	return sendSignal(syscall.Getpid(), syscall.SIGINT)
}
