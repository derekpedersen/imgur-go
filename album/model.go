package album

// Album struct represents the imgur album
type Album struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description interface{} `json:"description"`
	Images      []struct {
		ID          string      `json:"id"`
		Link        string      `json:"link"`
		Title       interface{} `json:"title"`
		Name        interface{} `json:"name"`
		Description string      `json:"description"`
	} `json:"images"`
}

// AlbumResponse struct represents the imgur api response for an Album request
type AlbumResponse struct {
	Data    Album
	Status  int
	Success bool
}
