package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/linkedin"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	InternalError = "Something Bad Happened!"
)

func Start(logger zerolog.Logger, env config.Environment, linkedinService linkedin.Service) error {
	e := http.NewServeMux()

	logOutput := os.Stdout

	e.Handle("/", handlers.LoggingHandler(logOutput, rootHandler(logger, env)))
	e.Handle("/auth/linkedin/callback", handlers.LoggingHandler(logOutput, linkedInCallbackHandler(logger, env, linkedinService)))

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

func linkedInCallbackHandler(logger zerolog.Logger, env config.Environment, service linkedin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errDesc := r.URL.Query().Get("error_description")
			http.Error(w, errDesc, http.StatusUnauthorized)
			return
		}
		logger.Info().Str("Code", code).Msg("Successfully received code in redirect url")

		// get access token
		endpoint := "https://www.linkedin.com/oauth/v2/accessToken"
		redirectURL := fmt.Sprintf("%s://%s/auth/linkedin/callback", "http", r.Host)

		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("client_id", env[config.ClientID])
		data.Set("client_secret", env[config.ClientSecret])
		data.Set("redirect_uri", redirectURL)

		req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			logger.Err(err).Msg("Failed to create HTTP request")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Err(err).Msg("Failed to do request")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			logger.Err(fmt.Errorf("expected status code 200, got %d", resp.StatusCode)).Msg("Request failed")
			_, err := io.Copy(os.Stderr, resp.Body)
			if err != nil {
				logger.Err(err).Msg("Failed to write response error")
			}
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var payload AccessTokenResponse

		rawJSON, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Err(err).Msg("Failed to read response body")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(rawJSON, &payload)
		if err != nil {
			logger.Err(err).Msg("Failed to unmarshal response body")
			http.Error(w, "Failed to get access token", http.StatusInternalServerError)
			return
		}

		logger.Info().Str("Token", payload.AccessToken).Int("Expires", payload.ExpiresIn).Msg("Successfully received access token")

		// Get user email
		email, err := getUserEmail(payload.AccessToken)
		if err != nil {
			logger.Err(err).Msg("Failed to get user email")
			http.Error(w, "Failed to get user email", http.StatusInternalServerError)
			return
		}

		logger.Info().Str("Email", email).Msg("Successfully received user email")

		input := linkedin.GetProfileInput{Email: email}
		out, err := service.GetProfile(input)
		if err != nil {
			logger.Err(err).Msg("Failed to get user profile data")
			http.Error(w, "Failed to get data from proxycurl", http.StatusInternalServerError)
			return
		}

		d := map[string]interface{}{
			"fullName":       out.FullName,
			"email":          out.Email,
			"location":       out.Location,
			"phone":          out.Phone,
			"photo":          out.ProfilePhoto,
			"workExperience": out.WorkExperience,
		}
		dpayload, err := json.Marshal(d)
		if err != nil {
			logger.Err(err).Msg("Failed to marshal final payload")
			http.Error(w, InternalError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(dpayload))
	}
}

func getUserEmail(token string) (string, error) {
	endpoint := "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(rawJSON))

	var payload EmailResponse
	err = json.Unmarshal(rawJSON, &payload)
	if err != nil {
		return "", err
	}

	if len(payload.Elements) <= 0 {
		return "", fmt.Errorf("got empty email list")
	}

	return payload.Elements[0].HandleContent.EmailAddress, nil
}
