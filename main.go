package main

import (
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/server"
	"os"
)

var defaultWriter = zerolog.ConsoleWriter{Out: os.Stdout}

func main() {
	appLogger := zerolog.New(defaultWriter).With().Timestamp().Logger()

	if err := server.Start(); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
