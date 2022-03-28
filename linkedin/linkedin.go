package linkedin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"net/http"
)

var NonExistentProfile = errors.New("No LinkedIn Profile Found")

type (
	Service interface {
		GetProfile(GetProfileInput, config.Environment) (GetProfileOutput, error)
	}

	GetProfileInput struct {
		Email string
	}

	GetProfileOutput struct {
		Payload map[string]interface{}
	}

	lkd struct {
		logger zerolog.Logger
	}
)

func New(logger zerolog.Logger) Service {
	return &lkd{
		logger: logger,
	}
}

func (l *lkd) GetProfile(input GetProfileInput, env config.Environment) (GetProfileOutput, error) {
	return l.getProfile(input.Email, env[config.OutlookToken])
}

func (l *lkd) getProfile(email, token string) (GetProfileOutput, error) {
	endpoint := "https://eur.loki.delve.office.com/api/v1/linkedin/profiles/full"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return GetProfileOutput{}, err
	}
	q := req.URL.Query()
	q.Add("Smtp", email)
	req.URL.RawQuery = q.Encode()
	// headers
	req.Header.Add("X-ClientArchitectureVersion", "v1")
	req.Header.Add("X-ClientFeature", "LivePersonaCard")
	req.Header.Add("X-ClientType", "OwaPeopleHub")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetProfileOutput{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GetProfileOutput{}, fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var payload map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return GetProfileOutput{}, err
	}

	if len(payload) <= 0 {
		return GetProfileOutput{}, NonExistentProfile
	}

	return GetProfileOutput{
		payload,
	}, nil
}
