package graceful

import (
	"io"
	"os"
	"testing"
)

func Test_configerLogger(t *testing.T) {
	err := os.Setenv("GO_SHUTDOWN_GRACEFUL_LOG", "true")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	configerLogger()
	if logger.Writer() != os.Stderr {
		t.Error("expected logger to be configured")
	}
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
