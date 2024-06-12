package worker

import (
	"context"
	"time"

	"github.com/anilsenay/message-sending-system/pkg/logger"
)

type ticker interface {
	Start(ctx context.Context, period time.Duration, callback func())
	Stop()
	IsRunning() bool
}

type MessageSender struct {
	ticker ticker
	period time.Duration
}

func NewMessageSender(ticker ticker, period time.Duration) *MessageSender {
	return &MessageSender{
		ticker: ticker,
		period: period,
	}
}

func (ms *MessageSender) Start(ctx context.Context, processFn func() error) {
	if ms.ticker.IsRunning() { // prevent calling Start again before Stop
		return
	}
	go ms.ticker.Start(ctx, ms.period, func() {
		if err := processFn(); err != nil {
			logger.Error().Msgf("error while processing messages: %s", err.Error())
		}
	})
}

func (ms *MessageSender) Stop() {
	if ms.ticker.IsRunning() {
		ms.ticker.Stop()
	}
}

func (ms *MessageSender) IsRunning() bool {
	return ms.ticker.IsRunning()
}
