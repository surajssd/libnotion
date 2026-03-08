package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/surajssd/libnotion/api"
)

func TestQueryDatabase_Success_SinglePage(t *testing.T) {
	expectedPages := []api.Page{
		{CommonObject: api.CommonObject{ID: "page-1", Object: "page"}},
		{CommonObject: api.CommonObject{ID: "page-2", Object: "page"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/data_sources/db-123/query" {
			t.Errorf("expected path /v1/data_sources/db-123/query, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.PageResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results:  expectedPages,
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	query := &api.QueryDB{
		Sorts: []api.Sort{{Property: "Name", Direction: &api.SortDirectionAscending}},
	}

	pages, err := client.QueryDatabase("db-123", query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
	if pages[0].ID != "page-1" {
		t.Errorf("expected page ID %q, got %q", "page-1", pages[0].ID)
	}
	if pages[1].ID != "page-2" {
		t.Errorf("expected page ID %q, got %q", "page-2", pages[1].ID)
	}
}

func TestQueryDatabase_Success_Pagination(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)

		body, _ := io.ReadAll(r.Body)
		var query api.QueryDB
		json.Unmarshal(body, &query)

		w.Header().Set("Content-Type", "application/json")

		if count == 1 {
			// First page — verify no start_cursor
			if query.StartCursor != "" {
				t.Errorf("first request should have empty start_cursor, got %q", query.StartCursor)
			}
			json.NewEncoder(w).Encode(api.PageResponseList{
				Response: api.Response{Object: "list", HasMore: true, NextCursor: "cursor-abc"},
				Results:  []api.Page{{CommonObject: api.CommonObject{ID: "page-1"}}},
			})
		} else {
			// Second page — verify start_cursor
			if query.StartCursor != "cursor-abc" {
				t.Errorf("expected start_cursor %q, got %q", "cursor-abc", query.StartCursor)
			}
			json.NewEncoder(w).Encode(api.PageResponseList{
				Response: api.Response{Object: "list", HasMore: false},
				Results:  []api.Page{{CommonObject: api.CommonObject{ID: "page-2"}}},
			})
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	pages, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

func TestQueryDatabase_NilQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.PageResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results:  []api.Page{{CommonObject: api.CommonObject{ID: "page-1"}}},
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	pages, err := client.QueryDatabase("db-123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
}

func TestQueryDatabase_Non200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(api.FailureResponse{
			Object:  "error",
			Status:  403,
			Code:    "restricted_resource",
			Message: "insufficient permissions",
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") || !contains(got, "insufficient permissions") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestQueryDatabase_Non200Response_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestQueryDatabase_InvalidResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "unmarshal response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestQueryDatabase_RequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := newTestClient(server.URL)
	_, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "listing database entries") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestQueryDatabase_InvalidBaseURL(t *testing.T) {
	client := NewNotionClient(WithBaseURL("://invalid-url"))
	_, err := client.QueryDatabase("db-123", &api.QueryDB{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "parsing the APIURL") {
		t.Errorf("unexpected error message: %s", got)
	}
}
