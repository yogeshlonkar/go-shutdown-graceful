package graceful

import (
	"io"
	"os"
	"testing"
)

func Test_configerLogger(t *testing.T) {
	t.Run("should enable logger if GO_SHUTDOWN_GRACEFUL_LOG is set to true", func(t *testing.T) {
		t.Setenv("GO_SHUTDOWN_GRACEFUL_LOG", "true")
		configerLogger()
		if logger.Writer() != os.Stderr {
			t.Error("expected logger to be configured")
		}
	})
	t.Run("should disable logger if GO_SHUTDOWN_GRACEFUL_LOG is not set false", func(t *testing.T) {
		t.Setenv("GO_SHUTDOWN_GRACEFUL_LOG", "false")
		configerLogger()
		if logger.Writer() != io.Discard {
			t.Error("expected logger to discard logs")
		}
	})
}

func TestEnableLogging(t *testing.T) {
	EnableLogging()
	if logger.Writer() != os.Stderr {
		t.Error("expected logger to be configured")
	}
}

func TestDisableLogging(t *testing.T) {
	DisableLogging()
	if logger.Writer() != io.Discard {
		t.Error("expected logger to be configured")
	}
}
