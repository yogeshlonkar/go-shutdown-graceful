package observer

import (
	"io"
	"log"
	"testing"
)

func TestPool(t *testing.T) {
	pool := NewPool(log.New(io.Discard, "", 0))
	t.Run("Add", func(t *testing.T) {
		pool.Add()
		if pool.Pending() != 1 {
			t.Errorf("Pending() = %v, want %v", pool.Pending(), 1)
		}
	})
	t.Run("newCloser", func(t *testing.T) {
		closer := pool.newCloser()
		closer()
		if pool.Pending() != 0 {
			t.Errorf("Pending() = %v, want %v", pool.Pending(), 0)
		}
		closer()
		if pool.Pending() != 0 {
			t.Errorf("Pending() = %v, want %v", pool.Pending(), 0)
		}
	})
}
