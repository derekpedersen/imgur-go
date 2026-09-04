package album_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/derekpedersen/imgur-go/album"
	"github.com/derekpedersen/imgur-go/authorization"
	"github.com/derekpedersen/imgur-go/imgur"
)

func TestNewService(t *testing.T) {
	validClient, err := imgur.NewClient(imgur.Config{BaseURL: "https://api.imgur.com", ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	tests := []struct {
		name    string
		client  *imgur.Client
		wantErr bool
	}{
		{name: "nil client", client: nil, wantErr: true},
		{name: "valid client", client: validClient, wantErr: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc, err := album.NewService(tc.client)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected constructor error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected constructor error: %v", err)
			}
			if svc == nil {
				t.Fatal("expected service to be initialized")
			}
		})
	}
}

func hasIntegrationEnv() bool {
	return os.Getenv("IMGUR_CLIENT_ID") != "" &&
		os.Getenv("IMGUR_CLIENT_SECRET") != "" &&
		os.Getenv("IMGUR_REFRESH_TOKEN") != ""
}

func newOAuthClient(t *testing.T) *imgur.Client {
	t.Helper()

	auth, err := authorization.NewAuthorization()
	if err != nil {
		t.Fatalf("failed to initialize authorization: %v", err)
	}
	if auth.ImgurTokenResponse == nil || auth.ImgurTokenResponse.AccessToken == "" {
		t.Fatal("expected oauth access token to be initialized")
	}

	client, err := imgur.NewClient(imgur.Config{
		ClientID:    auth.ClientID,
		AccessToken: auth.ImgurTokenResponse.AccessToken,
		Mode:        imgur.AuthModeOAuth,
	})
	if err != nil {
		t.Fatalf("failed to create oauth client: %v", err)
	}

	return client
}

func TestGetAlbum(t *testing.T) {
	if !hasIntegrationEnv() {
		t.Skip("skipping integration test: IMGUR_* env vars are not set")
	}

	// Arrange
	albumID := "4TZhhtk"
	client := newOAuthClient(t)
	alSvc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("failed to initialize album service: %v", err)
	}

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
	if !hasIntegrationEnv() {
		t.Skip("skipping integration test: IMGUR_* env vars are not set")
	}

	// Arrange
	albumID := "PIRuI"
	client := newOAuthClient(t)
	alSvc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("failed to initialize album service: %v", err)
	}

	// Act
	albumJSON, err := alSvc.QueryAlbum(albumID)

	// Assert
	if err != nil {
		t.Errorf("Experienced an error: %v", err)
	}
	if len(albumJSON) == 0 {
		t.Errorf("No Album Returned")
	}
}

func TestQueryAlbumWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/3/album/testhash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("Authorization"); got != "Client-ID test-client-id" {
			t.Fatalf("unexpected auth header: %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"testhash","title":"Album Test","description":null,"images":[]},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	body, err := svc.QueryAlbum("testhash")
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestGetAlbumWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":"abc123","title":"My Album","description":null,"images":[{"id":"img1","link":"https://i.imgur.com/1.png","title":null,"name":null,"description":""}]},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	got, err := svc.GetAlbum("abc123")
	if err != nil {
		t.Fatalf("unexpected get album error: %v", err)
	}

	if got.ID != "abc123" {
		t.Fatalf("expected album id abc123, got %s", got.ID)
	}

	if len(got.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(got.Images))
	}
}

func TestCreateAlbumWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/3/album" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		bodyString := string(body)
		if !strings.Contains(bodyString, "title=New+Album") {
			t.Fatalf("expected title in body, got %s", bodyString)
		}
		if !strings.Contains(bodyString, "ids=img1%2Cimg2") {
			t.Fatalf("expected ids in body, got %s", bodyString)
		}

		_, _ = w.Write([]byte(`{"data":"new-album-id","status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	id, err := svc.CreateAlbum(album.CreateAlbumRequest{Title: "New Album", IDs: []string{"img1", "img2"}})
	if err != nil {
		t.Fatalf("unexpected create album error: %v", err)
	}

	if id == nil || *id != "new-album-id" {
		t.Fatalf("expected new-album-id, got %v", id)
	}
}

func TestUpdateAlbumWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/3/album/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	ok, err := svc.UpdateAlbum("abc123", album.UpdateAlbumRequest{Title: "Renamed"})
	if err != nil {
		t.Fatalf("unexpected update album error: %v", err)
	}

	if !ok {
		t.Fatal("expected update response to be true")
	}
}

func TestDeleteAlbumWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/3/album/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	ok, err := svc.DeleteAlbum("abc123")
	if err != nil {
		t.Fatalf("unexpected delete album error: %v", err)
	}

	if !ok {
		t.Fatal("expected delete response to be true")
	}
}

func TestFavoriteAndImageMembershipWithClient(t *testing.T) {
	var seenFavorite bool
	var seenAdd bool
	var seenRemove bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/3/album/abc123/favorite":
			seenFavorite = true
		case r.Method == http.MethodPost && r.URL.Path == "/3/album/abc123/add":
			seenAdd = true
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "ids=img1%2Cimg2") {
				t.Fatalf("missing ids in add body: %s", string(body))
			}
		case r.Method == http.MethodDelete && r.URL.Path == "/3/album/abc123/remove_images":
			seenRemove = true
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "ids=img3") {
				t.Fatalf("missing ids in remove body: %s", string(body))
			}
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}

		_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{
		BaseURL:  server.URL,
		ClientID: "test-client-id",
		Mode:     imgur.AuthModeAnonymous,
	})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := album.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}

	favorited, err := svc.FavoriteAlbum("abc123")
	if err != nil || !favorited {
		t.Fatalf("favorite failed: ok=%v err=%v", favorited, err)
	}

	added, err := svc.AddImagesToAlbum("abc123", []string{"img1", "img2"})
	if err != nil || !added {
		t.Fatalf("add images failed: ok=%v err=%v", added, err)
	}

	removed, err := svc.RemoveImagesFromAlbum("abc123", []string{"img3"})
	if err != nil || !removed {
		t.Fatalf("remove images failed: ok=%v err=%v", removed, err)
	}

	if !seenFavorite || !seenAdd || !seenRemove {
		t.Fatalf("expected all operations to be observed: favorite=%v add=%v remove=%v", seenFavorite, seenAdd, seenRemove)
	}
}
