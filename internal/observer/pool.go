package observer

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// DefaultTimeout for waiting of observer/ goroutines to be done.
const DefaultTimeout = 30 * time.Second

var ErrTimeout = errors.New("timeout waiting")

// Pool to manage goroutines monitoring the graceful shutdown.
type Pool struct {
	count  *atomic.Int64
	logger *log.Logger
	done   chan struct{}
	wg     *sync.WaitGroup
}

func NewPool(logger *log.Logger) *Pool {
	return &Pool{
		wg:     &sync.WaitGroup{},
		count:  &atomic.Int64{},
		logger: logger,
		done:   make(chan struct{}),
	}
}

// Add will add a shutdown observer (goroutine) to pool.
func (o *Pool) Add() func() {
	o.wg.Add(1)
	o.count.Add(1)
	return o.newCloser()
}

// Pending returns the number of pending observers.
func (o *Pool) Pending() int {
	return int(o.count.Load())
}

// Close will wait for all observers to finish. returns a channel closed channel.
func (o *Pool) Close(timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	select {
	case <-time.After(timeout):
		if o.Pending() > 0 {
			return fmt.Errorf("graceful shutodwn: %d observers not closed: %w", o.Pending(), ErrTimeout)
		}
	case <-o.wait():
	}
	o.logger.Println("all observers closed")
	return nil
}

// newCloser will return a function to be called when routine is done.
// The function should be called only once.
func (o *Pool) newCloser() func() {
	closed := &atomic.Bool{}
	return func() {
		if closed.Load() {
			o.logger.Println("ignoring close call, observer already closed")
			return
		}
		closed.Store(true)
		o.wg.Done()
		o.count.Add(-1)
	}
}

func (o *Pool) wait() <-chan struct{} {
	go func() {
		o.wg.Wait()
		close(o.done)
	}()
	return o.done
}
