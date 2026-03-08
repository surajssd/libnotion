package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/surajssd/libnotion/api"
)

func newTestClient(serverURL string) *NotionClient {
	return NewNotionClient(
		WithSecretToken("test-token"),
		WithBaseURL(serverURL),
	)
}

func TestAddPage_Success(t *testing.T) {
	expectedPage := api.Page{
		CommonObject: api.CommonObject{
			ID:     "page-123",
			Object: "page",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/pages" {
			t.Errorf("expected path /v1/pages, got %s", r.URL.Path)
		}

		// Verify headers
		if r.Header.Get("Notion-Version") != NotionVersion {
			t.Errorf("expected Notion-Version %s, got %s", NotionVersion, r.Header.Get("Notion-Version"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify request body is valid JSON
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body: %v", err)
		}
		var pg api.Page
		if err := json.Unmarshal(body, &pg); err != nil {
			t.Fatalf("unmarshalling request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedPage)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	inputPage := api.Page{
		Parent: api.Parent{
			Type:       api.ParentTypeDatabase,
			DatabaseID: "db-123",
		},
	}

	result, err := client.AddPage(inputPage)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "page-123" {
		t.Errorf("expected page ID %q, got %q", "page-123", result.ID)
	}
	if result.Object != "page" {
		t.Errorf("expected object %q, got %q", "page", result.Object)
	}
}

func TestAddPage_Non200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.FailureResponse{
			Object:  "error",
			Status:  400,
			Code:    "validation_error",
			Message: "invalid page properties",
		})
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.AddPage(api.Page{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") || !contains(got, "invalid page properties") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestAddPage_Non200Response_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.AddPage(api.Page{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "non-200 response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestAddPage_InvalidResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.AddPage(api.Page{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "unmarshal response") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestAddPage_RequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close() // close immediately to cause connection error

	client := newTestClient(server.URL)
	_, err := client.AddPage(api.Page{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "adding a new page") {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestAddPage_InvalidBaseURL(t *testing.T) {
	client := NewNotionClient(WithBaseURL("://invalid-url"))
	_, err := client.AddPage(api.Page{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); !contains(got, "parsing the APIURL") {
		t.Errorf("unexpected error message: %s", got)
	}
}

// contains is a helper to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
