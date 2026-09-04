package images_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/derekpedersen/imgur-go/authorization"
	"github.com/derekpedersen/imgur-go/images"
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
			svc, err := images.NewService(tc.client)
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

func TestGetImage(t *testing.T) {
	if !hasIntegrationEnv() {
		t.Skip("skipping integration test: IMGUR_* env vars are not set")
	}

	imageID := "2nCt3Sbl"
	client := newOAuthClient(t)
	svc, err := images.NewService(client)
	if err != nil {
		t.Fatalf("failed to initialize image service: %v", err)
	}

	image, err := svc.GetImage(imageID)
	if err != nil {
		t.Errorf("experienced an error: %v", err)
	}
	if image == nil || len(image.ID) == 0 {
		t.Errorf("no image returned")
	}
}

func TestQueryImageWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/3/image/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Client-ID test-client-id" {
			t.Fatalf("unexpected auth header: %s", got)
		}
		_, _ = w.Write([]byte(`{"data":{"id":"abc123","link":"https://i.imgur.com/abc123.png","animated":false},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := images.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	body, err := svc.QueryImage("abc123")
	if err != nil {
		t.Fatalf("unexpected query image error: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("expected non-empty body")
	}
}

func TestGetImageWithClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"id":"abc123","link":"https://i.imgur.com/abc123.png","animated":false},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := images.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	got, err := svc.GetImage("abc123")
	if err != nil {
		t.Fatalf("unexpected get image error: %v", err)
	}
	if got.ID != "abc123" {
		t.Fatalf("expected image id abc123, got %s", got.ID)
	}
}

func TestImageMutationsWithClient(t *testing.T) {
	var seenUpload, seenUpdate, seenDelete, seenFavorite bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/3/image":
			seenUpload = true
			body, _ := io.ReadAll(r.Body)
			bodyString := string(body)
			if !strings.Contains(bodyString, "title=Test+Image") {
				t.Fatalf("missing title in upload payload: %s", bodyString)
			}
			_, _ = w.Write([]byte(`{"data":"img-new-1","status":200,"success":true}`))
		case r.Method == http.MethodPost && r.URL.Path == "/3/image/abc123":
			seenUpdate = true
			_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/3/image/abc123":
			seenDelete = true
			_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
		case r.Method == http.MethodPost && r.URL.Path == "/3/image/abc123/favorite":
			seenFavorite = true
			_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := images.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}

	uploadID, err := svc.UploadImage(images.UploadImageRequest{Image: "base64-data", Type: "base64", Title: "Test Image"})
	if err != nil {
		t.Fatalf("unexpected upload error: %v", err)
	}
	if uploadID == nil || *uploadID != "img-new-1" {
		t.Fatalf("expected upload id img-new-1, got %v", uploadID)
	}

	updated, err := svc.UpdateImage("abc123", images.UpdateImageRequest{Title: "Updated"})
	if err != nil || !updated {
		t.Fatalf("update failed: ok=%v err=%v", updated, err)
	}

	deleted, err := svc.DeleteImage("abc123")
	if err != nil || !deleted {
		t.Fatalf("delete failed: ok=%v err=%v", deleted, err)
	}

	favorited, err := svc.FavoriteImage("abc123")
	if err != nil || !favorited {
		t.Fatalf("favorite failed: ok=%v err=%v", favorited, err)
	}

	if !seenUpload || !seenUpdate || !seenDelete || !seenFavorite {
		t.Fatalf("expected all mutations observed upload=%v update=%v delete=%v favorite=%v", seenUpload, seenUpdate, seenDelete, seenFavorite)
	}
}
