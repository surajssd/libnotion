package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/surajssd/libnotion/api"
)

func TestFindDatabase_Success_FirstPage(t *testing.T) {
	expectedDB := api.Database{
		CommonObject: api.CommonObject{ID: "db-123", Object: "database"},
		Title: []api.Title{
			{Text: api.Text{Content: "My Database"}, PlainText: "My Database"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/databases" {
			t.Errorf("expected path /v1/databases, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.DatabaseResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results:  []api.Database{expectedDB},
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	db, err := client.FindDatabase("My Database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.ID != "db-123" {
		t.Errorf("expected database ID %q, got %q", "db-123", db.ID)
	}
}

func TestFindDatabase_Success_WithPagination(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)

		w.Header().Set("Content-Type", "application/json")

		if count == 1 {
			// Verify no start_cursor on first request
			if r.URL.Query().Get("start_cursor") != "" {
				t.Errorf("first request should have no start_cursor")
			}
			json.NewEncoder(w).Encode(api.DatabaseResponseList{
				Response: api.Response{Object: "list", HasMore: true, NextCursor: "cursor-xyz"},
				Results: []api.Database{
					{
						CommonObject: api.CommonObject{ID: "db-other"},
						Title:        []api.Title{{Text: api.Text{Content: "Other DB"}}},
					},
				},
			})
		} else {
			// Verify start_cursor on second request
			if got := r.URL.Query().Get("start_cursor"); got != "cursor-xyz" {
				t.Errorf("expected start_cursor %q, got %q", "cursor-xyz", got)
			}
			json.NewEncoder(w).Encode(api.DatabaseResponseList{
				Response: api.Response{Object: "list", HasMore: false},
				Results: []api.Database{
					{
						CommonObject: api.CommonObject{ID: "db-target"},
						Title:        []api.Title{{Text: api.Text{Content: "Target DB"}}},
					},
				},
			})
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	db, err := client.FindDatabase("Target DB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.ID != "db-target" {
		t.Errorf("expected database ID %q, got %q", "db-target", db.ID)
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

func TestFindDatabase_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.DatabaseResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results: []api.Database{
				{
					CommonObject: api.CommonObject{ID: "db-other"},
					Title:        []api.Title{{Text: api.Text{Content: "Other DB"}}},
				},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.FindDatabase("Nonexistent DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, `database "Nonexistent DB" not found`) {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestFindDatabase_SkipsEmptyTitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.DatabaseResponseList{
			Response: api.Response{Object: "list", HasMore: false},
			Results: []api.Database{
				// Database with empty Title slice
				{
					CommonObject: api.CommonObject{ID: "db-empty-title"},
				},
				// Database with empty Content
				{
					CommonObject: api.CommonObject{ID: "db-empty-content"},
					Title:        []api.Title{{Text: api.Text{Content: ""}}},
				},
				// Target database
				{
					CommonObject: api.CommonObject{ID: "db-target"},
					Title:        []api.Title{{Text: api.Text{Content: "Target DB"}}},
				},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	db, err := client.FindDatabase("Target DB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.ID != "db-target" {
		t.Errorf("expected database ID %q, got %q", "db-target", db.ID)
	}
}

func TestFindDatabase_Non200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.FailureResponse{
			Object:  "error",
			Status:  401,
			Code:    "unauthorized",
			Message: "API token is invalid",
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.FindDatabase("My DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") || !contains(got, "API token is invalid") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestFindDatabase_Non200Response_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.FindDatabase("My DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestFindDatabase_InvalidResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.FindDatabase("My DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "unmarshal response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestFindDatabase_RequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := newTestClient(server.URL)
	_, err := client.FindDatabase("My DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "listing databases") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestFindDatabase_InvalidBaseURL(t *testing.T) {
	client := NewNotionClient(WithBaseURL("://invalid-url"))
	_, err := client.FindDatabase("My DB")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "parsing the APIURL") {
		t.Errorf("unexpected error message: %s", got)
	}
}
