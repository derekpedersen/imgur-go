package authorization_test

import (
	"testing"

	"github.com/derekpedersen/imgur-go/authorization"
)

func Test_GenerateAccessToken(t *testing.T) {

	// Arrange

	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatal(err)
	}

	// Act

	accessToken, err := auth.GenerateAccessToken()
	if err != nil {
		t.Fatal(err)
	}

	// Assert

	if accessToken == nil || len(accessToken.AccessToken) <= 0 {
		t.Fatalf("expected access token to be set")
	}
}
