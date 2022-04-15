package main

import (
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/controllers"
	"github.com/thealamu/linkedinsignin/email"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/repository"
	"github.com/thealamu/linkedinsignin/server"
	"os"
)

var defaultWriter = zerolog.ConsoleWriter{Out: os.Stdout}

func main() {
	appLogger := zerolog.New(defaultWriter).With().Timestamp().Logger()

	env, err := config.New()
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to load configs")
	}

	cts := controllers.NewContainer(appLogger)
	rc := repository.NewContainer(appLogger)
	service := linkedin.New(appLogger, env)
	ses := email.New(appLogger)

	if err := server.Start(appLogger, env, cts, rc, service, ses); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
