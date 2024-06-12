package main

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/app"
	"github.com/anilsenay/message-sending-system/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := app.NewServer(ctx)
	err := server.Listen()
	if err != nil {
		logger.Error().Msgf("Listen error: %s", err.Error())
	}
}
