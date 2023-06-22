package graceful

import (
	"context"
	"errors"
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestHandleHandleSignalsWithContext_shutdown(t *testing.T) {
	initGrace()
	tested := &atomic.Bool{}
	_, done := NewShutdownObserver()
	go func() {
		err := HandleSignalsWithContext(context.Background(), 0)
		defer tested.Store(true)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if err := Shutdown(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	done()
	time.Sleep(1 * time.Second)
	if !tested.Load() {
		t.Error("expected to complete HandleSignalsWithContext")
	}
}

func TestHandleHandleSignalsWithContext_context_cancel(t *testing.T) {
	initGrace()
	tested := &atomic.Bool{}
	_, done := NewShutdownObserver()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := HandleSignalsWithContext(ctx, 0)
		defer tested.Store(true)
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
	if !tested.Load() {
		t.Error("expected to complete HandleSignalsWithContext")
	}
}

func TestHandleHandleSignalsWithContext_timout(t *testing.T) {
	initGrace()
	tested := &atomic.Bool{}
	NewShutdownObserver()
	go func() {
		err := HandleSignalsWithContext(context.Background(), 50*time.Millisecond)
		defer tested.Store(true)
		if err == nil {
			t.Error("expected err, got nil")
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
		}
	}()
	time.Sleep(20 * time.Millisecond)
	if err := Shutdown(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	time.Sleep(1 * time.Second)
	if !tested.Load() {
		t.Error("expected to complete HandleSignalsWithContext")
	}
}

func TestHandleHandleSignals_shutdown(t *testing.T) {
	initGrace()
	tested := &atomic.Bool{}
	_, done := NewShutdownObserver()
	go func() {
		err := HandleSignals(0)
		defer tested.Store(true)
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
	time.Sleep(1 * time.Second)
	if !tested.Load() {
		t.Error("expected to complete HandleSignals")
	}
}

func TestHandleHandleSignals_shutdown_timeout(t *testing.T) {
	initGrace()
	tested := &atomic.Bool{}
	NewShutdownObserver()
	go func() {
		err := HandleSignals(50 * time.Millisecond)
		defer tested.Store(true)
		if err == nil {
			t.Error("expected err, got nil")
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
		}
	}()
	time.Sleep(20 * time.Millisecond)
	if err := Shutdown(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	time.Sleep(1 * time.Second)
	if !tested.Load() {
		t.Error("expected to complete HandleSignals")
	}
}

func TestHandleHandleSignals_signal_args(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}
	initGrace()
	tested := &atomic.Bool{}
	NewShutdownObserver()
	go func() {
		err := HandleSignals(50*time.Millisecond, syscall.SIGABRT)
		defer tested.Store(true)
		if err == nil {
			t.Error("expected err, got nil")
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
		}
	}()
	time.Sleep(20 * time.Millisecond)

	if err := sendSignal(syscall.Getpid(), syscall.SIGABRT); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if !tested.Load() {
		t.Error("expected to complete HandleSignals")
	}
}

func TestShutdown_failure(t *testing.T) {
	observedSignals = []os.Signal{}
	err := Shutdown()
	if err == nil {
		t.Error("expected err, got nil")
		return
	}
	if !errors.Is(err, ErrSigTermNotObserved) {
		t.Errorf("expected '%v' to be in error tree, got '%v'", ErrSigTermNotObserved, err)
	}
}
