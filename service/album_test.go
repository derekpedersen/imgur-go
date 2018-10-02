package service_test

import (
	"os"
	"testing"

	"github.com/derekpedersen/imgur-go/service"
)

func TestGetAlbum(t *testing.T) {
	// Arrange
	albumID := "PIRuI"
	apiKey := os.Getenv("IMGUR_API_KEY")
	alSvc := service.NewAlbumService(apiKey)

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
	apiKey := os.Getenv("IMGUR_API_KEY")
	alSvc := service.NewAlbumService(apiKey)

	// Act
	album, err := alSvc.QueryAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if len(album) == 0 {
		t.Errorf("No Album Returned")
	}
}
