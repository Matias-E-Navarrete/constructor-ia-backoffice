package infrastructure

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"rimu/backend/internal/billing/domain"
)

// StripeClient is a minimal wrapper around the parts of the Stripe REST API
// this app needs (creating a Checkout Session, verifying a webhook
// signature) — no SDK dependency, same "plain net/http" convention as the
// Anthropic client.
type StripeClient struct {
	SecretKey     string
	PriceID       string
	WebhookSecret string
	BaseURL       string
	Client        *http.Client
}

func NewStripeClient(secretKey, priceID, webhookSecret string) *StripeClient {
	return &StripeClient{
		SecretKey:     secretKey,
		PriceID:       priceID,
		WebhookSecret: webhookSecret,
		BaseURL:       "https://api.stripe.com/v1",
		Client:        &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *StripeClient) Configured() bool {
	return c.SecretKey != "" && c.PriceID != ""
}

type checkoutSessionResponse struct {
	URL   string `json:"url"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// CreateCheckoutSession opens a subscription Checkout Session for the Pro
// plan, tagging it with clientReferenceID (our user ID) so the webhook can
// tell which account to upgrade once payment completes.
func (c *StripeClient) CreateCheckoutSession(ctx context.Context, clientReferenceID, customerEmail, successURL, cancelURL string) (string, error) {
	if !c.Configured() {
		return "", domain.ErrNotConfigured
	}

	form := url.Values{}
	form.Set("mode", "subscription")
	form.Set("line_items[0][price]", c.PriceID)
	form.Set("line_items[0][quantity]", "1")
	form.Set("client_reference_id", clientReferenceID)
	form.Set("customer_email", customerEmail)
	form.Set("success_url", successURL)
	form.Set("cancel_url", cancelURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/checkout/sessions", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.SecretKey, "")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed checkoutSessionResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("billing: decoding checkout session response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil {
			return "", fmt.Errorf("billing: stripe error: %s", parsed.Error.Message)
		}
		return "", errors.New("billing: stripe returned a non-200 status")
	}
	if parsed.URL == "" {
		return "", errors.New("billing: stripe returned no checkout url")
	}
	return parsed.URL, nil
}

// VerifyAndParseWebhook checks the Stripe-Signature header against the raw
// request body (per Stripe's documented HMAC scheme) and, only if it
// matches, parses out the event type and the checkout session's
// client_reference_id.
func (c *StripeClient) VerifyAndParseWebhook(payload []byte, sigHeader string) (domain.WebhookEvent, error) {
	if c.WebhookSecret == "" {
		return domain.WebhookEvent{}, domain.ErrNotConfigured
	}
	if err := verifyStripeSignature(payload, sigHeader, c.WebhookSecret, 5*time.Minute); err != nil {
		return domain.WebhookEvent{}, err
	}

	var parsed struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ClientReferenceID string `json:"client_reference_id"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return domain.WebhookEvent{}, fmt.Errorf("billing: decoding webhook payload: %w", err)
	}
	return domain.WebhookEvent{Type: parsed.Type, ClientReferenceID: parsed.Data.Object.ClientReferenceID}, nil
}

func verifyStripeSignature(payload []byte, sigHeader, secret string, tolerance time.Duration) error {
	var timestamp string
	var signatures []string
	for _, part := range strings.Split(sigHeader, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signatures = append(signatures, kv[1])
		}
	}
	if timestamp == "" || len(signatures) == 0 {
		return domain.ErrInvalidSignature
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return domain.ErrInvalidSignature
	}
	if age := time.Since(time.Unix(ts, 0)); age > tolerance || age < -tolerance {
		return domain.ErrInvalidSignature
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if hmac.Equal([]byte(sig), []byte(expected)) {
			return nil
		}
	}
	return domain.ErrInvalidSignature
}
