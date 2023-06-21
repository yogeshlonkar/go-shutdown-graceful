package graceful

import (
	"sync"
	"sync/atomic"
	"testing"
)

func Test_observerPool(t *testing.T) {
	observers := &observerPool{wg: &sync.WaitGroup{}, count: &atomic.Int64{}}
	t.Run("Add", func(t *testing.T) {
		observers.Add()
		if observers.Pending() != 1 {
			t.Errorf("Pending() = %v, want %v", observers.Pending(), 1)
		}
	})
	t.Run("newCloser", func(t *testing.T) {
		closer := observers.newCloser()
		closer()
		if observers.Pending() != 0 {
			t.Errorf("Pending() = %v, want %v", observers.Pending(), 0)
		}
		closer()
		if observers.Pending() != 0 {
			t.Errorf("Pending() = %v, want %v", observers.Pending(), 0)
		}
	})
}
