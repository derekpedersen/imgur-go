package authorization

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ImgurTokenResponse struct {
	AccessToken     string  `json:"access_token"`
	ExpiresIn       int64   `json:"expires_in"`
	TokenType       string  `json:"token_type"`
	Scope           *string `json:"scope"`
	RefreshToken    string  `json:"refresh_token"`
	AccountId       int64   `json:"account_id"`
	AccountUsername string  `json:"account_username"`
}

func (a *Authorization) GenerateAccessToken() (
	imgurTokenResponse *ImgurTokenResponse,
	err error,
) {

	form := new(bytes.Buffer)

	writer := multipart.NewWriter(form)

	formField, err := writer.CreateFormField("refresh_token")
	if err != nil {
		return nil, err
	}

	_, err = formField.Write([]byte(a.RefreshToken))
	if err != nil {
		return nil, err
	}

	formField, err = writer.CreateFormField("client_id")
	if err != nil {
		return nil, err
	}

	_, err = formField.Write([]byte(a.ClientID))
	if err != nil {
		return nil, err
	}

	formField, err = writer.CreateFormField("client_secret")
	if err != nil {
		return nil, err
	}

	_, err = formField.Write([]byte(a.ClientSecret))
	if err != nil {
		return nil, err
	}

	formField, err = writer.CreateFormField("grant_type")
	if err != nil {
		return nil, err
	}

	_, err = formField.Write([]byte("refresh_token"))
	if err != nil {
		return nil, err
	}

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.imgur.com/oauth2/token", form)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	imgurTokenResponse = &ImgurTokenResponse{}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bodyText, imgurTokenResponse); err != nil {
		return nil, err
	}

	logrus.Debugf("imgur auth: %v", string(bodyText))

	return imgurTokenResponse, nil
}

func (a *Authorization) SetAccessToken() (
	imgurTokenResponse *ImgurTokenResponse,
	err error,
) {

	resp, err := a.GenerateAccessToken()
	if err != nil {
		return nil, err
	}

	a.ImgurTokenResponse = resp

	return imgurTokenResponse, nil
}
