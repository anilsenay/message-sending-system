package graceful

import (
	"os"
	"os/signal"
	"syscall"
)

var serverShutdown = make(chan struct{})

// Wait termination signal and execute user functions when signal received
func OnShutdown(callback func()) {
	// create a channel to fill with a termination signal when received
	sig := make(chan os.Signal, 1)
	// when an interrupt signal received, send it to "sig" channel
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		// wait Notify function to release sig channel
		<-sig

		// execute user functions
		callback()

		// fill serverShutdown channel to block termination
		serverShutdown <- struct{}{}
	}()
}

// Shutdown application
func Shutdown() {
	// release serverShutdown channel and allow application to terminate
	<-serverShutdown
}

func Terminate() {
	serverShutdown <- struct{}{}
}
