package album_test

import (
	"testing"

	"github.com/derekpedersen/imgur-go/album"
	"github.com/derekpedersen/imgur-go/authorization"
)

func TestGetAlbum(t *testing.T) {
	// Arrange
	albumID := "PIRuI"
	auth, _ := authorization.NewAuthorization()
	alSvc := album.NewAlbumService(*auth, "https://api.imgur.com/3/album/")

	// Act
	album, err := alSvc.GetAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if len(album.ID) == 0 {
		t.Errorf("No Album Returned")
	}
}

func TestQueryAlbum(t *testing.T) {
	// Arrange
	albumID := "PIRuI"
	auth, _ := authorization.NewAuthorization()
	alSvc := album.NewAlbumService(*auth, "https://api.imgur.com/3/album/")

	// Act
	album, err := alSvc.QueryAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if len(*album) == 0 {
		t.Errorf("No Album Returned")
	}
}
