package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"html/template"
	"net/http"
	"os"
	"time"
)

const (
	InternalError = "Something Bad Happened!"
)

func Start(logger zerolog.Logger, env config.Environment, linkedinService linkedin.Service) error {
	e := http.NewServeMux()

	logOutput := os.Stdout

	e.Handle("/", handlers.LoggingHandler(logOutput, rootHandler(logger, env)))
	e.Handle("/auth/linkedin/callback", handlers.LoggingHandler(logOutput, linkedInCallbackHandler(logger)))

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

func rootHandler(logger zerolog.Logger, env config.Environment) http.HandlerFunc {
	tmpl := template.New("index.html")

	tmpl, err := template.ParseFiles("server/templates/index.html")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse template file")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := fmt.Sprintf("%s://%s/auth/linkedin/callback", "http", r.Host)
		authAPI := fmt.Sprintf("https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&scope=r_emailaddress",
			env[config.ClientID], redirectURL)

		data := IndexTemplateData{
			AuthAPI: authAPI,
		}

		err := tmpl.Execute(w, data)
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
