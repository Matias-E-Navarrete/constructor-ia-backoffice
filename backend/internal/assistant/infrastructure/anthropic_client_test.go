package infrastructure

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"rimu/backend/internal/assistant/domain"
)

func TestAnthropicClient_NoAPIKeyReturnsNotConfigured(t *testing.T) {
	c := NewAnthropicClient("")
	_, err := c.Complete(context.Background(), []domain.Message{{Role: domain.RoleUser, Content: "hola"}})
	if err != domain.ErrNotConfigured {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestAnthropicClient_ParsesSuccessfulReply(t *testing.T) {
	var gotAuth, gotVersion string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("x-api-key")
		gotVersion = r.Header.Get("anthropic-version")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"content":[{"type":"text","text":"Hola! ¿En qué te ayudo hoy?"}]}`))
	}))
	defer server.Close()

	c := &AnthropicClient{APIKey: "test-key", Model: "claude-sonnet-5", BaseURL: server.URL, Client: server.Client()}
	reply, err := c.Complete(context.Background(), []domain.Message{{Role: domain.RoleUser, Content: "hola"}})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	if reply != "Hola! ¿En qué te ayudo hoy?" {
		t.Fatalf("unexpected reply: %q", reply)
	}
	if gotAuth != "test-key" {
		t.Fatalf("expected x-api-key header to be set, got %q", gotAuth)
	}
	if gotVersion == "" {
		t.Fatalf("expected anthropic-version header to be set")
	}
}

func TestAnthropicClient_SurfacesProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"invalid x-api-key"}}`))
	}))
	defer server.Close()

	c := &AnthropicClient{APIKey: "bad-key", Model: "claude-sonnet-5", BaseURL: server.URL, Client: server.Client()}
	_, err := c.Complete(context.Background(), []domain.Message{{Role: domain.RoleUser, Content: "hola"}})
	if err == nil {
		t.Fatal("expected an error for a 401 response")
	}
}
