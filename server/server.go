package server

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"net/http"
	"time"
)

func Start(logger zerolog.Logger, env config.Environment, linkedinService linkedin.Service) error {
	e := http.NewServeMux()
	e.HandleFunc("/", rootHandler)

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      e,
		Addr:         fmt.Sprintf(":%s", env[config.Port]),
	}

	return srv.ListenAndServe()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}
