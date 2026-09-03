package authorization_test

import (
	"os"
	"testing"

	"github.com/derekpedersen/imgur-go/authorization"
)

func hasIntegrationEnv() bool {
	return os.Getenv("IMGUR_CLIENT_ID") != "" &&
		os.Getenv("IMGUR_CLIENT_SECRET") != "" &&
		os.Getenv("IMGUR_REFRESH_TOKEN") != ""
}

func Test_GenerateAccessToken(t *testing.T) {
	if !hasIntegrationEnv() {
		t.Skip("skipping integration test: IMGUR_* env vars are not set")
	}

	// Arrange

	// Act

	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatal(err)
	}

	// Assert

	if auth == nil || len(auth.ImgurTokenResponse.AccessToken) <= 0 {
		t.Fatalf("expected access token to be set")
	}
}
