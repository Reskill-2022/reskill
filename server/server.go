package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	InternalError = "Something Bad Happened!"
)

func registerRoutes(e *echo.Echo, linkedinService linkedin.Service) {
}

func Start(logger zerolog.Logger, env config.Environment, linkedinService linkedin.Service) error {
	e := echo.New()

	registerRoutes(e, linkedinService)

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Addr:         fmt.Sprintf(":%s", env[config.Port]),
	}

	// gracefully start server
	go func() {
		if err := e.StartServer(srv); err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// listen for ctrl+c and gracefully shutdown server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to shutdown server")
	}

	logger.Info().Msg("Server gracefully shutdown")
	return nil
}
