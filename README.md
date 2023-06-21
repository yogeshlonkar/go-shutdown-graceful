# go-shutdown-graceful [![Go Reference](https://pkg.go.dev/badge/github.com/yogeshlonkar/go-shutdown-graceful.svg)](https://pkg.go.dev/github.com/yogeshlonkar/go-shutdown-graceful) [![Continuous Integration](https://github.com/yogeshlonkar/go-shutdown-graceful/actions/workflows/on-push.yaml/badge.svg)](https://github.com/yogeshlonkar/go-shutdown-graceful/actions/workflows/on-push.yaml) [![Go Report Card](https://goreportcard.com/badge/github.com/yogeshlonkar/go-shutdown-graceful)](https://goreportcard.com/report/github.com/yogeshlonkar/go-shutdown-graceful)


Handle application graceful shutdown.

## ✏️ [Example]

## 🧑‍💻 Usage

### Functions

- `NewShutdownObserver` to create a new observer goroutine which will be notified when os signal is received via shutdown (1st return value). This routine should clean up itself, and other spawn routines. It should call `done()` (2nd return value) to notify `go-shutdown-graceful` that cleanup is done.
- `HandleSignals` to hold the main goroutine until os signal is received to the process
- `HandleSignalsWithContext` to hold the main goroutine until os signal is received to the process or passed context is canceled
- `Shutdown` to trigger shutdown signal to all observers. This is useful when you want to shut down based on API hook or goroutine other than the main goroutine.

### Logging

If `GO_SHUTDOWN_GRACEFUL_LOG` environment variable is set to `true`, then logging is enabled.
This can be overridden by calling the following functions:

- `EnableLogging` to enable logging in this module. By default, logging is disabled.
- `DisableLogging` to disable logging in this module.


```go
package main

import (
    "github.com/yogeshlonkar/go-shutdown-graceful"
)

func main() {
    go someGoroutine()
    // if INT or TERM signal is received, go-shutdown-graceful will trigger shutdown signal to all observers.
    // Observers can do cleanup and call done() to notify go-shutdown-graceful that they are done.
    // Default timeout for cleanup is 30 seconds. This can be changed by calling HandleOsSignals with a time.Duration value.
    graceful.HandleSignals(0)
}

func someGoroutine() {
    // do something in separate goroutine
    shutdown, done := graceful.NewShutdownObserver()
    <-shutdown
    // close the background goroutine started before
    done()
}
```

[Example]: ./example/README.md

