package server

type (
	IndexTemplateData struct {
		AuthAPI string
	}

	AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	EmailResponse struct {
		Handle        string `json:"handle"`
		HandleContent struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"handle~"`
	}
)
