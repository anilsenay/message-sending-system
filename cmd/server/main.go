package main

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := app.NewServer(ctx)
	_ = server.Listen()
}
