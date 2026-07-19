package http_test

import (
	"net/http"
	"testing"

	"rimu/backend/internal/platform/testutil"
)

func TestBilling_NotConfiguredByDefault(t *testing.T) {
	s := testutil.NewServer(t)
	token := s.RegisterAndLogin("alice@example.com", "password123")

	var cfg struct {
		Configured bool `json:"configured"`
	}
	status := s.Do(http.MethodGet, "/api/billing/config", token, nil, &cfg)
	if status != http.StatusOK || cfg.Configured {
		t.Fatalf("expected configured=false with no Stripe env set: status=%d body=%+v", status, cfg)
	}

	var checkoutBody struct {
		Error string `json:"error"`
	}
	status = s.Do(http.MethodPost, "/api/billing/checkout", token, nil, &checkoutBody)
	if status != http.StatusServiceUnavailable || checkoutBody.Error != "billing_not_configured" {
		t.Fatalf("checkout with no Stripe env: status=%d body=%+v", status, checkoutBody)
	}

	// The webhook route requires no auth token at all — Stripe calls it
	// directly — but without a webhook secret configured it must still
	// reject rather than silently trust the payload.
	status = s.Do(http.MethodPost, "/api/billing/webhook", "", map[string]string{"type": "checkout.session.completed"}, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("webhook with no signing secret: expected 400, got %d", status)
	}
}
