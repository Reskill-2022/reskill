package linkedin

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/errors"
	"net/http"
	"strings"
)

type (
	Service interface {
		GetProfile(email string) (GetProfileOutput, error)
	}

	GetProfileInput struct {
		Email string
	}

	GetProfileOutput struct {
		Email         string
		Name          string
		Photo         string
		ProfileURL    string
		Location      string
		Phone         string
		HasExperience bool
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
		logger  zerolog.Logger
		MSAAUTH string
		token   string
	}
)

func New(logger zerolog.Logger, env config.Environment) Service {
	return &lkd{
		logger:  logger,
		MSAAUTH: env["MSAAUTH"],
	}
}

func (l *lkd) GetProfile(email string) (GetProfileOutput, error) {
	l.logger.Debug().Msg("Getting LinkedIn Profile for " + email)
	return l.getProfile(email)
}

func (l *lkd) getProfile(email string) (GetProfileOutput, error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", l.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetProfileOutput{}, err
	}
	if resp.StatusCode != http.StatusOK {
		l.logger.Debug().Msgf("Got Status Code when getting profile: %d", resp.StatusCode)
		if resp.StatusCode != http.StatusUnauthorized {
			return GetProfileOutput{}, errors.New(fmt.Sprintf("expected 200, got %d", resp.StatusCode), 500)
		}

		token, err := l.getToken()
		if err != nil {
			return GetProfileOutput{}, errors.From(err, "failed to get token", 500)
		}
		l.token = token

		return l.getProfile(email)
	}
	defer resp.Body.Close()

	var payload UserProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return GetProfileOutput{}, errors.From(err, "failed to decode response", 500)
	}

	if len(payload.Persons) <= 0 {
		return GetProfileOutput{}, errors.New("No LinkedIn Profile Found For That Email", 404)
	}

	person := payload.Persons[0]
	return GetProfileOutput{
		Email:         email,
		Name:          person.DisplayName,
		Photo:         person.PhotoURL,
		Location:      person.Location,
		ProfileURL:    person.LinkedInURL,
		Phone:         firstOrEmpty(person.PhoneNumbers),
		HasExperience: len(person.Positions.PositionHistory) > 0,
	}, nil
}

func (l *lkd) getToken() (string, error) {
	endpoint := "https://login.live.com/oauth20_authorize.srf"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("response_type", "token")
	q.Add("prompt", "none")
	q.Add("redirect_uri", "https://outlook.live.com/owa/auth/dt.aspx")
	q.Add("scope", "liveprofilecard.access")
	q.Add("client_id", "292841")
	req.URL.RawQuery = q.Encode()

	// headers
	l.logger.Debug().Msgf("Requesting token")
	// add a cookie
	req.AddCookie(&http.Cookie{
		Name:  "__Host-MSAAUTH",
		Value: l.MSAAUTH,
		Path:  "/",
	})

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		return "", err
	}

	// extract Location header
	location, err := resp.Location()
	if err != nil {
		return "", err
	}

	locationStr := location.String()
	// strip out the access token
	locationStr = locationStr[strings.Index(locationStr, "access_token=")+13:]
	locationStr = locationStr[:strings.Index(locationStr, "&")]

	l.logger.Debug().Msgf("access token: %s", locationStr)
	return locationStr, nil
}

func firstOrEmpty(s []string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}
