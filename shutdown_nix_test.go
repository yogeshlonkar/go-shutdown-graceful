//go:build !windows

package graceful

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/yogeshlonkar/go-shutdown-graceful/internal/observer"
)

func TestShutdown_args(t *testing.T) {
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
