package images

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/derekpedersen/imgur-go/authorization"
	"github.com/derekpedersen/imgur-go/imgur"
	log "github.com/sirupsen/logrus"
)

// ImageService describes image operations supported by this package.
type ImageService interface {
	QueryImage(imageHash string) (json *string, err error)
	GetImage(imageHash string) (*Image, error)
	UploadImage(req UploadImageRequest) (*string, error)
	UpdateImage(imageHash string, req UpdateImageRequest) (bool, error)
	DeleteImage(imageHash string) (bool, error)
	FavoriteImage(imageHash string) (bool, error)
}

// ImageServiceImpl is the default image service implementation.
type ImageServiceImpl struct {
	auth   authorization.Authorization
	url    string
	client *imgur.Client
}

// NewImageService creates a new image service.
func NewImageService(auth authorization.Authorization, apiURL string) *ImageServiceImpl {
	return &ImageServiceImpl{auth: auth, url: apiURL}
}

// NewImageServiceWithClient creates an image service backed by shared Imgur client.
func NewImageServiceWithClient(client *imgur.Client) *ImageServiceImpl {
	return &ImageServiceImpl{client: client}
}

// QueryImage queries an image and returns raw JSON.
func (svc *ImageServiceImpl) QueryImage(imageHash string) (*string, error) {
	log.Infof("Querying Image: %s", imageHash)

	if svc.client != nil {
		req, err := svc.client.NewRequest(context.Background(), http.MethodGet, "/3/image/"+imageHash, nil)
		if err != nil {
			return nil, err
		}

		res, err := svc.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		body := string(b)
		return &body, nil
	}

	req, err := http.NewRequest(http.MethodGet, svc.url+imageHash, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", " Bearer "+svc.auth.ImgurTokenResponse.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := string(b)
	return &body, nil
}

// GetImage gets a typed image response.
func (svc *ImageServiceImpl) GetImage(imageHash string) (*Image, error) {
	res := ImageResponse{}
	imageJSON, err := svc.QueryImage(imageHash)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(*imageJSON), &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (svc *ImageServiceImpl) requireClient() error {
	if svc.client == nil {
		return fmt.Errorf("this operation requires NewImageServiceWithClient")
	}

	return nil
}

func (svc *ImageServiceImpl) doForm(method, path string, values url.Values) (*string, error) {
	if err := svc.requireClient(); err != nil {
		return nil, err
	}

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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := string(b)
	return &result, nil
}

func parseIDResponse(responseBody *string) (*string, error) {
	res := IDResponse{}
	if err := json.Unmarshal([]byte(*responseBody), &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func parseBoolResponse(responseBody *string) (bool, error) {
	res := BoolResponse{}
	if err := json.Unmarshal([]byte(*responseBody), &res); err != nil {
		return false, err
	}

	return res.Data, nil
}

// UploadImage uploads an image and returns the created image ID.
func (svc *ImageServiceImpl) UploadImage(req UploadImageRequest) (*string, error) {
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
func (svc *ImageServiceImpl) UpdateImage(imageHash string, req UpdateImageRequest) (bool, error) {
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
func (svc *ImageServiceImpl) DeleteImage(imageHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodDelete, "/3/image/"+imageHash, url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// FavoriteImage favorites/unfavorites an image for the authenticated user.
func (svc *ImageServiceImpl) FavoriteImage(imageHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodPost, "/3/image/"+imageHash+"/favorite", url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}
