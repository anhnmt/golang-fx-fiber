package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	app *fiber.App
}

func NewServer(lc fx.Lifecycle) *Server {
	// Creates a new Fiber instance.
	app := fiber.New(fiber.Config{
		AppName: "Docker MariaDB Clean Arch",
	})

	// Use global middlewares.
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(limiter.New(limiter.Config{
		Max: 100,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		errorMessage := fmt.Sprintf("Route '%s' does not exist in this API!", c.OriginalURL())

		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "fail",
			"message": errorMessage,
		})
	})

	srv := &Server{
		app: app,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return srv.Run()
		},
		OnStop: func(context.Context) error {
			return srv.Close()
		},
	})

	return srv
}

func (s *Server) Run() error {
	go func() {
		err := s.app.Listen(":3000")
		if err != nil {
			log.Err(err).Msg("Listen failed")
			return
		}
	}()

	return nil
}

func (s *Server) Close() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		return s.app.Shutdown()
	})

	return group.Wait()
}
