package service_test

import (
	"testing"

	"github.com/derekpedersen/imgur-go/service"
)

func TestGetAlbum(t *testing.T) {
	// Arrange
	albumID := "PIRuI"
	service.NewAlbumService("")

	// Act
	album, err := service.GetAlbum(albumID)

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
	service.NewAlbumService("")

	// Act
	album, err := service.QueryAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if len(album) == 0 {
		t.Errorf("No Album Returned")
	}
}
