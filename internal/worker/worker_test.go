package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/anilsenay/message-sending-system/internal/worker"
	"github.com/anilsenay/message-sending-system/pkg/ticker"
	"github.com/stretchr/testify/assert"
)

func TestMessageSender(t *testing.T) {
	ticker := ticker.NewTimeTicker()
	period := 1 * time.Second

	sender := worker.NewMessageSender(ticker, period)

	sender.Start(context.Background())
	time.Sleep(5 * time.Second)
	sender.Stop()
	time.Sleep(5 * time.Second)
	assert.False(t, sender.IsRunning())
}
