package model

// AlbumResponse struct represents the imgur api response for an Album request
type AlbumResponse struct {
	Data    Album
	Status  int
	Success bool
}
