package authorization_test

import (
	"os"
	"testing"

	"github.com/derekpedersen/imgur-go/authorization"
)

func Test_GenerateAccessToken(t *testing.T) {
	if len(os.Getenv("IMGUR_REFRESH_TOKEN")) == 0 ||
		len(os.Getenv("IMGUR_CLIENT_ID")) == 0 ||
		len(os.Getenv("IMGUR_CLIENT_SECRET")) == 0 {
		t.Skip("skipping integration test: Imgur credentials are not set")
	}

	// Arrange

	// Act

	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatal(err)
	}

	// Assert

	if auth == nil || len(auth.ImgurTokenResponse.AccessToken) <= 0 {
		if auth == nil || auth.ImgurTokenResponse == nil {
			t.Fatalf("expected token response to be set")
		}
		t.Fatalf("expected access token to be set")
	}
}
