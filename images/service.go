package images

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/derekpedersen/imgur-go/imgur"
)

// Service describes image operations supported by this package.
type Service struct {
	client *imgur.Client
}

// NewService creates an image service backed by a shared Imgur client.
func NewService(client *imgur.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("image service requires a shared client")
	}

	return &Service{client: client}, nil
}

// QueryImage queries an image and returns raw JSON.
func (svc *Service) QueryImage(imageHash string) ([]byte, error) {
	req, err := svc.client.NewRequest(context.Background(), http.MethodGet, "/3/image/"+imageHash, nil)
	if err != nil {
		return nil, err
	}

	res, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

// GetImage gets a typed image response.
func (svc *Service) GetImage(imageHash string) (*Image, error) {
	res := ImageResponse{}
	imageJSON, err := svc.QueryImage(imageHash)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(imageJSON, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (svc *Service) doForm(method, path string, values url.Values) ([]byte, error) {
	body := strings.NewReader(values.Encode())
	req, err := svc.client.NewRequest(context.Background(), method, path, body)
	if err != nil {
		return nil, err
	}
	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func parseIDResponse(responseBody []byte) (*string, error) {
	res := IDResponse{}
	if err := json.Unmarshal(responseBody, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func parseBoolResponse(responseBody []byte) (bool, error) {
	res := BoolResponse{}
	if err := json.Unmarshal(responseBody, &res); err != nil {
		return false, err
	}

	return res.Data, nil
}

// UploadImage uploads an image and returns the created image ID.
func (svc *Service) UploadImage(req UploadImageRequest) (*string, error) {
	values := url.Values{}
	if req.Image != "" {
		values.Set("image", req.Image)
	}
	if req.Type != "" {
		values.Set("type", req.Type)
	}
	if req.Name != "" {
		values.Set("name", req.Name)
	}
	if req.Title != "" {
		values.Set("title", req.Title)
	}
	if req.Description != "" {
		values.Set("description", req.Description)
	}
	if req.Album != "" {
		values.Set("album", req.Album)
	}

	responseBody, err := svc.doForm(http.MethodPost, "/3/image", values)
	if err != nil {
		return nil, err
	}

	return parseIDResponse(responseBody)
}

// UpdateImage updates an image's metadata.
func (svc *Service) UpdateImage(imageHash string, req UpdateImageRequest) (bool, error) {
	values := url.Values{}
	if req.Title != "" {
		values.Set("title", req.Title)
	}
	if req.Description != "" {
		values.Set("description", req.Description)
	}

	responseBody, err := svc.doForm(http.MethodPost, "/3/image/"+imageHash, values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// DeleteImage deletes an image.
func (svc *Service) DeleteImage(imageHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodDelete, "/3/image/"+imageHash, url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// FavoriteImage favorites/unfavorites an image for the authenticated user.
func (svc *Service) FavoriteImage(imageHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodPost, "/3/image/"+imageHash+"/favorite", url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}
