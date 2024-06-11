package ticker

import (
	"context"
	"time"
)

type TimeTicker struct {
	ticker   *time.Ticker
	done     chan struct{}
	runnning bool
}

func NewTimeTicker() *TimeTicker {
	return &TimeTicker{
		done: make(chan struct{}),
	}
}

func (t *TimeTicker) Start(ctx context.Context, period time.Duration, callback func()) {
	t.runnning = true
	t.done = make(chan struct{})

	// initial run
	callback()

	t.ticker = time.NewTicker(period)

	for loop := true; loop; {
		select {
		case <-t.ticker.C:
			callback()
		case <-t.done:
			t.runnning = false
			loop = false
			break
		case <-ctx.Done():
			t.runnning = false
			loop = false
			break
		}
	}

	t.ticker.Stop()
	close(t.done)
}

func (t *TimeTicker) Stop() {
	t.done <- struct{}{}
}

func (t *TimeTicker) IsRunning() bool {
	return t.runnning
}

func (t *TimeTicker) IsStopped() bool {
	return !t.runnning
}
