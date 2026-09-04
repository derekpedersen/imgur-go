package gallery

import "github.com/derekpedersen/imgur-go/imgurtypes"

// GalleryItem represents a gallery item summary from Imgur.
type GalleryItem struct {
	ID          string      `json:"id"`
	Title       interface{} `json:"title"`
	Description interface{} `json:"description"`
	Link        string      `json:"link"`
	Type        string      `json:"type"`
	IsAlbum     bool        `json:"is_album"`
	Views       int64       `json:"views"`
	Ups         int64       `json:"ups"`
	Downs       int64       `json:"downs"`
	Score       int64       `json:"score"`
}

// GalleryAlbumImage represents an image nested in a gallery album.
type GalleryAlbumImage struct {
	ID          string      `json:"id"`
	Title       interface{} `json:"title"`
	Description interface{} `json:"description"`
	Link        string      `json:"link"`
	Type        string      `json:"type"`
}

// GalleryAlbum represents album detail from gallery endpoints.
type GalleryAlbum struct {
	ID          string             `json:"id"`
	Title       interface{}        `json:"title"`
	Description interface{}        `json:"description"`
	Link        string             `json:"link"`
	Images      []GalleryAlbumImage `json:"images"`
}

// GalleryListResponse is the response for list/search endpoints.
type GalleryListResponse = imgurtypes.APIResponse[[]GalleryItem]

// GalleryAlbumResponse is the response for gallery item detail endpoints.
type GalleryAlbumResponse = imgurtypes.APIResponse[GalleryAlbum]

// BoolResponse is the response for mutation endpoints.
type BoolResponse = imgurtypes.APIResponse[bool]

// ListOptions controls optional list/search query parameters.
type ListOptions struct {
	Section   string
	Sort      string
	Window    string
	Page      int
	ShowViral *bool
}
