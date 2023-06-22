package graceful

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/yogeshlonkar/go-shutdown-graceful/internal/observer"
)

const delay = 10 * time.Millisecond

func TestShutdownWithContext_trigger(t *testing.T) {
	go func() {
		time.Sleep(delay)
		_, done := NewObserver()
		TriggerShutdown()
		done()
	}()
	if err := ShutdownWithContext(context.Background(), 0); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestShutdownWithContext_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(delay)
		_, done := NewObserver()
		cancel()
		done()
	}()
	if err := ShutdownWithContext(ctx, 0); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestShutdownWithContext_timout(t *testing.T) {
	go func() {
		time.Sleep(delay)
		NewObserver()
		TriggerShutdown()
	}()
	err := ShutdownWithContext(context.Background(), 50*time.Millisecond)
	if err == nil {
		t.Error("expected err, got nil")
	}
	if !errors.Is(err, observer.ErrTimeout) {
		t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
	}
}

func TestShutdown_trigger(t *testing.T) {
	go func() {
		time.Sleep(delay)
		_, done := NewObserver()
		TriggerShutdown()
		done()
	}()
	if err := Shutdown(0); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestShutdown_timeout(t *testing.T) {
	go func() {
		time.Sleep(delay)
		NewObserver()
		TriggerShutdown()
	}()
	err := Shutdown(delay)
	if err == nil {
		t.Error("expected err, got nil")
	}
	if !errors.Is(err, observer.ErrTimeout) {
		t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
	}
}

func TestShutdown_args(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}
	go func() {
		time.Sleep(delay)
		NewObserver()
		if err := trigger(syscall.SIGABRT); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	}()
	err := Shutdown(50*time.Millisecond, syscall.SIGABRT, syscall.SIGINT)
	if err == nil {
		t.Error("expected err, got nil")
	}
	if !errors.Is(err, observer.ErrTimeout) {
		t.Errorf("expected '%v' to be in error tree, got '%v'", context.Canceled, err)
	}
}

// trigger will send a syscall.SIGTERM signal to given process id.
func trigger(sig syscall.Signal) error {
	pid := os.Getpid()
	if err := syscall.Kill(pid, sig); err != nil {
		return fmt.Errorf("failed to send signal %s(%#v) to current process %d: %w", syscall.SIGTERM.String(), syscall.SIGTERM, pid, err)
	}
	return nil
}
