package graceful

import (
	"io"
	"log"
	"os"
	"strings"
)

var logger = log.New(io.Discard, "[go-shutdown-graceful] ", log.LstdFlags|log.LUTC)

func init() {
	lvl, _ := os.LookupEnv("GO_SHUTDOWN_GRACEFUL_LOG")
	if strings.ToLower(lvl) == "true" {
		EnableLogging()
	}
}

// EnableLogging for this module.
func EnableLogging() {
	logger.SetOutput(os.Stderr)
}

// DisableLogging for this module.
func DisableLogging() {
	logger.SetOutput(io.Discard)
}
