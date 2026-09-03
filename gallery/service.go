package gallery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/derekpedersen/imgur-go/imgur"
)

// Service provides gallery operations.
type Service struct {
	client *imgur.Client
}

// NewService creates a gallery service backed by a shared Imgur client.
func NewService(client *imgur.Client) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("gallery service requires a shared client")
	}

	return &Service{client: client}, nil
}

// do executes a request against the gallery endpoints and returns the raw body.
func (svc *Service) do(method, path string, body io.Reader, contentType string) ([]byte, error) {
	req, err := svc.client.NewRequest(context.Background(), method, path, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// defaultListOptions applies Imgur API defaults to unset list options.
func defaultListOptions(options ListOptions) ListOptions {
	if options.Section == "" {
		options.Section = "hot"
	}
	if options.Sort == "" {
		options.Sort = "viral"
	}
	if options.Window == "" {
		options.Window = "day"
	}
	if options.Page < 0 {
		options.Page = 0
	}
	return options
}

// buildGalleryPath builds a path for the /3/gallery listing endpoint.
func buildGalleryPath(base string, options ListOptions) string {
	o := defaultListOptions(options)
	path := base + "/" + o.Section + "/" + o.Sort + "/" + o.Window + "/" + strconv.Itoa(o.Page)
	if o.ShowViral != nil {
		query := url.Values{}
		query.Set("showViral", strconv.FormatBool(*o.ShowViral))
		path = path + "?" + query.Encode()
	}
	return path
}

// GetGallery returns gallery items by section, sort, window, and page.
func (svc *Service) GetGallery(options ListOptions) ([]GalleryItem, error) {
	body, err := svc.do(http.MethodGet, buildGalleryPath("/3/gallery", options), nil, "")
	if err != nil {
		return nil, err
	}

	res := GalleryListResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// SearchGallery searches gallery items by free-text query and list options.
func (svc *Service) SearchGallery(query string, options ListOptions) ([]GalleryItem, error) {
	o := defaultListOptions(options)
	q := url.Values{}
	q.Set("q", query)

	path := "/3/gallery/search/" + o.Sort + "/" + o.Window + "/" + strconv.Itoa(o.Page) + "?" + q.Encode()
	if o.ShowViral != nil {
		path += "&showViral=" + strconv.FormatBool(*o.ShowViral)
	}

	body, err := svc.do(http.MethodGet, path, nil, "")
	if err != nil {
		return nil, err
	}

	res := GalleryListResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetGalleryAlbum returns a gallery album detail by hash.
func (svc *Service) GetGalleryAlbum(albumHash string) (*GalleryAlbum, error) {
	body, err := svc.do(http.MethodGet, "/3/gallery/album/"+albumHash, nil, "")
	if err != nil {
		return nil, err
	}

	res := GalleryAlbumResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// VoteGalleryItem votes on a gallery item. Accepted values are up, down, or veto.
func (svc *Service) VoteGalleryItem(itemHash, vote string) (bool, error) {
	vote = strings.ToLower(strings.TrimSpace(vote))
	if vote != "up" && vote != "down" && vote != "veto" {
		return false, fmt.Errorf("invalid vote value: %s", vote)
	}

	body, err := svc.do(http.MethodPost, "/3/gallery/"+itemHash+"/vote/"+vote, strings.NewReader(""), "application/x-www-form-urlencoded")
	if err != nil {
		return false, err
	}

	res := BoolResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return false, err
	}

	return res.Data, nil
}
