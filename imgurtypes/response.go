package imgurtypes

// APIResponse is the common response wrapper from Imgur APIs.
type APIResponse[T any] struct {
	Data    T    `json:"data"`
	Status  int  `json:"status"`
	Success bool `json:"success"`
}
