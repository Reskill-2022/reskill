package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"net/http"
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
		Handler:      e,
		Addr:         fmt.Sprintf(":%s", env[config.Port]),
	}

	logger.Info().Str("Address", srv.Addr).Msg("Starting HTTP server")
	return srv.ListenAndServe()
}
