package album

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

// AlbumService interface
type AlbumService interface {
	QueryAlbum(albumHash string) (json *string, err error)
	GetAlbum(albumHash string) (*Album, error)
}

// AlbumServiceImpl struct
type AlbumServiceImpl struct {
	auth   authorization.Authorization
	url    string
	client *imgur.Client
}

// NewAlbumService creates a new album service
func NewAlbumService(
	auth authorization.Authorization,
	apiURL string,
) *AlbumServiceImpl {
	return &AlbumServiceImpl{
		auth:   auth,
		url:    apiURL,
		client: nil,
	}
}

// NewAlbumServiceWithClient creates a new album service using the shared Imgur client.
func NewAlbumServiceWithClient(client *imgur.Client) *AlbumServiceImpl {
	return &AlbumServiceImpl{client: client}
}

// QueryAlbum queries an album
func (svc *AlbumServiceImpl) QueryAlbum(
	albumHash string,
) (
	json *string,
	err error,
) {

	log.Infof("Querying Album: %s", albumHash)

	if svc.client != nil {
		req, err := svc.client.NewRequest(context.Background(), http.MethodGet, "/3/album/"+albumHash, nil)
		if err != nil {
			log.Errorf("Error creating client request:\n %v", err)
			return nil, err
		}

		res, err := svc.client.Do(req)
		if err != nil {
			log.Errorf("Error making client request:\n %v", err)
			return nil, err
		}

		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Errorf("Error reading client res.Body:\n %v", err)
			return nil, err
		}

		body := string(b)
		json = &body

		log.Debugf("Album JSON: %v", *json)

		return json, nil
	}

	req, err := http.NewRequest("GET", svc.url+albumHash, nil)
	if err != nil {
		log.Errorf("Error creating request:\n %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", " Bearer "+svc.auth.ImgurTokenResponse.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error making request:\n %v", err)
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading res.Body:\n %v", err)
		return nil, err
	}

	body := string(b)

	json = &body

	log.Debugf("Album JSON: %v", *json)

	return json, nil
}

// GetAlbum gets the album
func (svc *AlbumServiceImpl) GetAlbum(albumHash string) (*Album, error) {
	res := AlbumResponse{}
	albumJSON, err := svc.QueryAlbum(albumHash)
	if err != nil {
		log.Errorf("Error querying album:\n %v", err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(*albumJSON), &res); err != nil {
		log.Errorf("Error unmarshalling album response:\n %v", err)
		return nil, err
	}

	return &res.Data, nil
}

func (svc *AlbumServiceImpl) requireClient() error {
	if svc.client == nil {
		return fmt.Errorf("this operation requires NewAlbumServiceWithClient")
	}

	return nil
}

func (svc *AlbumServiceImpl) doForm(method, path string, values url.Values) (*string, error) {
	if err := svc.requireClient(); err != nil {
		return nil, err
	}

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
func (svc *AlbumServiceImpl) CreateAlbum(req CreateAlbumRequest) (*string, error) {
	values := url.Values{}
	setCommonAlbumValues(values, req.Title, req.Description, req.Privacy, req.Cover, req.IDs, req.Layout)

	responseBody, err := svc.doForm(http.MethodPost, "/3/album", values)
	if err != nil {
		return nil, err
	}

	return parseIDResponse(responseBody)
}

// UpdateAlbum updates metadata or image membership for an existing album.
func (svc *AlbumServiceImpl) UpdateAlbum(albumHash string, req UpdateAlbumRequest) (bool, error) {
	values := url.Values{}
	setCommonAlbumValues(values, req.Title, req.Description, req.Privacy, req.Cover, req.IDs, req.Layout)

	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash, values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// DeleteAlbum removes an existing album.
func (svc *AlbumServiceImpl) DeleteAlbum(albumHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodDelete, "/3/album/"+albumHash, url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// FavoriteAlbum toggles the favorite status of an album.
func (svc *AlbumServiceImpl) FavoriteAlbum(albumHash string) (bool, error) {
	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash+"/favorite", url.Values{})
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// AddImagesToAlbum adds image IDs to an album.
func (svc *AlbumServiceImpl) AddImagesToAlbum(albumHash string, imageIDs []string) (bool, error) {
	values := url.Values{}
	values.Set("ids", strings.Join(imageIDs, ","))

	responseBody, err := svc.doForm(http.MethodPost, "/3/album/"+albumHash+"/add", values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}

// RemoveImagesFromAlbum removes image IDs from an album.
func (svc *AlbumServiceImpl) RemoveImagesFromAlbum(albumHash string, imageIDs []string) (bool, error) {
	values := url.Values{}
	values.Set("ids", strings.Join(imageIDs, ","))

	responseBody, err := svc.doForm(http.MethodDelete, "/3/album/"+albumHash+"/remove_images", values)
	if err != nil {
		return false, err
	}

	return parseBoolResponse(responseBody)
}
