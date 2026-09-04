package album

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

// Service provides album operations.
type Service struct {
	client *imgur.Client
}

// NewService creates an album service using a shared Imgur client.
func NewService(client *imgur.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("album service requires a shared client")
	}

	return &Service{client: client}, nil
}

// QueryAlbum fetches a raw album API payload for the provided album hash.
func (svc *Service) QueryAlbum(albumHash string) ([]byte, error) {
	req, err := svc.client.NewRequest(context.Background(), http.MethodGet, "/3/album/"+albumHash, nil)
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

// GetAlbum fetches and unmarshals album details for the provided album hash.
func (svc *Service) GetAlbum(albumHash string) (*Album, error) {
	res := AlbumResponse{}
	albumJSON, err := svc.QueryAlbum(albumHash)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(albumJSON, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// doForm executes form-encoded album mutation requests and returns the raw body.
func (svc *Service) doForm(method, path string, values url.Values) ([]byte, error) {
	body := strings.NewReader(values.Encode())
	req, err := svc.client.NewRequest(context.Background(), method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// parseIDResponse parses standard Imgur string ID responses.
func parseIDResponse(responseBody []byte) (*string, error) {
	res := IDResponse{}
	if err := json.Unmarshal(responseBody, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// parseBoolResponse parses standard Imgur boolean responses.
func parseBoolResponse(responseBody []byte) (bool, error) {
	res := BoolResponse{}
	if err := json.Unmarshal(responseBody, &res); err != nil {
		return false, err
	}

	return res.Data, nil
}

// setCommonAlbumValues adds non-empty request values to a form payload.
func setCommonAlbumValues(values url.Values, title, description, privacy, cover string, ids []string, layout string) {
	if title != "" {
		values.Set("title", title)
	}
	if description != "" {
		values.Set("description", description)
	}
	if privacy != "" {
		values.Set("privacy", privacy)
	}
	if cover != "" {
		values.Set("cover", cover)
	}
	if len(ids) > 0 {
		values.Set("ids", strings.Join(ids, ","))
	}
	if layout != "" {
		values.Set("layout", layout)
	}
}

// CreateAlbum creates a new album and returns the created album ID.
func (svc *Service) CreateAlbum(req CreateAlbumRequest) (*string, error) {
	values := url.Values{}
	setCommonAlbumValues(values, req.Title, req.Description, req.Privacy, req.Cover, req.IDs, req.Layout)

	responseBody, err := svc.doForm(http.MethodPost, "/3/album", values)
	if err != nil {
		return nil, err
	}

	return parseIDResponse(responseBody)
}

// UpdateAlbum updates metadata or image membership for an existing album.
func (svc *Service) UpdateAlbum(albumHash string, req UpdateAlbumRequest) (bool, error) {
	values := url.Values{}
	setCommonAlbumValues(values, req.Title, req.Description, req.Privacy, req.Cover, req.IDs, req.Layout)

	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash, values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// DeleteAlbum removes an existing album.
func (svc *Service) DeleteAlbum(albumHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodDelete, "/3/album/"+albumHash, url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// FavoriteAlbum toggles the favorite status of an album.
func (svc *Service) FavoriteAlbum(albumHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash+"/favorite", url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// AddImagesToAlbum adds image IDs to an album.
func (svc *Service) AddImagesToAlbum(albumHash string, imageIDs []string) (bool, error) {
	values := url.Values{}
	values.Set("ids", strings.Join(imageIDs, ","))

	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash+"/add", values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// RemoveImagesFromAlbum removes image IDs from an album.
func (svc *Service) RemoveImagesFromAlbum(albumHash string, imageIDs []string) (bool, error) {
	values := url.Values{}
	values.Set("ids", strings.Join(imageIDs, ","))

	responseBody, err := svc.doForm(http.MethodDelete, "/3/album/"+albumHash+"/remove_images", values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}
