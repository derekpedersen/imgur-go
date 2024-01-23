package authorization

import "os"

// Authorization struct for imgur credentials
type Authorization struct {
	ClientID     string
	ClientSecret string
	RefreshToken string

	ImgurTokenResponse *ImgurTokenResponse
}

func NewAuthorization() (
	auth *Authorization,
	err error,
) {

	auth = &Authorization{
		ClientID:     os.Getenv("IMGUR_CLIENT_ID"),
		ClientSecret: os.Getenv("IMGUR_CLIENT_SECRET"),
		RefreshToken: os.Getenv("IMGUR_REFRESH_TOKEN"),
	}

	if len(auth.RefreshToken) > 0 {

		if auth.ImgurTokenResponse, err = auth.GenerateAccessToken(); err != nil {
			return nil, err
		}
	}

	return auth, nil
}
