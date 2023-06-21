package graceful

import (
	"sync"
	"sync/atomic"
)

// observerPool to manage goroutines monitoring the graceful shutdown.
type observerPool struct {
	wg    *sync.WaitGroup
	count *atomic.Int64
}

// Add will add a shutdown observer (goroutine) to pool.
func (o *observerPool) Add() func() {
	o.wg.Add(1)
	o.count.Add(1)
	return o.newCloser()
}

// newCloser will return a function to be called when routine is done.
// The function should be called only once.
func (o *observerPool) newCloser() func() {
	closed := &atomic.Bool{}
	return func() {
		if closed.Load() {
			logger.Println("ignoring close call, observer already closed")
			return
		}
		closed.Store(true)
		o.wg.Done()
		o.count.Add(-1)
	}
}

// Pending returns the number of pending observers.
func (o *observerPool) Pending() int {
	return int(o.count.Load())
}

// Wait will wait for all observers to finish.
func (o *observerPool) Wait() {
	o.wg.Wait()
}
