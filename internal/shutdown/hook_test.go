package shutdown

import (
	"context"
	"errors"
	"testing"
)

func TestHook_Observe(t *testing.T) {
	t.Run("don't observe if context is cancelled", func(t *testing.T) {
		h := NewHook()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := h.Observe(ctx)
		if err == nil {
			t.Error("expected err, got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected '%v', got '%v'", context.Canceled, err)
		}
	})
}
