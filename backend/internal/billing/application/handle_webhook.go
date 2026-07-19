package application

import (
	"context"

	"rimu/backend/internal/billing/domain"
)

// StripeWebhookVerifier is the boundary to Stripe's signature scheme —
// implemented by infrastructure/stripe_client.go.
type StripeWebhookVerifier interface {
	VerifyAndParseWebhook(payload []byte, sigHeader string) (domain.WebhookEvent, error)
}

// HandleWebhook verifies a raw Stripe webhook request and, for a completed
// checkout, upgrades the referenced user to Pro. Events this app doesn't
// care about (anything but checkout.session.completed) are verified but
// otherwise ignored.
type HandleWebhook struct {
	Verifier StripeWebhookVerifier
	Upgrade  *ApplyProUpgrade
}

func (uc *HandleWebhook) Execute(ctx context.Context, payload []byte, sigHeader string) error {
	event, err := uc.Verifier.VerifyAndParseWebhook(payload, sigHeader)
	if err != nil {
		return err
	}
	if event.Type != "checkout.session.completed" || event.ClientReferenceID == "" {
		return nil
	}
	return uc.Upgrade.Execute(ctx, event.ClientReferenceID)
}
