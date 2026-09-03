package album

import "github.com/derekpedersen/imgur-go/imgurtypes"

// Image struct represents an image in an Imgur album.
type Image struct {
	ID          string      `json:"id"`
	Link        string      `json:"link"`
	Title       interface{} `json:"title"`
	Name        interface{} `json:"name"`
	Description string      `json:"description"`
}

// Album struct represents the imgur album
type Album struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description interface{} `json:"description"`
	Images      []Image     `json:"images"`
}

// AlbumResponse struct represents the imgur api response for an Album request
type AlbumResponse = imgurtypes.APIResponse[Album]

// CreateAlbumRequest defines fields accepted by Imgur create album endpoint.
type CreateAlbumRequest struct {
	Title       string
	Description string
	Privacy     string
	Cover       string
	IDs         []string
	Layout      string
}

// UpdateAlbumRequest defines mutable album fields.
type UpdateAlbumRequest struct {
	Title       string
	Description string
	Privacy     string
	Cover       string
	IDs         []string
	Layout      string
}

// BoolResponse represents mutation endpoints that return a boolean data value.
type BoolResponse = imgurtypes.APIResponse[bool]

// IDResponse represents create endpoints that return an ID string.
type IDResponse = imgurtypes.APIResponse[string]
