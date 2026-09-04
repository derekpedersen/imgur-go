package imgurtypes

// TokenResponse represents a response from Imgur OAuth token endpoints.
type TokenResponse struct {
	AccessToken     string  `json:"access_token"`
	ExpiresIn       int64   `json:"expires_in"`
	TokenType       string  `json:"token_type"`
	Scope           *string `json:"scope"`
	RefreshToken    string  `json:"refresh_token"`
	AccountId       int64   `json:"account_id"`
	AccountUsername string  `json:"account_username"`
}
