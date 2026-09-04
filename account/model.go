package account

import "github.com/derekpedersen/imgur-go/imgurtypes"

// Account represents an Imgur account profile.
type Account struct {
	ID          int64       `json:"id"`
	URL         string      `json:"url"`
	Bio         interface{} `json:"bio"`
	Reputation  int64       `json:"reputation"`
	Created     int64       `json:"created"`
	Pro         bool        `json:"pro"`
	Avatar      string      `json:"avatar"`
	Cover       string      `json:"cover"`
	IsBlocked   bool        `json:"is_blocked"`
	DisplayName string      `json:"display_name"`
}

// AccountSettings represents account settings.
type AccountSettings struct {
	AccountURL       string `json:"account_url"`
	Email            string `json:"email"`
	PublicImages     bool   `json:"public_images"`
	AlbumPrivacy     string `json:"album_privacy"`
	ProExpiration    bool   `json:"pro_expiration"`
	AcceptedGalleryTerms bool `json:"accepted_gallery_terms"`
}

// AccountItem is a summary item used by favorites/submissions lists.
type AccountItem struct {
	ID      string      `json:"id"`
	Title   interface{} `json:"title"`
	Link    string      `json:"link"`
	IsAlbum bool        `json:"is_album"`
	Views   int64       `json:"views"`
}

// AccountResponse is a typed response for account profile endpoints.
type AccountResponse = imgurtypes.APIResponse[Account]

// AccountSettingsResponse is a typed response for account settings endpoint.
type AccountSettingsResponse = imgurtypes.APIResponse[AccountSettings]

// AccountItemsResponse is a typed response for account list endpoints.
type AccountItemsResponse = imgurtypes.APIResponse[[]AccountItem]
