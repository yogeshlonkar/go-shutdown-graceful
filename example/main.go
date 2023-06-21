package main

import (
	"context"
	"github.com/yogeshlonkar/go-shutdown-graceful"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "[example] ", log.LstdFlags|log.LUTC)
	logger.Println("started")
	go someGoroutine()
	// if INT or TERM signal is received, go-shutdown-graceful will trigger shutdown signal to all observers.
	// Observers can do cleanup and call done() to notify go-shutdown-graceful that they are done.
	// Default timeout for cleanup is 30 seconds. This can be changed by calling HandleSignals with a time.Duration value.
	//graceful.HandleSignals(0)
	graceful.HandleSignalsWithContext(context.Background(), 0)
	logger.Println("graceful shutdown complete")
}

func someGoroutine() {
	stop := make(chan struct{})
	go actualGoRoutine(stop)
	shutdown, done := graceful.NewShutdownObserver()
	<-shutdown
	stop <- struct{}{}
	done()
}

func actualGoRoutine(stop chan struct{}) {
	// doing something
	<-stop
}
