package app

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/config"
	"github.com/anilsenay/message-sending-system/pkg/logger"
	"github.com/anilsenay/message-sending-system/pkg/rest"
)

type Server struct {
	restApp *rest.Rest
}

func NewServer(ctx context.Context) *Server {
	logger.SetLogLevel(config.LOG_LEVEL)

	restApp := rest.NewRest(rest.Config{
		AppName:          config.APP_NAME,
		Version:          config.APP_VERSION,
		Host:             config.SERVER_HOST,
		Port:             config.SERVER_PORT,
		SwaggerEndpoints: messageHandler.GetSwaggerEndpoints(),
		OnShutdown: func() {
			messageService.StopMessageSending()
			db.Close()
		},
	})

	// set routes
	fiber := restApp.GetFiber()
	messageHandler.SetRoutes(fiber)

	return &Server{
		restApp: restApp,
	}
}

func (s *Server) Listen() error {
	// start message sender
	messageService.StartMessageSending(context.Background())

	// start http listener
	return s.restApp.Listen()
}

func (s *Server) Stop() {
	s.restApp.Stop()
}
