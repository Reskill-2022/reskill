package server

type (
	IndexTemplateData struct {
		AuthAPI string
	}

	AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	EmailResponse struct {
		Elements []struct {
			Handle        string `json:"handle"`
			HandleContent struct {
				EmailAddress string `json:"emailAddress"`
			} `json:"handle~"`
		}
	}
)
