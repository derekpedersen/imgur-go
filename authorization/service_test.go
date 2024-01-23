package authorization_test

import (
	"testing"

	"github.com/derekpedersen/imgur-go/authorization"
)

func Test_GenerateAccessToken(t *testing.T) {

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
