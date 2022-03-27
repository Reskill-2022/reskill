package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/linkedin"
	"github.com/thealamu/linkedinsignin/server"
	"os"
)

var defaultWriter = zerolog.ConsoleWriter{Out: os.Stdout}

func main() {
	appLogger := zerolog.New(defaultWriter).With().Timestamp().Logger()

	apiKey, err := getAPIKey()
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to get API key")
	}
	linkedinService := linkedin.New(appLogger, apiKey)

	if err := server.Start(appLogger, linkedinService); err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func getAPIKey() (string, error) {
	if key, ok := os.LookupEnv("PROXYCURL_APIKEY"); ok {
		return key, nil
	}
	return "", fmt.Errorf("can't find proxycurl API key in environment")
}
