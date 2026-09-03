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

// Service defines account operations.
type Service interface {
	GetAccount(username string) (*Account, error)
	GetFavorites(username string, page int, sort string) ([]AccountItem, error)
	GetSubmissions(username string, page int) ([]AccountItem, error)
	GetSettings(username string) (*AccountSettings, error)
}

// ServiceImpl is the default account service implementation.
type ServiceImpl struct {
	client *imgur.Client
}

// NewService creates an account service backed by a shared Imgur client.
func NewService(client *imgur.Client) *ServiceImpl {
	return &ServiceImpl{client: client}
}

func (svc *ServiceImpl) requireClient() error {
	if svc.client == nil {
		return fmt.Errorf("account service requires a shared client")
	}
	return nil
}

func (svc *ServiceImpl) get(path string) (*string, error) {
	if err := svc.requireClient(); err != nil {
		return nil, err
	}

	req, err := svc.client.NewRequest(context.Background(), http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := string(b)
	return &result, nil
}

// GetAccount returns account profile information.
func (svc *ServiceImpl) GetAccount(username string) (*Account, error) {
	body, err := svc.get("/3/account/" + username)
	if err != nil {
		return nil, err
	}

	res := AccountResponse{}
	if err := json.Unmarshal([]byte(*body), &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// GetFavorites returns a paged list of user favorites.
func (svc *ServiceImpl) GetFavorites(username string, page int, sort string) ([]AccountItem, error) {
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
	if err := json.Unmarshal([]byte(*body), &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetSubmissions returns a paged list of user submissions.
func (svc *ServiceImpl) GetSubmissions(username string, page int) ([]AccountItem, error) {
	if page < 0 {
		page = 0
	}

	path := "/3/account/" + username + "/submissions/" + strconv.Itoa(page)
	body, err := svc.get(path)
	if err != nil {
		return nil, err
	}

	res := AccountItemsResponse{}
	if err := json.Unmarshal([]byte(*body), &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetSettings returns account settings for a user.
func (svc *ServiceImpl) GetSettings(username string) (*AccountSettings, error) {
	body, err := svc.get("/3/account/" + username + "/settings")
	if err != nil {
		return nil, err
	}

	res := AccountSettingsResponse{}
	if err := json.Unmarshal([]byte(*body), &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}
