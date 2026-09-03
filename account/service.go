package account

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

import "github.com/derekpedersen/imgur-go/imgur"

// Service provides account operations.
type Service struct {
	client *imgur.Client
}

// NewService creates an account service backed by a shared Imgur client.
func NewService(client *imgur.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("account service requires a shared client")
	}

	return &Service{client: client}, nil
}

// get executes a GET request against the Imgur API and returns the response body.
func (svc *Service) get(path string) ([]byte, error) {
	req, err := svc.client.NewRequest(context.Background(), http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetAccount returns profile information for the requested account username.
func (svc *Service) GetAccount(username string) (*Account, error) {
	body, err := svc.get("/3/account/" + username)
	if err != nil {
		return nil, err
	}

	res := AccountResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// GetFavorites returns a page of favorited items for a user.
//
// If page is negative it defaults to 0. If sort is empty it defaults to "newest".
func (svc *Service) GetFavorites(username string, page int, sort string) ([]AccountItem, error) {
	if page < 0 {
		page = 0
	}
	if sort == "" {
		sort = "newest"
	}

	path := "/3/account/" + username + "/favorites/" + strconv.Itoa(page) + "/" + sort
	body, err := svc.get(path)
	if err != nil {
		return nil, err
	}

	res := AccountItemsResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetSubmissions returns a page of uploaded items for a user.
//
// If page is negative it defaults to 0.
func (svc *Service) GetSubmissions(username string, page int) ([]AccountItem, error) {
	if page < 0 {
		page = 0
	}

	path := "/3/account/" + username + "/submissions/" + strconv.Itoa(page)
	body, err := svc.get(path)
	if err != nil {
		return nil, err
	}

	res := AccountItemsResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetSettings returns account settings for the requested username.
func (svc *Service) GetSettings(username string) (*AccountSettings, error) {
	body, err := svc.get("/3/account/" + username + "/settings")
	if err != nil {
		return nil, err
	}

	res := AccountSettingsResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}
