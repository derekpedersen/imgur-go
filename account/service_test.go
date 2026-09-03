package account_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekpedersen/imgur-go/account"
	"github.com/derekpedersen/imgur-go/imgur"
)

func TestGetAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/3/account/alice" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"id":12,"url":"alice","display_name":"Alice","avatar":"https://i.imgur.com/a.png"},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc := account.NewService(client)
	profile, err := svc.GetAccount("alice")
	if err != nil {
		t.Fatalf("unexpected get account error: %v", err)
	}
	if profile.URL != "alice" {
		t.Fatalf("expected account url alice, got %s", profile.URL)
	}
}

func TestGetFavorites(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/3/account/alice/favorites/2/top" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"fav1","title":"Favorite 1","link":"https://imgur.com/fav1","is_album":true}],"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc := account.NewService(client)
	items, err := svc.GetFavorites("alice", 2, "top")
	if err != nil {
		t.Fatalf("unexpected get favorites error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "fav1" {
		t.Fatalf("unexpected favorites: %+v", items)
	}
}

func TestGetSubmissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/3/account/alice/submissions/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"sub1","title":"Submission 1","link":"https://imgur.com/sub1","is_album":false}],"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc := account.NewService(client)
	items, err := svc.GetSubmissions("alice", 1)
	if err != nil {
		t.Fatalf("unexpected get submissions error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "sub1" {
		t.Fatalf("unexpected submissions: %+v", items)
	}
}

func TestGetSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/3/account/alice/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"account_url":"alice","email":"alice@example.com","public_images":true,"album_privacy":"hidden","pro_expiration":false,"accepted_gallery_terms":true},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc := account.NewService(client)
	settings, err := svc.GetSettings("alice")
	if err != nil {
		t.Fatalf("unexpected get settings error: %v", err)
	}
	if settings.AccountURL != "alice" {
		t.Fatalf("unexpected settings: %+v", settings)
	}
}
