package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestAssistant_GatingAndNotConfiguredState(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	status := s.Do(http.MethodPost, "/api/assistant/messages", token, map[string]string{"content": "Hola"}, nil)
	if status != http.StatusForbidden {
		t.Fatalf("send on free plan: expected 403, got %d", status)
	}

	s.Upgrade(token)

	// The test server never sets ANTHROPIC_API_KEY, so a pro user still gets
	// a clean "not configured" response instead of a crash or a fake reply.
	var body struct {
		Error string `json:"error"`
	}
	status = s.Do(http.MethodPost, "/api/assistant/messages", token, map[string]string{"content": "Hola"}, &body)
	if status != http.StatusServiceUnavailable || body.Error != "assistant_not_configured" {
		t.Fatalf("send with no API key: status=%d body=%+v", status, body)
	}

	// The user's message is still saved even though the reply failed, so a
	// conversation history exists once the assistant is later configured.
	var messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	status = s.Do(http.MethodGet, "/api/assistant/messages", token, nil, &messages)
	if status != http.StatusOK || len(messages) != 1 || messages[0].Role != "user" || messages[0].Content != "Hola" {
		t.Fatalf("list messages: status=%d body=%+v", status, messages)
	}

	status = s.Do(http.MethodDelete, "/api/assistant/messages", token, nil, nil)
	if status != http.StatusNoContent {
		t.Fatalf("clear: expected 204, got %d", status)
	}

	messages = nil
	s.Do(http.MethodGet, "/api/assistant/messages", token, nil, &messages)
	if len(messages) != 0 {
		t.Fatalf("expected empty history after clear, got %+v", messages)
	}
}
