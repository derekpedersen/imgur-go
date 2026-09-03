package gallery_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/derekpedersen/imgur-go/gallery"
	"github.com/derekpedersen/imgur-go/imgur"
)

func TestGetGallery(t *testing.T) {
	showViral := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/3/gallery/top/top/week/2" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("showViral") != "false" {
			t.Fatalf("unexpected showViral value: %s", r.URL.Query().Get("showViral"))
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"item1","title":"First","is_album":false,"link":"https://i.imgur.com/item1.png"}],"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := gallery.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	items, err := svc.GetGallery(gallery.ListOptions{Section: "top", Sort: "top", Window: "week", Page: 2, ShowViral: &showViral})
	if err != nil {
		t.Fatalf("unexpected get gallery error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "item1" {
		t.Fatalf("unexpected gallery items: %+v", items)
	}
}

func TestSearchGallery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/3/gallery/search/time/all/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "cats" {
			t.Fatalf("unexpected query q: %s", r.URL.Query().Get("q"))
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"search1","title":"Cats","is_album":true,"link":"https://imgur.com/a/search1"}],"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := gallery.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	items, err := svc.SearchGallery("cats", gallery.ListOptions{Sort: "time", Window: "all", Page: 1})
	if err != nil {
		t.Fatalf("unexpected search gallery error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "search1" {
		t.Fatalf("unexpected search results: %+v", items)
	}
}

func TestGetGalleryAlbum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/3/gallery/album/abc123" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"id":"abc123","title":"Album","images":[{"id":"img1","link":"https://i.imgur.com/img1.png"}]},"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := gallery.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	albumItem, err := svc.GetGalleryAlbum("abc123")
	if err != nil {
		t.Fatalf("unexpected get gallery album error: %v", err)
	}
	if albumItem.ID != "abc123" || len(albumItem.Images) != 1 {
		t.Fatalf("unexpected album detail: %+v", albumItem)
	}
}

func TestVoteGalleryItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/3/gallery/abc123/vote/up" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/x-www-form-urlencoded") {
			t.Fatalf("unexpected content type: %s", ct)
		}
		_, _ = w.Write([]byte(`{"data":true,"status":200,"success":true}`))
	}))
	defer server.Close()

	client, err := imgur.NewClient(imgur.Config{BaseURL: server.URL, ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := gallery.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	ok, err := svc.VoteGalleryItem("abc123", "up")
	if err != nil {
		t.Fatalf("unexpected vote error: %v", err)
	}
	if !ok {
		t.Fatal("expected vote response true")
	}
}

func TestVoteGalleryItemRejectsInvalidVote(t *testing.T) {
	client, err := imgur.NewClient(imgur.Config{BaseURL: "https://api.imgur.com", ClientID: "test-client-id", Mode: imgur.AuthModeAnonymous})
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	svc, err := gallery.NewService(client)
	if err != nil {
		t.Fatalf("unexpected service error: %v", err)
	}
	_, err = svc.VoteGalleryItem("abc123", "banana")
	if err == nil {
		t.Fatal("expected invalid vote error")
	}
}
