package ticker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeTicker(t *testing.T) {
	t.Parallel()

	t.Run("TimeTicker Created", func(t *testing.T) {
		ticker := NewTimeTicker()
		assert.NotNil(t, ticker)
	})

	t.Run("TimeTicker Start", func(t *testing.T) {
		tickCount := 0

		timeTicker := NewTimeTicker()
		go timeTicker.Start(context.Background(), time.Second*1, func() {
			t.Log("tick")
			tickCount = tickCount + 1
		})

		time.Sleep(time.Millisecond * 3500)
		assert.Equal(t, 3, tickCount)

		timeTicker.Stop()
	})

	t.Run("TimeTicker Stop", func(t *testing.T) {
		tickCount := 0

		timeTicker := NewTimeTicker()

		go func() {
			time.Sleep(time.Second * 2)
			timeTicker.Stop()
		}()

		timeTicker.Start(context.Background(), time.Second*1, func() {
			t.Log("tick")
			tickCount = tickCount + 1
		})

		assert.Greater(t, tickCount, 0)
		assert.True(t, timeTicker.IsStopped())
	})
}

func TestTimeTicker_Tick(t *testing.T) {
	tickCount := 0

	timeTicker := NewTimeTicker()

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	go timeTicker.Start(context, time.Second*1, func() {
		t.Log("tick")
		tickCount = tickCount + 1
	})

	time.Sleep(time.Millisecond * 2500)
	assert.Equal(t, 2, tickCount)
}

func TestTimeTickerWithContextCancel(t *testing.T) {
	tickCount := 0

	timeTicker := NewTimeTicker()

	context, cancel := context.WithCancel(context.Background())

	go timeTicker.Start(context, time.Second*1, func() {
		t.Log("tick")
		tickCount = tickCount + 1
	})

	time.Sleep(time.Millisecond * 1000)
	cancel()
	time.Sleep(time.Millisecond * 1000)

	assert.True(t, timeTicker.IsStopped())
}

func TestTimeTickerStopBeforeWork(t *testing.T) {
	tickCount := 0

	timeTicker := NewTimeTicker()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	timeTicker.Start(ctx, time.Second*1, func() {
		time.Sleep(3 * time.Second)
		t.Log("tick")
		tickCount = tickCount + 1
	})

	assert.True(t, timeTicker.IsStopped())

	assert.Greater(t, tickCount, 0)
}
