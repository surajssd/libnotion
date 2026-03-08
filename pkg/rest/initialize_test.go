package rest

import (
	"testing"
)

func TestNewNotionClient_NoOptions(t *testing.T) {
	client := NewNotionClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.token != "" {
		t.Errorf("expected empty token, got %q", client.token)
	}
	if client.baseURL != "" {
		t.Errorf("expected empty baseURL, got %q", client.baseURL)
	}
}

func TestNewNotionClient_WithSecretToken(t *testing.T) {
	client := NewNotionClient(WithSecretToken("test-token"))
	if client.token != "test-token" {
		t.Errorf("expected token %q, got %q", "test-token", client.token)
	}
}

func TestNewNotionClient_WithBaseURL(t *testing.T) {
	client := NewNotionClient(WithBaseURL("http://localhost:8080"))
	if client.baseURL != "http://localhost:8080" {
		t.Errorf("expected baseURL %q, got %q", "http://localhost:8080", client.baseURL)
	}
}

func TestNewNotionClient_MultipleOptions(t *testing.T) {
	client := NewNotionClient(
		WithSecretToken("my-token"),
		WithBaseURL("http://localhost:9090"),
	)
	if client.token != "my-token" {
		t.Errorf("expected token %q, got %q", "my-token", client.token)
	}
	if client.baseURL != "http://localhost:9090" {
		t.Errorf("expected baseURL %q, got %q", "http://localhost:9090", client.baseURL)
	}
}

func TestGetBaseURL_Default(t *testing.T) {
	client := NewNotionClient()
	if got := client.getBaseURL(); got != APIURL {
		t.Errorf("expected %q, got %q", APIURL, got)
	}
}

func TestGetBaseURL_Custom(t *testing.T) {
	customURL := "http://localhost:8080"
	client := NewNotionClient(WithBaseURL(customURL))
	if got := client.getBaseURL(); got != customURL {
		t.Errorf("expected %q, got %q", customURL, got)
	}
}
