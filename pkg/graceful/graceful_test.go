package graceful

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func SendSigtermAfter(t time.Duration) {
	time.AfterFunc(t, func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	})
}

var wg = sync.WaitGroup{}

func StartApplication() {
	wg.Add(1)
}
func StopApplication() {
	appStopped = true
	wg.Done()
}
func Listen() {
	wg.Wait()
}

var appStopped = false

func TestGraceful(t *testing.T) {
	StartApplication()

	OnShutdown(func() {
		t.Log("Shutdown")
		StopApplication()
	})

	SendSigtermAfter(3 * time.Second)

	Listen()
	Shutdown()
	assert.True(t, appStopped)
}

func TestGracefulTerminate(t *testing.T) {
	appStopped := false
	OnShutdown(func() {
		appStopped = true
		t.Log("Shutdown")
	})

	time.AfterFunc(1*time.Second, func() {
		Terminate() // app won't shutdown gracefully
	})

	Shutdown()
	assert.Equal(t, appStopped, false)
}
