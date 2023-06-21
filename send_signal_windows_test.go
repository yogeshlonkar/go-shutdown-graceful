//go:build windows

package graceful

import (
	"syscall"
	"testing"
)

func Test_sendSignal_failure(t *testing.T) {
	t.Run("should return err if signal is not sent", func(t *testing.T) {
		if err := sendSignal(-123, syscall.SIGHUP); err == nil {
			t.Error("expected err, got nil")
		}
	})
}
