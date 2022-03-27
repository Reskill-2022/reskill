package server

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"html/template"
	"net/http"
	"time"
)

const (
	InternalError = "Something Bad Happened!"
)

func Start(logger zerolog.Logger, env config.Environment, linkedinService linkedin.Service) error {
	e := http.NewServeMux()
	e.HandleFunc("/", rootHandler(logger))
	e.HandleFunc("/auth/linkedin/callback", linkedInCallbackHandler(logger))

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

func rootHandler(logger zerolog.Logger) http.HandlerFunc {
	tmpl := template.New("index.html")

	tmpl, err := template.ParseFiles("server/templates/index.html")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse template file")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			logger.Err(err).Msg("Can't execute template")
			http.Error(w, InternalError, http.StatusInternalServerError)
			return
		}
	}
}

func linkedInCallbackHandler(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
