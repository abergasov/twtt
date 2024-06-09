package routes

import (
	"fmt"
	"net/http"
	"twtt/internal/logger"
	"twtt/internal/service/indexator"
	"twtt/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	appAddr          string
	log              logger.AppLogger
	serviceIndexator *indexator.Service
	httpEngine       *fiber.App
}

// InitAppRouter initializes the HTTP Server.
func InitAppRouter(log logger.AppLogger, srv *indexator.Service, address string) *Server {
	app := &Server{
		appAddr:          address,
		httpEngine:       fiber.New(fiber.Config{}),
		serviceIndexator: srv,
		log:              log.With(logger.WithString("service", "http")),
	}
	app.httpEngine.Use(recover.New())
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.httpEngine.Get("/get_current_block", func(ctx *fiber.Ctx) error {
		return ctx.SendString(fmt.Sprintf("last observed block: %d", s.serviceIndexator.GetCurrentBlock()))
	})
	s.httpEngine.Get("/subscribe", func(ctx *fiber.Ctx) error {
		address := ctx.Query("address")
		if utils.ValidateETHAddress(address) != nil {
			return ctx.Status(http.StatusBadRequest).SendString("invalid address")
		}
		if s.serviceIndexator.Subscribe(address) {
			return ctx.Status(http.StatusOK).SendString("subscribed")
		}
		return ctx.Status(http.StatusNotModified).SendString("already subscribed")
	})
	s.httpEngine.Get("/get_transactions", func(ctx *fiber.Ctx) error {
		address := ctx.Query("address")
		if utils.ValidateETHAddress(address) != nil {
			return ctx.Status(http.StatusBadRequest).SendString("invalid address")
		}
		txList, err := s.serviceIndexator.GetTransactions(address)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return ctx.JSON(txList)
	})
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", logger.WithString("port", s.appAddr))
	return s.httpEngine.Listen(s.appAddr)
}

func (s *Server) Stop() error {
	return s.httpEngine.Shutdown()
}
