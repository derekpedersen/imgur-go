package imgur

import "fmt"

// APIError represents a non-2xx response from Imgur.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("imgur API error: status=%d body=%s", e.StatusCode, e.Body)
}
