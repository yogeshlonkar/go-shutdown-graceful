package graceful

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHandleHandleSignalsWithContext(t *testing.T) {
	t.Run("should return nil if shutdown on signal", func(t *testing.T) {
		initGrace()
		tested := false
		_, done := NewShutdownObserver()
		go func() {
			err := HandleSignalsWithContext(context.Background(), 0)
			tested = true
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		err := Shutdown()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		done()
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if !tested {
			t.Error("expected to complete HandleSignalsWithContext")
		}
	})
	t.Run("should return err if shutdown on context", func(t *testing.T) {
		initGrace()
		tested := false
		_, done := NewShutdownObserver()
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			err := HandleSignalsWithContext(ctx, 0)
			tested = true
			if err == nil {
				t.Error("expected err, got nil")
			}
			if !errors.Is(err, context.Canceled) {
				t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		cancel()
		time.Sleep(100 * time.Millisecond)
		done()
		time.Sleep(100 * time.Millisecond)
		if !tested {
			t.Error("expected to complete HandleSignalsWithContext")
		}
	})
	t.Run("should return err if shutdown on timeout", func(t *testing.T) {
		initGrace()
		tested := false
		NewShutdownObserver()
		go func() {
			err := HandleSignalsWithContext(context.Background(), 50*time.Millisecond)
			tested = true
			if err == nil {
				t.Error("expected err, got nil")
			}
			if !errors.Is(err, ErrTimeout) {
				t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
			}
		}()
		time.Sleep(20 * time.Millisecond)
		err := Shutdown()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		if !tested {
			t.Error("expected to complete HandleSignalsWithContext")
		}
	})
}
