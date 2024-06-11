package app

import (
	"context"

	"github.com/anilsenay/message-sending-system/internal/config"
	"github.com/anilsenay/message-sending-system/pkg/rest"
	"github.com/rs/zerolog"
)

type Server struct {
	restApp *rest.Rest
}

func NewServer(ctx context.Context) *Server {
	setLogLevel(config.LOG_LEVEL)

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

func setLogLevel(level string) {
	switch level {
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "DISABLED":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
