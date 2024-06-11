package rest

import (
	"fmt"

	"github.com/anilsenay/message-sending-system/pkg/graceful"
	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

type Rest struct {
	config Config
	app    *fiber.App
}

type Config struct {
	AppName          string
	Version          string
	Host             string
	Port             int
	SwaggerEndpoints []*endpoint.EndPoint
	OnShutdown       func()
}

func NewRest(cfg Config) *Rest {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
	})

	// recover middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Error().Str("url", c.OriginalURL()).Msgf("%v", e)
		},
	}))

	// healtcheck
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"name":    cfg.AppName,
			"version": cfg.Version,
		})
	})

	// swagger handler
	sw := swagno.New(swagno.Config{Title: cfg.AppName, Version: cfg.Version})
	sw.AddEndpoints(cfg.SwaggerEndpoints)
	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	return &Rest{
		app:    app,
		config: cfg,
	}
}

func (s *Rest) Listen() error {
	// start http listener
	graceful.OnShutdown(s.Stop)

	log.Log().Msgf("App starting on http://%s:%d", s.config.Host, s.config.Port)
	err := s.app.Listen(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		log.Panic().Msgf("Error on Listen: %s", err.Error())
		return err
	}
	graceful.Shutdown()
	log.Log().Msg("App stopped")
	return nil
}

func (s *Rest) Stop() {
	_ = s.app.Shutdown()
	if s.config.OnShutdown != nil {
		s.config.OnShutdown()
	}
}

func (s *Rest) GetFiber() *fiber.App {
	return s.app
}
