package worker

import (
	"context"
	"fmt"
	"time"
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

func (ms *MessageSender) Start(ctx context.Context) {
	if ms.ticker.IsRunning() { // prevent calling Start again before Stop
		return
	}
	go ms.ticker.Start(ctx, ms.period, ms.run)
}

func (ms *MessageSender) Stop() {
	if ms.ticker.IsRunning() {
		ms.ticker.Stop()
	}
}

func (ms *MessageSender) IsRunning() bool {
	return ms.ticker.IsRunning()
}

func (ms *MessageSender) run() {
	fmt.Println("retrieving 2 messages from db...")
	fmt.Println("sending messages...")
}
