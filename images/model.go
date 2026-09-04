package images

import "github.com/derekpedersen/imgur-go/imgurtypes"

// Image represents an image object from the Imgur API.
type Image struct {
	ID          string      `json:"id"`
	Title       interface{} `json:"title"`
	Description interface{} `json:"description"`
	Type        string      `json:"type"`
	Animated    bool        `json:"animated"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	Size        int64       `json:"size"`
	Views       int64       `json:"views"`
	Bandwidth   int64       `json:"bandwidth"`
	DeleteHash  string      `json:"deletehash"`
	Name        interface{} `json:"name"`
	Link        string      `json:"link"`
}

// ImageResponse represents a typed API response for a single image.
type ImageResponse = imgurtypes.APIResponse[Image]

// BoolResponse represents mutation endpoints that return true/false.
type BoolResponse = imgurtypes.APIResponse[bool]

// IDResponse represents create endpoints that return an ID.
type IDResponse = imgurtypes.APIResponse[string]

// UploadImageRequest defines fields accepted by Imgur image upload endpoint.
type UploadImageRequest struct {
	Image       string
	Type        string
	Name        string
	Title       string
	Description string
	Album       string
}

// UpdateImageRequest defines mutable fields for an image.
type UpdateImageRequest struct {
	Title       string
	Description string
}
