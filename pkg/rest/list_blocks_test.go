package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/surajssd/libnotion/api"
	"github.com/surajssd/libnotion/api/blocks"
)

func TestListBlocks_Success_SinglePage(t *testing.T) {
	bt := blocks.BTParagraph
	expectedBlocks := []blocks.Block{
		{
			CommonObject: api.CommonObject{ID: "block-1", Object: "block"},
			Type:         &bt,
			Paragraph:    &blocks.Property{Text: []blocks.FullText{{PlainText: "Hello"}}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/blocks/block-parent-123/children" {
			t.Errorf("expected path /v1/blocks/block-parent-123/children, got %s", r.URL.Path)
		}
		if r.Header.Get("Notion-Version") != NotionVersion {
			t.Errorf("expected Notion-Version %s, got %s", NotionVersion, r.Header.Get("Notion-Version"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blocks.BlockResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results:  expectedBlocks,
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.ListBlocks("block-parent-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 block, got %d", len(result))
	}
	if result[0].ID != "block-1" {
		t.Errorf("expected block ID %q, got %q", "block-1", result[0].ID)
	}
}

func TestListBlocks_Success_Pagination(t *testing.T) {
	var callCount int32

	bt := blocks.BTParagraph
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)

		w.Header().Set("Content-Type", "application/json")

		if count == 1 {
			// First page — no start_cursor query param
			if r.URL.Query().Get("start_cursor") != "" {
				t.Errorf("first request should have no start_cursor")
			}
			json.NewEncoder(w).Encode(blocks.BlockResponseList{
				Response: api.Response{Object: "list", HasMore: true, NextCursor: "cursor-blocks"},
				Results: []blocks.Block{
					{CommonObject: api.CommonObject{ID: "block-1"}, Type: &bt},
				},
			})
		} else {
			// Second page — verify start_cursor and page_size
			if got := r.URL.Query().Get("start_cursor"); got != "cursor-blocks" {
				t.Errorf("expected start_cursor %q, got %q", "cursor-blocks", got)
			}
			if got := r.URL.Query().Get("page_size"); got != "100" {
				t.Errorf("expected page_size %q, got %q", "100", got)
			}
			json.NewEncoder(w).Encode(blocks.BlockResponseList{
				Response: api.Response{Object: "list", HasMore: false},
				Results: []blocks.Block{
					{CommonObject: api.CommonObject{ID: "block-2"}, Type: &bt},
				},
			})
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.ListBlocks("parent-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(result))
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

func TestListBlocks_Non200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(api.FailureResponse{
			Object:  "error",
			Status:  404,
			Code:    "object_not_found",
			Message: "block not found",
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.ListBlocks("nonexistent-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") || !contains(got, "block not found") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestListBlocks_Non200Response_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.ListBlocks("some-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestListBlocks_InvalidResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.ListBlocks("some-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "unmarshal response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestListBlocks_RequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := newTestClient(server.URL)
	_, err := client.ListBlocks("some-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "listing block entries") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestListBlocks_InvalidBaseURL(t *testing.T) {
	client := NewNotionClient(WithBaseURL("://invalid-url"))
	_, err := client.ListBlocks("some-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "parsing the APIURL") {
		t.Errorf("unexpected error message: %s", got)
	}
}
