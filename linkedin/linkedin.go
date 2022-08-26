package linkedin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"github.com/thealamu/linkedinsignin/config"
	"github.com/thealamu/linkedinsignin/errors"
)

type (
	AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	Service interface {
		GetProfile(authCode, redirectURI string) (*GetProfileOutput, error)
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

	UserPhone struct {
		Number string `json:"number"`
	}

	UserProfileResponse struct {
		Persons []struct {
			DisplayName  string      `json:"displayName"`
			PhoneNumbers []UserPhone `json:"phoneNumbers"`
			Location     string      `json:"location"`
			PhotoURL     string      `json:"photoUrl"`
			LinkedInURL  string      `json:"linkedInUrl"`
			Positions    struct {
				PositionHistory []struct {
					Title string `json:"title"`
				} `json:"positionHistory"`
			} `json:"positions"`
		} `json:"persons"`
	}

	lkd struct {
		logger       zerolog.Logger
		MSAAUTH      string
		token        string
		clientID     string
		clientSecret string
	}

	EmailResponse struct {
		Elements []struct {
			Handle        string `json:"handle"`
			HandleContent struct {
				EmailAddress string `json:"emailAddress"`
			} `json:"handle~"`
		}
	}

	ProfileResponse struct {
		LocalizedLastName  string `json:"localizedLastName"`
		LocalizedFirstName string `json:"localizedFirstName"`
		ProfilePicture     struct {
			DisplayImage string `json:"displayImage"`
		} `json:"profilePicture"`
	}
)

func New(logger zerolog.Logger, env config.Environment) Service {
	return &lkd{
		logger:       logger,
		MSAAUTH:      env["MSAAUTH"],
		clientID:     env["ClientID"],
		clientSecret: env["ClientSecret"],
	}
}

func (l *lkd) GetProfile(authCode, redirectURI string) (*GetProfileOutput, error) {
	l.logger.Debug().Msg("Getting LinkedIn Profile for")
	return l.getProfileNew(authCode, redirectURI)
}

func (l *lkd) getProfileNew(authCode, redirectURI string) (*GetProfileOutput, error) {
	endpoint := "https://www.linkedin.com/oauth/v2/accessToken"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("client_id", l.clientID)
	data.Set("client_secret", l.clientSecret)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		l.logger.Err(err).Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to build request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.logger.Err(err).Msg("Failed to do request")
		return nil, fmt.Errorf("failed to get access token")
	}

	if resp.StatusCode != http.StatusOK {
		l.logger.Err(fmt.Errorf("expected status code 200, got %d", resp.StatusCode)).Msg("Request failed")
		_, err := io.Copy(os.Stderr, resp.Body)
		if err != nil {
			l.logger.Err(err).Msg("Failed to write response error")
		}
		// http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to get access token, not ok")
	}
	defer resp.Body.Close()

	var payload AccessTokenResponse

	rawJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.logger.Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body")
	}
	err = json.Unmarshal(rawJSON, &payload)
	if err != nil {
		l.logger.Err(err).Msg("Failed to unmarshal response body")
		// http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to unmarshal response body")
	}

	email, err := getUserEmail(payload.AccessToken)
	if err != nil {
		return nil, err
	}

	fname, lname, picture, err := getUserProfile(payload.AccessToken)
	if err != nil {
		return nil, err
	}

	return &GetProfileOutput{
		Email: email,
		Name:  fname + " " + lname,
		Photo: picture,
	}, nil
}

func getUserProfile(token string) (string, string, string, error) {
	endpoint := "https://api.linkedin.com/v2/me"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to build request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to do request")
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("failed to get access token, not ok")
	}
	defer resp.Body.Close()

	var payload ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to unmarshal response body")
	}

	return payload.LocalizedFirstName, payload.LocalizedLastName, payload.ProfilePicture.DisplayImage, nil
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

func firstOrEmpty(s []UserPhone) string {
	for _, v := range s {
		if v.Number != "" {
			return v.Number
		}
	}
	return ""
}
