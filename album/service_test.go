package album_test

import (
	"os"
	"testing"

	"github.com/derekpedersen/imgur-go/album"
	"github.com/derekpedersen/imgur-go/authorization"
)

func TestGetAlbum(t *testing.T) {
	if len(os.Getenv("IMGUR_REFRESH_TOKEN")) == 0 {
		t.Skip("skipping integration test: IMGUR_REFRESH_TOKEN is not set")
	}

	// Arrange
	albumID := "4TZhhtk"
	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatalf("failed to initialize authorization: %v", err)
	}
	alSvc := album.NewAlbumService(*auth, "https://api.imgur.com/3/album/")

	// Act
	album, err := alSvc.GetAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if album == nil {
		t.Fatalf("No Album Returned")
	}
	if len(album.ID) == 0 {
		t.Errorf("No Album Returned")
	}
}

func TestQueryAlbum(t *testing.T) {
	if len(os.Getenv("IMGUR_REFRESH_TOKEN")) == 0 {
		t.Skip("skipping integration test: IMGUR_REFRESH_TOKEN is not set")
	}

	// Arrange
	albumID := "PIRuI"
	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatalf("failed to initialize authorization: %v", err)
	}
	alSvc := album.NewAlbumService(*auth, "https://api.imgur.com/3/album/")

	// Act
	albumJSON, err := alSvc.QueryAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if albumJSON == nil {
		t.Fatalf("No Album Returned")
	}
	if len(*albumJSON) == 0 {
		t.Errorf("No Album Returned")
	}
}
