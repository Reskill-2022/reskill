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
		FullName       string
		ProfilePhoto   string
		Location       string
		Email          string
		Phone          []string
		WorkExperience interface{}
	}

	UserProfileResponse struct {
		Persons []struct {
			DisplayName  string   `json:"displayName"`
			PhoneNumbers []string `json:"phoneNumbers"`
			Location     string   `json:"location"`
			PhotoURL     string   `json:"photoUrl"`
			LinkedInURL  string   `json:"linkedInUrl"`
			Positions    struct {
				PositionHistory []struct {
					Title string `json:"title"`
				} `json:"positionHistory"`
			} `json:"positions"`
		} `json:"persons"`
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

	var payload UserProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return GetProfileOutput{}, err
	}

	if len(payload.Persons) <= 0 {
		return GetProfileOutput{}, NonExistentProfile
	}

	person := payload.Persons[0]
	return GetProfileOutput{
		person.DisplayName,
		person.PhotoURL,
		person.Location,
		email,
		person.PhoneNumbers,
		person.Positions.PositionHistory,
	}, nil
}
