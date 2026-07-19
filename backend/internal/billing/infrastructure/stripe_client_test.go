package infrastructure

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"rimu/backend/internal/billing/domain"
)

func TestStripeClient_NotConfiguredWithoutCredentials(t *testing.T) {
	c := NewStripeClient("", "", "")
	if c.Configured() {
		t.Fatal("expected Configured() to be false with no credentials")
	}
	_, err := c.CreateCheckoutSession(context.Background(), "user-1", "a@b.com", "https://x/success", "https://x/cancel")
	if err != domain.ErrNotConfigured {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestStripeClient_CreateCheckoutSession(t *testing.T) {
	var gotAuthUser string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()
		gotAuthUser = user
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.FormValue("client_reference_id") != "user-1" {
			t.Fatalf("expected client_reference_id=user-1, got %q", r.FormValue("client_reference_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"url":"https://checkout.stripe.com/c/session_123"}`))
	}))
	defer server.Close()

	c := &StripeClient{SecretKey: "sk_test_123", PriceID: "price_123", BaseURL: server.URL, Client: server.Client()}
	url, err := c.CreateCheckoutSession(context.Background(), "user-1", "a@b.com", "https://x/success", "https://x/cancel")
	if err != nil {
		t.Fatalf("create checkout session: %v", err)
	}
	if url != "https://checkout.stripe.com/c/session_123" {
		t.Fatalf("unexpected url: %q", url)
	}
	if gotAuthUser != "sk_test_123" {
		t.Fatalf("expected secret key as basic auth user, got %q", gotAuthUser)
	}
}

func signPayload(secret string, payload []byte, ts int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(payload)
	return fmt.Sprintf("t=%d,v1=%s", ts, hex.EncodeToString(mac.Sum(nil)))
}

func TestStripeClient_VerifyAndParseWebhook_Success(t *testing.T) {
	c := &StripeClient{WebhookSecret: "whsec_test"}
	payload := []byte(`{"type":"checkout.session.completed","data":{"object":{"client_reference_id":"user-42"}}}`)
	sig := signPayload("whsec_test", payload, time.Now().Unix())

	event, err := c.VerifyAndParseWebhook(payload, sig)
	if err != nil {
		t.Fatalf("verify webhook: %v", err)
	}
	if event.Type != "checkout.session.completed" || event.ClientReferenceID != "user-42" {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestStripeClient_VerifyAndParseWebhook_RejectsBadSignature(t *testing.T) {
	c := &StripeClient{WebhookSecret: "whsec_test"}
	payload := []byte(`{"type":"checkout.session.completed"}`)
	sig := signPayload("wrong-secret", payload, time.Now().Unix())

	_, err := c.VerifyAndParseWebhook(payload, sig)
	if err != domain.ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestStripeClient_VerifyAndParseWebhook_RejectsStaleTimestamp(t *testing.T) {
	c := &StripeClient{WebhookSecret: "whsec_test"}
	payload := []byte(`{"type":"checkout.session.completed"}`)
	sig := signPayload("whsec_test", payload, time.Now().Add(-1*time.Hour).Unix())

	_, err := c.VerifyAndParseWebhook(payload, sig)
	if err != domain.ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for a stale timestamp, got %v", err)
	}
}

func TestStripeClient_VerifyAndParseWebhook_RejectsTamperedPayload(t *testing.T) {
	c := &StripeClient{WebhookSecret: "whsec_test"}
	original := []byte(`{"type":"checkout.session.completed","data":{"object":{"client_reference_id":"user-42"}}}`)
	sig := signPayload("whsec_test", original, time.Now().Unix())

	tampered := []byte(`{"type":"checkout.session.completed","data":{"object":{"client_reference_id":"attacker"}}}`)
	_, err := c.VerifyAndParseWebhook(tampered, sig)
	if err != domain.ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for a tampered payload, got %v", err)
	}
}
