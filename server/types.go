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
	//emailAddress
	//{"elements":[
	//	{
	//	"handle~": {"emailAddress":"faithfulnessalamu@outlook.com"},
	//	"handle":"urn:li:emailAddress:7832573868"
	//	}
	//	]
	//	}
)
